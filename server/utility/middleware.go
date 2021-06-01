package utility

import (
	"net/http"

	"github.com/gorilla/mux"
)

func VerifyAPIKey(next http.Handler) http.Handler {
	newHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		// if key exist
		if key, exist := vars["key"]; exist {

			// if valid key
			if key == "secretkey" {
				next.ServeHTTP(w, r)
				return
			}

		}

		// else return unauthorized
		http.Error(w, "Forbidden Access - Invalid API_KEY provided", http.StatusUnauthorized) // 401

	})

	return newHandlerFunc
}
