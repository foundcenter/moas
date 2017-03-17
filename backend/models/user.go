package models

import (
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	//Id            bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	ID            bson.ObjectId            `json:"id" bson:"_id,omitempty"`
	Sub           string                   `json:"sub"`
	Name          string                   `json:"name"`
	GivenName     string                   `json:"given_name"`
	FamilyName    string                   `json:"family_name"`
	Profile       string                   `json:"profile"`
	Picture       string                   `json:"picture"`
	Email         string                   `json:"email"`
	Emails        []string                 `json:"emails"`
	EmailVerified bool                     `json:"email_verified"`
	Gender        string                   `json:"gender"`
	Locale        string                   `json:"locale"`
	Accounts      map[string]*oauth2.Token `json:"accounts"`
}
