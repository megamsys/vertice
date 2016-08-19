package provision

import (
	"encoding/json"
"fmt"
	log "github.com/Sirupsen/logrus"
	nsqc "github.com/crackcomm/nsqueue/consumer"
	nsqp "github.com/crackcomm/nsqueue/producer"
	"github.com/megamsys/vertice/meta"
)

const (
	maxInFlight = 300
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
	fmt.Println("***************LogListener*************")
	b := make(chan Boxlog, maxInFlight)
	fmt.Println(b)
	cons := nsqc.New()
	fmt.Printf("%#v", cons)
	go func() {
		defer close(b)
		if err := cons.Register(logQueue(a.Name), "clients", maxInFlight, dumpLog(b)); err != nil {
			fmt.Println("**************register*************")
			fmt.Println(a.Name)
			return
		}

		if err := cons.Connect(meta.MC.NSQd...); err != nil {
			fmt.Println("************connect")
			return
		}
		log.Debugf("%s: start", logQueue(a.Name))
		fmt.Println("^^^^^^^^^^^^^^^^^start^^^^^^^^^^")
		cons.Start(true)
		log.Debugf("%s: start OK", logQueue(a.Name))
	}()

	l := LogListener{B: b, c: cons}
	fmt.Println(&l)
	return &l, nil
}

func (l *LogListener) Close() (err error) {
	l.c.Stop()
	return nil
}

func dumpLog(b chan Boxlog) func(m *nsqc.Message) {
	fmt.Println("*****************dumpLog****************")

	return func(msg *nsqc.Message) {
		go func(b chan Boxlog, msg *nsqc.Message) {
			fmt.Println(msg.Body)
			bl := Boxlog{}
			if err := json.Unmarshal(msg.Body, &bl); err != nil {
				log.Errorf("Unparsable log message, ignoring: %s", string(msg.Body))
			} else {
				fmt.Println("+++++++++++++++++++++++++++")
				fmt.Println(bl)
				b <- bl
			}
		}(b, msg)
	}
}

func notify(boxName string, messages []interface{}) error {
	fmt.Println("*************************notify****************8")
	fmt.Println(boxName)
	fmt.Println(messages)
	pons := nsqp.New()

	if err := pons.Connect(meta.MC.NSQd[0]); err != nil {
		return err
	}

	defer pons.Stop()

	for _, msg := range messages {
		log.Debugf("%s:%s", logQueue(boxName), msg)
		if err := pons.PublishJSONAsync(logQueue(boxName), msg, nil); err != nil {
			fmt.Println("**********publish****************")
			fmt.Println(logQueue(boxName))
			fmt.Println(msg)
			log.Errorf("Error on publish: %s", err.Error())
		}
	}
	return nil
}
