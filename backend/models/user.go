package models


import (
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"
)

const (
	LOGIN = "login"
	REGISTER = "register"
)

type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name"`
	Picture  string        `json:"picture"`
	Emails   []string      `json:"emails"`
	Accounts []AccountInfo `json:"accounts"`
}

type AccountInfo struct {
	Type   string        `json:"type"`
	ID     string        `json:"id"`
	Data   interface{}   `json:"data"`
	Token  *oauth2.Token `json:"-"`
	Active bool          `json:"active"`
}
