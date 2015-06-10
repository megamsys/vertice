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
package chef

import (
	log "code.google.com/p/log4go"
	"github.com/megamsys/megamd/iaas"
	"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/provisioner"
)

func Init() {
	provisioner.Register("chef", &Chef{})
}

type Chef struct {
}

func (i *Chef) CreateCommand(assembly *global.AssemblyWithComponents, id string, instance bool, act_id string) (string, error) {
	// Iaas Provider
	provider := ""
	log.Info("Chef provisioner entry")
	if instance {
		provider = "megam"
	} else {
	   // this is hack for only 0.8 release and future we implements hybrid cloud
		provider = "megam"
	}
	
	log.Info(provider)
	iaas, pdc, err1 := iaas.GetIaasProvider(provider)
	if err1 != nil {
		log.Error("Error: Iaas Provider :\n%s.", err1)
		return "", err1
	}		
	str, iaaserr := iaas.CreateMachine(pdc, assembly, act_id)
	if iaaserr != nil {
		log.Error("Error: Iaas Provider doesn't create machine:\n%s.", iaaserr)
		return "", iaaserr
	}
	return str, nil
}

func (i *Chef) DeleteCommand(assembly *global.AssemblyWithComponents, id string) (string, error) {
	// Iaas Provider
	provider := "megam"
		
	log.Info(provider)
	iaas, pdc, err1 := iaas.GetIaasProvider(provider)
	if err1 != nil {
		log.Error("Error: Iaas Provider :\n%s.", err1)
		return "", err1
	}
	log.Info(iaas)
	str, iaaserr := iaas.DeleteMachine(pdc, assembly)
	if iaaserr != nil {
		log.Error("Error: Iaas Provider doesn't delete machine:\n%s.", iaaserr)
		return "", iaaserr
	}
	log.Info(str)
	return str, nil
}
