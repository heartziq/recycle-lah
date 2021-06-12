package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"frontend/errlog"
)

func dummyPostSook(w http.ResponseWriter, r *http.Request) {
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
