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
package app

import (
	"github.com/megamsys/libgo/db"
	log "code.google.com/p/log4go"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/megamd/provisioner"
//	"encoding/json"
	"github.com/tsuru/config"
	"github.com/megamsys/megamd/global"
)


func GetPredefClouds(host string) (*global.PredefClouds, error) {
	pdc := &global.PredefClouds{}

	predefBucket, perr := config.GetString("riak:predefclouds")
	if perr != nil {
		return pdc, perr
	}
	conn, err := db.Conn(predefBucket)

	if err != nil {
		return pdc, err
	}

	ferr := conn.FetchStruct(host, pdc)
	if ferr != nil {
		return pdc, ferr
	}
	return pdc, nil
}

func LaunchApp(asm *global.AssemblyWithComponents, id string, act_id string) error {
	log.Debug("Launch App entry")
	if len(asm.Components) > 0 {
		LauncherHelper(asm, id, false, act_id)
	} else {
    LauncherHelper(asm, id, true, act_id)
          }
    return nil
}

func LauncherHelper(asm *global.AssemblyWithComponents, id string, instance bool, act_id string) error {

	pair_host, perrscm := global.ParseKeyValuePair(asm.Inputs, "host")
	if perrscm != nil {
	  log.Error("Failed to get the host value : %s", perrscm)
	}


if pair_host.Value == "docker" {

	log.Debug("Docker provisioner entry")
	        		// Provisioner
		            p, err := provisioner.GetProvisioner("docker")
		            if err != nil {
		                return err
		             }
		            log.Info("Provisioner: %v", p)
		             _, perr := p.CreateCommand(asm, id, instance, act_id)
		            if perr != nil {
		               return perr
		             }
}

if pair_host.Value == "chef" {

	p, err := provisioner.GetProvisioner("chef")
	if err != nil {
	    return err
	}

	str, perr := p.CreateCommand(asm, id, instance, act_id)
    if perr != nil {
	  return perr
	}
	asm.Command = str
	actions := []*action.Action{&launchedApp}

	pipeline := action.NewPipeline(actions...)
	aerr := pipeline.Execute(asm)
	if aerr != nil {
         return aerr
     }
}
return nil
 }

func DeleteApp(asm *global.AssemblyWithComponents, id string) error {
       log.Debug("Delete App entry")

				pair_host, perrscm := global.ParseKeyValuePair(asm.Inputs, "host")
				if perrscm != nil {
				  log.Error("Failed to get the host value : %s", perrscm)
				}

				if pair_host.Value == "docker" {

					log.Debug("Docker provisioner entry")
											// Provisioner
												p, err := provisioner.GetProvisioner("docker")
												if err != nil {
														return err
												}
												log.Info("Provisioner: %v", p)
												_, perr := p.DeleteCommand(asm, id)
												if perr != nil {
													return perr
												}



				}
        if pair_host.Value == "chef" {

	    p, err := provisioner.GetProvisioner("chef")
	   if err != nil {
	         return err
	   }
	   str, perr := p.DeleteCommand(asm, id)
	   if perr != nil {
	     return perr
	    }
	    asm.Command = str
	    actions := []*action.Action{&updateStatus}
	   pipeline := action.NewPipeline(actions...)
	   aerr := pipeline.Execute(asm)
	   if aerr != nil {
		    return aerr
	     }
	  	}
		return nil
}
