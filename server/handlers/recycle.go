package handlers

import "net/http"

func Recylce(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET: List of [recycle bins]..."))

	case "POST":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("POST: adding [recycle bins]..."))

	}

}
