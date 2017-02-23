package marketplaces

import (
  "fmt"
	"github.com/megamsys/vertice/marketplaces/provision"
	 "bytes"
	 "github.com/megamsys/libgo/cmd"
	 "io"
	// "strings"
	 "time"
	  "gopkg.in/yaml.v2"
	log "github.com/Sirupsen/logrus"
  "github.com/megamsys/libgo/utils"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/pairs"
	"github.com/megamsys/vertice/meta"
	"encoding/json"
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
	Id           string            `json:"id"`
	AccountId    string            `json:"account_id"`
	SettingsName string            `json:"settings_name"`
	Inputs       pairs.JsonPairs   `json:"inputs"`
	Outputs      pairs.JsonPairs   `json:"outputs"`
	Envs         pairs.JsonPairs   `json:"envs"`
	Options      pairs.JsonPairs   `json:"options"`
	CatType      string            `json:"cattype"`
	Flavor       string            `json:"flavor"`
	Image        string            `json:"image"`
	CatOrder     string            `json:"catorder"`
	Plans        map[string]string `json:""plans`
	Status       string            `json:"status"`
	JsonClaz     string            `json:"json_claz"`
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

func (s *Marketplaces) Update() error {
	cl := api.NewClient(newArgs(s.AccountId, ""), APIMARKETPLACES+UPDATE)
	if _, err := cl.Post(s); err != nil {
		return err
	}
	return nil
}

func (m *Marketplaces) rawImageCustomize() error {
	box, err := m.mkBox()
	if err != nil {
		return err
	}
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := LogWriter{Box: box}
	logWriter.Async()
	 defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	deployer, ok := ProvisionerMap[box.Provider].(provision.MarketPlaceAccess)
	if !ok {
		return fmt.Errorf("cannot provision marketplaces")
	}
	err = deployer.CustomiseRawImage(box, writer)
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
	raw := new(RawImages)
	raw.AccountId = m.AccountId
	raw.Id = m.rawImageId()
	raw, err := raw.get()
	if err != nil {
		return nil, err
	}

  box := &provision.Box{
    CartonId: m.Id,
    AccountId: m.AccountId,
		Name: m.ImageName(),
    Region: m.Region(),
    Provider: raw.provider(),
		SourceImage: raw.Name,
  }
  return box, nil
}

func (s *Marketplaces) ImageName() string {
	return s.Inputs.Match(utils.IMAGE_NAME)
}

func (s *Marketplaces) rawImageId() string {
	return s.Inputs.Match(utils.RAW_IMAGE_ID)
}

func (s *Marketplaces) instanceId() string {
	return s.Inputs.Match(utils.INSTANCE_ID)
}

func (s *Marketplaces) Region() string {
	return s.Inputs.Match(utils.REGION)
}
