package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frontend/errlog"
	"io/ioutil"
	"net/http"
)

type RecycleBinDetails struct {
	ID              int     `json:"id"`              // auto increm and primary ID.
	BinID           string  `json:"binid"`           // need to assign "A00001" A for HDB recycling bin
	BinType         string  `json:"bintype"`         // A: Common Bins, E : E waste, C: Recycling center, M: Mix Bins , W: Workplace Bins
	BinLocationLat  float32 `json:"binlocationlat"`  // HC: 311.364587
	BinLocationLong float32 `json:"binlocationlong"` // HC: 1.364587
	BinAddress      string  `json:"locdescription"`  // Postcode 123456
	Postcode        string  `json:"postcode"`        // string but need to conv to int.
	UserID          string  `json:"userid"`          // from main site HC: Lanzshot
	FBoptions       string  `json:"fboption"`        // "Bin Fullness Status"
	ColorCode       string  `json:"colorcode"`       // "Yellow Half Full"
	Remarks         string  `json:"remarks"`         // "Please clear asap."
	BinStatusUpdate string  `json:"binstatusupdate"` // Completed / Rejected / Submitted
}

// var binFeedbacks recycleBinDetails

// WEB Server port and url.
const baseURLBin = "http://localhost:5000/api/v1/recyclebindetails"

func IndexBin(res http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(res, "indexBin.gohtml", nil)
}

// Get and send user feedback.
func recycleBinFB(res http.ResponseWriter, req *http.Request) {

	Data := struct {
		PageName  string
		UserName  string
		Collector string
	}{PageName: "Post Feedback"}

	user, err := getSession(req)

	if err != nil {
		http.Redirect(res, req, "index.gohtml", http.StatusFound)
		return
	}
	Data.UserName = user.userName
	if user.isCollector {
		Data.Collector = "Y"
	}

	if req.Method == http.MethodPost {
		sess, err := getSession(req)
		errlog.Trace.Println(sess)
		if err != nil {
			errlog.Error.Println("error getting session")
			return
		}
		userID := sess.userId

		binFeedbacks := RecycleBinDetails{
			ID:              0,
			BinID:           "NIL",
			BinType:         "NIL",
			BinLocationLat:  0,
			BinLocationLong: 0,
			BinAddress:      req.FormValue("Binaddress"),
			Postcode:        req.FormValue("postcode"),
			// UserID:          "Lanzs", //To be puck in.
			UserID:          userID, //To be puck in.
			FBoptions:       req.FormValue("FBoptions"),
			ColorCode:       req.FormValue("binfull"),
			Remarks:         req.FormValue("remarks"),
			BinStatusUpdate: "Submitted",
		}
		fmt.Println("User binFeedbacks : ", binFeedbacks)

		jsonString, err := json.Marshal(binFeedbacks)
		if err != nil {
			fmt.Println("Json Mashal error :", err)
		}
		// os.Stdout.Write(jsonString)

		// apiCode := binFeedbacks.BinID
		fmt.Println("Sending User FB Via POST")
		// response, err := http.Post(baseURLBin+"/feedback", "application/json", bytes.NewBuffer(jsonString)) //POST to create course.
		response, err := http.Post(baseURLBin+"/NIL", "application/json", bytes.NewBuffer(jsonString)) //POST to create course.
		// response, err := http.Post(baseURL+"/"+apiCode, "application/json", bytes.NewBuffer(jsonString)) //POST to create course.

		if err != nil {
			fmt.Printf("The HTTP POST request failed with error %s\n", err)
		} else {
			defer response.Body.Close()
			data, _ := ioutil.ReadAll(response.Body)
			// fmt.Println("Add POST Course data:", data)
			// fmt.Println("Status Code : ", response.StatusCode)
			fmt.Println("Bin Details added :\n", string(data))
		}
	}
	tpl.ExecuteTemplate(res, "recycleBinOptions.gohtml", Data)
}

// Get user FB from DB.
func queryFB(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Show user past Feedback with status.")

	Data := struct {
		PageName  string
		UserName  string
		Collector string
		FeedBacks []RecycleBinDetails
	}{PageName: "Query Feedback"}

	user, err := getSession(req)
	errlog.Trace.Println(user)
	if err != nil {
		errlog.Error.Println("error getting session")
		return
	}

	if err != nil {
		http.Redirect(res, req, "index.gohtml", http.StatusFound)
		return
	}
	Data.UserName = user.userName
	if user.isCollector {
		Data.Collector = "Y"
	}

	// var feedBacks []RecycleBinDetails
	userID := user.userId

	fmt.Println(baseURLBin + "/feedback/" + userID)
	// response, err := http.Get(baseURLBin + "/feedback/" + userID)
	response, err := http.Get(baseURLBin + "/" + userID)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		fmt.Println("FB Query Status Code : ", response.StatusCode) //200 OK.
	}
	defer response.Body.Close()

	JsonByteData, _ := ioutil.ReadAll(response.Body)
	fmt.Println("User FB Recieved from Server in Json:", string(JsonByteData))

	// convert JSON to object
	json.Unmarshal(JsonByteData, &Data.FeedBacks)
	fmt.Println("User FB  Details Recieved from Server in String:", Data.FeedBacks)

	tpl.ExecuteTemplate(res, "showUserFB.gohtml", Data)

}

// show only recyclebins.
func showRecycleBins(res http.ResponseWriter, req *http.Request) {
	// response, err := http.Get(baseURLBin)

	Data := struct {
		PageName  string
		UserName  string
		Collector string
		Bins      []RecycleBinDetails
	}{PageName: "Show Recycle Bins"}

	user, err := getSession(req)
	if err != nil {
		http.Redirect(res, req, "index.gohtml", http.StatusFound)
		return
	}
	Data.UserName = user.userName
	if user.isCollector {
		Data.Collector = "Y"
	}

	response, err := http.Get(baseURLBin + "/NIL")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		fmt.Println("Show all Bins Status Code : ", response.StatusCode) //200 OK.
	}
	defer response.Body.Close()

	JsonByteData, _ := ioutil.ReadAll(response.Body)
	fmt.Println("ALL Bin Details Recieved from Server in Json:", JsonByteData)

	// convert JSON to object
	json.Unmarshal(JsonByteData, &Data.Bins)
	errlog.Trace.Println("ALL Bin Details Recieved from Server in String:", Data.Bins)

	tpl.ExecuteTemplate(res, "showRecycleBinsDetails.gohtml", Data)
}
