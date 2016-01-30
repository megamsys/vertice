package events

import (
	"time"
)

const (
	EventMachine   EventType = "machine"
	EventContainer           = "container"
	EventBill                = "bill"
	EventUser                = "user"
	//EventWHMS                        = "whmsCreation"

	Add EventAction = iota
	Destroy
	Status
	Deduct
	Alert

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

	//the action can be
	//bill create, bill delete
	EventAction EventAction

	// the original event object and all of its extraneous data, ex. an
	// OomInstance
	EventData EventData `json:"event_data,omitempty"`
}

// EventType is an enumerated type which lists the categories under which
// events may fall. The Event field EventType is populated by this enum.
type EventType string

type EventChannel struct {
	// Watch ID. Can be used by the caller to request cancellation of watch events.
	watchId int
	// Channel on which the caller can receive watch events.
	channel chan *Event
}

type EventAction int

func (v *EventAction) String() string {
	switch *v {
	case Add:
		return "add"
	case Destroy:
		return "destroy"
	case Status:
		return "status"
	case Deduct:
		return "deduct"
	case Alert:
		return "alert"
	default:
		return "arrgh"
	}
}

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
