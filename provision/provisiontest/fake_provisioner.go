package provisiontest

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/megamsys/vertice/provision"
)

var uniqueIpCounter int32 = 0

// Fake implementation for carton.
type FakeCarton struct {
	id           string
	name         string
	cartonsId    string
	tosca        string
	imageVersion string
	Compute      provision.BoxCompute
	domainName   string
	provider     string
	publicIp     string
	boxs         *[]provision.Box
	logs         []string
	logMut       sync.Mutex
}

func NewFakeCarton(name, tosca string, lvl provision.BoxLevel, units int) *FakeCarton {
	carton := FakeCarton{
		id:        "CMP010101010101",
		cartonsId: "ASM010101010101",
		name:      name,
		tosca:     tosca,
	}
	b := make([]provision.Box, units)

	for i := 0; i < units; i++ {
		val := atomic.AddInt32(&uniqueIpCounter, 1)
		b[i] = provision.Box{
			Id:           "CMP010101010101",
			CartonsId:    "ASM010101010101",
			CartonId:     "AMS010101010101",
			Level:        lvl,
			Name:         fmt.Sprintf(name, val),
			DomainName:   "megambox.com",
			Tosca:        tosca,
			ImageVersion: "",
			Compute: provision.BoxCompute{
				Cpushare: "0.2",
				Memory:   "512",
				HDD:      "",
			},
			Status:   provision.StatusLaunching,
			Provider: "one",
			PublicIp: "",
		}
	}

	carton.boxs = &b
	return &carton
}

func (a *FakeCarton) Logs() []string {
	a.logMut.Lock()
	defer a.logMut.Unlock()
	logs := make([]string, len(a.logs))
	copy(logs, a.logs)
	return logs
}

func (a *FakeCarton) HasLog(source, unit, message string) bool {
	log := source + unit + message
	a.logMut.Lock()
	defer a.logMut.Unlock()
	for _, l := range a.logs {
		if l == log {
			return true
		}
	}
	return false
}

func (a *FakeCarton) Log(message, source, unit string) error {
	a.logMut.Lock()
	a.logs = append(a.logs, source+unit+message)
	a.logMut.Unlock()
	return nil
}

func (a *FakeCarton) GetName() string {
	return a.name
}

func (a *FakeCarton) Boxs() *[]provision.Box {
	return a.boxs
}

func (a *FakeCarton) GetIp() string {
	return ""
}
