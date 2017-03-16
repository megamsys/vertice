/*
** Copyright [2013-2017] [Megam Systems]
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
/*"reflect"
"sort"
"strings"

"github.com/megamsys/libgo/cmd"
"github.com/megamsys/libgo/safe"
"github.com/megamsys/vertice/carton"
"github.com/megamsys/vertice/provision"
"gopkg.in/check.v1"
*/
)

/*

func (s *S) TestUpdateContainerInDBName(c *check.C) {
	c.Assert(updateContainerInDB.Name, check.Equals, "update-database-container")
}

func (s *S) TestUpdateContainerInDBForward(c *check.C) {
	cont := container.Container{Name: "myName"}
	coll := s.p.Collection()
	defer coll.Close()
	err := coll.Insert(cont)
	c.Assert(err, check.IsNil)
	cont.ID = "myID"
	context := action.FWContext{Previous: cont, Params: []interface{}{runContainerActionsArgs{
		provisioner: s.p,
	}}}
	r, err := updateContainerInDB.Forward(context)
	c.Assert(r, check.FitsTypeOf, container.Container{})
	retrieved, err := s.p.GetContainer(cont.ID)
	c.Assert(err, check.IsNil)
	c.Assert(retrieved.ID, check.Equals, cont.ID)
}

func (s *S) TestCreateContainerName(c *check.C) {
	c.Assert(createContainer.Name, check.Equals, "create-container")
}

func (s *S) TestCreateContainerForward(c *check.C) {
	err := s.newFakeImage(s.p, "tsuru/python", nil)
	c.Assert(err, check.IsNil)
	client, err := docker.NewClient(s.server.URL())
	c.Assert(err, check.IsNil)
	images, err := client.ListImages(docker.ListImagesOptions{All: true})
	c.Assert(err, check.IsNil)
	cmds := []string{"ps", "-ef"}
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	cont := container.Container{Name: "myName", AppName: app.GetName(), Type: app.GetPlatform(), Status: "created"}
	args := runContainerActionsArgs{
		app:         app,
		imageID:     images[0].ID,
		commands:    cmds,
		provisioner: s.p,
	}
	context := action.FWContext{Previous: cont, Params: []interface{}{args}}
	r, err := createContainer.Forward(context)
	c.Assert(err, check.IsNil)
	cont = r.(container.Container)
	defer cont.Remove(s.p)
	c.Assert(cont, check.FitsTypeOf, container.Container{})
	c.Assert(cont.ID, check.Not(check.Equals), "")
	c.Assert(cont.HostAddr, check.Equals, "127.0.0.1")
	dcli, err := docker.NewClient(s.server.URL())
	c.Assert(err, check.IsNil)
	cc, err := dcli.InspectContainer(cont.ID)
	c.Assert(err, check.IsNil)
	c.Assert(cc.State.Running, check.Equals, false)
}

func (s *S) TestCreateContainerBackward(c *check.C) {
	dcli, err := docker.NewClient(s.server.URL())
	c.Assert(err, check.IsNil)
	err = s.newFakeImage(s.p, "tsuru/python", nil)
	c.Assert(err, check.IsNil)
	defer dcli.RemoveImage("tsuru/python")
	conta, err := s.newContainer(nil, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(conta)
	cont := *conta
	args := runContainerActionsArgs{
		provisioner: s.p,
	}
	context := action.BWContext{FWResult: cont, Params: []interface{}{args}}
	createContainer.Backward(context)
	_, err = dcli.InspectContainer(cont.ID)
	c.Assert(err, check.NotNil)
	c.Assert(err, check.FitsTypeOf, &docker.NoSuchContainer{})
}

func (s *S) TestAddNewRouteName(c *check.C) {
	c.Assert(addNewRoute.Name, check.Equals, "add-new-routes")
}

func (s *S) TestAddNewRouteForward(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	imageName := "tsuru/app-" + app.GetName()
	customData := map[string]interface{}{
		"procfile": "web: python myapi.py\nworker: tail -f /dev/null",
	}
	err := saveImageCustomData(imageName, customData)
	c.Assert(err, check.IsNil)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	cont1 := container.Container{ID: "ble-1", AppName: app.GetName(), ProcessName: "web", HostAddr: "127.0.0.1", HostPort: "1234"}
	cont2 := container.Container{ID: "ble-2", AppName: app.GetName(), ProcessName: "web", HostAddr: "127.0.0.2", HostPort: "4321"}
	cont3 := container.Container{ID: "ble-3", AppName: app.GetName(), ProcessName: "worker", HostAddr: "127.0.0.3", HostPort: "8080"}
	defer cont1.Remove(s.p)
	defer cont2.Remove(s.p)
	defer cont3.Remove(s.p)
	args := runMachineActionsArgs{
		app:         app,
		provisioner: s.p,
		imageId:     imageName,
	}
	context := action.FWContext{Previous: []container.Container{cont1, cont2, cont3}, Params: []interface{}{args}}
	r, err := addNewRoute.Forward(context)
	c.Assert(err, check.IsNil)
	containers := r.([]container.Container)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont1.Address().String())
	c.Assert(hasRoute, check.Equals, true)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.Address().String())
	c.Assert(hasRoute, check.Equals, true)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont3.Address().String())
	c.Assert(hasRoute, check.Equals, false)
	c.Assert(containers, check.HasLen, 3)
	c.Assert(containers[0].Routable, check.Equals, true)
	c.Assert(containers[0].ID, check.Equals, "ble-1")
	c.Assert(containers[1].Routable, check.Equals, true)
	c.Assert(containers[1].ID, check.Equals, "ble-2")
	c.Assert(containers[2].Routable, check.Equals, false)
	c.Assert(containers[2].ID, check.Equals, "ble-3")
}

func (s *S) TestAddNewRouteForwardNoWeb(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	imageName := "tsuru/app-" + app.GetName()
	customData := map[string]interface{}{
		"procfile": "api: python myapi.py",
	}
	err := saveImageCustomData(imageName, customData)
	c.Assert(err, check.IsNil)
	cont1 := container.Container{ID: "ble-1", AppName: app.GetName(), ProcessName: "api", HostAddr: "127.0.0.1", HostPort: "1234"}
	cont2 := container.Container{ID: "ble-2", AppName: app.GetName(), ProcessName: "api", HostAddr: "127.0.0.2", HostPort: "4321"}
	defer cont1.Remove(s.p)
	defer cont2.Remove(s.p)
	args := runMachineActionsArgs{
		app:         app,
		provisioner: s.p,
		imageId:     imageName,
	}
	context := action.FWContext{Previous: []container.Container{cont1, cont2}, Params: []interface{}{args}}
	r, err := addNewRoute.Forward(context)
	c.Assert(err, check.IsNil)
	containers := r.([]container.Container)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont1.Address().String())
	c.Assert(hasRoute, check.Equals, true)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.Address().String())
	c.Assert(hasRoute, check.Equals, true)
	c.Assert(containers, check.HasLen, 2)
	c.Assert(containers[0].Routable, check.Equals, true)
	c.Assert(containers[0].ID, check.Equals, "ble-1")
	c.Assert(containers[1].Routable, check.Equals, true)
	c.Assert(containers[1].ID, check.Equals, "ble-2")
}

func (s *S) TestAddNewRouteForwardFailInMiddle(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	cont := container.Container{ID: "ble-1", AppName: app.GetName(), ProcessName: "", HostAddr: "addr1"}
	cont2 := container.Container{ID: "ble-2", AppName: app.GetName(), ProcessName: "", HostAddr: "addr2"}
	defer cont.Remove(s.p)
	defer cont2.Remove(s.p)
	routertest.FakeRouter.FailForIp(cont2.Address().String())
	args := runMachineActionsArgs{
		app:         app,
		provisioner: s.p,
	}
	prevContainers := []container.Container{cont, cont2}
	context := action.FWContext{Previous: prevContainers, Params: []interface{}{args}}
	_, err := addNewRoute.Forward(context)
	c.Assert(err, check.Equals, routertest.ErrForcedFailure)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont.Address().String())
	c.Assert(hasRoute, check.Equals, false)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.Address().String())
	c.Assert(hasRoute, check.Equals, false)
	c.Assert(prevContainers[0].Routable, check.Equals, true)
	c.Assert(prevContainers[0].ID, check.Equals, "ble-1")
	c.Assert(prevContainers[1].Routable, check.Equals, false)
	c.Assert(prevContainers[1].ID, check.Equals, "ble-2")
}

func (s *S) TestAddNewRouteBackward(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	cont1 := container.Container{ID: "ble-1", AppName: app.GetName(), ProcessName: "web", HostAddr: "127.0.0.1", HostPort: "1234"}
	cont2 := container.Container{ID: "ble-2", AppName: app.GetName(), ProcessName: "web", HostAddr: "127.0.0.2", HostPort: "4321"}
	cont3 := container.Container{ID: "ble-3", AppName: app.GetName(), ProcessName: "worker", HostAddr: "127.0.0.3", HostPort: "8080"}
	defer cont1.Remove(s.p)
	defer cont2.Remove(s.p)
	defer cont3.Remove(s.p)
	err := routertest.FakeRouter.AddRoute(app.GetName(), cont1.Address())
	c.Assert(err, check.IsNil)
	err = routertest.FakeRouter.AddRoute(app.GetName(), cont2.Address())
	c.Assert(err, check.IsNil)
	args := runMachineActionsArgs{
		app:         app,
		provisioner: s.p,
	}
	cont1.Routable = true
	cont2.Routable = true
	context := action.BWContext{FWResult: []container.Container{cont1, cont2, cont3}, Params: []interface{}{args}}
	addNewRoute.Backward(context)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont1.Address().String())
	c.Assert(hasRoute, check.Equals, false)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.Address().String())
	c.Assert(hasRoute, check.Equals, false)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont3.Address().String())
	c.Assert(hasRoute, check.Equals, false)
}

func (s *S) TestRemoveOldRoutesName(c *check.C) {
	c.Assert(removeOldRoutes.Name, check.Equals, "remove-old-routes")
}

func (s *S) TestRemoveOldRoutesForward(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	cont1 := container.Container{ID: "ble-1", AppName: app.GetName(), ProcessName: "web", HostAddr: "127.0.0.1", HostPort: "1234"}
	cont2 := container.Container{ID: "ble-2", AppName: app.GetName(), ProcessName: "web", HostAddr: "127.0.0.2", HostPort: "4321"}
	cont3 := container.Container{ID: "ble-3", AppName: app.GetName(), ProcessName: "worker", HostAddr: "127.0.0.3", HostPort: "8080"}
	defer cont1.Remove(s.p)
	defer cont2.Remove(s.p)
	defer cont3.Remove(s.p)
	err := routertest.FakeRouter.AddRoute(app.GetName(), cont1.Address())
	c.Assert(err, check.IsNil)
	err = routertest.FakeRouter.AddRoute(app.GetName(), cont2.Address())
	c.Assert(err, check.IsNil)
	args := runMachineActionsArgs{
		app:         app,
		toRemove:    []container.Container{cont1, cont2, cont3},
		provisioner: s.p,
	}
	context := action.FWContext{Previous: []container.Container{}, Params: []interface{}{args}}
	r, err := removeOldRoutes.Forward(context)
	c.Assert(err, check.IsNil)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont1.Address().String())
	c.Assert(hasRoute, check.Equals, false)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.Address().String())
	c.Assert(hasRoute, check.Equals, false)
	containers := r.([]container.Container)
	c.Assert(containers, check.DeepEquals, []container.Container{})
	c.Assert(args.toRemove[0].Routable, check.Equals, true)
	c.Assert(args.toRemove[1].Routable, check.Equals, true)
	c.Assert(args.toRemove[2].Routable, check.Equals, false)
}

func (s *S) TestRemoveOldRoutesForwardFailInMiddle(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	cont := container.Container{ID: "ble-1", AppName: app.GetName(), ProcessName: "web", HostAddr: "addr1"}
	cont2 := container.Container{ID: "ble-2", AppName: app.GetName(), ProcessName: "web", HostAddr: "addr2"}
	defer cont.Remove(s.p)
	defer cont2.Remove(s.p)
	err := routertest.FakeRouter.AddRoute(app.GetName(), cont.Address())
	c.Assert(err, check.IsNil)
	err = routertest.FakeRouter.AddRoute(app.GetName(), cont2.Address())
	c.Assert(err, check.IsNil)
	routertest.FakeRouter.FailForIp(cont2.Address().String())
	args := runMachineActionsArgs{
		app:         app,
		toRemove:    []container.Container{cont, cont2},
		provisioner: s.p,
	}
	context := action.FWContext{Previous: []container.Container{}, Params: []interface{}{args}}
	_, err = removeOldRoutes.Forward(context)
	c.Assert(err, check.Equals, routertest.ErrForcedFailure)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont.Address().String())
	c.Assert(hasRoute, check.Equals, true)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.Address().String())
	c.Assert(hasRoute, check.Equals, true)
	c.Assert(args.toRemove[0].Routable, check.Equals, true)
	c.Assert(args.toRemove[1].Routable, check.Equals, false)
}

func (s *S) TestRemoveOldRoutesBackward(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	cont := container.Container{ID: "ble-1", AppName: app.GetName(), ProcessName: "web"}
	cont2 := container.Container{ID: "ble-2", AppName: app.GetName(), ProcessName: "web"}
	defer cont.Remove(s.p)
	defer cont2.Remove(s.p)
	cont.Routable = true
	cont2.Routable = true
	args := runMachineActionsArgs{
		app:         app,
		toRemove:    []container.Container{cont, cont2},
		provisioner: s.p,
	}
	context := action.BWContext{Params: []interface{}{args}}
	removeOldRoutes.Backward(context)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont.Address().String())
	c.Assert(hasRoute, check.Equals, true)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.Address().String())
	c.Assert(hasRoute, check.Equals, true)
}

func (s *S) TestSetNetworkInfoName(c *check.C) {
	c.Assert(setNetworkInfo.Name, check.Equals, "set-network-info")
}

func (s *S) TestSetNetworkInfoForward(c *check.C) {
	conta, err := s.newContainer(nil, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(conta)
	cont := *conta
	context := action.FWContext{Previous: cont, Params: []interface{}{runContainerActionsArgs{
		provisioner: s.p,
	}}}
	r, err := setNetworkInfo.Forward(context)
	c.Assert(err, check.IsNil)
	cont = r.(container.Container)
	c.Assert(cont, check.FitsTypeOf, container.Container{})
	c.Assert(cont.IP, check.Not(check.Equals), "")
	c.Assert(cont.HostPort, check.Not(check.Equals), "")
}

func (s *S) TestSetImage(c *check.C) {
	conta, err := s.newContainer(nil, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(conta)
	cont := *conta
	context := action.FWContext{Previous: cont, Params: []interface{}{runContainerActionsArgs{
		provisioner: s.p,
	}}}
	r, err := setNetworkInfo.Forward(context)
	c.Assert(err, check.IsNil)
	cont = r.(container.Container)
	c.Assert(cont, check.FitsTypeOf, container.Container{})
	c.Assert(cont.HostPort, check.Not(check.Equals), "")
}

func (s *S) TestStartContainerForward(c *check.C) {
	conta, err := s.newContainer(nil, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(conta)
	cont := *conta
	context := action.FWContext{Previous: cont, Params: []interface{}{runContainerActionsArgs{
		provisioner: s.p,
		app:         provisiontest.NewFakeApp("myapp", "python", 1),
	}}}
	r, err := startContainer.Forward(context)
	c.Assert(err, check.IsNil)
	cont = r.(container.Container)
	c.Assert(cont, check.FitsTypeOf, container.Container{})
}

func (s *S) TestStartContainerBackward(c *check.C) {
	dcli, err := docker.NewClient(s.server.URL())
	c.Assert(err, check.IsNil)
	err = s.newFakeImage(s.p, "tsuru/python", nil)
	c.Assert(err, check.IsNil)
	defer dcli.RemoveImage("tsuru/python")
	conta, err := s.newContainer(nil, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(conta)
	cont := *conta
	err = dcli.StartContainer(cont.ID, nil)
	c.Assert(err, check.IsNil)
	context := action.BWContext{FWResult: cont, Params: []interface{}{runContainerActionsArgs{
		provisioner: s.p,
	}}}
	startContainer.Backward(context)
	cc, err := dcli.InspectContainer(cont.ID)
	c.Assert(err, check.IsNil)
	c.Assert(cc.State.Running, check.Equals, false)
}

func (s *S) TestaddBoxsToHostName(c *check.C) {
	c.Assert(addBoxsToHost.Name, check.Equals, "add-boxs-to-host")
}

func (s *S) TestaddBoxsToHostForward(c *check.C) {
	p, err := s.startMultipleServersCluster()
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("myapp-2", "python", 0)
	defer p.Destroy(app)
	p.Provision(app)
	coll := p.Collection()
	defer coll.Close()
	coll.Insert(container.Container{ID: "container-id", AppName: app.GetName(), Version: "container-version", Image: "tsuru/python"})
	defer coll.RemoveAll(bson.M{"appname": app.GetName()})
	imageId, err := appNewImageName(app.GetName())
	c.Assert(err, check.IsNil)
	err = s.newFakeImage(p, imageId, nil)
	c.Assert(err, check.IsNil)
	args := runMachineActionsArgs{
		app:         app,
		toHost:      "localhost",
		toAdd:       map[string]*containersToAdd{"web": {Quantity: 2}},
		imageId:     imageId,
		provisioner: p,
	}
	context := action.FWContext{Params: []interface{}{args}}
	result, err := addBoxsToHost.Forward(context)
	c.Assert(err, check.IsNil)
	containers := result.([]container.Container)
	c.Assert(containers, check.HasLen, 2)
	c.Assert(containers[0].HostAddr, check.Equals, "localhost")
	c.Assert(containers[1].HostAddr, check.Equals, "localhost")
	count, err := coll.Find(bson.M{"appname": app.GetName()}).Count()
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, 3)
}

func (s *S) TestaddBoxsToHostForwardWithoutHost(c *check.C) {
	p, err := s.startMultipleServersCluster()
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("myapp-2", "python", 0)
	defer p.Destroy(app)
	p.Provision(app)
	coll := p.Collection()
	defer coll.Close()
	imageId, err := appNewImageName(app.GetName())
	c.Assert(err, check.IsNil)
	err = s.newFakeImage(p, imageId, nil)
	c.Assert(err, check.IsNil)
	args := runMachineActionsArgs{
		app:         app,
		toAdd:       map[string]*containersToAdd{"web": {Quantity: 3}},
		imageId:     imageId,
		provisioner: p,
	}
	context := action.FWContext{Params: []interface{}{args}}
	result, err := addBoxsToHost.Forward(context)
	c.Assert(err, check.IsNil)
	containers := result.([]container.Container)
	c.Assert(containers, check.HasLen, 3)
	addrs := []string{containers[0].HostAddr, containers[1].HostAddr, containers[2].HostAddr}
	sort.Strings(addrs)
	isValid := reflect.DeepEqual(addrs, []string{"127.0.0.1", "localhost", "localhost"}) ||
		reflect.DeepEqual(addrs, []string{"127.0.0.1", "127.0.0.1", "localhost"})
	if !isValid {
		clusterNodes, _ := p.Cluster().UnfilteredNodes()
		c.Fatalf("Expected multiple hosts, got: %#v\nAvailable nodes: %#v", containers, clusterNodes)
	}
	count, err := coll.Find(bson.M{"appname": app.GetName()}).Count()
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, 3)
}

func (s *S) TestaddBoxsToHostBackward(c *check.C) {
	err := s.newFakeImage(s.p, "tsuru/python", nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("myapp-xxx-1", "python", 0)
	defer s.p.Destroy(app)
	s.p.Provision(app)
	coll := s.p.Collection()
	defer coll.Close()
	cont := container.Container{ID: "container-id", AppName: app.GetName(), Version: "container-version", Image: "tsuru/python"}
	coll.Insert(cont)
	defer coll.RemoveAll(bson.M{"appname": app.GetName()})
	args := runMachineActionsArgs{
		provisioner: s.p,
	}
	context := action.BWContext{FWResult: []container.Container{cont}, Params: []interface{}{args}}
	addBoxsToHost.Backward(context)
	_, err = s.p.GetContainer(cont.ID)
	c.Assert(err, check.Equals, provision.ErrBoxNotFound)
}

func (s *S) TestProvisionRemoveOldUnitsName(c *check.C) {
	c.Assert(provisionRemoveOldUnits.Name, check.Equals, "provision-remove-old-units")
}

func (s *S) TestProvisionRemoveOldUnitsForward(c *check.C) {
	cont, err := s.newContainer(nil, nil)
	c.Assert(err, check.IsNil)
	defer routertest.FakeRouter.RemoveBackend(cont.AppName)
	client, err := docker.NewClient(s.server.URL())
	c.Assert(err, check.IsNil)
	err = client.StartContainer(cont.ID, nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp(cont.AppName, "python", 0)
	unit := cont.AsUnit(app)
	app.BindUnit(&unit)
	args := runMachineActionsArgs{
		app:         app,
		toRemove:    []container.Container{*cont},
		provisioner: s.p,
	}
	context := action.FWContext{Params: []interface{}{args}, Previous: []container.Container{}}
	result, err := provisionRemoveOldUnits.Forward(context)
	c.Assert(err, check.IsNil)
	resultContainers := result.([]container.Container)
	c.Assert(resultContainers, check.DeepEquals, []container.Container{})
	_, err = s.p.GetContainer(cont.ID)
	c.Assert(err, check.NotNil)
}

func (s *S) TestProvisionUnbindOldUnitsName(c *check.C) {
	c.Assert(provisionUnbindOldUnits.Name, check.Equals, "provision-unbind-old-units")
}

func (s *S) TestProvisionUnbindOldUnitsForward(c *check.C) {
	cont, err := s.newContainer(nil, nil)
	c.Assert(err, check.IsNil)
	defer routertest.FakeRouter.RemoveBackend(cont.AppName)
	client, err := docker.NewClient(s.server.URL())
	c.Assert(err, check.IsNil)
	err = client.StartContainer(cont.ID, nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp(cont.AppName, "python", 0)
	unit := cont.AsUnit(app)
	app.BindUnit(&unit)
	args := runMachineActionsArgs{
		app:         app,
		toRemove:    []container.Container{*cont},
		provisioner: s.p,
	}
	context := action.FWContext{Params: []interface{}{args}, Previous: []container.Container{}}
	result, err := provisionUnbindOldUnits.Forward(context)
	c.Assert(err, check.IsNil)
	resultContainers := result.([]container.Container)
	c.Assert(resultContainers, check.DeepEquals, []container.Container{})
	c.Assert(app.HasBind(&unit), check.Equals, false)
}

func (s *S) TestFollowLogsAndCommitName(c *check.C) {
	c.Assert(followLogsAndCommit.Name, check.Equals, "follow-logs-and-commit")
}

func (s *S) TestFollowLogsAndCommitForward(c *check.C) {
	err := s.newFakeImage(s.p, "tsuru/python", nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("mightyapp", "python", 1)
	nextImgName, err := appNewImageName(app.GetName())
	c.Assert(err, check.IsNil)
	cont := container.Container{AppName: "mightyapp", ID: "myid123", BuildingImage: nextImgName}
	err = cont.Create(&container.CreateArgs{
		App:         app,
		ImageID:     "tsuru/python",
		Commands:    []string{"foo"},
		Provisioner: s.p,
	})
	c.Assert(err, check.IsNil)
	buf := safe.NewBuffer(nil)
	args := runContainerActionsArgs{writer: buf, provisioner: s.p}
	context := action.FWContext{Params: []interface{}{args}, Previous: cont}
	imageId, err := followLogsAndCommit.Forward(context)
	c.Assert(err, check.IsNil)
	c.Assert(imageId, check.Equals, "tsuru/app-mightyapp:v1")
	c.Assert(buf.String(), check.Not(check.Equals), "")
	var dbCont container.Container
	coll := s.p.Collection()
	defer coll.Close()
	err = coll.Find(bson.M{"id": cont.ID}).One(&dbCont)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "not found")
	_, err = s.p.Cluster().InspectContainer(cont.ID)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Matches, "No such container.*")
	err = s.p.Cluster().RemoveImage("tsuru/app-mightyapp:v1")
	c.Assert(err, check.IsNil)
}

func (s *S) TestFollowLogsAndCommitForwardNonZeroStatus(c *check.C) {
	err := s.newFakeImage(s.p, "tsuru/python", nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	cont := container.Container{AppName: "mightyapp"}
	err = cont.Create(&container.CreateArgs{
		App:         app,
		ImageID:     "tsuru/python",
		Commands:    []string{"foo"},
		Provisioner: s.p,
	})
	c.Assert(err, check.IsNil)
	err = s.server.MutateContainer(cont.ID, docker.State{ExitCode: 1})
	c.Assert(err, check.IsNil)
	buf := safe.NewBuffer(nil)
	args := runContainerActionsArgs{writer: buf, provisioner: s.p}
	context := action.FWContext{Params: []interface{}{args}, Previous: cont}
	imageId, err := followLogsAndCommit.Forward(context)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "Exit status 1")
	c.Assert(imageId, check.IsNil)
}

func (s *S) TestFollowLogsAndCommitForwardWaitFailure(c *check.C) {
	s.server.PrepareFailure("failed to wait for the container", "/containers/./wait")
	defer s.server.ResetFailure("failed to wait for the container")
	err := s.newFakeImage(s.p, "tsuru/python", nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	cont := container.Container{AppName: "mightyapp"}
	err = cont.Create(&container.CreateArgs{
		App:         app,
		ImageID:     "tsuru/python",
		Commands:    []string{"foo"},
		Provisioner: s.p,
	})
	c.Assert(err, check.IsNil)
	buf := safe.NewBuffer(nil)
	args := runContainerActionsArgs{writer: buf, provisioner: s.p}
	context := action.FWContext{Params: []interface{}{args}, Previous: cont}
	imageId, err := followLogsAndCommit.Forward(context)
	c.Assert(err, check.ErrorMatches, `.*failed to wait for the container\n$`)
	c.Assert(imageId, check.IsNil)
}

*/
