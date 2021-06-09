package main

import (
	"database/sql"
	"errors"

	"frontend/errlog"

	_ "github.com/go-sql-driver/mysql"
)

type DBConnection struct { // to construct database connection string
	dbType   string
	user     string
	password string
	hostAddr string
	port     string
	name     string
}

var db *sql.DB //database driver
type DBConnectionError error

var (
	errDBOpen = DBConnectionError(errors.New("failed to open database"))
	errDBPing = DBConnectionError(errors.New("failed to ping database"))
)

// connects to database and returns connected driver
func XopenDB(driver, credential string) (*sql.DB, error) {
	// var db *sql.DB
	// db, err = sql.Open("mysql", "user1:password@tcp(127.0.0.1:3306)/MYSTOREDB")
	db, err := sql.Open(driver, credential)
	if err != nil {
		errlog.Error.Println(err.Error())
		return db, errDBOpen
	}
	errlog.Info.Println("DB opened")
	if err = db.Ping(); err != nil {
		errlog.Error.Println(err.Error())
		return db, errDBPing
	}
	return db, nil
}
