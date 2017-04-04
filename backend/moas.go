package main

import (
	"log"
	"net/http"

	"fmt"

	"github.com/foundcenter/moas/backend/config"
	"github.com/foundcenter/moas/backend/controllers"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	router := httprouter.New()
	controllers.Load(router)
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler(router)

	fmt.Printf("Starting server on: %s\n", config.Settings.App.URL)

	log.Fatal(http.ListenAndServe(config.Settings.App.URL, handler))
}
