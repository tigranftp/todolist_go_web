package API

import (
	"encoding/json"
	"io"
	"net/http"
	"site/API/utils"
	"site/types"
)

type signUpJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (a *API) handleSignUp() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "can't read body", http.StatusBadRequest)
			return
		}
		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}
		var su signUpJSON
		err = json.Unmarshal(body, &su)
		if err != nil {
			http.Error(writer, "can't parse body", http.StatusInternalServerError)
			return
		}
		err = a.store.RegisterUser(su.Username, utils.GeneratePasswordHash(su.Password), su.Name, types.UserRole)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
