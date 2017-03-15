package search

import (
	"github.com/foundcenter/moas/backend/middleware/jwt_auth"
	"github.com/foundcenter/moas/backend/middleware/logger"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"github.com/foundcenter/moas/backend/services/gmail"
)

func Load(router *httprouter.Router) {
	stdChain := alice.New(logger.Handler)
	router.Handler("GET", "/search", stdChain.Append(jwt_auth.Handler).ThenFunc(handleSearch))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	gmail.HandleGmailSearch(w,r)
}
