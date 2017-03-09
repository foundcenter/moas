package repo

import "gopkg.in/mgo.v2"

var masterSession *mgo.Session

func InitMasterSession(session *mgo.Session)  {
	masterSession = session
}

func GetSession() *mgo.Session {
	return masterSession.Copy()
}

type DB struct {
	Session *mgo.Session
	UserRepo *User
}

func New() *DB {
	db :=  &DB{Session:masterSession.Copy()}
	db.UserRepo = &User{Session:db.Session}
	return db
}

func (db *DB) Destroy() {
	db.Session.Close()
}