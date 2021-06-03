package main

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
	userhandler "github.com/heartziq/recycle-lah/server/userhandler"
	middleware "github.com/heartziq/recycle-lah/server/utility"
)

// map with key "login": handler type
var (
	handlersList = map[string]http.Handler{
		"test": &handlers.Test{},
	}
)

func createServer() http.Handler {
	db, err := sql.Open("mysql", "admin:password@tcp(127.0.0.1:3306)/recycle")
	if err != nil {
		panic(err)
	}

	if v, ok := handlersList["test"].(*handlers.Test); ok {
		v.UseDb(db)
		// v.SetTemplate("templates/test/*")
	}

	router := mux.NewRouter() // Main Router

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
		Handler(handlersList["test"])

	subR.Use(middleware.VerifyAPIKey)

	// users
	// URI: http://localhost:5000/api/v1/users/user2345?key=secretkey with json data struct: NewUser
	// curl -H "Content-Type: application/json" -X POST http://localhost:5000/api/v1/users/user4567?key=secretkey -d {\"password\":\"password\",\"email:\":\"mail\",\"collector\":false}
	// curl -H "Content-Type: application/json" -X DELETE http://localhost:5000/api/v1/users/user4567?key=secretkey
	// curl -X GET http://localhost:5000/api/v1/users/USER1234?key=secretkey
	// curl -H "Content-Type: application/json" -X GET http://localhost:5000/api/v1/users/USER4567?key=secretkey -d {\"password\":\"password\"}
	subR.
		Methods("GET", "PUT", "POST", "DELETE").
		Path("/api/v1/users/{id:[[:alnum:]]+}").
		Queries("key", "{key}").
		HandlerFunc(userhandler.Users)

	// reward points for user
	// PUT will update the points by the amount supplied in the request
	// curl -X GET http://localhost:5000/api/v1/rewards/USER1234?key=secretkey
	// curl -X PUT http://localhost:5000/api/v1/rewards/USER1234?key=secretkey -d {\"reward_points\":3}
	subR.
		Methods("GET", "PUT").
		Path("/api/v1/rewards/{id:[[:alnum:]]+}").
		Queries("key", "{key}").
		HandlerFunc(userhandler.Rewards)
	subR.Use(middleware.VerifyAPIKey)

	// test without middleware and with invalid api key
	// curl -H "Content-Type: application/json" -X POST http://localhost:5000/api/v1/testusers/user1234?key=wrongkey -d {\"password\":\"password\",\"email:\":\"mail\",\"collector\":\"false\"}
	/* test result:
	   curl -H "Content-Type: application/json" -X POST http://localhost:5000/api/v1/testusers/user1234?key=wrongkey -d {\"password\":\"password\",\"email:\":\"mail\",\"collector\":\"false\"}
	   Forbidden Access - Invalid API_KEY provided
	*/
	subR.
		Methods("GET", "PUT", "POST", "DELETE").
		Path("/api/v1/testusers/{id}").
		Queries("key", "{key}").
		HandlerFunc(userhandler.TestUsers)

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
		// fmt.Println("starting server at port 5000 WITHOUT TLS")
		// http.ListenAndServe(":5000", router)
		fmt.Println("starting server at port 5000 with TLS")
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
