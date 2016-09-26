package carton

import (
	"bytes"
	"strconv"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	ldb "github.com/megamsys/libgo/db"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/provision"
	"github.com/pivotal-golang/bytefmt"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
	"time"
)

const (
	SNAPSHOTBUCKET = "snapshots"
	DISKSBUCKET    = "disks"
)

type DiskOpts struct {
	B *provision.Box
}

//The grand elephant for megam cloud platform.
type Snaps struct {
	Id         string `json:"snap_id" cql:"snap_id"`
	ImageId    string `json:"image_id" cql:"image_id"`
	OrgId      string `json:"org_id" cql:"org_id"`
	AccountId  string `json:"account_id" cql:"account_id"`
	Name       string `json:"name" cql:"name"`
	AssemblyId string `json:"asm_id" cql:"asm_id"`
	JsonClaz   string `json:"json_claz" cql:"json_claz"`
	CreatedAt  string `json:"created_at" cql:"created_at"`
	Status     string `json:"status" cql:"status"`
}

type Disks struct {
	Id         string `json:"snap_id" cql:"snap_id"`
	DiskId     string `json:"disk_id" cql:"disk_id"`
	OrgId      string `json:"org_id" cql:"org_id"`
	AccountId  string `json:"account_id" cql:"account_id"`
	Name       string `json:"name" cql:"name"`
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
func GetSnap(id string) (*Snaps, error) {
	a := &Snaps{}
	//ops := vdb.ScyllaOptions("snapshot", []string{"Snap_Id"}, []string{"org_id"}, map[string]string{"Id": id, "org_id":"ORG123"})
	ops := ldb.Options{
		TableName:   SNAPSHOTBUCKET,
		Pks:         []string{"snap_id"},
		Ccms:        []string{},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{"snap_id": id},
		CcmsClauses: make(map[string]interface{}),
	}
	if err := ldb.Fetchdb(ops, a); err != nil {
		return nil, err
	}
	log.Debugf("Snaps %v", a)
	return a, nil
}

func (a *Snaps) UpdateSnap(update_fields map[string]interface{}) error {
	ops := ldb.Options{
		TableName:   SNAPSHOTBUCKET,
		Pks:         []string{},
		Ccms:        []string{"id"},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{},
		CcmsClauses: map[string]interface{}{"id": a.Id},
	}
	if err := ldb.Updatedb(ops, update_fields); err != nil {
		return err
	}

	return nil

}

/** A public function which pulls the disks that attached to vm.
and any others we do. **/
func GetDisks(id string) (*Disks, error) {
	d := &Disks{}
	ops := ldb.Options{
		TableName:   DISKSBUCKET,
		Pks:         []string{"id"},
		Ccms:        []string{},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{"id": id},
		CcmsClauses: make(map[string]interface{}),
	}
	if err := ldb.Fetchdb(ops, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (a *Disks) GetDisks() (*[]Disks, error) {
	ds := &[]Disks{}
	d := &Disks{}
	ops := ldb.Options{
		TableName:   DISKSBUCKET,
		Pks:         []string{},
		Ccms:        []string{"account_id", "asm_id"},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{"account_id": a.AccountId, "asm_id": a.AssemblyId},
		CcmsClauses: make(map[string]interface{}),
	}
	if err := ldb.FetchListdb(ops, 10, d, ds); err != nil {
		return nil, err
	}

	return ds, nil
}

func (a *Disks) RemoveDisk() error {
	ops := ldb.Options{
		TableName:   DISKSBUCKET,
		Pks:         []string{"id"},
		Ccms:        []string{},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{"id": a.Id},
		CcmsClauses: map[string]interface{}{},
	}
	if err := ldb.Deletedb(ops, Disks{}); err != nil {
		return err
	}
	return nil
}

func (a *Snaps) RemoveSnap() error {
	ops := ldb.Options{
		TableName:   SNAPSHOTBUCKET,
		Pks:         []string{"id"},
		Ccms:        []string{},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{"id": a.Id},
		CcmsClauses: map[string]interface{}{},
	}
	if err := ldb.Deletedb(ops, Snaps{}); err != nil {
		return err
	}
	return nil
}

func (a *Disks) UpdateDisk(update_fields map[string]interface{}) error {
	ops := ldb.Options{
		TableName:   DISKSBUCKET,
		Pks:         []string{},
		Ccms:        []string{"id"},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{},
		CcmsClauses: map[string]interface{}{"id": a.Id},
	}
	if err := ldb.Updatedb(ops, update_fields); err != nil {
		return err
	}

	return nil

}

//make cartons from snaps.
func (a *Snaps) MkCartons() (Cartons, error) {
	newCs := make(Cartons, 0, 1)
	if len(strings.TrimSpace(a.AssemblyId)) > 1 {
		if ca, err := mkCarton(a.Id, a.AssemblyId); err != nil {
			return nil, err
		} else {
			ca.toBox()                //on success, make a carton2box if BoxLevel is BoxZero
			newCs = append(newCs, ca) //on success append carton
		}
	}
	log.Debugf("Cartons %v", newCs)
	return newCs, nil
}

//make cartons from snaps.
func (d *Disks) MkCartons() (Cartons, error) {
	newCs := make(Cartons, 0, 1)
	if len(strings.TrimSpace(d.AssemblyId)) > 1 {
		if ca, err := mkCarton(d.Id, d.AssemblyId); err != nil {
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
		return strconv.FormatUint(0, 64)
	} else {
		return strconv.FormatUint(cp, 64)
	}
}
