package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frontend/errlog"
	"log"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type nUser struct {
	UserName  string
	Email     string
	Password  string
	uuid      string
	Collector bool
}

//Web Main Pages func from below---------------------------------------------------------------------------//
//New User
// JUne 8 - SOOK modified lines around these
// registeration(ToDB)             // send data to server
// matchPassword[email] = password //**temporary use
// mapSessions[id] = email
// mapUsers[email] = user{username, id}
func newUser(res http.ResponseWriter, req *http.Request) {

	Data := struct {
		PageName  string
		UserName  string
		MsgToUser string
	}{PageName: "New User Registration", UserName: ""}
	tpl.ExecuteTemplate(res, "RL_NewUser.gohtml", Data)

	//this struct is to store into DB

	// var currentUser user
	if req.Method == http.MethodPost {
		//---save user's information in the map ---
		userId := req.FormValue("userid")
		userName := req.FormValue("username")
		password := req.FormValue("password")
		confirmpassword := req.FormValue("confirmpassword")
		email := req.FormValue("email")
		cBool := req.FormValue("collector")
		var collector bool
		if cBool == "true" {
			collector = true
		} else {
			collector = false
		}
		// leave it to backend server to return error
		// _, nameFound := matchPassword[email] //**need to get from database
		// if nameFound {
		// 	// fmt.Println("Htmlmain.newUser - name in existing data")
		// 	fmt.Scanf(Data.MsgToUser, "Please use other email! '%v' has been taken!", email)
		// 	defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		// }
		if err := checkInputUserName(userId); err != nil {
			Data.MsgToUser = "Please enter correct format for user id!"
			fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			return
		}
		if strings.TrimSpace(userName) == "" {
			userName = userId
		}
		UFok := isEmailValid(email)
		if !UFok {
			Data.MsgToUser = "Please enter correct email format!"
			fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			return
		}
		if password != confirmpassword {
			Data.MsgToUser = "Confirm Password is not same!"
			fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			return
		}
		// fmt.Println("Htmlmain.newUser - no same name in existing data")
		// id := uuid.NewV4().String()
		// setCookie(res, id)
		// ToDB := nUser{userId, userName, email, password, collector}
		// log.Println(ToDB, confirmpassword)
		// registration(ToDB) // send data to server
		// SOOKMODIFIED matchPassword[email] = password //**temporary use
		// mapSessions[id] = email
		// mapUsers[email] = user{username, id}
		Data.MsgToUser = "New User Registration Done! You may proceed to log in."
		fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		var newUser NewUser
		newUser.Password = password
		newUser.Email = email
		newUser.UserName = userName
		newUser.Collector = collector
		err := addUser(userId, newUser)
		if err != nil {
			Data.MsgToUser = "Failed to add user"
			fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		}

	}

}

// func newUser(res http.ResponseWriter, req *http.Request) {

// 	Data := struct {
// 		PageName  string
// 		UserName  string
// 		MsgToUser string
// 	}{PageName: "New User Registration", UserName: ""}
// 	tpl.ExecuteTemplate(res, "RL_NewUser.gohtml", Data)

// 	//this struct is to store into DB

// 	// var currentUser user
// 	if req.Method == http.MethodPost {
// 		//---save user's information in the map ---
// 		username := req.FormValue("username")
// 		password := req.FormValue("password")
// 		confirmpassword := req.FormValue("confirmpassword")
// 		email := req.FormValue("email")
// 		C_Bool := req.FormValue("collector")
// 		var Collector bool
// 		if C_Bool == "true" {
// 			Collector = true
// 		} else {
// 			Collector = false
// 		}
// 		_, nameFound := matchPassword[email] //**need to get from database
// 		UFok := isEmailValid(email)
// 		if !nameFound && password == confirmpassword && UFok {
// 			// fmt.Println("Htmlmain.newUser - no same name in existing data")
// 			id := uuid.NewV4().String()
// 			setCookie(res, id)
// 			ToDB := nUser{username, email, password, id, Collector}
// 			log.Println(ToDB, confirmpassword)
// 			registeration(ToDB)             // send data to server
// 			matchPassword[email] = password //**temporary use
// 			mapSessions[id] = email
// 			mapUsers[email] = user{username, id}
// 			Data.MsgToUser = "New User Registration Done! You may process to log in."
// 			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
// 		} else if !UFok {
// 			// fmt.Println("Htmlmain.newUser - email format not correct")
// 			Data.MsgToUser = "Please enter correct email!"
// 			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
// 		} else if nameFound {
// 			// fmt.Println("Htmlmain.newUser - name in existing data")
// 			fmt.Scanf(Data.MsgToUser, "Please use other email! '%v' has been taken!", email)
// 			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
// 		} else if password != confirmpassword {
// 			// fmt.Println("Htmlmain.newUser - confirm password not match")
// 			Data.MsgToUser = "Confirm Password is not same!"
// 			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
// 		}
// 	}

// }

//log in
// sooK modified
func logIn(res http.ResponseWriter, req *http.Request) {

	Data := struct {
		PageName  string
		UserName  string
		MsgToUser string
	}{PageName: "Log In"}
	myCookie, err := req.Cookie("RecycleLah")
	if err == nil {
		Data.UserName = mapUsers[mapSessions[myCookie.Value]].userId
	} else {
		Data.UserName = ""
	}

	if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			errlog.Error.Println("err in ParseForm", err)
			Data.MsgToUser = "processing error"
			return
		}
		userId := req.FormValue("userid")
		password := req.FormValue("password")
		log.Println(userId, password)
		if strings.TrimSpace(userId) == "" || strings.TrimSpace(password) == "" {
			Data.MsgToUser = "All fields are mandatory!"
			fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		} else {
			var reqData UserVerification
			reqData.Password = password
			apiResp, err := verifyUser(userId, reqData)
			if err != nil {
				Data.MsgToUser = err.Error()
				errlog.Trace.Println("validateUser", err)
				return
			}
			id := uuid.NewV4()
			cookie := http.Cookie{
				Name:     "RecycleLah",
				Value:    id.String(),
				HttpOnly: true,
			}
			http.SetCookie(res, &cookie)
			errlog.Trace.Printf("cookie set: %+v\n", cookie)
			mapSessions[cookie.Value] = userId
			var currentUser user
			currentUser.userId = userId
			// userSession.sessionCreatedTime = int(time.Now().Add(time.Minute * 2).Unix())
			currentUser.sessionCreatedTime = time.Now().UnixNano() / int64(time.Second)
			// i, err := strconv.ParseInt(userSession.sessionCreatedTime, 10, 64)
			// fmt.Println(i)
			// tm := time.Unix(userSession.sessionCreatedTime, 0)
			currentUser.isCollector = apiResp.IsCollector
			currentUser.email = apiResp.Email
			currentUser.userName = apiResp.UserName
			currentUser.token = apiResp.token
			mapUsers[userId] = currentUser
			errlog.Trace.Println("    !!!!!!       currentUser=", currentUser)
			http.Redirect(res, req, "/menu", http.StatusSeeOther)
			return
		}
	}
	tpl.ExecuteTemplate(res, "Login.gohtml", Data)
}

//Main Menu
func mainMenu(res http.ResponseWriter, req *http.Request) {

	Data := struct {
		PageName string
		UserName string
	}{PageName: "Main Menu"}
	myCookie, err := req.Cookie("RecycleLah")
	if err != nil {
		// fmt.Println("Htmlmain.mainMenu - Cookie cannot find")
		http.Redirect(res, req, "/login", http.StatusSeeOther)
	} else if err == nil {
		//check data from database
		//>> mapSessions[myCookie.Value] = checkUser[myCookie.Value]
		//>> mapUsers[mapSessions[myCookie.Value]] = user{checkUser[mapSessions[myCookie.Value]], myCookie.Value}

		id := myCookie.Value
		matchUser, ok := mapUsers[mapSessions[id]] //**temporary match
		errlog.Trace.Println("matchUser:", matchUser)
		if !ok {
			// fmt.Println("Htmlmain.MainMenu - cookie fount with no record match in data")
			http.Redirect(res, req, "/login", http.StatusSeeOther)
			return
		} else {
			// fmt.Println("Htmlmain.MainMenu - cookie fount with matched record in data")
			errlog.Trace.Println("UserName:", matchUser.userName)
			Data.UserName = matchUser.userName
		}
	}

	tpl.ExecuteTemplate(res, "RL_MainMenu.gohtml", Data)
}

//Log Out
func logOut(res http.ResponseWriter, req *http.Request) {

	Data := struct {
		PageName string
		UserName string
	}{PageName: "Log Out", UserName: "bye-bye"}
	Cookie, err := req.Cookie("RecycleLah")
	if err == nil {
		Cookie.MaxAge = -1
		delete(mapUsers, mapSessions[Cookie.Value])
		delete(mapSessions, Cookie.Value)
		http.SetCookie(res, Cookie)
		// fmt.Println("Cookie deleted")
	} else {
		// fmt.Println("No Cookie found and deleted")
		http.Redirect(res, req, "/logIn", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(res, "LogOut.gohtml", Data)
}

//Web Sub Pages func start from below---------------------------------------------------------------------------//

//Book A Car
func pickUp(res http.ResponseWriter, req *http.Request) {

	Data := struct {
		PageName   string
		UserName   string
		CarDisplay []string
	}{PageName: "Jobs List"}

	cookie, err := req.Cookie("RecycleLah")
	Data.UserName = mapSessions[cookie.Value]
	// Data.CarDisplay = vhcs.GetCarNames()
	if _, ok := mapSessions[cookie.Value]; err != nil || !ok {
		http.Redirect(res, req, "/menu", http.StatusSeeOther)
		return
	} else {
		// if req.Method == http.MethodPost {

		// }
		tpl.ExecuteTemplate(res, "RL_Jobs.gohtml", Data)
	}
}

//Current Booking
func viewStatus(res http.ResponseWriter, req *http.Request) {

	Data := struct {
		PageName string
		UserName string
		Bookings [][]string
	}{PageName: "Current Booking"}
	cookie, err := req.Cookie("RecycleLah")

	if _, ok := mapSessions[cookie.Value]; err != nil || !ok {
		http.Redirect(res, req, "/menu", http.StatusSeeOther)
	} else {
		currentUser := mapUsers[mapSessions[cookie.Value]]
		Data.UserName = currentUser.UserName
		// Data.Bookings = bks.ShowAllUserBookings(Data.UserName)
		// fmt.Println(Data.Bookings)
	}
	tpl.ExecuteTemplate(res, "ViewBooking.gohtml", Data)
}

//User Update Detial
// June 08 updated by Sook - commented line
// mapUsers[email] = user{newusername, mapUsers[email].Key}
func userDetailUpdate(res http.ResponseWriter, req *http.Request) {

	Data := struct {
		PageName  string
		UserName  string
		MsgToUser string
	}{PageName: "Updated Password"}

	myCookie, err := req.Cookie("CRA")
	if err != nil {
		http.Redirect(res, req, "/logIn", http.StatusSeeOther)
	} else {
		Data.UserName = mapSessions[myCookie.Value]
	}
	if req.Method == http.MethodPost {
		//get user name and current password
		email := req.FormValue("email")
		password := req.FormValue("oldpassword")
		//get user new password and confirm the new password
		newusername := req.FormValue("newusername")
		newpassword := req.FormValue("newpassword")
		confirmpassword := req.FormValue("confirmpassword")

		pw, mOk := matchPassword[email] //**need to modify
		log.Println(pw, mOk, password == pw)
		if !mOk || !(password == pw) { //**need to modify
			Data.MsgToUser = "The user name and password is not match! "
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			// http.Redirect(res, req, "/changepassword", http.StatusSeeOther)
		} else if newpassword != confirmpassword {
			Data.MsgToUser = "New password and confrim password is not same!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			// http.Redirect(res, req, "/changepassword", http.StatusSeeOther)
		} else if email == "" || password == "" {
			Data.MsgToUser = "Please insert username and password for verification!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		} else {
			//start update DB
			log.Println(email, password, newpassword, confirmpassword, newusername)
			// SOOKMODIFIED mapUsers[email] = user{newusername, mapUsers[email].Key}
			matchPassword[email] = newpassword
			//end update DB
			Data.MsgToUser = "Detail is updated!"
			defer fmt.Fprintf(res, "<h4 class='Application'><a href='/menu'>Main Menu</a></h4><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		}
	}
	tpl.ExecuteTemplate(res, "RL_UpdateUserDetail.gohtml", Data)
}

//Web pull out func ---------------------------------------------------------------------------//
//set cookie on client computer
func setCookie(res http.ResponseWriter, id string) error {
	//name of cookies = "cookie" for 1hr & "RecycleLah" for 2yrs
	co := &http.Cookie{
		Name:     "RecycleLah",
		Value:    id,
		HttpOnly: false,
		Expires:  time.Now().AddDate(2, 0, 0),
	}
	http.SetCookie(res, co)
	// fmt.Println("Htmlmain.setCookie - done with set id = ", id)
	return nil
}

func isEmailValid(e string) bool {
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

const baseURL = "http://localhost:5000/api/v1"

func registration(user nUser) {
	log.Println("registeration")

	log.Println(user)
	jsonValue, _ := json.Marshal(user)
	http.Post(baseURL+"/"+"registration", "application/json", bytes.NewBuffer(jsonValue))

	// response.Header.Set("Content-Type", "application/json")

	// if err != nil {
	// 	fmt.Printf("The HTTP request failed with error %s\n", err)
	// } else {
	// 	ky, _ := ioutil.ReadAll(response.Body)
	// 	response.Body.Close()
	// 	if len(string(ky)) == 36 {
	// 		key = string(ky)
	// 		fmt.Println("Please copy and keep your \"Key\" safe!\nKey: " + string(ky))
	// 		localKeyRecord(key)
	// 		// data := map[string]string{"key": string(ky), "date": time.Now().Format("2006-01-02")}
	// 		// jv, _ := json.Marshal(data)
	// 		// err := ioutil.WriteFile("Key", jv, 0644)
	// 		// if err != nil {
	// 		// 	fmt.Println(err)
	// 		// }
	// 	}
	// 	userInput()
	// }
}
