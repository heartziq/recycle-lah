package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

// func viewCompletedJobs() retrieves list of (assigned) jobs from the api server
// It filters to display only completed job
func viewCompletedJobs(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("\n\n***viewCompletedJobs***")

	data := struct {
		PageName   string
		UserName   string
		PickupList map[string]pickup
		MsgToUser  string
	}{PageName: "Requested List"}

	//  get data from session
	sess, err := getSession(r)
	errlog.Trace.Println(sess)
	if err != nil {
		errlog.Error.Println("error getting session")
		setFlashCookie(w, "Unauthorized access")
		message(w, r)
		return
	}

	// send request to api server
	url := "http://localhost:5000/api/v1/pickups/4" + "?key=secretkey&role="
	if sess.isCollector {
		url = url + "collector"
	} else {
		url = url + "user"
	}

	apiReq, err := http.NewRequest("GET", url, nil)

	bearer := "Bearer " + sess.token
	apiReq.Header.Add("Authorization", bearer)
	apiReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	errlog.Trace.Println("bearer=", bearer)
	response, err := client.Do(apiReq)

	if err != nil {
		errlog.Error.Println("client.Get")
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}

	data.UserName = sess.userName

	// process response body
	data1, err := ioutil.ReadAll(response.Body)

	errlog.Trace.Println("data1=", string(data1))
	defer response.Body.Close()
	if err != nil {
		errlog.Error.Printf("ReadAll: response status code:%+v err:%s\n", response.StatusCode, err.Error())
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}
	// errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))

	// storing the return results into map
	var allListitems []pickup
	var mapPickup = map[string]pickup{}
	json.Unmarshal(data1, &allListitems)
	for _, v := range allListitems {
		if v.Completed {
			mapPickup[v.Id] = v
		}
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

		// for each job in the map, calls completedJob() to send request to the api server

		for _, v := range jobs {
			var job map[string]string
			job = make(map[string]string)
			job["pickup_id"] = v
			job["collector_id"] = ""
			completedJob(job)
		}

		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
	}

	executeTemplate(w, "view_completed_jobs.gohtml", data)
	// }
}

// func userPickupList() shows list of jobs that have been assigned
// It allows user to select multipile job and indicate as completed
// It then calls completedJob() repeated to fulfill the request
// This could be implemented in Go Routine
func userPickupList(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("\n\n***userPickupList***")

	data := struct {
		PageName   string
		UserName   string
		PickupList map[string]pickup
		MsgToUser  string
	}{PageName: "Requested List"}

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

	// send request to api server
	apiReq, err := http.NewRequest("GET", url, nil)

	bearer := "Bearer " + sess.token
	apiReq.Header.Add("Authorization", bearer)
	apiReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	errlog.Trace.Println("bearer=", bearer)
	response, err := client.Do(apiReq)

	if err != nil {
		errlog.Error.Println("client.Get")
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}

	data.UserName = sess.userName

	// get response body and process the body
	data1, err := ioutil.ReadAll(response.Body)

	errlog.Trace.Println("data1=", string(data1))
	defer response.Body.Close()
	if err != nil {
		errlog.Error.Printf("ReadAll: response status code:%+v err:%s\n", response.StatusCode, err.Error())
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}
	// errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))

	// storing the return results into map
	var allListitems []pickup
	var mapPickup = map[string]pickup{}
	json.Unmarshal(data1, &allListitems)
	for _, v := range allListitems {
		if !v.Completed {
			mapPickup[v.Id] = v
		}
	}
	data.PickupList = mapPickup

	errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, data.PickupList)

	if r.Method == http.MethodPost {
		// stores user selected item into map
		var jobs []string
		for k, v := range mapPickup {
			errlog.Trace.Printf("key, value: %v, %v", k, v)
			tf := r.FormValue(string(k))
			if tf != "" {
				jobs = append(jobs, tf)
			}
		}
		errlog.Trace.Println("getback", jobs)

		// for each job in the map, calls completedJob() to send request to the api server
		for _, v := range jobs {
			var job map[string]string
			job = make(map[string]string)
			job["pickup_id"] = v
			job["collector_id"] = ""
			completedJob(job)
		}

		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
	}

	executeTemplate(w, "user_pickup_list.gohtml", data)

}

// func completedJob() send request to api server to approve the job completion
func completedJob(jobs map[string]string) (bool, error) {
	errlog.Trace.Println("completedJob: ", jobs)

	// send request to api server
	url := "http://localhost:5000/api/v1/pickups/" + jobs["pickup_id"] + "?key=secretkey&role=user"
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, url, nil)
	request.Header.Set("Content-Type", "application/json")
	errlog.Trace.Printf("put REQUEST=%v\n", request)
	response, err := client.Do(request)
	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)
		return false, err
	}

	// process the response
	data1, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	errlog.Error.Printf("response status code:%+v err:%s\n", response.StatusCode, string(data1))
	if err != nil {
		errlog.Error.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())
		return false, err
	}
	// interpret results, could have used status code
	if string(data1) == "Update record(s)" {
		return true, nil
	}
	return false, errors.New(string(data1))
}

// func requestPickup() gets input from user and send request to api server to submit the request
// The template is designed based on api Version 2 of the api server and
// some of the input fiels are not relevant for api Version 1
// As version 2 was not ready for the time of integration, the new fields would not be sent to the client
// For current version, we are not ready to integrate with google map api.
// It calls func getCoordinate() that returns some dummy latitude and longitude.
func requestPickup(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("\n\n***requestPickup***")
	data := struct {
		PageName  string
		UserName  string
		MsgToUser string
	}{PageName: "Request Pickup Order"}

	sess, err := getSession(r)
	errlog.Trace.Println(sess)
	if err != nil {
		errlog.Error.Println("error getting session")
		return
	}
	data.UserName = sess.userName

	if r.Method == http.MethodPost {

		var newOrder pickup

		newOrder.CreatedBy = sess.userId
		description := r.FormValue("description")
		weight_range := r.FormValue("weightrange")
		newOrder.Address = r.FormValue("address")
		// newOrder.Collector := r.FormValue("collector")
		// newOrder.Completed := r.FormValue("completed")
		postCode := r.FormValue("postcode")
		// to get coordinates of a post code
		lat, lng := getCoordinate(postCode)
		newOrder.Lat = lat
		newOrder.Lng = lng
		errlog.Trace.Println("Html Form :", newOrder, description, weight_range)

		// send request to api server
		url := "http://localhost:5000/api/v1/pickups/4" + "?key=secretkey&role="
		if sess.isCollector {
			url = url + "collector"
		} else {
			url = url + "user"
		}

		jsonValue, err := json.Marshal(newOrder)
		if err != nil {
			errlog.Error.Println("error in marshal", err)
			return
		}

		client := &http.Client{}
		errlog.Trace.Println("url=", url)
		response, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			errlog.Error.Println("client.Post")
			setFlashCookie(w, "Applicatin error (500)")
			message(w, r)
			return
		}

		// process response
		data1, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))
		errlog.Trace.Println("response.Header=", response.Header)
		errlog.Trace.Println("response.Body=", string(data1))
		if string(data1) == "inserted" {
			http.Redirect(w, r, "/welcome", http.StatusSeeOther)
		} else {
			data.MsgToUser = "We are sorry that your request has not been accepted!"
			fmt.Fprintf(w, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", data.MsgToUser)
		}
	}
	executeTemplate(w, "user_requested_form.gohtml", data)
}

// func getCoordinate() returns some dummy values based on the post code entered
func getCoordinate(postCode string) (lat float64, lng float64) {
	if postCode == "" {
		return 1.33221, 103.77466 //  599489 Ngee Ann Poly
	}
	if len(postCode) < 3 {
		return 1.33221, 103.77466 //  599489 Ngee Ann Poly
	}
	if postCode[0:3] == "560" { // amk hub 569933
		return 1.36974, 103.84873
	}
	if postCode[0:3] == "310" { // toa payoh 310480 toa payoh hub
		return 1.33458, 103.84903
	}
	return 1.33221, 103.77466 //  599489 Ngee Ann Poly
}
