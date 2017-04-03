package drive

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
	"google.golang.org/api/drive/v2"
)

const AccountType = "drive"

var conf *oauth2.Config

type UserDriveInfo struct {
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
			"https://www.googleapis.com/auth/drive.readonly",
		},
		Endpoint: google.Endpoint,
	}

	return config
}

func Login(ctx context.Context, code string, redirectURL string) (models.User, error) {

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
	var gu UserDriveInfo
	err = decoder.Decode(&gu)
	if err != nil {
		return user, err
	}

	db := repo.New()
	defer db.Destroy()

	user, _ = db.UserRepo.FindByAccount(gu.Email, AccountType)

	// If user is already registered merge data
	if !user.ID.Valid() {
		user.Name = gu.Name
		user.Picture = gu.Picture
	}

	addAccount(ctx, &user, &gu, accessToken)

	user, err = db.UserRepo.Upsert(user)

	return user, err
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
	var gu UserDriveInfo
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

	addAccount(ctx, &user, &gu, accessToken)
	user, err = db.UserRepo.Update(user)

	return user, nil
}

func addAccount(ctx context.Context, user *models.User, res *UserDriveInfo, token *oauth2.Token) {
	a := models.AccountInfo{
		Type:   AccountType,
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

func Search(ctx context.Context, account models.AccountInfo, query string) ([]models.SearchResult,error) {

	var searchResult []models.SearchResult = make([]models.SearchResult, 0)
	driveService := CreateDriveService(ctx, account.Token)


	ref, err := driveService.Files.List().Q("fullText contains '" + query + "'").Do()
	if err != nil {
		fmt.Printf("Unable to retrieve files. %v", err)
		return searchResult, err
	}

	if len(ref.Items) > 0 {
		for _, f := range ref.Items {
			s := models.SearchResult{}
			s.AccountID = account.ID
			s.Service = "drive"
			s.Resource = "file"
			s.Description = f.Description
			s.Url = f.AlternateLink
			s.Title = f.Title
			searchResult = append(searchResult, s)
		}
	} else {
		fmt.Print("No files found. \n")
	}

	return searchResult, nil
}

func CreateDriveService(ctx context.Context, token *oauth2.Token) *drive.Service {

	client := conf.Client(ctx, token)
	driveService, err := drive.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve drive client %v", err)
	}

	return driveService
}
