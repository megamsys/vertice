package auth

import (
	"errors"
	"strings"
)

var ErrInvalidToken = errors.New("Invalid token")

type Token interface {
	GetValue() string

	User() (*User, error)

	GetUserName() string
}

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
