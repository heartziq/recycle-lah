package handlers

import "net/http"

func Collectors(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET: i am here!"))

	case "POST":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("POST: completed job!"))

	}

}
