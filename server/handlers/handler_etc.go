package handlers

import (
	"runtime/debug"

	errlog "github.com/heartziq/recycle-lah/server/utility"
	"golang.org/x/crypto/ssh/terminal"

	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// func cleanUp() calls functions to perform the necessary cleanup before the application terminates.
// It is normally called from the defer function in the main().
func cleanUp() {
	errlog.Trace.Println("in cleanUp()")
	errlog.Info.Println("    start cleaning up, closing database")
	// dbCleanUp()
	errlog.Info.Println("    database closed.  Cleanup complete.")
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
		errlog.Error.Printf("recovered, r=%+v\n", r)
		debug.PrintStack()
	}
}

// var UserNameExp = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_\\.]*[a-zA-Z0-9]$")

// // func UserNamePattern() matches a given string with UserNameExp.
// // It returns true if it matches.
// // Otherwise, it returns false.
// func UserNamePattern(s string) bool {
// 	matched := UserNameExp.MatchString(s)
// 	return matched
// }
