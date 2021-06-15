// Web application with sign up and login functions
// To access: https://localhost:5221
package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"frontend/errlog"

	_ "github.com/go-sql-driver/mysql"
)

//application initialization
func init() {
	errlog.Trace.Println("main.go - init()")
	tpl = template.Must(template.New("").Funcs(tplFuncs).ParseGlob("templates/*"))
	files := http.FileServer(http.Dir("public"))
	http.Handle("/static/", http.StripPrefix("/static/", files))
	http.Handle("/Pictures/", http.StripPrefix("/Pictures", http.FileServer(http.Dir("Pictures"))))
	http.Handle("/Stuff/", http.StripPrefix("/Stuff", http.FileServer(http.Dir("Stuff"))))

}

func main() {
	defer func() {
		if r := recover(); r != nil {
			errlog.Info.Println("recovered in main()")
		}
		cleanUp()
	}()

	// setup handlers
	setupHandles()

	// create http server
	server := &http.Server{
		Addr: ":5221",
	}

	var wg sync.WaitGroup
	wg.Add(1)
	//  goroutine to capture os signal in order to perform orderly shutdown
	go func() {
		chanSignal := make(chan os.Signal, 1)
		signal.Notify(chanSignal, os.Interrupt)
		<-chanSignal
		errlog.Info.Println("Signal from OS, starting shutdown")
		// shutting down server
		if err := server.Shutdown(context.Background()); err != nil {
			errlog.Info.Println("error in server.Shutdown-", err)
		} else {
			errlog.Info.Println("server.Shutdown successful")
		}
		close(chanSignal)
		wg.Done()
	}()
	//  listen and serve
	errlog.Info.Printf("Listerning at port server :%s\n", server.Addr)
	fmt.Println("Listening at port", server.Addr, "since", time.Now().String())
	err := server.ListenAndServeTLS("batch3cert.pem", "batch3key.pem")
	if err == http.ErrServerClosed {
		errlog.Info.Println("server.ListenAndServer() - http.ErrServerClosed detected, wait for wg")
		//  will wait for the anonymous goroutine to complete
		wg.Wait()
	} else if err != nil {
		errlog.Info.Println("server.ListenAndServe()- ", err)
	}
}

// func setupHandles() setups handlers used in the application
func setupHandles() {
	//  home page
	http.HandleFunc("/", httpLog(index))
	// login, logout, main menu
	http.HandleFunc("/signin", httpLog(login))
	http.HandleFunc("/logout", httpLog(logout))
	http.HandleFunc("/welcome", httpLog(checkAccess(welcome)))
	// user
	http.HandleFunc("/newuser", httpLog(newUser))
	http.HandleFunc("/userupdate", userDetailUpdate)
	http.HandleFunc("/signupsuccess", httpLog(signUpSuccess))
	// reward points
	http.HandleFunc("/view_points", httpLog(checkAccess(viewPoints)))
	// pickups
	http.HandleFunc("/request_pickup", httpLog(requestPickup))
	http.HandleFunc("/user_pickup_list", httpLog(checkAccess(userPickupList)))
	http.HandleFunc("/view_Completed_Jobs", httpLog(checkAccess(viewCompletedJobs)))
	// http.HandleFunc("/req_pickup", httpLog(checkAccess(requestPickup)))

	//collector
	http.HandleFunc("/view_available_jobs", httpLog(showJobsAvailable))
	http.HandleFunc("/view_my_jobs", httpLog(showMyJobs))

	// recycle Bin
	http.HandleFunc("/indexBin", IndexBin)
	http.HandleFunc("/recyclebinsFB", recycleBinFB)
	// Get user query pass feedback
	http.HandleFunc("/queryFB", queryFB)
	// show ALL recyclebins and FeedBack.
	http.HandleFunc("/showrecyclebins", showRecycleBins)

	// Sook for testing only
	// http.HandleFunc("/index_sook", httpLog(indexSook))
	// http.HandleFunc("/dummyDeleteTask", httpLog(dummyCalldeletePickup))
	// http.HandleFunc("/dummyShowAllAvailJobsforCollector", httpLog(dummyCallCollectorShowJobsAvailable))
	// http.HandleFunc("/dummyAcceptJob", httpLog(dummyCalldummyacceptJob))
	// http.HandleFunc("/dummyViewAssignedJob", httpLog(dummyViewAttendingJob))

	// http.HandleFunc("/menu", mainMenu)

	// http.HandleFunc("/logout", logOut)
	// can delete ???

	// http.HandleFunc("/login", logIn)

}
