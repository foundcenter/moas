package repo

import "gopkg.in/mgo.v2"
import "github.com/foundcenter/moas/backend/config"
import "log"

var masterSession *mgo.Session

type DB struct {
	Session  *mgo.Session
	UserRepo *User
}

func init() {
	log.Printf("Connecting to mongo %+v...", config.Settings.Mongo)
	session, err := mgo.DialWithInfo(config.Settings.Mongo)
	if err != nil {
		log.Panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	masterSession = session
}

func New() *DB {
	db := &DB{Session: masterSession.Copy()}
	db.UserRepo = &User{Database: db.Session.DB(config.Settings.Mongo.Database)}
	return db
}

func (db *DB) Destroy() {
	db.Session.Close()
}
