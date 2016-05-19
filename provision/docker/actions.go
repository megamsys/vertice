package docker

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	lb "github.com/megamsys/vertice/logbox"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/provision/docker/container"
	"github.com/megamsys/vertice/router"
)

type runContainerActionsArgs struct {
	box              *provision.Box
	imageId          string
	containerStatus  utils.Status
	destinationHosts []string
	writer           io.Writer
	isDeploy         bool
	buildingImage    string
	provisioner      *dockerProvisioner
}

type containersToAdd struct {
	Quantity int
	Status   utils.Status
}

type changeUnitsPipelineArgs struct {
	box         *provision.Box
	writer      io.Writer
	toAdd       map[string]*containersToAdd
	toRemove    []container.Container
	toHost      string
	imageId     string
	provisioner *dockerProvisioner
	boxDestroy  bool
}

type callbackFunc func(*container.Container, chan *container.Container) error

type rollbackFunc func(*container.Container)

func runInContainers(containers []container.Container, callback callbackFunc, rollback rollbackFunc, parallel bool) error {

	if len(containers) == 0 {
		return nil
	}
	//workers, _ := config.GetInt("docker:max-workers")
	workers := 0
	if workers == 0 {
		workers = len(containers)
	}
	step := len(containers)/workers + 1
	toRollback := make(chan *container.Container, len(containers))
	errors := make(chan error, len(containers))
	var wg sync.WaitGroup
	runFunc := func(start, end int) error {
		defer wg.Done()
		for i := start; i < end; i++ {
			err := callback(&containers[i], toRollback)
			if err != nil {
				errors <- err
				return err
			}
		}
		return nil
	}
	for i := 0; i < len(containers); i += step {
		end := i + step
		if end > len(containers) {
			end = len(containers)
		}
		wg.Add(1)
		if parallel {
			go runFunc(i, end)
		} else {
			err := runFunc(i, end)
			if err != nil {
				break
			}
		}
	}
	wg.Wait()
	close(errors)
	close(toRollback)
	if err := <-errors; err != nil {
		if rollback != nil {
			for c := range toRollback {
				rollback(c)
			}
		}
		return err
	}

	return nil
}

var updateStatusInScylla = action.Action{
	Name: "update-status-scylla",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runContainerActionsArgs)
		var cont container.Container
		if ctx.Previous != nil {
			cont = ctx.Previous.(container.Container)
		} else {
			ncont, _ := args.provisioner.GetContainerByBox(args.box)
			cont = *ncont
			cont.Image = args.imageId
		}

		if err := cont.SetStatus(args.containerStatus); err != nil {
			return nil, err
		}
		return cont, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(container.Container)
		c.SetStatus(constants.StatusError)
	},
}


var createContainer = action.Action{
	Name: "create-container",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		cont := ctx.Previous.(container.Container)
		args := ctx.Params[0].(runContainerActionsArgs)
		log.Debugf("  create container for box (%s, image:%s)/%s", args.box.GetFullName(), args.imageId, args.box.Compute)

		err := cont.Create(&container.CreateArgs{
			ImageId:     args.imageId,
			Box:         args.box,
			Deploy:      args.isDeploy,
			Provisioner: args.provisioner,
		})

		if err != nil {
			return nil, err
		}
		cont.Status = constants.StatusLaunched
		return cont, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(container.Container)
		args := ctx.Params[0].(runContainerActionsArgs)

		fmt.Fprintf(args.writer, lb.W(lb.CONTAINER_DEPLOY, lb.INFO, fmt.Sprintf("\n---- Removing old container %s ----\n", c.Name)))
		err := args.provisioner.Cluster().RemoveContainer(docker.RemoveContainerOptions{ID: c.Id})
		if err != nil {
			log.Errorf("---- [start-container:Backward]\n     %s", err.Error())
		}
	},
}

var startContainer = action.Action{
	Name: "start-container",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		c := ctx.Previous.(container.Container)
		log.Debugf("  starting container (%s)", c.Id)
		args := ctx.Params[0].(runContainerActionsArgs)
		err := c.Start(&container.StartArgs{
			Provisioner: args.provisioner,
			Box:         args.box,
			Deploy:      args.isDeploy,
		})
		if err != nil {
			return nil, err
		}
		c.Status = constants.StatusStarted
		return c, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(container.Container)
		args := ctx.Params[0].(runContainerActionsArgs)

		fmt.Fprintf(args.writer, lb.W(lb.CONTAINER_DEPLOY, lb.INFO, fmt.Sprintf("\n---- Stopping old container %s ----", c.Id)))
		err := args.provisioner.Cluster().StopContainer(c.Id, 10)
		if err != nil {
			log.Errorf("---- [stop-container:Backward]\n     %s", err.Error())
		}
	},
}

var stopContainer = action.Action{
	Name: "stop-container",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runContainerActionsArgs)
		cont := ctx.Previous.(container.Container)
		log.Debugf("  stopping container (%s)", cont.Id)
		if err := cont.Stop(args.provisioner); err != nil {
			return nil, err
		}

		cont.Status = constants.StatusStopped
		return cont, nil
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(runContainerActionsArgs)
		c := ctx.FWResult.(container.Container)
		c.Status = constants.StatusStopping

		fmt.Fprintf(args.writer, lb.W(lb.CONTAINER_DEPLOY, lb.INFO, fmt.Sprintf("\n---- Skip Stopping old container %s ----", c.Id)))
	},
}

var destroyOldContainers = action.Action{
	Name: "destroy-old-containers",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(changeUnitsPipelineArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		total := len(args.toRemove)
		var plural string
		if total > 1 {
			plural = "s"
		}

		fmt.Fprintf(writer, lb.W(lb.CONTAINER_DEPLOY, lb.INFO, fmt.Sprintf("\n---- Destroying %d old containers%s ----", total, plural)))
		runInContainers(args.toRemove, func(c *container.Container, toRollback chan *container.Container) error {

			err := c.Remove(args.provisioner)
			if err != nil {
				log.Errorf("Ignored error trying to remove old container %q: %s", c.Id, err)
			}

			fmt.Fprintf(writer, lb.W(lb.CONTAINER_DEPLOY, lb.INFO, fmt.Sprintf(" ---> Destroyed old container (%s, %s)", c.BoxName, c.ShortId())))
			return nil
		}, nil, true)
		return ctx.Previous, nil
	},
	Backward: func(ctx action.BWContext) {
	},
	MinParams: 1,
}

var setNetworkInfo = action.Action{
	Name: "set-network-info",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		c := ctx.Previous.(container.Container)
		args := ctx.Params[0].(runContainerActionsArgs)
		info, err := c.NetworkInfo(args.provisioner)
		if err != nil {
			return nil, err
		}
		c.PublicIp = info.IP
		c.HostPort = info.HTTPHostPort
		return c, nil
	},
}

var followLogsAndCommit = action.Action{
	Name: "follow-logs-and-commit",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		c, ok := ctx.Previous.(container.Container)
		if !ok {
			return nil, errors.New("Previous result must be a container.")
		}
		args := ctx.Params[0].(runContainerActionsArgs)
		status, err := c.Logs(args.provisioner)
		if err != nil {
			log.Errorf("---- follow logs for container\n     %s", err.Error())
			return nil, err
		}
		if status != 0 {
			return nil, fmt.Errorf("Exit status %d", status)
		}
		/*fmt.Fprintf(args.writer, "\n---- Building application image ----\n")
		imageId, err := c.Commit(args.provisioner, args.writer)
		if err != nil {
			log.Errorf("error on commit container %s - %s", c.Id, err)
			return nil, err
		}
		fmt.Fprintf(args.writer, " ---> Cleaning up\n")
		c.Remove(args.provisioner)
		return imageId, nil
		*/
		return "", nil
	},
	Backward: func(ctx action.BWContext) {
	},
	MinParams: 1,
}

var addNewRoute = action.Action{
	Name: "add-new-route",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(changeUnitsPipelineArgs)
		c, ok := ctx.Previous.(container.Container)
		if !ok {
			return nil, errors.New("Previous result must be a container.")
		}
		r, err := getRouterForBox(args.box)
		if err != nil {
			return nil, err
		}
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		newContainers := []container.Container{c}
		if len(newContainers) > 0 {

			fmt.Fprintf(writer, lb.W(lb.CONTAINER_DEPLOY, lb.INFO, fmt.Sprintf("---- Adding routes to new containers ----")))
		}

		return newContainers, runInContainers(newContainers, func(c *container.Container, toRollback chan *container.Container) error {
			err = r.SetCName(c.BoxName, c.PublicIp)
			if err != nil {
				return err
			}
			c.Routable = true
			toRollback <- c
			
			fmt.Fprintf(writer, lb.W(lb.CONTAINER_DEPLOY, lb.INFO, fmt.Sprintf("---> Added route to container (%s,%s)", c.BoxName, c.ShortId())))
			return nil
		}, func(c *container.Container) {
			r.UnsetCName(c.BoxName, c.PublicIp)
		}, false)
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(changeUnitsPipelineArgs)
		c := ctx.FWResult.(container.Container)
		r, err := getRouterForBox(args.box)
		if err != nil {
			log.Errorf("---- [add-new-routes:Backward]\n     %s", err.Error())
		}

		newContainers := []container.Container{c}

		fmt.Fprintf(args.writer, lb.W(lb.CONTAINER_DEPLOY, lb.INFO, fmt.Sprintf("---- Destroying routes from created containers  (%s, %s)", c.BoxName, c.ShortId())))
		for _, cont := range newContainers {
			if !cont.Routable {
				continue
			}
			err = r.UnsetCName(cont.BoxName, cont.PublicIp)
			if err != nil {
				log.Errorf("---- [add-new-routes:Backward] (%s, %s)\n    %s", c.BoxName, args.box.PublicIp, err.Error())
			}

			fmt.Fprintf(args.writer, lb.W(lb.CONTAINER_DEPLOY, lb.INFO, fmt.Sprintf("---- Destroyed route from machine (%s, %s)", c.BoxName, c.ShortId())))
		}
	},
}

var removeOldRoutes = action.Action{
	Name: "remove-old-routes",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(changeUnitsPipelineArgs)
		r, err := getRouterForBox(args.box)
		if err != nil {
			return nil, err
		}
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		if len(args.toRemove) > 0 {

			fmt.Fprintf(writer, lb.W(lb.CONTAINER_DEPLOY, lb.INFO, fmt.Sprintf("---- Removing routes from old containers ----")))

		}
		return ctx.Previous, runInContainers(args.toRemove, func(c *container.Container, toRollback chan *container.Container) error {
			err = r.UnsetCName(c.BoxName, c.PublicIp)
			if err == router.ErrCNameNotFound {
				return nil
			}
			if err != nil {
				if !args.boxDestroy {
					return err
				}
				log.Errorf("---- ignored error removing route for %q during box %q destroy: %s", c.Address(), c.BoxName, err)
			}
			c.Routable = true
			toRollback <- c

			fmt.Fprintf(args.writer, lb.W(lb.CONTAINER_DEPLOY, lb.INFO, fmt.Sprintf("---> Removed route from container (%s, %s)", c.BoxName, c.ShortId())))
			return nil
		}, func(c *container.Container) {
			r.SetCName(c.BoxName, c.PublicIp)
		}, false)
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(changeUnitsPipelineArgs)
		r, err := getRouterForBox(args.box)
		if err != nil {
			log.Errorf("---- [remove-old-routes:Backward] Error geting router: %s", err.Error())
		}
		for _, cont := range args.toRemove {
			if !cont.Routable {
				continue
			}
			err = r.SetCName(cont.BoxName, cont.PublicIp)
			if err != nil {
				log.Errorf("---- [remove-old-routes:Backward] Error adding back route for (%s,%s): %s", cont.BoxName, cont.Id, err.Error())
			}
		}
	},
	MinParams: 1,
}

/*var bindAndHealthcheck = action.Action{
	Name: "bind-and-healthcheck",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(changeUnitsPipelineArgs)
		webProcessName, err := getImageWebProcessName(args.imageId)
		if err != nil {
			log.Errorf("[WARNING] cannot get the name of the web process: %s", err)
		}
		newContainers := ctx.Previous.([]container)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		doHealthcheck := true
		for _, c := range args.toRemove {
			if c.Status == provision.StatusError.String() || c.Status == provision.StatusStopped.String() {
				doHealthcheck = false
				break
			}
		}
		fmt.Fprintf(writer, "\n---- Binding and checking %d new units ----\n", len(newContainers))
		return newContainers, runInContainers(newContainers, func(c *container, toRollback chan *container) error {
			unit := c.asUnit(args.app)
			err := args.box.BindUnit(&unit)
			if err != nil {
				return err
			}
			toRollback <- c
			if doHealthcheck && c.ProcessName == webProcessName {
				err = runHealthcheck(c, writer)
				if err != nil {
					return err
				}
			}
			err = args.provisioner.runRestartAfterHooks(c, writer)
			if err != nil {
				return err
			}
			fmt.Fprintf(writer, " ---> Bound and checked unit %s [%s]\n", c.shortID(), c.ProcessName)
			return nil
		}, func(c *container) {
			unit := c.asUnit(args.app)
			err := args.box.UnbindUnit(&unit)
			if err != nil {
				log.Errorf("Unable to unbind unit %q: %s", c.ID, err)
			}
		}, true)
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(changeUnitsPipelineArgs)
		newContainers := ctx.FWResult.([]container)
		for _, c := range newContainers {
			unit := c.asUnit(args.app)
			err := args.box.UnbindUnit(&unit)
			if err != nil {
				log.Errorf("Removed binding for unit %q: %s", c.ID, err)
			}
		}
	},
}
*/
