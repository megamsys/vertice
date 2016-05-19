package logbox

import (
	"encoding/json"
)

const (
	INFO             = "Info"
	ERROR            = "Error"
	WARN             = "Warning"
	VM_DEPLOY        = "VM_Deploying"
	CONTAINER_DEPLOY = "Container_Deploying"
)

type LogBox struct {
	Source  string `json:"Source"`
	Type    string `json:"Type"`
	Message string `json:"Message"`
}

func (a *LogBox) String() string {
	if d, err := json.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

func W(source, typ, message string) string {
	a := LogBox{
		Source:  source,
		Type:    typ,
		Message: message}
	return a.String()
}
