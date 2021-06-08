package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	uuid "github.com/satori/go.uuid"
)

type user struct {
	UserName string
	// Password string
	Key string
}

var (
	tpl           *template.Template
	emailRegex    = regexp.MustCompile("^[\\w!#$%&'*+/=?`{|}~^-]+(?:\\.[\\w!#$%&'*+/=?`{|}~^-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,6}$") // regular expression
	mapUsers      = map[string]user{"r@l.com": user{"recycle", "278d0e77-76c2-4447-bbfb-6fb032f57414"}}                                 //**temporary use data
	mapSessions   = map[string]string{"278d0e77-76c2-4447-bbfb-6fb032f57414": "r@l.com"}
	matchPassword = map[string]string{"r@l.com": "password"} //**need to get from Database
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	log.SetFlags(log.Lshortfile)
}

func main() {

	//Main Pages
	http.HandleFunc("/", mainMenu)
	http.HandleFunc("/menu", mainMenu)
	http.HandleFunc("/newuser", newUser)
	http.HandleFunc("/logout", logOut)
	http.HandleFunc("/login", logIn)
	//Sub Pages
	http.HandleFunc("/userupdate", userDetailUpdate)
	http.HandleFunc("/pickup", pickUp)
	http.HandleFunc("/viewstatus", viewStatus)
	//Load Files
	http.Handle("/Pictures/", http.StripPrefix("/Pictures", http.FileServer(http.Dir("Pictures"))))
	http.Handle("/Stuff/", http.StripPrefix("/Stuff", http.FileServer(http.Dir("Stuff"))))
	//run server
	log.Println(http.ListenAndServe(":5221", nil))

}

// uuid "github.com/satori/go.uuid"

type nUser struct {
	UserName  string
	Email     string
	Password  string
	uuid      string
	Collector bool
}

//Web Main Pages func from below---------------------------------------------------------------------------//
//New User
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
		username := req.FormValue("username")
		password := req.FormValue("password")
		confirmpassword := req.FormValue("confirmpassword")
		email := req.FormValue("email")
		C_Bool := req.FormValue("collector")
		var Collector bool
		if C_Bool == "true" {
			Collector = true
		} else {
			Collector = false
		}
		_, nameFound := matchPassword[email] //**need to get from database
		UFok := isEmailValid(email)
		if !nameFound && password == confirmpassword && UFok {
			// fmt.Println("Htmlmain.newUser - no same name in existing data")
			id := uuid.NewV4().String()
			setCookie(res, id)
			ToDB := nUser{username, email, password, id, Collector}
			log.Println(ToDB, confirmpassword)
			registeration(ToDB)             // send data to server
			matchPassword[email] = password //**temporary use
			mapSessions[id] = email
			mapUsers[email] = user{username, id}
			Data.MsgToUser = "New User Registration Done! You may process to log in."
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		} else if !UFok {
			// fmt.Println("Htmlmain.newUser - email format not correct")
			Data.MsgToUser = "Please enter correct email!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		} else if nameFound {
			// fmt.Println("Htmlmain.newUser - name in existing data")
			fmt.Scanf(Data.MsgToUser, "Please use other email! '%v' has been taken!", email)
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		} else if password != confirmpassword {
			// fmt.Println("Htmlmain.newUser - confirm password not match")
			Data.MsgToUser = "Confirm Password is not same!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		}
	}

}

//log in
func logIn(res http.ResponseWriter, req *http.Request) {

	Data := struct {
		PageName  string
		UserName  string
		MsgToUser string
	}{PageName: "Log In"}
	myCookie, err := req.Cookie("RecycleLah")
	if err == nil {
		Data.UserName = mapUsers[mapSessions[myCookie.Value]].UserName
	} else {
		Data.UserName = ""
	}

	if req.Method == http.MethodPost {
		email := req.FormValue("email")
		password := req.FormValue("password")
		log.Println(email, password)
		if email == "" {
			Data.MsgToUser = "No value found is fould!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		} else {
			pw, mOk := matchPassword[email] //**need to get from database
			if pw == password && mOk {      //**need to modify also
				log.Println("ok", email, password)
				currentUser := user{mapUsers[email].UserName, mapUsers[email].Key} //**need to get from database
				mapUsers[email] = currentUser
				mapSessions[currentUser.Key] = email
				// log.Println("Htmlmain.logIn - currentUser", currentUser)
				setCookie(res, currentUser.Key)
				http.Redirect(res, req, "/menu", http.StatusSeeOther)
			} else {
				Data.MsgToUser = "The User Name or Password is incorrect!"
				defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			}
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
		if !ok {
			// fmt.Println("Htmlmain.MainMenu - cookie fount with no record match in data")
			http.Redirect(res, req, "/login", http.StatusSeeOther)
		} else {
			// fmt.Println("Htmlmain.MainMenu - cookie fount with matched record in data")
			Data.UserName = matchUser.UserName

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

func pickUp(res http.ResponseWriter, req *http.Request) {

	type Job struct {
		JobID       string
		Address     string
		Postcode    string
		Description string
		WeidthRange string
	}

	Data := struct {
		PageName string
		UserName string
		Jobs     []Job
	}{PageName: "Jobs List"}
	//get the jobs list from API Server for bellow data.
	Job1 := Job{"001", "70 Woodlands Avenue 7", "738344", "2 Bags of Plastic Bottles", "1) up to 1kg"}
	var Jobs []Job
	for i := 5; i > 0; i-- {
		Jobs = append(Jobs, Job1)
	}
	Data.Jobs = Jobs
	log.Println(Data.Jobs)
	cookie, err := req.Cookie("RecycleLah")
	Data.UserName = mapUsers[mapSessions[cookie.Value]].UserName
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
	}
	tpl.ExecuteTemplate(res, "ViewBooking.gohtml", Data)
}

//User Update Detial
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
			mapUsers[email] = user{newusername, mapUsers[email].Key}
			matchPassword[email] = newpassword
			//end update DB
			Data.MsgToUser = "Detail is updated!"
			defer fmt.Fprintf(res, "<h4 class='Application'><a href='/menu'>Main Menu</a></h4><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		}
	}
	tpl.ExecuteTemplate(res, "RL_UpdateUserDetail.gohtml", Data)
}

//Pull out func ---------------------------------------------------------------------------//
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

func registeration(user nUser) {
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
