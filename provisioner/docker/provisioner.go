package docker

import (
	log "code.google.com/p/log4go"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/megamd/provisioner"
	"encoding/json"
)

func Init() {
	provisioner.Register("docker", &Docker{})
}

type Message struct {
	Id string `json:"id"`
}

type Docker struct {
}

func (i *Docker) CreateCommand(assembly *provisioner.AssemblyResult, id string) (string, error) {
    address := "Docker."+assembly.Name+"."+assembly.Components[0].Inputs.Domain
    com := &Message{Id: id}
	mapB, err := json.Marshal(com)  
	if err != nil {
        log.Error(err)
        return "", err
    }
    go publisher(address, string(mapB))
	return "", nil
}


func publisher(key string, json string) {
	factor, aerr := amqp.Factory()
	if aerr != nil {
		log.Error("Failed to get the queue instance: %s", aerr)
	}
	
	pubsub, perr := factor.Get(key)
	if perr != nil {
		log.Error("Failed to get the queue instance: %s", perr)
	}

	serr := pubsub.Pub([]byte(json))
	if serr != nil {
		log.Error("Failed to publish the queue instance: %s", serr)
	}
}