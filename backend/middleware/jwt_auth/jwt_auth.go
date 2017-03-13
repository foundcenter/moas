package jwt_auth

import (
	"errors"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/foundcenter/moas/backend/controllers/response"
	"net/http"
)

const (
	hmacSampleSecret = "818DC95A5C27370654E087E0CFEFC13C876F7B3D0B5BF9ACE7F3FBE385D16EF9"
)

// Handler will log the HTTP requests.
func Handler(next http.Handler) http.Handler {

	jwtMiddleware := createJwtMiddleware()

	return jwtMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}))
}

func createJwtMiddleware() *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(hmacSampleSecret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			response.Reply(w).Unauthorized(errors.New(err))
		},
	})
}
