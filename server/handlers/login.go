package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type LoginHandler struct {
	Db  *sql.DB
	Tpl *template.Template
}

type user struct {
	Name string `json:"username"`
	Pwd  string `json:"password"`
	Role string `json:"role"`
}

func CreateLoginHandler(db *sql.DB, templatePath string) *LoginHandler {
	newLoginHandler := &LoginHandler{Db: db}
	if templatePath != "" {
		newLoginHandler.SetTemplate(templatePath)
	}

	return newLoginHandler
}

// set up DB conn
func (l *LoginHandler) UseDb(db *sql.DB) error {
	l.Db = db
	return nil // return future potential error(s)
}

// setting up template
func (l *LoginHandler) SetTemplate(path string) {

	l.Tpl = template.Must(template.ParseGlob(path))
}

func (l *LoginHandler) addUser(u *user) error {

	query := "INSERT INTO your_db.user VALUES (?, ?, ?, ?)"
	userId := uuid.NewString()

	result, err := l.Db.Exec(
		query,
		userId,
		u.Name,
		u.Pwd,
		u.Role,
	)

	if err != nil {
		return errors.New("error insert into user table")
	}

	if rows, err := result.RowsAffected(); err == nil {
		log.Printf("INsert successful\t(%d) rows affected\n", rows)
		return nil
	} else {
		return err
	}
}

func (l *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryParam := mux.Vars(r)

	newUser := new(user) // add more fields for registration
	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, newUser)

	// if empty body
	if len(data) < 3 {
		http.Error(w, "400 - Empty body", http.StatusBadRequest)
		return
	}

	// Ensure required fields are NOT EMPTY
	if newUser.Name == "" || newUser.Pwd == "" {
		http.Error(w, "406 - Empty username OR password", http.StatusNotAcceptable)
		return
	}

	if queryParam["mode"] == "login" {

	} else {
		// proceed with registration
		// w.Header().Set("Content-Type", "application/json")s
		if err := l.addUser(newUser); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("insert successful!"))
	}
}
