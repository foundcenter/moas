package account

import (
	"context"
	"github.com/foundcenter/moas/backend/controllers/response"
	"github.com/foundcenter/moas/backend/middleware/jwt_auth"
	"github.com/foundcenter/moas/backend/middleware/logger"
	"github.com/foundcenter/moas/backend/repo"
	"github.com/foundcenter/moas/backend/services/auth"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func WrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		ctx := context.WithValue(r.Context(), "params", ps)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

func Load(router *httprouter.Router) {

	standardChain := alice.New(logger.Handler)
	extendedChain := standardChain.Append(jwt_auth.Handler)

	router.Handle("DELETE", "/account/:type/:id", WrapHandler(extendedChain.ThenFunc(handleAccountDelete)))

}

func handleAccountDelete(w http.ResponseWriter, r *http.Request) {

	params := r.Context().Value("params").(httprouter.Params)
	accountType := params.ByName("type")
	accountID := params.ByName("id")

	token := r.Header.Get("Authorization")
	if token == "" {
		response.Reply(w).BadRequest()
	}

	user_id, err := auth.ParseToken(token[7:])
	if err != nil {
		response.Reply(w).ServerInternalError()
		return
	}

	db := repo.New()
	defer db.Destroy()

	user, err := db.UserRepo.FindById(user_id)
	if err != nil {
		response.Reply(w).Unauthorized(err)
		return
	}

	//find account & check if can be deleted
	canBeDeleted := len(user.Accounts) > 1 && accountType !="gmail" && accountType!="drive"

	if canBeDeleted {
		for i, a := range user.Accounts {
			if a.Type == accountType && a.ID == accountID {
				user.Accounts = append(user.Accounts[:i], user.Accounts[i+1:]...)
				user, err = db.UserRepo.Update(user)
				if err != nil {
					response.Reply(w).ServerInternalError()
					return
				}

				response.Reply(w).Ok(user)
				return
			}
		}
	}

	response.Reply(w).BadRequest()

}