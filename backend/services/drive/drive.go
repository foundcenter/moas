package drive

import (
	"fmt"
	"github.com/foundcenter/moas/backend/models"
	"github.com/foundcenter/moas/backend/repo"
	authService "github.com/foundcenter/moas/backend/services/auth"
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v2"
	"log"
)

func Search(user_sub string, query string) []models.ResultResponse {

	var searchResult []models.ResultResponse = make([]models.ResultResponse, 0)
	driveService := CreateDriveService(user_sub)

	user, err := FindUserById(user_sub)
	userEmail := user.Email

	ref, err := driveService.Files.List().Q("fullText contains '" + query + "'").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files. %v", err)
	}

	if len(ref.Items) > 0 {
		for _, f := range ref.Items {
			s := models.ResultResponse{}
			s.AccountID = userEmail
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

	return searchResult
}

func CreateDriveService(user_sub string) *drive.Service {

	ctx := context.Background()

	//get user from db with user_sub=sub
	user, err := FindUserById(user_sub)
	if err != nil {
		log.Fatalf("Unable to get user: %v", err)
	}

	config := authService.GetConfig()
	client := config.Client(ctx, user.Accounts["google"])

	driveService, err := drive.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	return driveService
}

func FindUserById(id string) (models.User, error) {
	db := repo.New()
	defer db.Destroy()
	err, user := db.UserRepo.FindById(id)
	return user, err
}
