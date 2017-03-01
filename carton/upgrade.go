package carton

import (
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	lw "github.com/megamsys/libgo/writer"
	"github.com/megamsys/vertice/provision"
)

type Upgradeable struct {
	B             *provision.Box
	w             io.Writer
	ShouldRestart bool
}

func NewUpgradeable(box *provision.Box) *Upgradeable {
	u := &Upgradeable{
		B:             box,
		ShouldRestart: true,
	}
	u.register()
	return u
}

func (u *Upgradeable) register() {}

// this is for CI/BIND for docker only
func (u *Upgradeable) Upgrade() error {
	logWriter := lw.NewLogWriter(u.B)
	defer logWriter.Close()
	writer := io.MultiWriter(&logWriter)
	err := u.operateBox(writer)
	if err != nil {
		return err
	}
	return nil
}

func (u *Upgradeable) operateBox(writer io.Writer) error {
	u.w = writer
	start := time.Now()
	/*opsRan, err := upgrade.Run(upgrade.RunArgs{
		Name:   u.B.GetFullName(),
		O:      u.B.Operations,
		Writer: writer,
		Force:  false,
	})
	if err != nil {
		return err
	}*/

	elapsed := time.Since(start)

	if err := u.saveData(elapsed); err != nil {
		log.Errorf("WARNING: couldn't save ops data, ops opts: %#v", u)
		return err
	}

	if !u.ShouldRestart {
		return nil
	}

	return nil
}

func (u *Upgradeable) saveData(duration time.Duration) error {
	return nil
}
