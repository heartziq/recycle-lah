package utility

import (
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
)

// func VerifyAPIKey() verify API key provided in the param
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

func Logging(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.Header.Set("Authorization", "Bearer somerandomtokenstring")
		next.ServeHTTP(w, r)
	})
}

// func HttpLog provides logging of the name of the next function/handler
func HttpLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name()
		Info.Println("httpLog: Handler called - " + name)
		next.ServeHTTP(w, r)
	})

}

//  func VerifyHdrToken() get token from header and call VerifyToken - not in used
//  replaced by ValidateJWTToken()
// func VerifyHdrToken(next http.Handler) http.Handler {
// 	newHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		token := r.Header.Get("Authorization")
// 		Trace.Println("Token=", token)
// 		verified, err := VerifyToken(token)
// 		if err != nil {
// 			Error.Println(err)
// 			http.Error(w, err.Error(), http.StatusUnauthorized) // 401	- unable to verified???
// 			return
// 		}
// 		Trace.Println("verified=", verified)
// 		if verified {
// 			next.ServeHTTP(w, r)
// 			return
// 		}
// 		http.Error(w, "Forbidden Access - Invalid TOKEN provided", http.StatusUnauthorized) // 401
// 	})

// 	return newHandlerFunc
// }

// func ValidateJWTToken() verifies JWT token
func ValidateJWTToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Reject anything other than Bearer
		// expect var 'mechanism' to be []string{"Bearer", "TOKEN_STRING"}
		mechanism := strings.Split(r.Header.Get("Authorization"), " ")
		if len(mechanism) > 1 && mechanism[0] == "Bearer" {
			if token := mechanism[1]; token != "" {

				// validate token
				if _, err := VerifyToken(token); err != nil {
					log.Printf("Validate token err: %v\n", err.Error())
					if err.Error() == "token expired" {
						http.Error(w, err.Error(), http.StatusUnauthorized)
						return
					}

					http.Error(w, "Invalid Token - Authorization Failed", http.StatusUnauthorized)

					return

				}

				next.ServeHTTP(w, r)
				return

			}

		}

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Status Unauthorized"))

	})
}
