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

var tpl *template.Template
var tplFuncs = template.FuncMap{"rangeStruct": RangeStructer, "fShortDate": fShortDate, "fmtDate": fmtDate}

//application initialization
func init() {
	errlog.Trace.Println("main.go - init()")
	// dbConn := DBConnection{
	// 	dbType:   "mysql",
	// 	user:     "goappuser",
	// 	password: "password",
	// 	hostAddr: "127.0.0.1",
	// 	port:     "3306",
	// 	name:     "goinaction2",
	// }
	// s := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConn.user, dbConn.password, dbConn.hostAddr, dbConn.port, dbConn.name)
	// var err error
	// db, err = openDB(dbConn.dbType, s)
	// if err != nil {
	// 	errlog.Panic.Panicln(err)
	// }
	// errlog.Trace.Printf("DB Name: %s, DB User:%s\n", dbConn.name, dbConn.user)

	tpl = template.Must(template.New("").Funcs(tplFuncs).ParseGlob("templates/*"))
	files := http.FileServer(http.Dir("public"))
	http.Handle("/static/", http.StripPrefix("/static/", files))
	http.Handle("/Pictures/", http.StripPrefix("/Pictures", http.FileServer(http.Dir("Pictures"))))
	http.Handle("/Stuff/", http.StripPrefix("/Stuff", http.FileServer(http.Dir("Stuff"))))

}

func main() {
	fmt.Println("in main - merging with Sin Yaw's")
	// now := int(time.Now().Unix())
	// nextMin := int(time.Now().Add(time.Hour * 1).Unix())
	// fmt.Println("now=", now)
	// fmt.Println("now=", nextMin)
	// fmt.Println("diff=", nextMin-now)

	// i := int(time.Now().Add(time.Minute * 2).Unix())
	// fmt.Println("i=", i)
	// dt := time.Now()
	// fmt.Println("current dt", dt.String())
	// now1 := time.Now().UnixNano() / int64(time.Millisecond)
	// fmt.Println("now1 (millisecond)=", now1)
	// now1 = time.Now().UnixNano() / int64(time.Second)
	// fmt.Println("now1 (Second)=", now1)
	// t1 := time.Unix(now1, 0)
	// fmt.Println("t1 (time.Unix(now1, 0))=", t1)
	// now2 := time.Now().UnixNano() / int64(time.Minute)
	// fmt.Println("now2 (minute =", now2)
	// t1 = time.Unix(now2, 0)
	// fmt.Println("t1 (time.Unix(now2, 0))=", t1)

	// i, err := strconv.ParseInt(userSession.sessionCreatedTime, 10, 64)
	// fmt.Println(i)
	// tm := time.Unix(userSession.sessionCreatedTime, 0)

	// ok
	//recycle()
	// ok
	// recyclePost()
	// ok
	// recyclePostUser()
	// err1 := testaddUser() - not ok
	// if err1 != nil {
	// 	fmt.Println(err1)
	// }
	// os.Exit(1)

	defer func() {
		if r := recover(); r != nil {
			errlog.Info.Println("recovered in main()")
		}
		errlog.Trace.Println("=====defer main() before cleanUp()=====")
		cleanUp()
		errlog.Trace.Println("======defer main() cleanUp() ended.=====")
	}()
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

// func XsetupHandlesSook() {
// 	http.HandleFunc("/", httpLog(index))
// 	http.HandleFunc("/contact", httpLog(contact))
// 	http.HandleFunc("/getcookie", httpLog(httpLog(getCookie1)))
// 	http.HandleFunc("/newuser", httpLog(httpLog(createAccount)))
// 	http.HandleFunc("/signupsuccess", httpLog(httpLog(signUpSuccess)))
// 	http.HandleFunc("/signin", httpLog(httpLog(login)))
// 	http.HandleFunc("/welcome", httpLog(checkUserAccess(welcome)))
// 	http.HandleFunc("/view_points", httpLog(checkUserAccess(viewPoints)))
// 	// http.HandleFunc("/email", httpLog(checkAccess(email)))

// 	http.HandleFunc("/logout", httpLog(logout))
// 	http.HandleFunc("/unauthorized", httpLog(unauthorized))
// 	http.HandleFunc("/test1", httpLog(checkUserAccess(testToken)))
// 	http.HandleFunc("/message", httpLog(message))

// 	http.HandleFunc("/collector_signin", httpLog(httpLog(collectorLogin)))
// 	http.HandleFunc("/collector_welcome", httpLog(checkCollectorAccess(collectorWelcome)))
// }

func setupHandles() {
	// Sook
	http.HandleFunc("/", httpLog(index))
	// Sook - create new account
	http.HandleFunc("/newusersook", httpLog(httpLog(createAccount)))
	http.HandleFunc("/signupsuccess", httpLog(httpLog(signUpSuccess)))
	// Sook - login/logout
	http.HandleFunc("/signin", httpLog(httpLog(login)))

	http.HandleFunc("/welcome", httpLog(checkUserAccess(welcome)))
	http.HandleFunc("/logout", httpLog(logout))
	// reward points
	http.HandleFunc("/view_points", httpLog(checkUserAccess(viewPoints)))
	//Main Pages
	// http.HandleFunc("/mainmenu", httpLog(mainMenu))
	http.HandleFunc("/menu", mainMenu)
	http.HandleFunc("/newuser", httpLog(newUser))

	// http.HandleFunc("/logout", logOut)
	http.HandleFunc("/login", logIn)
	// //Sub Pages
	// http.HandleFunc("/userupdate", userDetailUpdate)
	// http.HandleFunc("/pickup", pickUp)
	// http.HandleFunc("/viewstatus", viewStatus)
	//Load Files
	// http.Handle("/Pictures/", http.StripPrefix("/Pictures", http.FileServer(http.Dir("Pictures"))))
	// http.Handle("/Stuff/", http.StripPrefix("/Stuff", http.FileServer(http.Dir("Stuff"))))
	//run server

}
