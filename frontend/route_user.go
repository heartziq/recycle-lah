package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"frontend/errlog"

	uuid "github.com/satori/go.uuid"
)

//  login() authenticates userid and password
//  call validateUser() to perform the authentication
//  if authenticated, page would be redirected to welcome page
func login(w http.ResponseWriter, r *http.Request) {
	Data := struct {
		PageName  string
		UserName  string
		MsgToUser string
		Dt        map[string]string
	}{PageName: "Log In"}

	Data.Dt = map[string]string{}
	// data = make(map[string]string)
	if r.Method == http.MethodPost {
		validateUser(w, r, Data.Dt)
		_, foundMsg := Data.Dt["Message"]
		if foundMsg {
			Data.MsgToUser = Data.Dt["Error"]
		}
		_, foundErr := Data.Dt["Error"]
		if foundErr {
			Data.MsgToUser = Data.MsgToUser + ", " + Data.Dt["Error"]
		}
		if !foundErr && !foundMsg {
			http.Redirect(w, r, "welcome", http.StatusFound)
			return
		}

	}
	errlog.Trace.Println("test data", Data.Dt)
	// executeTemplate(w, "Login.gohtml", Data)
	// if r.Method == http.MethodPost {
	// 	validateUser(w, r, data)
	// 	_, foundErr := data["Error"]
	// 	_, foundMsg := data["Message"]
	// 	if !foundErr && !foundMsg {
	// 		http.Redirect(w, r, "welcome", http.StatusFound)
	// 		return
	// 	}
	// }
	executeTemplate(w, "login.gohtml", Data)

}

// welcome() displays user menu
// only authenticated user has access to this page
func welcome(w http.ResponseWriter, r *http.Request) {
	data := struct {
		PageName  string
		UserName  string
		Since     string
		Token     string
		Collector string
	}{PageName: "Welcome"}

	user, err := getSession(r)
	if err != nil {
		http.Redirect(w, r, "index.gohtml", http.StatusFound)
		return
	}
	data.UserName = user.userName
	data.Since = string(time.Unix(user.sessionCreatedTime, 0).String()[0:19])
	data.Token = user.token
	if user.isCollector {
		data.Collector = "Y"
	}
	executeTemplate(w, "welcome.gohtml", data)
}

// logout() deletes session details from the session map
// redirect to index page
// to use SinYaw's logOut
func logout(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("in logout")
	clearSession(w, r)
	http.Redirect(w, r, "/", http.StatusFound)

}

// func validateUser() authenticates user
// it calls verifyUser() to send request/process response from api server
// if successfully authenticated, it creates cookie and sets session and user details in map
func validateUser(w http.ResponseWriter, r *http.Request, data map[string]string) {
	// process submitted form
	// data = make(map[string]interface{})
	if r.Method != http.MethodPost {
		return
	}

	//  get data from form
	err := r.ParseForm()
	if err != nil {
		errlog.Error.Println("err in ParseForm", err)
		data["Error"] = "processing error"
		return
	}

	userId := r.FormValue("userid")
	password := r.FormValue("password")
	data["UserName"] = userId
	err = checkInputUserName(userId)
	// var msg []string
	if err != nil {

		// msg = append(msg, err.Error())
		data["Message"] = err.Error()
		errlog.Trace.Println("validateUser", data)
		return
	}
	var reqData UserVerification
	reqData.Password = password
	apiResp, err := verifyUser(userId, reqData)
	if err != nil {

		// msg = append(msg, err.Error())
		data["Message"] = err.Error()
		errlog.Trace.Println("validateUser", data)
		return
	}
	//  set cookie and create session - set cookie, add db session record
	errlog.Trace.Println("apiResp=", apiResp)

	// set session information
	// create cookie
	id := uuid.NewV4()
	cookie := http.Cookie{
		Name:     "RecycleLah",
		Value:    id.String(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	errlog.Trace.Printf("cookie set: %+v\n", cookie)

	// add to map - session and user
	mapSessions[cookie.Value] = userId

	var currentUser user
	currentUser.userId = userId
	currentUser.sessionCreatedTime = time.Now().UnixNano() / int64(time.Second)
	currentUser.isCollector = apiResp.IsCollector
	currentUser.email = apiResp.Email
	currentUser.userName = apiResp.UserName
	currentUser.token = apiResp.token
	mapUsers[userId] = currentUser
	// errlog.Trace.Println("    !!!!!!       currentUser=", currentUser)
	return
}

// func verifyUser() sends request to api server and processes response
// it returns user information and error if any
func verifyUser(id string, reqData UserVerification) (UserInfoResponse, error) {
	var rsp UserInfoResponse
	rsp.Id = id
	errlog.Trace.Println("verifyUser: ", id, reqData)
	jsonValue, err := json.Marshal(reqData)
	if err != nil {
		errlog.Error.Println("error in marshal", err)
		return rsp, err
	}

	// response, err := client.Post(config.BaseURL+"/"+code+"?key="+config.APIKey,
	// "application/json", bytes.NewBuffer(jsonValue))
	url := "http://localhost:5000/api/v1/users/" + id + "?key=secretkey"
	apiReq, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonValue))
	apiReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(apiReq)

	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)
		return rsp, err
	}
	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	token := response.Header.Get("Authorization")
	errlog.Trace.Println("Authorization token from header=", token)
	errlog.Trace.Println("response.Header", response.Header)
	if err != nil {
		errlog.Error.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())
		return rsp, err
	}
	errlog.Trace.Printf("response status code:%+v\nstring(data):%+v\n", response.StatusCode, string(data))

	if err := json.Unmarshal(data, &rsp); err != nil {
		errlog.Error.Println("unmarshal error", err)
		return rsp, err
	} else {
		errlog.Info.Println("response body (unmarshal)=", rsp)
		if rsp.Success {
			rsp.token = token
			return rsp, nil
		}
		return rsp, errors.New(rsp.Message)

	}
}

// func validateNewUser() validates user input, call addUser() to send request to api server
// it returns status to the caller function in the map
func validateNewUser(w http.ResponseWriter, r *http.Request, data map[string]interface{}) error {
	if r.Method != http.MethodPost {
		errlog.Error.Println("Wrong method in validateNewUser")
		return errors.New("Wrong method")
	}
	if err := r.ParseForm(); err != nil {
		errlog.Error.Println("err in ParseForm", err)
		data["Error"] = "processing error"
		return err
	}
	//  get data from form
	userId := r.FormValue("userid")
	userName := r.FormValue("username")
	email := r.FormValue("email")
	password1 := r.FormValue("pwd1")
	password2 := r.FormValue("pwd2")

	// validate input
	if err := checkInputUserName(userId); err != nil {
		data["Error"] = err.Error()
		errlog.Trace.Println("validateNewUser()", data)
		return err
	}
	if userName == "" {
		userName = userId
	}
	data["UserName"] = userName
	if err := checkInputNewPassword(password1); err != nil {
		data["Error"] = err.Error()
		errlog.Trace.Println("validateNewUser() - password check", data)
		return err
	}
	matched := confirmPassword(password1, password2)
	if !matched {
		data["Error"] = "Passwords not matching"
		errlog.Trace.Println("validateNewUser() - passwords not matched", data)
		return errors.New("Passwords not matching")
	}

	// prepares data required by addUser
	var newUser NewUser
	newUser.Password = password1
	newUser.Email = email
	// if email == "" {
	// 	newUser.Email = userId + "@gmail.com"
	// } else {
	// 	newUser.Email = email
	// }

	newUser.UserName = userName
	newUser.Collector = false
	err := addUser(userId, newUser)
	if err != nil {
		data["Error"] = "Failed to add user"
		errlog.Trace.Println("validateNewUser() - failed to add user", data)
		return errors.New("Failed to add user")
	}
	fmt.Println("in validate before return")
	return nil
}

// func addUser() sends request to api server to create new account
// it interprets the response from server and send error if creation fails
func addUser(code string, newUser NewUser) error {
	errlog.Trace.Println("addUser: ", code, newUser)
	jsonValue, err := json.Marshal(newUser)
	if err != nil {
		errlog.Error.Println("error in marshal", err)
		return err
	}

	// set usr and send requests
	url := "http://localhost:5000/api/v1/users/" + code + "?key=secretkey"
	client := &http.Client{}

	response, err := client.Post(url,
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)
		return err
	}

	// processes response from server
	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		errlog.Error.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())
		return err
	}
	errlog.Trace.Printf("response status code:%+v\nstring(data):%+v\n", response.StatusCode, string(data))

	var rsp Response
	if err := json.Unmarshal(data, &rsp); err != nil {
		errlog.Error.Println("unmarshal error", err)
		return err
	} else {
		errlog.Info.Println("response body (unmarshal)=", rsp)
		if rsp.Success {
			return nil
		}
		return errors.New(rsp.Message)

	}
}

// func message() displays message from flash cooke
func message(w http.ResponseWriter, r *http.Request) {
	d := struct {
		Message string
	}{}
	var err error
	d.Message, err = getFlashCookie(w, r)
	if err != nil {
		d.Message = "Message unavailable"
	}
	executeTemplate(w, "message.gohtml", d)
}

// func newUser() creates new user account
// it gets new user particulars, perform validation and send request to api server
// it interprets the response from server and display the correspondng message
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

		// perform input validation
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

		// to check response
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

//Web Sub Pages func start from below---------------------------------------------------------------------------//

//User Update Detial
// June 08 updated by Sook - commented line
// mapUsers[email] = user{newusername, mapUsers[email].Key}
func userDetailUpdate(res http.ResponseWriter, req *http.Request) {
	errlog.Trace.Println("userDetailUpdate")
	Data := struct {
		PageName  string
		UserName  string
		MsgToUser string
		Token     string
		Collector string
	}{PageName: "Updated Password"}

	user, err := getSession(req)
	if err != nil {
		http.Redirect(res, req, "index.gohtml", http.StatusFound)
		return
	}
	Data.UserName = user.userName
	Data.Token = user.token
	if user.isCollector {
		Data.Collector = "Y"
	}

	if req.Method == http.MethodPost {
		//get user name and current password
		userId := req.FormValue("userid")
		//get user new password and confirm the new password
		newusername := req.FormValue("newusername")
		newpassword := req.FormValue("newpassword")
		confirmpassword := req.FormValue("confirmpassword")

		errlog.Trace.Println("Ken In Value=", userId, newusername, newpassword, confirmpassword)

		if userId != user.userId {
			Data.MsgToUser = "Please log in before do updated detail!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			// http.Redirect(res, req, "/changepassword", http.StatusSeeOther)
		} else if newpassword != confirmpassword {
			Data.MsgToUser = "New password and confrim password is not same!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			// http.Redirect(res, req, "/changepassword", http.StatusSeeOther)
		} else if userId == "" {
			Data.MsgToUser = "Please insert user id and password for verification!"
			defer fmt.Fprintf(res, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		} else {
			//start update DB
			var dataUpdate NewUser
			dataUpdate.UserName = newusername
			dataUpdate.Password = newpassword

			id := userId
			// log.Println(email, password, newpassword, confirmpassword, newusername)
			changeUser(dataUpdate, id)
			// SOOKMODIFIED mapUsers[email] = user{newusername, mapUsers[email].Key}
			// matchPassword[email] = newpassword
			//end update DB
			Data.MsgToUser = "Detail is updated!"
			defer fmt.Fprintf(res, "<h4 class='Application'><a href='/'>Main Menu</a></h4><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
		}
	}
	tpl.ExecuteTemplate(res, "RL_UpdateUserDetail.gohtml", Data)
}

//Web pull out func ---------------------------------------------------------------------------//

// func isEmailValid() performs input validation for email
func isEmailValid(e string) bool {
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

// func changeUser() gets user particular changes and send request to api server
func changeUser(user NewUser, id string) error {
	errlog.Trace.Println("checgeUser: ")
	jsonValue, err := json.Marshal(user)
	if err != nil {
		errlog.Error.Println("error in marshal", err)
		return err
	}

	url := "http://localhost:5000/api/v1/users/" + id + "?key=secretkey"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, url,
		bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	errlog.Trace.Printf("put REQUEST=%v\n", request)
	response, err := client.Do(request)
	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)
		return err
	}
	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		errlog.Error.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())
		return err
	}
	errlog.Trace.Printf("response status code:%+v\nstring(data):%+v\n", response.StatusCode, string(data))
	var rsp Response
	if err := json.Unmarshal(data, &rsp); err != nil {
		errlog.Error.Println("unmarshal error", err)
		return err
	} else {
		errlog.Info.Println("response body (unmarshal)=", rsp)
		if rsp.Success {
			return nil
		}
		return errors.New(rsp.Message)
	}
}
