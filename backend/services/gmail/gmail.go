package gmail

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/foundcenter/moas/backend/config"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/foundcenter/moas/backend/utils"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

const (
	AccountTypeGmail = "gmail"
	AccountTypeDrive = "drive"
)

type UserGmailInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func initOAuthConfig(redirectURL string) *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     config.Settings.Google.ClientID,
		ClientSecret: config.Settings.Google.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"profile",
			"email",
			"https://www.googleapis.com/auth/gmail.readonly",
			"https://www.googleapis.com/auth/drive.readonly",
		},
		Endpoint: google.Endpoint,
	}

	return config
}

func Login(ctx context.Context, code string, redirectURL string) (models.User, error, string) {

	var user models.User
	config := initOAuthConfig(redirectURL)
	accessToken, err := config.Exchange(ctx, code)
	if err != nil {
		return user, err, ""
	}
	client := config.Client(ctx, accessToken)

	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return user, err, ""
	}
	defer userInfo.Body.Close()

	decoder := json.NewDecoder(userInfo.Body)
	var gu UserGmailInfo
	err = decoder.Decode(&gu)
	if err != nil {
		return user, err, ""
	}

	db := repo.New()
	defer db.Destroy()

	user, _ = db.UserRepo.FindByAccount(gu.Email, AccountTypeGmail)

	// If user is already registered merge data
	if !user.ID.Valid() {
		user.Name = gu.Name
		user.Picture = gu.Picture
	}

	addAccount(ctx, &user, &gu, accessToken, AccountTypeGmail)
	addAccount(ctx, &user, &gu, accessToken, AccountTypeDrive)

	user, err, action := db.UserRepo.Upsert(user)

	return user, err, action
}

func Connect(ctx context.Context, userID string, code string, redirectURL string) (models.User, error) {
	var user models.User
	config := initOAuthConfig(redirectURL)
	accessToken, err := config.Exchange(ctx, code)

	if err != nil {
		return user, err
	}

	client := config.Client(ctx, accessToken)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return user, err
	}
	defer userInfo.Body.Close()

	decoder := json.NewDecoder(userInfo.Body)
	var gu UserGmailInfo
	err = decoder.Decode(&gu)
	if err != nil {
		return user, err
	}

	db := repo.New()
	defer db.Destroy()

	user, err = db.UserRepo.FindById(userID)

	if err != nil {
		return user, err
	}

	addAccount(ctx, &user, &gu, accessToken, AccountTypeGmail)
	user, err = db.UserRepo.Update(user)

	return user, nil
}

func addAccount(ctx context.Context, user *models.User, res *UserGmailInfo, token *oauth2.Token, account_type string) {
	a := models.AccountInfo{
		Type:   account_type,
		ID:     res.Email,
		Data:   res,
		Token:  token,
		Active: true,
	}

	for i, acc := range user.Accounts {
		if acc.ID == a.ID && acc.Type == a.Type {
			user.Accounts[i].Active = true
			return
		}
	}

	user.Accounts = append(user.Accounts, a)

	if res.Email != "" && !utils.Contains(user.Emails, res.Email) {
		user.Emails = append(user.Emails, res.Email)
	}
}

func Search(ctx context.Context, account models.AccountInfo, query string) ([]models.SearchResult, error) {

	searchResult := make([]models.SearchResult, 0)
	gmailService := CreateGmailService(ctx, account.Token)
	userEmail := account.ID

	ref, err := gmailService.Users.Threads.List(userEmail).Q(query).MaxResults(100).Do()
	if err != nil {
		fmt.Printf("Unable to retrieve threads. %v", err)
		return searchResult, err
	}

	if len(ref.Threads) > 0 {
		for _, m := range ref.Threads {
			s := models.SearchResult{}
			s.Service = "gmail"
			s.Resource = "thread"
			s.AccountID = userEmail
			s.Description = m.Snippet
			s.Url = "https://mail.google.com/mail/u/" + userEmail + "/#inbox/" + m.Id
			searchResult = append(searchResult, s)
		}
	} else {
		fmt.Print("No threads found. \n")
	}

	return searchResult, nil
}

func CreateGmailService(ctx context.Context, token *oauth2.Token) *gmail.Service {

	config := initOAuthConfig("")
	client := config.Client(ctx, token)
	gmailService, err := gmail.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	return gmailService
}
