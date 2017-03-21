package gmail

import (
	"encoding/json"
	"fmt"
	"github.com/foundcenter/moas/backend/config"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/foundcenter/moas/backend/utils"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"log"
)

const (
	AccountTypeGmail = "gmail"
	AccountTypeDrive = "drive"

)

var conf *oauth2.Config

type UserGmailInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func init() {
	conf = &oauth2.Config{
		ClientID:     config.Settings.Google.ClientID,
		ClientSecret: config.Settings.Google.ClientSecret,
		RedirectURL:  config.Settings.Google.RedirectURL,
		Scopes: []string{
			"profile",
			"email",
			"https://www.googleapis.com/auth/gmail.readonly",
			"https://www.googleapis.com/auth/drive.readonly",
		},
		Endpoint: google.Endpoint,
	}
}

func Login(ctx context.Context, code string) (models.User, error) {

	var user models.User
	accessToken, err := conf.Exchange(ctx, code)
	if err != nil {
		return user, err
	}

	client := conf.Client(ctx, accessToken)

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

	user, _ = db.UserRepo.FindByAccount(gu.Email, AccountTypeGmail)

	// If user is already registered merge data
	if !user.ID.Valid() {
		user.Name = gu.Name
		user.Picture = gu.Picture
	}

	addAccount(ctx, &user, &gu, accessToken, AccountTypeGmail)
	addAccount(ctx, &user, &gu, accessToken, AccountTypeDrive)

	user, err = db.UserRepo.Upsert(user)

	return user, err
}

func Connect(ctx context.Context, userID string, code string) (models.User, error) {
	var user models.User
	accessToken, err := conf.Exchange(ctx, code)

	if err != nil {
		return user, err
	}

	client := conf.Client(ctx, accessToken)
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

	addAccount(ctx, &user, &gu, accessToken,AccountTypeGmail)
	db.UserRepo.Update(user)

	return user, nil
}

func addAccount(ctx context.Context, user *models.User, res *UserGmailInfo, token *oauth2.Token, account_type string ) {
	a := models.AccountInfo{
		Type:  account_type,
		ID:    res.Email,
		Data:  res,
		Token: token,
	}

	for _, acc := range user.Accounts {
		if acc.ID == a.ID && acc.Type == a.Type {
			return
		}
	}

	user.Accounts = append(user.Accounts, a)

	if res.Email != "" && !utils.Contains(user.Emails, res.Email) {
		user.Emails = append(user.Emails, res.Email)
	}
}

func Search(ctx context.Context, account models.AccountInfo, query string) []models.SearchResult {

	searchResult := make([]models.SearchResult, 0)
	gmailService := CreateGmailService(ctx, account.Token)
	userEmail := account.ID

	ref, err := gmailService.Users.Threads.List(userEmail).Q(query).MaxResults(100).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages. %v", err)
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

	return searchResult
}

func CreateGmailService(ctx context.Context, token *oauth2.Token) *gmail.Service {

	client := conf.Client(ctx, token)
	gmailService, err := gmail.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	return gmailService
}
