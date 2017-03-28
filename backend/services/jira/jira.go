package jira

import (
	"context"
	"github.com/andygrunwald/go-jira"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/foundcenter/moas/backend/utils"
)

const AccountType = "jira"

type JiraUser struct {
	Self         string `json:"self"`
	Key          string `json:"key"`
	AccountId    string `json:"accountId"`
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
}

func Connect(ctx context.Context, userID string, url string, username string, password string) (models.User, error) {

	db := repo.New()
	defer db.Destroy()

	user, err := db.UserRepo.FindById(userID)
	if err != nil {
		return user, err
	}

	jiraClient, err := jira.NewClient(nil, url)
	if err != nil {
		return user, err
	}

	jiraClient.Authentication.SetBasicAuth(username, password)

	jiraUser, err := GetMyself(jiraClient)
	if err != nil {
		return user, err
	}

	addAccount(ctx, &user, jiraUser)

	user, err = db.UserRepo.Update(user)

	return user, nil
}

func addAccount(ctx context.Context, user *models.User, res *JiraUser) {
	a := models.AccountInfo{
		Type: AccountType,
		ID:   res.AccountId,
		Data: res,
	}

	for _, acc := range user.Accounts {
		if acc.ID == a.ID && acc.Type == a.Type {
			return
		}
	}

	user.Accounts = append(user.Accounts, a)

	if res.EmailAddress != "" && !utils.Contains(user.Emails, res.EmailAddress) {
		user.Emails = append(user.Emails, res.EmailAddress)
	}
}

func GetMyself(client *jira.Client) (*JiraUser, error) {

	req, _ := client.NewRequest("GET", "/rest/api/2/myself", nil)

	myself := new(JiraUser)
	_, err := client.Do(req, myself)

	if err != nil {
		return nil, err
	}

	return myself, err
}
