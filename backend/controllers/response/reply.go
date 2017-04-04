package response

import (
	"github.com/alioygur/gores"
	"log"
	"net/http"
)

type Response struct {
	writer http.ResponseWriter
}

type ResponseMeta struct {
	Type string `json:"type"`
}

func Reply(w http.ResponseWriter) *Response {
	return &Response{writer: w}
}

func (r *Response) BadRequest() {
	err := &Error{"Bad request", http.StatusBadRequest}
	gores.JSON(r.writer, err.Status, map[string]interface{}{"error": err})
}

func (r *Response) ServerInternalError(err error) {
	log.Print(err.Error())
	errorResponse := &Error{"Internal server error: " + err.Error(), http.StatusInternalServerError}
	gores.JSON(r.writer, errorResponse.Status, map[string]interface{}{"error": errorResponse})
}

func (r *Response) Unauthorized(e error) {
	err := &Error{e.Error(), http.StatusUnauthorized}
	gores.JSON(r.writer, err.Status, map[string]interface{}{"error": err})
}

func (r *Response) Ok(data interface{}, meta ...interface{}) {
	if len(meta) == 0 {
		gores.JSON(r.writer, http.StatusOK, map[string]interface{}{"data": data})
	} else {
		gores.JSON(r.writer, http.StatusOK, map[string]interface{}{"data": data, "meta": meta[0]})
	}
}

func (r *Response) Created(data interface{}) {
	gores.JSON(r.writer, http.StatusCreated, map[string]interface{}{"data": data})
}
