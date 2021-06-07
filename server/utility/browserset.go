package utility

import (
	"net/http"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[\\w!#$%&'*+/=?`{|}~^-]+(?:\\.[\\w!#$%&'*+/=?`{|}~^-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,6}$") // regular expression

func IsEmailValid(e string) bool {
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

//set cookie on client computer
func SetCookie(w http.ResponseWriter, id string) {
	co := &http.Cookie{
		Name:     "Recycle Lah",
		Value:    id,
		HttpOnly: false,
		// Expires:  time.Now().AddDate(2, 0, 0),
	}
	http.SetCookie(w, co)
}
