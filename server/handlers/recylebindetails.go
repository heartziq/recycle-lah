package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// var BinData []RecycleBinDetails
var DBCon *sql.DB

// // dfdfd
// type RecycleBinDetails struct {
// 	ID              int     `json:"id"`              // auto increm and primary ID.
// 	BinID           string  `json:"binid"`           // need to assign "A00001" A for HDB recycling bin
// 	BinType         string  `json:"bintype"`         // A: Common Bins, E : E waste, C: Recycling center, M: Mix Bins , W: Workplace Bins
// 	BinLocationLat  float32 `json:"binlocationlat"`  // HC: 311.364587
// 	BinLocationLong float32 `json:"binlocationlong"` // HC: 1.364587
// 	BinAddress      string  `json:"locdescription"`  // Postcode 123456
// 	Postcode        string  `json:"postcode"`        // string but need to conv to int.
// 	UserID          string  `json:"userid"`          // from main site HC: Lanzshot
// 	FBoptions       string  `json:"fboption"`        // "Bin Fullness Status"
// 	ColorCode       string  `json:"colorcode"`       // "Yellow Half Full"
// 	Remarks         string  `json:"remarks"`         // "Please clear asap."
// 	BinStatusUpdate string  `json:"binstatusupdate"` // Completed / Rejected / Submitted
// }

//Get User Past FB from DB.
func getUserFBFromDB(db *sql.DB, userID string) (BinData []RecycleBinDetails) {

	sqlStatement := "SELECT * FROM recycle.RecycleBinsDetails WHERE UserID=?"
	results, err := db.Query(sqlStatement, userID)

	if err != nil {
		fmt.Println("Unable to query and access DB.", err)
		return nil
	}
	defer results.Close()

	for results.Next() {
		var binDet RecycleBinDetails
		if err = results.Scan(
			&binDet.ID,
			&binDet.BinID, &binDet.BinType,
			&binDet.BinLocationLat, &binDet.BinLocationLong,
			&binDet.BinAddress, &binDet.Postcode,
			&binDet.UserID, &binDet.FBoptions,
			&binDet.ColorCode, &binDet.Remarks,
			&binDet.BinStatusUpdate); err != nil {
			fmt.Println("Unable to find any UserID past FB.", err)
			return nil
		}
		fmt.Println("Able to find User FB :", binDet)
		BinData = append(BinData, binDet)
		fmt.Println("FB under USERID : ", BinData)
	}
	return
}

// Client query FB by userID.
func QueryBinFeedback(w http.ResponseWriter, r *http.Request) {

	userid := mux.Vars(r)
	// fmt.Println("PARAMETERS: ", userid)
	v := userid["userID"]
	fmt.Println("User ID Value:", v)

	if r.Method == "GET" {

		jsonBinData := getUserFBFromDB(DBCon, v)
		if jsonBinData == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No Feedback found under this UserID found."))
			return
		}

		json.NewEncoder(w).Encode(jsonBinData)
		// if _, ok := courses[params["courseid"]]; ok {
		// 	json.NewEncoder(w).Encode(courses[params["courseid"]]) //Write out json Map of the wanted courseID
		// } else {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	w.Write([]byte("404 - No course under this ID found."))
		// }
	}
}

func UpdateBinFeedback(w http.ResponseWriter, r *http.Request) {
	var feedBacks RecycleBinDetails
	// if r.Header.Get("Content-type") == "application/json" {

	if r.Method == "POST" {
		reqBody, err := ioutil.ReadAll(r.Body)
		fmt.Println("FB from Client Read is :", string(reqBody))
		if err == nil {
			json.Unmarshal(reqBody, &feedBacks)
			fmt.Println("FB from Client unMarshal :", feedBacks)

			postBinsFeedbackDB(DBCon, feedBacks)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - FeedBack added."))
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Unable to add POST."))
		}
	}
	// }
}

// Adding clients FB to DB
func postBinsFeedbackDB(db *sql.DB, UserFeedBacks RecycleBinDetails) {

	fmt.Println("Adding User FB to DB.")

	sqlStatement := "INSERT INTO recycle.RecycleBinsDetails VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := db.Exec(sqlStatement,
		UserFeedBacks.ID,
		UserFeedBacks.BinID, UserFeedBacks.BinType,
		UserFeedBacks.BinLocationLat, UserFeedBacks.BinLocationLong,
		UserFeedBacks.BinAddress, UserFeedBacks.Postcode,
		UserFeedBacks.UserID, UserFeedBacks.FBoptions,
		UserFeedBacks.ColorCode,
		UserFeedBacks.Remarks,
		UserFeedBacks.BinStatusUpdate)
	if err != nil {
		panic(err)
	} else {
		rows, _ := result.RowsAffected()
		fmt.Println("User Feedbacks successfully added to DB with Rows added:", rows)
	}
}

//client get all bins detail with UserID as NIL.
func GetAllBinDetails(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "List of all courses.\n")
	jsonBinData := getAllBinsFromDB(DBCon)
	json.NewEncoder(w).Encode(jsonBinData)

}

//Get all bins details, NIL User ID in DB.
func getAllBinsFromDB(db *sql.DB) (BinData []RecycleBinDetails) {
	// sqlStatement := "SELECT * FROM recycle.RecycleBinsDetails"
	sqlStatement := "SELECT * FROM recycle.RecycleBinsDetails WHERE UserID = 'NIL'"
	results, err := db.Query(sqlStatement)

	if err != nil {
		fmt.Println("Unable to query and access DB.", err)
		return nil
	}
	defer results.Close()

	fmt.Println("Getting all Bins details from DB.")
	// BinData := make([]RecycleBinDetails, 0)
	for results.Next() {
		var binDet RecycleBinDetails
		if err = results.Scan(
			&binDet.ID,
			&binDet.BinID, &binDet.BinType,
			&binDet.BinLocationLat, &binDet.BinLocationLong,
			&binDet.BinAddress, &binDet.Postcode,
			&binDet.UserID, &binDet.FBoptions,
			&binDet.ColorCode, &binDet.Remarks,
			&binDet.BinStatusUpdate); err != nil {

			fmt.Println("Unable to query all bin data from DB result next.", err)
			return nil
		}
		BinData = append(BinData, binDet)
	}
	return
	// fmt.Println("Append Bins Details from DB : ", BinData)

	// //Marshall the value to JSON
	// // p, err := json.Marshal(binsDataJson)
	// p, err := json.Marshal(BinData)
	// if err != nil {
	// 	fmt.Println("Unable to Marshal to JSON.", err)
	// 	return nil
	// }
	// fmt.Println("Json in string to be send : ", string(p))
	// fmt.Println("Json in bytes to be send : ", p)
	// return p

}
