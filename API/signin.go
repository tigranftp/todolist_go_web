package API

import (
	"encoding/json"
	"io"
	"net/http"
	"site/API/utils"
)

type signInJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInResponseJSON struct {
	SessionToken string `json:"session_token"`
	RefreshToken string `json:"refresh_token"`
}

func (a *API) handleSignIn() http.HandlerFunc {
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
		var si signInJSON
		err = json.Unmarshal(body, &si)
		if err != nil {
			http.Error(writer, "can't parse body", http.StatusInternalServerError)
			return
		}
		usr, err := a.store.GetUser(si.Username, utils.GeneratePasswordHash(si.Password))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		st, err := utils.NewSessionUserToken(usr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		rt, rtEAT, err := utils.NewRefreshToken()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		err = a.store.UpdateRefreshToken(usr.Id, rt, rtEAT)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := json.Marshal(signInResponseJSON{
			SessionToken: st,
			RefreshToken: rt,
		})
		writer.WriteHeader(http.StatusOK)
		writer.Write(res)
	}
}
