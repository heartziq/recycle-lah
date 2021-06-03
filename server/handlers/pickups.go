package handlers

import (
	"database/sql"
	"encoding/json"
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

// DB queries
func (p *Pickup) ListPickup() (users []*user) {
	// access db
	results, err := p.Db.Query("SELECT username FROM my_db.users")

	if err != nil {

		panic(err.Error())

	}

	for results.Next() {

		// map this type to the record in the table

		c := user{}

		err = results.Scan(&c.UserName)

		if err != nil {

			panic(err.Error())

		}

		users = append(users, &c)

	}
	return
}

func (p *Pickup) ShowPickup() http.HandlerFunc {
	d := func(w http.ResponseWriter, r *http.Request) {
		data := p.ListPickup()
		json.NewEncoder(w).Encode(data)
	}

	return http.HandlerFunc(d)
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
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("[collector] View All Pickups that I Accepted"))
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
