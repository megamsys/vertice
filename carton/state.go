package carton

import (
	"bytes"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/repository"
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
			return changer.SetState(opts.B, writer, opts.Changed)
		}
	}
	return nil
}

func saveStateData(opts *StateChangeOpts, slog string, duration time.Duration, changeError error) error {
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(opts.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(duration.String(), "green", "", "bold"),
		cmd.Colorfy(slog, "yellow", "", ""))

	if opts.B.Level == provision.BoxSome && opts.B.Repo.Enabled {
		hookId, err := repository.Manager(opts.B.Repo.GetSource()).CreateHook(opts.B.Repo)
		if err != nil {
			return nil
		}

		comp, err := NewComponent(opts.B.Id)
		if err != nil {
			return err
		}

		if err = comp.setDeployData(DeployData{
			Timestamp: time.Now(),
			Duration:  duration,
			HookId:    hookId,
		}); err != nil {
			return err
		}
	}
	return nil
}

func markStatesAsRemoved(boxName string) error {
	//Riak: code to nuke state changes out of a box
	return nil
}
