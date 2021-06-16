package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"frontend/errlog"
)

// func viewPoints() retrieves reward points from api server and display
// the rewoard points on the page
func viewPoints(w http.ResponseWriter, r *http.Request) {
	// errlog.Trace.Println("\n\n***getRewardPoints***")
	data := struct {
		PageName     string
		UserName     string
		Since        string
		RewardPoints int
		Token        string
		MsgToUser    string
	}{PageName: "My Reward Points"}

	//  get data from session

	sess, err := getSession(r)
	if err != nil {
		errlog.Error.Println("error getting session")
		setFlashCookie(w, "Unauthorized access")
		message(w, r)
		return
	}

	data.UserName = sess.userName

	var reward RewardPointsRequest2
	reward.Token = sess.token

	jsonValue, err := json.Marshal(reward)
	if err != nil {
		errlog.Error.Println("marshal error")
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}

	url := "http://localhost:5000/api/v1/rewards/" + sess.userId + "?key=secretkey"
	apiReq, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonValue))
	// set header to send token over to the server
	bearer := "Bearer " + sess.token
	apiReq.Header.Add("Authorization", bearer)
	apiReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	errlog.Trace.Println("bearer=", bearer)
	// sends the http request
	response, err := client.Do(apiReq)

	if err != nil {
		errlog.Error.Println("client.Do", err)
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}

	// retrieve the response body
	data1, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		errlog.Error.Printf("ReadAll: response status code:%+v err:%s\n", response.StatusCode, err.Error())
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}
	errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))

	// unmarshal the response body
	var rsp RewardPointsResponse
	if err := json.Unmarshal(data1, &rsp); err != nil {
		errlog.Error.Println("unmarshal error", err)
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	} else {
		errlog.Info.Println("response body (unmarshal)=", rsp)
		if !rsp.Success {
			errlog.Error.Println("api returns false", err)
			setFlashCookie(w, rsp.Message)
			message(w, r)
			return
		}
		data.RewardPoints = rsp.Points
		executeTemplate(w, "view_points.gohtml", data)
		return
	}
}
