package marketplaces

import (
	"bytes"
	"fmt"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/vertice/provision"
	"io"
	// "strings"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/events"
	"github.com/megamsys/libgo/events/alerts"
	"github.com/megamsys/libgo/pairs"
	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	lw "github.com/megamsys/libgo/writer"
	"github.com/megamsys/vertice/meta"
	"gopkg.in/yaml.v2"
	"time"
)

const (
	APIMARKETPLACES = "/marketplaces"
	UPDATE          = "/update"
)

type apiMarketplaces struct {
	JsonClaz string         `json:"json_claz" cql:"json_claz"`
	Results  []Marketplaces `json:"results" cql:"results"`
}

//Global provisioners set by the subd daemons.
var ProvisionerMap map[string]provision.Provisioner = make(map[string]provision.Provisioner)

// category: rawimage             //Get rawimages build market
// action: rawimage.iso.create
//
// category: marketplaces        // Get Marketplaces build market
// action: marketplaces.iso.customize
//
// category: marketplaces        // Get Marketplaces build market
// action: marketplaces.iso.finished
//
// category: marketplaces       // get marketplaces build market
// action:   marketplaces.image.add

// struct for marketplaces and rawimages
type Marketplaces struct {
	Id         string            `json:"id"`
	AccountId  string            `json:"account_id"`
	ProvidedBy string            `json:"provided_by"`
	Inputs     pairs.JsonPairs   `json:"inputs"`
	Outputs    pairs.JsonPairs   `json:"outputs"`
	Envs       pairs.JsonPairs   `json:"envs"`
	Options    pairs.JsonPairs   `json:"options"`
	CatType    string            `json:"cattype"`
	Flavor     string            `json:"flavor"`
	Image      string            `json:"image"`
	CatOrder   string            `json:"catorder"`
	Plans      map[string]string `json:""plans`
	//	Status       string            `json:"status"`
	JsonClaz string `json:"json_claz"`
}

func NewArgs(email, org string) api.ApiArgs {
	return newArgs(email, org)
}

func newArgs(email, org string) api.ApiArgs {
	return api.ApiArgs{
		Master_Key: meta.MC.MasterKey,
		Url:        meta.MC.Api,
		Email:      email,
		Org_Id:     org,
	}
}

// marketplaces json string
func (s *Marketplaces) String() string {
	if d, err := yaml.Marshal(s); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

func (m *Marketplaces) Get() (*Marketplaces, error) {
	return m.get()
}

func GetMarketplace(email, id string) (*Marketplaces, error) {
	r := new(Marketplaces)
	r.AccountId = email
	r.Id = id
	return r.get()
}

/** A public function which pulls the snapshot for disk save as image.
and any others we do. **/
func (m *Marketplaces) get() (*Marketplaces, error) {
	cl := api.NewClient(newArgs(m.AccountId, ""), APIMARKETPLACES+"/"+m.Id)

	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	res := &apiMarketplaces{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	a := &res.Results[0]
	log.Debugf("Marketplaces %v", a)
	return a, nil
}

/** A public function which pulls all marketplaces items under particular settings_name like vertice,bitnami.**/
func GetBySettingsName(settings_name, email string) ([]Marketplaces, error) {
	cl := api.NewClient(newArgs(email, ""), APIMARKETPLACES+settings_name)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}
	res := &apiMarketplaces{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}

	log.Debugf("Marketplaces of current Assemmbly %v", &res.Results)
	return res.Results, nil
}

/** A public function which pulls all marketplaces items.**/
func (s *Marketplaces) Gets() ([]Marketplaces, error) {
	cl := api.NewClient(newArgs(meta.MC.MasterUser, ""), APIMARKETPLACES)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	res := &apiMarketplaces{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}

	return res.Results, nil
}

func (m *Marketplaces) Update() error {
	return m.update()
}

func (m *Marketplaces) update() error {
	cl := api.NewClient(newArgs(m.AccountId, ""), APIMARKETPLACES+UPDATE)
	if _, err := cl.Post(m); err != nil {
		return err
	}
	return nil
}

func (m *Marketplaces) UpdateStatus(status utils.Status) error {
	lastStatusUpdate := time.Now().Local().Format(time.RFC822)
	i := make(map[string][]string, 2)
	i["lastsuccessstatusupdate"] = []string{lastStatusUpdate}
	i["status"] = []string{status.String()}
	m.Inputs.NukeAndSet(i) //just nuke the matching output key:
	err := m.update()
	if err != nil {
		return err
	}
	return m.trigger_event(status)
}

func (m *Marketplaces) Trigger_event(status utils.Status) error {
	return m.trigger_event(status)
}

func (m *Marketplaces) trigger_event(status utils.Status) error {
	mi := make(map[string]string)
	js := make(pairs.JsonPairs, 0)
	in := make(map[string][]string, 2)
	in["status"] = []string{status.String()}
	in["description"] = []string{status.Description(m.ImageName())}
	js.NukeAndSet(in) //just nuke the matching output key:

	mi[constants.MARKETPLACE_ID] = m.Id
	mi[constants.ACCOUNT_ID] = m.AccountId
	mi[constants.EVENT_TYPE] = status.MkEvent_type()

	newEvent := events.NewMulti(
		[]*events.Event{
			&events.Event{
				AccountsId:  m.AccountId,
				EventAction: alerts.STATUS,
				EventType:   constants.EventUser,
				EventData:   alerts.EventData{M: mi, D: js.ToString()},
				Timestamp:   time.Now().Local(),
			},
		})
	return newEvent.Write()
}

func (m *Marketplaces) rawImageCustomize() error {
	box, err := m.mkBox()
	if err != nil {
		return err
	}
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := lw.LogWriter{Box: box}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	deployer, ok := ProvisionerMap[box.Provider].(provision.MarketPlaceAccess)
	if !ok {
		return fmt.Errorf("cannot provision marketplaces")
	}
	err = deployer.CustomizeImage(box, writer)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	log.Debugf("%s in (%s)\n%s", cmd.Colorfy(box.Name, "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(outBuffer.String(), "yellow", "", ""))
	return nil
}

func (m *Marketplaces) saveImage() error {
	box, err := m.mkBox()
	if err != nil {
		return err
	}
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := lw.LogWriter{Box: box}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	deployer, ok := ProvisionerMap[box.Provider].(provision.MarketPlaceAccess)
	if !ok {
		return fmt.Errorf("cannot provision marketplaces")
	}
	err = deployer.SaveMarketplaceImage(box, writer)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	log.Debugf("%s in (%s)\n%s", cmd.Colorfy(box.Name, "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(outBuffer.String(), "yellow", "", ""))
	return nil
}

func (m *Marketplaces) mkBox() (*provision.Box, error) {
	box := &provision.Box{
		CartonId:   m.Id,
		AccountId:  m.AccountId,
		Name:       m.ImageName(),
		Region:     m.Region(),
		Provider:   m.provider(),
		InstanceId: m.instanceId(),
	}
	return box, nil
}

func (m *Marketplaces) NukeAndSetOutputs(out map[string][]string) error {
	log.Debugf("nuke and set outputs in scylla [%s]", out)
	m.Outputs.NukeAndSet(out) //just nuke the matching output key:
	return m.update()
}

func (s *Marketplaces) ImageName() string {
	return s.Inputs.Match(utils.IMAGE_NAME)
}

func (s *Marketplaces) ImageId() string {
	return s.Outputs.Match(utils.IMAGE_ID)
}

func (s *Marketplaces) RemoveVM() string {
	return s.Inputs.Match("remove_vm")
}

func (s *Marketplaces) RawImageId() string {
	return s.Inputs.Match(utils.RAW_IMAGE_ID)
}

func (s *Marketplaces) instanceId() string {
	return s.Inputs.Match(utils.INSTANCE_ID)
}

func (s *Marketplaces) Region() string {
	return s.Inputs.Match(utils.REGION)
}

func (s *Marketplaces) provider() string {
	return s.Inputs.Match(utils.PROVIDER)
}

func (m *Marketplaces) GetVMCpuCost() string {
	return ""
}

func (m *Marketplaces) GetVMMemoryCost() string {
	return ""
}

func (m *Marketplaces) GetVMHDDCost() string {
	return ""
}
