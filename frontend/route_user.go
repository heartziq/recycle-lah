package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"frontend/errlog"

	uuid "github.com/satori/go.uuid"
)

//  login() authenticates userid and password
//  call validateUser() to perform the authentication
//  if authenticated, page would be redirected to staffmenu
func login(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data = make(map[string]interface{})
	if r.Method == http.MethodPost {
		validateUser(w, r, data)
		_, foundErr := data["Error"]
		_, foundMsg := data["Message"]
		if !foundErr && !foundMsg {
			http.Redirect(w, r, "welcome", http.StatusFound)
			return
		}
	}
	executeTemplate(w, "login_sook.gohtml", data)
}

func collectorLogin(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data = make(map[string]interface{})
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
func welcome(w http.ResponseWriter, r *http.Request) {
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
	executeTemplate(w, "welcome.gohtml", data)
}

// welcome() displays login time
// only authenticated user has access to this page
func collectorWelcome(w http.ResponseWriter, r *http.Request) {
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

// logout() deletes session details from the session map
// redirect to index page
// to use SinYaw's logOut
func logout(w http.ResponseWriter, r *http.Request) {
	// p("in logout")
	clearSession(w, r)
	http.Redirect(w, r, "/", http.StatusFound)

}

//  validateUser() is called from login()
//  set session cookie if not found
//  authenticate user and add session id and user access to session map
//  route to staffmenu when verified
func validateUser(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	// process submitted form
	// data = make(map[string]interface{})
	if r.Method != http.MethodPost {
		return
	}
	err := r.ParseForm()
	if err != nil {
		errlog.Error.Println("err in ParseForm", err)
		data["Error"] = "processing error"
		return
	}
	//  get data from form
	userId := r.FormValue("userid")
	password := r.FormValue("pwd")
	data["UserName"] = userId
	err = checkInputUserName(userId)
	var msg []string
	if err != nil {

		msg = append(msg, err.Error())
		data["Message"] = msg
		errlog.Trace.Println("validateUser", data)
		return
	}
	var reqData UserVerification
	reqData.Password = password
	apiResp, err := verifyUser(userId, reqData)
	if err != nil {

		msg = append(msg, err.Error())
		data["Message"] = msg
		errlog.Trace.Println("validateUser", data)
		return
	}
	//  set cookie and create session - set cookie, add db session record
	errlog.Trace.Println("apiResp=", apiResp)
	errlog.Trace.Println("end of validateUser")
	// ok := createSession(w, userName)
	// if !ok {
	// 	data["Error"] = "Login unsuceesful"
	// }
	id := uuid.NewV4()
	cookie := http.Cookie{
		Name:     "RecycleLah",
		Value:    id.String(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
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
	return
}

// func validateUser1(w http.ResponseWriter, r *http.Request, collector bool, data map[string]interface{}) {
// 	// process submitted form
// 	// data = make(map[string]interface{})
// 	if r.Method != http.MethodPost {
// 		return
// 	}
// 	err := r.ParseForm()
// 	if err != nil {
// 		errlog.Error.Println("err in ParseForm", err)
// 		data["Error"] = "processing error"
// 		return
// 	}
// 	//  get data from form
// 	userName := r.FormValue("username")
// 	password := r.FormValue("pwd")
// 	data["UserName"] = userName
// 	err = checkInputUserName(userName)
// 	var msg []string
// 	if err != nil {

// 		msg = append(msg, err.Error())
// 		data["Message"] = msg
// 		errlog.Trace.Println("validateUser", data)
// 		return
// 	}
// 	var reqData UserVerification
// 	reqData.Password = password
// 	apiResp, err := verifyUser1(userName, collector, reqData)
// 	if err != nil {

// 		msg = append(msg, err.Error())
// 		data["Message"] = msg
// 		errlog.Trace.Println("validateUser", data)
// 		return
// 	}
// 	//  set cookie and create session - set cookie, add db session record
// 	errlog.Trace.Println("apiResp=", apiResp)
// 	errlog.Trace.Println("end of validateUser")
// 	// ok := createSession(w, userName)
// 	// if !ok {
// 	// 	data["Error"] = "Login unsuceesful"
// 	// }
// 	id := uuid.NewV4()
// 	cookie := http.Cookie{
// 		Name:     "RecycleLah",
// 		Value:    id.String(),
// 		HttpOnly: true,
// 	}
// 	http.SetCookie(w, &cookie)
// 	errlog.Trace.Printf("cookie set: %+v\n", cookie)
// 	var userSession Session
// 	userSession.userId = userName
// 	userSession.uuid = cookie.Value
// 	userSession.updatedDate = apiResp.UpdatedDate
// 	// userSession.sessionCreatedTime = int(time.Now().Add(time.Minute * 2).Unix())
// 	userSession.sessionCreatedTime = time.Now().UnixNano() / int64(time.Second)
// 	// i, err := strconv.ParseInt(userSession.sessionCreatedTime, 10, 64)
// 	// fmt.Println(i)
// 	// tm := time.Unix(userSession.sessionCreatedTime, 0)
// 	userSession.isCollector = apiResp.IsCollector
// 	userSession.email = apiResp.Email
// 	userSession.token = apiResp.token
// 	mapSession[cookie.Value] = userSession
// 	return
// } // validateUser1

// curl -H "Content-Type: application/json" -X GET http://localhost:5000/api/v1/users/USER4567?key=secretkey -d {\"password\":\"password\"}
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

func verifyUser1(id string, collector bool, reqData UserVerification) (UserInfoResponse, error) {
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

func createAccount(w http.ResponseWriter, r *http.Request) {
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

	if err := checkInputUserName(userId); err != nil {
		data["Error"] = err.Error()
		errlog.Trace.Println("validateNewUser()", data)
		return err
	}
	if userName == "" {
		userName = userId
	}
	data["UserName"] = userName
	// //  check if username already taken
	// avail := user.UserNameAvailable(db, userName)
	// if !avail {
	// 	data["Error"] = "username is not available"
	// 	errlog.Trace.Println("validateNewUser() - not avail", data)
	// 	return errors.New("username is not available")
	// }
	// err = checkInputNewPassword(password1)
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
	// added := user.AddUser(db, userName, password1)
	// if !added {
	// 	data["Error"] = "Failed to add user"
	// 	errlog.Trace.Println("validateNewUser() - failed to add user", data)
	// 	return errors.New("Failed to add user")
	// }
	var newUser NewUser
	newUser.Password = password1
	if email == "" {
		newUser.Email = userId + "@gmail.com"
	} else {
		newUser.Email = email
	}

	newUser.UserName = userName
	newUser.Collector = false
	err := addUser(userId, newUser)
	if err != nil {
		data["Error"] = "Failed to add user"
		errlog.Trace.Println("validateNewUser() - failed to add user", data)
		return errors.New("Failed to add user")
	}
	// recyclePost() //- ok error with TLS
	// recyclePostUser() - ok
	// recyclePostUserData(newUser) - ok
	// recyclePostUserDataReal(userName, newUser) - ok with no TLS
	// recyclePostUserDataReal(userName, newUser) //- error with 400 with TLS
	fmt.Println("in validate before return")
	return nil
}

func addUser(code string, newUser NewUser) error {
	errlog.Trace.Println("addUser: ", code, newUser)
	jsonValue, err := json.Marshal(newUser)
	if err != nil {
		errlog.Error.Println("error in marshal", err)
		return err
	}

	// response, err := client.Post(config.BaseURL+"/"+code+"?key="+config.APIKey,
	// "application/json", bytes.NewBuffer(jsonValue))
	url := "http://localhost:5000/api/v1/users/" + code + "?key=secretkey"
	client := &http.Client{}

	response, err := client.Post(url,
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

func testaddUser() error {
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

// ok
func recycle() {
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

func recyclePost() {
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

func recyclePostUser() {
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

func recyclePostUserData(user NewUser) {
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

func recyclePostUserDataReal(id string, user NewUser) {
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

// message() displays message from flash cooke
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
