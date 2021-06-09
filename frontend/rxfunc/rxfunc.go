// Package contains regular expressions functions for user name and password
package rxfunc

import "regexp"

// Must start and end with letters or digits, may contain . and underscore
// true  aa aaa   a..a a_._a
// false  aa. a a[space]a
var UserNameExp = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_\\.]*[a-zA-Z0-9]$")

// Must start and end with letters or digits, may contain . and underscore
// true  aa aaa   a..a a_._a
// false  aa. a a[space]a
var PasswordExp = regexp.MustCompile("^[\\S]+$")

// to test user name pattern
func UserNamePattern(s string) bool {
	matched := UserNameExp.MatchString(s)
	return matched
}

// to test password pattern
func PasswordPattern(s string) bool {
	matched := PasswordExp.MatchString(s)
	return matched
}
