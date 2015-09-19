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
package provision

import (
	"github.com/megamsys/megamd/repository"
	"gopkg.in/yaml.v2"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	CPU = "cpu"
	RAM = "ram"
	HDD = "hdd"

	// BoxSome indicates that there is atleast one box to deploy or delete.
	BoxSome BoxLevel = iota

	// BoxNone indicates that there are no boxes to deploy or delete but its parent can be.
	BoxNone
)

var cnameRegexp = regexp.MustCompile(`^(\*\.)?[a-zA-Z0-9][\w-.]+$`)

// Boxlevel represents the deployment level.
type BoxLevel int

// Boxlog represents a log entry.
type Boxlog struct {
	Date    time.Time
	Message string
	Source  string
	Name    string
	Unit    string
}

type BoxCompute struct {
	Cpushare string
	Memory   string
	Swap     string
	HDD      string
}

func (bc *BoxCompute) numCpushare() int64 {
	if cs, err := strconv.ParseInt(bc.Cpushare, 10, 64); err != nil {
		return 0
	} else {
		return cs
	}
}

func (bc *BoxCompute) numMemory() int64 {
	if cp, err := strconv.ParseInt(bc.Memory, 10, 64); err != nil {
		return 0
	} else {
		return cp
	}

}

func (bc *BoxCompute) numHDD() int64 {
	if cp, err := strconv.ParseInt(bc.HDD, 10, 64); err != nil {
		return 10
	} else {
		return cp
	}
}

func (bc *BoxCompute) String() string {
	return "[" + strings.Join([]string{
		CPU + ":" + bc.Cpushare,
		RAM + ":" + bc.Memory,
		HDD + ":" + bc.HDD},
		",") + " ]"
}

// Box represents a provision unit. Can be a machine, container or anything
// IP-addressable.
type Box struct {
	Id         string
	CartonId   string
	Level      BoxLevel
	Name       string
	DomainName string
	Tosca      string
	Compute    BoxCompute
	Image      string
	Repo       repository.Repo
	Status     Status
	Provider   string
	Commit     string
	Address    *url.URL
	Ip         string
}

func (b *Box) String() string {
	if d, err := yaml.Marshal(b); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

// GetName returns the name of the box.
func (b *Box) GetFullName() string {
	return b.Name + "." + b.DomainName
}

// GetTosca returns the tosca type of the box.
func (b *Box) GetTosca() string {
	return b.Tosca
}

// GetIp returns the Unit.IP.
func (b *Box) GetIp() string {
	return b.Ip
}

// Available returns true if the unit is available. It will return true
// whenever the unit itself is available, even when the application process is
// not.
func (b *Box) Available() bool {
	return b.Status == StatusDeploying ||
		b.Status == StatusCreating ||
		b.Status == StatusError
}

// AddCName adds a CName to box. It updates the attribute,
// calls the SetCName function on the provisioner and saves
// the box in the database, returning an error when it cannot save the change
// in the database or add the CName on the provisioner.
func (b *Box) AddCName(cnames ...string) error {
	/*	for _, cname := range cnames {
		if cname != "" && !cnameRegexp.MatchString(cname) {
			return stderr.New("Invalid cname")
		}

		if s, ok := Provisioner.(provision.CNameManager); ok {
			if err := s.SetCName(app, cname); err != nil {
				return err
			}
		}
		//Riak: append the ip/cname in the component.
		//here (or) can be handled as an action.
	}*/
	return nil
}

func (b *Box) RemoveCName(cnames ...string) error {
	/*for _, cname := range cnames {
		count := 0
		for _, appCname := range app.CName {
			if cname == appCname {
				count += 1
			}
		}
		if count == 0 {
			return stderr.New("cname not exists!")
		}
		if s, ok := Provisioner.(provision.CNameManager); ok {
			if err := s.UnsetCName(app, cname); err != nil {
				return err
			}
		}
		//Riak: append the ip/cname in the component available in the box.
		//or handle it as an action
		if err != nil {
			return err
		}
	}*/
	return nil
}

func (box *Box) GetRouter() (string, error) {
	return "route53", nil //dns.LoadConfig()
}

// Log adds a log message to the app. Specifying a good source is good so the
// user can filter where the message come from.
func (box *Box) Log(message, source, unit string) error {
	messages := strings.Split(message, "\n")
	logs := make([]interface{}, 0, len(messages))
	for _, msg := range messages {
		if msg != "" {
			bl := Boxlog{
				Date:    time.Now().In(time.UTC),
				Message: msg,
				Source:  source,
				Name:    box.Name,
				Unit:    box.Id,
			}
			logs = append(logs, bl)
		}
	}
	if len(logs) > 0 {
		//notify(bl.Name, logs)
	}
	return nil
}
