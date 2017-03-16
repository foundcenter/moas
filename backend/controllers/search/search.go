package search

import (
	"github.com/foundcenter/moas/backend/controllers/response"
	"github.com/foundcenter/moas/backend/middleware/jwt_auth"
	"github.com/foundcenter/moas/backend/middleware/logger"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/services/gmail"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/foundcenter/moas/backend/services/auth"
)

func Load(router *httprouter.Router) {
	stdChain := alice.New(logger.Handler)
	router.Handler("GET", "/search", stdChain.Append(jwt_auth.Handler).ThenFunc(handleSearch))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {

	resultOfSearch := make([]models.ResultResponse, 0)
	var wg sync.WaitGroup
	query:=r.URL.Query().Get("q")

	wg.Add(1)
	// gmail search
	go func() {
		token := context.Get(r, "user").(*jwt.Token).Raw
		_, user_sub:= auth.ParseToken(token)

		result := gmail.Search(user_sub, query)
		resultOfSearch = append(resultOfSearch, result...)
		wg.Done()
	}()

	// drive search
	//go func() {
	//	result := drive.Search(user_sub, query)
	//	resultOfSearch = append(resultOfSearch, result...)
	//	wg.Done()
	//}()

	// and other providers...

	wg.Wait()

	response.Reply(w).SearchResult(resultOfSearch)

}
