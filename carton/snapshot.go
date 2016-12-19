package carton

import (
	"bytes"
	"strconv"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/pairs"
	"github.com/megamsys/vertice/provision"
	"github.com/pivotal-golang/bytefmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"encoding/json"
	"strings"
	"time"
)

const (
	SNAPSHOTBUCKET = "snapshots"
	DISKSBUCKET    = "disks"
	ACCOUNTID      = "account_id"
	ASSEMBLYID     = "asm_id"
)

type DiskOpts struct {
	B *provision.Box
}

type ApiSnaps struct {
	JsonClaz   string `json:"json_claz" cql:"json_claz"`
	Results    []Snaps  `json:"results" cql:"results"`
}

type ApiDisks struct {
	JsonClaz   string `json:"json_claz" cql:"json_claz"`
	Results    []Disks  `json:"results" cql:"results"`
}


//The grand elephant for megam cloud platform.
type Snaps struct {
	Id         string `json:"id" cql:"id"`
	ImageId    string `json:"image_id" cql:"image_id"`
	OrgId      string `json:"org_id" cql:"org_id"`
	AccountId  string `json:"account_id" cql:"account_id"`
	Name       string `json:"name" cql:"name"`
	AssemblyId string `json:"asm_id" cql:"asm_id"`
	JsonClaz   string `json:"json_claz" cql:"json_claz"`
	CreatedAt  string `json:"created_at" cql:"created_at"`
	Status     string `json:"status" cql:"status"`
	Inputs     pairs.JsonPairs `json:"inputs" cql:"inputs"`
	Outputs    pairs.JsonPairs `json:"inputs" cql:"inputs"`
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

func (a *Snaps) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

// ChangeState runs a state increment of a machine or a container.
func SaveImage(opts *DiskOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := LogWriter{Box: opts.B}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := ProvisionerMap[opts.B.Provider].SaveImage(opts.B, writer)
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
func DeleteImage(opts *DiskOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := LogWriter{Box: opts.B}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := ProvisionerMap[opts.B.Provider].DeleteImage(opts.B, writer)
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
func AttachDisk(opts *DiskOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := LogWriter{Box: opts.B}
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
	logWriter := LogWriter{Box: opts.B}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := ProvisionerMap[opts.B.Provider].DetachDisk(opts.B, writer)
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

/** A public function which pulls the snapshot for disk save as image.
and any others we do. **/
func GetSnap(id , email string) (*Snaps, error) {
	args := newArgs(email, "")
	args.Path = "/snapshots/" + id
	cl := api.NewClient(args)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}
	htmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	res := &ApiSnaps{}
	err = json.Unmarshal(htmlData, res)
	if err != nil {
		return nil, err
	}
	a := &res.Results[0]
	log.Debugf("Snaps %v", a)
	return a, nil
}

func (s *Snaps) UpdateSnap() error {
	args := newArgs(s.AccountId, s.OrgId)
	args.Path = "/snapshots/update"
	cl := api.NewClient(args)
	if _, err := cl.Post(s); err != nil {
		return err
	}
	return nil

}
/** A public function which pulls the disks that attached to vm.
and any others we do. **/
func GetDisks(id, email string) (*Disks, error) {
	args := newArgs(email,"")
	args.Path = "/disks" + id
	cl := api.NewClient(args)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}
	htmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	res := &ApiDisks{}
	err = json.Unmarshal(htmlData, res)
	if err != nil {
		return nil, err
	}
	d := &res.Results[0]
	log.Debugf("Disks %v", d)
	return d, nil
}

func (a *Disks) RemoveDisk() error {
	args := newArgs(a.AccountId, a.OrgId)
	args.Path = "/disks/" + a.Id
	cl := api.NewClient(args)
	if	_, err := cl.Delete(); err != nil {
		return err
	}
	return nil
}

func (a *Snaps) RemoveSnap() error {
	args := newArgs(a.AccountId, a.OrgId)
	args.Path = "/snapshots/" + a.Id
	cl := api.NewClient(args)
	if	_, err := cl.Delete(); err != nil {
		return err
	}
	return nil
}

func (d *Disks) UpdateDisk() error {
	args := newArgs(d.AccountId, d.OrgId)
	args.Path = "/disks/update"
	cl := api.NewClient(args)
	if _, err := cl.Post(d); err != nil {
		return err
	}
	return nil

}

//make cartons from snaps.
func (a *Snaps) MkCartons() (Cartons, error) {
	newCs := make(Cartons, 0, 1)
	if len(strings.TrimSpace(a.AssemblyId)) > 1 {
		if ca, err := mkCarton(a.Id, a.AssemblyId, a.AccountId); err != nil {
			return nil, err
		} else {
			ca.toBox()                //on success, make a carton2box if BoxLevel is BoxZero
			newCs = append(newCs, ca) //on success append carton
		}
	}
	log.Debugf("Cartons %v", newCs)
	return newCs, nil
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
