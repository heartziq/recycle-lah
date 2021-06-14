package main

import (
	"encoding/base64"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"frontend/errlog"
)

// func httpLog() accepts a handlerfunc and returns anonymous (handlerfunc) function
// that print out the name of the HandlerFunc and call the HandlerFunc
// It is part of the chaining handlers
func httpLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		errlog.Info.Println("httpLog: Handler called - " + name)
		h(w, r)
	}
}

// func executeTemplate() executes templates and raises panic if error occurs
func executeTemplate(w http.ResponseWriter, filename string, data interface{}) {
	err := tpl.ExecuteTemplate(w, filename, data)
	if err != nil {
		errlog.Panic.Println(err)
	}
}

// func getCookie to get cookie and returns cookie.value
func getCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("RecycleLah")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// // gets and displays cookies (for debugging and testing purposes)
// func getCookie1(w http.ResponseWriter, r *http.Request) {
// 	cookie, err := r.Cookie("RecycleLah")
// 	if err != nil {
// 		fmt.Fprintln(w, "RecycleLah not found", err)
// 	} else {
// 		fmt.Fprintln(w, "RecycleLah :", cookie)
// 	}
// }

// func setFlashCookie() sets flash cookie to store message
// the cookie is used to display message in subsequent page
func setFlashCookie(w http.ResponseWriter, msg string) {
	cookieValue := []byte(msg)
	cookie := http.Cookie{
		Name:     "FlashCookie",
		Value:    base64.URLEncoding.EncodeToString(cookieValue),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	errlog.Trace.Printf("Set FlashCookie : %v\n", cookie)
}

// func getFlashCookie returns cookie value (message)
func getFlashCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	errlog.Trace.Println("In getFlashCookie")
	c, err := r.Cookie("FlashCookie")
	if err != nil {
		errlog.Trace.Println(err)
		if err == http.ErrNoCookie {
			errlog.Trace.Println("no cookie")
			return "", err
		}
		return "", err
	}
	errlog.Trace.Printf("in getFlashCookie %v %s\n", c.Value, c.Expires.Local().Format("02-Jan-2006 03:01:300 UTC"))
	msg, err := base64.URLEncoding.DecodeString(c.Value)
	if err != nil {
		errlog.Error.Println(err)
		msg = []byte("Error in decoding")
	}
	clearFlashCookie(w, r)
	return string(msg), err
}

// func clearFlashCookie() set flash cookie to expire
func clearFlashCookie(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("FlashCookie")
	if err != nil {
		if err == http.ErrNoCookie {
			errlog.Info.Println("FlashCookie ")
			return
		}
	}
	msg, _ := base64.URLEncoding.DecodeString(c.Value)
	errlog.Trace.Printf("clear FlashCookie with message: %s\n", msg)
	cookie := http.Cookie{
		Name:     "FlashCookie",
		HttpOnly: true,
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
	}
	http.SetCookie(w, &cookie)

}

// func getSession() returns cookie value which is the uuid
func getSession(r *http.Request) (user, error) {
	var currentUser user
	c, err := r.Cookie("RecycleLah")
	if err != nil {
		if err == http.ErrNoCookie {
			return currentUser, err
		}
		errlog.Error.Println(err)
		return currentUser, err
	}
	errlog.Trace.Printf("%v %s\n", c.Value, c.Expires.Local().Format("02-Jan-2006 03:01:300 UTC"))
	key, ok := mapSessions[c.Value]
	if ok {
		currentUser, _ = mapUsers[key]
	}
	return currentUser, err
}

// func clearSession() clears session details (delete from Map)
func clearSession(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("going to get Cookie")
	cookie, err := r.Cookie("RecycleLah")
	if err != nil {
		errlog.Trace.Println("getting cookie:", err)
		if err == http.ErrNoCookie {
			errlog.Error.Println("RecycleLah not found when logout", err)
			return
		}
		return
	}
	// delete the session data
	errlog.Trace.Println("delete session from map")
	key, ok := mapSessions[cookie.Value]
	if ok {
		delete(mapUsers, key)
	}
	delete(mapSessions, cookie.Value)
	// make the session cookie expires
	cookie1 := http.Cookie{
		Name:     "RecycleLah",
		HttpOnly: true,
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
	}
	http.SetCookie(w, &cookie1)
	return
}

//  reference: https://stackoverflow.com/questions/19991124/go-template-html-iteration-to-generate-table-from-struct
//  RangeStructer()
//  this function enables {{ range .collectionOfStruct }} to print without specifying individual field name
//  it is used in listqueue
func RangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}
	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}
	output := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		output[i] = v.Field(i).Interface()
	}
	return output
}

//  fmtDate gets time.Time and return formatted date
//  however, not working properly in html template
//  using .DateConsulted.Format "02-Jan-2006 15:04:050 UTC" in .gohtml instead
func fmtDate(input time.Time) string {
	return input.Format("02-Jan-2006 03:01:300 UTC")
	// return input.Local().Format("02-Jan-2006 03:01:300 PM")
	// return input.String()
}

//  fShortDate gets time.Time and return short date
func fShortDate(input time.Time) string {
	return input.Local().Format("02-Jan-2006")
}

// func isSessionExpired() checks if session has expired
// return true if expired
// return false if it is not
func isSessionExpired(sessionKey string) (expired bool, collector bool, err error) {
	key, ok := mapSessions[sessionKey]
	if !ok {
		// not in mapusers
		return true, false, errSessionNotFound
	}
	if user, ok := mapUsers[key]; ok {
		errlog.Trace.Println("session sessionCreatedTime date:", user.sessionCreatedTime)
		// 	now1 = time.Now().UnixNano() / int64(time.Second)
		collector = user.isCollector
		currentTime := time.Now().UnixNano() / int64(time.Second)
		errlog.Trace.Println("time.Now():", int(time.Now().Unix()))
		errlog.Trace.Println("currentTime:", currentTime)
		if currentTime > (user.sessionCreatedTime + 120*60) { // in second
			// delete the session data
			delete(mapSessions, sessionKey)
			delete(mapUsers, key)
			errlog.Info.Println("Session has expired")
			// setFlashCookie(w, "Your session has expired, please re-login")
			// message(w, r)
			return true, collector, nil
		}
		return false, collector, nil
	}
	return true, collector, errSessionNotFound
}

// func checkAccess() check if a user (session details) has the access right
// It is used as a middleware in the handlers authentication
func checkAccess(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//  get cookie
		key, err := getCookie(r)
		if err != nil {
			// have not login in - access not allowed
			errlog.Error.Println("Cookie not found, no access")
			setFlashCookie(w, "Unauthorized access")
			message(w, r)
			return
		}
		//  check if session is expired
		expired, _, err := isSessionExpired(key)
		if err != nil {
			//  session data not found
			errlog.Error.Println("Session data not found, no access")
			setFlashCookie(w, "Unauthorized access")
			message(w, r)
			return
		}
		if expired {
			setFlashCookie(w, "Your session has expired, please re-login")
			message(w, r)
			return
		}

		errlog.Trace.Println("will route to the requested page")
		h(w, r)
		return

	} // return func()
}
