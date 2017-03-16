package gmail

import (
	"fmt"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	authService "github.com/foundcenter/moas/backend/services/auth"
	"golang.org/x/net/context"
	"google.golang.org/api/gmail/v1"
	"log"
	"sync"
)

func Search(user_sub string, query string) []models.ResultResponse {

	var wg sync.WaitGroup
	var searchResult []models.ResultResponse = make([]models.ResultResponse, 0)
	gmailService := CreateGmailService(user_sub)

	user, _ := FindUserById(user_sub)
	userEmail := user.Email

	wg.Add(2)
	go func() {
		result := SearchMessages(gmailService, userEmail, query)
		searchResult = append(searchResult, result...)
		wg.Done()
	}()

	go func() {
		result := SearchThreads(gmailService, userEmail, query)
		searchResult = append(searchResult, result...)
		wg.Done()
	}()
	wg.Wait()

	return searchResult
}

func SearchMessages(gmailService *gmail.Service, userEmail string, query string) []models.ResultResponse {

	var searchResult []models.ResultResponse = make([]models.ResultResponse, 0)

	ref, err := gmailService.Users.Messages.List(userEmail).Q(query).MaxResults(50).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages. %v", err)
	}

	if len(ref.Messages) > 0 {
		for _, m := range ref.Messages {
			s := models.ResultResponse{}
			s.Service = "gmail"
			s.Resource = "messages"
			s.AccountID = userEmail
			s.Description = m.Snippet
			s.Url = "https://mail.google.com/mail/u/" + userEmail + "/#inbox/" + m.Id
			searchResult = append(searchResult, s)
		}
	} else {
		fmt.Print("No messages found.")
	}
	return searchResult
}

func SearchThreads(gmailService *gmail.Service, userEmail string, query string) []models.ResultResponse {

	var searchResult []models.ResultResponse = make([]models.ResultResponse, 0)

	ref, err := gmailService.Users.Threads.List(userEmail).Q(query).MaxResults(50).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages. %v", err)
	}

	if len(ref.Threads) > 0 {
		for _, m := range ref.Threads {
			s := models.ResultResponse{}
			s.Service = "gmail"
			s.Resource = "thread"
			s.AccountID = userEmail
			s.Description = m.Snippet
			s.Url = "https://mail.google.com/mail/u/" + userEmail + "/#inbox/" + m.Id
			searchResult = append(searchResult, s)
		}
	} else {
		fmt.Print("No threads found.")
	}
	return searchResult
}

func CreateGmailService(user_sub string) *gmail.Service {

	ctx := context.Background()

	//get user from db with user_sub=sub
	user, err := FindUserById(user_sub)
	if err != nil {
		log.Fatalf("Unable to get user: %v", err)
	}

	config := authService.GetConfig()
	client := config.Client(ctx, user.Accounts["google"])

	gmailService, err := gmail.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	return gmailService
}

func FindUserById(id string) (models.User, error) {
	db := repo.New()
	defer db.Destroy()
	err, user := db.UserRepo.FindById(id)
	return user, err
}
