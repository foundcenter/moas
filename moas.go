package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"io"
)

func main() {
	router := httprouter.New()
	router.GET("/", hello)

	log.Fatal(http.ListenAndServe(":8081", router))

}

func hello (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	io.WriteString(w, "Hello world!")
}
