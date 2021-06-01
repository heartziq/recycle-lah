package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	handlers "github.com/heartziq/recycle-lah/server/handlers"
	middleware "github.com/heartziq/recycle-lah/server/utility"
)

// map with key "login": handler type

func createServer() http.Handler {
	db, err := sql.Open("mysql", "user1:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()

	// Protected route - need to supply API_KEY
	subR := router.NewRoute().Subrouter()

	// pickups
	subR.
		Methods("GET", "PUT", "POST", "DELETE").
		Path("/api/v1/pickups/{id:\\d+}").
		Queries("key", "{key}").
		HandlerFunc(handlers.Pickups)

	subR.Use(middleware.VerifyAPIKey)

	// test
	subR.
		Methods("GET", "PUT", "POST", "DELETE").
		Path("/api/v1/test/{id:\\d+}").
		Queries("key", "{key}").
		Handler(&handlers.Test{
			Db: db,
		})

	subR.Use(middleware.VerifyAPIKey)

	// Public route

	// recycle
	router.HandleFunc("/api/v1/recycle", handlers.Recylce)

	return router
}

func main() {
	// Instantiate server
	router := createServer()

	// Serve
	go func() {
		http.ListenAndServeTLS(":5000", "cert/cert.pem", "cert/key.pem", router)
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt) // User abruptly quit - Ctrl-C
	<-c

	// Do some cleaning ups before shutdown
	log.Println("INterrupt.. closing connection...")
	log.Println("Doing cleanup...")
	log.Println("done cleaning up")

}
