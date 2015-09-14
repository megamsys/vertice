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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"sync/atomic"
	"time"

)

func (s *S) TestShouldBeRegistered(c *check.C) {
	p, err := provision.Get("one")
	c.Assert(err, check.IsNil)
	c.Assert(p, check.FitsTypeOf, &oneProvisioner{})
}

func (s *S) TestProvisionerProvision(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	err := s.p.Provision(app)
	c.Assert(err, check.IsNil)
	c.Assert(routertest.FakeRouter.HasBackend("myapp"), check.Equals, true)
}

func (s *S) TestProvisionerRestart(c *check.C) {
	app := provisiontest.NewFakeApp("almah", "static", 1)
	customData := map[string]interface{}{
		"procfile": "web: python web.py\nworker: python worker.py\n",
	}
	cont1, err := s.newContainer(&newContainerOpts{
		AppName:         app.GetName(),
		ProcessName:     "web",
		ImageCustomData: customData,
		Image:           "github.com/megamsys/megamd/app-" + app.GetName(),
	}, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(cont1)
	cont2, err := s.newContainer(&newContainerOpts{
		AppName:         app.GetName(),
		ProcessName:     "worker",
		ImageCustomData: customData,
		Image:           "github.com/megamsys/megamd/app-" + app.GetName(),
	}, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(cont2)
	err = s.p.Start(app, "")
	c.Assert(err, check.IsNil)
	dockerContainer, err := s.p.Cluster().InspectContainer(cont1.ID)
	c.Assert(err, check.IsNil)
	c.Assert(dockerContainer.State.Running, check.Equals, true)
	dockerContainer, err = s.p.Cluster().InspectContainer(cont2.ID)
	c.Assert(err, check.IsNil)
	c.Assert(dockerContainer.State.Running, check.Equals, true)
	err = s.p.Restart(app, "", nil)
	c.Assert(err, check.IsNil)
	dbConts, err := s.p.listAllContainers()
	c.Assert(err, check.IsNil)
	c.Assert(dbConts, check.HasLen, 2)
	c.Assert(dbConts[0].ID, check.Not(check.Equals), cont1.ID)
	c.Assert(dbConts[0].AppName, check.Equals, app.GetName())
	c.Assert(dbConts[0].Status, check.Equals, provision.StatusStarting.String())
	c.Assert(dbConts[1].ID, check.Not(check.Equals), cont2.ID)
	c.Assert(dbConts[1].AppName, check.Equals, app.GetName())
	c.Assert(dbConts[1].Status, check.Equals, provision.StatusStarting.String())
	dockerContainer, err = s.p.Cluster().InspectContainer(dbConts[0].ID)
	c.Assert(err, check.IsNil)
	c.Assert(dockerContainer.State.Running, check.Equals, true)
	expectedIP := dockerContainer.NetworkSettings.IPAddress
	expectedPort := dockerContainer.NetworkSettings.Ports["8888/tcp"][0].HostPort
	c.Assert(dbConts[0].IP, check.Equals, expectedIP)
	c.Assert(dbConts[0].HostPort, check.Equals, expectedPort)
}

func (s *S) TestProvisionerRestartProcess(c *check.C) {
	app := provisiontest.NewFakeApp("almah", "static", 1)
	customData := map[string]interface{}{
		"procfile": "web: python web.py\nworker: python worker.py\n",
	}
	cont1, err := s.newContainer(&newContainerOpts{
		AppName:         app.GetName(),
		ProcessName:     "web",
		ImageCustomData: customData,
		Image:           "github.com/megamsys/megamd/app-" + app.GetName(),
	}, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(cont1)
	cont2, err := s.newContainer(&newContainerOpts{
		AppName:         app.GetName(),
		ProcessName:     "worker",
		ImageCustomData: customData,
		Image:           "github.com/megamsys/megamd/app-" + app.GetName(),
	}, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(cont2)
	err = s.p.Start(app, "")
	c.Assert(err, check.IsNil)
	dockerContainer, err := s.p.Cluster().InspectContainer(cont1.ID)
	c.Assert(err, check.IsNil)
	c.Assert(dockerContainer.State.Running, check.Equals, true)
	dockerContainer, err = s.p.Cluster().InspectContainer(cont2.ID)
	c.Assert(err, check.IsNil)
	c.Assert(dockerContainer.State.Running, check.Equals, true)
	err = s.p.Restart(app, "web", nil)
	c.Assert(err, check.IsNil)
	dbConts, err := s.p.listAllContainers()
	c.Assert(err, check.IsNil)
	c.Assert(dbConts, check.HasLen, 2)
	c.Assert(dbConts[0].ID, check.Equals, cont2.ID)
	c.Assert(dbConts[1].ID, check.Not(check.Equals), cont1.ID)
	c.Assert(dbConts[1].AppName, check.Equals, app.GetName())
	c.Assert(dbConts[1].Status, check.Equals, provision.StatusStarting.String())
	dockerContainer, err = s.p.Cluster().InspectContainer(dbConts[1].ID)
	c.Assert(err, check.IsNil)
	c.Assert(dockerContainer.State.Running, check.Equals, true)
	expectedIP := dockerContainer.NetworkSettings.IPAddress
	expectedPort := dockerContainer.NetworkSettings.Ports["8888/tcp"][0].HostPort
	c.Assert(dbConts[1].IP, check.Equals, expectedIP)
	c.Assert(dbConts[1].HostPort, check.Equals, expectedPort)
}

func (s *S) stopContainers(endpoint string, n uint) <-chan bool {
	ch := make(chan bool)
	go func() {
		defer close(ch)
		client, err := docker.NewClient(endpoint)
		if err != nil {
			return
		}
		for n > 0 {
			opts := docker.ListContainersOptions{All: false}
			containers, err := client.ListContainers(opts)
			if err != nil {
				return
			}
			if len(containers) > 0 {
				for _, cont := range containers {
					if cont.ID != "" {
						client.StopContainer(cont.ID, 1)
						n--
					}
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return ch
}

func (s *S) TestDeploy(c *check.C) {
	stopCh := s.stopContainers(s.server.URL(), 1)
	defer func() { <-stopCh }()
	err := s.newFakeImage(s.p, "github.com/megamsys/megamd/python:latest", nil)
	c.Assert(err, check.IsNil)
	a := app.App{
		Name:     "otherapp",
		Platform: "python",
	}
	err = s.storage.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	repository.Manager().CreateRepository(a.Name, nil)
	s.p.Provision(&a)
	defer s.p.Destroy(&a)
	w := safe.NewBuffer(make([]byte, 2048))
	var serviceBodies []string
	rollback := s.addServiceInstance(c, a.Name, nil, func(w http.ResponseWriter, r *http.Request) {
		data, _ := ioutil.ReadAll(r.Body)
		serviceBodies = append(serviceBodies, string(data))
		w.WriteHeader(http.StatusOK)
	})
	defer rollback()
	customData := map[string]interface{}{
		"procfile": "web: python myapp.py",
	}
	err = saveImageCustomData("github.com/megamsys/megamd/app-"+a.Name+":v1", customData)
	c.Assert(err, check.IsNil)
	err = app.Deploy(app.App{
		App:          &a,
		Version:      "master",
		Commit:       "123",
		OutputStream: w,
	})
	c.Assert(err, check.IsNil)
	units, err := a.Units()
	c.Assert(err, check.IsNil)
	c.Assert(units, check.HasLen, 1)
	c.Assert(serviceBodies, check.HasLen, 1)
	c.Assert(serviceBodies[0], check.Matches, ".*unit-host="+units[0].Ip)
}

func (s *S) TestDeployErasesOldImages(c *check.C) {
	config.Set("docker:image-history-size", 1)
	defer config.Unset("docker:image-history-size")
	stopCh := s.stopContainers(s.server.URL(), 3)
	defer func() { <-stopCh }()
	err := s.newFakeImage(s.p, "github.com/megamsys/megamd/python:latest", nil)
	c.Assert(err, check.IsNil)
	a := app.App{
		Name:     "appdeployimagetest",
		Platform: "python",
	}
	err = s.storage.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	repository.Manager().CreateRepository(a.Name, nil)
	err = s.p.Provision(&a)
	c.Assert(err, check.IsNil)
	defer s.p.Destroy(&a)
	w := safe.NewBuffer(make([]byte, 2048))
	customData := map[string]interface{}{
		"procfile": "web: python myapp.py",
	}
	err = saveImageCustomData("github.com/megamsys/megamd/app-"+a.Name+":v1", customData)
	c.Assert(err, check.IsNil)
	err = saveImageCustomData("github.com/megamsys/megamd/app-"+a.Name+":v2", customData)
	c.Assert(err, check.IsNil)
	err = app.Deploy(app.App{
		App:          &a,
		Version:      "master",
		Commit:       "123",
		OutputStream: w,
	})
	c.Assert(err, check.IsNil)
	imgs, err := s.p.Cluster().ListImages(docker.ListImagesOptions{All: true})
	c.Assert(err, check.IsNil)
	c.Assert(imgs, check.HasLen, 2)
	c.Assert(imgs[0].RepoTags, check.HasLen, 1)
	c.Assert(imgs[1].RepoTags, check.HasLen, 1)
	expected := []string{"github.com/megamsys/megamd/app-appdeployimagetest:v1", "github.com/megamsys/megamd/python:latest"}
	got := []string{imgs[0].RepoTags[0], imgs[1].RepoTags[0]}
	sort.Strings(got)
	c.Assert(got, check.DeepEquals, expected)
	err = app.Deploy(app.App{
		App:          &a,
		Version:      "master",
		Commit:       "123",
		OutputStream: w,
	})
	c.Assert(err, check.IsNil)
	imgs, err = s.p.Cluster().ListImages(docker.ListImagesOptions{All: true})
	c.Assert(err, check.IsNil)
	c.Assert(imgs, check.HasLen, 2)
	c.Assert(imgs[0].RepoTags, check.HasLen, 1)
	c.Assert(imgs[1].RepoTags, check.HasLen, 1)
	got = []string{imgs[0].RepoTags[0], imgs[1].RepoTags[0]}
	sort.Strings(got)
	expected = []string{"github.com/megamsys/megamd/app-appdeployimagetest:v2", "github.com/megamsys/megamd/python:latest"}
	c.Assert(got, check.DeepEquals, expected)
}

func (s *S) TestDeployErasesOldImagesIfFailed(c *check.C) {
	config.Set("docker:image-history-size", 1)
	defer config.Unset("docker:image-history-size")
	err := s.newFakeImage(s.p, "github.com/megamsys/megamd/python:latest", nil)
	c.Assert(err, check.IsNil)
	a := app.App{
		Name:     "appdeployimagetest",
		Platform: "python",
	}
	err = s.storage.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	err = s.p.Provision(&a)
	c.Assert(err, check.IsNil)
	defer s.p.Destroy(&a)
	s.server.CustomHandler("/containers/create", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		var result docker.Config
		err := json.Unmarshal(data, &result)
		if err == nil {
			if result.Image == "github.com/megamsys/megamd/app-appdeployimagetest:v1" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		s.server.DefaultHandler().ServeHTTP(w, r)
	}))
	w := safe.NewBuffer(make([]byte, 2048))
	err = app.Deploy(app.App{
		App:          &a,
		Version:      "master",
		Commit:       "123",
		OutputStream: w,
	})
	c.Assert(err, check.NotNil)
	imgs, err := s.p.Cluster().ListImages(docker.ListImagesOptions{All: true})
	c.Assert(err, check.IsNil)
	c.Assert(imgs, check.HasLen, 1)
	c.Assert(imgs[0].RepoTags, check.HasLen, 1)
	c.Assert("github.com/megamsys/megamd/python:latest", check.Equals, imgs[0].RepoTags[0])
}

func (s *S) TestImageDeploy(c *check.C) {
	err := s.newFakeImage(s.p, "github.com/megamsys/megamd/app-otherapp:v1", nil)
	c.Assert(err, check.IsNil)
	err = appendAppImageName("otherapp", "github.com/megamsys/megamd/app-otherapp:v1")
	c.Assert(err, check.IsNil)
	a := app.App{
		Name:     "otherapp",
		Platform: "python",
	}
	err = s.storage.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	s.p.Provision(&a)
	defer s.p.Destroy(&a)
	w := safe.NewBuffer(make([]byte, 2048))
	err = app.Deploy(app.App{
		App:          &a,
		OutputStream: w,
		Image:        "github.com/megamsys/megamd/app-otherapp:v1",
	})
	c.Assert(err, check.IsNil)
	units, err := a.Units()
	c.Assert(err, check.IsNil)
	c.Assert(units, check.HasLen, 1)
}

func (s *S) TestImageDeployInvalidImage(c *check.C) {
	a := app.App{
		Name:     "otherapp",
		Platform: "python",
	}
	err := s.storage.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	s.p.Provision(&a)
	defer s.p.Destroy(&a)
	w := safe.NewBuffer(make([]byte, 2048))
	err = app.Deploy(app.App{
		App:          &a,
		OutputStream: w,
		Image:        "github.com/megamsys/megamd/app-otherapp:v1",
	})
	c.Assert(err, check.ErrorMatches, "invalid image for app otherapp: github.com/megamsys/megamd/app-otherapp:v1")
	units, err := a.Units()
	c.Assert(err, check.IsNil)
	c.Assert(units, check.HasLen, 0)
}

func (s *S) TestProvisionerDestroy(c *check.C) {
	cont, err := s.newContainer(nil, nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp(cont.AppName, "python", 1)
	unit := cont.AsUnit(app)
	app.BindUnit(&unit)
	s.p.Provision(app)
	err = s.p.Destroy(app)
	c.Assert(err, check.IsNil)
	coll := s.p.Collection()
	defer coll.Close()
	count, err := coll.Find(bson.M{"appname": cont.AppName}).Count()
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, 0)
	c.Assert(routertest.FakeRouter.HasBackend("myapp"), check.Equals, false)
	c.Assert(app.HasBind(&unit), check.Equals, false)
}

func (s *S) TestProvisionerDestroyRemovesImage(c *check.C) {
	var registryRequests []*http.Request
	registryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registryRequests = append(registryRequests, r)
		w.WriteHeader(http.StatusOK)
	}))
	defer registryServer.Close()
	registryURL := strings.Replace(registryServer.URL, "http://", "", 1)
	config.Set("docker:registry", registryURL)
	defer config.Unset("docker:registry")
	stopCh := s.stopContainers(s.server.URL(), 1)
	defer func() { <-stopCh }()
	a := app.App{
		Name:     "mydoomedapp",
		Platform: "python",
	}
	err := s.storage.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	repository.Manager().CreateRepository(a.Name, nil)
	s.p.Provision(&a)
	defer s.p.Destroy(&a)
	w := safe.NewBuffer(make([]byte, 2048))
	customData := map[string]interface{}{
		"procfile": "web: python myapp.py",
	}
	err = saveImageCustomData(registryURL+"/github.com/megamsys/megamd/app-"+a.Name+":v1", customData)
	c.Assert(err, check.IsNil)
	err = app.Deploy(app.App{
		App:          &a,
		Version:      "master",
		Commit:       "123",
		OutputStream: w,
	})
	c.Assert(err, check.IsNil)
	err = s.p.Destroy(&a)
	c.Assert(err, check.IsNil)
	coll := s.p.Collection()
	defer coll.Close()
	count, err := coll.Find(bson.M{"appname": a.Name}).Count()
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, 0)
	c.Assert(routertest.FakeRouter.HasBackend(a.Name), check.Equals, false)
	c.Assert(registryRequests, check.HasLen, 1)
	c.Assert(registryRequests[0].Method, check.Equals, "DELETE")
	c.Assert(registryRequests[0].URL.Path, check.Equals, "/v1/repositories/github.com/megamsys/megamd/app-mydoomedapp:v1/")
	imgs, err := s.p.Cluster().ListImages(docker.ListImagesOptions{All: true})
	c.Assert(err, check.IsNil)
	c.Assert(imgs, check.HasLen, 1)
	c.Assert(imgs[0].RepoTags, check.HasLen, 1)
	c.Assert(imgs[0].RepoTags[0], check.Equals, registryURL+"/github.com/megamsys/megamd/python:latest")
}

func (s *S) TestProvisionerDestroyEmptyUnit(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 0)
	s.p.Provision(app)
	err := s.p.Destroy(app)
	c.Assert(err, check.IsNil)
}

func (s *S) TestProvisionerDestroyRemovesRouterBackend(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 0)
	err := s.p.Provision(app)
	c.Assert(err, check.IsNil)
	err = s.p.Destroy(app)
	c.Assert(err, check.IsNil)
	c.Assert(routertest.FakeRouter.HasBackend("myapp"), check.Equals, false)
}

func (s *S) TestProvisionerAddr(c *check.C) {
	cont, err := s.newContainer(nil, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(cont)
	app := provisiontest.NewFakeApp(cont.AppName, "python", 1)
	addr, err := s.p.Addr(app)
	c.Assert(err, check.IsNil)
	r, err := getRouterForApp(app)
	c.Assert(err, check.IsNil)
	expected, err := r.Addr(cont.AppName)
	c.Assert(err, check.IsNil)
	c.Assert(addr, check.Equals, expected)
}

func (s *S) TestProvisionerAddUnits(c *check.C) {
	err := s.newFakeImage(s.p, "github.com/megamsys/megamd/app-myapp", nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("myapp", "python", 0)
	app.Deploys = 1
	s.p.Provision(app)
	defer s.p.Destroy(app)
	_, err = s.newContainer(&newContainerOpts{AppName: app.GetName()}, nil)
	c.Assert(err, check.IsNil)
	units, err := s.p.AddUnits(app, 3, "web", nil)
	c.Assert(err, check.IsNil)
	coll := s.p.Collection()
	defer coll.Close()
	defer coll.RemoveAll(bson.M{"appname": app.GetName()})
	c.Assert(units, check.HasLen, 3)
	count, err := coll.Find(bson.M{"appname": app.GetName()}).Count()
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, 4)
}

func (s *S) TestProvisionerAddUnitsInvalidProcess(c *check.C) {
	err := s.newFakeImage(s.p, "github.com/megamsys/megamd/app-myapp", nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("myapp", "python", 0)
	app.Deploys = 1
	s.p.Provision(app)
	defer s.p.Destroy(app)
	_, err = s.newContainer(&newContainerOpts{AppName: app.GetName()}, nil)
	c.Assert(err, check.IsNil)
	_, err = s.p.AddUnits(app, 3, "bogus", nil)
	c.Assert(err, check.FitsTypeOf, provision.InvalidProcessError{})
	c.Assert(err, check.ErrorMatches, `process error: no command declared in Procfile for process "bogus"`)
}

func (s *S) TestProvisionerAddUnitsWithErrorDoesntLeaveLostUnits(c *check.C) {
	callCount := 0
	s.server.CustomHandler("/containers/create", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.server.DefaultHandler().ServeHTTP(w, r)
	}))
	defer s.server.CustomHandler("/containers/create", s.server.DefaultHandler())
	err := s.newFakeImage(s.p, "github.com/megamsys/megamd/python:latest", nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("myapp", "python", 0)
	s.p.Provision(app)
	defer s.p.Destroy(app)
	coll := s.p.Collection()
	defer coll.Close()
	coll.Insert(container.Container{ID: "c-89320", AppName: app.GetName(), Version: "a345fe", Image: "github.com/megamsys/megamd/python:latest"})
	defer coll.RemoveId(bson.M{"id": "c-89320"})
	_, err = s.p.AddUnits(app, 3, "web", nil)
	c.Assert(err, check.NotNil)
	count, err := coll.Find(bson.M{"appname": app.GetName()}).Count()
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, 1)
}

func (s *S) TestProvisionerAddZeroUnits(c *check.C) {
	err := s.newFakeImage(s.p, "github.com/megamsys/megamd/python:latest", nil)
	c.Assert(err, check.IsNil)
	app := provisiontest.NewFakeApp("myapp", "python", 0)
	app.Deploys = 1
	s.p.Provision(app)
	defer s.p.Destroy(app)
	coll := s.p.Collection()
	defer coll.Close()
	coll.Insert(container.Container{ID: "c-89320", AppName: app.GetName(), Version: "a345fe", Image: "github.com/megamsys/megamd/python:latest"})
	defer coll.RemoveId(bson.M{"id": "c-89320"})
	units, err := s.p.AddUnits(app, 0, "web", nil)
	c.Assert(units, check.IsNil)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "Cannot add 0 units")
}

func (s *S) TestProvisionerRemoveUnitsFailRemoveOldRoute(c *check.C) {
	a1 := app.App{Name: "impius", Teams: []string{"github.com/megamsys/megamdteam", "nodockerforme"}, Pool: "pool1"}
	cont1 := container.Container{ID: "1", Name: "impius1", AppName: a1.Name, ProcessName: "web", HostAddr: "url0", HostPort: "1"}
	cont2 := container.Container{ID: "2", Name: "mirror1", AppName: a1.Name, ProcessName: "worker", HostAddr: "url0", HostPort: "2"}
	cont3 := container.Container{ID: "3", Name: "dedication1", AppName: a1.Name, ProcessName: "web", HostAddr: "url0", HostPort: "3"}
	err := s.storage.Apps().Insert(a1)
	c.Assert(err, check.IsNil)
	defer s.storage.Apps().RemoveAll(bson.M{"name": a1.Name})
	p := provision.Pool{Name: "pool1", Teams: []string{
		"github.com/megamsys/megamdteam",
		"nodockerforme",
	}}
	o := provision.AddPoolOptions{Name: p.Name}
	err = provision.AddPool(o)
	c.Assert(err, check.IsNil)
	err = provision.AddTeamsToPool(p.Name, p.Teams)
	defer provision.RemovePool(p.Name)
	contColl := s.p.Collection()
	defer contColl.Close()
	err = contColl.Insert(
		cont1, cont2, cont3,
	)
	c.Assert(err, check.IsNil)
	scheduler := segregatedScheduler{provisioner: s.p}
	s.p.storage = &cluster.MapStorage{}
	clusterInstance, err := cluster.New(&scheduler, s.p.storage)
	c.Assert(err, check.IsNil)
	s.p.cluster = clusterInstance
	s.p.scheduler = &scheduler
	err = clusterInstance.Register(cluster.Node{
		Address:  "http://url0:1234",
		Metadata: map[string]string{"pool": "pool1"},
	})
	c.Assert(err, check.IsNil)
	customData := map[string]interface{}{
		"procfile": "web: python myapp.py",
	}
	err = saveImageCustomData("github.com/megamsys/megamd/app-"+a1.Name, customData)
	c.Assert(err, check.IsNil)
	papp := provisiontest.NewFakeApp(a1.Name, "python", 0)
	s.p.Provision(papp)
	conts := []container.Container{cont1, cont2, cont3}
	units := []provision.Unit{cont1.AsUnit(papp), cont2.AsUnit(papp), cont3.AsUnit(papp)}
	for i := range conts {
		err = routertest.FakeRouter.AddRoute(a1.Name, conts[i].Address())
		c.Assert(err, check.IsNil)
		err = papp.BindUnit(&units[i])
		c.Assert(err, check.IsNil)
	}
	routertest.FakeRouter.FailForIp(conts[2].Address().String())
	err = s.p.RemoveUnits(papp, 2, "web", nil)
	c.Assert(err, check.ErrorMatches, "error removing routes, units weren't removed: Forced failure")
	_, err = s.p.GetContainer(conts[0].ID)
	c.Assert(err, check.IsNil)
	_, err = s.p.GetContainer(conts[1].ID)
	c.Assert(err, check.IsNil)
	_, err = s.p.GetContainer(conts[2].ID)
	c.Assert(err, check.IsNil)
	c.Assert(s.p.scheduler.ignoredContainers, check.IsNil)
	c.Assert(routertest.FakeRouter.HasRoute(a1.Name, conts[0].Address().String()), check.Equals, true)
	c.Assert(routertest.FakeRouter.HasRoute(a1.Name, conts[1].Address().String()), check.Equals, true)
	c.Assert(routertest.FakeRouter.HasRoute(a1.Name, conts[2].Address().String()), check.Equals, true)
	c.Assert(papp.HasBind(&units[0]), check.Equals, true)
	c.Assert(papp.HasBind(&units[1]), check.Equals, true)
	c.Assert(papp.HasBind(&units[2]), check.Equals, true)
}

func (s *S) TestProvisionerRemoveUnitsNotFound(c *check.C) {
	err := s.p.RemoveUnits(nil, 1, "web", nil)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "remove units: app should not be nil")
}

func (s *S) TestProvisionerRemoveUnitsZeroUnits(c *check.C) {
	err := s.p.RemoveUnits(provisiontest.NewFakeApp("something", "python", 0), 0, "web", nil)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "cannot remove zero units")
}

func (s *S) TestProvisionerRemoveUnitsTooManyUnits(c *check.C) {
	a1 := app.App{Name: "impius", Teams: []string{"github.com/megamsys/megamdteam", "nodockerforme"}, Pool: "pool1"}
	cont1 := container.Container{ID: "1", Name: "impius1", AppName: a1.Name, ProcessName: "web"}
	cont2 := container.Container{ID: "2", Name: "mirror1", AppName: a1.Name, ProcessName: "web"}
	cont3 := container.Container{ID: "3", Name: "dedication1", AppName: a1.Name, ProcessName: "web"}
	err := s.storage.Apps().Insert(a1)
	c.Assert(err, check.IsNil)
	defer s.storage.Apps().RemoveAll(bson.M{"name": a1.Name})
	p := provision.Pool{Name: "pool1", Teams: []string{
		"github.com/megamsys/megamdteam",
		"nodockerforme",
	}}
	o := provision.AddPoolOptions{Name: p.Name}
	err = provision.AddPool(o)
	c.Assert(err, check.IsNil)
	err = provision.AddTeamsToPool(p.Name, p.Teams)
	defer provision.RemovePool(p.Name)
	contColl := s.p.Collection()
	defer contColl.Close()
	err = contColl.Insert(
		cont1, cont2, cont3,
	)
	c.Assert(err, check.IsNil)
	defer contColl.RemoveAll(bson.M{"name": bson.M{"$in": []string{cont1.Name, cont2.Name, cont3.Name}}})
	scheduler := segregatedScheduler{provisioner: s.p}
	s.p.storage = &cluster.MapStorage{}
	clusterInstance, err := cluster.New(&scheduler, s.p.storage)
	s.p.scheduler = &scheduler
	s.p.cluster = clusterInstance
	c.Assert(err, check.IsNil)
	err = clusterInstance.Register(cluster.Node{
		Address:  "http://url0:1234",
		Metadata: map[string]string{"pool": "pool1"},
	})
	c.Assert(err, check.IsNil)
	customData := map[string]interface{}{
		"procfile": "web: python myapp.py",
	}
	err = saveImageCustomData("github.com/megamsys/megamd/app-"+a1.Name, customData)
	papp := provisiontest.NewFakeApp(a1.Name, "python", 0)
	s.p.Provision(papp)
	c.Assert(err, check.IsNil)
	err = s.p.RemoveUnits(papp, 4, "web", nil)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "cannot remove 4 units from process \"web\", only 3 available")
}

func (s *S) TestProvisionerSetUnitStatus(c *check.C) {
	err := s.newFakeImage(s.p, "github.com/megamsys/megamd/python:latest", nil)
	c.Assert(err, check.IsNil)
	opts := newContainerOpts{Status: provision.StatusStarted.String(), AppName: "someapp"}
	container, err := s.newContainer(&opts, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(container)
	err = s.p.SetUnitStatus(provision.Unit{Name: container.ID, AppName: container.AppName}, provision.StatusError)
	c.Assert(err, check.IsNil)
	container, err = s.p.GetContainer(container.ID)
	c.Assert(err, check.IsNil)
	c.Assert(container.Status, check.Equals, provision.StatusError.String())
}

func (s *S) TestProvisionerSetUnitStatusUpdatesIp(c *check.C) {
	err := s.storage.Apps().Insert(&app.App{Name: "myawesomeapp"})
	c.Assert(err, check.IsNil)
	err = s.newFakeImage(s.p, "github.com/megamsys/megamd/python:latest", nil)
	c.Assert(err, check.IsNil)
	opts := newContainerOpts{Status: provision.StatusStarted.String(), AppName: "myawesomeapp"}
	container, err := s.newContainer(&opts, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(container)
	container.IP = "xinvalidx"
	coll := s.p.Collection()
	defer coll.Close()
	err = coll.Update(bson.M{"id": container.ID}, container)
	c.Assert(err, check.IsNil)
	err = s.p.SetUnitStatus(provision.Unit{Name: container.ID, AppName: container.AppName}, provision.StatusStarted)
	c.Assert(err, check.IsNil)
	container, err = s.p.GetContainer(container.ID)
	c.Assert(err, check.IsNil)
	c.Assert(container.Status, check.Equals, provision.StatusStarted.String())
	c.Assert(container.IP, check.Matches, `\d+.\d+.\d+.\d+`)
}

func (s *S) TestProvisionerSetUnitStatusWrongApp(c *check.C) {
	err := s.newFakeImage(s.p, "github.com/megamsys/megamd/python:latest", nil)
	c.Assert(err, check.IsNil)
	opts := newContainerOpts{Status: provision.StatusStarted.String(), AppName: "someapp"}
	container, err := s.newContainer(&opts, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(container)
	err = s.p.SetUnitStatus(provision.Unit{Name: container.ID, AppName: container.AppName + "a"}, provision.StatusError)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "wrong app name")
	container, err = s.p.GetContainer(container.ID)
	c.Assert(err, check.IsNil)
	c.Assert(container.Status, check.Equals, provision.StatusStarted.String())
}

func (s *S) TestProvisionSetUnitStatusNoAppName(c *check.C) {
	err := s.newFakeImage(s.p, "github.com/megamsys/megamd/python:latest", nil)
	c.Assert(err, check.IsNil)
	opts := newContainerOpts{Status: provision.StatusStarted.String(), AppName: "someapp"}
	container, err := s.newContainer(&opts, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(container)
	err = s.p.SetUnitStatus(provision.Unit{Name: container.ID}, provision.StatusError)
	c.Assert(err, check.IsNil)
	container, err = s.p.GetContainer(container.ID)
	c.Assert(err, check.IsNil)
	c.Assert(container.Status, check.Equals, provision.StatusError.String())
}

func (s *S) TestProvisionerSetUnitStatusUnitNotFound(c *check.C) {
	err := s.p.SetUnitStatus(provision.Unit{Name: "mycontainer", AppName: "myapp"}, provision.StatusError)
	c.Assert(err, check.Equals, provision.ErrUnitNotFound)
}

func (s *S) TestProvisionSetCName(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend("myapp")
	addr, _ := url.Parse("http://127.0.0.1")
	routertest.FakeRouter.AddRoute("myapp", addr)
	cname := "mycname.com"
	err := s.p.SetCName(app, cname)
	c.Assert(err, check.IsNil)
	c.Assert(routertest.FakeRouter.HasCName(cname), check.Equals, true)
	c.Assert(routertest.FakeRouter.HasRoute(cname, addr.String()), check.Equals, true)
}

func (s *S) TestProvisionUnsetCName(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	routertest.FakeRouter.AddBackend("myapp")
	addr, _ := url.Parse("http://127.0.0.1")
	routertest.FakeRouter.AddRoute("myapp", addr)
	cname := "mycname.com"
	err := s.p.SetCName(app, cname)
	c.Assert(err, check.IsNil)
	c.Assert(routertest.FakeRouter.HasCName(cname), check.Equals, true)
	c.Assert(routertest.FakeRouter.HasRoute(cname, addr.String()), check.Equals, true)
	err = s.p.UnsetCName(app, cname)
	c.Assert(err, check.IsNil)
	c.Assert(routertest.FakeRouter.HasCName(cname), check.Equals, false)
	c.Assert(routertest.FakeRouter.HasRoute(cname, addr.String()), check.Equals, false)
}

func (s *S) TestProvisionerIsCNameManager(c *check.C) {
	var _ provision.CNameManager = &dockerProvisioner{}
}

func (s *S) TestProvisionerPlatformAdd(c *check.C) {
	var requests []*http.Request
	server, err := testing.NewServer("127.0.0.1:0", nil, func(r *http.Request) {
		requests = append(requests, r)
	})
	c.Assert(err, check.IsNil)
	defer server.Stop()
	config.Set("docker:registry", "localhost:3030")
	defer config.Unset("docker:registry")
	var p dockerProvisioner
	err = p.Initialize()
	c.Assert(err, check.IsNil)
	p.cluster, _ = cluster.New(nil, &cluster.MapStorage{},
		cluster.Node{Address: server.URL()})
	args := make(map[string]string)
	args["dockerfile"] = "http://localhost/Dockerfile"
	err = p.PlatformAdd("test", args, bytes.NewBuffer(nil))
	c.Assert(err, check.IsNil)
	c.Assert(requests, check.HasLen, 3)
	c.Assert(requests[0].URL.Path, check.Equals, "/build")
	queryString := requests[0].URL.Query()
	c.Assert(queryString.Get("t"), check.Equals, platformImageName("test"))
	c.Assert(queryString.Get("remote"), check.Equals, "http://localhost/Dockerfile")
	c.Assert(requests[1].URL.Path, check.Equals, "/images/localhost:3030/github.com/megamsys/megamd/test:latest/json")
	c.Assert(requests[2].URL.Path, check.Equals, "/images/localhost:3030/github.com/megamsys/megamd/test/push")
}

func (s *S) TestProvisionerPlatformRemove(c *check.C) {
	registryServer := httptest.NewServer(nil)
	defer registryServer.Close()
	u, _ := url.Parse(registryServer.URL)
	config.Set("docker:registry", u.Host)
	defer config.Unset("docker:registry")
	var requests []*http.Request
	server, err := testing.NewServer("127.0.0.1:0", nil, func(r *http.Request) {
		requests = append(requests, r)
	})
	c.Assert(err, check.IsNil)
	defer server.Stop()
	var p dockerProvisioner
	err = p.Initialize()
	c.Assert(err, check.IsNil)
	p.cluster, _ = cluster.New(nil, &cluster.MapStorage{},
		cluster.Node{Address: server.URL()})
	var buf bytes.Buffer
	err = p.PlatformAdd("test", map[string]string{"dockerfile": "http://localhost/Dockerfile"}, &buf)
	c.Assert(err, check.IsNil)
	err = p.PlatformRemove("test")
	c.Assert(err, check.IsNil)
	c.Assert(requests, check.HasLen, 4)
	c.Assert(requests[3].Method, check.Equals, "DELETE")
	c.Assert(requests[3].URL.Path, check.Matches, "/images/[^/]+")
}

func (s *S) TestRunRestartAfterHooks(c *check.C) {
	a := &app.App{Name: "myrestartafterapp"}
	customData := map[string]interface{}{
		"hooks": map[string]interface{}{
			"restart": map[string]interface{}{
				"after": []string{"cmd1", "cmd2"},
			},
		},
	}
	err := saveImageCustomData("github.com/megamsys/megamd/python:latest", customData)
	c.Assert(err, check.IsNil)
	err = s.storage.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	opts := newContainerOpts{AppName: a.Name}
	container, err := s.newContainer(&opts, nil)
	c.Assert(err, check.IsNil)
	defer s.removeTestContainer(container)
	var reqBodies [][]byte
	s.server.CustomHandler("/containers/"+container.ID+"/exec", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		reqBodies = append(reqBodies, data)
		s.server.DefaultHandler().ServeHTTP(w, r)
	}))
	defer container.Remove(s.p)
	var buf bytes.Buffer
	err = s.p.runRestartAfterHooks(container, &buf)
	c.Assert(err, check.IsNil)
	c.Assert(buf.String(), check.Equals, "")
	c.Assert(reqBodies, check.HasLen, 2)
	var req1, req2 map[string]interface{}
	err = json.Unmarshal(reqBodies[0], &req1)
	c.Assert(err, check.IsNil)
	err = json.Unmarshal(reqBodies[1], &req2)
	c.Assert(err, check.IsNil)
	c.Assert(req1, check.DeepEquals, map[string]interface{}{
		"AttachStdout": true,
		"AttachStderr": true,
		"Cmd":          []interface{}{"/bin/bash", "-lc", "cmd1"},
		"Container":    container.ID,
		"User":         "root",
	})
	c.Assert(req2, check.DeepEquals, map[string]interface{}{
		"AttachStdout": true,
		"AttachStderr": true,
		"Cmd":          []interface{}{"/bin/bash", "-lc", "cmd2"},
		"Container":    container.ID,
		"User":         "root",
	})
}

func (s *S) TestDryMode(c *check.C) {
	err := s.newFakeImage(s.p, "github.com/megamsys/megamd/app-myapp", nil)
	c.Assert(err, check.IsNil)
	appInstance := provisiontest.NewFakeApp("myapp", "python", 0)
	defer s.p.Destroy(appInstance)
	s.p.Provision(appInstance)
	imageId, err := appCurrentImageName(appInstance.GetName())
	c.Assert(err, check.IsNil)
	_, err = addContainersWithHost(&changeUnitsPipelineArgs{
		toHost:      "127.0.0.1",
		toAdd:       map[string]*containersToAdd{"web": {Quantity: 5}},
		app:         appInstance,
		imageId:     imageId,
		provisioner: s.p,
	})
	c.Assert(err, check.IsNil)
	newProv, err := s.p.dryMode(nil)
	c.Assert(err, check.IsNil)
	contsNew, err := newProv.listAllContainers()
	c.Assert(err, check.IsNil)
	c.Assert(contsNew, check.HasLen, 5)
}

func (s *S) TestMetricEnvs(c *check.C) {
	err := bs.SaveEnvs(bs.EnvMap{}, bs.PoolEnvMap{
		"mypool": bs.EnvMap{
			"METRICS_BACKEND":      "LOGSTASH",
			"METRICS_LOGSTASH_URI": "localhost:2222",
		},
	})
	c.Assert(err, check.IsNil)
	appInstance := &app.App{
		Name: "impius",
		Pool: "mypool",
	}
	envs := s.p.MetricEnvs(appInstance)
	expected := map[string]string{
		"METRICS_LOGSTASH_URI": "localhost:2222",
		"METRICS_BACKEND":      "LOGSTASH",
	}
	c.Assert(envs, check.DeepEquals, expected)
}
