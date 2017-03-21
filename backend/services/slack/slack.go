package slack

import (
	"context"

	"fmt"
	"github.com/foundcenter/moas/backend/config"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/foundcenter/moas/backend/utils"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
	slackAuth "golang.org/x/oauth2/slack"
)

const AccountType = "slack"

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

func Login(ctx context.Context, code string) (models.User, error) {

	var user models.User
	accessToken, err := conf.Exchange(ctx, code)
	if err != nil {
		return user, err
	}

	client := slack.New(accessToken.AccessToken)

	res, err := client.GetUserIdentity()

	if err != nil {
		return user, err
	}

	db := repo.New()
	defer db.Destroy()

	user, err = db.UserRepo.FindByAccount(AccountType, res.User.ID)

	if err != nil {
		return user, err
	}

	// If user is already registered merge data
	if !user.ID.Valid() {
		user.Name = res.User.Name
		user.Picture = res.User.Image512
	}

	addAccount(ctx, &user, res, accessToken)

	user, err = db.UserRepo.Upsert(user)

	return user, err
}

func Connect(ctx context.Context, userID string, code string) (models.User, error) {
	var user models.User
	accessToken, err := conf.Exchange(ctx, code)

	if err != nil {
		return user, err
	}

	client := slack.New(accessToken.AccessToken)
	res, err := client.GetUserIdentity()

	db := repo.New()
	defer db.Destroy()

	user, err = db.UserRepo.FindById(userID)

	if err != nil {
		return user, err
	}

	addAccount(ctx, &user, res, accessToken)
	db.UserRepo.Update(user)

	return user, nil
}

func addAccount(ctx context.Context, user *models.User, res *slack.UserIdentityResponse, token *oauth2.Token) {
	a := models.AccountInfo{
		Type:  AccountType,
		ID:    res.User.ID,
		Data:  res,
		Token: token,
	}

	for _, acc := range user.Accounts {
		if acc.ID == a.ID && acc.Type == a.Type {
			return
		}
	}

	user.Accounts = append(user.Accounts, a)

	if res.User.Email != "" && !utils.Contains(user.Emails, res.User.Email) {
		user.Emails = append(user.Emails, res.User.Email)
	}
}

func Search(ctx context.Context, accountInfo models.AccountInfo, query string) ([]models.SearchResult, error) {
	db := repo.New()
	defer db.Destroy()

	if accountInfo.Type != AccountType {
		return nil, fmt.Errorf("AccountInfo type %s not valid. Should be %s.", accountInfo.Type, AccountType)
	}

	client := slack.New(accountInfo.Token.AccessToken)

	messages, files, _ := client.Search(query, slack.SearchParameters{})

	var results []models.SearchResult

	for _, m := range messages.Matches {
		results = append(results, models.SearchResult{
			AccountID:   "slack",
			Service:     "slack",
			Resource:    "message",
			Title:       m.Username,
			Description: m.Text,
			Url:         m.Permalink,
		})
	}

	for _, f := range files.Matches {
		results = append(results, models.SearchResult{
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
