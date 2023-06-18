package API

import (
	"encoding/json"
	"io"
	"net/http"
	"site/API/utils"
)

type UpdateListItemJSON struct {
	Token       string `json:"token"`
	TaskID      int64  `json:"taskID"`
	Taskname    string `json:"taskname"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func (a *API) handleUpdateTodoListItem() http.HandlerFunc {
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
		var uli UpdateListItemJSON
		err = json.Unmarshal(body, &uli)
		if err != nil {
			http.Error(writer, "can't parse body", http.StatusInternalServerError)
			return
		}
		usr, err := utils.ParseUserToken(uli.Token)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		err = a.store.UpdateTodoListItem(uli.TaskID, usr.Id, uli.Taskname, uli.Description, uli.Done)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
