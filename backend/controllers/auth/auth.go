package auth

import (
	"encoding/json"
	"github.com/foundcenter/moas/backend/controllers/response"
	"github.com/foundcenter/moas/backend/middleware/jwt_auth"
	"github.com/foundcenter/moas/backend/middleware/logger"
	"github.com/foundcenter/moas/backend/services/auth"
	"github.com/foundcenter/moas/backend/services/drive"
	"github.com/foundcenter/moas/backend/services/gmail"
	"github.com/foundcenter/moas/backend/services/slack"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Load(router *httprouter.Router) {
	standardChain := alice.New(logger.Handler)
	extendedChain := standardChain.Append(jwt_auth.Handler)
	router.Handler("POST", "/auth/google", standardChain.ThenFunc(handleGoogleAuth))
	//router.Handler("POST", "/auth/slack", standardChain.ThenFunc(handleSlackAuth))
	//router.Handler("POST", "/auth/gmail", standardChain.ThenFunc(handleGmailAuth))
	//router.Handler("POST", "/auth/drive", standardChain.ThenFunc(handleDriveAuth))
	router.Handler("POST", "/connect/slack", extendedChain.ThenFunc(handleSlackConnect))
	router.Handler("POST", "/connect/gmail", extendedChain.ThenFunc(handleGmailConnect))
	router.Handler("POST", "/connect/drive", extendedChain.ThenFunc(handleDriveConnect))
}

func handleGoogleAuth(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)
	if err != nil {
		panic(err)
	}

	user, err := gmail.Login(r.Context(), ga.Code)
	if err != nil {
		response.Reply(w).ServerInternalError()
		return
	}

	err, tokenString := auth.IssueToken(user)
	if err != nil {
		response.Reply(w).BadRequest()
		return
	}

	response.Reply(w).Logged(map[string]interface{}{"user": user, "token": tokenString})

}

func handleGmailAuth(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)

	user, err := gmail.Login(r.Context(), ga.Code)
	if err != nil {
		response.Reply(w).ServerInternalError()
		return
	}

	err, tokenString := auth.IssueToken(user)
	if err != nil {
		response.Reply(w).BadRequest()
		return
	}

	response.Reply(w).Logged(map[string]interface{}{"user": user, "token": tokenString})

}

func handleGmailConnect(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)

	token := r.Header.Get("Authorization")
	user_id, err := auth.ParseToken(token[7:])
	if err != nil {
		response.Reply(w).ServerInternalError()
		return
	}

	user, err := gmail.Connect(r.Context(), user_id, ga.Code)
	if err != nil {
		response.Reply(w).ServerInternalError()
		return
	}

	response.Reply(w).Logged(map[string]interface{}{"user": user, "token": token})

}

func handleDriveAuth(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)

	user, err := drive.Login(r.Context(), ga.Code)
	if err != nil {
		response.Reply(w).ServerInternalError()
		return
	}

	err, tokenString := auth.IssueToken(user)
	if err != nil {
		response.Reply(w).BadRequest()
		return
	}

	response.Reply(w).Logged(map[string]interface{}{"user": user, "token": tokenString})

}

func handleDriveConnect(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)

	token := r.Header.Get("Authorization")
	user_id, err := auth.ParseToken(token[7:])
	if err != nil {
		response.Reply(w).ServerInternalError()
		return
	}

	user, err := drive.Connect(r.Context(), user_id, ga.Code)
	if err != nil {
		response.Reply(w).ServerInternalError()
		return
	}

	response.Reply(w).Logged(map[string]interface{}{"user": user, "token": token})

}

func handleSlackAuth(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.SlackAuth
	err := decoder.Decode(&ga)

	user, err := slack.Login(r.Context(), ga.Code)
	if err != nil {
		response.Reply(w).ServerInternalError()
		return
	}

	err, tokenString := auth.IssueToken(user)
	if err != nil {
		response.Reply(w).BadRequest()
		return
	}

	response.Reply(w).Logged(map[string]interface{}{"user": user, "token": tokenString})

}

func handleSlackConnect(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)

	token := r.Header.Get("Authorization")
	user_id, err := auth.ParseToken(token[7:])
	if err != nil {
		response.Reply(w).ServerInternalError()
		return
	}

	user, err := slack.Connect(r.Context(), user_id, ga.Code)
	if err != nil {
		response.Reply(w).ServerInternalError()
		return
	}

	response.Reply(w).Logged(map[string]interface{}{"user": user, "token": token})

}
