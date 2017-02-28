package one

import (
	"github.com/megamsys/libgo/action"
	constants "github.com/megamsys/libgo/utils"
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
		mach.Status = constants.StatusCreating
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
		err := mach.UpdateMarketStatus()
		if err != nil {
			return mach, err
		}
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}
