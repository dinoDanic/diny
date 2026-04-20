package auth

import "errors"

type Session struct {
	UserID int
	Token  string
}

func Login(userID int, password string) (*Session, error) {
	if password == "" {
		return nil, errors.New("password required")
	}
	return &Session{UserID: userID, Token: "dummy-token"}, nil
}
