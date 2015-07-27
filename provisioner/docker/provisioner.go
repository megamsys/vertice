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
package docker

import (
	"encoding/json"
	"fmt"
	"strings"
	"net"
	"net/http"
	"io/ioutil"
	"bytes"
	log "code.google.com/p/log4go"
	"github.com/fsouza/go-dockerclient"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/provisioner"
	"github.com/megamsys/seru/cmd"
	"github.com/megamsys/seru/cmd/seru"
	"github.com/tsuru/config"
	"github.com/megamsys/megamd/utils"
)

/*
*
* Registers docker as provisioner in provisioner interface.
*
*/

func Init() {
	provisioner.Register("docker", &Docker{})
}

type Docker struct{}

const BAREMETAL = "baremetal"

/*
* Create provisioner is called to launch docker containers by
* talking to swarm cluster. Common provisioner for both
* Baremetal and VM-docker launch. Specify endpoint
* Swarm Host IP is added into the conf file.
*
*/

func (i *Docker) Create(assembly *global.AssemblyWithComponents, id string, instance bool, act_id string) (string, error) {
	log.Info("%q", assembly)
	pair_endpoint, perrscm := global.ParseKeyValuePair(assembly.Inputs, "endpoint")
	if perrscm != nil {
		log.Error("Failed to get the endpoint value : %s", perrscm)
		return "", perrscm
	}

	pair_img, perrscm := global.ParseKeyValuePair(assembly.Components[0].Inputs, "source")
	if perrscm != nil {
		log.Error("Failed to get the image value : %s", perrscm)
		return "", perrscm
	}

	pair_domain, perrdomain := global.ParseKeyValuePair(assembly.Components[0].Inputs, "domain")
	if perrdomain != nil {
		log.Error("Failed to get the image value : %s", perrdomain)
		return "", perrdomain
	}

	var endpoint string
	if pair_endpoint.Value == BAREMETAL {

		/*
		 * swarm host is obtained from conf file. Swarm host is considered
		 * only when the 'endpoint' is baremetal in the Component JSON
		 */
		api_host, _ := config.GetString("swarm:host")
		endpoint = api_host
	} else {
		endpoint = pair_endpoint.Value
	}
	/*
	 * Docker API client to connect to swarm. Swarm supports all docker API endpoints
	 */
	client, _ := docker.NewClient(endpoint)

	opts := docker.PullImageOptions{
		Repository: pair_img.Value,
	}
	pullerr := client.PullImage(opts, docker.AuthConfiguration{})
	if pullerr != nil {
		log.Error(pullerr)
	}

	/*
	 * Inspect image to get the default internal port to ExposedPorts
	 * the running internal port to external port
	 *
	 */

	img, err := client.InspectImage(pair_img.Value)
	if err != nil {
		log.Error("Inspect image failed : %s", err)
	}
	InspectImg := &docker.Image{}
	mapFP, _ := json.Marshal(img)
	json.Unmarshal([]byte(string(mapFP)), InspectImg)
	conf := InspectImg.Config

	var Iport string

	for k, _ := range conf.ExposedPorts {
		port := strings.Split(string(k), "/")
		Iport = port[0]

	}

	
	dconfig := docker.Config{Image: pair_img.Value, NetworkDisabled: true}
	copts := docker.CreateContainerOptions{Name: fmt.Sprint(assembly.Components[0].Name, ".", pair_domain.Value), Config: &dconfig}

	/*
	 * Creation of the container with copts.
	 */

	container, conerr := client.CreateContainer(copts)
	if conerr != nil {
		log.Error("Container creation failed : %s", conerr)
		return "", conerr
	}

	cont := &docker.Container{}
	mapP, _ := json.Marshal(container)
	json.Unmarshal([]byte(string(mapP)), cont)

	/*
	 * hostConfig{} struct for portbindings - to expose visible ports
	 *  Also for specfying the container configurations (memory, cpuquota etc)
	 */

	hostConfig := docker.HostConfig{}
	log.Info(Iport)
	//hostConfig.PortBindings = map[docker.Port][]docker.PortBinding{
	//	docker.Port(Iport + "/tcp"): {{HostIP: "", HostPort: ""}},
	//}	
	
	/*
	 *   Starting container once the container is created - container ID &
	 *   hostConfig is proivided to start the container.
	 *
	 */
	serr := client.StartContainer(cont.ID, &hostConfig)
	if serr != nil {
		log.Error("Start container was failed : %s", serr)
		return "", serr
	}
	
	/*
	* generate the ip 
	*/
	subnetip, _ := config.GetString("swarm:subnet")
	_, subnet, _ := net.ParseCIDR(subnetip)
	ip, pos, iperr := utils.IPRequest(*subnet)
	if iperr != nil {
		log.Error("Ip generation was failed : %s", iperr)
		return "", iperr
	}
	
	/*
	* configure ip to container
	*/
	ipperr := postnetwork(cont.ID, ip.String())
	log.Info(ipperr)
	
	uerr := updateIndex(ip.String(), pos)
	if uerr != nil {
		log.Error("Ip index update was failed : %s", uerr)
	}

	/*
	 * Inspect API is called to fetch the data about the launched container
	 *
	 */
	inscontainer, _ := client.InspectContainer(cont.ID)
	contain := &docker.Container{}
	mapC, _ := json.Marshal(inscontainer)
	json.Unmarshal([]byte(string(mapC)), contain)

	container_network := &docker.NetworkSettings{}
	mapN, _ := json.Marshal(contain.NetworkSettings)
	json.Unmarshal([]byte(string(mapN)), container_network)

	configs := &docker.Config{}
	mapPort, _ := json.Marshal(contain.Config)
	json.Unmarshal([]byte(string(mapPort)), configs)
   
	var port string

	for k, _ := range container_network.Ports {
		porti := strings.Split(string(k), "/")
		port = porti[0]
	}
	fmt.Println(port)

	updatecomponent(assembly, ip.String(), cont.ID, port)

	herr := setHostName(fmt.Sprint(assembly.Components[0].Name, ".", pair_domain.Value), ip.String())
	if herr != nil {
		log.Error("Failed to set the host name : %s", herr)
		return "", herr
	}
	

	return "", nil
}

/*
* Register a hostname on AWS Route53 using megam seru -
*        www.github.com/megamsys/seru
*/
func setHostName(name string, ip string) error {

	s := make([]string, 4)
	s = strings.Split(name, ".")

	accesskey, _ := config.GetString("aws:accesskey")
	secretkey, _ := config.GetString("aws:secretkey")

	seru := &main.NewSubdomain{
		Accesskey: accesskey,
		Secretid:  secretkey,
		Domain:    fmt.Sprint(s[1], ".", s[2], "."),
		Subdomain: s[0],
		Ip:        ip,
	}

	seruerr := seru.ApiRun(&cmd.Context{})
	if seruerr != nil {
		log.Error("Failed to seru run : %s", seruerr)
	}

	return nil
}

/*
* Delete command kills the container by talking to swarm cluster and giving
* the container ID.
*
 */
func (i *Docker) Delete(assembly *global.AssemblyWithComponents, id string) (string, error) {

	pair_endpoint, perrscm := global.ParseKeyValuePair(assembly.Inputs, "endpoint")
	if perrscm != nil {
		log.Error("Failed to get the endpoint value : %s", perrscm)
	}

	pair_id, iderr := global.ParseKeyValuePair(assembly.Components[0].Outputs, "id")
	if iderr != nil {
		log.Error("Failed to get the endpoint value : %s", iderr)
	}

	var endpoint string
	if pair_endpoint.Value == BAREMETAL {

		api_host, _ := config.GetString("swarm:host")
		endpoint = api_host

	} else {
		endpoint = pair_endpoint.Value
	}

	client, _ := docker.NewClient(endpoint)
	kerr := client.KillContainer(docker.KillContainerOptions{ID: pair_id.Value})
	if kerr != nil {
		log.Error("Failed to kill the container : %s", kerr)
		return "", kerr
	}
	log.Info("Container is killed")
	return "", nil
}

/*
*
* UpdateComponent updates the ipaddress that is bound to the container
* It talks to riakdb and updates the respective component(s)
 */
func updatecomponent(assembly *global.AssemblyWithComponents, ipaddress string, id string, port string) {
	log.Debug("Update process for component with ip and container id")
	mySlice := make([]*global.KeyValuePair, 3)
	mySlice[0] = &global.KeyValuePair{Key: "ip", Value: ipaddress}
	mySlice[1] = &global.KeyValuePair{Key: "id", Value: id}
	mySlice[2] = &global.KeyValuePair{Key: "port", Value: port}

	update := global.Component{
		Id:                assembly.Components[0].Id,
		Name:              assembly.Components[0].Name,
		ToscaType:         assembly.Components[0].ToscaType,
		Inputs:            assembly.Components[0].Inputs,
		Outputs:           mySlice,
		Artifacts:         assembly.Components[0].Artifacts,
		RelatedComponents: assembly.Components[0].RelatedComponents,
		Operations:        assembly.Components[0].Operations,
		Status:            assembly.Components[0].Status,
		CreatedAt:         assembly.Components[0].CreatedAt,
	}

	conn, connerr := db.Conn("components")
	if connerr != nil {
		log.Error("Failed to riak connection : %s", connerr)
	}

	err := conn.StoreStruct(assembly.Components[0].Id, &update)
	if err != nil {
		log.Error("Failed to store the update component data : %s", err)
	}
	log.Info("Container component update was successfully.")
}

func GetMemory() int64 {
	memory, _ := config.GetInt("dockerconfig:memory")
	return int64(memory)
}

func GetSwap() int64 {
	swap, _ := config.GetInt("dockerconfig:swap")
	return int64(swap)

}

func GetCpuPeriod() int64 {
	cpuPeriod, _ := config.GetInt("dockerconfig:cpuperiod")
	return int64(cpuPeriod)

}

func GetCpuQuota() int64 {
	cpuQuota, _ := config.GetInt("dockerconfig:cpuquota")
	return int64(cpuQuota)

}

func postnetwork(containerid string, ip string) error {
	url := "http://192.168.1.100:8084/docker/networks"
    fmt.Println("URL:>", url)

    data := &global.DockerNetworksInfo{Bridge: "one", ContainerId: containerid, IpAddr: ip, Gateway: "103.56.92.1"} 
	res2B, _ := json.Marshal(data)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(res2B))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
    return nil
}

func updateIndex(ip string, pos uint) error{

	index := global.IPIndex{}
	res, err := index.Get(global.IPINDEXKEY)
	if err != nil {
		log.Error("Error: Riak didn't cooperate:\n%s.", err)
		return err
	}

	update := global.IPIndex{
		Ip:			ip, 			
		Subnet: 	res.Subnet,
		Index:		pos,
	}

	conn, connerr := db.Conn("ipindex")
	if connerr != nil {
		log.Error("Failed to riak connection : %s", connerr)
		return connerr
	}

	serr := conn.StoreStruct(global.IPINDEXKEY, &update)
	if serr != nil {
		log.Error("Failed to store the update index value : %s", serr)
		return serr
	}
	log.Info("Docker network index update was successfully.")
	return nil
}
