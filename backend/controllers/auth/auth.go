package auth

import (
	"encoding/json"
	"errors"
	"github.com/alioygur/gores"
	"github.com/dgrijalva/jwt-go"
	"github.com/foundcenter/moas/backend/controllers/response"
	"github.com/foundcenter/moas/backend/middleware/jwt_auth"
	"github.com/foundcenter/moas/backend/middleware/logger"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/foundcenter/moas/backend/services/auth"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"time"
)

const (
	hmacSampleSecret = "818DC95A5C27370654E087E0CFEFC13C876F7B3D0B5BF9ACE7F3FBE385D16EF9"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MyClaims struct {
	User_sub string `json:"user_sub"`
	jwt.StandardClaims
}

type GoogleAuth struct {
	Code string `json:"code"`
}

func Load(router *httprouter.Router) {
	router.Handler("GET", "/auth", alice.New(logger.Handler).ThenFunc(handleAuth))
	router.Handler("POST", "/auth/login", alice.New(logger.Handler).ThenFunc(handleAuth))
	router.Handler("POST", "/auth/google", alice.New(logger.Handler ).ThenFunc(handleGoogleAuth))
	router.Handler("GET", "/auth/check", alice.New(logger.Handler).ThenFunc(handleAuthMock))
	router.Handler("GET", "/jwt", alice.New(logger.Handler, jwt_auth.Handler).ThenFunc(handleJwtMock))
}

func handleGoogleAuth(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var ga GoogleAuth
	err := decoder.Decode(&ga)
	if err != nil {
		panic(err)
	}

	err, user := auth.Exchange(ga.Code)
	if err != nil {
		response.Reply(w).Unauthorized(err)
		return
	}

	err, tokenString := IssueToken(user)
	if err != nil {
		response.Reply(w).BadRequest()
		return
	}

	//find or insert in DB
	db := repo.New()
	defer db.Destroy()
	user, action := db.UserRepo.FindByIdOrInsert(user)
	if action == "login" {
		response.Reply(w).Logged(map[string]interface{}{"user": user, "token": tokenString})
		return
	}
	// register
	response.Reply(w).Created(map[string]interface{}{"user": user, "token": tokenString})

}

func IssueToken(user models.User) (error, string) {

	mc := MyClaims{user.Sub, jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), Issuer: "moas"}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mc)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))

	if err != nil {
		return err, ""
	}
	return nil, tokenString
}

func ParseToken(tokenString string) (error, interface{}) {

	myClaims := MyClaims{}
	token, _ := jwt.ParseWithClaims(tokenString, &myClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(hmacSampleSecret), nil
	})

	if !token.Valid {
		return errors.New("Token not valid!"), ""
	} else {
		return nil, myClaims.User_sub
	}
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

	err, parsed := ParseToken(token)

	if err != nil {
		response.Reply(w).Unauthorized(err)
		return
	}

	gores.JSON(w, 200, map[string]interface{}{"user_sub": parsed})
}

func handleJwtMock(w http.ResponseWriter, r *http.Request) {

	t := r.Header.Get("Authorization")
	token := t[7:len(t)]

	err, parsed := ParseToken(token)

	if err != nil {
		response.Reply(w).Unauthorized(err)
		return
	}

	gores.JSON(w, 200, map[string]interface{}{"user_sub": parsed})
}