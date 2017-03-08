package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"github.com/foundcenter/moas/backend/controllers"
	"gopkg.in/mgo.v2"
	"github.com/foundcenter/moas/backend/repo"
)

func main() {
	initDatabase()

	router := httprouter.New()
	controllers.Load(router)
	log.Fatal(http.ListenAndServe(":8081", router))
}

func initDatabase() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	//defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	repo.InitMasterSession(session)
}