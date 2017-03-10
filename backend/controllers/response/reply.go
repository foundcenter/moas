package response

import (
	"github.com/alioygur/gores"
	"net/http"
	"reflect"
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

func (r *Response) Unauthorized() {
	err := &Error{"Bad credentials", http.StatusUnauthorized}
	gores.JSON(r.writer, err.Status, map[string]interface{}{"error": err})
}

func (r *Response) Ok(data interface{}) {
	meta := &ResponseMeta{Type: reflect.TypeOf(data).Name()}
	gores.JSON(r.writer, http.StatusOK, map[string]interface{}{"data": data, "meta": meta})
}

func (r *Response) Logged(data interface{}) {
	//meta := &ResponseMeta{Type: reflect.TypeOf(data).Name()}
	gores.JSON(r.writer, http.StatusOK, map[string]interface{}{"data": data})
}

func (r *Response) Created(data interface{}) {
	//meta := &ResponseMeta{Type: reflect.TypeOf(data).Name()}
	gores.JSON(r.writer, http.StatusCreated, map[string]interface{}{"data": data})
}
