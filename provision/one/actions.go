/*
** Copyright [2013-2015] [Megam Systems]
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
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/provision/one/machine"
)

type runMachineActionsArgs struct {
	box           *provision.Box
	writer        io.Writer
	imageId       string
	isDeploy      bool
	machineStatus provision.Status
	provisioner   *oneProvisioner
}

var updateStatusInRiak = action.Action{
	Name: "update-status-riak",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		log.Debugf("update status for machine %s image %s for %s", args.box.GetFullName(), args.imageId, args.box.Compute)

		mach := machine.Machine{
			Id:    args.box.Id,
			Level: args.box.Level,
			Name:  args.box.GetFullName(),
			Image: args.imageId,
		}

		mach.SetStatus(args.machineStatus)
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		c.SetStatus(provision.StatusError)
	},
}

var createMachine = action.Action{
	Name: "create-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		log.Debugf("create machine for box %s, on image %s, with %s", args.box.GetFullName(), args.imageId, args.box.Compute)

		err := mach.Create(&machine.CreateArgs{
			Box:         args.box,
			Compute:     args.box.Compute,
			Deploy:      true,
			Provisioner: args.provisioner,
		})
		if err != nil {
			return nil, err
		}
		args.machineStatus = provision.StatusCreating
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, "\n---- Removing %d old machine %s ----\n", c.Name)

		err := c.Remove(args.provisioner)
		if err != nil {
			log.Errorf("Failed to remove the machine %q: %s", c.Name, c.Id, err)
		}
	},
}

var removeOldMachine = action.Action{
	Name: "remove-old-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		fmt.Fprintf(writer, "\n---- Removing old machine %s ----\n", mach.Name)

		err := mach.Remove(args.provisioner)
		if err != nil {
			log.Errorf("Ignored error trying to remove old machine %q: %s", mach.Id, err)
		}
		fmt.Fprintf(writer, "\n---- Removed old machine %s [%s]\n", mach.Id, mach.Name)
		return ctx.Previous, nil
	},
	Backward: func(ctx action.BWContext) {
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var addNewRoutes = action.Action{
	Name: "add-new-routes",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		r, err := getRouterForBox(args.box)
		if err != nil {
			return nil, err
		}

		if _, err := r.Addr(args.box.GetFullName()); err != nil {
			log.Errorf("[WARNING] route attached: %s", err)
		}

		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}

		fmt.Fprintf(writer, "\n---- Adding route to new machine ----\n")
		err = r.SetCName(mach.Name, args.box.Ip)
		if err != nil {
			return mach, err
		}
		mach.Routable = true
		fmt.Fprintf(writer, " ---> Added route to machine %s [%s]\n", mach.Name, args.box.Ip)

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.FWResult.(machine.Machine)
		r, err := getRouterForBox(args.box)
		if err != nil {
			log.Errorf("[add-new-routes:Backward] Error geting router: %s", err.Error())
		}
		w := args.writer
		if w == nil {
			w = ioutil.Discard
		}
		fmt.Fprintf(w, "\n---- Removing routes from created machine ----\n")
		if mach.Routable {
			err = r.UnsetCName(mach.Name, args.box.Ip)
			if err != nil {
				log.Errorf("[add-new-routes:Backward] Error removing route for %s [%s]: %s", mach.Name, args.box.Ip, err.Error())
			}
			fmt.Fprintf(w, " ---> Removed route from machine %s [%s]\n", mach.Id, mach.Name)
		}
	},
	OnError: rollbackNotice,
}

var removeOldRoutes = action.Action{
	Name: "remove-old-routes",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)

		r, err := getRouterForBox(args.box)
		if err != nil {
			log.Errorf("[remove-old-routes] Error geting router: %s", err.Error())
			return mach, err
		}
		w := args.writer
		if w == nil {
			w = ioutil.Discard
		}
		fmt.Fprintf(w, "\n---- Removing routes from created machine ----\n")
		if mach.Routable {
			err = r.UnsetCName(mach.Name, args.box.Ip)
			if err != nil {
				log.Errorf("[add-new-routes:Backward] Error removing route for %s [%s]: %s", mach.Name, args.box.Ip, err.Error())
			}
			fmt.Fprintf(w, " ---> Removed route from unit %s [%s]\n", mach.Name, args.box.Ip)
		}
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.FWResult.(machine.Machine)
		r, err := getRouterForBox(args.box)
		if err != nil {
			log.Errorf("[remove-old-routes:Backward] Error geting router: %s", err.Error())
		}
		w := args.writer
		if w == nil {
			w = ioutil.Discard
		}
		fmt.Fprintf(w, "\n---- Adding back routes to old machines ----\n")
		if mach.Routable {

			err = r.SetCName(mach.Name, args.box.Ip)
			if err != nil {
				log.Errorf("[remove-old-routes:Backward] Error adding back route for %s [%s]: %s", mach.Name, args.box.Ip, err.Error())
			}
			fmt.Fprintf(w, " ---> Added route to unit %s [%s]\n", mach.Name, args.box.Ip)
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
		fmt.Fprintf(args.writer, "\n**** ROLLING BACK AFTER FAILURE ****\n ---> %s <---\n", err)
	}
}
