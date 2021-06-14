// Package errlog provides four customised loggers: Trace, Info, Warning, Error and Panic
package utility

import (
	"encoding/json"
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
	Trace = log.New(io.MultiWriter(os.Stdout), "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	// if discardTrace() {
	// 	Trace = log.New(ioutil.Discard, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	// } else {
	// 	Trace = log.New(io.MultiWriter(logFile, os.Stdout), "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	// }
	Info = log.New(io.MultiWriter(os.Stdout), "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	// Info = log.New(logFile, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(errFile, logFile, os.Stderr), "Warning :", log.Ldate|log.Ltime|log.Lshortfile)
	// Error = log.New(io.MultiWriter(errFile, logFile), "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(errFile, os.Stderr), "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	Panic = log.New(io.MultiWriter(panicFile, errFile, logFile, os.Stderr), "Panic: ", log.Ldate|log.Ltime|log.Lshortfile)
	Trace.Println("errlog:logger.go - init()")
}

// func getYYYYMMDDSS() returns a string with "YYYYMMDDSSS" format.
func getYYYYMMDDSS() string {
	upToMS := strings.Split(strings.Split(time.Now().String(), "+")[0], ".")[0]
	replacer := strings.NewReplacer("-", "", " ", "-", ":", "")
	return replacer.Replace(upToMS)
}

const errlogCfgFile = "config_errlog.json"

// func discardTrace() checks the value of attribute "discard"
// in the configuration file for errlog.  It returns the boolean value
// from the file.
func discardTrace() bool {
	var config struct {
		TraceOutput struct {
			Discard bool `json:"discard"`
		} `json:"trace_output"`
	}
	if _, err := os.Stat(errlogCfgFile); err != nil {
		if os.IsNotExist(err) {
			Trace.Println("Errlog Configuration file not found!", err)
			return false
		} else { // other err
			Trace.Println("Error when checking errlog config file", err)
			return false
		}
	}
	f, err := os.Open(errlogCfgFile)
	if err != nil {
		Trace.Println("Error opening errlog configuration file", err)
		return false
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		Trace.Println("Error reading errlog configuration file", err)
		return false
	}
	if config.TraceOutput.Discard {
		return true
	}
	return false
}

// // func p() users fmt.Printf and prints parameters received.
// func p(p ...interface{}) {
// 	fmt.Println(p)
// }
