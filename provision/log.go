package provision

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	nsqc "github.com/crackcomm/nsqueue/consumer"
	nsqp "github.com/crackcomm/nsqueue/producer"
	"github.com/megamsys/megamd/meta"
)

const (
	maxInFlight = 150
)

var LogPubSubQueueSuffix = "_log"

type LogListener struct {
	b chan Boxlog
	B <-chan Boxlog
	c *nsqc.Consumer
}

func logQueue(boxName string) string {
	return boxName + LogPubSubQueueSuffix
}

func NewLogListener(a *Box) (*LogListener, error) {
	bo := make(chan Boxlog, 10)

	cons := nsqc.New()
	if err := cons.Register(logQueue(a.Name), "clients", maxInFlight, dumpLog(bo)); err != nil {
		return nil, err
	}

	if err := cons.Connect(meta.MC.NSQd...); err != nil {
		return nil, err
	}

	go cons.Start(true)

	l := LogListener{B: bo, c: cons, b: bo}
	return &l, nil
}

func (l *LogListener) Close() (err error) {
	l.c.Stop()
	defer close(l.b)
	return nil
}

func dumpLog(b chan Boxlog) func(m *nsqc.Message) {
	return func(msg *nsqc.Message) {
		go func(b chan Boxlog, msg *nsqc.Message) {
			bl := Boxlog{}
			if err := json.Unmarshal(msg.Body, &bl); err != nil {
				log.Errorf("Unparsable log message, ignoring: %s", string(msg.Body))
			} else {
				b <- bl
			}
		}(b, msg)
	}
}

func notify(boxName string, messages []interface{}) error {
	log.Debugf("  notify %s", logQueue(boxName))
	pons := nsqp.New()
	if err := pons.Connect(meta.MC.NSQd[0]); err != nil {
		return err
	}

	for _, msg := range messages {
		bytes, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Error on logs notify: %s", err.Error())
			continue
		}

		if err = pons.PublishAsync(logQueue(boxName), bytes, nil); err != nil {
			log.Errorf("Error on logs notify: %s", err.Error())
		}
	}
	defer pons.Stop()
	return nil
}
