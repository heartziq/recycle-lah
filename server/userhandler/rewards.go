package userhandler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	errlog "github.com/heartziq/recycle-lah/server/utility"
)

func Rewards(w http.ResponseWriter, r *http.Request) {
	p("Sook in Rewards()")
	defer recoverFunc()
	id, reqBody, err := userPreProcessRequest(w, r)
	if err != nil {
		errlog.Error.Println(err)
		return
	}
	errlog.Trace.Println("afeter pre-processing request", id, string(reqBody))

	switch r.Method {
	case "GET":
		methodGetRewards(w, r, id)
		return
		// w.WriteHeader(http.StatusOK)
		// w.Write([]byte("DELETE: mark user record as deleted"))

	case "PUT":
		p("Sook in Rewards() - case PUT ")
		methodUpdateRewards(w, r, id, reqBody)
		return
		// w.WriteHeader(http.StatusOK)
		// w.Write([]byte("POST: create user in package userhandler"))
	}

	errlog.Error.Println("Uncaterred request")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("405 - Method Not Allowed"))
} // func Users()

func methodGetRewards(w http.ResponseWriter, r *http.Request, id string) {
	var rsp RewardPointsResponse

	// proceed to get reward points from db
	if rewardPoints, err := getRewardPoints(db, id); err != nil {
		rsp.Message = appUserError(err)
		w.WriteHeader((http.StatusUnprocessableEntity))
		encodeRewardsJsonAndWrite(w, rsp)
		return
	} else {
		rsp.Points = rewardPoints
	}

	rsp.Success = true
	w.WriteHeader(http.StatusAccepted)
	encodeRewardsJsonAndWrite(w, rsp)
} // func methodDeleteUser

// func encodeRewardsJsonAndWrite() sets the header of the http.ResponseWriter
// to "application/json".  It then encodes the data into json and writes
// to the http.ResponseWriter.
func encodeRewardsJsonAndWrite(w http.ResponseWriter, rsp RewardPointsResponse) {
	errlog.Trace.Printf("response=:%+v\n", rsp)
	w.Header().Set("Content-Type", "application/json")
	rsp.Timestamp = int(time.Now().Unix())
	json.NewEncoder(w).Encode(rsp)
}

func getRewardPoints(db *sql.DB, id string) (int, error) {
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
func updateRewardPoints(db *sql.DB, id string, points int) (int, error) {
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

func methodUpdateRewards(w http.ResponseWriter, r *http.Request, id string, reqBody []byte) {
	errlog.Trace.Println(id, string(reqBody))
	var reward AdditionalRewardPoints
	var rsp Response
	// unmarshal data into couseInfo data structure
	err := json.Unmarshal(reqBody, &reward)

	errlog.Trace.Println("====================additional reward", id, reward)
	if err != nil {
		errlog.Error.Println(err)
		rsp.Message = userErrUnmarshalReqBody.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		encodeJsonAndWrite(w, rsp)
		return
	}
	// validates if details presents

	// proceed to update and prepare response accordingly
	errlog.Trace.Printf("Going to update points: %s %d\n", id, reward.Points)
	if _, err := updateRewardPoints(db, id, reward.Points); err != nil {
		rsp.Message = appUserError(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		encodeJsonAndWrite(w, rsp)
		return
	} else { // else no err from processAddUser
		errlog.Trace.Println(">>>>>>201 - User added:", id)
		rsp.Success = true
		w.WriteHeader(http.StatusCreated)
		encodeJsonAndWrite(w, rsp)
		return
	} // else no err from processAddUser
} // methodPostUser()
