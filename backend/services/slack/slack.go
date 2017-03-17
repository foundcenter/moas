package slack

import (
	"context"
	"errors"
	"fmt"

	"github.com/foundcenter/moas/backend/config"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
	slackAuth "golang.org/x/oauth2/slack"
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
		Endpoint: slackAuth.Endpoint,
	}
}

func Connect(ctx context.Context, code string) (error, models.User) {
	accessToken, err := conf.Exchange(ctx, code)
	if err != nil {
		fmt.Printf("Error with auth handler %s with code %s \n", err, code)
		return errors.New("Could not exchange token"), models.User{}
	}

	client := slack.New(accessToken.AccessToken)

	res, err := client.GetUserIdentity()

	if err != nil {
		return err, models.User{}
	}

	db := repo.New()
	defer db.Destroy()

	err, user := db.UserRepo.FindByEmail(res.User.Email)

	if err != nil {
		return err, user
	}

	// If user is already registered merge data
	if !user.ID.Valid() {
		user.Name = res.User.Name
		user.Email = res.User.Email
		user.Picture = res.User.Image512
	}

	user.Emails = append(user.Emails, user.Email)
	user.Accounts["slack"] = accessToken

	return nil, user
}

func Search(ctx context.Context, userID string, query string) ([]models.ResultResponse, error) {
	db := repo.New()
	defer db.Destroy()

	user, err := db.UserRepo.FindById(userID)
	if err != nil {
		return nil, err
	}

	token := user.Accounts["slack"]

	if token == nil {
		return nil, errors.New("Slack account not defined")
	}

	messages, files, err := slack.Search(query, slack.SearchParameters{})

	var results []models.ResultResponse

	for m := range messages.Matches {
		results = append(results, models.ResultResponse{
			AccountID:   "slack",
			Service:     "slack",
			Resource:    "message",
			Title:       m.Username,
			Description: m.Text,
			Url:         m.Permalink,
		})
	}

	for f := range files.Matches {
		results = append(results, models.ResultResponse{
			AccountID:   "slack",
			Service:     "slack",
			Resource:    "file",
			Title:       f.Title,
			Description: "",
			Url:         f.Permalink,
		})
	}

	return results, nil
}
