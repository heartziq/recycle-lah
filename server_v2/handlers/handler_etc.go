package handlers

import (
	"regexp"
	"runtime/debug"

	errlog "github.com/heartziq/recycle-lah/server/utility"

	_ "github.com/go-sql-driver/mysql"
)

// func recoverFunc() uses recover built-in function to
// recover from panic.  It allows the calling function to recover from panic
// when it is called a defer function.
func recoverFunc() {
	if r := recover(); r != nil {
		errlog.Error.Printf("recovered, r=%+v\n", r)
		debug.PrintStack()
	}
}

// this regular expression UserNameExp and func UserNamePattern()
// could be in a package shared by the front-end web server
// var UserNameExp contains regular expression for User Id
var UserNameExp = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_\\.]*[a-zA-Z0-9]$")

// func UserNamePattern() matches a given string with UserNameExp.
// It returns true if it matches.
// Otherwise, it returns false.
func UserNamePattern(s string) bool {
	matched := UserNameExp.MatchString(s)
	return matched
}

// // const cfgFile defines the name of the configuration file
// // for database connection and port number of the web server.
// // The file is in json format.
// const cfgFile = "./config.json"

// // var config Configuration is a global variable that stores the configuration details.
// var config Configuration

// // struct Configuration defines attributes in the configuration file
// // and its corresponding json attributes name.
// type Configuration struct {
// 	DB struct {
// 		Name   string `json:"name"` // database name
// 		User   string `json:"user"` // database user
// 		Host   string `json:"host"`
// 		Port   string `json:"port"`
// 		driver string
// 	} `json:"database"`
// 	Server struct {
// 		Port string `json:"port"`
// 	} `json:"server"`
// }
