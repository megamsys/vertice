package auth

import (
//	"gopkg.in/check.v1"
)

/*
func (s *S) TestGetUserByEmail(c *check.C) {
	u := User{Email: "wolverine@xmen.com", Password: "123456"}
	err := u.Create()
	c.Assert(err, check.IsNil)
	defer u.Delete()
	u2, err := GetUserByEmail(u.Email)
	c.Assert(err, check.IsNil)
	c.Check(u2.Email, check.Equals, u.Email)
	c.Check(u2.Password, check.Equals, u.Password)
}

func (s *S) TestGetUserByEmailReturnsErrorWhenNoUserIsFound(c *check.C) {
	u, err := GetUserByEmail("unknown@globo.com")
	c.Assert(u, check.IsNil)
	c.Assert(err, check.Equals, ErrUserNotFound)
}

func (s *S) TestGetUserByEmailWithInvalidEmail(c *check.C) {
	u, err := GetUserByEmail("unknown")
	c.Assert(u, check.IsNil)
	c.Assert(err, check.NotNil)
	e, ok := err.(*errors.ValidationError)
	c.Assert(ok, check.Equals, true)
	c.Assert(e.Message, check.Equals, "invalid email")
}
*/
