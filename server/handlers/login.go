package handlers

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type LoginHandler struct {
}

func CreateLoginHandler() http.Handler {
	return &LoginHandler{}
}

func (l *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryParam := mux.Vars(r)
	if queryParam["mode"] == "login" {

	} else {
		// register
		if data, err := io.ReadAll(r.Body); err == nil {
			w.Header().Set("Content-Type", "application/json")

			w.Write([]byte(data))

		}
	}
}
