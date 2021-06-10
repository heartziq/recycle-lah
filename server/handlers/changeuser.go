package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	errlog "github.com/heartziq/recycle-lah/server/utility"
	middleware "github.com/heartziq/recycle-lah/server/utility"
)

func (p *UserHandler) methodPutUser(w http.ResponseWriter, r *http.Request, id string, reqBody []byte) {
	errlog.Trace.Println(id, string(reqBody))
	var updateUser NewUser
	var dbData dbNewUser
	var rsp Response

	err := json.Unmarshal(reqBody, &updateUser)

	errlog.Trace.Println("methodPutUser", updateUser, id, updateUser.Password, updateUser.Email, updateUser.UserName, updateUser.Collector)
	if err != nil {
		errlog.Error.Println(err)
		rsp.Message = userErrUnmarshalReqBody.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return
	}

	errlog.Trace.Println(updateUser)

	// dbData.id = id
	dbData.password = string(middleware.HashPassword(updateUser.Password))
	// dbData.email = updateUser.Email
	dbData.userName = updateUser.UserName
	// dbData.collector = updateUser.Collector
	// proceed to process deletion and prepare response accordingly
	if _, err := p.updateUserDetail(p.Db, dbData, id); err != nil {
		rsp.Message = appUserError(err)
		w.WriteHeader((http.StatusUnprocessableEntity))
		p.encodeJsonAndWrite(w, rsp)
		return
	}

	rsp.Success = true
	rsp.Message = "User Detial Updated"
	w.WriteHeader(http.StatusAccepted)
	p.encodeJsonAndWrite(w, rsp)
}

// update User Name
// func (p *UserHandler) updateUserDetail(db *sql.DB, user dbNewUser, id string) (int, error) {

// 	var userName, password string
// 	userInfo, err := p.getUserSensitiveInfo(db, id)
// 	if err != nil {
// 		errlog.Info.Println(err)
// 		rsp.Success = false
// 		rsp.Message = errErrAuthenticate.Error()
// 		w.WriteHeader((http.StatusUnauthorized))
// 		p.encodeUserInfoJsonAndWrite(w, rsp)
// 		return 0, err
// 	}
// 	if user.userName == "" {
// 		userName = userInfo.UserName
// 	} else {
// 		userName = user.userName
// 	}
// 	if user.password == "" {
// 		password = userInfo.password
// 	} else {
// 		password = string(middleware.HashPassword(user.password))
// 	}

// 	results, err := db.Exec("UPDATE user SET user_name=?, password=? WHERE id=?", userName, password, id)
// 	if err != nil {
// 		errlog.Error.Println("Error in db.Exec - Update into user", err)
// 		// Error 1062: Duplicate entry
// 		return 0, errSQLStmt
// 	}
// 	if rows, err := results.RowsAffected(); err != nil {
// 		errlog.Error.Println("Error in updating user", err)
// 		return int(rows), err
// 	} else { // no err
// 		errlog.Info.Println("Number of rows added:", rows)
// 		return int(rows), nil
// 	} // no err
// }
func (p *UserHandler) updateUserDetail(db *sql.DB, user dbNewUser, id string) (int, error) {

	var userName, password string
	var rsp UserInfoResponse
	userInfo, err := p.getUserSensitiveInfo(db, id)
	if err != nil {
		errlog.Info.Println(err)
		rsp.Success = false
		rsp.Message = errErrAuthenticate.Error()
		// w.WriteHeader((http.StatusUnauthorized))
		// p.encodeUserInfoJsonAndWrite(w, rsp)
		return 0, err
	}
	if user.userName == "" {
		userName = userInfo.UserName
	} else {
		userName = user.userName
	}
	if user.password == "" {
		password = userInfo.password
	} else {
		password = string(middleware.HashPassword(user.password))
	}
	results, err := db.Exec("UPDATE user SET user_name=?, password=? WHERE id=?", userName, password, id)
	if err != nil {
		errlog.Error.Println("Error in db.Exec - Update into user", err)
		// Error 1062: Duplicate entry
		return 0, errSQLStmt
	}
	if rows, err := results.RowsAffected(); err != nil {
		errlog.Error.Println("Error in updating user", err)
		return int(rows), err
	} else { // no err
		errlog.Info.Println("Number of rows added:", rows)
		return int(rows), nil
	} // no err
}
