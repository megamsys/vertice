package carton

import (
	"bytes"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/utils"
	lw "github.com/megamsys/libgo/writer"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/repository"
	_ "github.com/megamsys/vertice/repository/github"
	_ "github.com/megamsys/vertice/repository/gitlab"
)

type StateChangeOpts struct {
	B       *provision.Box
	Changed utils.Status
}

// ChangeState runs a state increment of a machine or a container.
func ChangeState(opts *StateChangeOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := lw.LogWriter{Box: opts.B}
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
		if changer, ok := ProvisionerMap[opts.B.Provider].(provision.StateChanger); ok {
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

	if opts.B.Level == provision.BoxSome && opts.B.Repo.IsEnabled() {
		hookId, err := repository.Manager(opts.B.Repo.GetSource()).CreateHook(opts.B.Repo)
		if err != nil {
			return nil
		}

		comp, err := NewComponent(opts.B.Id, opts.B.AccountId, "")
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
	//Scylla: code to nuke state changes out of a box
	return nil
}
