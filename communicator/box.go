package communicator

import (
	"time"
	"strings"
)

// Boxlog represents a log entry.
type Boxlog struct {
	Timestamp string
	Message   string
	Source    string
	Name      string
	Unit      string
}

type Box struct {
	Name        string
	AccountId   string
	HostId      string
}

// GetName returns the assemblyname.domain(assembly001YeahBoy.megambox.com) of the box.
func (b *Box) GetFullName() string {
	return strings.TrimSpace(b.Name)
}

// Log adds a log message to the app. Specifying a good source is good so the
// user can filter where the message come from.
func (box *Box) Log(message, source, unit string) error {
	messages := strings.Split(message, "\n")
	logs := make([]interface{}, 0, len(messages))
	for _, msg := range messages {
		if len(strings.TrimSpace(msg)) > 0 {
			bl := Boxlog{
				Timestamp: time.Now().Local().Format(time.RFC822),
				Message:   msg,
				Source:    source,
				Name:      box.Name,
				Unit:      box.AccountId,
			}
			logs = append(logs, bl)
		}
	}
	if len(logs) > 0 {
		_ = notify(box.GetFullName(), logs)
	}
	return nil
}
