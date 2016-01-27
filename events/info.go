package events

import (
	"time"
)

const (
	// 10ms, i.e. 0.01s
	timePrecision time.Duration = 10 * time.Millisecond
)

// Event contains information general to events such as the time at which they
// occurred, their specific type, and the actual event. Event types are
// differentiated by the EventType field of Event.
type Event struct {
	// the absolute container name for which the event occurred
	ContainerName string `json:"container_name"`

	// the time at which the event occurred
	Timestamp time.Time `json:"timestamp"`

	// the type of event. EventType is an enumerated type
	EventType EventType `json:"event_type"`

	// the original event object and all of its extraneous data, ex. an
	// OomInstance
	EventData EventData `json:"event_data,omitempty"`
}

// EventType is an enumerated type which lists the categories under which
// events may fall. The Event field EventType is populated by this enum.
type EventType string

const (
	EventVMCreation        EventType = "vmCreation"
	EventVMDestroy                   = "vmDestroy"
	EventContainerCreation           = "containerCreation"
	EventContainerDestroy            = "containerDestroy"
	EventBill                        = "bill"
	EventUserAlert                   = "userAlert"
	EventDev                         = "dev"
	//EventWHMS                        = "whmsCreation"
)

// Extra information about an event. Only one type will be set.
type EventData struct {
	// Information about an OOM kill event.
	OomKill *OomKillEventData `json:"oom,omitempty"`
}

// Information related to an OOM kill instance
type OomKillEventData struct {
	// process id of the killed process
	Pid int `json:"pid"`

	// The name of the killed process
	ProcessName string `json:"process_name"`
}
