package provision

import (
	//	"encoding/json"
	log "github.com/Sirupsen/logrus"
	//	"github.com/megamsys/megamd/meta"
)

var LogPubSubQueueSuffix = "_log"

type LogListener struct {
	B <-chan Boxlog
	//q *nsq.Consumer
}

func logQueue(boxName string) string {
	return boxName + LogPubSubQueueSuffix
}

func NewLogListener(a *Box) (*LogListener, error) {
	/*if c := seg.NewConsumer(TOPIC, "scheduler"); c != nil {
		c.Set("nsqd", s.Meta.NSQd)
		c.Set("nsqlookupd", s.Meta.NSQLookupd)
		c.Set("concurrency", 15)
		c.Set("max_attempts", 10)
		c.Set("max_in_flight", 150)
		c.Set("default_requeue_delay", "15s")
		s.Consumer = c
		c.Start(nsq.HandlerFunc(s.processNSQ))
	}

	b := make(chan Boxlog, 10)
	go func() {
		defer close(b)
		for msg := range subChan {
			bl := Boxlog{}
			err := json.Unmarshal(msg, &bl)
			if err != nil {
				log.Errorf("Unparsable log message, ignoring: %s", string(msg))
				continue
			}
			//write box logs to the  "channel - b" as json.
			b <- bl
		}
	}()

	l := LogListener{B: b, q: pubSubQ}
	return &l, nil

	pubSubQ, err := amqp.NewRabbitMQ(meta.MC.AMQP, logQueue(a.Name))
	if err != nil {
		return nil, err
	}
	subChan, err := pubSubQ.Sub()
	if err != nil {
		return nil, err
	}

	b := make(chan Boxlog, 10)
	go func() {
		defer close(b)
		for msg := range subChan {
			bl := Boxlog{}
			err := json.Unmarshal(msg, &bl)
			if err != nil {
				log.Errorf("Unparsable log message, ignoring: %s", string(msg))
				continue
			}
			//write box logs to the  "channel - c" as json.
			b <- bl
		}
	}()
	l := LogListener{B: b, q: pubSubQ}
	return &l, nil
	*/
	return nil, nil
}

func (l *LogListener) Close() (err error) {
	//err = l.q.UnSub()
	return nil
}

func notify(boxName string, messages []interface{}) error {
	log.Debugf("  notify %s", logQueue(boxName))
	/*pubSubQ, err := amqp.NewRabbitMQ("testing", logQueue(boxName))
	if err != nil {
		return err
	}
	defer pubSubQ.Close()
	err = pubSubQ.Connect()
	if err != nil {
		return err
	}

	for _, msg := range messages {
		bytes, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Error on logs notify: %s", err.Error())
			continue
		}
		err = pubSubQ.Pub(bytes)
		if err != nil {
			log.Errorf("Error on logs notify: %s", err.Error())
		}
	}
	*/
	return nil
}
