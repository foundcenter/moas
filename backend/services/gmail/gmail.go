package gmail

import (
	"fmt"
	"github.com/foundcenter/moas/backend/controllers/response"
	"github.com/foundcenter/moas/backend/repo"
	authService "github.com/foundcenter/moas/backend/services/auth"
	"golang.org/x/net/context"
	"google.golang.org/api/gmail/v1"
	"log"
	"net/http"
)


func HandleGmailSearch(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	//get user subject
	jwt_token := r.Header.Get("Authorization")
	err, user_sub := authService.ParseToken(jwt_token[7:])

	//get user from db with user_sub=sub
	db := repo.New()
	defer db.Destroy()
	err, user := db.UserRepo.FindById(user_sub)
	if err != nil {
		log.Fatalf("Unable to get user: %v", err)
	}

	config, err := authService.GetConfig()
	//google.ConfigFromJSON()
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := config.Client(ctx, user.Accounts["google"])

	gmailService, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	userEmail := "me"
	ref, err := gmailService.Users.Threads.List(userEmail).MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages. %v", err)
	}

	if len(ref.Threads) > 0 {
		fmt.Print("Messages:\n")
		for _, m := range ref.Threads {
			fmt.Printf("- %s\n", m.Snippet)
		}
	} else {
		fmt.Print("No messages found.")
	}

	response.Reply(w).Ok("ok")
}
