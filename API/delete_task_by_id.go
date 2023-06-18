package API

import (
	"net/http"
	"site/API/utils"
	"strconv"
)

func (a *API) handleDeleteTaskByID() http.HandlerFunc {
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
		taskID, err := strconv.ParseInt(request.Header.Get("TaskID"), 10, 64)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		err = a.store.DeleteTaskByID(taskID, usr.Id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
