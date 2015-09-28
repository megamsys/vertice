package api

import "github.com/megamsys/megamd/auth"

type Token struct {
	Token     string
	UserEmail string
}

func (t *Token) GetValue() string {
	return t.Token
}

func (t *Token) User() (*auth.User, error) {
	return auth.GetUserByEmail(t.UserEmail)
}

func (t *Token) GetUserName() string {
	return t.UserEmail
}

func Auth(t string) (auth.Token, error) {
	tt := Token{
		Token:     "aaaa",
		UserEmail: "info@megam.io",
	}
	return &tt, nil
}
