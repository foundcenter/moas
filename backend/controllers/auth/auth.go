package auth

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"github.com/foundcenter/moas/backend/middleware/logger"
	"encoding/json"
	"github.com/foundcenter/moas/backend/services/auth"
	"github.com/foundcenter/moas/backend/controllers/response"
	"github.com/foundcenter/moas/backend/repo"
	"fmt"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Load(router *httprouter.Router) {
	router.Handler("GET", "/auth", alice.New(logger.Handler).ThenFunc(handleAuth))
	router.Handler("POST", "/auth/login", alice.New(logger.Handler).ThenFunc(handleAuth))
	router.Handler("POST", "/auth/google", alice.New(logger.Handler).ThenFunc(handleGoogleAuth))
	//router.Handler("POST", "/auth/login/mock", alice.New(logger.Handler).ThenFunc(handleAuthMock))
	//router.Handler("GET", "/user/test", alice.New(logger.Handler).ThenFunc(insertDummyUser))
}

type GoogleAuth struct {
	Code string `json:"code"`
}
func handleGoogleAuth(w http.ResponseWriter, r *http.Request)  {

	decoder := json.NewDecoder(r.Body)
	var ga GoogleAuth
	err := decoder.Decode(&ga)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Code in request is %s \n", ga.Code)

	// error, user, registe/login
	// error, data (user . register/login)

	err, user := auth.Exchange(ga.Code)
	if err != nil {
		//didnt exchange user
		response.Reply(w).Unauthorized()
		return
	}

	//find or insert in DB
	db := repo.New()
	defer db.Destroy()
	user, action := db.UserRepo.FindByIdOrInsert(user)
	if action == "login" {
		response.Reply(w).Ok(ga.Code)
		return
	}
	// register
	response.Reply(w).Created(ga.Code)

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
		response.Reply(w).Unauthorized()
		return
	}

	//issue jwt

	response.Reply(w).Ok(user)
}

//func handleAuthMock(w http.ResponseWriter, r *http.Request) {
//	decoder := json.NewDecoder(r.Body)
//	var l loginRequest
//	err := decoder.Decode(&l)
//	if err != nil {
//		panic(err)
//	}
//
//	err, user := auth.LoginMock(l.Email, l.Password)
//
//	if err != nil {
//		// later use switch
//		// if there are more reasons
//		// err.Error() == auth.BadCredentials
//		response.Reply(w).Unauthorized()
//		return
//	}
//
//	//issue jwt
//
//	response.Reply(w).Ok(user)
//}

//func insertDummyUser(w http.ResponseWriter, r *http.Request)  {
//	db := repo.New()
//	defer db.Destroy()
//
//	err, user := db.UserRepo.Insert()
//
//	if err != nil {
//		response.Reply(w).BadRequest()
//		return
//	}
//
//	//issue jwt
//
//	response.Reply(w).Ok(user)
//}