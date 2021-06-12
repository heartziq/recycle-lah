package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"frontend/errlog"
)

// func viewPoints1(w http.ResponseWriter, r *http.Request) {
// 	data := struct {
// 		UserName     string
// 		Since        string
// 		RewardPoints int
// 		Token        string
// 	}{}
// 	if err := getRewardPoints(w, r, data); err == nil {
// 		http.Redirect(w, r, "/welcome", http.StatusFound)
// 		return
// 	}
// 	executeTemplate(w, "view_points.gohtml", data)
// }

// curl -X GET http://localhost:5000/api/v1/rewards/USER1234?key=secretkey
func viewPoints(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("\n\n***getRewardPoints***")

	//  get data from session

	sess, err := getSession(r)
	if err != nil {
		errlog.Error.Println("error getting session")
		setFlashCookie(w, "Unauthorized access")
		message(w, r)
		return
	}
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
	// bearer := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyMTIzNCIsImV4cCI6MTYyMzQ3NTMwMywiaXNzIjoidGVzdCJ9.USu2NiQ9vcWHGCeV2m1JhZ23P6r5yCL7UY-m-zeVLBg"
	bearer := "Bearer " + sess.token
	apiReq.Header.Add("Authorization", bearer)
	apiReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	errlog.Trace.Println("bearer=", bearer)
	response, err := client.Do(apiReq)

	if err != nil {
		errlog.Error.Println("client.Do")
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}
	data1, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		errlog.Error.Printf("ReadAll: response status code:%+v err:%s\n", response.StatusCode, err.Error())
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}
	errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))
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
		data := struct {
			UserName     string
			Since        string
			RewardPoints int
			Token        string
		}{}
		data.RewardPoints = rsp.Points
		executeTemplate(w, "view_points.gohtml", data)
		return
	}
}
