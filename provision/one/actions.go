/*
** Copyright [2013-2017] [Megam Systems]
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
	"github.com/megamsys/libgo/events/alerts"
	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	vm "github.com/megamsys/opennebula-go/virtualmachine"
	"github.com/megamsys/vertice/carton"
	lb "github.com/megamsys/vertice/logbox"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/provision/one/machine"
)

type runMachineActionsArgs struct {
	box           *provision.Box
	writer        io.Writer
	imageId       string
	isDeploy      bool
	machineStatus utils.Status
	machineState  utils.State
	provisioner   *oneProvisioner
	process       string
}

//If there is a previous machine created and it has a status, we use that.
// eg: if it we have deployed, then make it created after a machine is created in ONE.

var machCreating = action.Action{
	Name: "machine-struct-creating",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" creating struct machine (%s, %s)", args.box.GetFullName(), args.machineStatus.String())))
		mach := machine.Machine{
			Id:           args.box.Id,
			AccountId:    args.box.AccountId,
			CartonId:     args.box.CartonId,
			CartonsId:    args.box.CartonsId,
			Level:        args.box.Level,
			Name:         args.box.GetFullName(),
			Status:       args.machineStatus,
			State:        args.machineState,
			Image:        args.imageId,
			StorageType:  args.box.StorageType,
			PublicUrl:    args.box.PublicUrl,
			Region:       args.box.Region,
			VMId:         args.box.InstanceId,
			VCPUThrottle: args.provisioner.vcpuThrottle,
		}
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" creating struct machine (%s, %s)OK", args.box.GetFullName(), args.machineStatus.String())))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var updateStatusInScylla = action.Action{
	Name: "update-status-scylla",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		mach := ctx.Previous.(machine.Machine)
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" update status for machine (%s, %s)", args.box.GetFullName(), args.machineStatus.String())))
		if err := mach.SetStatus(mach.Status); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" update status for machine (%s, %s)OK", args.box.GetFullName(), args.machineStatus.String())))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		w := args.writer
		if w == nil {
			w = ioutil.Discard
		}
		var status constants.Status
		if args.isDeploy {
			status = constants.StatusPreError
			_ = carton.DoneNotify(args.box, args.writer, alerts.FAILURE)
		} else {
			status = constants.StatusError
		}

		c.SetStatusErr(status, ctx.CauseOf)
	},
}

var checkBalances = action.Action{
	Name: "balance-check",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" check balance for user (%s) machine (%s)", args.box.AccountId, args.box.GetFullName())))
		err := mach.CheckCredits(args.box, writer)
		if err != nil {
			_ = mach.SetMileStone(constants.StateMachineParked)
			_ = mach.SetStatus(constants.StatusInsufficientFund)
			return nil, err
		}
		mach.SetStatus(constants.StatusBalanceVerified)
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" check balance for user (%s) machine (%s) OK", args.box.AccountId, args.box.GetFullName())))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
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
		err := mach.SetStatus(mach.Status)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" create machine for box (%s, image:%s)/%s", args.box.GetFullName(), args.imageId, args.box.Compute)))
		err = mach.Create(&machine.CreateArgs{
			Box:         args.box,
			Compute:     args.box.Compute,
			Deploy:      true,
			Provisioner: args.provisioner,
		})
		if err != nil {
			return nil, err
		}
		mach.State = constants.StateInitialized
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" create machine for box (%s, image:%s)/%s OK", args.box.GetFullName(), args.imageId, args.box.Compute)))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		log.Debugf("create machine backward state : %v", c.State)
		if c.State != constants.StateInitialized {
			log.Debugf(" backward removing machine")
			err := c.Remove(args.provisioner)
			if err != nil {
				fmt.Fprintf(args.writer, lb.W(lb.DESTORYING, lb.ERROR, fmt.Sprintf("  removing err machine %s", err.Error())))
			}
		}
	},
}

var getVmHostIpPort = action.Action{
	Name: "gethost-port",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.VmHostIpPort(&machine.CreateArgs{Provisioner: args.provisioner})
		if err != nil {
			return nil, err
		}
		mach.Status = constants.StatusVncHostUpdating

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {

	},
}

var updateVnchostPostInScylla = action.Action{
	Name: "update-vnc-host-port",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		err := mach.UpdateVncHostPost()
		if err != nil {
			return nil, err
		}
		mach.Status = constants.StatusVncHostUpdated
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {

	},
}

var updateNetworkIps = action.Action{
	Name: "update-vm-assigned-ips",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		err := mach.UpdateVMIps(args.provisioner)
		if err != nil {
			return nil, err
		}
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {

	},
}

var setFinalStatus = action.Action{
	Name: "set-final-status",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		mach.Status = constants.StatusVMBooting
		return mach, nil
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

		fmt.Fprintf(writer, lb.W(lb.DESTORYING, lb.INFO, fmt.Sprintf("  destroying old machine %s ----", mach.Name)))
		err := mach.Remove(args.provisioner)
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(writer, lb.W(lb.DESTORYING, lb.INFO, fmt.Sprintf("  destroyed old machine (%s, %s) OK", mach.Id, mach.Name)))
		return ctx.Previous, nil
	},
	Backward: func(ctx action.BWContext) {

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

		fmt.Fprintf(writer, lb.W(lb.STARTING, lb.INFO, fmt.Sprintf("  starting  machine %s", mach.Name)))
		err := mach.LifecycleOps(args.provisioner, args.process)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.STARTING, lb.ERROR, fmt.Sprintf("  error start machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		err = mach.WaitUntillVMState(args.provisioner, vm.ACTIVE, vm.RUNNING)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.STARTING, lb.ERROR, fmt.Sprintf("  error start machine ( %s)", args.box.GetFullName())))
			return nil, err
		}

		mach.Status = constants.StatusStarted
		mach.State = constants.StateRunning
		fmt.Fprintf(writer, lb.W(lb.STARTING, lb.INFO, fmt.Sprintf("  starting  machine (%s, %s) OK", mach.Id, mach.Name)))
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

		fmt.Fprintf(writer, lb.W(lb.STOPPING, lb.INFO, fmt.Sprintf("\n   stopping  machine %s", mach.Name)))
		err := mach.LifecycleOps(args.provisioner, args.process)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.STOPPING, lb.ERROR, fmt.Sprintf("  error stop machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		err = mach.WaitUntillVMState(args.provisioner, vm.POWEROFF, vm.LCM_INIT)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.STOPPING, lb.ERROR, fmt.Sprintf("  error stop machine ( %s)", args.box.GetFullName())))
			return nil, err
		}

		mach.Status = constants.StatusStopped
		mach.State = constants.StateStopped
		fmt.Fprintf(writer, lb.W(lb.STOPPING, lb.INFO, fmt.Sprintf("\n   stopping  machine (%s, %s)OK", mach.Id, mach.Name)))
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

		fmt.Fprintf(writer, lb.W(lb.RESTARTING, lb.INFO, fmt.Sprintf("restarting  machine %s", mach.Name)))
		err := mach.LifecycleOps(args.provisioner, args.process)
		if err != nil {
			return nil, err
		}
		mach.Status = constants.StatusRunning

		fmt.Fprintf(writer, lb.W(lb.RESTARTING, lb.INFO, fmt.Sprintf("restarting  machine (%s, %s)OK", mach.Id, mach.Name)))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var suspendMachine = action.Action{
	Name: "suspend-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}

		fmt.Fprintf(writer, lb.W(lb.STOPPING, lb.INFO, fmt.Sprintf("\n   suspending  machine %s", mach.Name)))
		err := mach.LifecycleOps(args.provisioner, args.process)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.STOPPING, lb.ERROR, fmt.Sprintf("  error suspend machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		err = mach.WaitUntillVMState(args.provisioner, vm.POWEROFF, vm.LCM_INIT)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.STOPPING, lb.ERROR, fmt.Sprintf("  error suspend machine ( %s)", args.box.GetFullName())))
			return nil, err
		}

		mach.Status = constants.StatusSuspended
		mach.State = constants.StateStopped
		fmt.Fprintf(writer, lb.W(lb.STOPPING, lb.INFO, fmt.Sprintf("\n   suspending  machine (%s, %s)OK", mach.Id, mach.Name)))
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

		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("  change state of machine (%s, %s)", args.box.GetFullName(), args.machineStatus.String())))
		mach := machine.Machine{
			Id:       args.box.Id,
			CartonId: args.box.CartonId,
			Level:    args.box.Level,
			Name:     args.box.GetFullName(),
		}
		err := mach.SetStatus(constants.StatusStateupping)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  error change state of machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		err = mach.ChangeState(args.machineStatus)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  error publish state change of machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		if args.box.PublicIp != "" {
			mach.Status = constants.StatusNetworkCreating

		} else {
			mach.Status = constants.StatusNetworkSkipped
		}

		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("  change state of machine (%s, %s)OK", args.box.GetFullName(), args.machineStatus.String())))
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

		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("adding route to machine (%s, %s)", mach.Name, args.box.PublicIp)))
		err = r.SetCName(mach.Name, args.box.PublicIp)
		if err != nil {
			return mach, err
		}
		mach.SetRoutable(args.box.PublicIp)
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("adding route to machine (%s, %s)OK", mach.Name, args.box.PublicIp)))
		mach.Status = constants.StatusNetworkCreated
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

			fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf(" destroy route error    %s", err.Error())))
		}

		fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("   destroy routes from created machine  (%s, %s)", mach.Id, mach.Name)))
		if mach.Routable {
			err = r.UnsetCName(mach.Name, args.box.PublicIp)
			if err != nil {

				fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("   destroy route error (%s, %s)    %s", mach.Name, args.box.PublicIp, err.Error())))
			}

			fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("destroy route from machine (%s, %s) OK", mach.Id, mach.Name)))
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
		fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("destroy routes from created machine")))
		if mach.Routable {
			err = r.UnsetCName(mach.Name, args.box.PublicIp)
			if err != nil {

				fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("destroy route error (%s, %s)   %s", mach.Name, args.box.PublicIp, err.Error())))
			}

			fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("  destroy route from machine (%s, %s)", mach.Name, args.box.PublicIp)))
		} else {

			fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("  skip destroy routes from created machine (%s, %s) OK", mach.Name, args.box.PublicIp)))
		}
		mach.Status = constants.StatusDestroyed
		mach.State = constants.StateDestroyed
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

			fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("   destroy route error (%s, %s)   %s", mach.Name, args.box.PublicIp, err.Error())))
		}

		fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("  addding back routes to old machine")))
		if mach.Routable {
			err = r.SetCName(mach.Name, args.box.PublicIp)
			if err != nil {

				fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("destroy error (%s, %s)     %s", mach.Name, args.box.PublicIp, err.Error())))
			}

			fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("   adding route to machine (%s, %s) OK", mach.Name, args.box.PublicIp)))
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

		fmt.Fprintf(args.writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("==> ROLLBACK     %s", err)))

	}
}

var createSnapshot = action.Action{
	Name: "create-snapshot-disk",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}

		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf("  creating snapshot machine %s ----", mach.Name)))
		err := mach.CreateDiskSnap(args.provisioner)
		if err != nil {
			return mach, err
		}

		mach.Status = constants.StatusSnapCreated

		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" creating snapshot machine (%s, %s) OK", mach.Id, mach.Name)))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.FWResult.(machine.Machine)
		w := args.writer
		if w == nil {
			w = ioutil.Discard
		}
		if err := mach.RemoveSnapshot(args.provisioner); err != nil {
			fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  snapshot remove failure error (%s)   %s", mach.Name, err.Error())))
		}
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var restoreVirtualMachine = action.Action{
	Name: "restore-restore-snapshot",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" remove snapshot for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.RestoreSnapshot(args.provisioner); err != nil {
			return nil, err
		}
		err := mach.WaitUntillVMState(args.provisioner, vm.POWEROFF, vm.LCM_INIT)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.STARTING, lb.ERROR, fmt.Sprintf("  error start machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		mach.Status = constants.StatusSnapRestored
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" remove snapshot for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var makeActiveSnap = action.Action{
	Name: "activate-current-snapshot",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" remove snapshot for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.MakeActiveSnapshot(); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" remove snapshot for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var removeSnapShot = action.Action{
	Name: "remove-snap-shot",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" remove snapshot for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.RemoveSnapshot(args.provisioner); err != nil {
			return nil, err
		}
		mach.Status = constants.StatusSnapDeleted
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" remove snapshot for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var createBackupImage = action.Action{
	Name: "create-backup-image",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}

		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf("  creating backup machine %s ----", mach.Name)))
		err := mach.CreateDiskImage(args.provisioner)
		if err != nil {
			return nil, err
		}

		mach.Status = constants.StatusBackupCreated

		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" creating backup machine (%s, %s) OK", mach.Id, mach.Name)))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.FWResult.(machine.Machine)
		mach.Status = constants.Status("error")
		w := args.writer
		if w == nil {
			w = ioutil.Discard
		}
		if err := mach.RemoveBackupImage(args.provisioner); err != nil {
			fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  backup remove failure error (%s)   %s", mach.Name, err.Error())))
		}
		err := mach.UpdateBackupStatus(mach.Status)
		if err != nil {
			fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  backup create failure update error (%s)   %s", mach.Name, err.Error())))
		}
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var removeBackup = action.Action{
	Name: "remove-backup-image",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" remove snapshot for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.RemoveBackupImage(args.provisioner); err != nil {
			return nil, err
		}
		mach.Status = constants.StatusBackupDeleted
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" remove snapshot for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var mileStoneUpdate = action.Action{
	Name: "change-milestone-state",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" update milestone state for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.SetMileStone(mach.State); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" update milestone state for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		var err error
		c := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("\n---- State Changing Backward for %s ----", args.box.GetFullName())))
		var state constants.State
		if args.isDeploy {
			state = constants.StatePreError
			_ = carton.DoneNotify(args.box, args.writer, alerts.FAILURE)
		} else {
			state = constants.StateError
		}

		err = c.SetMileStone(state)
		if err != nil {
			log.Errorf("---- [state-change:Backward]\n     %s", err.Error())
		}
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var setFinalState = action.Action{
	Name: "set-final-state",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		mach.Status = constants.StatusStateupped
		return mach, nil
	},
}

var addNewStorage = action.Action{
	Name: "add-new-storage",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf("  attaching new disk to machine %s ----", mach.Name)))
		err := mach.AttachNewDisk(args.provisioner)
		if err != nil {
			return nil, err
		}
		mach.Status = constants.StatusDiskAttaching
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf("  attaching new disk to machine (%s, %s) OK", mach.Id, mach.Name)))
		return ctx.Previous, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var updateIdInSnapTable = action.Action{
	Name: "update-snap-table",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update snapshot status for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.UpdateSnap(); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update snapshot status for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var updateSnapStatus = action.Action{
	Name: "update-snap-status",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update snapshot status for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.UpdateSnapStatus(mach.Status); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update snapshot status for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.FWResult.(machine.Machine)
		w := args.writer
		if w == nil {
			w = ioutil.Discard
		}
		err := mach.UpdateSnapStatus(constants.StatusError)
		if err != nil {
			fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  snapshot create failure update error (%s)   %s", mach.Name, err.Error())))
		}
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var updateIdInBackupTable = action.Action{
	Name: "update-backups-table",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update backups status for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.UpdateBackup(); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update backups status for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var updateSourceVMIdIps = action.Action{
	Name: "update-backups-vm-ips",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update backups source machine ips (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.UpdateBackupVMIps(); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update backups source machine ips  (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var createBackupMachine = action.Action{
	Name: "create-backup-machine",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.SetStatus(mach.Status)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" create machine for box (%s, image:%s)/%s", args.box.GetFullName(), args.imageId, args.box.Compute)))
		err = mach.CreateBackupVM(&machine.CreateArgs{
			Box:         args.box,
			Compute:     args.box.Compute,
			Deploy:      true,
			Provisioner: args.provisioner,
		})
		if err != nil {
			return nil, err
		}
		mach.State = constants.StateInitialized
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" create machine for box (%s, image:%s)/%s OK", args.box.GetFullName(), args.imageId, args.box.Compute)))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		log.Debugf("create machine backward state : %v", c.State)
		if c.State != constants.StateInitialized {
			log.Debugf(" backward removing machine")
			err := c.Remove(args.provisioner)
			if err != nil {
				fmt.Fprintf(args.writer, lb.W(lb.DESTORYING, lb.ERROR, fmt.Sprintf("  removing err machine %s", err.Error())))
			}
		}
	},
}

var updateBackupStatus = action.Action{
	Name: "update-backup-status",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update backups status for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.UpdateBackupStatus(mach.Status); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update backups status for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var waitUntillImageReady = action.Action{
	Name: "wait-for-image-ready",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" waiting to backups creating for machine (%s, %s)", args.box.GetFullName(), constants.SNAPSHOTTING)))
		if err := mach.IsImageReady(args.provisioner); err != nil {
			return nil, err
		}
		mach.Status = constants.StatusImageReady
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" waiting to backups creating  for machine (%s, %s)OK", args.box.GetFullName(), constants.SNAPSHOTTING)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}
var waitUntillSnapReady = action.Action{
	Name: "wait-for-snapshot-ready",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" waiting to backups creating for machine (%s, %s)", args.box.GetFullName(), constants.SNAPSHOTTING)))
		if err := mach.IsSnapReady(args.provisioner); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" waiting to backups creating  for machine (%s, %s)OK", args.box.GetFullName(), constants.SNAPSHOTTING)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var updateIdInDiskTable = action.Action{
	Name: "update-disk-table",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update disks status for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.UpdateDisk(args.provisioner); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update disks status for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var removeDiskStorage = action.Action{
	Name: "remove-disk-storage",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" remove disk from machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.RemoveDisk(args.provisioner); err != nil {
			return nil, err
		}
		mach.Status = constants.StatusDiskDetached
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" remove disk from machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var updateSnapQuotaCount = action.Action{
	Name: "update-quota-snapshots-count",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update quota for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.UpdateSnapQuotas(args.box.QuotaId); err != nil {
			return nil, err
		}
		mach.Status = constants.StatusQuotaUpdated
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update quota for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var updateVMQuota = action.Action{
	Name: "update-quota-for-vm",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update quota for machine (%s, %s)", args.box.GetFullName(), constants.LAUNCHED)))
		if err := mach.UpdateVMQuotas(args.box.QuotaId); err != nil {
			return nil, err
		}
		mach.Status = constants.StatusQuotaUpdated
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update quota for machine (%s, %s)OK", args.box.GetFullName(), constants.LAUNCHED)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("\n---- State Changing Backward for %s ----", args.box.GetFullName())))
		c.Status = constants.StatusRunning
		err := c.UpdateVMQuotas(args.box.QuotaId)
		if err != nil {
			log.Errorf("---- [state-change:Backward]\n     %s", err.Error())
		}
	},
	OnError:   rollbackNotice,
	MinParams: 1,
}

var checkQuotaState = action.Action{
	Name: "check-quota-state",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" check balance for user (%s) machine (%s)", args.box.AccountId, args.box.GetFullName())))
		err := mach.CheckQuotaState(args.box, writer)
		if err != nil {
			_ = mach.SetMileStone(constants.StateMachineParked)
			_ = mach.SetStatus(constants.StatusInsufficientFund)
			return nil, err
		}
		mach.SetStatus(constants.StatusBalanceVerified)
		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" check balance for user (%s) machine (%s) OK", args.box.AccountId, args.box.GetFullName())))
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("\n---- State Changing Backward for %s ----", args.box.GetFullName())))
		c.Status = constants.StatusDestroying
		err := c.UpdateVMQuotas(args.box.QuotaId)
		if err != nil {
			log.Errorf("---- [state-change:Backward]\n     %s", err.Error())
		}
	},
}

var updataPoliciesStatus = action.Action{
	Name: "update-policy-status",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update policy for machine (%s, %s)", args.box.GetFullName(), args.box.PolicyOps.Operation)))
		if err := mach.UpdatePolicyStatus(args.box.PolicyOps.Index); err != nil {
			return nil, err
		}
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" update policy for machine (%s, %s)OK", args.box.GetFullName(), args.box.PolicyOps.Operation)))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		c := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("\n---- update policy Backward for %s ----", args.box.GetFullName())))
		c.Status = constants.StatusPolicyFailure
		err := c.UpdatePolicyStatus(args.box.PolicyOps.Index)
		if err != nil {
			log.Errorf("---- [policy-update:Backward]\n     %s", err.Error())
		}
	},
}

var attachNetworks = action.Action{
	Name: "attach-networks",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" attach ips to machine (%s)", args.box.GetFullName())))
		if err := mach.AttachNetwork(args.box, args.provisioner); err != nil {
			return nil, err
		}
		mach.Status = constants.StatusAttachInprogres
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" attach ips to machine (%s)OK", args.box.GetFullName())))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var detachNetworks = action.Action{
	Name: "detach-networks",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" detach ip from machine (%s)", args.box.GetFullName())))
		if err := mach.DetachNetwork(args.box, args.provisioner); err != nil {
			return nil, err
		}
		mach.Status = constants.StatusDetachInprogres
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" detach ip from machine (%s)OK", args.box.GetFullName())))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var removeNetworkIps = action.Action{
	Name: "remove-network-ips",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" remove ips for machine (%s)", args.box.GetFullName())))
		if err := mach.RemoveNetworkIps(args.box); err != nil {
			return nil, err
		}
		mach.Status = constants.StatusActive
		fmt.Fprintf(writer, lb.W(lb.UPDATING, lb.INFO, fmt.Sprintf(" remove ips for machine (%s)OK", args.box.GetFullName())))

		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}
