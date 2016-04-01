/*
** Copyright [2013-2016] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package one

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/provision/one/machine"
)

const (
	START   = "start"
	STOP    = "stop"
	RESTART = "restart"
)

type runMachineActionsArgs struct {
	box           *provision.Box
	writer        io.Writer
	imageId       string
	isDeploy      bool
	machineStatus utils.Status
	provisioner   *oneProvisioner
}

//If there is a previous machine created and it has a status, we use that.
// eg: if it we have deployed, then make it created after a machine is created in ONE.
var updateStatusInScylla = action.Action{
	Name: "update-status-scylla",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		fmt.Fprintf(writer, "  update status for machine (%s, %s)\n", args.box.GetFullName(), args.machineStatus.String())

		var mach machine.Machine
		if ctx.Previous != nil {
			mach = ctx.Previous.(machine.Machine)
		} else {
			mach = machine.Machine{
				Id:         args.box.Id,
				AccountsId: args.box.AccountsId,
				CartonId:   args.box.CartonId,
				Level:      args.box.Level,
				Name:       args.box.GetFullName(),
				Status:     args.machineStatus,
				Image:      args.imageId,
			}
		}

		if err := mach.SetStatus(mach.Status); err != nil {
			return err, nil
		}
		fmt.Fprintf(writer, "  update status for machine (%s, %s) OK\n", args.box.GetFullName(), args.machineStatus.String())
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.SetStatus(constants.StatusError)
	},
}

var createMachine = action.Action{
	Name: "create-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		fmt.Fprintf(writer, "  create machine for box (%s, image:%s)/%s\n", args.box.GetFullName(), args.imageId, args.box.Compute)

		err := mach.Create(&machine.CreateArgs{
			Box:         args.box,
			Compute:     args.box.Compute,
			Deploy:      true,
			Provisioner: args.provisioner,
		})
		if err != nil {
			return nil, err
		}
		mach.Status = constants.StatusLaunched
		fmt.Fprintf(writer, "  create machine for box (%s, image:%s)/%s OK\n", args.box.GetFullName(), args.imageId, args.box.Compute)
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "\n   removing err machine %s\n", c.Name)

		err := c.Remove(args.provisioner)
		if err != nil {
			fmt.Fprintf(args.writer, "\n   removing err machine\n %s\n", err.Error())
		}
	},
}

var deductCons = action.Action{
	Name: "deduct-cons-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		fmt.Fprintf(writer, "  deduct cons of machine (%s, %s)\n", args.box.GetFullName(), args.machineStatus.String())
		mach.Deduct()
		fmt.Fprintf(writer, "  deduct cons of machine (%s, %s) OK\n", args.box.GetFullName(), args.machineStatus.String())
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.SetStatus(constants.StatusError)
	},
}

var destroyOldMachine = action.Action{
	Name: "destroy-old-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		fmt.Fprintf(writer, "\n   destroying old machine %s ----\n", mach.Name)
		err := mach.Remove(args.provisioner)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, "\n   destroyed old machine (%s, %s) OK\n", mach.Id, mach.Name)
		return ctx.Previous, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var startMachine = action.Action{
	Name: "start-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		fmt.Fprintf(writer, "\n   starting  machine %s\n", mach.Name)
		err := mach.LifecycleOps(args.provisioner, START)
		if err != nil {
			return nil, err
		}
		mach.Status = constants.StatusStarted
		fmt.Fprintf(writer, "\n   started machine (%s, %s) OK\n", mach.Id, mach.Name)
		return mach, nil
	},

	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var stopMachine = action.Action{
	Name: "stop-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		fmt.Fprintf(writer, "\n   stopping  machine %s\n", mach.Name)
		err := mach.LifecycleOps(args.provisioner, STOP)
		if err != nil {
			return nil, err
		}
		mach.Status = constants.StatusStopped
		fmt.Fprintf(writer, "\n   stopped machine (%s, %s) OK\n", mach.Id, mach.Name)
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var restartMachine = action.Action{
	Name: "restart-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		fmt.Fprintf(writer, "\n   restarting  machine %s\n", mach.Name)

		err := mach.LifecycleOps(args.provisioner, RESTART)
		if err != nil {
			return nil, err
		}
		mach.Status = constants.StatusRunning
		fmt.Fprintf(writer, "\n   restarted machine (%s, %s) OK\n", mach.Id, mach.Name)
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var changeStateofMachine = action.Action{
	Name: "change-state-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		fmt.Fprintf(writer, "  change state of machine (%s, %s)\n", args.box.GetFullName(), args.machineStatus.String())
		mach := machine.Machine{
			Id:       args.box.Id,
			CartonId: args.box.CartonId,
			Level:    args.box.Level,
			Name:     args.box.GetFullName(),
		}
		mach.ChangeState(args.machineStatus)
		mach.Status = args.machineStatus
		fmt.Fprintf(writer, "  change state of machine (%s, %s) OK\n", args.box.GetFullName(), args.machineStatus.String())
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.SetStatus(constants.StatusError)
	},
}

var addNewRoute = action.Action{
	Name: "add-new-route",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		r, err := getRouterForBox(args.box)
		if err != nil {
			return mach, err
		}
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}

		fmt.Fprintf(writer, "\n   adding route to machine (%s, %s)\n", mach.Name, args.box.PublicIp)
		err = r.SetCName(mach.Name, args.box.PublicIp)
		if err != nil {
			return mach, err
		}
		mach.SetRoutable(args.box.PublicIp)
		fmt.Fprintf(writer, "   added route to machine (%s, %s) OK\n", mach.Name, args.box.PublicIp)

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.FWResult.(machine.Machine)
		r, err := getRouterForBox(args.box)
		w := args.writer
		if w == nil {
			w = ioutil.Discard
		}
		if err != nil {
			fmt.Fprintf(w, "   destroy route error\n     %s", err.Error())
		}
		fmt.Fprintf(w, "\n   destroy routes from created machine  (%s, %s)\n", mach.Id, mach.Name)
		if mach.Routable {
			err = r.UnsetCName(mach.Name, args.box.PublicIp)
			if err != nil {
				fmt.Fprintf(w, "   destroy route error (%s, %s)\n    %s", mach.Name, args.box.PublicIp, err.Error())
			}
			fmt.Fprintf(w, "\n   destroy route from machine (%s, %s) OK\n", mach.Id, mach.Name)
		}
	},
	OnError: rollbackNotice,
}

var destroyOldRoute = action.Action{
	Name: "destroy-old-route",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		r, err := getRouterForBox(args.box)
		if err != nil {
			return mach, err
		}
		w := args.writer
		if w == nil {
			w = ioutil.Discard
		}
		mach.SetRoutable(args.box.PublicIp)
		fmt.Fprintf(w, "\n   destroy routes from created machine\n")
		if mach.Routable {
			err = r.UnsetCName(mach.Name, args.box.PublicIp)
			if err != nil {
				fmt.Fprintf(w, "   destroy route error (%s, %s)\n    %s", mach.Name, args.box.PublicIp, err.Error())
			}
			fmt.Fprintf(w, "\n   destroy route from machine (%s, %s)\n", mach.Name, args.box.PublicIp)
		} else {
			fmt.Fprintf(w, "\n   skip destroy routes from created machine (%s, %s) OK\n", mach.Name, args.box.PublicIp)
		}
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.FWResult.(machine.Machine)
		r, err := getRouterForBox(args.box)
		w := args.writer
		if w == nil {
			w = ioutil.Discard
		}
		if err != nil {
			fmt.Fprintf(w, "   destroy route error (%s, %s)\n    %s", mach.Name, args.box.PublicIp, err.Error())
		}
		fmt.Fprintf(w, "\n   addding back routes to old machine\n")
		if mach.Routable {
			err = r.SetCName(mach.Name, args.box.PublicIp)
			if err != nil {
				fmt.Fprintf(w, "   destroy error (%s, %s)\n     %s", mach.Name, args.box.PublicIp, err.Error())
			}
			fmt.Fprintf(w, "   adding route to machine (%s, %s) OK\n", mach.Name, args.box.PublicIp)
		}
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var followLogs = action.Action{
	Name: "follow-logs",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		c, ok := ctx.Previous.(machine.Machine)
		if !ok {
			return nil, errors.New("Previous result must be a machine.")
		}
		args := ctx.Params[0].(runMachineActionsArgs)
		err := c.Logs(args.provisioner, args.writer)
		if err != nil {
			log.Errorf("error on get logs for machine %s - %s", c.Name, err)
			return nil, err
		}

		return args.imageId, nil
	},
	Backward: func(ctx action.BWContext) {
	},
	MinParams: 1,
}

var rollbackNotice = func(ctx action.FWContext, err error) {
	args := ctx.Params[0].(runMachineActionsArgs)
	if args.writer != nil {
		fmt.Fprintf(args.writer, "\n==> ROLLBACK     \n%s\n", err)
	}
}
