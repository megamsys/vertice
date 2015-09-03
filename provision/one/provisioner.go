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

package one

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"sync"
	"time"
)

var mainOneProvisioner *oneProvisioner

type oneProvisioner struct {
	client *api.RPCClient
}

func init() {
	mainOneProvisioner = &oneProvisioner{}
	provision.Register("one", mainOneProvisioner)
}

func (p *oneProvisioner) Initialize() error {
	return p.initOneCluster()
}

func (p *oneProvisioner) initOneCluster() error {
	client, err := api.NewRPCClient("http://localhost:2633/RPC2", "oneadmin", "RaifZuewjoc4")
	p.client = client
	return err
}

func (p *oneProvisioner) StartupMessage() (string, error) {
	out := "One provisioner reports the following:\n"
	out += fmt.Sprintf("    One rpc inited: %#v\n", p.client)
	return out, nil
}

func (p *oneProvisioner) GitDeploy(box provision.Box, version string, w io.Writer) (string, error) {
	imageId, err := p.gitDeploy(app, version, w)
	if err != nil {
		return "", err
	}
	return imageId, p.deployAndClean(box, imageId, w)
}

func (p *dockerProvisioner) gitDeploy(app provision.App, version string, w io.Writer) (string, error) {
	commands, err := gitDeployCmds(app, version)
	if err != nil {
		return "", err
	}
	return p.deployPipeline(app, p.getBuildImage(app), commands, w)
}


func (p *dockerProvisioner) deployPipeline(app provision.App, imageId string, commands []string, w io.Writer) (string, error) {
	actions := []*action.Action{
		&insertEmptyContainerInDB,
		&createContainer,
		&startContainer,
		&updateContainerInDB,
		&followLogsAndCommit,
	}
	pipeline := action.NewPipeline(actions...)
	buildingImage, err := appNewImageName(app.GetName())
	if err != nil {
		return "", log.WrapError(fmt.Errorf("error getting new image name for app %s", app.GetName()))
	}
	args := runContainerActionsArgs{
		app:           app,
		imageID:       imageId,
		commands:      commands,
		writer:        w,
		isDeploy:      true,
		buildingImage: buildingImage,
		provisioner:   p,
	}
	err = pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute deploy pipeline for app %s - %s", app.GetName(), err)
		return "", err
	}
	return buildingImage, nil
}


func (p *oneProvisioner) deployAndClean(b provision.Box, imageId string, w io.Writer) error {
	err := p.deploy(b, imageId, w)
	if err != nil {
		p.cleanImage(b.GetName(), imageId)
	}
	return err
}

func (p *oneProvisioner) deploy(b provision.Box, imageId string, w io.Writer) error {
	containers, err := p.listContainersByApp(b.GetName())
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

func getContainersToAdd(data ImageMetadata, oldContainers []container.Container) map[string]*containersToAdd {
	processMap := make(map[string]*containersToAdd, len(data.Processes))
	for name := range data.Processes {
		processMap[name] = &containersToAdd{}
	}
	minCount := 0
	for _, container := range oldContainers {
		if container.ProcessName == "" {
			minCount++
		}
		if _, ok := processMap[container.ProcessName]; ok {
			processMap[container.ProcessName].Quantity++
		}
	}
	if minCount == 0 {
		minCount = 1
	}
	for name, cont := range processMap {
		if cont.Quantity == 0 {
			processMap[name].Quantity = minCount
		}
	}
	return processMap
}

func (p *oneProvisioner) Destroy(box provision.Box) error {
	/*
	Action 1.
	vm := VirtualMachine{Name: "yeshapp", Client: p.client}
	result, error := vmObj.Delete()

	Action 2:
	Riak: Nuke component

	Action 3: (handled in the callee)
	Riak: Nuke Assembly

	Action 4: (handled in the  callee)
	Riak: Nuke Assemblies

	Action 5: (handled in the callee)
	Clean up log directories.

	Action 6: (do it here)
	Notify stuff via events.
	*/

	args := changeUnitsPipelineArgs{
		app:         app,
		toRemove:    containers,
		writer:      ioutil.Discard,
		provisioner: p,
		appDestroy:  true,
	}
	pipeline := action.NewPipeline(
		&removeOldRoutes,
		&provisionRemoveOldUnits,
		&provisionUnbindOldUnits,
	)
	err = pipeline.Execute(args)
	if err != nil {
		return err
	}
	images, err := listAppImages(app.GetName())
	if err != nil {
		log.Errorf("Failed to get image ids for app %s: %s", app.GetName(), err.Error())
	}
	cluster := p.Cluster()
	for _, imageId := range images {
		err := cluster.RemoveImage(imageId)
		if err != nil {
			log.Errorf("Failed to remove image %s: %s", imageId, err.Error())
		}
		err = cluster.RemoveFromRegistry(imageId)
		if err != nil {
			log.Errorf("Failed to remove image %s from registry: %s", imageId, err.Error())
		}
	}
	err = deleteAllAppImageNames(app.GetName())
	if err != nil {
		log.Errorf("Failed to remove image names from storage for app %s: %s", app.GetName(), err.Error())
	}
	r, err := getRouterForApp(app)
	if err != nil {
		log.Errorf("Failed to get router: %s", err.Error())
		return err
	}
	err = r.RemoveBackend(app.GetName())
	if err != nil {
		log.Errorf("Failed to remove route backend: %s", err.Error())
		return err
	}
	return nil
}

func (*oneProvisioner) Addr(app provision.App) (string, error) {
	r, err := getRouterForApp(app)
	if err != nil {
		log.Errorf("Failed to get router: %s", err)
		return "", err
	}
	addr, err := r.Addr(app.GetName())
	if err != nil {
		log.Errorf("Failed to obtain app %s address: %s", app.GetName(), err)
		return "", err
	}
	return addr, nil
}

func (p *oneProvisioner) runRestartAfterHooks(cont *container.Container, w io.Writer) error {
	yamlData, err := getImage(cont.Image)
	if err != nil {
		return err
	}
	cmds := yamlData.Hooks.Restart.After
	for _, cmd := range cmds {
		err := cont.Exec(p, w, w, cmd)
		if err != nil {
			return fmt.Errorf("couldn't execute restart:after hook %q(%s): %s", cmd, cont.ShortID(), err.Error())
		}
	}
	return nil
}

func addContainersWithHost(args *changeUnitsPipelineArgs) ([]container.Container, error) {
	a := args.app
	w := args.writer
	var units int
	processMsg := make([]string, 0, len(args.toAdd))
	imageId := args.imageId
	for processName, v := range args.toAdd {
		units += v.Quantity
		if processName == "" {
			_, processName, _ = processCmdForImage(processName, imageId)
		}
		processMsg = append(processMsg, fmt.Sprintf("[%s: %d]", processName, v.Quantity))
	}
	var destinationHost []string
	if args.toHost != "" {
		destinationHost = []string{args.toHost}
	}
	if w == nil {
		w = ioutil.Discard
	}
	fmt.Fprintf(w, "\n---- Starting %d new %s %s ----\n", units, pluralize("unit", units), strings.Join(processMsg, " "))
	oldContainers := make([]container.Container, 0, units)
	for processName, cont := range args.toAdd {
		for i := 0; i < cont.Quantity; i++ {
			oldContainers = append(oldContainers, container.Container{
				ProcessName: processName,
				Status:      cont.Status.String(),
			})
		}
	}
	rollbackCallback := func(c *container.Container) {
		log.Errorf("Removing container %q due failed add units.", c.ID)
		errRem := c.Remove(args.provisioner)
		if errRem != nil {
			log.Errorf("Unable to destroy container %q: %s", c.ID, errRem)
		}
	}
	var (
		createdContainers []*container.Container
		m                 sync.Mutex
	)
	err := runInContainers(oldContainers, func(c *container.Container, toRollback chan *container.Container) error {
		c, err := args.provisioner.start(c, a, imageId, w, destinationHost...)
		if err != nil {
			return err
		}
		toRollback <- c
		m.Lock()
		createdContainers = append(createdContainers, c)
		m.Unlock()
		fmt.Fprintf(w, " ---> Started unit %s [%s]\n", c.ShortID(), c.ProcessName)
		return nil
	}, rollbackCallback, true)
	if err != nil {
		return nil, err
	}
	result := make([]container.Container, len(createdContainers))
	i := 0
	for _, c := range createdContainers {
		result[i] = *c
		i++
	}
	return result, nil
}

func (p *oneProvisioner) AddUnits(a provision.App, units uint, process string, w io.Writer) ([]provision.Unit, error) {
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
		result[i] = c.AsUnit(a)
	}
	return result, nil
}

func (p *oneProvisioner) RemoveUnits(a provision.App, units uint, processName string, w io.Writer) error {
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
	fmt.Fprintf(w, "\n---- Removing %d %s ----\n", units, pluralize("unit", int(units)))
	p, err = p.cloneProvisioner(nil)
	if err != nil {
		return err
	}
	toRemove := make([]container.Container, 0, units)
	for i := 0; i < int(units); i++ {
		containerID, err := p.scheduler.GetRemovableContainer(a.GetName(), processName)
		if err != nil {
			return err
		}
		cont, err := p.GetContainer(containerID)
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

func (p *oneProvisioner) SetUnitStatus(unit provision.Unit, status provision.Status) error {
	container, err := p.GetContainer(unit.Name)
	if err != nil {
		return err
	}
	if unit.AppName != "" && container.AppName != unit.AppName {
		return errors.New("wrong app name")
	}
	err = container.SetStatus(p, status.String(), true)
	if err != nil {
		return err
	}
	return p.checkContainer(container)
}

func (p *oneProvisioner) ExecuteCommandOnce(stdout, stderr io.Writer, app provision.App, cmd string, args ...string) error {
	containers, err := p.listRunnableContainersByApp(app.GetName())
	if err != nil {
		return err
	}
	if len(containers) == 0 {
		return provision.ErrEmptyApp
	}
	container := containers[0]
	return container.Exec(p, stdout, stderr, cmd, args...)
}

func (p *oneProvisioner) ExecuteCommand(stdout, stderr io.Writer, app provision.App, cmd string, args ...string) error {
	containers, err := p.listRunnableContainersByApp(app.GetName())
	if err != nil {
		return err
	}
	if len(containers) == 0 {
		return provision.ErrEmptyApp
	}
	for _, c := range containers {
		err = c.Exec(p, stdout, stderr, cmd, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *oneProvisioner) SetCName(app provision.App, cname string) error {
	r, err := getRouterForApp(app)
	if err != nil {
		return err
	}
	return r.SetCName(cname, app.GetName())
}

func (p *oneProvisioner) UnsetCName(app provision.App, cname string) error {
	r, err := getRouterForApp(app)
	if err != nil {
		return err
	}
	return r.UnsetCName(cname, app.GetName())
}

// PlatformAdd build and push a new docker platform to register
func (p *oneProvisioner) PlatformAdd(name string, args map[string]string, w io.Writer) error {
	if args["dockerfile"] == "" {
		return errors.New("Dockerfile is required.")
	}
	if _, err := url.ParseRequestURI(args["dockerfile"]); err != nil {
		return errors.New("dockerfile parameter should be an url.")
	}
	imageName := platformImageName(name)
	cluster := p.Cluster()
	buildOptions := docker.BuildImageOptions{
		Name:           imageName,
		NoCache:        true,
		RmTmpContainer: true,
		Remote:         args["dockerfile"],
		InputStream:    nil,
		OutputStream:   w,
	}
	err := cluster.BuildImage(buildOptions)
	if err != nil {
		return err
	}
	parts := strings.Split(imageName, ":")
	var tag string
	if len(parts) > 2 {
		imageName = strings.Join(parts[:len(parts)-1], ":")
		tag = parts[len(parts)-1]
	} else if len(parts) > 1 {
		imageName = parts[0]
		tag = parts[1]
	} else {
		imageName = parts[0]
		tag = "latest"
	}
	return p.PushImage(imageName, tag)
}

func (p *oneProvisioner) PlatformUpdate(name string, args map[string]string, w io.Writer) error {
	return p.PlatformAdd(name, args, w)
}

func (p *oneProvisioner) PlatformRemove(name string) error {
	err := p.Cluster().RemoveImage(platformImageName(name))
	if err != nil && err == docker.ErrNoSuchImage {
		log.Errorf("error on remove image %s from docker.", name)
		return nil
	}
	return err
}

func (p *oneProvisioner) Units(app provision.App) ([]provision.Unit, error) {
	containers, err := p.listContainersByApp(app.GetName())
	if err != nil {
		return nil, err
	}
	units := make([]provision.Unit, len(containers))
	for i, container := range containers {
		units[i] = container.AsUnit(app)
	}
	return units, nil
}

func (p *oneProvisioner) RoutableUnits(app provision.App) ([]provision.Unit, error) {
	imageId, err := appCurrentImageName(app.GetName())
	if err != nil {
		return nil, err
	}
	webProcessName, err := getImageWebProcessName(imageId)
	if err != nil {
		return nil, err
	}
	containers, err := p.listContainersByApp(app.GetName())
	if err != nil {
		return nil, err
	}
	units := make([]provision.Unit, 0, len(containers))
	for _, container := range containers {
		if container.ProcessName == webProcessName {
			units = append(units, container.AsUnit(app))
		}
	}
	return units, nil
}

func (p *oneProvisioner) RegisterUnit(unit provision.Unit, customData map[string]interface{}) error {
	container, err := p.GetContainer(unit.Name)
	if err != nil {
		return err
	}
	if container.Status == provision.StatusBuilding.String() {
		if container.BuildingImage != "" && customData != nil {
			return saveImageCustomData(container.BuildingImage, customData)
		}
		return nil
	}
	err = container.SetStatus(p, provision.StatusStarted.String(), true)
	if err != nil {
		return err
	}
	return p.checkContainer(container)
}

func (p *oneProvisioner) MetricEnvs(app provision.App) map[string]string {
	envMap := map[string]string{}
	bsConf, err := bs.LoadConfig()
	if err != nil {
		return envMap
	}
	envs, err := bsConf.EnvListForEndpoint("", app.GetPool())
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

func pluralize(str string, sz int) string {
	if sz == 0 || sz > 1 {
		str = str + "s"
	}
	return str
}
