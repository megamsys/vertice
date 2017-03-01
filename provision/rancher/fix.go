package rancher

/*
import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/vertice/provision/docker/container"
)

func (p *rancherProvisioner) fixContainers() error {
	/*containers, err := p.listAllContainers()
	if err != nil {
		return err
	}
	err = runInContainers(containers, func(c *container.Container, _ chan *container.Container) error {
		return p.checkContainer(c)
	}, nil, true)
	if err != nil {
		log.Errorf("error checking containers for fixing: %s", err.Error())
	}
	return err
*/
//	return nil
//}
/*
func (p *rancherProvisioner) checkContainer(container *container.Container) error {
	if container.Available() {
		info, err := container.NetworkInfo(p)
		if err != nil {
			return err
		}
		if info.HTTPHostPort != container.HostPort || info.IP != container.PublicIp {
			err = p.fixContainer(container, info)
			if err != nil {
				log.Errorf("error on fix container hostport for (container %s)", container.Id)
				return err
			}
		}
	}
	return nil
}

func (p *rancherProvisioner) fixContainer(container *container.Container, info container.NetworkInfo) error {
	/*if info.HTTPHostPort == "" {
		return nil
	}
	appInstance, err := app.GetByName(container.AppName)
	if err != nil {
		return err
	}
	r, err := getRouterForApp(appInstance)
	if err != nil {
		return err
	}
	err = r.RemoveRoute(container.AppName, container.Address())
	if err != nil && err != router.ErrRouteNotFound {
		return err
	}
	container.IP = info.IP
	container.HostPort = info.HTTPHostPort
	err = r.AddRoute(container.AppName, container.Address())
	if err != nil && err != router.ErrRouteExists {
		return err
	}*/
//return nil
//}
