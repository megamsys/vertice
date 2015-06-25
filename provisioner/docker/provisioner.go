package docker

import (
	"encoding/json"
	log "code.google.com/p/log4go"
	"github.com/megamsys/megamd/global"
	"github.com/fsouza/go-dockerclient"
	"github.com/megamsys/megamd/provisioner"
	"github.com/tsuru/config"
	"github.com/megamsys/seru/cmd/seru"
	"fmt"

)

func Init() {
	provisioner.Register("docker", &Docker{})
}

type Docker struct {
}

const BAREMETAL = "baremetal"

func (i *Docker) Create(assembly *global.AssemblyWithComponents, id string, instance bool, act_id string) (string, error) {
	//Creates containers into the specificed  endpoint provided in the assembly.
	log.Info("%q", assembly)
	pair_endpoint, perrscm := global.ParseKeyValuePair(assembly.Inputs, "endpoint")
	if perrscm != nil {
		log.Error("Failed to get the endpoint value : %s", perrscm)
	}

	pair_img, perrscm := global.ParseKeyValuePair(assembly.Components[0].Inputs, "source")
	if perrscm != nil {
		log.Error("Failed to get the image value : %s", perrscm)
	}
	
	pair_domain, perrdomain := global.ParseKeyValuePair(assembly.Components[0].Inputs, "domain")
	if perrdomain != nil {
		log.Error("Failed to get the image value : %s", perrdomain)
	}
	
	var endpoint string
	if pair_endpoint.Value == BAREMETAL {

		api_host, _ := config.GetString("swarm:host")
		endpoint = api_host

	} else {
		endpoint = pair_endpoint.Value
	}

	client, _ := docker.NewClient(endpoint)

	config := docker.Config{Image: pair_img.Value}
	copts := docker.CreateContainerOptions{Name: fmt.Sprint(assembly.Components[0].Name, ".", pair_domain.Value), Config: &config}

	container, conerr := client.CreateContainer(copts)
	if conerr != nil {
		log.Error(conerr)
	}
    
	cont := &docker.Container{}
	mapP, _ := json.Marshal(container)
	json.Unmarshal([]byte(string(mapP)), cont)	
 
	serr := client.StartContainer(cont.ID, &docker.HostConfig{})
	if serr != nil {
		log.Error(serr)
	}	
	
	herr := setHostName(client, cont)
	if herr != nil {
		log.Error(herr)
	}	
	
	return "", nil
}

func setHostName(client *docker.Client, c *docker.Container) error {
	container, _ := client.InspectContainer(c.ID)
	cont := &docker.Container{}
	mapP, _ := json.Marshal(container)
	json.Unmarshal([]byte(string(mapP)), cont)
	
	container_network := &docker.NetworkSettings{}
	mapN, _ := json.Marshal(cont.NetworkSettings)
	json.Unmarshal([]byte(string(mapN)), container_network)
	fmt.Println(container_network.IPAddress)
	
	seru := &main.NewSubdomain{}
	
	fmt.Println(seru.Info())
	
	return nil
}


func (i *Docker) Delete(assembly *global.AssemblyWithComponents, id string) (string, error) {
	return "", nil
}
