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
	// "errors"
	// "fmt"
	 "io"
	 "io/ioutil"

	//log "github.com/Sirupsen/logrus"
	 "github.com/megamsys/libgo/action"
	// "github.com/megamsys/libgo/events/alerts"
	"github.com/megamsys/libgo/utils"
	//"github.com/megamsys/vertice/marketplaces"
	 constants "github.com/megamsys/libgo/utils"
	// vm "github.com/megamsys/opennebula-go/virtualmachine"
//	 lb "github.com/megamsys/vertice/logbox"
	 "github.com/megamsys/vertice/marketplaces/provision"
	 "github.com/megamsys/vertice/marketplaces/provision/one/machine"
)

type runMachineActionsArgs struct {
	box           *provision.Box
	writer        io.Writer
	imageId       string
	isRaw         bool
	machineStatus utils.Status
	provisioner   *oneProvisioner
}

//If there is a previous machine created and it has a status, we use that.
// eg: if it we have deployed, then make it created after a machine is created in ONE.
//
var machCreating = action.Action{
	Name: "machine-struct-creating",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		mach := machine.Machine{
			AccountId: args.box.AccountId,
			CartonId: args.box.Id,
			Name: args.box.Name,
			Region: args.box.Region,
			PublicUrl: args.box.PublicUrl,
		}
		  mach.Status =  args.machineStatus
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var createRawISOImage = action.Action{
	Name: "create-rawimage-iso",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
   err := mach.CreateISO(args.provisioner)
	 if err != nil {
		 return mach, err
	 }
	 mach.Status = constants.StatusCreating
	 return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}

var waitUntillImageReady = action.Action{
	Name: "wait-for-image-ready",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		mach := ctx.Previous.(machine.Machine)
		args := ctx.Params[0].(runMachineActionsArgs)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
		if err := mach.IsImageReady(args.provisioner); err != nil {
			return nil, err
		}
		mach.Status = constants.StatusActive
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
		//do you want to add it back.
	},
}

var updateRawImageId = action.Action{
	Name: "update-rawimage-id",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
   err := mach.UpdateRawImageId()
	 if err != nil {
		 return mach, err
	 }
		return mach, nil
	},
	Backward: func(ctx action.BWContext) {
	},
}


var updateRawStatus = action.Action{
	Name: "update-rawimage-status",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runMachineActionsArgs)
		mach := ctx.Previous.(machine.Machine)
		writer := args.writer
		if writer == nil {
			writer = ioutil.Discard
		}
   err := mach.UpdateRawStatus()
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
   err := mach.CreateDatablock(args.provisioner)
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
