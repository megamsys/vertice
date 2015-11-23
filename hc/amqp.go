package hc

import (
	"fmt"

	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/libgo/hc"
	"github.com/megamsys/megamd/meta"
)

func init() {
	hc.AddChecker("megamd:amqp", healthCheckAMQP)
}

func healthCheckAMQP() (interface{}, error) {
	pq, err := amqp.NewRabbitMQ(meta.MC.AMQP, "ping")
	if err != nil {
		return nil, err
	}
	if err = pq.Pub([]byte(`{"megamoja": "ping success"}`)); err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s up", meta.MC.AMQP), nil
}
