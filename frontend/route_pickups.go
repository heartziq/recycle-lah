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

func viewCompletedJobs(w http.ResponseWriter, r *http.Request) {
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

	apiReq, err := http.NewRequest("GET", url, nil)

	bearer := "Bearer " + sess.token
	apiReq.Header.Add("Authorization", bearer)
	apiReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	errlog.Trace.Println("bearer=", bearer)
	response, err := client.Do(apiReq)

	// client := &http.Client{}
	// errlog.Trace.Println("url=", url)
	// response, err := client.Get(url)
	if err != nil {
		errlog.Error.Println("client.Get")
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}

	data.UserName = sess.userName
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
		for k, v := range mapPickup {
			errlog.Trace.Printf("key, value: %v, %v", k, v)
			tf := r.FormValue(string(k))
			if tf != "" {
				jobs = append(jobs, tf)
			}
		}
		errlog.Trace.Println("getback", jobs)
		for _, v := range jobs {
			var job map[string]string
			job = make(map[string]string)
			job["pickup_id"] = v
			job["collector_id"] = ""
			completedJob(job)
		}

		// update backend

		// if string(data1) == "Update record(s)" {
		// 	return
		// }
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
	}

	executeTemplate(w, "view_completed_jobs.gohtml", data)
	// }
}

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

	apiReq, err := http.NewRequest("GET", url, nil)

	bearer := "Bearer " + sess.token
	apiReq.Header.Add("Authorization", bearer)
	apiReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	errlog.Trace.Println("bearer=", bearer)
	response, err := client.Do(apiReq)

	// client := &http.Client{}
	// errlog.Trace.Println("url=", url)
	// response, err := client.Get(url)
	if err != nil {
		errlog.Error.Println("client.Get")
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}

	data.UserName = sess.userName
	data1, err := ioutil.ReadAll(response.Body)

	errlog.Trace.Println("data1=", data1)
	defer response.Body.Close()
	if err != nil {
		errlog.Error.Printf("ReadAll: response status code:%+v err:%s\n", response.StatusCode, err.Error())
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}
	// errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))
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
		var jobs []string
		for k, v := range mapPickup {
			errlog.Trace.Printf("key, value: %v, %v", k, v)
			tf := r.FormValue(string(k))
			if tf != "" {
				jobs = append(jobs, tf)
			}
		}
		errlog.Trace.Println("getback", jobs)
		for _, v := range jobs {
			var job map[string]string
			job = make(map[string]string)
			job["pickup_id"] = v
			job["collector_id"] = ""
			completedJob(job)
		}

		// update backend

		// if string(data1) == "Update record(s)" {
		// 	return
		// }
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
	}

	executeTemplate(w, "user_pickup_list.gohtml", data)
	// }
}

func completedJob(jobs map[string]string) (bool, error) {
	errlog.Trace.Println("completedJob: ", jobs)

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
	data1, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	errlog.Error.Printf("response status code:%+v err:%s\n", response.StatusCode, string(data1))
	if err != nil {
		errlog.Error.Printf("response status code:%+v err:%s\n", response.StatusCode, err.Error())
		return false, err
	}
	if string(data1) == "Update record(s)" {
		return true, nil
	}
	return false, errors.New(string(data1))
}

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
		// //  get data from form
		// dt := time.Now()
		// hr, mi, se := dt.Clock()
		// y, m, d := dt.Date()
		// ID := int(y) + int(m) + d + hr + mi + se
		// newOrder.Id = sess.userId + "pickup_" + strconv.Itoa(ID)
		// errlog.Error.Println("data: ", data)

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
