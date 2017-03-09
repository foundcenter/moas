package auth

import (
	"errors"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
)

const (
	BadCredentials = "Bad Credentials"
)

func Login(email string, password string) (error, models.User) {
	db := repo.New()
	defer db.Destroy()

	err, user := db.UserRepo.FindByEmailPassword(email, password)

	if err != nil {
		return errors.New(BadCredentials), models.User{}
	}

	return nil, user
}

//func LoginMock(email string, password string) (error, models.User) {
//	if email != "moas@foundcenter.com" || password != "moas123" {
//		return errors.New(BadCredentials), models.User{}
//	}
//
//	return nil, models.User{Email: "moas@foundcenter.com", Id: "1"}
//}

