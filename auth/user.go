
package auth

type User struct {
	Email    string
	Password string
	APIKey   string
}

func GetUserByEmail(email string) (*User, error) {
	/*
		var u User
		conn, err := db.Conn()
		if err != nil {
			return nil, err
		}

		return &u, nil
	*/
	return nil, nil
}
