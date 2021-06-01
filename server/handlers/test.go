package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type Test struct {
	Db  *sql.DB
	Tpl *template.Template
	// error logging
	// Info  *log.Logger
	Error *log.Logger
}

type user struct {
	Id       string `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	APIKey   string `json:"api_key"`
	Count    int    `json:"count"`
	Expiry   string `json:"expiry"`
}

// DB queries
func (t *Test) ListPickup() (users []*user) {
	// access db
	results, err := t.Db.Query("SELECT username FROM my_db.users")

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

func (t *Test) UseDb(db *sql.DB) error {
	t.Db = db
	return nil // return future potential error(s)
}

func (t *Test) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// template
	if t.Tpl != nil {
		err := t.Tpl.ExecuteTemplate(w, "index.gohtml", struct{}{})
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}
	if r.Method == "POST" {
		list := t.ListPickup()
		json.NewEncoder(w).Encode(list)
		return
	}

	w.Write([]byte("get test only"))

}

// setting up template
func (t *Test) SetTemplate(path string) {

	t.Tpl = template.Must(template.ParseGlob(path))
}

// example of logging
func (t *Test) SetErrorLog() {
	// Set up output file
	var f *os.File
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic")
			path, _ := os.Getwd()
			f, _ = os.Create(path + "/error.log")
		}
	}()

	f = loadFile("errors.log")
	t.Error = log.New(io.MultiWriter(os.Stdout, f), "Error: ", log.Ldate|log.Ltime|log.Lshortfile)

}

func loadFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic("directory not found")
	}

	return f
}
