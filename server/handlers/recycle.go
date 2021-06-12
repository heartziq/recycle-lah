package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var BinData []RecycleBinDetails

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

type RBinHandler struct {
	Db  *sql.DB
	Tpl *template.Template
	// error logging
	// Info  *log.Logger
	Error *log.Logger
}

// set up DB conn
func (p *RBinHandler) UseDb(db *sql.DB) error {
	p.Db = db
	return nil // return future potential error(s)
}

// setting up template
func (p *RBinHandler) SetTemplate(path string) {

	p.Tpl = template.Must(template.ParseGlob(path))
}

func CreateRBinHandler(db *sql.DB, templatePath string) *RBinHandler {
	newPickup := &RBinHandler{Db: db}
	if templatePath != "" {
		newPickup.SetTemplate(templatePath)
	}

	return newPickup
}

func (rBin *RBinHandler) GetAllBinDetails() http.HandlerFunc {
	d := func(w http.ResponseWriter, r *http.Request) {

		jsonBinData := rBin.getAllBinsFromDB()
		json.NewEncoder(w).Encode(jsonBinData)

	}

	return http.HandlerFunc(d)
}

func (rBin *RBinHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userid := mux.Vars(r)
	v := userid["userID"]

	switch r.Method {
	case "GET":

		jsonBinData := rBin.getUserFBFromDB(v)
		if jsonBinData == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No Feedback found under this UserID found."))
			return
		}

		json.NewEncoder(w).Encode(jsonBinData)

	case "POST":
		var feedBacks RecycleBinDetails

		reqBody, err := ioutil.ReadAll(r.Body)
		log.Println("FB from Client Read is :", string(reqBody))
		if err == nil {
			json.Unmarshal(reqBody, &feedBacks)
			log.Println("FB from Client unMarshal :", feedBacks)

			rBin.postBinsFeedbackDB(feedBacks)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - FeedBack added."))
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Unable to add POST."))
		}

	}
}

//Get User Past FB from DB.
func (rb *RBinHandler) getUserFBFromDB(userID string) (BinData []RecycleBinDetails) {

	sqlStatement := "SELECT * FROM RecycleBinsDetails WHERE UserID=?"
	results, err := rb.Db.Query(sqlStatement, userID)

	if err != nil {
		log.Println("Unable to query and access DB.", err)
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
			log.Println("Unable to find any UserID past FB.", err)
			return nil
		}
		log.Println("Able to find User FB :", binDet)
		BinData = append(BinData, binDet)
		log.Println("FB under USERID : ", BinData)
	}
	return
}

// Adding clients FB to DB
func (rb *RBinHandler) postBinsFeedbackDB(UserFeedBacks RecycleBinDetails) {

	log.Println("Adding User FB to DB.")

	sqlStatement := "INSERT INTO RecycleBinsDetails VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := rb.Db.Exec(sqlStatement,
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
		log.Println("User Feedbacks successfully added to DB with Rows added:", rows)
	}
}

//Get all bins details, NIL User ID in DB.
func (rBin *RBinHandler) getAllBinsFromDB() (BinData []RecycleBinDetails) {
	sqlStatement := "SELECT * FROM RecycleBinsDetails WHERE UserID = 'NIL'"
	results, err := rBin.Db.Query(sqlStatement)

	if err != nil {
		log.Println("Unable to query and access DB.", err)
		return nil
	}
	defer results.Close()

	log.Println("Getting all Bins details from DB.")

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

			log.Println("Unable to query all bin data from DB result next.", err)
			return nil
		}
		BinData = append(BinData, binDet)
	}
	return

}
