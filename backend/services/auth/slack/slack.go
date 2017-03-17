package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/amirilovic/snoop/models"
	"github.com/foundcenter/moas/backend/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/slack"
)

var conf *oauth2.Config

func init() {
	conf = &oauth2.Config{
		ClientID:     config.Settings.Slack.ClientID,
		ClientSecret: config.Settings.Slack.ClientSecret,
		RedirectURL:  config.Settings.Slack.RedirectURL,
		Scopes: []string{
			"search:read",
			"identity.basic",
			"identity.email",
		},
		Endpoint: slack.Endpoint,
	}
}

func Exchange(ctx Context, code string) (error, models.User) {

	accessToken, err := conf.Exchange(ctx, code)
	if err != nil {
		fmt.Printf("Error with auth handler %s with code %s \n", err, code)
		return errors.New("Could not exchange token"), models.User{}
	}

	client := conf.Client(ctx, accessToken)

	userInfoResponse, err := client.Get("https://slack.com/api/users.identity")
	if err != nil {
		log.Printf("Error with auth handler %s \n", err)
		return errors.New("Could not get user info"), models.User{}
	}
	defer userInfoResponse.Body.Close()

	decoder := json.NewDecoder(userInfoResponse.Body)
	err = decoder.Decode(&userInfo)
	if err != nil {
		panic(err)
	}

	gu.Accounts = map[string]*oauth2.Token{"google": accessToken}

	return nil, gu
}
func GetConfig() *oauth2.Config {
	return conf
}
