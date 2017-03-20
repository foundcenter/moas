package auth

import (
	"encoding/json"
	"net/http"

	"github.com/alioygur/gores"
	"github.com/foundcenter/moas/backend/controllers/response"
	"github.com/foundcenter/moas/backend/middleware/jwt_auth"
	"github.com/foundcenter/moas/backend/middleware/logger"
	"github.com/foundcenter/moas/backend/services/auth"
	"github.com/foundcenter/moas/backend/services/slack"
	"github.com/foundcenter/moas/backend/services/gmail"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Load(router *httprouter.Router) {
	stdChain := alice.New(logger.Handler)
	router.Handler("GET", "/auth", stdChain.ThenFunc(handleAuth))
	router.Handler("POST", "/auth/login", stdChain.ThenFunc(handleAuth))
	router.Handler("POST", "/auth/google", stdChain.ThenFunc(handleGoogleAuth))
	router.Handler("POST", "/auth/slack", stdChain.ThenFunc(handleSlackAuth))
	router.Handler("POST", "/auth/gmail", stdChain.ThenFunc(handleGmailAuth))
	router.Handler("GET", "/auth/check", stdChain.ThenFunc(handleAuthMock))
	router.Handler("GET", "/jwt", stdChain.Append(jwt_auth.Handler).ThenFunc(handleJwtMock))
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

func handleAuth(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var l loginRequest
	err := decoder.Decode(&l)
	if err != nil {
		panic(err)
	}

	err, user := auth.Login(l.Email, l.Password)

	if err != nil {
		// later use switch
		// if there are more reasons
		// err.Error() == auth.BadCredentials
		response.Reply(w).Unauthorized(err)
		return
	}

	response.Reply(w).Ok(user)
}

func handleAuthMock(w http.ResponseWriter, r *http.Request) {

	t := r.Header.Get("Authorization")
	token := t[7:len(t)]

	err, parsed := auth.ParseToken(token)

	if err != nil {
		response.Reply(w).Unauthorized(err)
		return
	}

	gores.JSON(w, 200, map[string]interface{}{"user_sub": parsed})
}

func handleJwtMock(w http.ResponseWriter, r *http.Request) {

	t := r.Header.Get("Authorization")
	token := t[7:len(t)]

	err, parsed := auth.ParseToken(token)

	if err != nil {
		response.Reply(w).Unauthorized(err)
		return
	}

	gores.JSON(w, 200, map[string]interface{}{"user_sub": parsed})
}
