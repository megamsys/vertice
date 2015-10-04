package route53

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/karlentwistle/route53"
	"github.com/megamsys/megamd/router"
	"github.com/megamsys/megamd/subd/dns"
)

const (
	routerName = "route53"
	CREATE     = "CREATE"
	DELETE     = "DELETE"
)

func init() {
	router.Register(routerName, createRouter)
}

type route53Router struct {
	cname  string
	ip     string
	client route53.AccessIdentifiers
	zone   route53.HostedZone
	choped string
}

func createRouter(name string) (router.Router, error) {
	vRouter := route53Router{
		client: route53.AccessIdentifiers{
			AccessKey: dns.R53.AccessKey,
			SecretKey: dns.R53.SecretKey,
		},
	}
	log.Debugf("%s ready", routerName)
	return vRouter, nil
}

//cname is the fullname eg: test.megambox.com
func (r route53Router) SetCName(cname, ip string) error {
	r.cname = cname
	if len(strings.TrimSpace(r.cname)) <= 0 || len(strings.TrimSpace(ip)) <= 0 {
		return router.ErrCNameMissingArgs
	}
	if _, err := r.Addr(cname); err != nil {
		return err
	}
	r.ip = ip

	if err := r.createOrNuke(CREATE); err != nil {
		return err
	}
	return nil
}

//unset cname for a fullname : test.megambox.com
func (r route53Router) UnsetCName(cname string, ip string) error {
	r.cname = cname
	if len(strings.TrimSpace(r.cname)) <= 0 {
		return router.ErrCNameMissingArgs
	}

	if _, err := r.Addr(cname); err != nil {
		return err
	}

	if err := r.createOrNuke(DELETE); err != nil {
		return err
	}
	return nil
}

func (r route53Router) Addr(cname string) (string, error) {
	r.cname = cname
	chp, err := r.chopIt()
	if err != nil {
		return "", err
	}
	if err := r.zoneMatch(chp); err != nil {
		return "", err
	}

	rr, err := r.zone.ResourceRecordSets(r.client)
	if err != nil {
		return "", err
	}

	for i := range rr.ResourceRecordSets {
		rrp := strings.TrimSpace(rr.ResourceRecordSets[i].Name)
		if strings.HasSuffix(rrp, ".") {
			rrp = strings.TrimRight(rrp, ".")
		}
		if strings.Compare(rrp, chp) == 0 {
			return rr.ResourceRecordSets[i].Name, nil
		}
	}
	return "", router.ErrCNameNotFound
}

func (r *route53Router) createOrNuke(action string) error {
	log.Debugf("%s (cname, ip)", action)

	var u = route53.ChangeResourceRecordSetsRequest{
		ZoneID:  r.zone.HostedZoneId(),
		Comment: "",
		Changes: []route53.Change{
			{
				Action: action,
				Name:   r.cname,
				Type:   "A",
				TTL:    300,
				Value:  r.ip,
			},
		},
	}

	if _, err := u.Create(r.client); err != nil {
		return err
	}
	return nil
}

func (r *route53Router) StartupMessage() (string, error) {
	return "R53 router ok!", nil
}

func (r *route53Router) chopIt() (string, error) {
	chp, err := router.ChopDomain(r.cname)
	if err != nil {
		return "", router.ErrInvalidCName
	}
	return chp, err
}

//get the hosted zones and match the domain name with it.
func (r *route53Router) zoneMatch(chop string) error {
	zones := r.client.Zones().HostedZones
	noMatch := true
	for i := range zones {
		p := zones[i].Name
		if strings.HasSuffix(zones[i].Name, ".") {
			p = strings.TrimRight(p, ".")
		}
		if strings.Compare(p, chop) == 0 {
			r.zone = zones[i]
			noMatch = false
			break
		}
	}
	if noMatch {
		return router.ErrDomainNotFound
	}
	return nil
}
