package marketplaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/events"
	"github.com/megamsys/libgo/events/alerts"
	"github.com/megamsys/libgo/pairs"
	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	lw "github.com/megamsys/libgo/writer"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/provision"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
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
	Id          string          `json:"id"`
	AccountId   string          `json:"account_id"`
	ProvidedBy  string          `json:"provided_by"`
	Inputs      pairs.JsonPairs `json:"inputs"`
	Outputs     pairs.JsonPairs `json:"outputs"`
	Envs        pairs.JsonPairs `json:"envs"`
	Options     pairs.JsonPairs `json:"options"`
	AclPolicies pairs.JsonPairs `json:"acl_policies"`
	CatType     string          `json:"cattype"`
	Flavor      string          `json:"flavor"`
	Image       string          `json:"image"`
	CatOrder    string          `json:"catorder"`
	Plans       pairs.JsonPairs `json:"plans"`
	Status      string          `json:"status"`
	Url         string          `json:"url"`
	JsonClaz    string          `json:"json_claz"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
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
	if m.AccountId != "" && m.Id != "" {
		return m.get()
	}
	return nil, fmt.Errorf("Get credentials missing email (%s) id(%s)", m.AccountId, m.Id)
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
	a.AccountId = m.AccountId
	return a, nil
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
	m.Status = status.String()
	err := m.update()
	if err != nil {
		return err
	}
	return m.trigger_event(status)
}

func (m *Marketplaces) UpdateError(status utils.Status, cause error) error {
	m.Status = status.String()
	err := m.update()
	if err != nil {
		return err
	}
	return m.trigger_error_event(status, cause)
}

func (m *Marketplaces) trigger_error_event(status utils.Status, causeof error) error {
	mi := make(map[string]string)
	js := make(pairs.JsonPairs, 0)
	a := make(map[string][]string, 2)
	a["status"] = []string{status.String()}
	a["description"] = []string{status.Description(causeof.Error())}
	js.NukeAndSet(a) //just nuke the matching output key:

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
	log.Debugf("marketplaces  %v", box)
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := lw.LogWriter{Box: box}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	// have to get provisioner from User
	deployer, ok := ProvisionerMap[utils.PROVIDER_ONE].(provision.MarketPlaceAccess)
	if !ok {
		log.Debugf("cannot provision marketplaces  %s", utils.PROVIDER_ONE)
		return fmt.Errorf("cannot provision marketplaces %s", utils.PROVIDER_ONE)
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
		CartonId:    m.Id,
		AccountId:   m.AccountId,
		Name:        m.Flavor,
		CartonName:  m.Flavor,
		Region:      m.Region(),
		Provider:    m.provider(),
		InstanceId:  m.instanceId(),
		Compute:     m.newCompute(),
		StorageType: m.storageType(),
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

func (m *Marketplaces) newCompute() provision.BoxCompute {
	return provision.BoxCompute{
		Cpushare: m.getCpushare(),
		Memory:   m.getMemory(),
		Swap:     m.getSwap(),
		HDD:      m.getHDD(),
	}
}

func (m *Marketplaces) getCpushare() string {
	return m.Inputs.Match(provision.CPU)
}

func (m *Marketplaces) getMemory() string {
	return m.Inputs.Match(provision.RAM)
}

func (m *Marketplaces) getSwap() string {
	return ""
}

//The default HDD is 10. we should configure it in the vertice.conf
func (m *Marketplaces) getHDD() string {
	if len(strings.TrimSpace(m.Inputs.Match(provision.HDD))) <= 0 {
		return "10"
	}
	return m.Inputs.Match(provision.HDD)
}

func (m *Marketplaces) storageType() string {
	return strings.ToLower(m.Inputs.Match(utils.STORAGE_TYPE))
}

func (m *Marketplaces) NukeKeysOutputs(k string) error {
	if len(k) > 0 {
		log.Debugf("nuke keys from outputs in cassandra [%s]", k)
		m.Outputs.NukeKeys(k) //just nuke the matching output key:
		err := m.update()
		if err != nil {
			return err
		}
	} else {
		return provision.ErrNoOutputsFound
	}
	return nil
}
