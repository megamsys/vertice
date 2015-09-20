package auth

import (
	"errors"
	"strings"
)

type Token interface {
	GetValue() string
	GetAppName() string
	GetUserName() string
	IsAppToken() bool
	User() (*User, error)
}

var ErrInvalidToken = errors.New("Invalid token")

// ParseToken extracts token from a header:
// 'type token' or 'token'
func ParseToken(header string) (string, error) {
	s := strings.Split(header, " ")
	var value string
	if len(s) < 3 {
		value = s[len(s)-1]
	}
	if value != "" {
		return value, nil
	}
	return value, ErrInvalidToken
}
