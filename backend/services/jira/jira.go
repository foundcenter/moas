package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/foundcenter/moas/backend/utils"
)

const AccountType = "jira"

type JiraUser struct {
	Url          string `json:"url"`
	Self         string `json:"self"`
	Key          string `json:"key"`
	Password     string `json:"password"`
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
	jiraUser.Url = url
	jiraUser.Password = password

	addAccount(ctx, &user, jiraUser)

	user, err = db.UserRepo.Update(user)

	return user, nil
}

func addAccount(ctx context.Context, user *models.User, res *JiraUser) {
	a := models.AccountInfo{
		Type: AccountType,
		ID:   res.AccountId,
		Data: res,
		Active: true,
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

func Search(ctx context.Context, accountInfo models.AccountInfo, query string) ([]models.SearchResult, error) {

	searchResult := make([]models.SearchResult, 0)

	user := JiraUser{}
	data, _ := json.Marshal(accountInfo.Data)
	json.Unmarshal(data, &user)

	client, err := jira.NewClient(nil, user.Url)
	if err != nil {
		return searchResult, err
	}

	client.Authentication.SetBasicAuth(user.Key, user.Password)

	jqlQuery := "assignee=" + user.Key + "&(description~" + query + "|summary~" + query + ")"
	issue, _, err := client.Issue.Search(jqlQuery, nil)
	if err != nil {
		return searchResult, err
	}

	if len(issue) > 0 {
		for _, i := range issue {
			s := models.SearchResult{}
			s.Service = "jira"
			s.Resource = "issue"
			s.AccountID = accountInfo.ID
			s.Title = i.Key + " " + i.Fields.Summary
			s.Description = i.Fields.Description
			s.Url = user.Url + "browse/" + i.Key
			searchResult = append(searchResult, s)
		}
	} else {
		fmt.Print("No issues in Jira found. \n")
	}


	project, _, _ := client.Project.Get(query)

	if project != nil {
		s := models.SearchResult{}
		s.Service = "jira"
		s.Resource = "project"
		s.AccountID = accountInfo.ID
		s.Title = project.Name
		s.Description = project.Description
		s.Url = user.Url + "projects/" + project.Key
		searchResult = append(searchResult, s)

	} else {
		fmt.Print("No projects in Jira found. \n")
	}

	return searchResult, nil
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
