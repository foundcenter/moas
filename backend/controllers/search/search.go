package search

import (
	"github.com/foundcenter/moas/backend/controllers/response"
	"github.com/foundcenter/moas/backend/middleware/jwt_auth"
	"github.com/foundcenter/moas/backend/middleware/logger"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/services/gmail"
	"github.com/foundcenter/moas/backend/services/drive"
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

	queueOfResults := make(chan []models.ResultResponse,2)

	token := context.Get(r, "user").(*jwt.Token).Raw
	_, user_sub:= auth.ParseToken(token)

	wg.Add(2)

	// gmail search
	go func() {
		result := gmail.Search(user_sub, query)
		queueOfResults<-result
	}()

	// drive search
	go func() {
		result := drive.Search(user_sub, query)
		queueOfResults<-result
	}()

	// and other providers...

	//here we wait result of search from all services
	go func() {
		for r := range queueOfResults {
			resultOfSearch = append(resultOfSearch, r...)
			wg.Done()
		}
	}()

	wg.Wait()

	response.Reply(w).SearchResult(resultOfSearch)

}
