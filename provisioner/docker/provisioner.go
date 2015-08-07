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
	"strconv"
	"strings"
	log "code.google.com/p/log4go"
	"github.com/fsouza/go-dockerclient"
	"github.com/megamsys/libgo/log"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/provisioner"
	"github.com/megamsys/megamd/swarmc"
	"github.com/megamsys/seru/cmd"
	"github.com/megamsys/seru/cmd/seru"
	"github.com/tsuru/config"
)

/*
*
* Registers docker as provisioner in provisioner interface.
*
 */
type dockerProvisioner struct {
	cluster        *swarmc.Cluster
	storage        swarmc.Storage
	isDryMode      bool
}

var mainDockerProvisioner *dockerProvisioner

// Registers docker as provisioner in provisioner interface.
func Init() {
	mainDockerProvisioner = &dockerProvisioner{}
	provisioner.Register("docker", mainDockerProvisioner)
}

func (p *dockerProvisioner) initDockerCluster() error {
	var err error
	if p.storage == nil {
		p.storage, err = buildClusterStorage()
		if err != nil {
			return err
		}
	}

	var nodes []cluster.Node
	totalMemoryMetadata, _ := config.GetString("docker:scheduler:total-memory-metadata")
	maxUsedMemory, _ := config.GetFloat("docker:scheduler:max-used-memory")
	p.scheduler = &segregatedScheduler{
		maxMemoryRatio:      float32(maxUsedMemory),
		totalMemoryMetadata: totalMemoryMetadata,
		provisioner:         p,
	}
	p.cluster, err = cluster.New(p.scheduler, p.storage, nodes...)
	if err != nil {
		return err
	}
	return nil
}

func (p *dockerProvisioner) stopDryMode() {
	if p.isDryMode {
		p.cluster.StopDryMode()
	}
}

func (p *dockerProvisioner) dryMode(ignoredContainers []container) (*dockerProvisioner, error) {
	var err error
	overridenProvisioner := &dockerProvisioner{
		collectionName: "containers_dry_" + randomString(),
		isDryMode:      true,
	}
	containerIds := make([]string, len(ignoredContainers))
	for i := range ignoredContainers {
		containerIds[i] = ignoredContainers[i].ID
	}
	overridenProvisioner.scheduler = &segregatedScheduler{
		maxMemoryRatio:      p.scheduler.maxMemoryRatio,
		totalMemoryMetadata: p.scheduler.totalMemoryMetadata,
		provisioner:         overridenProvisioner,
		ignoredContainers:   containerIds,
	}
	overridenProvisioner.cluster, err = cluster.New(overridenProvisioner.scheduler, p.storage)
	if err != nil {
		return nil, err
	}
	overridenProvisioner.cluster.DryMode()
	containersToCopy, err := p.listAllContainers()
	if err != nil {
		return nil, err
	}
	coll := overridenProvisioner.collection()
	defer coll.Close()
	toInsert := make([]interface{}, len(containersToCopy))
	for i := range containersToCopy {
		toInsert[i] = containersToCopy[i]
	}
	if len(toInsert) > 0 {
		err = coll.Insert(toInsert...)
		if err != nil {
			return nil, err
		}
	}
	return overridenProvisioner, nil
}

func (p *dockerProvisioner) getCluster() *cluster.Cluster {
	if p.cluster == nil {
		panic("nil cluster")
	}
	return p.cluster
}



func (p *dockerProvisioner) Initialize() error {
	return p.initDockerCluster()
}



func (p *dockerProvisioner) StartupMessage() (string, error) {
	nodeList, err := p.getCluster().UnfilteredNodes()
	if err != nil {
		return "", err
	}
	out := `
	*-------------------------------------------------*
	| megamswarm      --->     fished nodes for you.  |
	*-------------------------------------------------*\n`

	for _, node := range nodeList {
		out += fmt.Sprintf("    Docker node: %s\n", node.Address)
	}
	return out, nil
}

/*
* Create provisioner is called to launch docker containers by
* talking to swarm cluster. Common provisioner for both
* Baremetal and VM-docker launch. Specify endpoint
* Swarm Host IP is added into the conf file.
*
*/
func (i *dockerProvisioner) Create(assembly *global.AssemblyWithComponents, id string, instance bool, act_id string) (string, error) {
	log.Info("%q", assembly)
	pair_endpoint, perrscm := global.ParseKeyValuePair(assembly.Inputs, "endpoint")
	if perrscm != nil {
		log.Error("Failed to get the endpoint value : %s", perrscm)
		return "", perrscm
	}

	var endpoint string
	if pair_endpoint.Value == BAREMETAL {

		/*
		 * swarm host is obtained from conf file. Swarm host is considered
		 * only when the 'endpoint' is baremetal in the Component JSON
		 */
		api_host, _ := config.GetString("docker:swarm_host")
		endpoint = api_host

		container, cerr := create(assembly, endpoint)
		if cerr != nil {
			log.Error("container creation was failed : %s", cerr)
			return "", cerr
		}

		pair_cpu, perrscm := global.ParseKeyValuePair(assembly.Inputs, "cpu")
		if perrscm != nil {
			log.Error("Failed to get the endpoint value : %s", perrscm)
		}

		pair_memory, iderr := global.ParseKeyValuePair(assembly.Inputs, "ram")
		if iderr != nil {
			log.Error("Failed to get the endpoint value : %s", iderr)
		}

        container.Cpu = pair_cpu.Value
        container.Ram = pair_memory.Value

		serr := StartContainer(container, endpoint)
		if serr != nil {
			log.Error("container starting error : %s", serr)
			return "", serr
		}

		ipaddress, iperr := setContainerNAL(container)
		if iperr != nil {
			log.Error("set container network was failed : %s", iperr)
			return "", iperr
		}

		container.IPAddress = ipaddress

		herr := setHostName(container)
		if herr != nil {
			log.Error("set host name error : %s", herr)
		}

		updateContainerJSON(assembly, container, endpoint)
	} else {
		endpoint = pair_endpoint.Value
		create(assembly, endpoint)
	}

	return "", nil
}

/*
* Delete command kills the container by talking to swarm cluster and giving
* the container ID.
*
 */
func (i *dockerProvisioner) Delete(assembly *global.AssemblyWithComponents, id string) (string, error) {

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

		api_host, _ := config.GetString("docker:swarm_host")
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
* Docker API client to connect to swarm/docker VM.
* Swarm supports all docker API endpoints
 */

func create(assembly *global.AssemblyWithComponents, endpoint string) (*global.Container, error) {

	pair_img, perrscm := global.ParseKeyValuePair(assembly.Components[0].Inputs, "source")
	if perrscm != nil {
		log.Error("Failed to get the image value : %s", perrscm)
		return &global.Container{}, perrscm
	}

	pair_domain, perrdomain := global.ParseKeyValuePair(assembly.Components[0].Inputs, "domain")
	if perrdomain != nil {
		log.Error("Failed to get the image value : %s", perrdomain)
		return &global.Container{}, perrdomain
	}

	client, _ := docker.NewClient(endpoint)

	opts := docker.PullImageOptions{
		Repository: pair_img.Value,
	}
	pullerr := client.PullImage(opts, docker.AuthConfiguration{})
	if pullerr != nil {
		log.Error("Image pulled failed : %s", pullerr)
		return &global.Container{}, pullerr
	}

	dconfig := docker.Config{Image: pair_img.Value, NetworkDisabled: true}
	copts := docker.CreateContainerOptions{Name: fmt.Sprint(assembly.Components[0].Name, ".", pair_domain.Value), Config: &dconfig}

	/*
	 * Creation of the container with copts.
	 */

	container, conerr := client.CreateContainer(copts)
	if conerr != nil {
		log.Error("Container creation failed : %s", conerr)
		return &global.Container{}, conerr
	}

	cont := &docker.Container{}
	mapP, _ := json.Marshal(container)
	json.Unmarshal([]byte(string(mapP)), cont)

	inspect, _ := client.InspectContainer(cont.ID)
	contI := &docker.Container{}
	mapInsp, _ := json.Marshal(inspect)
	json.Unmarshal([]byte(string(mapInsp)), contI)

	ins, _ := client.InspectImage(pair_img.Value)
	conti := &docker.Image{}
	mapins, _ := json.Marshal(ins)
	json.Unmarshal([]byte(string(mapins)), conti)

	confi := &docker.Config{}
	mapc, _ := json.Marshal(conti.ContainerConfig)
	json.Unmarshal([]byte(string(mapc)), confi)


	//ports := make([]string, len(confi.ExposedPorts))
	var ports string
	for k, _ := range confi.ExposedPorts {
		//i := 0
		porti := strings.Split(string(k), "/")
		//ports[i] = porti[0]
		ports = porti[0] + ", " + ports
		//i++
	}

	swarmNode := &docker.SwarmNode{}
	mapSwarm, _ := json.Marshal(contI.Node)
	json.Unmarshal([]byte(string(mapSwarm)), swarmNode)

	return &global.Container{ContainerID: cont.ID, ContainerName: cont.Name, SwarmNode: swarmNode.IP, Ports: ports}, nil
}

/*
* start the container using docker endpoint
*/
func (p *dockerProvisioner) StartContainer(container *global.Container, endpoint string) error {
		client, _ := docker.NewClient(endpoint)

	/*
	 * hostConfig{} struct for portbindings - to expose visible ports
	 *  Also for specifying the container configurations (memory, cpuquota etc)
	 */
	mem, _ := strconv.Atoi(container.Ram)
	var memory int64
	memory = int64(mem)

	cpuq, _ := strconv.Atoi(container.Cpu)
	var cpuqo int64
	cpuqo = int64(cpuq)
	cpuQuota := cpuqo * 25000

	period := 50000
	var cpuPeriod int64
	cpuPeriod = int64(period)

	hostConfig := docker.HostConfig{Memory: memory, CPUPeriod: cpuPeriod, CPUQuota: cpuQuota}

	/*
	 *   Starting container once the container is created - container ID &
	 *   hostConfig is provided to start the container.
	 *
	 */
	serr := client.StartContainer(container.ContainerID, &hostConfig)
	if serr != nil {
		log.Error("Start container was failed : %s", serr)
		return serr
	}
	return nil
}

/*
* stop the container using docker endpoint
 */
func StopContainer(container *global.Container, endpoint string) error {

	client, _ := docker.NewClient(endpoint)
	serr := client.StopContainer(container.ContainerID, 10)
	if serr != nil {
		log.Error("container was not stopped - Error : %s", serr)
		return serr
	}
	return nil
}

/*
* restart the container using docker endpoint
 */
func RestartContainer(container *global.Container, endpoint string) error {
	client, _ := docker.NewClient(endpoint)
	rerr := client.RestartContainer(container.ContainerID, 10)
	if rerr != nil {
		log.Error("container was not restarted - Error : %s", rerr)
		return rerr
	}
	return nil
}
