package one

import (
	"fmt"
	"github.com/megamsys/libgo/action"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/images"
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
		err := mach.CreateImage(args.provisioner, images.CD_ROM)
		if err != nil {
			return mach, err
		}
		mach.Status = constants.StatusCreating
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		mach := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		err := mach.RemoveImage(args.provisioner)
		if err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.DESTORYING, lb.ERROR, fmt.Sprintf("  removing err image %s", err.Error())))
		}
	},
}

var removeImage = action.Action{
	Name: "remove-rawimage-iso",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.RemoveImage(args.provisioner)
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
		mach.Status = constants.StatusDataBlockCreating
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		mach := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		err := mach.RemoveDatablock(args.provisioner)
		if err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.DESTORYING, lb.ERROR, fmt.Sprintf("  removing err datablock %s", err.Error())))
		}
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
		mach := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		fmt.Fprintf(args.writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  removing instance %s", ctx.CauseOf.Error())))
		err := mach.Remove(args.provisioner)
		if err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  removing err instance %s", err.Error())))
		}
	},
}

var attachDatablockImage = action.Action{
	Name: "attach-datablock-image-to-vm",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.AttachDatablock(args.provisioner, args.box)
		if err != nil {
			return mach, err
		}
		mach.Status = constants.StatusDataBlockCreated
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var updateMarketplaceStatus = action.Action{
	Name: "update-marketplaces-status",
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
		mach := ctx.FWResult.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		mach.Status = constants.StatusPreError
		err := mach.UpdateMarketplaceError(ctx.CauseOf)
		if err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  failure update marketplace status  %s", err.Error())))
		}
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

var stopMachineIfRunning = action.Action{
	Name: "shutdown-if-machine-running",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.StopMarkplaceInstance(args.provisioner)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  error stop machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		mach.Status = constants.StatusImageSaving
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var waitForsaveImage = action.Action{
	Name: "wait-for-save-image",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.CheckSaveImage(args.provisioner)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  error start machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		mach.Status = constants.StatusImageSaved
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var makeImageAsPersistent = action.Action{
	Name: "make-image-as-persistent",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.ImagePersistent(args.provisioner)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  error start machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var changeAsOsImage = action.Action{
	Name: "make-as-os-image",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.ImageTypeChange(args.provisioner)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  error start machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		mach.Status = constants.StatusImageReady
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var removeInstance = action.Action{
	Name: "remove-instance-vm",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		err := mach.RemoveInstance(args.provisioner)
		if err != nil {
			fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("  error start machine ( %s)", args.box.GetFullName())))
			return nil, err
		}
		mach.Status = constants.StatusImageReady
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}
