package carton

import (
	"bytes"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/provision"
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
	err := Provisioner.Destroy(opts.B, writer)
	elapsed := time.Since(start)
	saveErr := saveDestroyedData(opts, outBuffer.String(), elapsed, err)
	if saveErr != nil {
		log.Errorf("WARNING: couldn't save destroyed data, destroy opts: %#v", opts)
	}
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

	return nil
}
