package API

import (
	"encoding/json"
	"io"
	"net/http"
	"site/API/utils"
	"strconv"
)

type AddTaskForUserJSON struct {
	Token       string `json:"token"`
	Taskname    string `json:"taskname"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func (a *API) handleAddTaskForUser() http.HandlerFunc {
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
		var atfu AddTaskForUserJSON
		err = json.Unmarshal(body, &atfu)
		if err != nil {
			http.Error(writer, "can't parse body", http.StatusInternalServerError)
			return
		}
		usr, err := utils.ParseUserToken(atfu.Token)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		lastID, err := a.store.AddUserListItem(usr.Id, atfu.Taskname, atfu.Description, atfu.Done)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(strconv.FormatInt(lastID, 10)))
	}
}
