package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	errlog "github.com/heartziq/recycle-lah/server/utility"
)

type RewardHandler struct {
	Db  *sql.DB
	Tpl *template.Template
	// Error *log.Logger
}

func CreateRewardHandler(db *sql.DB, templatePath string) *RewardHandler {
	newReward := &RewardHandler{Db: db}
	if templatePath != "" {
		newReward.SetTemplate(templatePath)
	}
	return newReward
}

// setting up template, currently not using template
func (p *RewardHandler) SetTemplate(path string) {
	p.Tpl = template.Must(template.ParseGlob(path))
}

// func ServeHTTP calls the respective handlers based on requeset header's content type and method
// it implements a panic recovery function to better handling of client connection
func (p *RewardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer recoverFunc()
	errlog.Trace.Println("r.Host", r.Host)
	errlog.Trace.Println("r.URL", r.URL)
	errlog.Trace.Println("r.RequestURI", r.RequestURI)
	errlog.Trace.Println("r.Context", r.Context())
	errlog.Trace.Println("r.Header", r.Header)
	errlog.Trace.Println("r.Header.Get(Authorization)", r.Header.Get("Authorization"))

	// processing request
	id, reqBody, err := p.userPreProcessRequest(w, r)
	if err != nil {
		errlog.Error.Println(err)
		return
	}
	errlog.Trace.Println("afeter pre-processing request", id, string(reqBody))

	switch r.Method {
	case "GET": // get reward points
		p.methodGetRewards(w, r, id)
		return
	case "PUT": // update reward points
		errlog.Trace.Println("reward PUT")
		p.methodUpdateRewards(w, r, id, reqBody)
		return
	}
	errlog.Error.Println("Uncaterred request")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("405 - Method Not Allowed"))
}

// func methodGetRewards() process request and call db function to retrieve the reward points
// It then send a response back to the client
func (p *RewardHandler) methodGetRewards(w http.ResponseWriter, r *http.Request, id string) {
	var rsp RewardPointsResponse

	// proceed to get reward points from db
	if rewardPoints, err := p.getRewardPoints(p.Db, id); err != nil {
		rsp.Message = appUserError(err)
		w.WriteHeader((http.StatusUnprocessableEntity))
		p.encodeRewardsJsonAndWrite(w, rsp)
		return
	} else {
		rsp.Points = rewardPoints
	}

	rsp.Success = true
	w.WriteHeader(http.StatusAccepted)
	p.encodeRewardsJsonAndWrite(w, rsp)
} // func methodDeleteUser

// func encodeRewardsJsonAndWrite() sets the header of the http.ResponseWriter
// to "application/json".  It then encodes the data into json and writes
// to the http.ResponseWriter.
func (p *RewardHandler) encodeRewardsJsonAndWrite(w http.ResponseWriter, rsp RewardPointsResponse) {
	errlog.Trace.Printf("response=:%+v\n", rsp)
	w.Header().Set("Content-Type", "application/json")
	rsp.Timestamp = int(time.Now().Unix())
	json.NewEncoder(w).Encode(rsp)
}

// func getRewardPoints() retrieves reward points from database table and return
// reward points and error if any
func (p *RewardHandler) getRewardPoints(db *sql.DB, id string) (int, error) {
	var reward struct {
		points int
	}

	err := db.QueryRow("Select reward_points FROM user WHERE id=?", id).Scan(&reward.points)
	if err != nil {
		errlog.Error.Println(err.Error())
		return 0, errSQLStmt
	}
	if err != nil {
		if err == sql.ErrNoRows {
			// record not found
			return 0, userErrNoRecord
		}
		errlog.Error.Println(err.Error())
		return 0, userErrGeneral
	}
	return reward.points, nil
}

// func updateRewardPoints() update the reward points with the new value.  It returns
// the number of rows affected and any error encountered.
func (p *RewardHandler) updateRewardPoints(db *sql.DB, id string, points int) (int, error) {
	results, err := db.Exec("UPDATE user SET reward_points=? WHERE id=?", points, id)
	if err != nil {
		errlog.Error.Println("Error in db.Exec - updateRewardPoints", err)
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

// func methodUpdateRewards processes request and calls database function to perform update
// if there is no error. It will send response back to the client
func (p *RewardHandler) methodUpdateRewards(w http.ResponseWriter, r *http.Request, id string, reqBody []byte) {
	errlog.Trace.Println(id, string(reqBody))
	var reward AdditionalRewardPoints
	var rsp Response

	// unmarshal data into couseInfo data structure
	err := json.Unmarshal(reqBody, &reward)
	if err != nil {
		errlog.Error.Println(err)
		rsp.Message = userErrUnmarshalReqBody.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return
	}
	// validates if details presents

	// proceed to update and prepare response accordingly
	// errlog.Trace.Printf("Going to update points: %s %d\n", id, reward.Points)
	if _, err := p.updateRewardPoints(p.Db, id, reward.Points); err != nil {
		rsp.Message = appUserError(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		p.encodeJsonAndWrite(w, rsp)
		return
	} else { // else no err
		errlog.Trace.Println(">>>>>>201 - User added:", id)
		rsp.Success = true
		w.WriteHeader(http.StatusCreated)
		p.encodeJsonAndWrite(w, rsp)
		return
	} // else no err
} // methodUpdateRewards()

// func encodeJsonAndWrite() sets the header of the http.ResponseWriter
// to "application/json".  It then encodes the data into json and writes
// to the http.ResponseWriter.
func (p *RewardHandler) encodeJsonAndWrite(w http.ResponseWriter, rsp Response) {
	errlog.Trace.Printf("response=:%+v\n", rsp)
	w.Header().Set("Content-Type", "application/json")
	rsp.Timestamp = int(time.Now().Unix())
	json.NewEncoder(w).Encode(rsp)
}

// func userPreProcessRequest() calls functions to extract parameters and request body
// It returns these values when there is no error occurred.
func (p *RewardHandler) userPreProcessRequest(w http.ResponseWriter, r *http.Request) (id string, reqBody []byte, err error) {
	// params
	id, err = p.userPreProcessParam(w, r)
	if err != nil {
		return "", nil, err
	}
	reqBody, err = p.userPreProcessBody(w, r)
	errlog.Info.Println(id, string(reqBody), err)
	return id, reqBody, nil
}

// func userPreProcessParam() checks for the parameters and performs input validation
// It writes to the client when an error is detected.
// It returns the parameter if there is no error.
func (p *RewardHandler) userPreProcessParam(w http.ResponseWriter, r *http.Request) (string, error) {
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
	// if len(id) > 30 {
	// 	rsp.Message = errUserNameLength.Error()
	// 	w.WriteHeader(http.StatusUnprocessableEntity)
	// 	p.encodeJsonAndWrite(w, rsp)
	// 	return "", errUserNameLength
	// }
	// ok = UserNamePattern(id)
	// if !ok {
	// 	rsp.Message = errUserNameFmt.Error()
	// 	w.WriteHeader(http.StatusUnprocessableEntity)
	// 	p.encodeJsonAndWrite(w, rsp)
	// 	return "", errUserNameFmt

	// }
	errlog.Trace.Printf("id supplied=%s\n", id)
	return strings.ToUpper(id), nil
}

// func preProcessBody() reads the request's body and returns
// the body if there is no error.  It writes to the client when an error is detected.
func (p *RewardHandler) userPreProcessBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
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
