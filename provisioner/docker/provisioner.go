
package docker

import (
 log "code.google.com/p/log4go"
 "encoding/json"
 "github.com/megamsys/megamd/global"
 "github.com/fsouza/go-dockerclient"
 "github.com/megamsys/megamd/provisioner"
 "github.com/tsuru/config"
)


func Init() {
	provisioner.Register("docker", &Docker{})
}

type Docker struct {
}

func (i *Docker) CreateCommand(assembly *global.AssemblyWithComponents, id string, instance bool, act_id string) (string, error) {
//Creates containers into the specificed  endpoint provided in the assembly.

pair_endpoint, perrscm := global.ParseKeyValuePair(assembly.Inputs, "endpoint")
   if perrscm != nil {
     log.Error("Failed to get the domain value : %s", perrscm)
   }

pair_img, perrscm := global.ParseKeyValuePair(assembly.Components[0].Inputs, "source")
    if perrscm != nil {
      log.Error("Failed to get the domain value : %s", perrscm)
    }
var endpoint string
if pair_endpoint.Value == "baremetal" {

  api_host, _:= config.GetString("swarm:host")
//  if apierr != nil {
//    return apierr
//  }
      endpoint = api_host

} else {
   endpoint = pair_endpoint.Value
 }

 client, _ := docker.NewClient(endpoint)


   config := docker.Config{Image: pair_img.Value}
 copts := docker.CreateContainerOptions{Name: pair_img.Value, Config: &config}
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

   //contt, _ := client.ListContainers(docker.ListContainersOptions{})
   return "",nil
}



func (i *Docker) DeleteCommand(assembly *global.AssemblyWithComponents, id string) (string, error) {return "",nil}
