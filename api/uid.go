package api

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/satori/go.uuid"
)

func Uid(prefix string) string {
	return prefix + strings.Replace(uuid.NewV1().String(), "-", "", -1)
}

func PP(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
