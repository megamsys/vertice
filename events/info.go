package events

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/megamsys/vertice/carton/bind"
	"github.com/megamsys/vertice/events/alerts"
)

const (
	EventMachine   EventType = "machine"
	EventContainer           = "container"
	EventBill                = "bill"
	EventUser                = "user"
	// 10ms, i.e. 0.01s
	timePrecision time.Duration = 10 * time.Millisecond
)

type StoredEvent struct {
	Id         string         `json:"id"`
	AccountsId string         `json:"accounts" riak:"index"`
	Type       string         `json:"type"`
	Action     string         `json:"action"`
	Inputs     bind.JsonPairs `json:"inputs"`
	CreatedAt  string         `json:"created_at"`
}

func NewParseEvent(b []byte) (*StoredEvent, error) {
	st := &StoredEvent{}
	err := json.Unmarshal(b, &st)
	if err != nil {
		return nil, err
	}
	return st, err
}

func (st *StoredEvent) AsEvent() (*Event, error) {
	ea, err := strconv.Atoi(st.Action)
	if err != nil {
		return nil, err
	}

	e := Event{
		AccountsId:  st.AccountsId,
		EventType:   toEventType(st.Type),
		EventAction: toEventAction(ea),
		EventData:   EventData{M: st.Inputs.ToMap()},
		Timestamp:   time.Now().Local(),
	}

	if err != nil {
		return nil, err
	}
	return &e, err
}

// Event contains information general to events such as the time at which they
// occurred, their specific type, and the actual event. Event types are
// differentiated by the EventType field of Event.
type Event struct {
	AccountsId string
	// the time at which the event occurred
	Timestamp time.Time

	// the type of event. EventType is an enumerated type
	EventType EventType

	//the action can be
	//bill create, bill delete
	EventAction alerts.EventAction

	// the original event object and all of its extraneous data, ex. an
	// OomInstance
	EventData EventData
}

// EventType is an enumerated type which lists the categories under which
// events may fall. The Event field EventType is populated by this enum.
type EventType string

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

func toEventAction(a int) alerts.EventAction {
	return alerts.EventAction(a)
}

func toEventType(t string) EventType {
	switch t {
	case "machine":
		return EventMachine
	case "container":
		return EventContainer
	case "bill":
		return EventBill
	case "user":
		return EventUser
	default:
		return ""
	}
}
