package opennebula

import (
   "github.com/megamsys/megamd/iaas"
   "bytes"
   "fmt"
   "strings"
   "github.com/tsuru/config"
   "github.com/megamsys/megamd/provisioner"
)

func Init() {
	iaas.RegisterIaasProvider("opennebula", &OpenNebulaIaaS{})
}


type OpenNebulaIaaS struct{}

func (i *OpenNebulaIaaS) DeleteMachine(string) error {
	
	return nil
}

func (i *OpenNebulaIaaS) CreateMachine(pdc *iaas.PredefClouds, assembly *provisioner.AssemblyResult)  (string, error) {
	  keys, err_keys := iaas.GetAccessKeys(pdc)
	  if err_keys != nil {
	  	return "", err_keys
	  }
	  
	  str, err := buildCommand(iaas.GetPlugins("opennebula"), pdc, "create")
	   if err != nil {
	   	return "", err
	   }   
	   str = str + " -N " + assembly.Name+"."+assembly.Components[0].Inputs.Domain
	   str = str + " -A " + keys.AccessKey
	   str = str + " -K " + keys.SecretKey
	//strings.Replace(str,"-c","-c "+assembly.Name+"."+assembly.Components[0].Inputs.Domain,-1)  
	 knifePath, kerr := config.GetString("knife:path")
	  if kerr != nil {
	  	return "", kerr
	  }
	str = strings.Replace(str,"-c","-c "+knifePath,-1)  
	return str, nil
}

func buildCommand(plugin *iaas.Plugins, pdc *iaas.PredefClouds, command string) (string, error) {
	var buffer bytes.Buffer
	if len(plugin.Tool) > 0 { 
	     buffer.WriteString(plugin.Tool)
	} else {
		return "", fmt.Errorf("Plugin tool doesn't loaded")
	}
	
	if command == "create" {
	  if len(plugin.Command.Create) > 0 { 
	       buffer.WriteString(" "+plugin.Command.Create)
	   } else {
	    	return "", fmt.Errorf("Plugin commands doesn't loaded")
	   }
	}

	if len(pdc.Spec.Image) > 0 { 
	     buffer.WriteString(" --template-name "+pdc.Spec.Image)
	} else {
		return "", fmt.Errorf("Image doesn't loaded")
	}
	
	if len(pdc.Access.Sshuser) > 0 { 
	     buffer.WriteString(" -x "+pdc.Access.Sshuser)
	} else {
		return "", fmt.Errorf("Ssh user value doesn't loaded")
	}
	
	if len(pdc.Access.IdentityFile) > 0 { 
	    ifile, err := iaas.GetIdentityFileLocation(pdc.Access.IdentityFile)
	    if err != nil {
	    	return "", fmt.Errorf("Identity file doesn't loaded")
	    }
	    buffer.WriteString(" --identity-file "+ifile)
	} else {
		return "", fmt.Errorf("Identity file doesn't loaded")
	}
	
	if len(pdc.Access.Zone) > 0 { 
	     buffer.WriteString(" --endpoint "+pdc.Access.Zone)
	} else {
		return "", fmt.Errorf("Zone doesn't loaded")
	}
	
	return buffer.String(), nil
}



