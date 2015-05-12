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
package megam

import (
	"github.com/megamsys/megamd/iaas"
	"github.com/megamsys/megamd/global"
	"github.com/tsuru/config"
	"encoding/json"
	"bytes"
	"fmt"
	"strings"
)

func Init() {
	iaas.RegisterIaasProvider("megam", &MegamIaaS{})
}

type MegamIaaS struct{}

func (i *MegamIaaS) CreateMachine(pdc *global.PredefClouds, assembly *global.AssemblyWithComponents, act_id string) (string, error) {
  global.LOG.Info("Megam provider create entry")
  accesskey, err_accesskey := config.GetString("ACCESS_KEY")
	if err_accesskey != nil {
		return "", err_accesskey
	}
	
	secretkey, err_secretkey := config.GetString("SECRET_KEY")
	if err_secretkey != nil {
		return "", err_secretkey
	}
	
	str, err := buildCommand(assembly)
	if err != nil {
		return "", err
	}
	
	pair, perr := global.ParseKeyValuePair(assembly.Inputs, "domain")
	if perr != nil {
		global.LOG.Error("Failed to get the domain value : %s", perr)
	}
		
	knifePath, kerr := config.GetString("knife:path")
	if kerr != nil {
		return "", kerr
	}	
	
	str = str + " -c " + knifePath
	str = str + " -N " + assembly.Name + "." + pair.Value
	str = str + " -A " + accesskey
	str = str + " -K " + secretkey
	
	recipe, err_recipe := config.GetString("knife:recipe")
	if err_recipe != nil {
		return "", err_recipe
	}
	
	riakHost, err_riakHost := config.GetString("hosts:riak_host")
	if err_riakHost != nil {
		return "", err_riakHost
	}
	
	rabbitmqHost, err_rabbitmq := config.GetString("hosts:rabbitmq_host")
	if err_rabbitmq != nil {
		return "", err_rabbitmq
	}
	
	monitor, err_monitor := config.GetString("hosts:monitor_host")
	if err_monitor != nil {
		return "", err_monitor
	}
	
	kibana, err_kibana := config.GetString("hosts:kibana_host")
	if err_kibana != nil {
		return "", err_kibana
	}
	
	etcdHost, err_etcd := config.GetString("hosts:etcd_host")
	if err_etcd != nil {
		return "", err_etcd
	}

	
	str = str + " --run-list recipe[" + recipe + "]"
	attributes := &iaas.Attributes{RiakHost: riakHost, AccountID: act_id, AssemblyID: assembly.Id, RabbitMQ: rabbitmqHost, MonitorHost: monitor, KibanaHost: kibana, EtcdHost: etcdHost}
    b, aerr := json.Marshal(attributes)
    if aerr != nil {        
        return "", aerr
    }
	str = str + " --json-attributes " + string(b)
	
	return str, nil
 
}

func (i *MegamIaaS) DeleteMachine(pdc *global.PredefClouds, assembly *global.AssemblyWithComponents) (string, error) {
  
	keys, err_keys := iaas.GetAccessKeys(pdc)
     if err_keys != nil {
     	return "", err_keys
     }
     
     str, err := buildDelCommand(iaas.GetPlugins("opennebula"), pdc, "delete")
	if err != nil {
	return "", err
	 }
	//str = str + " -P " + " -y "
	pair, perr := global.ParseKeyValuePair(assembly.Components[0].Inputs, "domain")
		if perr != nil {
			global.LOG.Error("Failed to get the domain value : %s", perr)
		}
	str = str + " -N " + assembly.Name + "." + pair.Value
	str = str + " -A " + keys.AccessKey
	str = str + " -K " + keys.SecretKey

   knifePath, kerr := config.GetString("knife:path")
	if kerr != nil {
		return "", kerr
	}
	str = strings.Replace(str, " -c ", " -c "+knifePath+" ", -1)
	str = strings.Replace(str, "<node_name>", assembly.Name + "." + pair.Value, -1 )
   
    if len(pdc.Access.Zone) > 0 {
		   str = str + " --endpoint" + pdc.Access.Zone
	} else {
		return "", fmt.Errorf("Zone doesn't loaded")
	}

return str, nil	
}


func buildCommand(assembly *global.AssemblyWithComponents) (string, error) {
	var buffer bytes.Buffer
	buffer.WriteString("knife ")
	buffer.WriteString("opennebula ")
	buffer.WriteString("server ")
	buffer.WriteString("create")	
	
	templatekey := ""
	if len(assembly.Components) > 0 {
	   megamtemplatekey, err_templatekey := config.GetString("MEGAM_TEMPLATE_NAME")
		if err_templatekey != nil {
			return "", err_templatekey
		}	
		templatekey = megamtemplatekey
	} else {
		atype := make([]string, 3)
		atype = strings.Split(assembly.ToscaType, ".")
    	templatekey = "megam_" + atype[2]
	}
	
	if len(templatekey) > 0 {
		buffer.WriteString(" --template-name " + templatekey)
	} else {
		return "", fmt.Errorf("Template doesn't loaded")
	}

    sshuserkey, err_sshuserkey := config.GetString("SSH_USER")
	if err_sshuserkey != nil {
		return "", err_sshuserkey
	}
	if len(sshuserkey) > 0 {
		buffer.WriteString(" -x " + sshuserkey)
	} else {
		return "", fmt.Errorf("Ssh user value doesn't loaded")
	}

	identityfilekey, err_identityfilekey := config.GetString("IDENTITY_FILE")
	if err_identityfilekey != nil {
		return "", err_identityfilekey
	}
	if len(identityfilekey) > 0 {
		buffer.WriteString(" --identity-file " + identityfilekey)
	} else {
		return "", fmt.Errorf("Identity file doesn't loaded")
	}

	zonekey, err_zonekey := config.GetString("ZONE")
	if err_zonekey != nil {
		return "", err_zonekey
	}
	if len(zonekey) > 0 {
		buffer.WriteString(" --endpoint " + zonekey)
	} else {
		return "", fmt.Errorf("Zone doesn't loaded")
	}

	return buffer.String(), nil
}

func buildDelCommand(plugin *iaas.Plugins, pdc *global.PredefClouds, command string) (string, error) {
	var buffer bytes.Buffer
	if len(plugin.Tool) > 0 {
		buffer.WriteString(plugin.Tool)
	} else {
		return "", fmt.Errorf("Plugin tool doesn't loaded")
	}
	if command == "delete" {
		if len(plugin.Command.Delete) > 0 {
			buffer.WriteString(" " + plugin.Command.Delete)
		} else {
			return "", fmt.Errorf("Plugin commands doesn't loaded")
		}
	}
	return buffer.String(), nil 
	
}	