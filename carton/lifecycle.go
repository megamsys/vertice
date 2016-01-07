/*
** copyright [2013-2016] [Megam Systems]
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
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/provision"
)

type LifecycleOpts struct {
	B         *provision.Box
	start     time.Time
	logWriter LogWriter
	writer    io.Writer
}

func (cy *LifecycleOpts) setLogger() {
	cy.start = time.Now()
	cy.logWriter = NewLogWriter(cy.B)
	cy.writer = io.MultiWriter(&cy.logWriter)
}

//if the state is in running, started, stopped, restarted then allow it to be lcycled.
// to-do: allow states that ends with "*ing or *ed" that should fix this generically.
func (cy *LifecycleOpts) canCycle() bool {
	return cy.B.Status == provision.StatusRunning ||
		cy.B.Status == provision.StatusStarted ||
		cy.B.Status == provision.StatusStopped ||
		cy.B.Status == provision.StatusUpgraded
}

// Starts  the box.
func Start(cy *LifecycleOpts) error {
	log.Debugf("  start cycle for box (%s, %s)", cy.B.Id, cy.B.GetFullName())
	cy.setLogger()
	defer cy.logWriter.Close()
	if cy.canCycle() {
		if err := ProvisionerMap[cy.B.Provider].Start(cy.B, "", cy.writer); err != nil {
			return err
		}
	}
	fmt.Fprintf(cy.writer, "    start (%s, %s, %s) OK\n", cy.B.GetFullName(), cy.B.Status.String(), time.Since(cy.start))
	return nil
}

// Stops the box
func Stop(cy *LifecycleOpts) error {
	log.Debugf("  stop cycle for box (%s, %s)", cy.B.Id, cy.B.GetFullName())
	cy.setLogger()
	defer cy.logWriter.Close()
	if cy.canCycle() {
		if err := ProvisionerMap[cy.B.Provider].Stop(cy.B, "", cy.writer); err != nil {
			return err
		}
	}
	fmt.Fprintf(cy.writer, "    stop (%s, %s, %s) OK\n", cy.B.GetFullName(), cy.B.Status.String(), time.Since(cy.start))
	return nil
}

// Restart the box.
func Restart(cy *LifecycleOpts) error {
	log.Debugf("  restart cycle for box (%s, %s)", cy.B.Id, cy.B.GetFullName())
	cy.setLogger()
	defer cy.logWriter.Close()
	if cy.canCycle() {
		if err := ProvisionerMap[cy.B.Provider].Restart(cy.B, "", cy.writer); err != nil {
			return err
		}
	}
	fmt.Fprintf(cy.writer, "    restart (%s, %s, %s) OK\n", cy.B.GetFullName(), cy.B.Status.String(), time.Since(cy.start))
	return nil
}
