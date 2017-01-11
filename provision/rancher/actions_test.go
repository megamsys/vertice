package rancher

/*import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sort"
	"strings"
)


func (s *S) TestCreateContainerName(c *check.C) {
	c.Assert(createContainer.Name, check.Equals, "create-container")
}

func (s *S) TestCreateContainerForward(c *check.C) {
	err := s.newFakeImage(s.p, "github.com/megamsys/vertice/python", nil)
	c.Assert(err, check.IsNil)
	client, err := docker.NewClient(s.server.URL())
	c.Assert(err, check.IsNil)
	images, err := client.ListImages(docker.ListImagesOptions{All: true})
	c.Assert(err, check.IsNil)
	cmds := []string{"ps", "-ef"}
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	cont := container{Name: "myName", AppName: app.GetName(), Type: app.GetPlatform(), Status: "created"}
	args := runContainerActionsArgs{
		app:         app,
		imageID:     images[0].ID,
		commands:    cmds,
		provisioner: s.p,
	}
	context := action.FWContext{Previous: cont, Params: []interface{}{args}}
	r, err := createContainer.Forward(context)
	c.Assert(err, check.IsNil)
	cont = r.(container)
	defer cont.remove(s.p)
	c.Assert(cont, check.FitsTypeOf, container{})
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
	err = s.newFakeImage(s.p, "github.com/megamsys/vertice/python", nil)
	c.Assert(err, check.IsNil)
	defer dcli.RemoveImage("github.com/megamsys/vertice/python")
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
	imageName := "github.com/megamsys/vertice/app-" + app.GetName()
	customData := map[string]interface{}{
		"procfile": "web: python myapi.py\nworker: tail -f /dev/null",
	}
	err := saveImageCustomData(imageName, customData)
	c.Assert(err, check.IsNil)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	cont1 := container{ID: "ble-1", AppName: app.GetName(), ProcessName: "web", HostAddr: "127.0.0.1", HostPort: "1234"}
	cont2 := container{ID: "ble-2", AppName: app.GetName(), ProcessName: "web", HostAddr: "127.0.0.2", HostPort: "4321"}
	cont3 := container{ID: "ble-3", AppName: app.GetName(), ProcessName: "worker", HostAddr: "127.0.0.3", HostPort: "8080"}
	defer cont1.remove(s.p)
	defer cont2.remove(s.p)
	defer cont3.remove(s.p)
	args := changeUnitsPipelineArgs{
		app:         app,
		provisioner: s.p,
		imageId:     imageName,
	}
	context := action.FWContext{Previous: []container{cont1, cont2, cont3}, Params: []interface{}{args}}
	r, err := addNewRoute.Forward(context)
	c.Assert(err, check.IsNil)
	containers := r.([]container)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont1.getAddress().String())
	c.Assert(hasRoute, check.Equals, true)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.getAddress().String())
	c.Assert(hasRoute, check.Equals, true)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont3.getAddress().String())
	c.Assert(hasRoute, check.Equals, false)
	c.Assert(containers, check.HasLen, 3)
	c.Assert(containers[0].routable, check.Equals, true)
	c.Assert(containers[0].ID, check.Equals, "ble-1")
	c.Assert(containers[1].routable, check.Equals, true)
	c.Assert(containers[1].ID, check.Equals, "ble-2")
	c.Assert(containers[2].routable, check.Equals, false)
	c.Assert(containers[2].ID, check.Equals, "ble-3")
}

func (s *S) TestAddNewRouteForwardNoWeb(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	imageName := "github.com/megamsys/vertice/app-" + app.GetName()
	customData := map[string]interface{}{
		"procfile": "api: python myapi.py",
	}
	err := saveImageCustomData(imageName, customData)
	c.Assert(err, check.IsNil)
	cont1 := container{ID: "ble-1", AppName: app.GetName(), ProcessName: "api", HostAddr: "127.0.0.1", HostPort: "1234"}
	cont2 := container{ID: "ble-2", AppName: app.GetName(), ProcessName: "api", HostAddr: "127.0.0.2", HostPort: "4321"}
	defer cont1.remove(s.p)
	defer cont2.remove(s.p)
	args := changeUnitsPipelineArgs{
		app:         app,
		provisioner: s.p,
		imageId:     imageName,
	}
	context := action.FWContext{Previous: []container{cont1, cont2}, Params: []interface{}{args}}
	r, err := addNewRoute.Forward(context)
	c.Assert(err, check.IsNil)
	containers := r.([]container)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont1.getAddress().String())
	c.Assert(hasRoute, check.Equals, true)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.getAddress().String())
	c.Assert(hasRoute, check.Equals, true)
	c.Assert(containers, check.HasLen, 2)
	c.Assert(containers[0].routable, check.Equals, true)
	c.Assert(containers[0].ID, check.Equals, "ble-1")
	c.Assert(containers[1].routable, check.Equals, true)
	c.Assert(containers[1].ID, check.Equals, "ble-2")
}

func (s *S) TestAddNewRouteForwardFailInMiddle(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	cont := container{ID: "ble-1", AppName: app.GetName(), ProcessName: "", HostAddr: "addr1"}
	cont2 := container{ID: "ble-2", AppName: app.GetName(), ProcessName: "", HostAddr: "addr2"}
	defer cont.remove(s.p)
	defer cont2.remove(s.p)
	routertest.FakeRouter.FailForIp(cont2.getAddress().String())
	args := changeUnitsPipelineArgs{
		app:         app,
		provisioner: s.p,
	}
	prevContainers := []container{cont, cont2}
	context := action.FWContext{Previous: prevContainers, Params: []interface{}{args}}
	_, err := addNewRoute.Forward(context)
	c.Assert(err, check.Equals, routertest.ErrForcedFailure)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont.getAddress().String())
	c.Assert(hasRoute, check.Equals, false)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.getAddress().String())
	c.Assert(hasRoute, check.Equals, false)
	c.Assert(prevContainers[0].routable, check.Equals, true)
	c.Assert(prevContainers[0].ID, check.Equals, "ble-1")
	c.Assert(prevContainers[1].routable, check.Equals, false)
	c.Assert(prevContainers[1].ID, check.Equals, "ble-2")
}

func (s *S) TestAddNewRouteBackward(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	cont1 := container{ID: "ble-1", AppName: app.GetName(), ProcessName: "web", HostAddr: "127.0.0.1", HostPort: "1234"}
	cont2 := container{ID: "ble-2", AppName: app.GetName(), ProcessName: "web", HostAddr: "127.0.0.2", HostPort: "4321"}
	cont3 := container{ID: "ble-3", AppName: app.GetName(), ProcessName: "worker", HostAddr: "127.0.0.3", HostPort: "8080"}
	defer cont1.remove(s.p)
	defer cont2.remove(s.p)
	defer cont3.remove(s.p)
	err := routertest.FakeRouter.AddRoute(app.GetName(), cont1.getAddress())
	c.Assert(err, check.IsNil)
	err = routertest.FakeRouter.AddRoute(app.GetName(), cont2.getAddress())
	c.Assert(err, check.IsNil)
	args := changeUnitsPipelineArgs{
		app:         app,
		provisioner: s.p,
	}
	cont1.routable = true
	cont2.routable = true
	context := action.BWContext{FWResult: []container{cont1, cont2, cont3}, Params: []interface{}{args}}
	addNewRoute.Backward(context)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont1.getAddress().String())
	c.Assert(hasRoute, check.Equals, false)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.getAddress().String())
	c.Assert(hasRoute, check.Equals, false)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont3.getAddress().String())
	c.Assert(hasRoute, check.Equals, false)
}

func (s *S) TestRemoveOldRoutesName(c *check.C) {
	c.Assert(removeOldRoutes.Name, check.Equals, "remove-old-routes")
}

func (s *S) TestRemoveOldRoutesForward(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	cont1 := container{ID: "ble-1", AppName: app.GetName(), ProcessName: "web", HostAddr: "127.0.0.1", HostPort: "1234"}
	cont2 := container{ID: "ble-2", AppName: app.GetName(), ProcessName: "web", HostAddr: "127.0.0.2", HostPort: "4321"}
	cont3 := container{ID: "ble-3", AppName: app.GetName(), ProcessName: "worker", HostAddr: "127.0.0.3", HostPort: "8080"}
	defer cont1.remove(s.p)
	defer cont2.remove(s.p)
	defer cont3.remove(s.p)
	err := routertest.FakeRouter.AddRoute(app.GetName(), cont1.getAddress())
	c.Assert(err, check.IsNil)
	err = routertest.FakeRouter.AddRoute(app.GetName(), cont2.getAddress())
	c.Assert(err, check.IsNil)
	args := changeUnitsPipelineArgs{
		app:         app,
		toRemove:    []container{cont1, cont2, cont3},
		provisioner: s.p,
	}
	context := action.FWContext{Previous: []container{}, Params: []interface{}{args}}
	r, err := removeOldRoutes.Forward(context)
	c.Assert(err, check.IsNil)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont1.getAddress().String())
	c.Assert(hasRoute, check.Equals, false)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.getAddress().String())
	c.Assert(hasRoute, check.Equals, false)
	containers := r.([]container)
	c.Assert(containers, check.DeepEquals, []container{})
	c.Assert(args.toRemove[0].routable, check.Equals, true)
	c.Assert(args.toRemove[1].routable, check.Equals, true)
	c.Assert(args.toRemove[2].routable, check.Equals, false)
}

func (s *S) TestRemoveOldRoutesForwardFailInMiddle(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	cont := container{ID: "ble-1", AppName: app.GetName(), ProcessName: "web", HostAddr: "addr1"}
	cont2 := container{ID: "ble-2", AppName: app.GetName(), ProcessName: "web", HostAddr: "addr2"}
	defer cont.remove(s.p)
	defer cont2.remove(s.p)
	err := routertest.FakeRouter.AddRoute(app.GetName(), cont.getAddress())
	c.Assert(err, check.IsNil)
	err = routertest.FakeRouter.AddRoute(app.GetName(), cont2.getAddress())
	c.Assert(err, check.IsNil)
	routertest.FakeRouter.FailForIp(cont2.getAddress().String())
	args := changeUnitsPipelineArgs{
		app:         app,
		toRemove:    []container{cont, cont2},
		provisioner: s.p,
	}
	context := action.FWContext{Previous: []container{}, Params: []interface{}{args}}
	_, err = removeOldRoutes.Forward(context)
	c.Assert(err, check.Equals, routertest.ErrForcedFailure)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont.getAddress().String())
	c.Assert(hasRoute, check.Equals, true)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.getAddress().String())
	c.Assert(hasRoute, check.Equals, true)
	c.Assert(args.toRemove[0].routable, check.Equals, true)
	c.Assert(args.toRemove[1].routable, check.Equals, false)
}

func (s *S) TestRemoveOldRoutesBackward(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend(app.GetName())
	defer routertest.FakeRouter.RemoveBackend(app.GetName())
	cont := container{ID: "ble-1", AppName: app.GetName(), ProcessName: "web"}
	cont2 := container{ID: "ble-2", AppName: app.GetName(), ProcessName: "web"}
	defer cont.remove(s.p)
	defer cont2.remove(s.p)
	cont.routable = true
	cont2.routable = true
	args := changeUnitsPipelineArgs{
		app:         app,
		toRemove:    []container{cont, cont2},
		provisioner: s.p,
	}
	context := action.BWContext{Params: []interface{}{args}}
	removeOldRoutes.Backward(context)
	hasRoute := routertest.FakeRouter.HasRoute(app.GetName(), cont.getAddress().String())
	c.Assert(hasRoute, check.Equals, true)
	hasRoute = routertest.FakeRouter.HasRoute(app.GetName(), cont2.getAddress().String())
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
	cont = r.(container)
	c.Assert(cont, check.FitsTypeOf, container{})
	c.Assert(cont.IP, check.Not(check.Equals), "")
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
	cont = r.(container)
	c.Assert(cont, check.FitsTypeOf, container{})
}

func (s *S) TestStartContainerBackward(c *check.C) {
	dcli, err := docker.NewClient(s.server.URL())
	c.Assert(err, check.IsNil)
	err = s.newFakeImage(s.p, "github.com/megamsys/vertice/python", nil)
	c.Assert(err, check.IsNil)
	defer dcli.RemoveImage("github.com/megamsys/vertice/python")
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

func (s *S) TestAddBoxToHostName(c *check.C) {
	c.Assert(addBoxsToHost.Name, check.Equals, "add-boxs-to-host")
}

func (s *S) TestAddBoxsToHostForward(c *check.C) {
	p, err := s.startMultipleServersCluster()
	c.Assert(err, check.IsNil)
	defer s.stopMultipleServersCluster(p)
	app := provisiontest.NewFakeApp("myapp-2", "python", 0)
	defer p.Destroy(app)
	p.Provision(app)
	coll := p.collection()
	defer coll.Close()
	coll.Insert(container{ID: "container-id", AppName: app.GetName(), Version: "container-version", Image: "github.com/megamsys/vertice/python"})
	defer coll.RemoveAll(bson.M{"appname": app.GetName()})
	imageId, err := appNewImageName(app.GetName())
	c.Assert(err, check.IsNil)
	err = s.newFakeImage(p, imageId, nil)
	c.Assert(err, check.IsNil)
	args := changeUnitsPipelineArgs{
		app:         app,
		toHost:      "localhost",
		toAdd:       map[string]*containersToAdd{"web": {Quantity: 2}},
		imageId:     imageId,
		provisioner: p,
	}
	context := action.FWContext{Params: []interface{}{args}}
	result, err := addBoxsToHost.Forward(context)
	c.Assert(err, check.IsNil)
	containers := result.([]container)
	c.Assert(containers, check.HasLen, 2)
	c.Assert(containers[0].HostAddr, check.Equals, "localhost")
	c.Assert(containers[1].HostAddr, check.Equals, "localhost")
	count, err := coll.Find(bson.M{"appname": app.GetName()}).Count()
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, 3)
}

func (s *S) TestAddBoxsToHostForwardWithoutHost(c *check.C) {
	p, err := s.startMultipleServersCluster()
	c.Assert(err, check.IsNil)
	defer s.stopMultipleServersCluster(p)
	app := provisiontest.NewFakeApp("myapp-2", "python", 0)
	defer p.Destroy(app)
	p.Provision(app)
	coll := p.collection()
	defer coll.Close()
	imageId, err := appNewImageName(app.GetName())
	c.Assert(err, check.IsNil)
	err = s.newFakeImage(p, imageId, nil)
	c.Assert(err, check.IsNil)
	args := changeUnitsPipelineArgs{
		app:         app,
		toAdd:       map[string]*containersToAdd{"web": {Quantity: 3}},
		imageId:     imageId,
		provisioner: p,
	}
	context := action.FWContext{Params: []interface{}{args}}
	result, err := addBoxsToHost.Forward(context)
	c.Assert(err, check.IsNil)
	containers := result.([]container)
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

func (s *S) TestAddBoxsToHostBackward(c *check.C) {
	err := s.newFakeImage(s.p, "github.com/megamsys/vertice/python", nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("myapp-xxx-1", "python", 0)
	defer s.p.Destroy(app)
	s.p.Provision(app)
	coll := s.p.collection()
	defer coll.Close()
	cont := container{ID: "container-id", AppName: app.GetName(), Version: "container-version", Image: "github.com/megamsys/vertice/python"}
	coll.Insert(cont)
	defer coll.RemoveAll(bson.M{"appname": app.GetName()})
	args := changeUnitsPipelineArgs{
		provisioner: s.p,
	}
	context := action.BWContext{FWResult: []container{cont}, Params: []interface{}{args}}
	addBoxsToHost.Backward(context)
	_, err = s.p.getContainer(cont.ID)
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
	unit := cont.asUnit(app)
	app.BindUnit(&unit)
	args := changeUnitsPipelineArgs{
		app:         app,
		toRemove:    []container{*cont},
		provisioner: s.p,
	}
	context := action.FWContext{Params: []interface{}{args}, Previous: []container{}}
	result, err := provisionRemoveOldUnits.Forward(context)
	c.Assert(err, check.IsNil)
	resultContainers := result.([]container)
	c.Assert(resultContainers, check.DeepEquals, []container{})
	_, err = s.p.getContainer(cont.ID)
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
	unit := cont.asUnit(app)
	app.BindUnit(&unit)
	args := changeUnitsPipelineArgs{
		app:         app,
		toRemove:    []container{*cont},
		provisioner: s.p,
	}
	context := action.FWContext{Params: []interface{}{args}, Previous: []container{}}
	result, err := provisionUnbindOldUnits.Forward(context)
	c.Assert(err, check.IsNil)
	resultContainers := result.([]container)
	c.Assert(resultContainers, check.DeepEquals, []container{})
	c.Assert(app.HasBind(&unit), check.Equals, false)
}

func (s *S) TestFollowLogsAndCommitName(c *check.C) {
	c.Assert(followLogsAndCommit.Name, check.Equals, "follow-logs-and-commit")
}

func (s *S) TestFollowLogsAndCommitForward(c *check.C) {
	err := s.newFakeImage(s.p, "github.com/megamsys/vertice/python", nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("mightyapp", "python", 1)
	nextImgName, err := appNewImageName(app.GetName())
	c.Assert(err, check.IsNil)
	cont := container{AppName: "mightyapp", ID: "myid123", BuildingImage: nextImgName}
	err = cont.create(runContainerActionsArgs{
		app:         app,
		imageID:     "github.com/megamsys/vertice/python",
		commands:    []string{"foo"},
		provisioner: s.p,
	})
	c.Assert(err, check.IsNil)
	buf := safe.NewBuffer(nil)
	args := runContainerActionsArgs{writer: buf, provisioner: s.p}
	context := action.FWContext{Params: []interface{}{args}, Previous: cont}
	imageId, err := followLogsAndCommit.Forward(context)
	c.Assert(err, check.IsNil)
	c.Assert(imageId, check.Equals, "github.com/megamsys/vertice/app-mightyapp:v1")
	c.Assert(buf.String(), check.Not(check.Equals), "")
	var dbCont container
	coll := s.p.collection()
	defer coll.Close()
	err = coll.Find(bson.M{"id": cont.ID}).One(&dbCont)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "not found")
	_, err = s.p.Cluster().InspectContainer(cont.ID)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Matches, "No such container.*")
	err = s.p.Cluster().RemoveImage("github.com/megamsys/vertice/app-mightyapp:v1")
	c.Assert(err, check.IsNil)
}

func (s *S) TestFollowLogsAndCommitForwardNonZeroStatus(c *check.C) {
	err := s.newFakeImage(s.p, "github.com/megamsys/vertice/python", nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	cont := container{AppName: "mightyapp"}
	err = cont.create(runContainerActionsArgs{
		app:         app,
		imageID:     "github.com/megamsys/vertice/python",
		commands:    []string{"foo"},
		provisioner: s.p,
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


func (s *S) TestBindAndHealthcheckName(c *check.C) {
	c.Assert(bindAndHealthcheck.Name, check.Equals, "bind-and-healthcheck")
}


*/
