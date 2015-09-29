package route53

import (
	"fmt"
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
	if _, err := r.Addr(cname); err != nil {
		return err
	}

	r.cname = cname
	r.ip = ip

	if err := r.createOrNuke(CREATE); err != nil {
		return err
	}
	return nil
}

//unset cname for a fullname : test.megambox.com
func (r route53Router) UnsetCName(cname string, ip string) error {
	if _, err := r.Addr(cname); err != nil {
		return err
	}

	r.cname = cname

	if err := r.createOrNuke(DELETE); err != nil {
		return err
	}
	return nil
}

func (r route53Router) Addr(cname string) (string, error) {
	zones := r.client.Zones().HostedZones

	for i := range zones {
		if strings.Compare(zones[i].Name, cname) == 0 {
			r.zone = zones[i]
			break
		}
	}

	rr, err := r.zone.ResourceRecordSets(r.client)
	if err != nil {
		return "", err
	}

	for i := range rr.ResourceRecordSets {
		if strings.Compare(strings.TrimSpace(rr.ResourceRecordSets[i].Name), cname) == 0 {
			return rr.ResourceRecordSets[i].Name, nil
		}
	}
	return "", router.ErrCNameExists
}

func (r *route53Router) createOrNuke(action string) error {
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
	return fmt.Sprintf("route53 router => %s", r.cname), nil
}
