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

package carton

import (
	"bytes"
)

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
	Hard bool
}

func (s StopProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("STOP CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s StopProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Stop(s.Hard); err != nil {
			return err
		}
	}
	return nil
}

// SuspendProcess represents a command for suspend  cartons.
type SuspendProcess struct {
	Name string
}

func (s SuspendProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SUSPEND CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s SuspendProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Suspend(); err != nil {
			return err
		}
	}
	return nil
}

// RestartProcs represents a command for restarting  cartons.
type RestartProcess struct {
	Name string
	Hard bool
}

func (s RestartProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("RESTART CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s RestartProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Restart(s.Hard); err != nil {
			return err
		}
	}
	return nil
}

// UpgradeProcs represents a command for starting  cartons.
type UpgradeProcess struct {
	Name string
}

func (s UpgradeProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("UPGRADE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s UpgradeProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Upgrade(); err != nil {
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

// SnapCreateProcess represents a command for delete cartons.
type SnapCreateProcess struct {
	Name string
}

func (s SnapCreateProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SNAPSHOT CREATE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s SnapCreateProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.CreateSnapshot(); err != nil {
			return err
		}
	}
	return nil
}

// DiskSaveProcs represents a command for delete cartons.
type SnapDestroyProcess struct {
	Name string
}

func (s SnapDestroyProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SNAP DELETE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s SnapDestroyProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.DeleteSnapshot(); err != nil {
			return err
		}
	}
	return nil
}

// DiskSaveProcs represents a command for delete cartons.
type SnapRestoreProcess struct {
	Name string
}

func (s SnapRestoreProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SNAP RESTORE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s SnapRestoreProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.RestoreSnapshot(); err != nil {
			return err
		}
	}
	return nil
}

// SnapCreateProcess represents a command for delete cartons.
type SnapSaveAsProcess struct {
	Name string
}

func (s SnapSaveAsProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SNAPSHOT SAVE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s SnapSaveAsProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.SnapshotSaveAs(); err != nil {
			return err
		}
	}
	return nil
}

// ImageCreateProcess represents a command for create backup box.
type ImageCreateProcess struct {
	Name string
}

func (s ImageCreateProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("BACKUP CREATE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s ImageCreateProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.CreateImage(); err != nil {
			return err
		}
	}
	return nil
}

// DiskSaveProcs represents a command for delete cartons.
type ImageDestroyProcess struct {
	Name string
}

func (s ImageDestroyProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("BACKUP DELETE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s ImageDestroyProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.DeleteImage(); err != nil {
			return err
		}
	}
	return nil
}

// DiskAttachProcess represents a command for delete cartons.
type DiskAttachProcess struct {
	Name string
}

func (s DiskAttachProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("DISK ADD CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s DiskAttachProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.AttachDisk(); err != nil {
			return err
		}
	}
	return nil
}

// DiskDetachProcess represents a command for delete cartons.
type DiskDetachProcess struct {
	Name string
}

func (s DiskDetachProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("DISK REMOVE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s DiskDetachProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.DetachDisk(); err != nil {
			return err
		}
	}
	return nil
}

// UpdateProcess represents a command for update operations based on poliecs cartons.
type UpdateNetworkProcess struct {
	Name string
}

func (s UpdateNetworkProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("NETWORK CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s UpdateNetworkProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.NetworkUpdate(); err != nil {
			return err
		}
	}
	return nil
}

// UpgradeProcs represents a command for starting  cartons.
type RunningProcess struct {
	Name             string
	InstancePassword string
}

func (s RunningProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("RUNNING CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s RunningProcess) Process(ca Cartons) error {
	for _, c := range ca {
		if err := c.Running(); err != nil {
			return err
		}
	}
	return nil
}

// UpgradeProcs represents a command for starting  cartons.
type FailureProcess struct {
	Name string
}

func (s FailureProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("FAILURE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s FailureProcess) Process(ca Cartons) error {
	return nil
}
