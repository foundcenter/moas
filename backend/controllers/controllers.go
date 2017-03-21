package controllers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/foundcenter/moas/backend/controllers/auth"
	"github.com/foundcenter/moas/backend/controllers/search"

)

func Load(router *httprouter.Router) {
	auth.Load(router)
	search.Load(router)
}
