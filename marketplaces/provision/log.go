package provision

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	nsqc "github.com/crackcomm/nsqueue/consumer"
	nsqp "github.com/crackcomm/nsqueue/producer"
	"github.com/megamsys/vertice/meta"
)

const (
	maxInFlight = 300
)

var LogPubSubQueueSuffix = "_log"


// Boxlevel represents the deployment level.
type BoxLevel int

// Boxlog represents a log entry.
type Boxlog struct {
	Timestamp string
	Message   string
	Source    string
	Name      string
	Unit      string
}

type Box struct {
	Name string
}

type LogListener struct {
	B <-chan Boxlog
	c *nsqc.Consumer
}

type BoxSSH struct {
	User   string
	Prefix string
}

func (bs *BoxSSH) Pub() string {
	return bs.Prefix + "_pub"
}


func logQueue(boxName string) string {
	return boxName + LogPubSubQueueSuffix
}

func NewLogListener(a *Box) (*LogListener, error) {
	b := make(chan Boxlog, maxInFlight)
	cons := nsqc.New()
	go func() {
		defer close(b)
		if err := cons.Register(logQueue(a.Name), "clients", maxInFlight, dumpLog(b)); err != nil {
			return
		}

		if err := cons.Connect(meta.MC.NSQd...); err != nil {
			return
		}
		log.Debugf("%s: start", logQueue(a.Name))
		cons.Start(true)
		log.Debugf("%s: start OK", logQueue(a.Name))
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
