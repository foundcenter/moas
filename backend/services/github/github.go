package github

import (
	"context"
	"fmt"
	"github.com/foundcenter/moas/backend/config"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/foundcenter/moas/backend/utils"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"sync"
	"os/user"
)

const AccountType = "slack"

var conf *oauth2.Config

func init() {
	conf = &oauth2.Config{
		ClientID:     config.Settings.Github.ClientID,
		ClientSecret: config.Settings.Github.ClientSecret,
		RedirectURL:  config.Settings.Github.RedirectURL,
		Scopes: []string{
			"user:email",
		},
		//Endpoint: github.Endpoint,
	}
}

func Login(ctx context.Context, code string) (models.User, error) {

	var user models.User
	accessToken, err := conf.Exchange(ctx, code)
	if err != nil {
		return user, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	github_user, _, err := client.Users.Get(ctx, "")
	github_account_id := string(github_user.ID)

	if err != nil {
		return user, err
	}

	db := repo.New()
	defer db.Destroy()

	user, err = db.UserRepo.FindByAccount(AccountType, github_account_id)

	if err != nil {
		return user, err
	}

	// If user is already registered merge data
	if !user.ID.Valid() {
		user.Name = github_user.GetName()
		user.Picture = github_user.GetAvatarURL()
	}

	addAccount(ctx, &user, github_user, accessToken)

	db.UserRepo.Upsert(user)

	return user, err
}

func Connect(ctx context.Context, userID string, code string) (models.User, error) {
	var user models.User
	accessToken, err := conf.Exchange(ctx, code)
	if err != nil {
		return user, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	github_user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return user, err
	}

	db := repo.New()
	defer db.Destroy()

	user, err = db.UserRepo.FindById(userID)
	if err != nil {
		return user, err
	}

	addAccount(ctx, &user, github_user, accessToken)
	db.UserRepo.Update(user)

	return user, nil
}

func addAccount(ctx context.Context, user *models.User, res *github.User, token *oauth2.Token) {
	a := models.AccountInfo{
		Type:  AccountType,
		ID:    string(res.GetID()),
		Data:  res,
		Token: token,
	}

	for _, acc := range user.Accounts {
		if acc.ID == a.ID && acc.Type == a.Type {
			return
		}
	}

	user.Accounts = append(user.Accounts, a)

	if res.GetEmail() != "" && !utils.Contains(user.Emails, res.GetEmail()) {
		user.Emails = append(user.Emails, res.GetEmail())
	}
}

func Search(ctx context.Context, accountInfo models.AccountInfo, query string) ([]models.SearchResult, error) {

	var wg sync.WaitGroup
	searchResult := make([]models.SearchResult, 0)
	resultOfSearch := make([]models.SearchResult, 0)
	queueOfResults := make(chan []models.SearchResult, 2)

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accountInfo.Token.AccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	wg.Add(3)
	go func() {
		result, _, _ := client.Search.Commits(ctx, query, nil)
		if len(result.Commits) > 0 {
			for _, c := range result.Commits {
				s := models.SearchResult{}
				s.Service = "github"
				s.Resource = "commit"
				s.AccountID = accountInfo.ID
				s.Description = c.GetMessage()
				s.Url = "https://api.github.com/repos/" + c.GetAuthorName() + "/" + c.Repository.GetName() + "/git/commits/" + c.GetHash()
				searchResult = append(searchResult, s)
			}
		} else {
			fmt.Print("No commits found. \n")
		}
	}()

	go func() {
		result, _, _ := client.Search.Issues(ctx, query, nil)
		if len(result.Issues) > 0 {
			for _, i := range result.Issues {
				s := models.SearchResult{}
				s.Service = "github"
				s.Resource = "issue"
				s.AccountID = accountInfo.ID
				s.Description = i.String()
				s.Url = i.GetURL()
				searchResult = append(searchResult, s)
			}
		} else {
			fmt.Print("No issues found. \n")
		}
	}()

	go func() {
		result, _, _ := client.Search.Repositories(ctx, query, nil)
		if len(result.Repositories) > 0 {
			for _, r := range result.Repositories {
				s := models.SearchResult{}
				s.Service = "github"
				s.Resource = "repository"
				s.AccountID = accountInfo.ID
				s.Description = r.GetDescription()
				s.Url = r.GetURL()
				searchResult = append(searchResult, s)
			}
		} else {
			fmt.Print("No repositories found. \n")
		}
	}()

	go func() {
		for r := range queueOfResults {
			resultOfSearch = append(resultOfSearch, r...)
			wg.Done()
		}
	}()

	wg.Wait()

	return searchResult, nil
}
