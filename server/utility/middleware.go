package utility

import (
	"log"
	"net/http"
	"reflect"
	"runtime"
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

func Logging(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.Header.Set("Authorization", "Bearer somerandomtokenstring")
		next.ServeHTTP(w, r)
	})
}

// handler
func HttpLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name()
		Info.Println("httpLog: Handler called - " + name)
		next.ServeHTTP(w, r)
	})

}

// handlerFunc
func HttpLog1(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		Info.Println("httpLog: Handler called - " + name)
		h(w, r)
	}
}

func VerifyHdrToken(next http.Handler) http.Handler {
	Trace.Println("============ verifying token =====================")
	newHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		Trace.Println("Token=", token)
		verified, err := VerifyToken(token)
		if err != nil {
			Error.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized) // 401	- unable to verified???
			return
		}
		Trace.Println("verified=", verified)
		if verified {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Forbidden Access - Invalid TOKEN provided", http.StatusUnauthorized) // 401
	})

	return newHandlerFunc
}

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
						// http.Redirect(w, r, "/gimme", http.StatusPermanentRedirect)
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

func ValidateJWTToken_Sook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Trace.Println("Sook in ValidateJWTToken")
		// Reject anything other than Bearer
		// expect var 'mechanism' to be []string{"Bearer", "TOKEN_STRING"}
		mechanism := strings.Split(r.Header.Get("Authorization"), " ")
		Trace.Println("Sook in mechanism=", mechanism)
		if len(mechanism) > 1 && mechanism[0] == "Bearer" {
			Trace.Println("Sook in if len")
			if token := mechanism[1]; token != "" {
				Trace.Println("Sook in if token")
				// validate token
				if _, err := VerifyToken(token); err != nil {
					log.Printf("Validate token err: %v\n", err.Error())
					if err.Error() == "token expired" {
						http.Error(w, err.Error(), http.StatusUnauthorized)
						// http.Redirect(w, r, "/gimme", http.StatusPermanentRedirect)
						return
					}
					Trace.Println("Sook in veriftoken err not nil ")
					http.Error(w, "Invalid Token - Authorization Failed", http.StatusUnauthorized)

					return

				}
				Trace.Println("Sook in ValidateJWTToken - token verified")
				next.ServeHTTP(w, r)
				return

			}

		}
		Trace.Println("Sook in ValidateJWTToken -last part")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Status Unauthorized"))

	})
}
