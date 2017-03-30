package auth

import (
	"errors"
	"time"

	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"

	"github.com/dgrijalva/jwt-go"
)

const (
	BadCredentials   = "Bad Credentials"
	hmacSampleSecret = "818DC95A5C27370654E087E0CFEFC13C876F7B3D0B5BF9ACE7F3FBE385D16EF9"
)

type MyClaims struct {
	User_ID string `json:"user_id"`
	jwt.StandardClaims
}

type GoogleAuth struct {
	Code        string `json:"code"`
	RedirectURL string `json:"redirectUri"`
}

type SlackAuth struct {
	Code        string `json:"code"`
	RedirectURL string `json:"redirectUri"`
}

type JiraAuth struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(email string, password string) (error, models.User) {
	db := repo.New()
	defer db.Destroy()

	user, err := db.UserRepo.FindByEmailPassword(email, password)

	if err != nil {
		return errors.New(BadCredentials), models.User{}
	}

	return nil, user
}

func IssueToken(user models.User) (error, string) {

	mc := MyClaims{user.ID.Hex(), jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), Issuer: "moas"}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mc)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))

	if err != nil {
		return err, ""
	}
	return nil, tokenString
}

func ParseToken(tokenString string) (string, error) {

	myClaims := MyClaims{}
	token, _ := jwt.ParseWithClaims(tokenString, &myClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(hmacSampleSecret), nil
	})

	if !token.Valid {
		return "", errors.New("Token not valid!")
	} else {
		return myClaims.User_ID, nil
	}
}
