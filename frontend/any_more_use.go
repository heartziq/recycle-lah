package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"frontend/errlog"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

//Log Out
func XlogOut(res http.ResponseWriter, req *http.Request) {

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

//log in
// sooK modified
func XlogIn(res http.ResponseWriter, req *http.Request) {

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

//set cookie on client computer
func XsetCookie(res http.ResponseWriter, id string) error {
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

//PickUp
func XpickUp(res http.ResponseWriter, req *http.Request) {

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
func XviewStatus(res http.ResponseWriter, req *http.Request) {

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

func XcollectorLogin(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	data = make(map[string]string)
	if r.Method == http.MethodPost {
		validateUser(w, r, data)
		_, foundErr := data["Error"]
		_, foundMsg := data["Message"]
		if !foundErr && !foundMsg {
			http.Redirect(w, r, "collector_welcome", http.StatusFound)
			return
		}
	}
	executeTemplate(w, "collector_login.gohtml", data)
}

// welcome() displays login time
// only authenticated user has access to this page
func XcollectorWelcome(w http.ResponseWriter, r *http.Request) {
	data := struct {
		UserName string
		Since    string
		Token    string
	}{}
	user, err := getSession(r)
	if err != nil {
		http.Redirect(w, r, "index.gohtml", http.StatusFound)
		return
	}

	data.UserName = user.userName
	data.Since = string(time.Unix(user.sessionCreatedTime, 0).String()[0:19])
	data.Token = user.token
	executeTemplate(w, "collector_welcome.gohtml", data)
}

func XverifyUser1(id string, collector bool, reqData UserVerification) (UserInfoResponse, error) {
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
} // verifyUser1

func XcreateAccount(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data = make(map[string]interface{})
	if r.Method == http.MethodPost {
		if err := validateNewUser(w, r, data); err == nil {
			http.Redirect(w, r, "/signupsuccess", http.StatusFound)
			return
		}
	}
	executeTemplate(w, "signup.gohtml", data)
}

func XtestaddUser() error {
	code := "new88888"
	newUser := NewUser{"password", "email", "username", false}
	errlog.Trace.Println("addUser: ", code, newUser)

	fmt.Println("addUser: ", code, newUser)
	jsonValue, err := json.Marshal(newUser)
	if err != nil {
		errlog.Error.Println("error in marshal", err)
		return err
	}

	client := &http.Client{}
	// response, err := client.Post(config.BaseURL+"/"+code+"?key="+config.APIKey,
	// "application/json", bytes.NewBuffer(jsonValue))

	response, err := client.Post(config.BaseURL+"/"+code+"?key=secretkey",
		"application/json", bytes.NewBuffer(jsonValue))

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

func Xrecycle() {
	url := "http://localhost:5000/api/v1/recycle"
	errlog.Info.Println("client url=" + url)

	client := &http.Client{}
	response, err := client.Get(url)
	fmt.Println("resonse", response.Body, "\nheader=", response.Header)
	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)

	}

	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())

	}

	fmt.Printf("response status code:%+v\nstring(data):%+v\n", response.StatusCode, string(data))

}

func XrecyclePost() {
	url := "http://localhost:5000/api/v1/pickups/3333?key=secretkey"
	errlog.Info.Println("client url=" + url)

	client := &http.Client{}
	response, err := client.Post(url, "application/json", nil)
	fmt.Println("resonse", response.Body, "\nheader=", response.Header)
	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)

	}

	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())

	}

	fmt.Printf("response status code:%+v\nstring(data):%+v\n", response.StatusCode, string(data))

}

func XrecyclePostUser() {
	url := "http://localhost:5000/api/v1/users/3333?key=secretkey"
	errlog.Info.Println("client url=" + url)

	client := &http.Client{}
	response, err := client.Post(url, "application/json", nil)
	fmt.Println("resonse", response.Body, "\nheader=", response.Header)
	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)

	}

	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())

	}

	fmt.Printf("response status code:%+v\nstring(data):%+v\n", response.StatusCode, string(data))

}

func XrecyclePostUserData(user NewUser) {
	jsonValue, err := json.Marshal(user)
	if err != nil {
		errlog.Error.Println("error in marshal", err)
	}
	url := "http://localhost:5000/api/v1/users/3333?key=secretkey"
	errlog.Info.Println("client url=" + url)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	fmt.Println("resonse", response.Body, "\nheader=", response.Header)
	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)

	}

	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())

	}

	fmt.Printf("response status code:%+v\nstring(data):%+v\n", response.StatusCode, string(data))

}

func XrecyclePostUserDataReal(id string, user NewUser) {
	jsonValue, err := json.Marshal(user)
	if err != nil {
		errlog.Error.Println("error in marshal", err)
	}
	// url := "http://localhost:5000/api/v1/users/3333?key=secretkey"
	url := "http://localhost:5000/api/v1/users/" + id + "?key=secretkey"
	errlog.Info.Println("client url=" + url)

	client := &http.Client{}
	response, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	fmt.Println("resonse", response.Body, "\nheader=", response.Header)
	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)

	}

	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())

	}

	fmt.Printf("response status code:%+v\nstring(data):%+v\n", response.StatusCode, string(data))

}
