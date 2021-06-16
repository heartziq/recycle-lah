package main

import (
	"html/template"
	"regexp"
	"time"
)

// this file contains the variables used in the application

// user defines the data structure to store user details
type user struct {
	UserName string
	// Password string
	Key                string
	userId             string
	userName           string
	token              string
	isCollector        bool // true for collector, false for user
	email              string
	sessionCreatedTime int64
}

type NewUser struct {
	Password  string `json:"password"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	Collector bool   `json:"collector"`
}

var (
	tpl        *template.Template
	tplFuncs   = template.FuncMap{"rangeStruct": RangeStructer, "fShortDate": fShortDate, "fmtDate": fmtDate}
	emailRegex = regexp.MustCompile("^[\\w!#$%&'*+/=?`{|}~^-]+(?:\\.[\\w!#$%&'*+/=?`{|}~^-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,6}$") // regular expression
	// mapUsers      = map[string]user{"r@l.com": user{"recycle", "278d0e77-76c2-4447-bbfb-6fb032f57414"}}                                 //**temporary use data
	mapUsers      = map[string]user{} //**temporary use data
	mapSessions   = map[string]string{}
	matchPassword = map[string]string{"r@l.com": "password"} //**need to get from Database
)

// var tpl *template.Template

// var mapSession = map[string]Session{}

// const cfgFile defines the name of configuration file
// for host and account related information.
// The file is in json format.
const cfgFile = "./config.json"

// var config Configuration is a global variable that contains Host and
// user configuration information.
var config Configuration

// type Configuration struct defines attributes in the configuration file
// and its corresponding json attributes name.
// APIKey must belongs to the same account.
type Configuration struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
	Account string `json:"account"`
}

// type Response struct defines the data transfer format from
// REST API server.
type Response struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
}

// type RewardPointsRequest2 contains token to be sent to api server
type RewardPointsRequest2 struct {
	Token string `json:"token"`
}

// type RewardPointsResponse defines response data structure for
// GET reward points from the rest api server
type RewardPointsResponse struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
	Points    int    `json:"reward_points"`
}

// type UserVerification struct defines the data format
// sent to the REST API server for user verification
type UserVerification struct {
	Password string `json:"password"`
}

// type UserInfoResponse struct defines the data format
// received from api server upon successful authenticated
type UserInfoResponse struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
	UserInfo
}

// type UserInfo defines data structure for user details
// It is part of the response from api server (UserInfoResponse)
type UserInfo struct { // follows database structure
	Id           string    `json:"id"`
	Email        string    `json:"code"`
	UserName     string    `json:"user_name"`
	IsCollector  bool      `json:"collector"`
	RewardPoints int       `json:"reward_points"`
	CreatedDate  time.Time `json:"create_date"`
	UpdatedDate  time.Time `json:"updated_date"`
	token        string
}
