package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"github.com/foundcenter/moas/backend/controllers"
	"gopkg.in/mgo.v2"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/rs/cors"
)

func main() {
	initDatabase()

	router := httprouter.New()
	controllers.Load(router)
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8081", handler))
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