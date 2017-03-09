package controllers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/foundcenter/moas/backend/controllers/auth"
)

func Load(router *httprouter.Router) {
	auth.Load(router)
}
