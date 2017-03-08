package repo

import (

	"github.com/foundcenter/moas/backend/models"
	"gopkg.in/mgo.v2"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Session *mgo.Session
}

func (u *User) Insert() (error, models.User) {

	c := u.Session.DB("local").C("users")

	model := models.User{"55552", "neb.vojvodic@gmail.com", "kikiriki123"}
	err := c.Insert(model)

	if err != nil {
		if mgo.IsDup(err) {
			fmt.Println("User is duplicate")
			return err, models.User{}
		}
		return err, models.User{}
	}

	return nil, model
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

