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
package google

import (
	"bytes"
	"encoding/json"
	log "code.google.com/p/log4go"
	"fmt"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/iaas"
	"github.com/megamsys/megamd/provisioner"
	"github.com/tsuru/config"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func Init() {
	iaas.RegisterIaasProvider("google", &GoogleIaaS{})
}

type GoogleIaaS struct{}

type GoogleCredential struct {
	AuthorizationURI   string `json:"authorization_uri"`
	TokenCredentailURI string `json:"token_credential_uri"`
	Scope              string `json:"scope"`
	RedirectURI        string `json:"redirect_uri"`
	ClientID           string `json:"client_id"`
	ClientSecret       string `json:"client_secret"`
	AccessToken        string `json:"access_token"`
	ExpiresIN          int64 `json:"expires_in"`
	RefreshToken       string `json:"refresh_token"`
	Project            string `json:"project"`
}

func (i *GoogleIaaS) DeleteMachine(pdc *global.PredefClouds, assembly  *provisioner.AssemblyResult) (string, error) {

	keys, err_keys := iaas.GetAccessKeys(pdc)
     if err_keys != nil {
     	return "", err_keys
     }
     
     str, err := buildDelCommand(iaas.GetPlugins("google"), pdc, "delete")
	if err != nil {
	return "", err
	 }
	//str = str + " -P " + " -y "
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


func (i *GoogleIaaS) CreateMachine(pdc *global.PredefClouds, assembly *provisioner.AssemblyResult) (string, error) {
	cre, derr := downloadCredentials(pdc)
	if derr != nil {
		return "", derr
	}
	log.Info("================after download")
	str, err := buildCommand(iaas.GetPlugins("google"), pdc, "create")
	if err != nil {
		return "", err
	}
	str = str + " -N " + assembly.Name + "." + assembly.Components[0].Inputs.Domain
	
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
	attributes := &iaas.Attributes{RiakHost: riakHost, AccountID: pdc.Accounts_id, AssemblyID: assembly.Id, RabbitMQ: rabbitmqHost, MonitorHost: monitor, KibanaHost: kibana, EtcdHost: etcdHost}
    b, aerr := json.Marshal(attributes)
    if aerr != nil {
        fmt.Println(aerr)
        return "", aerr
    }
	str = str + " --json-attributes " + string(b)
	str = strings.Replace(str, "create -f", "create "+assembly.Name+"."+assembly.Components[0].Inputs.Domain+" -f "+cre, -1)
	knifePath, kerr := config.GetString("knife:path")
	if kerr != nil {
		return "", kerr
	}
	str = strings.Replace(str, " -c ", " -c "+knifePath+" ", -1)
	return str, nil
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
		buffer.WriteString(" -n " + pdc.Spec.Groups)
	} else {
		return "", fmt.Errorf("Groups doesn't loaded")
	}

	if len(pdc.Spec.Image) > 0 {
		buffer.WriteString(" -I " + pdc.Spec.Image)
	} else {
		return "", fmt.Errorf("Image doesn't loaded")
	}

	if len(pdc.Spec.Flavor) > 0 {
		buffer.WriteString(" -m " + pdc.Spec.Flavor)
	} else {
		return "", fmt.Errorf("Flavor doesn't loaded")
	}

	if len(pdc.Access.Sshuser) > 0 {
		buffer.WriteString(" -x " + pdc.Access.Sshuser)
	} else {
		return "", fmt.Errorf("Ssh user value doesn't loaded")
	}

	if len(pdc.Access.IdentityFile) > 0 {
		buffer.WriteString(" --identity-file " + pdc.Access.IdentityFile)
	} else {
		return "", fmt.Errorf("Identity file doesn't loaded")
	}

	if len(pdc.Access.Zone) > 0 {
		buffer.WriteString(" -Z " + pdc.Access.Zone)
	} else {
		return "", fmt.Errorf("Zone doesn't loaded")
	}

	return buffer.String(), nil
}

func downloadCredentials(pdc *global.PredefClouds) (string, error) {
	s := make([]string, 2)
	s = strings.Split(pdc.Access.VaultLocation, "_")
	email, name := s[0], s[1]
	cloudkeysBucket, ckberr := config.GetString("buckets:CLOUDKEYS")
	if ckberr != nil {
		return "", ckberr
	}

	megam_home, ckberr := config.GetString("MEGAM_HOME")
	if ckberr != nil {
		return "", ckberr
	}
	basePath := megam_home + cloudkeysBucket
	cloudaccesskeysBucket, cakberr := config.GetString("buckets:CLOUDACCESSKEYS")
	if cakberr != nil {
		return "", cakberr
	}
	conn, err := db.Conn(cloudaccesskeysBucket)
	keys := &GoogleCredential{}
	if err != nil {
		return "", err
	}
    log.Info("+++++++++++++keys1+++++++++++++++++")
    log.Info(conn)
    log.Info(keys)
    log.Info(pdc.Access.VaultLocation)
	ferr := conn.FetchStruct(pdc.Access.VaultLocation, keys)
	if ferr != nil {
		log.Info(ferr)
		return "", ferr
	}
    log.Info("+++++++++++++keys2+++++++++++++++++")
    log.Info(keys)
	b, errk := json.Marshal(keys)
	if errk != nil {
		return "", errk
	}

	dir := path.Join(basePath, email)
	filePath := path.Join(dir, name+".json")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", dir)

		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return "", errm
		}
		// open output file
		_, err := os.Create(filePath)
		if err != nil {
			return "", err
		}
	}
	errf := ioutil.WriteFile(filePath, []byte(string(b)), 0644)
	if errf != nil {
		return "", errf
	}

	return cloudkeysBucket + "/" + email + "/" + name + ".json", nil
}
