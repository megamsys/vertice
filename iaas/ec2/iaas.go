package ec2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/megamsys/megamd/iaas"
	"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/provisioner"
	"github.com/tsuru/config"
	"strings"
)

func Init() {
	iaas.RegisterIaasProvider("ec2", &EC2IaaS{})
}

type EC2IaaS struct{}

func (i *EC2IaaS) DeleteMachine(pdc *global.PredefClouds, assembly *provisioner.AssemblyResult) (string, error) {
     keys, err_keys := iaas.GetAccessKeys(pdc)
     if err_keys != nil {
     	return "", err_keys
     }
     
     str, err := buildDelCommand(iaas.GetPlugins("ec2"), pdc, "delete")
	if err != nil {
	return "", err
	 }
	str = str + " -P " + " -y "
	str = str + " -N " + assembly.Name + "." + assembly.Components[0].Inputs.Domain
	str = str + " -A " + keys.AccessKey
	str = str + " -K " + keys.SecretKey
	

   knifePath, kerr := config.GetString("knife:path")
	if kerr != nil {
		return "", kerr
	}
	str = strings.Replace(str, " -c ", " -c "+knifePath+" ", -1)
	str = strings.Replace(str, "<node_name>", assembly.Name + "." + assembly.Components[0].Inputs.Domain, -1 )
   

return str, nil	
}

func (i *EC2IaaS) CreateMachine(pdc *global.PredefClouds, assembly *provisioner.AssemblyResult) (string, error) {
	keys, err_keys := iaas.GetAccessKeys(pdc)
	if err_keys != nil {
		return "", err_keys
	}

	str, err := buildCommand(iaas.GetPlugins("ec2"), pdc, "create")
	if err != nil {
		return "", err
	}
	str = str + " -N " + assembly.Name + "." + assembly.Components[0].Inputs.Domain
	str = str + " -A " + keys.AccessKey
	str = str + " -K " + keys.SecretKey
	
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

	recipe, err_recipe := config.GetString("knife:recipe")
	if err_recipe != nil {
		return "", err_recipe
	}

	str = str + " --run-list recipe[" + recipe + "]" 
	attributes := &iaas.Attributes{RiakHost: riakHost, AccountID: pdc.Accounts_id, AssemblyID: assembly.Id, RabbitMQ: rabbitmqHost, MonitorHost: monitor, KibanaHost: kibana, EtcdHost: etcdHost}
	b, aerr := json.Marshal(attributes)
	if aerr != nil {
		
		fmt.Println(aerr)
		return "", aerr
	}
	str = str + " --json-attributes " + string(b)

	//strings.Replace(str,"-c","-c "+assembly.Name+"."+assembly.Components[0].Inputs.Domain,-1)
	knifePath, kerr := config.GetString("knife:path")
	if kerr != nil {
		return "", kerr
	}
	str = strings.Replace(str, " -c ", " -c "+knifePath+" ", -1)
	return str, nil
	
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
	
	
	if len(pdc.Access.Region) > 0 {
		buffer.WriteString(" --region " + pdc.Access.Region)
	} else {
		return "", fmt.Errorf("Region does not get loaded")
	}
	
	return buffer.String(), nil 
	
}	

func buildCommand(plugin *iaas.Plugins, pdc *global.PredefClouds, command string) (string, error) {
	var buffer bytes.Buffer
	if len(plugin.Tool) > 0 {
		buffer.WriteString(plugin.Tool)
	} else {
		return "", fmt.Errorf("Plugin tool doesn't loaded")
	}

	if command == "create" {
		if len(plugin.Command.Create) > 0 {
			buffer.WriteString(" " + plugin.Command.Create)
		} else {
			return "", fmt.Errorf("Plugin commands doesn't loaded")
		}
	}
	
	

	if len(pdc.Spec.Groups) > 0 {
		buffer.WriteString(" -G " + pdc.Spec.Groups)
	} else {
		return "", fmt.Errorf("Groups doesn't loaded")
	}

	if len(pdc.Spec.Image) > 0 {
		buffer.WriteString(" -I " + pdc.Spec.Image)
	} else {
		return "", fmt.Errorf("Image doesn't loaded")
	}

	if len(pdc.Spec.Flavor) > 0 {
		buffer.WriteString(" -f " + pdc.Spec.Flavor)
	} else {
		return "", fmt.Errorf("Flavor doesn't loaded")
	}

	if len(pdc.Access.Sshkey) > 0 {
		buffer.WriteString(" -S " + pdc.Access.Sshkey)
	} else {
		return "", fmt.Errorf("Ssh key value doesn't loaded")
	}

	if len(pdc.Access.Sshuser) > 0 {
		buffer.WriteString(" -x " + pdc.Access.Sshuser)
	} else {
		return "", fmt.Errorf("Ssh user value doesn't loaded")
	}

	if len(pdc.Access.IdentityFile) > 0 {
		ifile, err := iaas.GetIdentityFileLocation(pdc.Access.IdentityFile)
		if err != nil {
			return "", fmt.Errorf("Identity file doesn't loaded")
		}
		buffer.WriteString(" --identity-file " + ifile + ".key")
	} else {
		return "", fmt.Errorf("Identity file doesn't loaded")
	}

	if len(pdc.Access.Region) > 0 {
		buffer.WriteString(" --region " + pdc.Access.Region)
	} else {
		return "", fmt.Errorf("Zone doesn't loaded")
	}
	return buffer.String(), nil
}
