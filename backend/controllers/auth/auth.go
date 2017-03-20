package auth

import (
	"encoding/json"
	"net/http"
	"github.com/foundcenter/moas/backend/controllers/response"
	"github.com/foundcenter/moas/backend/middleware/logger"
	"github.com/foundcenter/moas/backend/services/auth"
	"github.com/foundcenter/moas/backend/services/slack"
	"github.com/foundcenter/moas/backend/services/gmail"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/foundcenter/moas/backend/services/drive"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Load(router *httprouter.Router) {
	stdChain := alice.New(logger.Handler)
	router.Handler("POST", "/auth/google", stdChain.ThenFunc(handleGoogleAuth))
	router.Handler("POST", "/auth/slack", stdChain.ThenFunc(handleSlackAuth))
	router.Handler("POST", "/auth/gmail", stdChain.ThenFunc(handleGmailAuth))
	router.Handler("POST", "/auth/drive", stdChain.ThenFunc(handleDriveAuth))
}

func handleGoogleAuth(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga auth.GoogleAuth
	err := decoder.Decode(&ga)
	if err != nil {
		panic(err)
	}

	//user, err := google.Exchange(ga.Code)
	//if err != nil {
	//	response.Reply(w).Unauthorized(err)
	//	return
	//}
	//
	////find or insert in DB
	//db := repo.New()
	//defer db.Destroy()
	//user, action, err := db.UserRepo.FindByIdOrInsert(user)
	//
	//if err != nil {
	//	response.Reply(w).ServerInternalError()
	//	return
	//}
	//
	//err, tokenString := auth.IssueToken(user)
	//if err != nil {
	//	response.Reply(w).BadRequest()
	//	return
	//}
	//
	//if action == "login" {
	//	response.Reply(w).Logged(map[string]interface{}{"user": user, "token": tokenString})
	//	return
	//}
	//// register
	//response.Reply(w).Created(map[string]interface{}{"user": user, "token": tokenString})

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

func handleDriveAuth(w http.ResponseWriter, r *http.Request) {

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