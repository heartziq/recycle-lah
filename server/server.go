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

func createServer() http.Handler {
	// db, err := sql.Open("mysql", "admin:password@tcp(127.0.0.1:3306)/recycle")
	db, err := sql.Open("mysql", "admin:password@tcp(127.0.0.1:3306)/recycle?parseTime=true")
	if err != nil {
		panic(err)
	}

	// Initialize handlers
	pickup := handlers.CreatePickupHandler(db, "")
	reward := handlers.CreateRewardHandler(db, "")
	user := handlers.CreateUserHandler(db, "")
	router := mux.NewRouter() // Main Router

	// Protected route - need to supply API_KEY
	subR := router.NewRoute().Subrouter()
	// URI: https://localhost:5000/api/v1/pickups/4?key=secretkey&limit=true&role=collector
	subR.
		Methods("GET", "PUT", "POST", "DELETE").
		Path("/api/v1/pickups/{id:\\d+}").
		Queries("key", "{key}").
		// Queries("limit", "{limit}").
		Queries("role", "{role:user|collector}").
		Handler(pickup)

	subR.Use(middleware.VerifyAPIKey)

	// test
	if v, ok := handlersList["test"].(*handlers.Test); ok {
		v.UseDb(db)
		// v.SetTemplate("templates/test/*")
	}

	//
	// A route "test" that use AddAuthHeader middleware
	//
	subR = router.NewRoute().Subrouter()
	subR.
		Methods("GET", "PUT", "POST", "DELETE").
		Path("/api/v1/test/{id:\\d+}").
		Queries("key", "{key}").
		Handler(handlersList["test"])

	subR.Use(middleware.AddAuthHeader)

	// /*

	// users
	// URI: http://localhost:5000/api/v1/users/user2345?key=secretkey with json data struct: NewUser
	// curl -H "Content-Type: application/json" -X POST http://localhost:5000/api/v1/users/user4567?key=secretkey -d {\"password\":\"password\",\"email:\":\"mail\",\"collector\":false}
	// v2: curl -H "Content-Type: application/json" -X POST http://localhost:5000/api/v1/users/curl1111?key=secretkey -d {\"password\":\"password\",\"email:\":\"curl1111@gmail.com\",\"user_name\":\"curl1111_cat\",\"collector\":false}
	// curl -H "Content-Type: application/json" -X DELETE http://localhost:5000/api/v1/users/user4567?key=secretkey
	// curl -X GET http://localhost:5000/api/v1/users/USER1234?key=secretkey
	// curl -H "Content-Type: application/json" -X GET http://localhost:5000/api/v1/users/USER4567?key=secretkey -d {\"password\":\"password\"}
	// PUT for change password and username
	subRUser := router.NewRoute().Subrouter()
	subRUser.
		Methods("GET", "PUT", "POST", "DELETE").
		Path("/api/v1/users/{id:[[:alnum:]]+}").
		Queries("key", "{key}").
		// HandlerFunc(userhandler.Users)
		Handler(user)
	subRUser.Use(middleware.HttpLog)
	subRUser.Use(middleware.VerifyAPIKey)

	// reward points for user
	// PUT will update the points by the amount supplied in the request
	// curl -X GET http://localhost:5000/api/v1/rewards/USER1234?key=secretkey
	// curl -X PUT http://localhost:5000/api/v1/rewards/USER1234?key=secretkey -d {\"reward_points\":3}
	subRReward := router.NewRoute().Subrouter()
	subRReward.
		Methods("GET", "PUT").
		Path("/api/v1/rewards/{id:[[:alnum:]]+}").
		Queries("key", "{key}").
		// HandlerFunc(userhandler.Rewards)
		Handler(reward)
	subRReward.Use(middleware.HttpLog)
	subRReward.Use(middleware.VerifyAPIKey)
	subRReward.Use(middleware.VerifyHdrToken)

	// Public route
	router.HandleFunc("/api/v1/pickups", pickup.ShowPickup())
	// recycle
	router.HandleFunc("/api/v1/recycle", handlers.Recylce)

	return router
}

func main() {
	// Instantiate server
	router := createServer()

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
	log.Println("Doing cleanup...")
	log.Println("done cleaning up")

}
