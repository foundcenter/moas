package github

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/foundcenter/moas/backend/config"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/foundcenter/moas/backend/utils"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githubAuth "golang.org/x/oauth2/github"
)

const AccountType = "github"

func initOAuthConfig(redirectURL string) *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     config.Settings.Github.ClientID,
		ClientSecret: config.Settings.Github.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"user:email",
		},
		Endpoint: githubAuth.Endpoint,
	}

	return config
}

func Login(ctx context.Context, code string, redirectURL string) (models.User, error) {

	var user models.User
	var github_account_id string
	config := initOAuthConfig(redirectURL)
	accessToken, err := config.Exchange(ctx, code)
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

	github_user_emails, _, err := client.Users.ListEmails(ctx, nil)
	if err != nil {
		return user, err
	}
	for _, e := range github_user_emails {
		if e.GetPrimary() {
			github_account_id = e.GetEmail()
		}
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

	addAccount(ctx, &user, github_user, github_account_id, accessToken)

	user, err = db.UserRepo.Upsert(user)

	return user, err
}

func Connect(ctx context.Context, userID string, code string, redirectURL string) (models.User, error) {

	var githubAccountId string
	var user models.User
	config := initOAuthConfig(redirectURL)
	accessToken, err := config.Exchange(ctx, code)
	if err != nil {
		return user, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	githubUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return user, err
	}

	github_user_emails, _, err := client.Users.ListEmails(ctx, nil)
	if err != nil {
		return user, err
	}
	for _, e := range github_user_emails {
		if e.GetPrimary() {
			githubAccountId = e.GetEmail()
		}
	}

	db := repo.New()
	defer db.Destroy()

	user, err = db.UserRepo.FindById(userID)
	if err != nil {
		return user, err
	}
	accessToken.Expiry = time.Now().Add(time.Hour * 24 * 365)

	addAccount(ctx, &user, githubUser, githubAccountId, accessToken)
	user, err = db.UserRepo.Update(user)

	return user, nil
}

func addAccount(ctx context.Context, user *models.User, res *github.User, primaryEmail string, token *oauth2.Token) {
	a := models.AccountInfo{
		Type:  AccountType,
		ID:    res.GetLogin(),
		Data:  res,
		Token: token,
		Active: true,
	}

	for _, acc := range user.Accounts {
		if acc.ID == a.ID && acc.Type == a.Type {
			return
		}
	}

	user.Accounts = append(user.Accounts, a)

	if primaryEmail != "" && !utils.Contains(user.Emails, primaryEmail) {
		user.Emails = append(user.Emails, primaryEmail)
	}
}

func Search(ctx context.Context, accountInfo models.AccountInfo, query string) ([]models.SearchResult, error) {

	var wg sync.WaitGroup
	resultOfSearch := make([]models.SearchResult, 0)
	queueOfResults := make(chan []models.SearchResult, 2)

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accountInfo.Token.AccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	wg.Add(3)
	go func() {
		searchResult := searchCommits(ctx, accountInfo, query, client)
		queueOfResults <- searchResult
	}()

	go func() {
		searchResult := searchIssues(ctx, accountInfo, query, client)
		queueOfResults <- searchResult
	}()

	go func() {
		searchResult := searchRepositories(ctx, accountInfo, query, client)
		queueOfResults <- searchResult
	}()

	go func() {
		for r := range queueOfResults {
			resultOfSearch = append(resultOfSearch, r...)
			wg.Done()
		}
	}()

	wg.Wait()

	return resultOfSearch, nil
}

func searchCommits(ctx context.Context, accountInfo models.AccountInfo, query string, client *github.Client) []models.SearchResult {

	userQuery := fmt.Sprintf("%s+user:%s", query, accountInfo.ID)
	result, _, _ := client.Search.Commits(ctx, userQuery, nil)
	searchResult := make([]models.SearchResult, 0)

	if len(result.Commits) > 0 {
		for _, c := range result.Commits {
			s := models.SearchResult{}
			s.Service = "github"
			s.Resource = "commit"
			s.AccountID = accountInfo.ID
			s.Title = c.Commit.GetMessage()
			s.Url = *c.HTMLURL
			searchResult = append(searchResult, s)
		}
	} else {
		fmt.Print("No commits found. \n")
	}

	return searchResult

}

func searchIssues(ctx context.Context, accountInfo models.AccountInfo, query string, client *github.Client) []models.SearchResult {
	var wg sync.WaitGroup
	resultOfSearch := make([]models.SearchResult, 0)
	queueOfResults := make(chan []models.SearchResult, 2)

	wg.Add(3)
	go func() {
		issuesQuery := fmt.Sprintf("%s+assignee:%s", query, accountInfo.ID)
		result, _, _ := client.Search.Issues(ctx, issuesQuery, nil)
		searchResult := make([]models.SearchResult, 0)
		if len(result.Issues) > 0 {
			for _, i := range result.Issues {
				s := models.SearchResult{}
				s.ID = i.GetID()
				s.Service = "github"
				s.Resource = "issue"
				s.AccountID = accountInfo.ID
				s.Description = i.Milestone.GetDescription()
				s.Url = i.GetHTMLURL()
				s.Title = i.GetTitle()
				searchResult = append(searchResult, s)
			}
		} else {
			fmt.Print("No assignee issues found. \n")
		}
		queueOfResults <- searchResult
	}()

	go func() {
		issuesQuery := fmt.Sprintf("%s+author:%s", query, accountInfo.ID)
		result, _, _ := client.Search.Issues(ctx, issuesQuery, nil)
		searchResult := make([]models.SearchResult, 0)
		if len(result.Issues) > 0 {
			for _, i := range result.Issues {
				s := models.SearchResult{}
				s.ID = i.GetID()
				s.Service = "github"
				s.Resource = "issue"
				s.AccountID = accountInfo.ID
				s.Description = i.Milestone.GetDescription()
				s.Url = i.GetHTMLURL()
				s.Title = i.GetTitle()
				searchResult = append(searchResult, s)
			}
		} else {
			fmt.Print("No issues of this author found. \n")
		}
		queueOfResults <- searchResult
	}()

	go func() {
		issuesQuery := fmt.Sprintf("%s+mentions:%s", query, accountInfo.ID)
		result, _, _ := client.Search.Issues(ctx, issuesQuery, nil)
		searchResult := make([]models.SearchResult, 0)
		if len(result.Issues) > 0 {
			for _, i := range result.Issues {
				s := models.SearchResult{}
				s.ID = i.GetID()
				s.Service = "github"
				s.Resource = "issue"
				s.AccountID = accountInfo.ID
				s.Description = i.Milestone.GetDescription()
				s.Url = i.GetHTMLURL()
				s.Title = i.GetTitle()
				searchResult = append(searchResult, s)
			}
		} else {
			fmt.Print("No mentions in issues found. \n")
		}
		queueOfResults <- searchResult
	}()

	go func() {
		for result := range queueOfResults {
			resultOfSearch = append(resultOfSearch, result...)
			wg.Done()
		}
	}()

	wg.Wait()

	//remove duplicates
	encountered := make(map[int]bool, 0)
	result := make([]models.SearchResult, 0)

	for _, v := range resultOfSearch {
		if encountered[v.ID] == true {
		} else {
			encountered[v.ID] = true
			result = append(result, v)
		}
	}

	return result

}

func searchRepositories(ctx context.Context, accountInfo models.AccountInfo, query string, client *github.Client) []models.SearchResult {

	reposQuery := fmt.Sprintf("%s user:%s", query, accountInfo.ID)
	result, _, _ := client.Search.Repositories(ctx, reposQuery, nil)
	searchResult := make([]models.SearchResult, 0)
	if len(result.Repositories) > 0 {
		for _, r := range result.Repositories {
			s := models.SearchResult{}
			s.Service = "github"
			s.Resource = "repository"
			s.AccountID = accountInfo.ID
			s.Description = r.GetDescription()
			s.Url = r.GetHTMLURL()
			searchResult = append(searchResult, s)
		}
	} else {
		fmt.Print("No repositories found. \n")
	}
	return searchResult
}
