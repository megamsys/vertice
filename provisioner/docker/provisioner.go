
package docker

import (
 log "code.google.com/p/log4go"
 "encoding/json"
//	"github.com/tsuru/config"
 //"github.com/megamsys/libgo/db"
 "github.com/megamsys/megamd/global"
 //"github.com/megamsys/gulp/policies"
 //"github.com/megamsys/gulp/app"
 "github.com/fsouza/go-dockerclient"
 "github.com/megamsys/megamd/provisioner"
 "strings"
 //"fmt"
 "bytes"
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


pair_img, perrscm := global.ParseKeyValuePair(assembly.Inputs, "source")
    if perrscm != nil {
      log.Error("Failed to get the domain value : %s", perrscm)
    }


 endpoint := pair_endpoint.Value
 client, _ := docker.NewClient(endpoint)



   var buf bytes.Buffer
   source := strings.Split(pair_img.Value, ":")
   var tag string

   if len(source) > 1 {
        tag = source[1]

    } else {
      tag = ""
   }


 opts := docker.PullImageOptions{
                     Repository:   source[0],
                     Registry:     "",
                     Tag:          tag,
                     OutputStream: &buf,
                    }
 pullerr := client.PullImage(opts, docker.AuthConfiguration{})
 if pullerr != nil {
      log.Error(pullerr)
   }

   config := docker.Config{Image: "gomegam/megamgateway:0.5.0"}
 copts := docker.CreateContainerOptions{Name: "redis", Config: &config}
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
