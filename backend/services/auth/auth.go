package auth

import (
	"errors"
	"github.com/foundcenter/moas/backend/models"
)

func Login(email string, password string) (error, models.User) {
	if email != "moas@foundcenter.com" || password != "moas123" {
		return errors.New(BadCredentials), models.User{}
	}

	return nil, models.User{Email: "moas@foundcenter.com", Id: "1"}
}

const (
	BadCredentials = "Bad Credentials"
)