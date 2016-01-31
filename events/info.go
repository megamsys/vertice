package events

import (
	"time"
)

const (
	EventMachine   EventType = "machine"
	EventContainer           = "container"
	EventBill                = "bill"
	EventUser                = "user"

	Add EventAction = iota
	Destroy
	Status
	Deduct
	Notify

	// 10ms, i.e. 0.01s
	timePrecision time.Duration = 10 * time.Millisecond
)

// Event contains information general to events such as the time at which they
// occurred, their specific type, and the actual event. Event types are
// differentiated by the EventType field of Event.
type Event struct {
	// the time at which the event occurred
	Timestamp time.Time

	// the type of event. EventType is an enumerated type
	EventType EventType

	//the action can be
	//bill create, bill delete
	EventAction EventAction

	// the original event object and all of its extraneous data, ex. an
	// OomInstance
	EventData EventData
}

// EventType is an enumerated type which lists the categories under which
// events may fall. The Event field EventType is populated by this enum.
type EventType string

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
	case Notify:
		return "alert"
	default:
		return "arrgh"
	}
}

// Extra information about an event.
type EventData struct {
	M map[string]string
}

type EventChannel struct {
	// Watch ID. Can be used by the caller to request cancellation of watch events.
	watchId int
	// Channel on which the caller can receive watch events.
	channel chan *Event
}
