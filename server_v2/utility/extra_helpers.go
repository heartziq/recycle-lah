package utility

import (
	"database/sql"
	"errors"
	"regexp"
	"runtime/debug"

	// errlog "github.com/heartziq/recycle-lah/server/utility"
	"golang.org/x/crypto/ssh/terminal"

	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// type DBConnectionError error defines error type for database connection
type DBConnectionError error

var (
	errDBOpen = DBConnectionError(errors.New("failed to open database"))
	errDBPing = DBConnectionError(errors.New("failed to ping database"))
)

// func openDB() connects to the database and returns a connected driver/
// It prompts for the database user password and uses the config database information
// to connect to the database.  It performs a ping to make sure the connection is in order.
func openDB1() (*sql.DB, error) {
	fmt.Println("Press enter to use the default password.")
	pwd, err := getPasswordInput()
	if err != nil {
		Error.Println("Error when getting password input", err)
	}
	if strings.TrimSpace(pwd) == "" {
		pwd = "password" // for convenience of code testing
		fmt.Printf("*** use default password *** \n")
	}

	s := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DB.User, pwd, config.DB.Host, config.DB.Port, config.DB.Name)
	db, err := sql.Open(config.DB.driver, s)
	if err != nil {
		Error.Println(err.Error())
		return db, errDBOpen
	}
	if err = db.Ping(); err != nil {
		Error.Println(err.Error())
		return db, errDBPing
	}
	Info.Println("DB opened")
	return db, nil
}

func OpenDB() (*sql.DB, error) {
	var err error
	db, err := sql.Open("mysql", "admin:password@tcp(127.0.0.1:3306)/recycle?parseTime=true")
	if err != nil {
		Error.Println(err.Error())
		return db, errDBOpen
	}
	if err = db.Ping(); err != nil {
		Error.Println(err.Error())
		return db, errDBPing
	}
	Info.Println("DB opened")
	return db, nil
}

// func dbCleanUp() closes the database.
func dbCleanUp(db *sql.DB) {
	defer recoverFunc()
	db.Close()
}

// const cfgFile defines the name of the configuration file
// for database connection and port number of the web server.
// The file is in json format.
const cfgFile = "./config.json"

// var db *sql.DB os a global variable for the database driver
var db *sql.DB

// var config Configuration is a global variable that stores the configuration details.
var config Configuration

// struct Configuration defines attributes in the configuration file
// and its corresponding json attributes name.
type Configuration struct {
	DB struct {
		Name   string `json:"name"` // database name
		User   string `json:"user"` // database user
		Host   string `json:"host"`
		Port   string `json:"port"`
		driver string
	} `json:"database"`
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}

// func printConfig() prints the configuration details.
func printConfig() {
	fmt.Printf("Database\n   DB Name: %s\n   Host: %s\n   Port: %s\n   User Name: %s\n", config.DB.Name, config.DB.Host, config.DB.Port, config.DB.User)
	fmt.Printf("HTTP Server\n   Port: %s\n", config.Server.Port)
}

// func loadConfig() reads the configuration file
// defined in const cfgFile.  The file is in json format.
// The content of the file is decoded and stored in const config.
// It also sets the database driver to "mysql".
func loadConfig() {
	if !isFileExist(cfgFile) {
		Panic.Println("Error opening configuration file")
	}
	f, err := os.Open(cfgFile)
	if err != nil {
		Panic.Println("Error reading configuration file", err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		Trace.Println("Error decoding configuration file", err)
	}
	config.DB.driver = "mysql"
	Trace.Println("after loading config:", config)
	printConfig()
}

// func getArgs() reads the arguments from the command line.
// It calls setServerPort() to set http server port if "http_port=" or "http_port:" is supplied.
// It also calls setDBPort() to set database server port if "db_port=" or "db_port:" is supplied.
func getArgs() {
	Info.Println("command string:", os.Args)
	args := os.Args[1:]
	for i, v := range args {
		Trace.Printf("arg %d :%v\n", i, v)
		arg := strings.ToLower(strings.TrimSpace(v))
		if strings.Contains(arg, "http_port=") || strings.Contains(arg, "http_port:") {
			setServerPort(arg)
		} else if strings.Contains(arg, "db_port=") || strings.Contains(arg, "db_port:") {
			setDBPort(arg)
		}
	}
}

// func setServerPort() checks for string to setup http server port
// and sets config.Server.Port when matching string is found.
func setServerPort(input string) {
	if strings.Contains(input, "http_port=") || strings.Contains(input, "http_port:") {
		var s []string
		if strings.Contains(input, "http_port=") {
			s = strings.Split(input, "http_port=")
		} else {
			s = strings.Split(input, "http_port:")
		}
		if len(s) == 2 { //[0]=blank [1]=[1234]
			if string(s[1]) == "" {
				Trace.Println("Empty server port number supplied")
			} else {
				if i, err := strconv.Atoi(string(s[1])); err != nil {
					Error.Println("atoi error when converting server port number", err)
				} else {
					Trace.Println("http port is ", i)
					if i > 0 {
						config.Server.Port = s[1]
						Info.Println("Setting Server port number to: ", config.Server.Port)
					}
				}
			}
		} else { // len(s) != 2
			Trace.Println("arguments not valid) != 2")
		}
	}
}

// func setDBPort() checks for string to setup database server port
// and sets config.DB.Port  when matching string is found.
func setDBPort(input string) {
	if strings.Contains(input, "db_port=") || strings.Contains(input, "db_port:") {
		var s []string
		if strings.Contains(input, "db_port=") {
			s = strings.Split(input, "db_port=")
		} else {
			s = strings.Split(input, "db_port:")
		}
		if len(s) == 2 { //[0]=blank [1]=[1234]
			if string(s[1]) == "" {
				Trace.Println("Empty database port number supplied")
			} else {
				if i, err := strconv.Atoi(string(s[1])); err != nil {
					Error.Println("atoi error when converting server port number", err)
				} else {
					Trace.Println("db port is ", i)
					if i > 0 {
						config.DB.Port = s[1]
						Info.Println("Setting Database Server port number to: ", config.DB.Port)
					}
				}
			}
		} else { // len(s) != 2
			Trace.Println("arguments not valid) != 2")
		}
	}
}

// func isFileExist() calls os.Stat() to determine if a file exists.
// It returns true if it does and returns false when the file is not found.
func isFileExist(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			Trace.Println("File:", fileName, " not found!", err)
			return false
		} else { // other err
			Trace.Println("Error when checking file:", fileName, err)
			return false
		}
	}
	return true
}

// func cleanUp() calls functions to perform the necessary cleanup before the application terminates.
// It is normally called from the defer function in the main().
func cleanUp() {
	Trace.Println("in cleanUp()")
	Info.Println("    start cleaning up, closing database")
	// dbCleanUp()
	Info.Println("    database closed.  Cleanup complete.")
}

// func p() uses fmt.Printf and prints parameters received.
func p(p ...interface{}) {
	fmt.Printf("%+v\n", p)
}

// func getPasswordInput() gets password for a db user.
// It returns the input string.
func getPasswordInput() (string, error) {
	var s string
	fmt.Printf("Please enter database password for user [%s]: ", config.DB.User)
	b, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("...")
	s = string(b)
	if err != nil {
		return "", err
	}
	return s, nil
}

// func recoverFunc() uses recover built-in function to
// recover from panic.  It allows the calling function to recover from panic
// when it is called a defer function.
func recoverFunc() {
	if r := recover(); r != nil {
		Error.Printf("recovered, r=%+v\n", r)
		debug.PrintStack()
	}
}

var UserNameExp = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_\\.]*[a-zA-Z0-9]$")

// func UserNamePattern() matches a given string with UserNameExp.
// It returns true if it matches.
// Otherwise, it returns false.
func UserNamePattern(s string) bool {
	matched := UserNameExp.MatchString(s)
	return matched
}
