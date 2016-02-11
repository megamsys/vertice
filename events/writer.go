package events

import (
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/vertice/events/alerts"
)

var W *EventsWriter
var eventStorageAgeLimit = "default=24h"
var eventStorageEventLimit = "default=100000"

type EventsConfigMap map[string]map[string]string

func (ec EventsConfigMap) Get(key string) map[string]string {
	return ec[key]
}

type EventsWriter struct {
	H *events
}

type eventWatcher struct {
	eventType EventType
	Watcher
}

func NewWrap(c EventsConfigMap) error {
	e := &EventsWriter{
		H: NewEventManager(parseEventsStoragePolicy()),
	}
	W = e
	return e.open(c)
}

func (e *EventsWriter) open(c EventsConfigMap) error {
	watchers := watchHandlers(c)
	for _, w := range watchers {
		ec, err := e.WatchForEvents(NewRequest(&eventReqOpts{etype: w.eventType}))
		if err != nil {
			log.Infof("error type  %#v %s\n", w.eventType, err.Error())
			return err
		}
		if err := w.Watch(ec); err != nil {
			return nil
		}
	}
	return nil
}

// can be called by the api which will take events returned on the channel
func (ew *EventsWriter) Write(e *Event) error {
	return ew.H.AddEvent(e)
}

// can be called by the api which will take events returned on the channel
func (ew *EventsWriter) WatchForEvents(request *Request) (*EventChannel, error) {
	return ew.H.WatchEvents(request)
}

// can be called by the api which will return all events satisfying the request
func (ew *EventsWriter) GetPastEvents(request *Request) ([]*Event, error) {
	return ew.H.GetEvents(request)
}

func (ew *EventsWriter) CloseEventChannel(watch_id int) {
	ew.H.StopWatch(watch_id)
}

func (ew *EventsWriter) Close() {
	//close all channels.
}

func watchHandlers(c EventsConfigMap) []*eventWatcher {
	watchers := make([]*eventWatcher, 0)
	watchers = append(watchers, &eventWatcher{eventType: EventMachine, Watcher: &Machine{}})
	watchers = append(watchers, &eventWatcher{eventType: EventContainer, Watcher: &Container{}})
	b := NewBill(c.Get(BILLMGR))
	watchers = append(watchers, &eventWatcher{eventType: EventBill, Watcher: b})
	watchers = append(watchers, &eventWatcher{eventType: EventUser, Watcher: NewUser(c, AfterFuncsMap{alerts.ONBOARD: AfterFuncs{b.OnboardFunc}})})
	return watchers
}

// Parses the events StoragePolicy from the flags.
func parseEventsStoragePolicy() StoragePolicy {
	policy := DefaultStoragePolicy()

	// Parse max age.
	parts := strings.Split(eventStorageAgeLimit, ",")
	for _, part := range parts {
		items := strings.Split(part, "=")
		if len(items) != 2 {
			log.Warningf("Unknown event storage policy %q when parsing max age", part)
			continue
		}
		dur, err := time.ParseDuration(items[1])
		if err != nil {
			log.Warningf("Unable to parse event max age duration %q: %v", items[1], err)
			continue
		}
		if items[0] == "default" {
			policy.DefaultMaxAge = dur
			continue
		}
		policy.PerTypeMaxAge[EventType(items[0])] = dur
	}

	// Parse max number.
	parts = strings.Split(eventStorageEventLimit, ",")
	for _, part := range parts {
		items := strings.Split(part, "=")
		if len(items) != 2 {
			log.Warningf("Unknown event storage policy %q when parsing max event limit", part)
			continue
		}
		val, err := strconv.Atoi(items[1])
		if err != nil {
			log.Warningf("Unable to parse integer from %q: %v", items[1], err)
			continue
		}
		if items[0] == "default" {
			policy.DefaultMaxNumEvents = val
			continue
		}
		policy.PerTypeMaxNumEvents[EventType(items[0])] = val
	}

	return policy
}
