package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"frontend/errlog"
	"io/ioutil"
	"net/http"
	"strings"
)

type nUser struct {
	UserName  string
	Email     string
	Password  string
	uuid      string
	Collector bool
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
