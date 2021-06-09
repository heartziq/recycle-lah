// Package errlog provides four customised loggers: Trace, Info, Warning, Error and Panic
package errlog

import (
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	Trace   *log.Logger // log for debugging purposes
	Info    *log.Logger // log for general important information
	Warning *log.Logger // log for situation requires attention
	Error   *log.Logger // error log
	Panic   *log.Logger // panic log
)

// initializes loggers
func init() {
	errFile, err := os.OpenFile("./log/errors.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicln(err)
	}
	panicFile, err := os.OpenFile("./log/panic.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicln(err)
	}
	logFile, err := os.OpenFile("./log/log"+getYYYYMMDDSS()+".txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicln(err)
	}
	// Trace: currently logs in logfile, may set to ioutil.Discard at a later stage
	Trace = log.New(io.MultiWriter(logFile, os.Stdout), "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	// Trace = log.New(logFile, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(io.MultiWriter(logFile, os.Stdout), "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	// Info = log.New(logFile, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(errFile, logFile, os.Stderr), "Warning :", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(errFile, logFile, os.Stderr), "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	Panic = log.New(io.MultiWriter(panicFile, errFile, logFile, os.Stderr), "Panic: ", log.Ldate|log.Ltime|log.Lshortfile)
	Trace.Println("errlog:logger.go - init()")
}

// returns string with YYYYMMDDSSS format
func getYYYYMMDDSS() string {
	upToMS := strings.Split(strings.Split(time.Now().String(), "+")[0], ".")[0]
	replacer := strings.NewReplacer("-", "", " ", "-", ":", "")
	return replacer.Replace(upToMS)
}
