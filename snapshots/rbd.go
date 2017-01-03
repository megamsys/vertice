package snapshots

import (
  "fmt"
  cmr "github.com/megamsys/vertice/communicator"
  "github.com/megamsys/vertice/carton"
)


var RBD_RUNNER = []string{"CephRbdImageSize"}

type CephRbdRunner struct {
	CephRbdImageSize bool              `json:"cephrbdimagesize"`
	Host             string            `json:"ipaddress"`
	Username         string            `json:"username"`
	Password         string            `json:"password"`
	PrivateKey       string            `json:"privatekey"`
	RbdImage         string            `json:"rbd_image"`
	Inputs           map[string]string `json:"inputs"`
}

type AsmSnaps struct {
	AssemblyId    string
	AccountId     string
	AssemblyName  string
	NumberofSnaps string
	TotalStorage  float64
}

func (i *CephRbdRunner) CephRbd() error {
	runner, _ := cmr.NewUrkRunner(i, i.Inputs)
	if runner != nil {
		if r, ok := runner.(cmr.UrkRunner); ok {
			s, err := r.Run(RBD_RUNNER)
			if err != nil {
				return err
			}
			fmt.Println(s)
		}
	}
	return nil
}

func (i *CephRbdRunner) GetUserSnaps(acts []carton.Account) []AsmSnaps {

  return nil
}
