package utility

import (
	"net/http"
	"reflect"
	"runtime"

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
		http.Error(w, "Forbidden Access - Invalid API_KEY provided", http.StatusUnauthorized) // 401
	})

	return newHandlerFunc
}
