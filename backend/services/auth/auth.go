package auth

import (
	"errors"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	BadCredentials = "Bad Credentials"
	hmacSampleSecret = "818DC95A5C27370654E087E0CFEFC13C876F7B3D0B5BF9ACE7F3FBE385D16EF9"
)
type MyClaims struct {
	User_sub string `json:"user_sub"`
	jwt.StandardClaims
}

type GoogleAuth struct {
	Code string `json:"code"`
}

func Login(email string, password string) (error, models.User) {
	db := repo.New()
	defer db.Destroy()

	err, user := db.UserRepo.FindByEmailPassword(email, password)

	if err != nil {
		return errors.New(BadCredentials), models.User{}
	}

	return nil, user
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

func ParseToken(tokenString string) (error, string) {

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