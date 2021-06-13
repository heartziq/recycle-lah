package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	errlog "github.com/heartziq/recycle-lah/server/utility"
	middleware "github.com/heartziq/recycle-lah/server/utility"
)

type UserHandler struct {
	Db  *sql.DB
	Tpl *template.Template
}

var UserNameExp = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_\\.]*[a-zA-Z0-9]$")

// func UserNamePattern() matches a given string with UserNameExp.
// It returns true if it matches.
// Otherwise, it returns false.
func UserNamePattern(s string) bool {
	matched := UserNameExp.MatchString(s)
	return matched
}

// setting up template
func (p *UserHandler) SetTemplate(path string) {
	p.Tpl = template.Must(template.ParseGlob(path))
}

// func init() calls various functions as well as getting command line arguments
// to setup server and database configuration details.  It then call openDB() to
// establish database connection.
// func init() {
// 	// loadConfig()
// 	// getArgs()
// 	var err error
// 	if db, err = middleware.OpenDB(); err != nil {
// 		errlog.Panic.Panicln(err)
// 	}
// 	// Tpl = template.Must(template.ParseGlob("templates/user/*"))
// }

func CreateUserHandler(db *sql.DB, templatePath string) *UserHandler {
	newUser := &UserHandler{Db: db}
	if templatePath != "" {
		newUser.SetTemplate(templatePath)
	}

	return newUser
}

// func Users(w http.ResponseWriter, r *http.Request) {
func (p *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errlog.Info.Println("r.Host", r.Host)
	errlog.Info.Println("r.URL", r.URL)
	errlog.Info.Println("r.RequestURI", r.RequestURI)
	errlog.Info.Println("r.Context", r.Context())
	errlog.Info.Println("r.Header", r.Header)

	defer recoverFunc()
	id, reqBody, err := p.userPreProcessRequest(w, r)
	if err != nil {
		errlog.Error.Println(err)
		return
	}
	errlog.Trace.Println("afeter pre-processing request", id, string(reqBody))

	if r.Header.Get("Content-type") == "application/json" {
		switch r.Method {
		case "GET":
			p.methodVerifyUser(w, r, id, reqBody)

			// w.WriteHeader(http.StatusOK)
			// w.Write([]byte("to verify user"))
			return
		case "DELETE":
			p.methodDeleteUser(w, r, id)
			return
			// w.WriteHeader(http.StatusOK)
			// w.Write([]byte("DELETE: mark user record as deleted"))
		case "PUT":
			errlog.Trace.Println("user PUT ")
			p.methodPutUser(w, r, id, reqBody)
			return
			// w.WriteHeader(http.StatusOK)
			// w.Write([]byte("POST: create user in package userhandler"))
		case "POST":
			errlog.Trace.Println("user POST ")
			p.methodPostUser(w, r, id, reqBody)
			return
			// w.WriteHeader(http.StatusOK)
			// w.Write([]byte("POST: create user in package userhandler"))
		}
	} else {
		switch r.Method {
		case "GET":
			p.methodGetUser(w, r, id)
			return
			// w.WriteHeader(http.StatusOK)
			// w.Write([]byte("DELETE: mark user record as deleted"))
		}
	}

	errlog.Error.Println("Uncaterred request")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("405 - Method Not Allowed"))
} // func Users()

// func methodPostUser() unmarshals the request body from its parameters list into
// msUser data structure format. It performs input validation on the data received.
// All responses returned to the response writer is in the format of Response data structure.
// When the data passes the validation check, it calls processAddUser() to perform the
// request. Based on the return values from processAddUser(), it prepares the appropriate
// response and calls encodeJsonAndWrite() to encode and send the data in json format.
func (p *UserHandler) methodPostUser(w http.ResponseWriter, r *http.Request, id string, reqBody []byte) {
	errlog.Trace.Println(id, string(reqBody))
	var newUser NewUser
	var dbData dbNewUser
	var rsp Response
	// unmarshal data into couseInfo data structure
	err := json.Unmarshal(reqBody, &newUser)

	errlog.Trace.Println("====================newUser", newUser, id, newUser.Password, newUser.Email, newUser.UserName, newUser.Collector)
	if err != nil {
		errlog.Error.Println(err)
		rsp.Message = userErrUnmarshalReqBody.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return
	}
	// validates if details presents

	errlog.Trace.Println(newUser)
	if strings.TrimSpace(id) == "" {
		errlog.Error.Println("id is blank")
		rsp.Message = userErrMissingAccount.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return
	}
	// validates if title presents
	if strings.TrimSpace(newUser.Password) == "" {
		rsp.Message = errNoPassword.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return
	}
	// validates format for email
	ok := true
	if !ok {
		rsp.Message = errEmailFormat.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return
	}
	dbData.id = id
	dbData.password = string(middleware.HashPassword(newUser.Password))
	dbData.email = newUser.Email
	dbData.userName = newUser.UserName
	dbData.collector = newUser.Collector
	// proceed to process creation and prepare response accordingly
	errlog.Trace.Printf("Going to add: %s %s\n", id, newUser.Email)
	if err := p.addUser(p.Db, dbData); err != nil {
		// if error shows duplicate record, set response message accordingly
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			rsp.Message = errUserNameNotAvail.Error()
			w.WriteHeader(http.StatusConflict)
			p.encodeJsonAndWrite(w, rsp)
			return
		} else { // else no duplicate
			rsp.Message = appUserError(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			p.encodeJsonAndWrite(w, rsp)
			return
		} // else no diplicate
	} else { // else no err from processAddUser
		errlog.Trace.Println(">>>>>>201 - User added:", id)
		rsp.Success = true
		w.WriteHeader(http.StatusCreated)
		p.encodeJsonAndWrite(w, rsp)
		return
	} // else no err from processAddUser
} // methodPostUser()

// func methodDeleteUser() performs input validation on the id received.
// It calls processMarkUserAsDelete() to perform the request.
// Based on the return values from processMarkUserAsDelete(), it prepares the appropriate
// response and calls encodeJsonAndWrite() to encode and send the data in json format.
func (p *UserHandler) methodDeleteUser(w http.ResponseWriter, r *http.Request, id string) {
	var rsp Response
	if strings.TrimSpace(id) == "" {
		errlog.Error.Println("id is blank")
		rsp.Message = userErrMissingAccount.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return
	}

	// proceed to process deletion and prepare response accordingly
	if err := p.processMarkUserAsDelete(p.Db, id); err != nil {
		rsp.Message = appUserError(err)
		w.WriteHeader((http.StatusUnprocessableEntity))
		p.encodeJsonAndWrite(w, rsp)
		return
	}
	rsp.Success = true
	w.WriteHeader(http.StatusAccepted)
	p.encodeJsonAndWrite(w, rsp)
} // func methodDeleteUser

// func addUser() inserts a new record of table courses.  It returns
// the number of rows affected and any error encountered.
func (p *UserHandler) addUser(db *sql.DB, user dbNewUser) error {
	results, err := db.Exec("INSERT INTO user (id, password, email, user_name, is_collector) VALUES (?,?,?,?,?)", user.id, user.password, user.email, user.userName, user.collector)
	if err != nil {
		errlog.Error.Println("Error in db.Exec - insert into user", err)
		// Error 1062: Duplicate entry
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return errUserNameNotAvail
		}
		return errSQLStmt
	}
	if rows, err := results.RowsAffected(); err != nil {
		errlog.Error.Println("Error in insert into courses", err, rows)
		if rows == 0 {
			return userErrNotAdded
		}
		return err
	} else { // no err
		errlog.Info.Println("Number of rows added:", rows)
		return nil
	} // no err
} // func addUser()

// func processMarkUserAsDelete() calls markUserAsDelete() to update the database record.
func (p *UserHandler) processMarkUserAsDelete(db *sql.DB, id string) error {
	i, err := p.markUserAsDelete(db, id)
	if err != nil {
		errlog.Trace.Println(err)
		errlog.Trace.Println("in processMarkUserAsDelete:", appUserError(err))

		return err
	}
	if i == 1 {
		errlog.Info.Printf("mark user id=%s as deleted\n", id)
		return nil
	} else if i == 0 {
		errlog.Info.Println("Record not found, failed to update")
		return userErrNoRecord
	} else { // i > 1
		errlog.Error.Println("Number of courses updated > 1")
		return errMoreThanOne
	}
} // func processMarkUserAsDelete()

// func markUserAsDelete() set user record to be inactive.  It returns
// the number of rows affected and any error encountered.
func (p *UserHandler) markUserAsDelete(db *sql.DB, id string) (int, error) {
	results, err := db.Exec("UPDATE user SET is_deleted=true WHERE id=? and is_deleted=false", id)
	if err != nil {
		errlog.Error.Println("Error in db.Exec - edit user (isDeleted flag)", err)
		return 0, errSQLStmt
	}
	if rows, err := results.RowsAffected(); err != nil {
		errlog.Error.Println("Error in updating user", err)
		return int(rows), err
	} else { // no err
		errlog.Info.Println("Number of rows updated:", rows)
		return int(rows), nil
	} // no err
} // func markUserAsDelete()

// func encodeJsonAndWrite() sets the header of the http.ResponseWriter
// to "application/json".  It then encodes the data into json and writes
// to the http.ResponseWriter.
func (p *UserHandler) encodeJsonAndWrite(w http.ResponseWriter, rsp Response) {
	errlog.Trace.Printf("response=:%+v\n", rsp)
	w.Header().Set("Content-Type", "application/json")
	rsp.Timestamp = int(time.Now().Unix())
	json.NewEncoder(w).Encode(rsp)
}

// func preProcessParam() checks for the course code and performs
// input validation (e.g. length, code format) and sanitization (converts to uppercase).
// It writes to the client when an error is detected.
// It returns course code if there is no error.
func (p *UserHandler) userPreProcessParam(w http.ResponseWriter, r *http.Request) (string, error) {
	rsp := Response{}
	params := mux.Vars(r)
	errlog.Trace.Println(params)

	id, ok := params["id"]
	if !ok {
		rsp.Message = userErrMissingCode.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return "", userErrMissingCode
	}
	// validate input
	if strings.TrimSpace(id) == "" {
		rsp.Message = errNoId.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return "", errNoId
	}
	//  should not check here - should move to individual method
	if len(id) > 30 {
		rsp.Message = errUserNameLength.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return "", errUserNameLength
	}
	ok = UserNamePattern(id)
	if !ok {
		rsp.Message = errUserNameFmt.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return "", errUserNameFmt

	}
	errlog.Trace.Printf("id supplied=%s\n", id)
	return id, nil
	// return strings.ToUpper(id), nil
}

// func preProcessBody() reads the request's body and returns
// the body if there is no error.  It writes to the client when an error is detected.
func (p *UserHandler) userPreProcessBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	rsp := Response{}

	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		rsp.Message = userErrReadReqBody.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return nil, err
	}
	errlog.Trace.Printf("request body:%v\n", string(reqBody))
	return reqBody, nil
}

// func preProcessRequest() calls few function to extract the api key,
// course code and request body from the http.Request.
// It returns these values when there is no error occurred.
func (p *UserHandler) userPreProcessRequest(w http.ResponseWriter, r *http.Request) (id string, reqBody []byte, err error) {
	errlog.Info.Println("r.Host", r.Host)
	errlog.Info.Println("r.URL", r.URL)
	errlog.Info.Println("r.RequestURI", r.RequestURI)
	errlog.Info.Println("r.Context", r.Context())
	errlog.Info.Println("r.Header", r.Header)
	// apikey, err = preProcessQuery(w, r)
	// if err != nil {
	// 	return "", "", nil, err
	// }
	// params
	id, err = p.userPreProcessParam(w, r)
	if err != nil {
		return "", nil, err
	}
	//
	reqBody, err = p.userPreProcessBody(w, r)
	errlog.Info.Println(id, string(reqBody), err)
	return id, reqBody, nil
}

func (p *UserHandler) getUserInfo(db *sql.DB, id string) (UserInfo, error) {
	var user UserInfo

	err := db.QueryRow("Select id, email, user_name, is_deleted, is_collector, reward_points, creation_date, updated_date FROM user WHERE id=?", id).
		Scan(&user.Id, &user.Email, &user.UserName, &user.isDeleted, &user.IsCollector, &user.RewardPoints, &user.CreatedDate, &user.UpdatedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			// record not found
			return user, userErrNoRecord
		}
		errlog.Error.Println(err.Error())
		return user, userErrGeneral
	}
	return user, nil
}

func (p *UserHandler) methodGetUser(w http.ResponseWriter, r *http.Request, id string) {
	var rsp UserInfoResponse
	if rewardPoints, err := p.getUserInfo(p.Db, id); err != nil {
		rsp.Message = appUserError(err)
		w.WriteHeader((http.StatusUnprocessableEntity))
		p.encodeUserInfoJsonAndWrite(w, rsp)
		return
	} else {
		rsp.UserInfo = rewardPoints
	}

	rsp.Success = true
	w.WriteHeader(http.StatusAccepted)
	p.encodeUserInfoJsonAndWrite(w, rsp)
} // func methodDeleteUser

// func encodeRewardsJsonAndWrite() sets the header of the http.ResponseWriter
// to "application/json".  It then encodes the data into json and writes
// to the http.ResponseWriter.
func (p *UserHandler) encodeUserInfoJsonAndWrite(w http.ResponseWriter, rsp UserInfoResponse) {
	errlog.Trace.Printf("response=:%+v\n", rsp)
	w.Header().Set("Content-Type", "application/json")
	rsp.Timestamp = int(time.Now().Unix())
	json.NewEncoder(w).Encode(rsp)
}

// func methodVerifyUser() calls functions to perform user authentication
func (p *UserHandler) methodVerifyUser(w http.ResponseWriter, r *http.Request, id string, reqBody []byte) {
	errlog.Trace.Println(id, string(reqBody))
	var authenticateInfo UserVerification

	var rsp UserInfoResponse
	// unmarshal data into couseInfo data structure
	err := json.Unmarshal(reqBody, &authenticateInfo)

	errlog.Trace.Println("=====authenticateInfo:", id, authenticateInfo.Password)
	if err != nil {
		errlog.Error.Println(err)
		rsp.Message = userErrUnmarshalReqBody.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeUserInfoJsonAndWrite(w, rsp)
		return
	}
	// validates data
	// read user data, get hashed password
	userInfo, err := p.getUserSensitiveInfo(p.Db, id)
	if err != nil {
		errlog.Info.Println(err)
		rsp.Success = false
		rsp.Message = errErrAuthenticate.Error()
		w.WriteHeader((http.StatusUnauthorized))
		p.encodeUserInfoJsonAndWrite(w, rsp)
		return
	}
	// verify Password
	passed := middleware.VerifyPassword([]byte(userInfo.password), authenticateInfo.Password)
	if !passed {
		rsp.Success = false
		rsp.Message = errErrAuthenticate.Error()
		w.WriteHeader((http.StatusUnauthorized))
		p.encodeUserInfoJsonAndWrite(w, rsp)
		return
	}
	rsp.Success = true
	rsp.UserInfo.Id = id
	rsp.UserInfo.Email = userInfo.Email
	rsp.UserInfo.UserName = userInfo.UserName
	rsp.UserInfo.IsCollector = userInfo.IsCollector
	rsp.UserInfo.isDeleted = userInfo.isDeleted
	rsp.UserInfo.RewardPoints = userInfo.RewardPoints
	rsp.UserInfo.CreatedDate = userInfo.CreatedDate
	rsp.UserInfo.UpdatedDate = userInfo.UpdatedDate
	token, err := middleware.GenToken(middleware.KEY, id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf(">>>>>Token: %d %s\n", len(token), token)
	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusAccepted)
	p.encodeUserInfoJsonAndWrite(w, rsp)
} // func methodDeleteUser

func (p *UserHandler) getUserSensitiveInfo(db *sql.DB, id string) (UserInfo, error) {
	var user UserInfo

	err := db.QueryRow("Select id, password, email, user_name, is_deleted, is_collector, reward_points, creation_date, updated_date FROM user WHERE id=? and is_deleted=false", id).
		Scan(&user.Id, &user.password, &user.Email, &user.UserName, &user.isDeleted, &user.IsCollector, &user.RewardPoints, &user.CreatedDate, &user.UpdatedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			// record not found
			return user, userErrNoRecord
		}
		errlog.Error.Println(err.Error())
		return user, userErrGeneral
	}
	errlog.Trace.Println(user)
	return user, nil
}
