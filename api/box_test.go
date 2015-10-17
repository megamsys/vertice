package api

import (
//"github.com/megamsys/megamd/config"
//"gopkg.in/check.v1"
)

/*
func (s *S) TestAppLogShouldReturnNotFoundWhenAppDoesNotExist(c *check.C) {
	request, err := http.NewRequest("GET", "/apps/unknown/log/?:app=unknown&lines=10", nil)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	err = appLog(recorder, request, s.token)
	c.Assert(err, check.NotNil)
	e, ok := err.(*errors.HTTP)
	c.Assert(ok, check.Equals, true)
	c.Assert(e.Code, check.Equals, http.StatusNotFound)
	c.Assert(e, check.ErrorMatches, "^App unknown not found.$")
}

func (s *S) TestAppLogReturnsForbiddenIfTheGivenUserDoesNotHaveAccessToTheApp(c *check.C) {
	a := app.App{Name: "lost", Platform: "vougan"}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	defer s.logConn.Logs(a.Name).DropCollection()
	url := fmt.Sprintf("/apps/%s/log/?:app=%s&lines=10", a.Name, a.Name)
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	err = appLog(recorder, request, s.token)
	c.Assert(err, check.NotNil)
	e, ok := err.(*errors.HTTP)
	c.Assert(ok, check.Equals, true)
	c.Assert(e.Code, check.Equals, http.StatusForbidden)
}

func (s *S) TestAppLogReturnsBadRequestIfNumberOfLinesIsMissing(c *check.C) {
	url := "/apps/something/log/?:app=doesntmatter"
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	err = appLog(recorder, request, s.token)
	c.Assert(err, check.NotNil)
	e, ok := err.(*errors.HTTP)
	c.Assert(ok, check.Equals, true)
	c.Assert(e.Code, check.Equals, http.StatusBadRequest)
	c.Assert(e.Message, check.Equals, `Parameter "lines" is mandatory.`)
}

func (s *S) TestAppLogReturnsBadRequestIfNumberOfLinesIsNotAnInteger(c *check.C) {
	url := "/apps/something/log/?:app=doesntmatter&lines=2.34"
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	err = appLog(recorder, request, s.token)
	c.Assert(err, check.NotNil)
	e, ok := err.(*errors.HTTP)
	c.Assert(ok, check.Equals, true)
	c.Assert(e.Code, check.Equals, http.StatusBadRequest)
	c.Assert(e.Message, check.Equals, `Parameter "lines" must be an integer.`)
}

func (s *S) TestAppLogFollowWithPubSub(c *check.C) {
	a := app.App{Name: "lost1", Platform: "zend", Teams: []string{s.team.Name}}
	err := app.CreateApp(&a, s.user)
	c.Assert(err, check.IsNil)
	url := "/apps/something/log/?:app=" + a.Name + "&lines=10&follow=1"
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, check.IsNil)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		recorder := httptest.NewRecorder()
		err := appLog(recorder, request, s.token)
		c.Assert(err, check.IsNil)
		body, err := ioutil.ReadAll(recorder.Body)
		c.Assert(err, check.IsNil)
		splitted := strings.Split(strings.TrimSpace(string(body)), "\n")
		c.Assert(splitted, check.HasLen, 2)
		c.Assert(splitted[0], check.Equals, "[]")
		logs := []app.Applog{}
		err = json.Unmarshal([]byte(splitted[1]), &logs)
		c.Assert(err, check.IsNil)
		c.Assert(logs, check.HasLen, 1)
		c.Assert(logs[0].Message, check.Equals, "x")
	}()
	var listener *app.LogListener
	timeout := time.After(5 * time.Second)
	for listener == nil {
		select {
		case <-timeout:
			c.Fatal("timeout after 5 seconds")
		case <-time.After(50 * time.Millisecond):
		}
		logTracker.Lock()
		for listener = range logTracker.conn {
		}
		logTracker.Unlock()
	}
	factory, err := queue.Factory()
	c.Assert(err, check.IsNil)
	q, err := factory.PubSub(app.LogPubSubQueuePrefix + a.Name)
	c.Assert(err, check.IsNil)
	err = q.Pub([]byte(`{"message": "x"}`))
	c.Assert(err, check.IsNil)
	time.Sleep(500 * time.Millisecond)
	listener.Close()
	wg.Wait()
}


func (s *S) TestAppLogShouldHaveContentType(c *check.C) {
	a := app.App{Name: "lost", Platform: "zend", Teams: []string{s.team.Name}}
	err := app.CreateApp(&a, s.user)
	c.Assert(err, check.IsNil)
	url := fmt.Sprintf("/apps/%s/log/?:app=%s&lines=10", a.Name, a.Name)
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request.Header.Set("Content-Type", "application/json")
	err = appLog(recorder, request, s.token)
	c.Assert(err, check.IsNil)
	c.Assert(recorder.Header().Get("Content-Type"), check.Equals, "application/json")
}


func (s *S) TestAppLogSelectBySource(c *check.C) {
	a := app.App{Name: "lost", Platform: "zend", Teams: []string{s.team.Name}}
	err := app.CreateApp(&a, s.user)
	c.Assert(err, check.IsNil)
	a.Log("mars log", "mars", "")
	a.Log("earth log", "earth", "")
	url := fmt.Sprintf("/apps/%s/log/?:app=%s&source=mars&lines=10", a.Name, a.Name)
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request.Header.Set("Content-Type", "application/json")
	err = appLog(recorder, request, s.token)
	c.Assert(err, check.IsNil)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body, err := ioutil.ReadAll(recorder.Body)
	c.Assert(err, check.IsNil)
	logs := []app.Applog{}
	err = json.Unmarshal(body, &logs)
	c.Assert(err, check.IsNil)
	c.Assert(logs, check.HasLen, 1)
	c.Assert(logs[0].Message, check.Equals, "mars log")
	c.Assert(logs[0].Source, check.Equals, "mars")
	action := rectest.Action{
		Action: "app-log",
		User:   s.user.Email,
		Extra:  []interface{}{"app=" + a.Name, "lines=10", "source=mars"},
	}
	c.Assert(action, rectest.IsRecorded)
}

func (s *S) TestAppLogSelectByUnit(c *check.C) {
	a := app.App{Name: "lost", Platform: "zend", Teams: []string{s.team.Name}}
	err := app.CreateApp(&a, s.user)
	c.Assert(err, check.IsNil)
	a.Log("mars log", "mars", "prospero")
	a.Log("earth log", "earth", "caliban")
	url := fmt.Sprintf("/apps/%s/log/?:app=%s&unit=caliban&lines=10", a.Name, a.Name)
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request.Header.Set("Content-Type", "application/json")
	err = appLog(recorder, request, s.token)
	c.Assert(err, check.IsNil)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body, err := ioutil.ReadAll(recorder.Body)
	c.Assert(err, check.IsNil)
	logs := []app.Applog{}
	err = json.Unmarshal(body, &logs)
	c.Assert(err, check.IsNil)
	c.Assert(logs, check.HasLen, 1)
	c.Assert(logs[0].Message, check.Equals, "earth log")
	c.Assert(logs[0].Source, check.Equals, "earth")
	c.Assert(logs[0].Unit, check.Equals, "caliban")
	action := rectest.Action{
		Action: "app-log",
		User:   s.user.Email,
		Extra:  []interface{}{"app=" + a.Name, "lines=10", "unit=caliban"},
	}
	c.Assert(action, rectest.IsRecorded)
}

func (s *S) TestAppLogSelectByLinesShouldReturnTheLastestEntries(c *check.C) {
	a := app.App{Name: "lost", Platform: "zend", Teams: []string{s.team.Name}}
	err := app.CreateApp(&a, s.user)
	c.Assert(err, check.IsNil)
	now := time.Now()
	coll := s.logConn.Logs(a.Name)
	defer coll.DropCollection()
	for i := 0; i < 15; i++ {
		l := app.Applog{
			Date:    now.Add(time.Duration(i) * time.Hour),
			Message: strconv.Itoa(i),
			Source:  "source",
			AppName: a.Name,
		}
		coll.Insert(l)
	}
	url := fmt.Sprintf("/apps/%s/log/?:app=%s&lines=3", a.Name, a.Name)
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request.Header.Set("Content-Type", "application/json")
	err = appLog(recorder, request, s.token)
	c.Assert(err, check.IsNil)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body, err := ioutil.ReadAll(recorder.Body)
	c.Assert(err, check.IsNil)
	var logs []app.Applog
	err = json.Unmarshal(body, &logs)
	c.Assert(err, check.IsNil)
	c.Assert(logs, check.HasLen, 3)
	c.Assert(logs[0].Message, check.Equals, "12")
	c.Assert(logs[1].Message, check.Equals, "13")
	c.Assert(logs[2].Message, check.Equals, "14")
}

func (s *S) TestAppLogShouldReturnLogByApp(c *check.C) {
	app1 := app.App{Name: "app1", Platform: "zend", Teams: []string{s.team.Name}}
	err := app.CreateApp(&app1, s.user)
	c.Assert(err, check.IsNil)
	app1.Log("app1 log", "source", "")
	app2 := app.App{Name: "app2", Platform: "zend", Teams: []string{s.team.Name}}
	err = app.CreateApp(&app2, s.user)
	c.Assert(err, check.IsNil)
	app2.Log("app2 log", "sourc ", "")
	app3 := app.App{Name: "app3", Platform: "zend", Teams: []string{s.team.Name}}
	err = app.CreateApp(&app3, s.user)
	c.Assert(err, check.IsNil)
	app3.Log("app3 log", "tsuru", "")
	url := fmt.Sprintf("/apps/%s/log/?:app=%s&lines=10", app3.Name, app3.Name)
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request.Header.Set("Content-Type", "application/json")
	err = appLog(recorder, request, s.token)
	c.Assert(err, check.IsNil)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body, err := ioutil.ReadAll(recorder.Body)
	c.Assert(err, check.IsNil)
	logs := []app.Applog{}
	err = json.Unmarshal(body, &logs)
	c.Assert(err, check.IsNil)
	var logged bool
	for _, log := range logs {
		// Should not show the app1 log
		c.Assert(log.Message, check.Not(check.Equals), "app1 log")
		// Should not show the app2 log
		c.Assert(log.Message, check.Not(check.Equals), "app2 log")
		if log.Message == "app3 log" {
			logged = true
		}
	}
	// Should show the app3 log
	c.Assert(logged, check.Equals, true)
}

*/
