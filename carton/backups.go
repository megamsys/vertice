package carton

import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/pairs"
	"github.com/megamsys/libgo/utils"
	lw "github.com/megamsys/libgo/writer"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/provision"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
	"time"
)

const (
	APIBACKUPS   = "/backups/"
	BACKUPS_SHOW = "/backups/show/"
	UPDATE       = "update"
	DELETE       = "delete/"
	ACCOUNTID    = "account_id"
	ASSEMBLYID   = "asm_id"
	PUBLIC_URL   = "public_url"
)

type ApiBackups struct {
	JsonClaz string    `json:"json_claz" cql:"json_claz"`
	Results  []Backups `json:"results" cql:"results"`
}

//The grand elephant for megam cloud platform.
type Backups struct {
	Id         string          `json:"id" cql:"id"`
	ImageId    string          `json:"image_id" cql:"image_id"`
	OrgId      string          `json:"org_id" cql:"org_id"`
	AccountId  string          `json:"account_id" cql:"account_id"`
	Name       string          `json:"name" cql:"name"`
	AssemblyId string          `json:"asm_id" cql:"asm_id"`
	JsonClaz   string          `json:"json_claz" cql:"json_claz"`
	CreatedAt  string          `json:"created_at" cql:"created_at"`
	Status     string          `json:"status" cql:"status"`
	Tosca      string          `json:"tosca_type" cql:"tosca_type"`
	Labels     pairs.JsonPairs `json:"labels" cql:"labels"`
	Inputs     pairs.JsonPairs `json:"inputs" cql:"inputs"`
	Outputs    pairs.JsonPairs `json:"outputs" cql:"outputs"`
}

func (s *Backups) String() string {
	if d, err := yaml.Marshal(s); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

// ChangeState runs a state increment of a machine or a container.
func CreateImage(opts *DiskOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := lw.LogWriter{Box: opts.B}
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
	logWriter := lw.LogWriter{Box: opts.B}
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

/** A public function which pulls the backup for disk save as image.
and any others we do. **/
func GetBackup(id, email string) (*Backups, error) {
	cl := api.NewClient(newArgs(email, ""), BACKUPS_SHOW+id)

	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	res := &ApiBackups{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	a := &res.Results[0]
	return a, nil
}

/** A public function which pulls the snapshot for disk save as image.
and any others we do. **/
func (s *Backups) GetBox() ([]Backups, error) {
	cl := api.NewClient(newArgs(meta.MC.MasterUser, ""), "/admin"+APIBACKUPS)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	res := &ApiBackups{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}

	return res.Results, nil
}

func (s *Backups) UpdateBackup() error {
	cl := api.NewClient(newArgs(s.AccountId, s.OrgId), APIBACKUPS+UPDATE)
	if _, err := cl.Post(s); err != nil {
		return err
	}
	return nil

}

func (s *Backups) RemoveBackup() error {
	cl := api.NewClient(newArgs(s.AccountId, s.OrgId), APIBACKUPS+s.AssemblyId+"/"+s.Id)
	if _, err := cl.Delete(); err != nil {
		return err
	}
	return nil
}

//make cartons from backups.
func (a *Backups) MkCartons() (Cartons, error) {
	newCs := make(Cartons, 0, 1)
	if len(strings.TrimSpace(a.AssemblyId)) > 1 && a.Tosca != utils.BACKUP_NEW {
		if ca, err := mkCarton(a.Id, a.AssemblyId, a.AccountId); err != nil {
			return nil, err
		} else {
			ca.toBox()                //on success, make a carton2box if BoxLevel is BoxZero
			newCs = append(newCs, ca) //on success append carton
		}
	} else {
		ca, err := a.mkCarton()
		if err != nil {
			return nil, err
		}
		ca.toBox()
		newCs = append(newCs, ca)
	}

	log.Debugf("Cartons %v", newCs)
	return newCs, nil
}

func (a *Backups) mkCarton() (*Carton, error) {
	act, err := new(Account).get(newArgs(a.AccountId, ""))
	if err != nil {
		return nil, err
	}
	b := make([]provision.Box, 0, 0)
	return &Carton{
		CartonsId:    a.Id,
		OrgId:        a.OrgId,
		Name:         a.Name,
		Tosca:        a.Tosca,
		AccountId:    a.AccountId,
		Authority:    act.States.Authority,
		ImageVersion: a.imageVersion(),
		Provider:     a.provider(),
		Region:       a.region(),
		ImageName:    a.imageName(),
		PublicUrl:    a.publicUrl(),
		Boxes:        &b,
		Status:       utils.Status(a.Status),
	}, nil
}

func (s *Backups) Sizeof() string {
	return s.Outputs.Match("image_size")
}

func (b *Backups) region() string {
	return b.Inputs.Match(REGION)
}

func (b *Backups) provider() string {
	p := b.Inputs.Match(utils.PROVIDER)
	if p != "" {
		return p
	}
	return utils.PROVIDER_ONE
}

func (b *Backups) imageVersion() string {
	return b.Inputs.Match(IMAGE_VERSION)
}

func (b *Backups) publicUrl() string {
	return b.Inputs.Match(PUBLIC_URL)
}

func (b *Backups) imageName() string {
	p := b.Inputs.Match(BACKUPNAME)
	if p != "" {
		return p
	}
	return b.Name
}
