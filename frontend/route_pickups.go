package main

import (
	"io/ioutil"
	"net/http"

	"frontend/errlog"
)

type pickup struct {
	Id        string  `json:"id"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Address   string  `json:"address"`
	CreatedBy string  `json:"created_by"`
	Collector string  `json:"attend_by"`
	Completed bool    `json:"completed"`
}

// curl -X GET http://localhost:5000/api/v1/rewards/USER1234?key=secretkey
func userPickupList(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("\n\n***userPickupList***")

	//  get data from session

	sess, err := getSession(r)
	errlog.Trace.Println(sess)
	if err != nil {
		errlog.Error.Println("error getting session")
		setFlashCookie(w, "Unauthorized access")
		message(w, r)
		return
	}
	url := "http://localhost:5000/api/v1/pickups/4" + "?key=secretkey&role="
	if sess.isCollector {
		url = url + "collector"
	} else {
		url = url + "user"
	}

	/*
		apiReq, err := http.NewRequest("GET", url, nil)
		// bearer := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyMTIzNCIsImV4cCI6MTYyMzQ3NTMwMywiaXNzIjoidGVzdCJ9.USu2NiQ9vcWHGCeV2m1JhZ23P6r5yCL7UY-m-zeVLBg"
		bearer := sess.token
		apiReq.Header.Add("Authorization", bearer)
		apiReq.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		errlog.Trace.Println("apiReq=", apiReq)
		response, err := client.Do(apiReq)
	*/
	client := &http.Client{}
	errlog.Trace.Println("url=", url)
	response, err := client.Get(url)
	if err != nil {
		errlog.Error.Println("client.Get")
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
	// var rsp RewardPointsResponse
	// // var pickupList []pickup
	// if err := json.Unmarshal(data1, &rsp); err != nil {
	// 	errlog.Error.Println("unmarshal error", err)
	// 	setFlashCookie(w, "Applicatin error (500)")
	// 	message(w, r)
	// 	return
	// } else {
	// 	errlog.Info.Println("response body (unmarshal)=", rsp)
	// 	if !rsp.Success {
	// 		errlog.Error.Println("api returns false", err)
	// 		setFlashCookie(w, rsp.Message)
	// 		message(w, r)
	// 		return
	// 	}
	// 	data := struct {
	// 		UserName     string
	// 		Since        string
	// 		RewardPoints int
	// 		Token        string
	// 	}{}
	// 	data.RewardPoints = rsp.Points
	executeTemplate(w, "user_pickup_list.gohtml", nil)
	return
	// }
}
