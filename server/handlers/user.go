package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
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

// method SetTemplate() to setting up template
func (p *UserHandler) SetTemplate(path string) {
	p.Tpl = template.Must(template.ParseGlob(path))
}

// func CreateUserHandler() sets database handler and template path
// template currently unused
func CreateUserHandler(db *sql.DB, templatePath string) *UserHandler {
	newUser := &UserHandler{Db: db}
	if templatePath != "" {
		newUser.SetTemplate(templatePath)
	}
	return newUser
}

// func ServeHTTP calls the respective handlers based on requeset header's content type and method
// it implements a panic recovery function to better handling of client connection
func (p *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer recoverFunc()
	errlog.Trace.Println("r.Host", r.Host)
	errlog.Trace.Println("r.URL", r.URL)
	errlog.Trace.Println("r.RequestURI", r.RequestURI)
	errlog.Trace.Println("r.Context", r.Context())
	errlog.Trace.Println("r.Header", r.Header)

	id, reqBody, err := p.userPreProcessRequest(w, r)
	if err != nil {
		errlog.Error.Println(err)
		return
	}
	errlog.Trace.Println("afeter pre-processing request", id, string(reqBody))

	if r.Header.Get("Content-type") == "application/json" {
		switch r.Method {
		case "GET": // verify user - authentication
			p.methodVerifyUser(w, r, id, reqBody)
			return
		case "DELETE": // mark a user record as deleted
			p.methodDeleteUser(w, r, id)
			return
		case "PUT": // update user particular
			p.methodPutUser(w, r, id, reqBody)
			return
		case "POST": // create new user
			p.methodPostUser(w, r, id, reqBody)
			return
		}
	} else {
		switch r.Method {
		case "GET": // get user details
			p.methodGetUser(w, r, id)
			return
		}
	}
	errlog.Error.Println("Uncaterred request")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("405 - Method Not Allowed"))
} // func Users()

// func methodPostUser() process request data from the client and call functions to
// insert new record into the database table.  It checks for error handling and will
// update database only when there is no error.
// It returns a response that contains the status of the operation.
func (p *UserHandler) methodPostUser(w http.ResponseWriter, r *http.Request, id string, reqBody []byte) {
	errlog.Trace.Println(id, string(reqBody))
	var newUser NewUser
	var dbData dbNewUser
	var rsp Response

	// unmarshal data into couseInfo data structure
	err := json.Unmarshal(reqBody, &newUser)
	// errlog.Trace.Println("====================newUser", newUser, id, newUser.Password, newUser.Email, newUser.UserName, newUser.Collector)
	if err != nil {
		errlog.Error.Println(err)
		rsp.Message = userErrUnmarshalReqBody.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return
	}

	// validates if details presents
	// errlog.Trace.Println(newUser)
	if strings.TrimSpace(id) == "" {
		errlog.Error.Println("id is blank")
		rsp.Message = userErrMissingAccount.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return
	}
	// validates if password presents
	if strings.TrimSpace(newUser.Password) == "" {
		rsp.Message = errNoPassword.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return
	}
	// validates format for email - currently not checking for email format
	ok := true
	if !ok {
		rsp.Message = errEmailFormat.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return
	}
	// assign details to dbData to be passed to addUser() function
	dbData.id = id
	dbData.password = string(middleware.HashPassword(newUser.Password))
	dbData.email = newUser.Email
	dbData.userName = newUser.UserName
	dbData.collector = newUser.Collector
	// proceed to process creation and prepare response accordingly
	errlog.Trace.Printf("Going to add: %s %s\n", id, newUser.UserName)
	if err := p.addUser(p.Db, dbData); err != nil {
		// if error shows duplicate record, set response message accordingly
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			rsp.Message = errUserNameNotAvail.Error()
			w.WriteHeader(http.StatusConflict)
			p.encodeJsonAndWrite(w, rsp)
			return
		} else { // else not duplicate error
			rsp.Message = appUserError(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			p.encodeJsonAndWrite(w, rsp)
			return
		} // else not diplicate error
	} else { // else no err from processAddUser - successful
		errlog.Trace.Println(">>>>>>201 - User added:", id)
		rsp.Success = true
		w.WriteHeader(http.StatusCreated)
		p.encodeJsonAndWrite(w, rsp)
		return
	} // else no err from processAddUser
} // methodPostUser()

// func methodDeleteUser() calls processMarkUserAsDelete() to perform the request.
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

// func addUser() inserts a new user record into the database.  It returns
// error if error encountered.
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

// func userPreProcessParam() checks for the parameters and performs input validation
// It writes to the client when an error is detected.
// It returns the parameter if there is no error.
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
}

// func userPreProcessBody() reads the request's body and returns
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
	// errlog.Trace.Printf("request body:%v\n", string(reqBody))
	return reqBody, nil
}

// func userPreProcessRequest() calls few function to extract data from the http.Request.
// It returns these values when there is no error occurred.
func (p *UserHandler) userPreProcessRequest(w http.ResponseWriter, r *http.Request) (id string, reqBody []byte, err error) {
	// get params
	id, err = p.userPreProcessParam(w, r)
	if err != nil {
		return "", nil, err
	}
	//
	reqBody, err = p.userPreProcessBody(w, r)
	errlog.Info.Println(id, string(reqBody), err)
	return id, reqBody, nil
}

// func getUserInfo() retrieve user record from database for a given user id
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

// func methodGetUser calls getUserInfo to get database details and return reward points to the client
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

// func encodeUserInfoJsonAndWrite() sets the header of the http.ResponseWriter
// to "application/json".  It then encodes the data into json and writes
// to the http.ResponseWriter.
func (p *UserHandler) encodeUserInfoJsonAndWrite(w http.ResponseWriter, rsp UserInfoResponse) {
	errlog.Trace.Printf("response=:%+v\n", rsp)
	w.Header().Set("Content-Type", "application/json")
	rsp.Timestamp = int(time.Now().Unix())
	json.NewEncoder(w).Encode(rsp)
}

// func methodVerifyUser() unmarshal input, and calls functions to perform user authentication and
// return user details in the response
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

// method getUserSensitiveInfo retrieves password and other details
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
