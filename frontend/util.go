package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	// "regexp"

	"frontend/errlog"

	"frontend/rxfunc"
)

type SessionError error

var (
	errSessionNotFound = SessionError(errors.New("Session not found"))
)

type ValidationError error //  input validation error type

var (
	errUserNameFmt = ValidationError(errors.New("Must start/end with letters or digits, may contain . and underscore.  e.g.:user1234, user.1234"))
	// errUserNameFmt      = ValidationError(errors.New("Only letters(a-z), numbers(0-9) and periods(.) are allowed"))
	errNoUserName       = ValidationError(errors.New("Username must not be empty"))
	errUserNameNotAvail = ValidationError(errors.New("Username is not available"))
	errUserNameLength   = ValidationError(errors.New("username length must be between 6 and 30"))
	errNoPassword       = ValidationError(errors.New("missing password"))
	errConfirmPassword  = ValidationError(errors.New("Confirmation password not matched"))
	errPasswordLength   = ValidationError(errors.New("password length must be between 8 and 64"))
	errPasswordFormat   = ValidationError(errors.New("use 8 or more characters with a mix of letters, numbers & symbols"))
	errEmailFormat      = ValidationError(errors.New("letters, numbers & periods"))
	errEmailFormat2     = ValidationError(errors.New("username must be between 6 and 30 characters"))
	errExceedAttempt    = ValidationError(errors.New("Exceed three attempts"))
)

// returns true if input string passes validation for username
// returns false otherwise
func checkInputUserName(s string) error {
	if strings.TrimSpace(s) == "" {
		return errNoUserName
	}
	if len(s) < 6 || len(s) > 30 {
		return errUserNameLength
	}
	// check for username format
	ok := rxfunc.UserNamePattern(s)
	if !ok {
		return errUserNameFmt
	}
	return nil
}

// returns true if input string passes validation for password
// returns false otherwise
func checkInputNewPassword(s string) error {
	if strings.TrimSpace(s) == "" {
		return errNoPassword
	}
	if len(s) < 8 || len(s) > 64 {
		return errPasswordLength
	}
	//  use regex to check for password combination
	ok := rxfunc.PasswordPattern(s)
	if !ok {
		return errPasswordFormat
	}
	return nil
}

// returns true if both input passwords are the same
// returns false otherwise
func confirmPassword(password1 string, password2 string) bool {
	if password1 != password2 {
		return false
	}
	return true
}

// simulates cleanUp process
func cleanUp() {
	errlog.Info.Println("    start cleaning up, closing database")
	errlog.Info.Println("    database closed.  Cleanup complete.")
}

// func printConfig() prints the configuration details.
func printConfig() {
	fmt.Printf("Base URL: %s\n", config.BaseURL)
	fmt.Printf("API Key: %s\n", config.APIKey)
	fmt.Printf("Account: %s\n", config.Account)
}

// func loadConfig() reads the configuration file
// defined in const cfgFile.  The file is in json format.
// The content of the file is decoded and stored in const config.
func loadConfig() {
	if !isFileExist(cfgFile) {
		errlog.Panic.Println("Error opening configuration file")
	}
	f, err := os.Open(cfgFile)
	if err != nil {
		errlog.Panic.Println("Error reading configuration file", err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		errlog.Trace.Println("Error decoding configuration file", err)
	}
	printConfig()
}

// func isFileExist() calls os.Stat() to determine if a file exists.
// It returns true if it does and returns false when the file is not found.
func isFileExist(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			errlog.Trace.Println("File:", fileName, " not found!", err)
			return false
		} else { // other err
			errlog.Trace.Println("Error when checking file:", fileName, err)
			return false
		}
	}
	return true
}
