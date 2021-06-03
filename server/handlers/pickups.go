package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Pickup struct {
	Db  *sql.DB
	Tpl *template.Template
	// error logging
	// Info  *log.Logger
	Error *log.Logger
}

func ShowPickup(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("List pickups [public]"))
}

func (p *Pickup) UseDb(db *sql.DB) error {
	p.Db = db
	return nil // return future potential error(s)
}

func (p *Pickup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// User
	vars := mux.Vars(r)
	// limit := vars["limit"]
	role := vars["role"]

	if role == "user" {
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("[user] View All Pickups that I Created"))
		case "POST":
			data := map[string]string{"address": "[user] Request for a pickup (Create new pickup)"}
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(data)
		case "PUT":
			data := map[string]string{"func": "[user] Approve a pickup"}
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(data)
		case "DELETE":
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("[user] Delete a pickup"))
		}
	} else {
		// collector
		switch r.Method {
		case "GET": // show current pickup that I am attending
			data := getUser()
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("[collector] View All Pickups that I Accepted"))
			json.NewEncoder(w).Encode(data)
			fmt.Println(data)

		case "PUT": // cancel or accept
			reqBody, err := ioutil.ReadAll(r.Body)
			newValue := make(map[string]string)
			if err == nil {
				json.Unmarshal(reqBody, &newValue)

				if newValue["collector_id"] == "" {
					// cancel accepted pickup
					log.Println("[collector] cancel a pickup")
				} else {
					// accept a pickup
					log.Println("[collector] accept a pickup")
				}
			}
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(newValue)

		}
	}
}

// setting up template
func (p *Pickup) SetTemplate(path string) {

	p.Tpl = template.Must(template.ParseGlob(path))
}
