package machine

import (
  "github.com/megamsys/libgo/utils"
 // log "github.com/Sirupsen/logrus"
 // constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/marketplaces/provision"
	"github.com/megamsys/vertice/marketplaces/provision/one/cluster"
)

type OneProvisioner interface {
	Cluster() *cluster.Cluster
}

type Machine struct {
	Name         string
	Region       string
	Id           string
	CartonId     string
	CartonsId    string
	AccountId    string
	Level        provision.BoxLevel
	SSH          provision.BoxSSH
	Image        string
	VCPUThrottle string
	VMId         string
	VNCHost      string
	VNCPort      string
	ImageId      string
	StorageType  string
	Routable     bool
	Status       utils.Status
	State        utils.State
}

type CreateArgs struct {
	Commands    []string
	Box         *provision.Box
	Compute     provision.BoxCompute
	Deploy      bool
	Provisioner OneProvisioner
}

// func (m *Machine) Create(args *CreateArgs) error {
// 	asm, err := carton.NewAssembly(m.CartonId, m.AccountId, "")
// 	if err != nil {
// 		return err
// 	}
//
//  if err = asm.NukeAndSetOutputs(id); err != nil {
// 		return err
// 	}
// 	return nil
// }
//
//
//
// func (m *Machine) IPs(nics []virtualmachine.Nic) map[string][]string {
// 	var ips = make(map[string][]string)
// 	pubipv4s := []string{}
// 	priipv4s := []string{}
// 	for _, nic := range nics {
// 		if nic.IPaddress != "" {
// 			ip4 := strings.Split(nic.IPaddress, ".")
// 			if len(ip4) == 4 {
// 				if ip4[0] == "192" || ip4[0] == "10" || ip4[0] == "172" {
// 					priipv4s = append(priipv4s, nic.IPaddress)
// 				} else {
// 					pubipv4s = append(pubipv4s, nic.IPaddress)
// 				}
// 			}
// 		}
// 	}
//
// 	ips[carton.PUBLICIPV4] = pubipv4s
// 	ips[carton.PRIVATEIPV4] = priipv4s
// 	return ips
// }
//
// func (m *Machine) mergeSameIPtype(mm map[string][]string) map[string][]string {
// 	for IPtype, ips := range mm {
// 		var sameIp string
// 		for _, ip := range ips {
// 			sameIp = sameIp + ip + ", "
// 		}
// 		if sameIp != "" {
// 			mm[IPtype] = []string{strings.TrimRight(sameIp, ", ")}
// 		}
// 	}
// 	return mm
// }
//
// func (m *Machine) Remove(p OneProvisioner, state constants.State) error {
// 	log.Debugf("  removing machine in one (%s)", m.Name)
// 	return nil
// }
//
//
// //it possible to have a Notifier interface that does this, duck typed b y Assembly, Components.
// func (m *Machine) SetStatus(status utils.Status) error {
// 	log.Debugf("  set status[%s] of machine (%s, %s)", m.Id, m.Name, status.String())
// 	return nil
// }
//
// func (m *Machine) SetMileStone(state utils.State) error {
// 	log.Debugf("  set state[%s] of machine (%s, %s)", m.Id, m.Name, state.String())
//
// 	return nil
// }
