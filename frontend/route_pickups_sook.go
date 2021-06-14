package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"frontend/errlog"
)

// delete user pickup job
// curl -x DELETE http://localhost:5000/api/v1/pickups/040be96d-20ae-4f05-8ef0-026983b613ea?key=secretkey&role=user
func dummyCalldeletePickup(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("\n\n***dummyCalldeletePickup***")
	type Data struct {
		Action   string
		UserId   string
		PickupId string
		Status   bool
		Message  string
	}
	userId := "sook6666"
	pickupId := "c736eb0a-a71a-4247-83fe-dabad2702ec8"
	var data Data
	data.Action = "delete pickup"
	data.UserId = userId
	data.PickupId = pickupId
	errlog.Trace.Println("data=", data)
	deleted, err := dummydeletePickup(userId, pickupId)
	if err != nil {
		data.Message = err.Error()
		errlog.Error.Println(err)
	}
	errlog.Trace.Println("deleted=", deleted)

	data.Status = deleted
	executeTemplate(w, "dummy.gohtml", data)
}

func dummydeletePickup(userId, pickupId string) (bool, error) {
	url := "http://localhost:5000/api/v1/pickups/" + pickupId + "?key=secretkey&role=user"
	request, err := http.NewRequest(http.MethodDelete,
		url, nil)
	errlog.Trace.Printf("delete REQUEST= %+v\n", request)
	client := &http.Client{}
	if err != nil {
		errlog.Error.Printf("The HTTP request failed with error %s\n", err)
		return false, err
	}
	errlog.Trace.Println("request=", request)
	response, err := client.Do(request)
	if err != nil {
		errlog.Error.Println(err)
		return false, err
	}

	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		errlog.Error.Println(err)
		return false, err
	}
	errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data))
	errlog.Trace.Println("response.Header=", response.Header)
	errlog.Trace.Println("response.Body=", string(data))
	if string(data) == "Deleted!" {
		return true, nil
	}
	return false, errors.New("Failed to deleted.")
}

// ===============  collector - show all jobs available
// curl http://localhost:5000/api/v1/pickups
func dummyCallCollectorShowJobsAvailable(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("\n\n***dummyCallCollectorShowJobsAvailable***")

	data := struct {
		PageName   string
		UserName   string
		PickupList []pickup
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
	url := "http://localhost:5000/api/v1/pickups"
	// if sess.isCollector {
	// 	url = url + "collector"
	// } else {
	// 	url = url + "user"
	// }

	client := &http.Client{}
	errlog.Trace.Println("url=", url)
	response, err := client.Get(url)
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
	json.Unmarshal(data1, &data.PickupList)
	errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, data.PickupList)

	executeTemplate(w, "dummy.gohtml", data)
	// }
}

// func dummySubmit(w http.ResponseWriter, r *http.Request) {
// 	errlog.Trace.Println("\n\n***requestPickup***")
// 	data := struct {
// 		Response string
// 	}{}

// 	if r.Method == http.MethodPost {
// 		if err := r.ParseForm(); err != nil {
// 			errlog.Error.Println("err in ParseForm", err)
// 			return
// 		}
// 		postCode := r.FormValue("dummy1")
// 		lat, lng := getCoordinate(postCode)
// 		var job pickup
// 		job.Id = "id"
// 		job.Lat = lat
// 		job.Lng = lng
// 		job.Address = "Marine Parade 1"
// 		job.CreatedBy = "sook6666"
// 		job.Collector = "testcol"
// 		job.Completed = false

// 		// sess, err := getSession(r)
// 		// if err != nil {
// 		// 	errlog.Error.Println("error getting session")
// 		// 	setFlashCookie(w, "Unauthorized access")
// 		// 	message(w, r)
// 		// 	return
// 		// }

// 		jsonValue, err := json.Marshal(job)
// 		if err != nil {
// 			errlog.Error.Println("marshal error")
// 			setFlashCookie(w, "Applicatin error (500)")
// 			message(w, r)
// 			return
// 		}
// 		errlog.Trace.Println("job=", job)
// 		url := "http://localhost:5000/api/v1/pickups/1" + "?key=secretkey&role=user"
// 		client := &http.Client{}

// 		response, err := client.Post(url,
// 			"application/json", bytes.NewBuffer(jsonValue))

// 		if err != nil {
// 			errlog.Error.Println(err)
// 			setFlashCookie(w, "Applicatin error (500)")
// 			message(w, r)
// 			return
// 		}
// 		data1, err := ioutil.ReadAll(response.Body)
// 		defer response.Body.Close()
// 		errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))
// 		errlog.Trace.Println("response.Header=", response.Header)
// 		errlog.Trace.Println("response.Body=", string(data1))
// 		data.Response = string(data1)
// 		if err != nil {
// 			errlog.Error.Printf("ReadAll: response status code:%+v err:%s\n", response.StatusCode, err.Error())
// 			setFlashCookie(w, "Applicatin error (500)")
// 			message(w, r)
// 			return
// 		}
// 	} // if methodPost
// 	executeTemplate(w, "dummy_submit.gohtml", data)
// }

func dummyPost(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("\n\n***requestPickup***")
	data := struct {
		Response string
	}{}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			errlog.Error.Println("err in ParseForm", err)
			return
		}
		postCode := r.FormValue("dummy1")
		lat, lng := getCoordinate(postCode)
		var job pickup
		job.Id = "id"
		job.Lat = lat
		job.Lng = lng
		job.Address = "Marine Parade 1"
		job.CreatedBy = "sook6666"
		job.Collector = "testcol"
		job.Completed = false

		// sess, err := getSession(r)
		// if err != nil {
		// 	errlog.Error.Println("error getting session")
		// 	setFlashCookie(w, "Unauthorized access")
		// 	message(w, r)
		// 	return
		// }

		jsonValue, err := json.Marshal(job)
		if err != nil {
			errlog.Error.Println("marshal error")
			setFlashCookie(w, "Applicatin error (500)")
			message(w, r)
			return
		}
		errlog.Trace.Println("job=", job)
		url := "http://localhost:5000/api/v1/pickups/1" + "?key=secretkey&role=user"
		client := &http.Client{}

		response, err := client.Post(url,
			"application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			errlog.Error.Println(err)
			setFlashCookie(w, "Applicatin error (500)")
			message(w, r)
			return
		}
		data1, err := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))
		errlog.Trace.Println("response.Header=", response.Header)
		errlog.Trace.Println("response.Body=", string(data1))
		data.Response = string(data1)
		if err != nil {
			errlog.Error.Printf("ReadAll: response status code:%+v err:%s\n", response.StatusCode, err.Error())
			setFlashCookie(w, "Applicatin error (500)")
			message(w, r)
			return
		}
	} // if methodPost
	executeTemplate(w, "dummy_submit.gohtml", data)

}

func dummyCalldummyacceptJob(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("\n\n***dummyCallCollectorShowJobsAvailable***")

	data := struct {
		PageName   string
		UserName   string
		PickupList []pickup
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
	// if sess.isCollector {
	// 	url = url + "collector"
	// } else {
	// 	url = url + "user"
	// }

	var jobs map[string]string
	jobs = make(map[string]string)

	jobId := "115e76da-75fc-4353-8e27-d216322f73a0"
	jobs["pickup_id"] = jobId
	jobs["collector_id"] = sess.userId

	ok, err := dummyacceptJob(jobs)
	if ok {
		data.MsgToUser = "job has been assigned"
	} else {
		data.MsgToUser = "Failed to assign job" + jobId
	}
	executeTemplate(w, "dummy.gohtml", data)
	// }
}

func dummyacceptJob(jobs map[string]string) (bool, error) {
	// var jobs map[string]string
	// jobs = make(map[string]string)
	// jobs["c736eb0a-a71a-4247-83fe-dabad2702ec8"] = "testing"

	errlog.Trace.Println("dummyacceptJob: ", jobs)
	jsonValue, err := json.Marshal(jobs)
	if err != nil {
		errlog.Error.Println("error in marshal", err)
		return false, err
	}

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

func dummyViewAttendingJob(w http.ResponseWriter, r *http.Request) {
	errlog.Trace.Println("\n\n***dummyViewAttendingJob***")
	data := struct {
		PageName  string
		UserName  string
		Since     string
		Token     string
		MsgToUser string
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

	url := "http://localhost:5000/api/v1/pickups/3" + "?key=secretkey&role=collector"
	apiReq, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
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
	errlog.Trace.Printf("response status code:%+v\nstring(data1):%+v\n", response.StatusCode, string(data1))
	if err != nil {
		errlog.Error.Printf("ReadAll: response status code:%+v err:%s\n", response.StatusCode, err.Error())
		setFlashCookie(w, "Applicatin error (500)")
		message(w, r)
		return
	}
	data.MsgToUser = string(data1)
	executeTemplate(w, "dummy.gohtml", data)
}
