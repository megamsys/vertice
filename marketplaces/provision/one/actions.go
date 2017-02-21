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

// import (
// 	"errors"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
//
// 	log "github.com/Sirupsen/logrus"
// 	"github.com/megamsys/libgo/action"
// 	"github.com/megamsys/libgo/events/alerts"
// 	"github.com/megamsys/libgo/utils"
// 	constants "github.com/megamsys/libgo/utils"
// 	vm "github.com/megamsys/opennebula-go/virtualmachine"
// 	"github.com/megamsys/vertice/carton"
// 	lb "github.com/megamsys/vertice/logbox"
// 	"github.com/megamsys/vertice/marketplaces/provision"
// 	"github.com/megamsys/vertice/marketplaces/provision/one/machine"
// )
//
// const (
// 	START   = "start"
// 	STOP    = "stop"
// 	RESTART = "restart"
// )
//
// type runMachineActionsArgs struct {
// 	box           *provision.Box
// 	writer        io.Writer
// 	imageId       string
// 	isDeploy      bool
// 	machineStatus utils.Status
// 	machineState  utils.State
// 	provisioner   *oneProvisioner
// }
//
// //If there is a previous machine created and it has a status, we use that.
// // eg: if it we have deployed, then make it created after a machine is created in ONE.
//
// var machCreating = action.Action{
// 	Name: "machine-struct-creating",
// 	Forward: func(ctx action.FWContext) (action.Result, error) {
// 		args := ctx.Params[0].(runMachineActionsArgs)
// 		writer := args.writer
// 		if writer == nil {
// 			writer = ioutil.Discard
// 		}
// 		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" creating struct machine (%s, %s)", args.box.GetFullName(), args.machineStatus.String())))
// 		mach := machine.Machine{
// 			Id:           args.box.Id,
// 			AccountId:    args.box.AccountId,
// 			CartonId:     args.box.CartonId,
// 			CartonsId:    args.box.CartonsId,
// 			Level:        args.box.Level,
// 			Name:         args.box.GetFullName(),
// 			Status:       args.machineStatus,
// 			State:        args.machineState,
// 			Image:        args.imageId,
// 			StorageType:  args.box.StorageType,
// 			Region:       args.box.Region,
// 			VMId:         args.box.InstanceId,
// 			VCPUThrottle: args.provisioner.vcpuThrottle,
// 		}
// 		fmt.Fprintf(writer, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf(" creating struct machine (%s, %s)OK", args.box.GetFullName(), args.machineStatus.String())))
// 		return mach, nil
// 	},
// 	Backward: func(ctx action.BWContext) {
// 	},
// }
