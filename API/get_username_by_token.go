package API

import (
	"net/http"
	"site/API/utils"
)

func (a *API) handleGetUsernameByToken() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		bearerToken := extractTokenFromHeader(request.Header.Get("Authorization"))
		if bearerToken == "" {
			http.Error(writer, "no session token passed", http.StatusUnauthorized)
			return
		}
		usr, err := utils.ParseUserToken(bearerToken)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(usr.Username))
	}
}
