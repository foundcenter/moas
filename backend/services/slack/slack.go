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

func initOAuthConfig(redirectURL string) *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     config.Settings.Slack.ClientID,
		ClientSecret: config.Settings.Slack.ClientSecret,
		Scopes: []string{
			"search:read",
			"identify",
			"users:read",
		},
		Endpoint: slackAuth.Endpoint,
	}

	return config
}

func Login(ctx context.Context, code string, redirectURL string) (models.User, error) {

	var user models.User
	config := initOAuthConfig(redirectURL)
	config.RedirectURL = redirectURL
	accessToken, err := config.Exchange(ctx, code)
	if err != nil {
		return user, err
	}
	client := slack.New(accessToken.AccessToken)

	auth, err := client.AuthTest()

	if err != nil {
		return user, err
	}

	userId := auth.UserID

	res, err := client.GetUserInfo(userId)

	if err != nil {
		return user, err
	}

	db := repo.New()
	defer db.Destroy()

	user, err = db.UserRepo.FindByAccount(AccountType, res.ID)

	if err != nil {
		return user, err
	}

	// If user is already registered merge data
	if !user.ID.Valid() {
		user.Name = res.Name
		user.Picture = res.Profile.Image192
	}

	addAccount(ctx, &user, res, accessToken)

	user, err = db.UserRepo.Upsert(user)

	return user, err
}

func Connect(ctx context.Context, userID string, code string, redirectURL string) (models.User, error) {
	var user models.User
	config := initOAuthConfig(redirectURL)
	config.RedirectURL = redirectURL
	accessToken, err := config.Exchange(ctx, code)

	if err != nil {
		return user, err
	}

	client := slack.New(accessToken.AccessToken)

	auth, err := client.AuthTest()

	if err != nil {
		return user, err
	}

	userId := auth.UserID

	res, err := client.GetUserInfo(userId)

	db := repo.New()
	defer db.Destroy()

	user, err = db.UserRepo.FindById(userID)

	if err != nil {
		return user, err
	}

	addAccount(ctx, &user, res, accessToken)
	user, err = db.UserRepo.Update(user)

	return user, nil
}

func addAccount(ctx context.Context, user *models.User, slackUser *slack.User, token *oauth2.Token) {
	a := models.AccountInfo{
		Type:  AccountType,
		ID:    slackUser.ID,
		Data:  slackUser,
		Token: token,
		Active: true,
	}

	for _, acc := range user.Accounts {
		if acc.ID == a.ID && acc.Type == a.Type {
			return
		}
	}

	user.Accounts = append(user.Accounts, a)

	if slackUser.Profile.Email != "" && !utils.Contains(user.Emails, slackUser.Profile.Email) {
		user.Emails = append(user.Emails, slackUser.Profile.Email)
	}
}

func Search(ctx context.Context, accountInfo models.AccountInfo, query string) ([]models.SearchResult, error) {
	db := repo.New()
	defer db.Destroy()

	if accountInfo.Type != AccountType {
		return nil, fmt.Errorf("AccountInfo type %s not valid. Should be %s.", accountInfo.Type, AccountType)
	}

	client := slack.New(accountInfo.Token.AccessToken)

	messages, files, _ := client.Search(query, slack.NewSearchParameters())

	results := make([]models.SearchResult, 0)

	for _, m := range messages.Matches {
		msg := models.SearchResult{
			AccountID:   "slack",
			Service:     "slack",
			Resource:    "message",
			Title:       m.Username,
			Description: m.Text,
			Url:         m.Permalink,
		}
		results = append(results, msg)
	}

	for _, f := range files.Matches {
		file := models.SearchResult{
			AccountID:   "slack",
			Service:     "slack",
			Resource:    "file",
			Title:       f.Title,
			Description: "",
			Url:         f.Permalink,
		}
		results = append(results, file)
	}

	return results, nil
}
