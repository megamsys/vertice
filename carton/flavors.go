package carton

import (
	"encoding/json"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/pairs"
	constants "github.com/megamsys/libgo/utils"
	"gopkg.in/yaml.v2"
	"strings"
)

const (
	GB = "GB"
)

type Flavor struct {
	Name       string          `json:"name"`
	Id         string          `json:"id"`
	Category   []string        `json:"category"`
	Cpu        string          `json:"cpu"`
	Disk       string          `json:"disk"`
	JsonClaz   string          `json:"json_claz"`
	Price      pairs.JsonPairs `json:"price"`
	Properties pairs.JsonPairs `json:"properties"`
	Ram        string          `json:"ram"`
	Regions    []string        `json:"regions"`
	Status     string          `json:"status"`
	CreatedAt  string          `json:"crerated_at"`
	UpdatedAt  string          `json:"updated_at"`
}

type ApiFlavor struct {
	JsonClaz string   `json:"json_claz"`
	Results  []Flavor `json:"results"`
}

func (f *Flavor) String() string {
	if d, err := yaml.Marshal(f); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

func GetFlavor(email, id string) (*Flavor, error) {
	fs, err := new(Flavor).gets(email, "/flavors/"+id)
	if err != nil {
		return nil, err
	}
	return &fs[0], nil
}

func GetFlavors(email string) ([]Flavor, error) {
	return new(Flavor).gets(email, "/flavors")
}

func (f *Flavor) gets(email, path string) ([]Flavor, error) {
	args := newArgs(email, "")
	cl := api.NewClient(args, path)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	ac := &ApiFlavor{}
	//log.Debugf("Response %s :  (%s)",cmd.Colorfy("[Body]", "green", "", "bold"),string(response))
	err = json.Unmarshal(response, ac)
	if err != nil {
		return nil, err
	}
	return ac.Results, nil
}

func (f *Flavor) getCpushare() string {
	return f.Cpu
}

func (f *Flavor) getMemory() string {
	return f.Ram + GB
}

func (f *Flavor) getSwap() string {
	return ""
}

//The default HDD is 10. we should configure it in the vertice.conf
func (f *Flavor) getHDD() string {
	if len(strings.TrimSpace(f.Disk)) <= 0 {
		return "10" + GB
	}
	return f.Disk + GB
}

func (f *Flavor) GetCpuCost() string {
	return f.Price.Match(constants.CPU_COST_HOUR)
}

func (f *Flavor) GetMemoryCost() string {
	return f.Price.Match(constants.MEMORY_COST_HOUR)
}

func (f *Flavor) GetHDDCost() string {
	return f.Price.Match(constants.DISK_COST_HOUR)
}
