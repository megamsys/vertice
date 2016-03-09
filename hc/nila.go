package hc

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/megamsys/libgo/hc"
	"github.com/megamsys/vertice/meta"
	"gopkg.in/yaml.v2"
)

func init() {
	hc.AddChecker("nilavu", healthCheck)
}

func healthCheck() (interface{}, error) {
	filename := filepath.Join(meta.MC.Home, "nilavu.conf")
	if _, err := os.Stat(filename); err != nil {
		return nil, hc.ErrDisabledComponent
	}
	var m MiniNilavu
	if dat, err := ioutil.ReadFile(filename); err == nil {
		n := nilavu{}
		if err = n.mkNilavu(dat); err != nil {
			return nil, err
		}
		m = mkMiniNilavu(&n)
	} else {
		return nil, err
	}
	return m, nil
}

func (n *nilavu) mkNilavu(data []byte) error {
	if err := yaml.Unmarshal(data, n); err != nil {
		return err
	}
	return nil
}

func mkMiniNilavu(n *nilavu) MiniNilavu {
	return MiniNilavu{
		Api:  n.Defaults.Api,
		Logs: n.Defaults.Logs,
	}
}

type MiniNilavu struct {
	Api  string
	Logs string
}

//a boring struct to parse nilavu.yml
type nilavu struct {
	Defaults struct {
		Api  string
		Logs string
	}
}
