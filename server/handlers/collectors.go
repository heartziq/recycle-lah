package handlers

import "net/http"

func collected(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET: list nearby pickup requests"))

	case "POST":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("POST: request pickups"))

	}

}
