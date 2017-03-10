package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/foundcenter/moas/backend/config"
	"github.com/foundcenter/moas/backend/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var conf *oauth2.Config

func init() {
	conf = &oauth2.Config{
		ClientID:     config.Settings.Google.ClientID,
		ClientSecret: config.Settings.Google.ClientSecret,
		RedirectURL:  config.Settings.Google.RedirectURL,
		Scopes: []string{
			"profile",
			"email",
			"https://www.googleapis.com/auth/gmail.readonly",
			"https://www.googleapis.com/auth/drive.readonly",
		},
		Endpoint: google.Endpoint,
	}
}

func Exchange(code string) (error, models.User) {

	accessToken, err := conf.Exchange(context.TODO(), code)
	if err != nil {
		fmt.Printf("Error with auth handler %s with code %s \n", err, code)
		return errors.New("Could not exchange token"), models.User{}
	}

	client := conf.Client(context.TODO(), accessToken)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Printf("Error with auth handler %s \n", err)
		return errors.New("Could not get user info"), models.User{}
	}
	defer userInfo.Body.Close()

	decoder := json.NewDecoder(userInfo.Body)
	var gu models.User
	err = decoder.Decode(&gu)
	if err != nil {
		panic(err)
	}

	return nil, gu
}
