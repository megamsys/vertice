package profitbricks

import (
   "github.com/megamsys/megamd/iaas"
   "bytes"
   "strings"
   "fmt"
   "github.com/megamsys/megamd/provisioner"
)

func Init() {
	iaas.RegisterIaasProvider("profitbricks", &ProfitBricksIaaS{})
}


type ProfitBricksIaaS struct{}

func (i *ProfitBricksIaaS) DeleteMachine(string) error {
	
	return nil
}

func (i *ProfitBricksIaaS) CreateMachine(pdc *iaas.PredefClouds, assembly *provisioner.AssemblyResult)  (string, error) {
	  keys, err_keys := iaas.GetAccessKeys(pdc)
	  if err_keys != nil {
	  	return "", err_keys
	  }
	  
	  str, err := buildCommand(iaas.GetPlugins("profitbricks"), pdc, "create")
	   if err != nil {
	   	return "", err
	   }   
	   str = str + " -N " + assembly.Name+"."+assembly.Components[0].Inputs.Domain
	   str = str + " -A " + keys.AccessKey
	   str = str + " -K " + keys.SecretKey
	//strings.Replace(str,"-c","-c "+assembly.Name+"."+assembly.Components[0].Inputs.Domain,-1)  
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
	     buffer.WriteString(" --image-name "+pdc.Spec.Image)
	} 
	
	if len(pdc.Spec.TenantID) > 0 { 
	     buffer.WriteString(" --data-center "+pdc.Spec.TenantID)
	} else {
		return "", fmt.Errorf("Data center doesn't loaded")
	}
	
	if len(pdc.Access.Sshuser) > 0 { 
	     buffer.WriteString(" -x "+pdc.Access.Sshuser)
	} else {
		return "", fmt.Errorf("Ssh user value doesn't loaded")
	}
	
	s := strings.Split(pdc.Spec.Flavor, ",")
	
	if len(pdc.Spec.Flavor) > 0 {
		 cpus := strings.Split(s[0], "=")  
	     buffer.WriteString(" --cpus "+cpus[1])
	} else {
		return "", fmt.Errorf("cpus doesn't loaded")
	}
	
	if len(pdc.Spec.Flavor) > 0 {  
		hddsize := strings.Split(s[2], "=") 
	     buffer.WriteString(" --hdd-size "+hddsize[1])
	} else {
		return "", fmt.Errorf("hdd-size doesn't loaded")
	}
	
	if len(pdc.Spec.Flavor) > 0 {  
		ram := strings.Split(s[1], "=") 
	     buffer.WriteString(" --ram "+ram[1])
	} else {
		return "", fmt.Errorf("ram doesn't loaded")
	}
	
	if len(pdc.Access.IdentityFile) > 0 { 
	    ifile, err := iaas.GetIdentityFileLocation(pdc.Access.IdentityFile)
	    if err != nil {
	    	return "", fmt.Errorf("Identity file doesn't loaded")
	    }
	    buffer.WriteString(" --identity-file "+ifile+".key")
	} 
	//else {
	//	return "", fmt.Errorf("Identity file doesn't loaded")
	//}
	
	if len(pdc.Access.SshpubLocation) > 0 { 
		 ifile, err := iaas.GetIdentityFileLocation(pdc.Access.SshpubLocation)
	    if err != nil {
	    	return "", fmt.Errorf("SshpubLocation file doesn't loaded")
	    }
	     buffer.WriteString(" --public-key-file "+ifile+".pub")
	} else {
		return "", fmt.Errorf("SshpubLocation doesn't loaded")
	}
	
	if len(pdc.Spec.Groups) > 0 { 
	     buffer.WriteString(" -S "+pdc.Spec.Groups)
	} else {
		return "", fmt.Errorf("Groups doesn't loaded")
	}
	
	if len(pdc.Access.Sshkey) > 0 { 
	     buffer.WriteString(" --image-password "+pdc.Access.Sshkey)
	} else {
		return "", fmt.Errorf("Image password doesn't loaded")
	}
	
	return buffer.String(), nil
}



