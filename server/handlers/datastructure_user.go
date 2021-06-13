package handlers

import (
	"time"
)

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
	UserName     string    `json:"user_name"`
	isDeleted    bool      `json:"deleted"`
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

// type NewUser struct defines the data format for Create New User
// sent from the api client
type NewUser struct {
	Password  string `json:"password"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	Collector bool   `json:"collector"`
}

// type dbNewUser contains id besides details from NewUser
// it is used to pass data for the database operation
type dbNewUser struct {
	id        string
	password  string
	email     string
	userName  string
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

// type UserVerification struct defines the data format
// sent to the REST API server for verification
type UserVerification struct {
	Password string `json:"password"`
}

// type UserInfoResponse struct defines the data format that contain user details
// that sent to the REST API client.
type UserInfoResponse struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
	UserInfo
}
