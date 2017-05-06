package metrix

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/events/alerts"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/carton"
	"time"
)

const (
	OPENNEBULA = "one"
	QUOTA      = "quota"
	ONDEMAND   = "ondemand"
	INSTANCE   = "instance"
	DOCKER     = "docker"
)

type InstanceHandler struct {
	Deployd        bool
	Dockerd        bool
	ContainerUnits map[string]string
	Flavors        map[string]*carton.Flavor
	VMUnits        map[string]string
	SkewsActions   map[string]string
}

func (on *InstanceHandler) Prefix() string {
	return INSTANCE
}

func (i *InstanceHandler) Collect(mc *MetricsCollection) (e error) {
	orgs, e := i.ReadOrgs()
	if e != nil {
		return
	}

	users, e := i.ReadUsers()
	if e != nil {
		return
	}
	// error not handled because < 1.5.2 version dont have flavors
	flvr, e := i.parseFlavors()
	i.Flavors = flvr
	instances, e := i.ParseAssemblies(orgs, users)
	if e != nil {
		return
	}

	if i.Deployd {
		i.CollectMetricsFromStats(mc, instances[OPENNEBULA], ONE_VM_SENSOR)
	}
	if i.Dockerd {
		i.CollectMetricsFromStats(mc, instances[DOCKER], DOCKER_CONTAINER_SENSOR)
	}

	return i.DeductBill(mc)
}

//actually the NewSensor can create trypes based on the event type.
func (i *InstanceHandler) CollectMetricsFromStats(mc *MetricsCollection, amies []carton.Assemblies, sensorType string) {
	for _, h := range amies {
		for _, ay := range h.Assemblys {
			resources := ay.Resources(i.Flavors[ay.FlavorId()])
			sc := NewSensor(sensorType)
			sc.QuotaId = ay.QuotaId()
			sc.AccountId = ay.AccountId
			sc.System = ay.Tosca
			sc.Node = ay.HostName()
			sc.AssemblyId = ay.Id
			sc.AssemblyName = ay.GetFullName()
			sc.AssembliesId = h.Id
			sc.Resources = resources[constants.RESOURCES]
			sc.Source = i.Prefix()
			sc.Message = "vm billing"
			sc.Status = ay.State
			sc.AuditPeriodBeginning = time.Now().Add(-MetricsInterval).Format(time.RFC3339) // time.Unix(h.PStime, 0).String()
			sc.AuditPeriodEnding = time.Now().Format(time.RFC3339)                          // time.Unix(h.PEtime, 0).String()
			sc.AuditPeriodDelta = ""
			sc.addMetric(constants.CPU_COST, resources[constants.CPU_COST], resources[constants.CPU], "delta")
			sc.addMetric(constants.MEMORY_COST, resources[constants.MEMORY_COST], resources[constants.RAM], "delta")
			sc.addMetric(constants.DISK_COST, resources[constants.DISK_COST], resources[constants.STORAGE], "delta")
			sc.CreatedAt = time.Now()
			mc.Add(sc)
		}
	}
	return
}

func (i *InstanceHandler) DeductBill(c *MetricsCollection) (e error) {
	var action alerts.EventAction
	defaultUnits := make(map[string]string)
	for _, mc := range c.Sensors {
		if mc.SensorType == ONE_VM_SENSOR {
			defaultUnits = i.VMUnits
		} else if mc.SensorType == DOCKER_CONTAINER_SENSOR {
			defaultUnits = i.ContainerUnits
		}
		if mc.QuotaId == "" {
			mkBalance(mc, defaultUnits)
		}

		if i.SkewsActions[constants.ENABLED] == constants.TRUE {
			if len(mc.QuotaId) > 0 {
				quota, err := carton.NewQuota(mc.AccountId, mc.QuotaId)
				if err != nil {
					log.Debugf("quota get error : %s", err.Error())
				}
				action = alerts.QUOTA_UNPAID
				if quota.Status == "paid" {
					continue
				}
				i.SkewsActions[constants.SKEWS_TYPE] = "instance.quota.unpaid"
			} else {
				action = alerts.SKEWS_ACTIONS
				i.SkewsActions[constants.SKEWS_TYPE] = "instance.ondemand.bills"
			}
			e = eventSkews(mc, action, i.SkewsActions)
		}

	}
	return
}

func (i *InstanceHandler) parseFlavors() (map[string]*carton.Flavor, error) {
	flavors := make(map[string]*carton.Flavor, 0)
	flvs, e := carton.GetFlavors()
	if e != nil {
		return flavors, e
	}
	for _, f := range flvs {
		flavors[f.Id] = &f
	}
	return flavors, nil
}

func (i *InstanceHandler) ParseAssemblies(orgs []carton.Organization, users map[string]*carton.Account) (map[string][]carton.Assemblies, error) {
	assembly := make(map[string]carton.Assembly, 0)
	one := make([]carton.Assemblies, 0)
	docker := make([]carton.Assemblies, 0)
	asms, e := carton.AssemblyBox()
	if e != nil {
		log.Debugf("%v ", e)
		return nil, e
	}
	for _, ay := range asms {
		assembly[ay.Id] = ay
	}

	for _, org := range orgs {
		usr, ok := users[org.AccountId]
		if ok && !usr.IsAdmin() {
			amies, e := carton.Gets(org.AccountId, org.Id)
			if e != nil {
				log.Debugf("%v ", e)
			} else {
				for _, ays := range amies {
					ays.Assemblys = make(map[string]carton.Assembly, 0)
					asm, ok := assembly[ays.AssemblysId[0]]
					if ok && asm.IsAlive() {
						switch true {
						case asm.IsTopedo():
							ays.Assemblys[asm.Id] = asm
							one = append(one, ays)
						case asm.IsContainer():
							ays.Assemblys[asm.Id] = asm
							docker = append(docker, ays)
						}
					}
				}
				log.Debugf("bill for %s - vms(%v) and containers(%v)", org.AccountId, len(one), len(docker))
			}

		}
	}
	return map[string][]carton.Assemblies{
		OPENNEBULA: one,
		DOCKER:     docker,
	}, nil
}

func (i *InstanceHandler) ReadOrgs() ([]carton.Organization, error) {
	return carton.OrgBox()
}

func (i *InstanceHandler) ReadUsers() (map[string]*carton.Account, error) {
	accounts := make(map[string]*carton.Account, 0)
	acts, e := new(carton.Account).GetUsers()
	if e != nil {
		return accounts, e
	}
	for _, a := range acts {
		accounts[a.Email] = a
	}
	return accounts, nil
}
