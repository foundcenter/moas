package controllers

import (
	"github.com/foundcenter/moas/backend/controllers/auth"
	"github.com/foundcenter/moas/backend/controllers/search"
	"github.com/julienschmidt/httprouter"

	"github.com/foundcenter/moas/backend/controllers/account"
)

func Load(router *httprouter.Router) {
	auth.Load(router)
	search.Load(router)
	account.Load(router)
}
