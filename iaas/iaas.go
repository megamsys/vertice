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
package iaas

import (
	log "code.google.com/p/log4go"
	"errors"
	"fmt"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/megamd/global"
	"github.com/tsuru/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Every Tsuru IaaS must implement this interface.
type IaaS interface {
	// Called when tsuru is creating a Machine.
	CreateMachine(*global.PredefClouds, *global.AssemblyWithComponents, string) (string, error)

	// Called when tsuru is destroying a Machine.
	DeleteMachine(*global.PredefClouds, *global.AssemblyWithComponents) (string, error)
}

const (
	defaultYAMPath = "conf/commands.yaml"
	PREDEFCLOUDS   = "predefclouds"
	CLOUDKEYS      = "cloudkeys"
	SSH_FILES      = "ssh_files"
)

type Attributes struct {
	RiakHost    string `json:"riak_host"`
	AccountID   string `json:"accounts_id"`
	AssemblyID  string `json:"assembly_id"`
	RabbitMQ    string `json:"rabbitmq_host"`
	MonitorHost string `json:"monitor_host"`
	KibanaHost  string `json:"kibana_host"`
}

type Plugins struct {
	Tool    string
	Command *Commands
}

type Commands struct {
	Create string
	Delete string
	List   string
	Data   string
}

var iaasProviders = make(map[string]IaaS)

func RegisterIaasProvider(name string, iaas IaaS) {
	iaasProviders[name] = iaas
}

func GetIaasProvider(name string) (IaaS, *global.PredefClouds, error) {
	pdc := &global.PredefClouds{}
	err := errors.New("")
	pdc_type := ""
	if name == "megam" {
		pdc_type = name
	} else {
		pdc, err = getProviderName(name)
		if err != nil {
			return nil, nil, fmt.Errorf("Error: Riak didn't cooperate:\n%s.", err)
		}
	}

	provider, ok := iaasProviders[pdc_type]
	if !ok {
		return nil, nil, fmt.Errorf("IaaS provider not registered")
	}
	return provider, pdc, nil
	//return nil
}

func getProviderName(host string) (*global.PredefClouds, error) {
	pdc := &global.PredefClouds{}

	conn, err := db.Conn(PREDEFCLOUDS)

	if err != nil {
		return pdc, err
	}

	ferr := conn.FetchStruct(host, pdc)
	if ferr != nil {
		return pdc, ferr
	}

	sshkeyerr := downloadSshFiles(pdc, "key", 0600)
	if sshkeyerr != nil {
		return pdc, sshkeyerr
	}
	sshpuberr := downloadSshFiles(pdc, "pub", 0644)
	if sshpuberr != nil {
		return pdc, sshpuberr
	}

	return pdc, nil
}

func GetPlugins(cloud string) *Plugins {
	p, _ := filepath.Abs(defaultYAMLPath)
	log.Info(fmt.Errorf("Conf: %s", p))

	data, err := ioutil.ReadFile(p)

	if err != nil {
		log.Info("error: %v", err)
	}

	m := make(map[interface{}]Plugins)

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Info("error: %v", err)
	}
	for key, value := range m {
		if key == cloud {
			return &value
		}
	}
	return &Plugins{}
}

func GetIdentityFileLocation(file string) (string, error) {
	s := make([]string, 2)
	s = strings.Split(file, "_")
	email, name := s[0], s[1]

	megam_home, err := config.GetString("megam:home")
	if err != nil {
		return "", err
	}

	return megam_home + CLOUDKEYS + "/" + email + "/" + name, nil
}

type SshFile struct {
	data string
}

func downloadSshFiles(pdc *global.PredefClouds, keyvalue string, permission os.FileMode) error {
	sa := make([]string, 2)
	sa = strings.Split(pdc.Access.IdentityFile, "_")
	email, name := sa[0], sa[1]
	ssh := &db.SshObject{}

	conn, err := db.Conn(SSH_FILES)
	if err != nil {
		return err
	}

	ferr := conn.FetchObject(pdc.Access.IdentityFile+"_"+keyvalue, ssh)
	if ferr != nil {
		return ferr
	}

	megam_home, ckberr := config.GetString("megam:home")
	if ckberr != nil {
		return ckberr
	}

	basePath := megam_home + CLOUDKEYS

	dir := path.Join(basePath, email)
	filePath := path.Join(dir, name+"."+keyvalue)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", dir)

		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return errm
		}
		// open output file
		_, err := os.Create(filePath)
		if err != nil {
			return err
		}
	}
	errf := ioutil.WriteFile(filePath, []byte(ssh.Data), permission)
	if errf != nil {
		return errf
	}
	return nil
}
