package repo

import (
	"github.com/foundcenter/moas/backend/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Database *mgo.Database
}

func (u *User) Insert(user models.User) (models.User, error) {

	c := u.Database.C("users")

	err := c.Insert(user)

	return user, err
}

func (u *User) FindByEmailPassword(email, password string) (models.User, error) {
	c := u.Database.C("users")

	model := models.User{}
	err := c.Find(bson.M{"email": email, "password": password}).One(&model)

	return model, err
}

func (u *User) FindById(id string) (models.User, error) {
	c := u.Database.C("users")

	model := models.User{}
	err := c.FindId(bson.ObjectIdHex(id)).One(&model)

	return model, err
}

func (u *User) FindByEmail(email string) (models.User, error) {
	c := u.Database.C("users")

	model := models.User{}
	err := c.Find(bson.M{"emails": bson.M{"$in": [1]string{email}}}).One(&model)

	return model, err
}

func (u *User) FindByIdOrInsert(user models.User) (models.User, string, error) {
	storedUser, err := u.FindById(user.ID.String())

	if err != nil {
		return storedUser, "", err
	}

	if !storedUser.ID.Valid() {
		u.Insert(user)
		return user, "register", nil
	}
	return storedUser, "login", nil
}
