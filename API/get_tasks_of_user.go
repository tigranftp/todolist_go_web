package API

import (
	"encoding/json"
	"net/http"
	"site/API/utils"
)

func (a *API) handleGetTasksOfUser() http.HandlerFunc {
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
		todoList, err := a.store.GetUserTodoList(usr.Id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		dataToWrite, err := json.Marshal(todoList)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write(dataToWrite)
	}
}
