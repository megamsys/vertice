package auth

import (
//	"github.com/megamsys/megamd/db"
)

type APIToken struct {
	Token     string `json:"api_key"`
	UserEmail string `json:"email"`
}

func (t *APIToken) GetValue() string {
	return t.Token
}

func (t *APIToken) User() (*User, error) {
	return GetUserByEmail(t.UserEmail)
}

func (t *APIToken) IsAppToken() bool {
	return false
}

func (t *APIToken) GetUserName() string {
	return t.UserEmail
}

func (t *APIToken) GetAppName() string {
	return ""
}

func getAPIToken(header string) (*APIToken, error) {
	/*

		token, err := ParseToken(header)
		if err != nil {
			return nil, err
		}
		var t APIToken
		t, err := db.Fetch(id)
		if err != nil {
			return nil, err
		}
		return &t, ni
	*/
	return &APIToken{Token: "aaaa",
		UserEmail: "info@megam.io"}, nil
}

func APIAuth(token string) (*APIToken, error) {
	return getAPIToken(token)
}
