package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	DBQuery = map[string]string{
		"GetAllPickups": `
		SELECT 
		id, ST_X(coord) as lat, ST_Y(coord) as lng,
		address,
		created_by, attend_by,
		completed
		
		FROM your_db.pickups
		WHERE attend_by='';`,
		"GetPickupInProgress": `
		SELECT 
		id, ST_X(coord) as lat, ST_Y(coord) as lng,
		address,
		created_by, attend_by,
		completed
		
		FROM your_db.pickups
		WHERE created_by=?
        AND attend_by!='';`,
		"PickIAccept": `
		SELECT 
		id, ST_X(coord) as lat, ST_Y(coord) as lng,
		address,
		created_by, attend_by,
		completed
		
		FROM your_db.pickups

		WHERE attend_by=? AND completed=false;
		`,
	}
)

type Pickup struct {
	Id        string  `json:"id"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Address   string  `json:"address"`
	CreatedBy string  `json:"created_by"`
	Collector string  `json:"attend_by"`
	Completed bool    `json:"completed"`
}

type PickupHandler struct {
	Db  *sql.DB
	Tpl *template.Template
	// error logging
	// Info  *log.Logger
	Error *log.Logger
}

func CreatePickupHandler(db *sql.DB, templatePath string) *PickupHandler {
	newPickup := &PickupHandler{Db: db}
	if templatePath != "" {
		newPickup.SetTemplate(templatePath)
	}

	return newPickup
}

// DB queries
func (p *PickupHandler) ListPickup() (users []*Pickup) {
	// access db
	results, err := p.Db.Query(DBQuery["GetAllPickups"])

	if err != nil {

		panic(err.Error())

	}

	for results.Next() {
		// map this type to the record in the table
		c := Pickup{}

		err = results.Scan(
			&c.Id,
			&c.Lat, &c.Lng,
			&c.Address, &c.CreatedBy,
			&c.Collector, &c.Completed,
		)

		if err != nil {

			panic(err.Error())

		}
		users = append(users, &c)

	}
	return
}

// ShowPickup list all pickups exist in DB
// Public access - no auth needed
func (p *PickupHandler) ShowPickup() http.HandlerFunc {
	d := func(w http.ResponseWriter, r *http.Request) {
		data := p.ListPickup()
		json.NewEncoder(w).Encode(data)
	}

	return http.HandlerFunc(d)
}

// RequestPickup create a new entry in the pickup table
// id will be auto-generated by uuid pkg
func (p *PickupHandler) requestPickup(pu *Pickup) error {
	newPickupId := uuid.NewString()

	query := "INSERT INTO your_db.pickups VALUES (?,POINT(?,?),?,?,?,?)"
	result, err := p.
		Db.
		Exec(
			query,
			newPickupId,
			pu.Lat, pu.Lng,
			pu.Address,
			pu.CreatedBy,
			pu.Collector,
			pu.Completed,
		)

	if err != nil {
		return errors.New("error inserting into your_db.pickups")
	}

	rows, _ := result.RowsAffected()
	log.Printf("Insert Successful\t(%d) rows affected", rows)
	return nil
}

func (p *PickupHandler) acceptPickup(pickup_id, collector_id string) error {
	results, err := p.Db.Exec("UPDATE your_db.pickups SET attend_by=? WHERE id=?;", collector_id, pickup_id)
	if err != nil {
		return errors.New("error updating record")
	}

	rows, _ := results.RowsAffected()
	log.Printf("Update Successful\t(%d) rows affected", rows)
	return nil
}

func (p *PickupHandler) approvePickup(pickup_id string) error {
	results, err := p.Db.Exec("UPDATE your_db.pickups SET completed=? WHERE id=?;", true, pickup_id)
	if err != nil {
		return errors.New("error updating record")
	}

	rows, _ := results.RowsAffected()
	log.Printf("Pickup completed\t(%d) rows affected", rows)
	return nil
}

func (p *PickupHandler) showPickupInProgress(user_id string) (users []*Pickup) {
	// access db
	results, err := p.Db.Query(DBQuery["GetPickupInProgress"], user_id)

	if err != nil {

		panic(err.Error())

	}

	for results.Next() {
		// map this type to the record in the table
		c := Pickup{}

		err = results.Scan(
			&c.Id,
			&c.Lat, &c.Lng,
			&c.Address, &c.CreatedBy,
			&c.Collector, &c.Completed,
		)

		if err != nil {

			panic(err.Error())

		}
		users = append(users, &c)

	}
	return
}

func (p *PickupHandler) showAcceptedPickups(collector_id string) (result []*Pickup) {
	// PickIAccept
	// access db
	results, err := p.Db.Query(DBQuery["PickIAccept"], collector_id)

	if err != nil {

		panic(err.Error())

	}

	for results.Next() {
		// map this type to the record in the table
		c := Pickup{}

		err = results.Scan(
			&c.Id,
			&c.Lat, &c.Lng,
			&c.Address, &c.CreatedBy,
			&c.Collector, &c.Completed,
		)

		if err != nil {

			panic(err.Error())

		}
		result = append(result, &c)

	}
	return

}

func (p *PickupHandler) deletePickup(pickup_id string) error {
	results, err := p.Db.Exec("DELETE FROM your_db.pickups WHERE id=?;", pickup_id)
	if err != nil {
		return errors.New("Error deleting record")
	}

	rows, _ := results.RowsAffected()

	if rows < 1 {
		return errors.New("record not found")
	}

	return nil
}

func (p *PickupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// User
	vars := mux.Vars(r)
	role := vars["role"]
	pickup_id := vars["id"]

	if role == "user" {
		switch r.Method {
		case "GET": // Show MY pickup in-progress
			result := p.showPickupInProgress("12345") // replace "12345" with user_id
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(result)
		case "POST": // Request for pickup (status, desc and weight_range, creation_date, updated_date)
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("error parsing body"))
				return
			}
			newPickup := new(Pickup)
			json.Unmarshal(reqBody, newPickup)
			log.Printf("newPickup %v\n", newPickup)
			if err := p.requestPickup(newPickup); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("server/db error"))
				return
			}
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("inserted"))
		case "PUT": // Approve a pickup
			p.approvePickup(pickup_id)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("approve successful!"))

		case "DELETE": // Cancel a pickup
			if err := p.deletePickup(pickup_id); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("sql/server error -- operation aborted!"))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Deleted!"))
		}
	} else {
		// collector
		switch r.Method {
		case "GET": // show current pickup that I am attending
			result := p.showAcceptedPickups("54321")
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(result)
		case "PUT": // cancel or accept (0,1,2,9)
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				payload := map[string]string{}
				json.Unmarshal(reqBody, &payload)
				log.Println("[collector] accept a pickup")

				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("error parsing body"))
					return
				}

				// log.Printf("payload: %v\n", payload)
				p.acceptPickup(payload["pickup_id"], payload["collector_id"])
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Update record(s)"))

				return
			}

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error parsing req.body"))

		}
	}
}

// set up DB conn
func (p *PickupHandler) UseDb(db *sql.DB) error {
	p.Db = db
	return nil // return future potential error(s)
}

// setting up template
func (p *PickupHandler) SetTemplate(path string) {

	p.Tpl = template.Must(template.ParseGlob(path))
}
