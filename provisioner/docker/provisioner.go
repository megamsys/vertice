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
	out := """
	*-------------------------------------------------*
	| Megamswarm           /-+-fishing nodes for you. |
	*-------------------------------------------------*\n"""

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
func (p *dockerProvisioner) Create(assembly *global.AssemblyWithComponents, id string, instance bool, act_id string) (string, error) {
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


	config := docker.Config{Image: pair_img.Value}

	copts := docker.CreateContainerOptions{Name: fmt.Sprint(assembly.Components[0].Name, ".", pair_domain.Value), Config: &config}

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
	 * hostConfig{} stuct for portbindings - to expose visible ports
	 *  Also for specfying the container configurations (memory, cpuquota etc)
	 */

	hostConfig := docker.HostConfig{
	//Memory: GetMemory(),
	//MemorySwap: GetMemory() + GetSwap(),
	//CPUQuota:  GetCpuQuota(),
	//CPUPeriod: GetCpuPeriod(),
	}
	hostConfig.PortBindings = map[docker.Port][]docker.PortBinding{
		docker.Port(Iport + "/tcp"): {{HostIP: "", HostPort: ""}},
	}

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

	updatecomponent(assembly, container_network.IPAddress, cont.ID, port)

	herr := setHostName(fmt.Sprint(assembly.Components[0].Name, ".", pair_domain.Value), container_network.IPAddress)
	if herr != nil {
		log.Error("Failed to set the host name : %s", herr)
		return "", herr
	}

	return "", nil
}


func (p *dockerProvisioner) SetCName(app provision.App, cname string) error {
	r, err := getRouterForApp(app)
	if err != nil {
		return err
	}
	return r.SetCName(cname, app.GetName())
}

func (p *dockerProvisioner) UnsetCName(app provision.App, cname string) error {
	r, err := getRouterForApp(app)
	if err != nil {
		return err
	}
	return r.UnsetCName(cname, app.GetName())
}



/*
* Register a hostname on AWS Route53 using megam seru -
*        www.github.com/megamsys/seru
*/
func setHostName(name string, ip string) error {

	s := make([]string, 4)
	s = strings.Split(name, ".")

	accesskey, _ := config.GetString("dns:aws_accesskey")
	secretkey, _ := config.GetString("dns:aws_secretkey")

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


func (p *dockerProvisioner) AddUnits(a provision.App, units uint, process string, w io.Writer) ([]provision.Unit, error) {
	if a.GetDeploys() == 0 {
		return nil, errors.New("New units can only be added after the first deployment")
	}
	if units == 0 {
		return nil, errors.New("Cannot add 0 units")
	}
	if w == nil {
		w = ioutil.Discard
	}
	writer := io.MultiWriter(w, &app.LogWriter{App: a})
	imageId, err := appCurrentImageName(a.GetName())
	if err != nil {
		return nil, err
	}
	conts, err := p.runCreateUnitsPipeline(writer, a, map[string]*containersToAdd{process: {Quantity: int(units)}}, imageId)
	if err != nil {
		return nil, err
	}
	result := make([]provision.Unit, len(conts))
	for i, c := range conts {
		result[i] = c.asUnit(a)
	}
	return result, nil
}

func (p *dockerProvisioner) RemoveUnits(a provision.App, units uint, processName string, w io.Writer) error {
	if a == nil {
		return errors.New("remove units: app should not be nil")
	}
	if units == 0 {
		return errors.New("cannot remove zero units")
	}
	var err error
	if w == nil {
		w = ioutil.Discard
	}
	imgId, err := appCurrentImageName(a.GetName())
	if err != nil {
		return err
	}
	_, processName, err = processCmdForImage(processName, imgId)
	if err != nil {
		return err
	}
	containers, err := p.listContainersByProcess(a.GetName(), processName)
	if err != nil {
		return err
	}
	if len(containers) < int(units) {
		return fmt.Errorf("cannot remove %d units from process %q, only %d available", units, processName, len(containers))
	}
	var plural string
	if units > 1 {
		plural = "s"
	}
	fmt.Fprintf(w, "\n---- Removing %d unit%s ----\n", units, plural)
	p, err = p.cloneProvisioner(nil)
	if err != nil {
		return err
	}
	toRemove := make([]container, 0, units)
	for i := 0; i < int(units); i++ {
		containerID, err := p.scheduler.GetRemovableContainer(a.GetName(), processName)
		if err != nil {
			return err
		}
		cont, err := p.getContainer(containerID)
		if err != nil {
			return err
		}
		p.scheduler.ignoredContainers = append(p.scheduler.ignoredContainers, cont.ID)
		toRemove = append(toRemove, *cont)
	}
	args := changeUnitsPipelineArgs{
		app:         a,
		toRemove:    toRemove,
		writer:      w,
		provisioner: p,
	}
	pipeline := action.NewPipeline(
		&removeOldRoutes,
		&provisionRemoveOldUnits,
		&provisionUnbindOldUnits,
	)
	err = pipeline.Execute(args)
	if err != nil {
		return fmt.Errorf("error removing routes, units weren't removed: %s", err)
	}
	return nil
}

func (p *dockerProvisioner) SetUnitStatus(unit provision.Unit, status provision.Status) error {
	container, err := p.getContainer(unit.Name)
	if err != nil {
		return err
	}
	if unit.AppName != "" && container.AppName != unit.AppName {
		return errors.New("wrong app name")
	}
	err = container.setStatus(p, status.String())
	if err != nil {
		return err
	}
	return p.checkContainer(container)
}


func (p *dockerProvisioner) ImageDeploy(app provision.App, imageId string, w io.Writer) (string, error) {
	isValid, err := isValidAppImage(app.GetName(), imageId)
	if err != nil {
		return "", err
	}
	if !isValid {
		return "", fmt.Errorf("invalid image for app %s: %s", app.GetName(), imageId)
	}
	return imageId, p.deploy(app, imageId, w)
}

func (p *dockerProvisioner) GitDeploy(app provision.App, version string, w io.Writer) (string, error) {
	imageId, err := p.gitDeploy(app, version, w)
	if err != nil {
		return "", err
	}
	return imageId, p.deployAndClean(app, imageId, w)
}


func (p *dockerProvisioner) deployAndClean(a provision.App, imageId string, w io.Writer) error {
	err := p.deploy(a, imageId, w)
	if err != nil {
		p.cleanImage(a.GetName(), imageId)
	}
	return err
}

func (p *dockerProvisioner) deploy(a provision.App, imageId string, w io.Writer) error {
	containers, err := p.listContainersByApp(a.GetName())
	if err != nil {
		return err
	}
	imageData, err := getImageCustomData(imageId)
	if err != nil {
		return err
	}
	if len(containers) == 0 {
		toAdd := make(map[string]*containersToAdd, len(imageData.Processes))
		for processName := range imageData.Processes {
			_, ok := toAdd[processName]
			if !ok {
				ct := containersToAdd{Quantity: 0}
				toAdd[processName] = &ct
			}
			toAdd[processName].Quantity++
		}
		_, err = p.runCreateUnitsPipeline(w, a, toAdd, imageId)
	} else {
		toAdd := getContainersToAdd(imageData, containers)
		_, err = p.runReplaceUnitsPipeline(w, a, toAdd, containers, imageId)
	}
	return err
}

func (p *dockerProvisioner) Nodes(app provision.App) ([]cluster.Node, error) {
	pool := app.GetPool()
	var (
		pools []provision.Pool
		err   error
	)
	if pool == "" {
		pools, err = provision.ListPools(bson.M{"$or": []bson.M{{"teams": app.GetTeamOwner()}, {"teams": bson.M{"$in": app.GetTeamsName()}}}})
	} else {
		pools, err = provision.ListPools(bson.M{"_id": pool})
	}
	if err != nil {
		return nil, err
	}
	if len(pools) == 0 {
		query := bson.M{"default": true}
		pools, err = provision.ListPools(query)
		if err != nil {
			return nil, err
		}
	}
	if len(pools) == 0 {
		return nil, errNoDefaultPool
	}
	for _, pool := range pools {
		nodes, err := p.getCluster().NodesForMetadata(map[string]string{"pool": pool.Name})
		if err != nil {
			return nil, errNoDefaultPool
		}
		if len(nodes) > 0 {
			return nodes, nil
		}
	}
	var nameList []string
	for _, pool := range pools {
		nameList = append(nameList, pool.Name)
	}
	poolsStr := strings.Join(nameList, ", pool=")
	return nil, fmt.Errorf("No nodes found with one of the following metadata: pool=%s", poolsStr)
}

func (p *dockerProvisioner) MetricEnvs(app provision.App) map[string]string {
	envMap := map[string]string{}
	bsConf, err := loadBsConfig()
	if err != nil {
		return envMap
	}
	envs, err := bsConf.envListForEndpoint("", app.GetPool())
	if err != nil {
		return envMap
	}
	for _, env := range envs {
		if strings.HasPrefix(env, "METRICS_") {
			slice := strings.SplitN(env, "=", 2)
			envMap[slice[0]] = slice[1]
		}
	}
	return envMap
}
