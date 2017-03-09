package auth

import (
	"io/ioutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
	"encoding/json"
	"log"
	"fmt"
	"context"
	"github.com/foundcenter/moas/backend/models"
	"errors"
)

type Credentials struct {
	Cid string `json:"cid"`
	Csecret string `json:"csecret"`
}

var cred Credentials
var conf *oauth2.Config

func init() {
	file, err := ioutil.ReadFile("./config/google-credentials.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &cred)

	conf = &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  "http://localhost:4200",
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

	fmt.Printf("Access token is %s \n", accessToken)

	client := conf.Client(context.TODO(), accessToken)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Printf("Error with auth handler %s \n", err)
		return errors.New("Could not get user info"), models.User{}
	}
	//defer userInfo.Body.Close()
	//data, _ := ioutil.ReadAll(userInfo.Body)

	decoder := json.NewDecoder(userInfo.Body)
	var gu models.User
	err = decoder.Decode(&gu)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("User info is: %s \n", string(data))

	return nil, gu
}