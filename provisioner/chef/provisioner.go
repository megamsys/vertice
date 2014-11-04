package chef

import (
	log "code.google.com/p/log4go"
	"github.com/megamsys/megamd/iaas"
	"github.com/megamsys/megamd/provisioner"
)

func Init() {
	provisioner.Register("chef", &Chef{})
}

type Chef struct {
}

func (i *Chef) CreateCommand(assembly *provisioner.AssemblyResult, id string) (string, error) {
	// Iaas Provider
	iaas, pdc, err1 := iaas.GetIaasProvider(assembly.Components[0].Requirements.Host)
	if err1 != nil {
		log.Error("Error: Iaas Provider :\n%s.", err1)
		return "", err1
	}
	log.Info(iaas)
	str, iaaserr := iaas.CreateMachine(pdc, assembly)
	if iaaserr != nil {
		log.Error("Error: Iaas Provider doesn't create machine:\n%s.", iaaserr)
		return "", iaaserr
	}
	log.Info(str)
	return str, nil
}
