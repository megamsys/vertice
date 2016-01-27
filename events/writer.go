package events

import (
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/subd/eventsd"
)

var EW *EventsWriter
var eventStorageAgeLimit = "default=24h"
var eventStorageEventLimit = "default=100000" //Max number of events to store (per type).

type EventsWriter struct {
	H *events
}

func NewWriter(e *eventsd.Config) error {
	ew := &EventsWriter{
		H: NewEventManager(parseEventsStoragePolicy()),
	}
	return ew.open(e)
}

func (ew *EventsWriter) open(e *eventsd.Config) error {
	//loop through all the event_handlers and start them.
	EW = ew
	return nil
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
