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

func createServer(db *sql.DB) http.Handler {

	// Initialize handlers
	pickUpHandler := handlers.CreatePickupHandler(db, "")
	recycleBinHandler := handlers.CreateRBinHandler(db, "")

	// Init Main Router
	router := mux.NewRouter()

	// Protected route - need to supply API_KEY
	subR := router.NewRoute().Subrouter()

	//
	// endpoint: Pickups //
	//
	subR.
		Methods("GET", "PUT", "POST", "DELETE").
		Path("/api/v1/pickups/{id}").
		Queries("key", "{key}").
		Queries("role", "{role:user|collector}").
		Handler(pickUpHandler)

	subR.Use(middleware.VerifyAPIKey)

	router.HandleFunc("/api/v1/pickups", pickUpHandler.ShowPickup())

	//
	// endpoint: RecycleBin //
	//
	router.
		Methods("GET", "POST").
		Path("/api/v1/recyclebindetails/{userID:\\w+|NIL}"). // set to NIL or integer
		Handler(recycleBinHandler)

	return router
}

func main() {
	// create db connect
	db, err := sql.Open("mysql", "user1:password@tcp(127.0.0.1:3306)/your_db")
	if err != nil {
		panic(err)
	}

	// create server
	router := createServer(db)

	// Serve
	go func() {
		http.ListenAndServeTLS(":5000", "cert/cert.pem", "cert/key.pem", router)
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt) // User abruptly quit - Ctrl-C
	<-c

	// Do some cleaning ups before shutdown
	log.Println("INterrupt.. closing DB connection...")
	if err := db.Close(); err != nil { // Close db connection
		log.Printf("error closing db connection: %v", err)
	}

	log.Println("done cleaning up")

}
