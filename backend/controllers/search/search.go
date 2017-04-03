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
	"github.com/foundcenter/moas/backend/services/jira"
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

	var wg sync.WaitGroup
	queueOfResults := make(chan models.SearchResultAndErrors, 2)
	resultOfSearch := make([]models.SearchResult, 0)
	searchErrors := models.SearchErrors{}
	query := r.URL.Query().Get("q")

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
	for _, account := range user.Accounts {
		if account.Active {
			switch account.Type {
			case "gmail":
				// gmail search
				go func(account models.AccountInfo) {
					searchError := models.SearchError{}
					result, err := gmail.Search(r.Context(), account, query)
					if err != nil {
						searchError.AccountID = account.ID
						searchError.AccountType = account.Type
						searchError.Error = err
					}
					queueOfResults <- models.SearchResultAndErrors{result, searchError}
				}(account)
			case "drive":
				// drive search
				go func(account models.AccountInfo) {
					e := models.SearchError{}
					result, err := drive.Search(r.Context(), account, query)
					if err != nil {
						e.AccountID = account.ID
						e.AccountType = account.Type
						e.Error = err
					}
					queueOfResults <- models.SearchResultAndErrors{result, e}
				}(account)
			case "slack":
				// slack search
				go func(account models.AccountInfo) {
					e := models.SearchError{}
					result, err := slack.Search(r.Context(), account, query)
					if err != nil {
						e.AccountID = account.ID
						e.AccountType = account.Type
						e.Error = err
					}
					queueOfResults <- models.SearchResultAndErrors{result, e}
				}(account)
			case "github":
				// github search
				go func(account models.AccountInfo) {
					e := models.SearchError{}
					result, err := github.Search(r.Context(), account, query)
					if err != nil {
						e.AccountID = account.ID
						e.AccountType = account.Type
						e.Error = err
					}
					queueOfResults <- models.SearchResultAndErrors{result, e}
				}(account)
			case "jira":
				//jira search
				go func(account models.AccountInfo) {
					e := models.SearchError{}
					result, err := jira.Search(r.Context(), account, query)
					if err != nil {
						e.AccountID = account.ID
						e.AccountType = account.Type
						e.Error = err
					}
					queueOfResults <- models.SearchResultAndErrors{result, e}
				}(account)

			}
		} else {
			wg.Done()
		}
	}
	//here we wait result of search from all services
	go func() {
		for r := range queueOfResults {
			resultOfSearch = append(resultOfSearch, r.SearchResult...)
			if r.SearchError.AccountID != "" {
				searchErrors.SearchError = append(searchErrors.SearchError, r.SearchError)
			}
			wg.Done()
		}
	}()

	wg.Wait()

	if len(searchErrors.SearchError) > 0 {
		response.Reply(w).Ok(resultOfSearch, searchErrors)
	} else {
		response.Reply(w).Ok(resultOfSearch)
	}

}
