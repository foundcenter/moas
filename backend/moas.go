package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"github.com/foundcenter/moas/backend/controllers"
)

func main() {
	router := httprouter.New()
	controllers.Load(router)
	log.Fatal(http.ListenAndServe(":8081", router))
}