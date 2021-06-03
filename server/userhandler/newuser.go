package userhandler

import (
	"strings"
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
	ok := UserNamePattern(s)
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
	// ok := rxfunc.PasswordPattern(s)
	ok := true
	if !ok {
		return errPasswordFormat
	}
	return nil
}
