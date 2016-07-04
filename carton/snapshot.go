package carton

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	ldb "github.com/megamsys/libgo/db"
		"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/provision"
	"strings"
	"gopkg.in/yaml.v2"
	"io"
	"time"
)

type DiskSaveOpts struct {
	B *provision.Box
}

//The grand elephant for megam cloud platform.
type Snaps struct {
	Id          string   `json:"snap_id" cql:"snap_id"`
	OrgId       string   `json:"org_id" cql:"org_id"`
	AccountId   string   `json:"account_id" cql:"account_id"`
	Name        string   `json:"name" cql:"name"`
	AssemblyId  string   `json:"asm_id" cql:"asm_id"`
	JsonClaz    string   `json:"json_claz" cql:"json_claz"`
	CreatedAt   string   `json:"created_at" cql:"created_at"`
}

func (a *Snaps) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

// ChangeState runs a state increment of a machine or a container.
func SaveImage(opts *DiskSaveOpts) error {
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

/** A public function which pulls the snapshot for disk save as image.
and any others we do. **/
func GetSnap(id string) (*Snaps, error) {

	a := &Snaps{}
	//ops := vdb.ScyllaOptions("snapshot", []string{"Snap_Id"}, []string{"org_id"}, map[string]string{"Id": id, "org_id":"ORG123"})
	ops := ldb.Options{
		TableName:   "snapshots",
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
