package auth

import (
	"context"
	"encoding/json"
	"github.com/foundcenter/moas/backend/controllers/response"
	"github.com/foundcenter/moas/backend/middleware/jwt_auth"
	"github.com/foundcenter/moas/backend/middleware/logger"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/foundcenter/moas/backend/services/auth"
	"github.com/foundcenter/moas/backend/services/drive"
	"github.com/foundcenter/moas/backend/services/github"
	"github.com/foundcenter/moas/backend/services/gmail"
	"github.com/foundcenter/moas/backend/services/jira"
	"github.com/foundcenter/moas/backend/services/slack"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"errors"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// WrapHandler for route params
func WrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		ctx := context.WithValue(r.Context(), "params", ps)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Load routes for router
func Load(router *httprouter.Router) {

	standardChain := alice.New(logger.Handler)
	extendedChain := standardChain.Append(jwt_auth.Handler)

	router.Handler("GET", "/auth/check", extendedChain.ThenFunc(handleAuthCheck))
	router.Handler("POST", "/auth/google", standardChain.ThenFunc(handleGoogleAuth))
	//router.Handler("POST", "/auth/slack", standardChain.ThenFunc(handleSlackAuth))
	//router.Handler("POST", "/auth/gmail", standardChain.ThenFunc(handleGmailAuth))
	//router.Handler("POST", "/auth/drive", standardChain.ThenFunc(handleDriveAuth))
	router.Handler("POST", "/connect/slack", extendedChain.ThenFunc(handleSlackConnect))
	router.Handler("POST", "/connect/gmail", extendedChain.ThenFunc(handleGmailConnect))
	router.Handler("POST", "/connect/drive", extendedChain.ThenFunc(handleDriveConnect))
	router.Handler("POST", "/connect/github", extendedChain.ThenFunc(handleGithubConnect))
	router.Handler("POST", "/connect/jira", extendedChain.ThenFunc(handleJiraConnect))
	router.Handle("PUT", "/connect/github/:username", WrapHandler(extendedChain.ThenFunc(handleGithubConnectWithApiToken)))

}

func handleGithubConnectWithApiToken(w http.ResponseWriter, r *http.Request) {

	params := r.Context().Value("params").(httprouter.Params)
	username := params.ByName("username")

	gitHubToken := auth.GitHubToken{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&gitHubToken)
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	if len(gitHubToken.Token) != 40 {
		response.Reply(w).Error(errors.New("Personal token must be 40 characters long"), http.StatusUnprocessableEntity)
		return
	}

	token := r.Header.Get("Authorization")
	userID, err := auth.ParseToken(token[7:])
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	user, err := github.ConnectWithApiToken(r.Context(), userID, gitHubToken.Token, username)
	if err != nil {
		response.Reply(w).Error(err, http.StatusUnprocessableEntity)
		return
	}

	response.Reply(w).Ok(map[string]interface{}{"user": user})
}

func handleAuthCheck(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	userID, err := auth.ParseToken(token[7:])
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	db := repo.New()
	defer db.Destroy()
	user, err := db.UserRepo.FindById(userID)
	if err != nil {
		response.Reply(w).Unauthorized(err)
		return
	}

	response.Reply(w).Ok(map[string]interface{}{"user": user})

}

func handleGoogleAuth(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	user, err := gmail.Login(r.Context(), ga.Code, ga.RedirectURL)
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	err, tokenString := auth.IssueToken(user)
	if err != nil {
		response.Reply(w).BadRequest()
		return
	}

	response.Reply(w).Ok(map[string]interface{}{"user": user, "token": tokenString})

}

func handleGmailAuth(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)

	user, err := gmail.Login(r.Context(), ga.Code, ga.RedirectURL)
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	err, tokenString := auth.IssueToken(user)
	if err != nil {
		response.Reply(w).BadRequest()
		return
	}

	response.Reply(w).Ok(map[string]interface{}{"user": user, "token": tokenString})

}

func handleGmailConnect(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)

	token := r.Header.Get("Authorization")
	userID, err := auth.ParseToken(token[7:])
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	user, err := gmail.Connect(r.Context(), userID, ga.Code, ga.RedirectURL)
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	response.Reply(w).Ok(map[string]interface{}{"user": user})
}

func handleDriveAuth(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)

	user, err := drive.Login(r.Context(), ga.Code, ga.RedirectURL)
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	err, tokenString := auth.IssueToken(user)
	if err != nil {
		response.Reply(w).BadRequest()
		return
	}

	response.Reply(w).Ok(map[string]interface{}{"user": user, "token": tokenString})
}

func handleDriveConnect(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)

	token := r.Header.Get("Authorization")
	userID, err := auth.ParseToken(token[7:])
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	user, err := drive.Connect(r.Context(), userID, ga.Code, ga.RedirectURL)
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	response.Reply(w).Ok(map[string]interface{}{"user": user})
}

func handleSlackAuth(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.SlackAuth
	err := decoder.Decode(&ga)

	user, err := slack.Login(r.Context(), ga.Code, ga.RedirectURL)
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	err, tokenString := auth.IssueToken(user)
	if err != nil {
		response.Reply(w).BadRequest()
		return
	}

	response.Reply(w).Ok(map[string]interface{}{"user": user, "token": tokenString})

}

func handleSlackConnect(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)

	token := r.Header.Get("Authorization")
	userID, err := auth.ParseToken(token[7:])
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	user, err := slack.Connect(r.Context(), userID, ga.Code, ga.RedirectURL)
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	response.Reply(w).Ok(map[string]interface{}{"user": user})

}

func handleGithubConnect(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)

	token := r.Header.Get("Authorization")
	userID, err := auth.ParseToken(token[7:])
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	user, err := github.Connect(r.Context(), userID, ga.Code, ga.RedirectURL)
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	response.Reply(w).Ok(map[string]interface{}{"user": user})

}

func handleJiraConnect(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var jiraAuth auth.JiraAuth
	err := decoder.Decode(&jiraAuth)

	token := r.Header.Get("Authorization")
	userID, err := auth.ParseToken(token[7:])
	if err != nil {
		response.Reply(w).ServerInternalError(err)
		return
	}

	user, err := jira.Connect(r.Context(), userID, jiraAuth.Url, jiraAuth.Username, jiraAuth.Password)

	if err != nil {
		response.Reply(w).BadRequest()
		return
	}

	response.Reply(w).Ok(map[string]interface{}{"user": user})

}
