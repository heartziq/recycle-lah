package main

// remove calling  middleare AddAuthHeader
import (
	"database/sql"
	"fmt"
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
var (
	handlersList = map[string]http.Handler{
		"test": &handlers.Test{},
	}
)

func createServer(db *sql.DB) http.Handler {
	// db, err := sql.Open("mysql", "admin:password@tcp(127.0.0.1:3306)/recycle")
	db, err := sql.Open("mysql", "admin:password@tcp(127.0.0.1:3306)/recycle?parseTime=true")
	if err != nil {
		panic(err)
	}
	// handlers.DBCon = db
	// Initialize handlers
	pickUpHandler := handlers.CreatePickupHandler(db, "")
	recycleBinHandler := handlers.CreateRBinHandler(db, "")
	reward := handlers.CreateRewardHandler(db, "")
	user := handlers.CreateUserHandler(db, "")
	router := mux.NewRouter() // Main Router

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

	//
	// endpoint: users //
	//
	subRUser := router.NewRoute().Subrouter()
	subRUser.
		Methods("GET", "PUT", "POST", "DELETE").
		Path("/api/v1/users/{id:[[:alnum:]]+}").
		Queries("key", "{key}").
		Handler(user)
	subRUser.Use(middleware.HttpLog)
	subRUser.Use(middleware.VerifyAPIKey)

	//
	// endpoint: rewards //
	//
	subRReward := router.NewRoute().Subrouter()
	subRReward.
		Methods("GET", "PUT").
		Path("/api/v1/rewards/{id:[[:alnum:]]+}").
		Queries("key", "{key}").
		// HandlerFunc(userhandler.Rewards)
		Handler(reward)
	subRReward.Use(middleware.HttpLog)
	subRReward.Use(middleware.VerifyAPIKey)
	subRReward.Use(middleware.ValidateJWTToken)

	return router
}

func main() {
	// create db connect
	db, err := sql.Open("mysql", "admin:password@tcp(127.0.0.1:3306)/recycle?parseTime=true")
	if err != nil {
		panic(err)
	}

	// Instantiate server
	router := createServer(db)

	// Serve
	go func() {
		fmt.Println("starting server at port 5000 WITHOUT TLS")
		http.ListenAndServe(":5000", router)
		// fmt.Println("starting server at port 5000 with TLS")
		// http.ListenAndServeTLS(":5000", "cert/cert.pem", "cert/key.pem", router)
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt) // User abruptly quit - Ctrl-C
	<-c

	// Do some cleaning ups before shutdown
	log.Println("INterrupt.. closing connection...")
	if err := db.Close(); err != nil { // Close db connection
		log.Printf("error closing db connection: %v", err)
	}
	log.Println("done cleaning up")
}
