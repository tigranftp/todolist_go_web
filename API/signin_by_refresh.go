package API

import (
	"encoding/json"
	"net/http"
	"site/API/utils"
)

func (a *API) handleSignInByRefresh() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		refreshToken := request.Header.Get("RefreshToken")
		usr, err := a.store.GetUserByRefresh(refreshToken)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
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
