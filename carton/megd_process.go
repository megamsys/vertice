package carton

import (
	"bytes"
)

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

// CreateProcs represents a command for creating new cartons.
type CreateProcess struct {
	Name string
}

func (s CreateProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("CREATE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s CreateProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Deploy(); err != nil {
			return err
		}
	}
	return nil
}

// DeleteProcs represents a command for delete cartons.
type DestroyProcess struct {
	Name string
}

func (s DestroyProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("DESTROY CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s DestroyProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Destroy(); err != nil {
			return err
		}
	}
	return nil
}

// StartProcs represents a command for starting  cartons.
type StartProcess struct {
	Name string
}

func (s StartProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("START CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

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

func (s StopProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("STOP CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

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

func (s RestartProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("RESTART CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

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

func (s StateupProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("STATEUP CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s StateupProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Stateup(); err != nil {
			return err
		}
	}
	return nil
}
