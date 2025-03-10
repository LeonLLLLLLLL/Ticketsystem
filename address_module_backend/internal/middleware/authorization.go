package middleware

import (
	"errors"
	"net/http"

	"address_module/api"

	log "github.com/sirupsen/logrus"
)

var ErrUnAuthorizedError = errors.New("invalid username or token")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var username string = r.URL.Query().Get("username")
		var token = r.Header.Get("Authorization")
		//var err error

		if username == "" || token == "" {
			log.Error(ErrUnAuthorizedError)
			api.RequestErrorHandler(w, ErrUnAuthorizedError)
			return
		}
		/*
			#####################
			TODO: Implement with auth with real database
			#####################

			var database *tools.DatabaseInterface
			database, err = tools.NewDatabase()
			if err != nil {
				api.InternalErrorHandler(w)
				return
			}

			var loginDetails *tools.LoginDetails
			loginDetails = (*database).GetUserLoginDetails(username)

			if loginDetails == nil || (token != (*loginDetails).AuthToken) {
				log.Error(UnAuthorizedError)
				api.RequestErrorHandler(w, UnAuthorizedError)
				return
			}
		*/

		next.ServeHTTP(w, r)
	})
}
