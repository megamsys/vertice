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
package carton

import (
	"io"
	"bytes"
)

type ReqOperator struct {
}

// NewParser returns a new instance of Parser.
func NewReqOperator(r io.Reader) *ReqOperator {
	return &ReqOperator{}
}

func (p *ReqOperator) AcceptRequest(r *MegdProcessor) error {
	mg := *r
	c, err := p.Get("nil")

	if err != nil {
		return err
	}

	return mg.Process(c)
}

func (p *ReqOperator) Get(cat_id string) (Cartons, error) {
	a, err := Get(cat_id)
	if err != nil {
		return nil, err
	}

	c, err := a.MkCartons()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Statement represents a single operation in Megamd.
type MegdProcessor interface {
	//Name() string
	Process(c Cartons) error
	//Required() ExecutionRequirements
}

// CreateProcs represents a command for creating new cartons.
type CreateProcess struct {
	Name string
}


// String returns a string representation of the create cartons
func (s CreateProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("CREATE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

// Process that creates a cartons
func (s CreateProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Deploy(); err != nil {
			return err
		}
	}
	return nil
}

// DeleteProcs represents a command for delete cartons.
type DeleteProcess struct {
	Name string
}

// String returns a string representation of the delete carton
func (s DeleteProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("DELETE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

// Process that deletes the cartons
func (s DeleteProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Delete(); err != nil {
			return err
		}
	}
	return nil
}

// StartProcs represents a command for starting  cartons.
type StartProcess struct {
	Name string
}

// String returns a string representation of the start procs.
func (s StartProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("START CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

// Process that start procs the cartons
func (s StartProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Start(); err != nil {
			return err
		}
	}
	return nil
}

// StopProcs represents a command for stoping  cartons.
type StopProcess struct {
	Name string
}

// String returns a string representation of the stop ops.
func (s StopProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("STOP CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

// Process that stops the cartons
func (s StopProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Stop(); err != nil {
			return err
		}
	}
	return nil
}

// RestartProcs represents a command for restarting  cartons.
type RestartProcess struct {
	Name string
}

// String returns a string representation of the restart ops.
func (s RestartProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("RESTART CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

// Process that restarts the cartons
func (s RestartProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Restart(); err != nil {
			return err
		}
	}
	return nil
}

// StateupProcess represents a command for restarting  cartons.
type StateupProcess struct {
	Name string
}

// String returns a string representation of the stateup ops.
func (s StateupProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("STATEUP CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

// Process that restarts the cartons
func (s StateupProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Stateup(); err != nil {
			return err
		}
	}
	return nil
}
