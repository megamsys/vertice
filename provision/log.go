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
	B <-chan Boxlog
	c *nsqc.Consumer
}

func logQueue(boxName string) string {
	return boxName + LogPubSubQueueSuffix
}

func NewLogListener(a *Box) (*LogListener, error) {
	b := make(chan Boxlog, 10)
	cons := nsqc.New()
	go func() {
		defer close(b)
		if err := cons.Register(logQueue(a.Name), "clients", maxInFlight, dumpLog(b)); err != nil {
			return
		}

		if err := cons.Connect(meta.MC.NSQd...); err != nil {
			return
		}
		log.Debugf("%s: start", logQueue(a.GetFullName()))
		cons.Start(true)
		log.Debugf("%s: start OK", logQueue(a.GetFullName()))
	}()

	l := LogListener{B: b, c: cons}
	return &l, nil
}

func (l *LogListener) Close() (err error) {
	l.c.Stop()
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
	pons := nsqp.New()

	if err := pons.Connect(meta.MC.NSQd[0]); err != nil {
		return err
	}

	defer pons.Stop()

	for _, msg := range messages {
		log.Debugf("%s:%s", logQueue(boxName), msg)
		if err := pons.PublishJSONAsync(logQueue(boxName), msg, nil); err != nil {
			log.Errorf("Error on publish: %s", err.Error())
		}
	}
	return nil
}
