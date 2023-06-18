package API

import (
	"encoding/json"
	"net/http"
	"site/API/utils"
	"strings"
)

func (a *API) testParseToken() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		bearerToken := extractTokenFromHeader(request.Header.Get("Authorization"))
		if bearerToken == "" {
			http.Error(writer, "no session token passed", http.StatusUnauthorized)
			return
		}
		usr, err := utils.ParseUserToken(bearerToken)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}
		res, err := json.Marshal(usr)
		if err != nil {
			http.Error(writer, "error during json marshaling user: "+err.Error(), http.StatusUnauthorized)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write(res)
	}
}

func extractTokenFromHeader(authHeader string) string {
	// Проверка формата заголовка
	if authHeader == "" {
		return ""
	}

	// Разделение заголовка на тип и токен
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}
