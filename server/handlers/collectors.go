package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
)

type Collection struct {
	Db  *sql.DB
	Tpl *template.Template
	// error logging
	// Info  *log.Logger
	Error *log.Logger
}

type Users struct {
	Id            string
	Password      string
	Email         string
	Is_collector  string
	Is_deleted    string
	Reward_points string
	Creation_date string
	Updated_date  string
}

var c Collection

func GetDb(db *sql.DB) error {
	c.Db = db
	return nil // return future potential error(s)
}

func getUser() map[int]interface{} {
	mapUsers := make(map[int]interface{})
	results, err := c.Db.Query("Select id FROM recycle.session")
	if err != nil {
		panic(err)
	}

	i := 0
	for results.Next() {
		var id int
		// results.Scan(&u.Id, &u.Password, &u.Email, &u.Is_collector, &u.Is_deleted, &u.Reward_points, &u.Creation_date, &u.Updated_date)
		results.Scan(&id)
		mapUsers[i] = id
		// fmt.Println(mapUsers)
		i++
	}

	fmt.Println(mapUsers)
	return mapUsers
}
