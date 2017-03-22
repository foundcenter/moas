package search

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/foundcenter/moas/backend/controllers/response"
	"github.com/foundcenter/moas/backend/middleware/jwt_auth"
	"github.com/foundcenter/moas/backend/middleware/logger"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/foundcenter/moas/backend/services/auth"
	"github.com/foundcenter/moas/backend/services/drive"
	"github.com/foundcenter/moas/backend/services/github"
	"github.com/foundcenter/moas/backend/services/gmail"
	"github.com/foundcenter/moas/backend/services/slack"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"sync"
)

func Load(router *httprouter.Router) {
	stdChain := alice.New(logger.Handler)
	router.Handler("GET", "/search", stdChain.Append(jwt_auth.Handler).ThenFunc(handleSearch))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {

	resultOfSearch := make([]models.SearchResult, 0)
	var wg sync.WaitGroup
	query := r.URL.Query().Get("q")

	queueOfResults := make(chan []models.SearchResult, 2)

	token := context.Get(r, "user").(*jwt.Token).Raw
	user_id, _ := auth.ParseToken(token)

	db := repo.New()
	defer db.Destroy()
	user, err := db.UserRepo.FindById(user_id)
	if err != nil {
		response.Reply(w).Unauthorized(err)
		return
	}

	wg.Add(len(user.Accounts))
	for _, provider := range user.Accounts {
		switch provider.Type {
		case "gmail":
			// gmail search
			go func() {
				result := gmail.Search(r.Context(), provider, query)
				queueOfResults <- result
			}()
		case "drive":
			// drive search
			go func() {
				result := drive.Search(r.Context(), provider, query)
				queueOfResults <- result
			}()
		case "slack":
			// slack search
			go func() {
				result, _ := slack.Search(r.Context(), provider, query)
				queueOfResults <- result
			}()
		case "github":
			// github search
			go func() {
				result, _ := github.Search(r.Context(), provider, query)
				queueOfResults <- result
			}()
			// and other providers...
		}
	}
	//here we wait result of search from all services
	go func() {
		for r := range queueOfResults {
			resultOfSearch = append(resultOfSearch, r...)
			wg.Done()
		}
	}()

	wg.Wait()

	response.Reply(w).Ok(resultOfSearch)

}
