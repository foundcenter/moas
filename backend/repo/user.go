package repo

import (
	"fmt"
	"github.com/foundcenter/moas/backend/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Session *mgo.Session
}

func (u *User) Insert(user models.User) (error, models.User) {

	c := u.Session.DB("local").C("users")

	err := c.Insert(user)

	if err != nil {
		if mgo.IsDup(err) {
			fmt.Println("User is duplicate")
			return err, models.User{}
		}
		return err, models.User{}
	}

	return nil, user
}

func (u *User) FindByEmailPassword(email, password string) (error, models.User) {
	c := u.Session.DB("local").C("users")

	model := models.User{}
	err := c.Find(bson.M{"email": email, "password": password}).One(&model)

	if err != nil {
		return err, models.User{}
	}

	return nil, model
}

func (u *User) FindById(id string) (error, models.User) {
	c := u.Session.DB("local").C("users")

	model := models.User{}
	err := c.Find(bson.M{"sub": id}).One(&model)

	if err != nil {
		return err, models.User{}
	}

	return nil, model
}

func (u *User) FindByIdOrInsert(user models.User) (models.User, string) {
	_, storedUser := u.FindById(user.Sub)
	if storedUser.Sub == "" {
		u.Insert(user)
		return user, "register"
	}
	return storedUser, "login"
}
