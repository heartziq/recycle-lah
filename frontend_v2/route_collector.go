package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"frontend/errlog"
	"io/ioutil"
	"net/http"
)

// func showJobsAvailable) calls function to retrieve list of available jobs from api server
// It allows user to accept multiple jobs
// It then call func changeJobCollector repeating to accept job
// This could be implemented in go routine
func showJobsAvailable(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("\n\n***showJobsAvailable***")

	data := struct {
		PageName   string
		UserName   string
		PickupList map[string]pickup
		MsgToUser  string
	}{PageName: "Show Jobs Available"}

	//  get data from session
	sess, err := getSession(r)
	errlog.Trace.Println(sess)
	if err != nil {
		errlog.Error.Println("error getting session")
		setFlashCookie(w, "Unauthorized access")
		message(w, r)
		return
	}
	data.UserName = sess.userName

	// sending request to api server
	// url := "http://localhost:5000/api/v1/pickups"
	url := "http://localhost:5000/api/v2/pickups"

	client := &http.Client{}
	errlog.Trace.Println("url=", url)
	response, err := client.Get(url)
	if err != nil {
		errlog.Error.Println("client.Get")
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}

	// get response body and process the body
	data1, err := ioutil.ReadAll(response.Body)
	errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))
	defer response.Body.Close()
	if err != nil {
		errlog.Error.Printf("ReadAll: response status code:%+v err:%s\n", response.StatusCode, err.Error())
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}

	// storing the return results into map
	var allListitems []pickup
	var mapPickup = map[string]pickup{}
	json.Unmarshal(data1, &allListitems)
	for _, v := range allListitems {
		mapPickup[v.Id] = v
	}
	data.PickupList = mapPickup

	errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, data.PickupList)

	if r.Method == http.MethodPost {
		var jobs []string

		// stores user selected item into map
		for k, v := range mapPickup {
			errlog.Trace.Printf("key, value: %v, %v", k, v)
			tf := r.FormValue(string(k))
			if tf != "" {
				jobs = append(jobs, tf)
			}
		}
		errlog.Trace.Println("getback", jobs)

		// for each job in the map, calls changeJobCollector() to send request to the api server
		for _, v := range jobs {
			var job map[string]string
			job = make(map[string]string)
			job["pickup_id"] = v
			job["collector_id"] = sess.userId
			changeJobCollector(job)
		}
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
	}

	executeTemplate(w, "collector_pickup_jobs.gohtml", data)

}

// func showMyJobs() calls function to retrieve collector's job
// It allows user to select multiple jobs for cancellation
// It then call func changeJobCollector repeating to cancel job
// This could be implemented in go routine
func showMyJobs(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("\n\n***showMyJobs***")

	data := struct {
		PageName   string
		UserName   string
		PickupList map[string]pickup
		MsgToUser  string
	}{PageName: "View My Jobs"}

	//  get data from session

	sess, err := getSession(r)
	if err != nil {
		errlog.Error.Println("error getting session")
		setFlashCookie(w, "Unauthorized access")
		message(w, r)
		return
	}

	data.UserName = sess.userName

	// sending request to api server

	url := "http://localhost:5000/api/v1/pickups/3" + "?key=secretkey&role=collector"
	apiReq, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
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

	// get response body and process the body
	data1, err := ioutil.ReadAll(response.Body)

	defer response.Body.Close()
	errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))
	if err != nil {
		errlog.Error.Printf("ReadAll: response status code:%+v err:%s\n", response.StatusCode, err.Error())
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}

	// storing the return results into map
	var allListitems []pickup
	var mapPickup = map[string]pickup{}
	json.Unmarshal(data1, &allListitems)
	for _, v := range allListitems {
		mapPickup[v.Id] = v
	}
	data.PickupList = mapPickup

	errlog.Trace.Println("data.pickuplist", data.PickupList)

	if r.Method == http.MethodPost {
		var jobs []string
		// stores user selected item into map

		for k, v := range mapPickup {
			errlog.Trace.Printf("key, value: %v, %v", k, v)
			tf := r.FormValue(string(k))
			if tf != "" {
				jobs = append(jobs, tf)
			}
		}
		errlog.Trace.Println("getback", jobs)

		// for each job in the map, calls changeJobCollector() to send request to the api server
		for _, v := range jobs {
			var job map[string]string
			job = make(map[string]string)
			job["pickup_id"] = v
			job["collector_id"] = ""
			changeJobCollector(job)
		}
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
	}

	executeTemplate(w, "collector_jobs.gohtml", data)
}

// func changeJobCollector() send request to api server to cancel/accept job
func changeJobCollector(jobs map[string]string) (bool, error) {

	errlog.Trace.Println("changeJobCollector: ", jobs)
	jsonValue, err := json.Marshal(jobs)
	if err != nil {
		errlog.Error.Println("error in marshal", err)
		return false, err
	}

	// send request to api server
	url := "http://localhost:5000/api/v1/pickups/5" + "?key=secretkey&role=collector"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, url,
		bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	errlog.Trace.Printf("put REQUEST=%v\n", request)
	response, err := client.Do(request)
	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)
		return false, err
	}

	// process request body
	data1, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	errlog.Error.Printf("response status code:%+v err:%s\n", response.StatusCode, string(data1))
	if err != nil {
		errlog.Error.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())
		return false, err
	}
	// interpret returned results, could check for status code
	if string(data1) == "Update record(s)" {
		return true, nil
	}
	return false, errors.New(string(data1))
}
