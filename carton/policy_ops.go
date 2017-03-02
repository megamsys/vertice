/*
** Copyright [2013-2016] [Megam Systems]
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
package carton

import (
	"github.com/megamsys/libgo/pairs"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/repository"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	lw "github.com/megamsys/libgo/writer"
	"bytes"
	"time"
	"io"
)

const (
	NETWORK_ATTACH = "NetworkAttachPolicy"
	NETWORK_DETACH = "NetworkDetachPolicy"
)


type Operations struct {
	Type        string          `json:"operation_type"`
	Description string          `json:"description"`
	Properties  pairs.JsonPairs `json:"properties"`
	Status      string          `json:"status"`
}

func (o *Operations) prepBuildHook() *repository.Hook {
	return &repository.Hook{
		Enabled:  true,
		Token:    o.Properties.Match(repository.TOKEN),
		UserName: o.Properties.Match(repository.USERNAME),
	}
}

func BuildHook(ops []*Operations, opsType string) *repository.Hook {
	for _, o := range ops {
		switch o.Type {
		case opsType:
			return o.prepBuildHook()
		}
	}
	return nil
}

func newOperations(typ string) *Operations {
	return &Operations{
		Type: typ,
		Status: "initialized",
	}
}

func NetworkUpdate(box *provision.Box) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := lw.LogWriter{Box: box}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)

	err := newOperations(NETWORK).network(box, writer)
	elapsed := time.Since(start)
	if err != nil {
		return err
	}
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(box.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(outBuffer.String(), "yellow", "", ""))
	return nil
}

func (o *Operations) network(box *provision.Box, w io.Writer) error {
	if box.IsPolicyOk() {
		if deployer, ok := ProvisionerMap[box.Provider].(provision.Network); ok {
			return deployer.NetworkUpdate(box, w)
		}
  }
	return nil
}
