package hc

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/megamsys/libgo/hc"
	"github.com/megamsys/megamd/meta"
	"gopkg.in/yaml.v2"
)

func init() {
	hc.AddChecker("nilavu", healthCheck)
}

func healthCheck() (interface{}, error) {
	filename := filepath.Join(meta.MC.Home, "nilavu.yml")
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
		Billings:     n.Defaults.Billings,
		Backup:       n.Defaults.Backup,
		Locale:       n.Defaults.Locale,
		Nilavu:       n.Defaults.Nilavu,
		Riak:         n.Defaults.Riak,
		Api:          n.Defaults.Api,
		Logs:         n.Defaults.Logs,
		Gitlab:       n.Defaults.Gitlab,
		Notification: n.Defaults.Notification,
	}
}

type MiniNilavu struct {
	Billings     string
	Backup       backup
	Locale       string
	Nilavu       string
	Riak         string
	Api          string
	Logs         string
	Gitlab       string
	Notification notification
}

type backup struct {
	Enable   string
	Host     string
	Username string
	Password string
}

type notification struct {
	Use    string
	Enable string
	Smtp   struct {
		Email    string
		Password string
	}
	Mailgun struct {
		Api_key string
		Domain  string
	}
	Hipchat struct {
		Api_key string
		Room    string
	}
}

//a boring struct to parse nilavu.yml
type nilavu struct {
	Defaults struct {
		Billings     string
		Backup       backup
		Locale       string
		Nilavu       string
		Riak         string
		Api          string
		Logs         string
		Gitlab       string
		Notification notification
		Oauth        struct {
			Facebook struct {
				Client_id  string
				Secret_key string
			}
			Github struct {
				Client_id  string
				Secret_key string
			}
			Google struct {
				Client_id  string
				Secret_key string
			}
		}
	}
}
