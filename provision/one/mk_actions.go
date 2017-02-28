package one

import (
	"fmt"
	"github.com/megamsys/libgo/action"
	constants "github.com/megamsys/libgo/utils"
	lb "github.com/megamsys/vertice/logbox"
	"github.com/megamsys/vertice/provision/one/machine"
	"io/ioutil"
)

var createImage = action.Action{
	Name: "create-rawimage-iso",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.CreateImage(args.provisioner)
		if err != nil {
			return mach, err
		}
		mach.Status = constants.StatusCreating
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var updateImage = action.Action{
	Name: "update-rawimage-id",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.UpdateImage()
		if err != nil {
			return mach, err
		}
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var updateImageStatus = action.Action{
	Name: "update-rawimage-status",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.UpdateImageStatus()
		if err != nil {
			return mach, err
		}
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var createDatablockImage = action.Action{
	Name: "create-datablock",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.CreateDatablock(args.provisioner, args.box)
		if err != nil {
			return mach, err
		}
		mach.Status = constants.StatusCreating
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var updateMarketplaceImageId = action.Action{
	Name: "update-marketplace-block-id",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.UpdateMarketImageId()
		if err != nil {
			return mach, err
		}
		mach.Status = constants.StatusDataBlockCreated
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var createInstanceForCustomize = action.Action{
	Name: "create-instance-to-customize",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.CreateInstance(args.provisioner, args.box)
		if err != nil {
			return mach, err
		}
		mach.Status = constants.StatusLaunching
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var updateMarketplaceStatus = action.Action{
	Name: "update-rawimage-status",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.UpdateMarketplaceStatus()
		if err != nil {
			return mach, err
		}
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var waitUntillVmReady = action.Action{
	Name: "update-for-vm-running",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.MarketplaceInstanceState(args.provisioner)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  error start machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		mach.Status = constants.StatusLaunched
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var getMarketplaceVncPost = action.Action{
	Name: "get-vnc-host-ip-port",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.GetMarketplaceVNC(args.provisioner)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  error start machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		mach.Status = constants.StatusVncHostUpdating
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var updateMarketplaceVnc = action.Action{
	Name: "update-vnc-host-ip-port",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.UpdateMarketplaceVNC()
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  error start machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		mach.Status = constants.StatusVncHostUpdated
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}
