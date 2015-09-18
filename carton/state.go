// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package carton

import (
	"bytes"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/provision"
)

type StateChangeOpts struct {
	B       *provision.Box
	Changed provision.Status
}

// ChangeState runs a state increment of a machine or a container.
func ChangeState(opts *StateChangeOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := LogWriter{Box: opts.B}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := moveState(opts, writer)
	elapsed := time.Since(start)
	saveErr := saveStateData(opts, outBuffer.String(), elapsed, err)
	if saveErr != nil {
		log.Errorf("WARNING: couldn't save changed data, statechange opts: %#v", opts)
	}
	if err != nil {
		return err
	}
	return nil
}

func moveState(opts *StateChangeOpts, writer io.Writer) error {
	if &opts.Changed != nil {
		if changer, ok := Provisioner.(provision.StateChanger); ok {
			return changer.SetState(opts.B,  writer, opts.Changed)
		}
	}
	return nil
}

func saveStateData(opts *StateChangeOpts, slog string, duration time.Duration, changeError error) error {
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(opts.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(duration.String(), "green", "", "bold"),
		cmd.Colorfy(slog, "yellow", "", ""))
	//if there are state changes to track as follows in outputs: {} then do it here.
	//Riak: code to save the status of a deploy (created.)
	// state :
	//     name:
	//     status:
	return nil
}

func markStatesAsRemoved(boxName string) error {
	//Riak: code to nuke state changes out of a box
	return nil
}
