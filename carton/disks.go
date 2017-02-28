package carton

import (
	"bytes"
	"code.cloudfoundry.org/bytefmt"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/cmd"
	lw "github.com/megamsys/libgo/writer"
	"github.com/megamsys/vertice/provision"
	"io"
	"strconv"
	"strings"
	"time"
)

type DiskOpts struct {
	B *provision.Box
}

type ApiDisks struct {
	JsonClaz string  `json:"json_claz" cql:"json_claz"`
	Results  []Disks `json:"results" cql:"results"`
}

type Disks struct {
	Id         string `json:"id" cql:"id"`
	DiskId     string `json:"disk_id" cql:"disk_id"`
	OrgId      string `json:"org_id" cql:"org_id"`
	AccountId  string `json:"account_id" cql:"account_id"`
	AssemblyId string `json:"asm_id" cql:"asm_id"`
	JsonClaz   string `json:"json_claz" cql:"json_claz"`
	CreatedAt  string `json:"created_at" cql:"created_at"`
	Size       string `json:"size" cql:"size"`
	Status     string `json:"status" cql:"status"`
}

func NewDisk(email, org, assembly, id string) *Disks {
	return &Disks{
		Id:         id,
		OrgId:      org,
		AccountId:  email,
		AssemblyId: assembly,
	}
}

// ChangeState runs a state increment of a machine or a container.
func AttachDisk(opts *DiskOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := lw.LogWriter{Box: opts.B}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := ProvisionerMap[opts.B.Provider].AttachDisk(opts.B, writer)
	elapsed := time.Since(start)

	if err != nil {
		return err
	}
	slog := outBuffer.String()
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(opts.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(slog, "yellow", "", ""))
	return nil
}

// ChangeState runs a state increment of a machine or a container.
func DetachDisk(opts *DiskOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := lw.LogWriter{Box: opts.B}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := ProvisionerMap[opts.B.Provider].DetachDisk(opts.B, writer)
	elapsed := time.Since(start)
	if err != nil {
		return err
	}
	err = destroyDiskData(opts)
	if err != nil {
		return err
	}
	slog := outBuffer.String()
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(opts.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(slog, "yellow", "", ""))
	return nil
}

func destroyDiskData(opts *DiskOpts) error {
	dsk := NewDisk(opts.B.AccountId, opts.B.OrgId, opts.B.CartonId, opts.B.CartonsId)
	err := dsk.RemoveDisk()
	if err != nil {
		return err
	}
	return nil
}

/** A public function which pulls the disks that attached to vm.
and any others we do. **/
func GetDisks(id, email string) (*Disks, error) {
	cl := api.NewClient(newArgs(email, ""), "/disks/show/"+id)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	res := &ApiDisks{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	d := &res.Results[0]
	log.Debugf("Disks %v", d)
	return d, nil
}

func (a *Disks) RemoveDisk() error {
	cl := api.NewClient(newArgs(a.AccountId, a.OrgId), "/disks/"+a.AssemblyId+"/"+a.Id)
	if _, err := cl.Delete(); err != nil {
		return err
	}
	return nil
}

func (d *Disks) UpdateDisk() error {
	cl := api.NewClient(newArgs(d.AccountId, d.OrgId), "/disks/update")
	if _, err := cl.Post(d); err != nil {
		return err
	}
	return nil

}

//make cartons from disks.
func (d *Disks) MkCartons() (Cartons, error) {
	newCs := make(Cartons, 0, 1)
	if len(strings.TrimSpace(d.AssemblyId)) > 1 {
		if ca, err := mkCarton(d.Id, d.AssemblyId, d.AccountId); err != nil {
			return nil, err
		} else {
			ca.toBox()                //on success, make a carton2box if BoxLevel is BoxZero
			newCs = append(newCs, ca) //on success append carton
		}
	}
	log.Debugf("Cartons %v", newCs)
	return newCs, nil
}

func (bc *Disks) NumMemory() string {
	if cp, err := bytefmt.ToMegabytes(strings.Replace(bc.Size, " ", "", -1)); err != nil {
		return strconv.FormatUint(0, 10)
	} else {
		return strconv.FormatUint(cp, 10)
	}
}
