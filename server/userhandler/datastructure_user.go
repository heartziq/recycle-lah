package userhandler

import (
	"database/sql"
	"time"
)

// const cfgFile defines the name of the configuration file
// for database connection and port number of the web server.
// The file is in json format.
const cfgFile = "./config.json"

// var db *sql.DB os a global variable for the database driver
var db *sql.DB

// var config Configuration is a global variable that stores the configuration details.
var config Configuration

// struct Configuration defines attributes in the configuration file
// and its corresponding json attributes name.
type Configuration struct {
	DB struct {
		Name   string `json:"name"` // database name
		User   string `json:"user"` // database user
		Host   string `json:"host"`
		Port   string `json:"port"`
		driver string
	} `json:"database"`
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}

// type Response struct defines the data transfer format from
// the REST API server to the REST API client.
type Response struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
}

// type UserInfo struct mimics database table - user.
type UserInfo struct { // follows database structure
	Id           string `json:"id"`
	password     string
	Email        string    `json:"code"`
	IsDeleted    bool      `json:"deleted"`
	IsCollector  bool      `json:"collector"`
	RewardPoints int       `json:"reward_points"`
	CreatedDate  time.Time `json:"create_date"`
	UpdatedDate  time.Time `json:"updated_date"`
}

type msUser struct { // micro services data structure for user
	id       string
	Password string `json:"password"`
	Email    string `json:"email"`
}

// type NewUser struct defines the data format
// sent to the REST API server.
type NewUser struct {
	Password  string `json:"password"`
	Email     string `json:"email"`
	Collector bool   `json:"collector"`
}

type dbNewUser struct {
	id        string
	password  string
	email     string
	collector bool
}

// type RewardPoints struct defines the data format
// sent to the REST API client.
type RewardPointsResponse struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
	Points    int    `json:"reward_points"`
}

// type AdditionalRewardPoints struct defines the data format
// sent to the REST API server for reward points update.
type AdditionalRewardPoints struct {
	Points int `json:"reward_points"`
}

// type AdditionalRewardPoints struct defines the data format
// sent to the REST API server for reward points update.
type UserVerification struct {
	Password string `json:"password"`
}

// type UserInfoResponse struct defines the data format
// sent to the REST API client.
type UserInfoResponse struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
	UserInfo
}
