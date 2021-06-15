package utility

import (
	"log"
	"net/http"
	"strings"

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

func AddAuthHeader(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.Header.Set("Authorization", "Bearer somerandomtokenstring")
		next.ServeHTTP(w, r)
	})
}

func ValidateJWTToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Reject anything other than Bearer
		// expect var 'mechanism' to be []string{"Bearer", "TOKEN_STRING"}
		mechanism := strings.Split(r.Header.Get("Authorization"), " ")
		if len(mechanism) > 1 && mechanism[0] == "Bearer" {
			if token := mechanism[1]; token != "" {

				// validate token
				userid, err := VerifyToken(token)
				if err != nil {
					log.Printf("Validate token err: %v\n", err.Error())
					if err.Error() == "token expired" {
						http.Redirect(w, r, "/gimme", http.StatusPermanentRedirect)
						return
					}

					http.Error(w, "Invalid Token - Authorization Failed", http.StatusUnauthorized)

					return

				}
				log.Printf("userid: %v", userid)
				next.ServeHTTP(w, r)
				return

			}

		}

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Status Unauthorized"))
	})
}
