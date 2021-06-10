package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"frontend/errlog"
)

// logs http calls
func httpLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		errlog.Info.Println("httpLog: Handler called - " + name)
		h(w, r)
	}
}

// func executeTemplate1(w http.ResponseWriter, data interface{}, files ...string) {
// 	var tmplList []string
// 	for _, file := range files {
// 		tmplList = append(tmplList, fmt.Sprintf("templates/%s.gohtml", file))
// 	}
// 	templates := template.Must(template.ParseFiles(tmplList...))
// 	templates.ExecuteTemplate(w, "layout", data)
// }

// executes templates and raises panic if error occurs
func executeTemplate(w http.ResponseWriter, filename string, data interface{}) {
	err := tpl.ExecuteTemplate(w, filename, data)
	if err != nil {
		errlog.Panic.Println(err)
	}
}

// err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
// if err != nil {
// 	// panicLog("Execute Template error ", err)
// 	log.Panic(err)
// }

// func XparseTemplateFiles(filenames ...string) (t *template.Template) {
// 	var files []string
// 	t = template.New("layout")
// 	for _, file := range filenames {
// 		files = append(files, fmt.Sprintf("templates/%s.gohtml", file))
// 	}
// 	t = template.Must(t.ParseFiles(files...))
// 	return
// }

func getCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("RecycleLah")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// gets and displays cookies (for debugging and testing purposes)
func getCookie1(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("RecycleLah")
	if err != nil {
		fmt.Fprintln(w, "RecycleLah not found", err)
	} else {
		fmt.Fprintln(w, "RecycleLah :", cookie)
	}
}

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

// returns session details
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

// // returns true when session information is created successfully
// // returns false otherwise
// func createSession(w http.ResponseWriter, userName string) bool {
// 	defer func() {
// 		errlog.Trace.Println("running defer func() in createSession():")
// 		if r := recover(); r != nil {
// 			errlog.Info.Println("Recovering from createSession()", r)
// 		}
// 	}()
// 	id := uuid.NewV4()
// 	cookie := http.Cookie{
// 		Name:     "RecycleLah",
// 		Value:    id.String(),
// 		HttpOnly: true,
// 	}
// 	http.SetCookie(w, &cookie)
// 	errlog.Trace.Printf("cookie set: %+v\n", cookie)
// 	//  create session record
// 	i := user.InsertSession(db, id.String(), userName)
// 	if i == 0 { //  record not added
// 		errlog.Error.Println("Failed to create session record for uuid:", id.String())
// 		return false
// 	}
// 	return true
// }

// // returns session details
// // func getSession(w http.ResponseWriter, r *http.Request) (Session, error) {
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

// clears session details
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

// func checkAccess(h http.HandlerFunc) http.HandlerFunc {
// 	errlog.Trace.Println("in checkAccess()")
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		sess, err := getSession(w, r)
// 		if err != nil {
// 			errlog.Error.Println("checkAdminAccess: Cookie not found, no access")
// 			unauthorized(w, r)
// 			return
// 		}
// 		errlog.Trace.Printf("checkAccess session:%+v", sess)
// 		//  should include check to if session expires
// 		h(w, r)
// 	}
// }

//  checkCollectorAccess check if a user (session details) has the access right
func checkCollectorAccess(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//  get cookie
		key, err := getCookie(r)
		if err != nil {
			// probably have not login in - access not allowed
			errlog.Error.Println("Cookie not found, no access")
			setFlashCookie(w, "Unauthorized access")
			message(w, r)
			return
		}
		//  check if session is expired
		expired, collector, err := isSessionExpired(key)
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
		if !collector {
			//  collector trying to access user web page
			errlog.Info.Println("account has no user access")
			setFlashCookie(w, "Unauthorized access")
			// route to a message page
			message(w, r)
			return
		}
		errlog.Trace.Println("will route to the requested page")
		h(w, r)
		return

	} // return func()
}

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
		if currentTime > (user.sessionCreatedTime + 30*60) {
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

// //  checkUserAccess check if a user (session details) has the access right
// func checkAccess(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		//  get cookie
// 		cookie, err := r.Cookie("RecycleLah")
// 		if err != nil {
// 			// probably have not login in - access not allowed
// 			errlog.Error.Println("checkCollectorAccess: Cookie not found, no access")
// 			index(w, r)
// 		} else {
// 			//  got cookie
// 			if userSession, ok := mapSession[cookie.Value]; ok {
// 				errlog.Trace.Println("session sessionCreatedTime date:", userSession.sessionCreatedTime)
// 				// 	now1 = time.Now().UnixNano() / int64(time.Second)
// 				currentTime := time.Now().UnixNano() / int64(time.Second)
// 				errlog.Trace.Println("time.Now():", int(time.Now().Unix()))
// 				errlog.Trace.Println("currentTime:", currentTime)
// 				if currentTime > (userSession.sessionCreatedTime + 1*60) {
// 					// delete the session data
// 					delete(mapSession, cookie.Value)
// 					errlog.Info.Println("Session has expired")
// 					// how to inform session has expired
// 					// route to a message page
// 					setFlashCookie(w, "Your session has expired, please re-login")
// 					message(w, r)
// 					return
// 				}
// 				//  found session in session map
// 				if !userSession.isCollector {
// 					//  call handler as user has access
// 					h(w, r)
// 					return
// 				} else {
// 					// log access violation and render index page
// 					errlog.Info.Println("account has no user access", userSession.userId)
// 					setFlashCookie(w, "Unauthorized access")
// 					// route to a message page
// 					message(w, r)
// 					return
// 				}

// 			} else {
// 				//  cannot find session and render index page
// 				errlog.Error.Println(w, "session data not found")
// 				setFlashCookie(w, "Please login through our home page")
// 				message(w, r)
// 				return
// 			}

// 		}

// 	}
// }

//  checkUserAccess check if a user (session details) has the access right
func checkUserAccess(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//  get cookie
		key, err := getCookie(r)
		if err != nil {
			// probably have not login in - access not allowed
			errlog.Error.Println("Cookie not found, no access")
			setFlashCookie(w, "Unauthorized access")
			message(w, r)
			return
		}
		//  check if session is expired
		expired, collector, err := isSessionExpired(key)
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
		if collector {
			//  collector trying to access user web page
			errlog.Info.Println("account has no user access")
			setFlashCookie(w, "Unauthorized access")
			// route to a message page
			message(w, r)
			return
		}
		errlog.Trace.Println("will route to the requested page")
		h(w, r)
		return

	} // return func()
}
