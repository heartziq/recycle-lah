package main

import "time"

// For Sin Yaw - changes on 7 June
//  - added UserName, mandatory for email and username
// /*
type NewUser struct {
	Password  string `json:"password"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	Collector bool   `json:"collector"`
}

// */

// type UserInfo struct mimics database table - user.
// /*
type UserInfo struct { // follows database structure
	Id           string `json:"id"`
	password     string
	Email        string    `json:"code"`
	UserName     string    `json:"user_name"`
	IsCollector  bool      `json:"collector"`
	RewardPoints int       `json:"reward_points"`
	CreatedDate  time.Time `json:"create_date"`
	UpdatedDate  time.Time `json:"updated_date"`
	token        string
}

// */

// //  Session defines fields to store session for a particular user
// type Session struct {
// 	uuid               string
// 	userId             string
// 	userName           string
// 	token              string
// 	isCollector        bool
// 	email              string
// 	rewardPoints       int
// 	updatedDate        time.Time
// 	sessionCreatedTime int64
// }
