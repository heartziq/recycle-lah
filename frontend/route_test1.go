package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"frontend/errlog"
)

// curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyMTIzNCIsImV4cCI6MTYyMzQ3NTMwMywiaXNzIjoidGVzdCJ9.USu2NiQ9vcWHGCeV2m1JhZ23P6r5yCL7UY-m-zeVLBg" -H "Content-Type: application/json" -X GET http://localhost:5000/api/v2/users/USER4567?key=secretkey -d {\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyMTIzNCIsImV4cCI6MTYyMzQ3NTMwMywiaXNzIjoidGVzdCJ9.USu2NiQ9vcWHGCeV2m1JhZ23P6r5yCL7UY-m-zeVLBg\"} -v
func testToken(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data = make(map[string]interface{})
	if r.Method == http.MethodPost {
		if err := testSendToken(w, r, data); err == nil {
			http.Redirect(w, r, "/welcome", http.StatusFound)
			return
		}
	}
	executeTemplate(w, "test1.gohtml", data)
}

// 	url := "http://localhost:5000/api/v2/users/" + "user1234" + "?key=secretkey" with header
func testSendTokenOk(w http.ResponseWriter, r *http.Request, data map[string]interface{}) error {
	errlog.Trace.Println("\n\n***testSendToken***")
	if r.Method != http.MethodPost {
		errlog.Error.Println("Wrong method in testSendToken")
		return errors.New("Wrong method")
	}
	if err := r.ParseForm(); err != nil {
		errlog.Error.Println("err in ParseForm", err)
		data["Error"] = "processing error"
		return err
	}
	//  get data from form

	var reward RewardPointsRequest2
	reward.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyMTIzNCIsImV4cCI6MTYyMzQ3NTMwMywiaXNzIjoidGVzdCJ9.USu2NiQ9vcWHGCeV2m1JhZ23P6r5yCL7UY-m-zeVLBg"

	jsonValue, err := json.Marshal(reward)
	if err != nil {
		errlog.Error.Println("error in marshal", err)
		return err
	}

	url := "http://localhost:5000/api/v2/users/" + "user1234" + "?key=secretkey"
	apiReq, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonValue))
	// bearer := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyMTIzNCIsImV4cCI6MTYyMzQ3NTMwMywiaXNzIjoidGVzdCJ9.USu2NiQ9vcWHGCeV2m1JhZ23P6r5yCL7UY-m-zeVLBg"
	bearer := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyMTIzNCIsImV4cCI6MTYyMzQ3NTMwMywiaXNzIjoidGVzdCJ9.USu2NiQ9vcWHGCeV2m1JhZ23P6r5yCL7UY-m-zeVLBg"
	apiReq.Header.Add("Authorization", bearer)
	apiReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(apiReq)

	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)
		return err
	}
	data1, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		errlog.Error.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())
		return err
	}
	errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))
	var rsp RewardPointsResponse
	if err := json.Unmarshal(data1, &rsp); err != nil {
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

// curl -X GET http://localhost:5000/api/v1/rewards/USER1234?key=secretkey
func testSendToken(w http.ResponseWriter, r *http.Request, data map[string]interface{}) error {
	errlog.Trace.Println("\n\n***testSendToken***")
	if r.Method != http.MethodPost {
		errlog.Error.Println("Wrong method in testSendToken")
		return errors.New("Wrong method")
	}
	if err := r.ParseForm(); err != nil {
		errlog.Error.Println("err in ParseForm", err)
		data["Error"] = "processing error"
		return err
	}
	//  get data from form

	var reward RewardPointsRequest2
	reward.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyMTIzNCIsImV4cCI6MTYyMzQ3NTMwMywiaXNzIjoidGVzdCJ9.USu2NiQ9vcWHGCeV2m1JhZ23P6r5yCL7UY-m-zeVLBg"

	jsonValue, err := json.Marshal(reward)
	if err != nil {
		errlog.Error.Println("error in marshal", err)
		return err
	}

	url := "http://localhost:5000/api/v1/rewards/" + "user1234" + "?key=secretkey"
	apiReq, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonValue))
	// bearer := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyMTIzNCIsImV4cCI6MTYyMzQ3NTMwMywiaXNzIjoidGVzdCJ9.USu2NiQ9vcWHGCeV2m1JhZ23P6r5yCL7UY-m-zeVLBg"
	bearer := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyMTIzNCIsImV4cCI6MTYyMzQ3NTMwMywiaXNzIjoidGVzdCJ9.USu2NiQ9vcWHGCeV2m1JhZ23P6r5yCL7UY-m-zeVLBg"
	apiReq.Header.Add("Authorization", bearer)
	apiReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(apiReq)

	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)
		return err
	}
	data1, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		errlog.Error.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())
		return err
	}
	errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))
	var rsp RewardPointsResponse
	if err := json.Unmarshal(data1, &rsp); err != nil {
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
