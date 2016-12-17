package carton

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/vertice/provision"
	"io"
	"time"
)

type DestroyOpts struct {
	B *provision.Box
}

// ChangeState runs a state increment of a machine or a container.
func Destroy(opts *DestroyOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := LogWriter{Box: opts.B}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := ProvisionerMap[opts.B.Provider].Destroy(opts.B, writer)
	elapsed := time.Since(start)
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(opts.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(outBuffer.String(), "yellow", "", ""))
	// saveErr := saveDestroyedData(opts, outBuffer.String(), elapsed, err)
	// if saveErr != nil {
	// 	log.Errorf("WARNING: couldn't save destroyed data, destroy opts: %#v", opts)
	// }
	if err != nil {
		return err
	}
	return nil
}

func saveDestroyedData(opts *DestroyOpts, slog string, duration time.Duration, destroyError error) error {
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(opts.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(duration.String(), "green", "", "bold"),
		cmd.Colorfy(slog, "yellow", "", ""))
	if destroyError == nil {
		markDeploysAsRemoved(opts)
	}
	return nil
}

func markDeploysAsRemoved(opts *DestroyOpts) {
	removedAssemblys := make([]string, 1)

	if _, err := NewAssembly(opts.B.CartonId,opts.B.AccountId,opts.B.OrgId); err == nil {
		removedAssemblys[0] = opts.B.CartonId
	}

	if opts.B.Level == provision.BoxSome {
		if comp, err := NewComponent(opts.B.Id,opts.B.AccountId,opts.B.OrgId); err == nil {
			comp.Delete(opts.B.AccountId,opts.B.OrgId)
		}
	}

	if asms, err := Get(opts.B.CartonsId, opts.B.AccountId); err == nil {
		asms.Delete(opts.B.CartonsId, opts.B.AccountId, removedAssemblys)
	}

}
