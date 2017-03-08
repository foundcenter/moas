package auth

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"github.com/foundcenter/moas/backend/middleware/logger"
	"encoding/json"
	"github.com/foundcenter/moas/backend/services/auth"
	"github.com/foundcenter/moas/backend/controllers/response"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Load(router *httprouter.Router) {
	router.Handler("GET", "/auth", alice.New(logger.Handler).ThenFunc(handleAuth))
	router.Handler("POST", "/auth/login", alice.New(logger.Handler).ThenFunc(handleAuthMock))
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{"data": map[string]interface{}{"age": 1, "johnson": "peanuts"}}
	json.NewEncoder(w).Encode(res)
}

func handleAuthMock(w http.ResponseWriter, r *http.Request) {
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
