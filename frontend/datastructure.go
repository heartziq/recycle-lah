package main

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

// type NewUser struct defines the data format
// sent to the REST API server.
// type NewUser struct {
// 	Password  string `json:"password"`
// 	Email     string `json:"email"`
// 	Collector bool   `json:"collector"`
// }

// type Response struct defines the data transfer format from
// REST API server.
type Response struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
}

// test token
type RewardPointsRequest2 struct {
	Token string `json:"token"`
}

type RewardPointsResponse struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
	Points    int    `json:"reward_points"`
}

// type AdditionalRewardPoints struct defines the data format
// sent to the REST API server for reward points update.
type UserVerification struct {
	Password string `json:"password"`
}

// type UserInfo received from api server
// type UserInfo struct {
// 	Id           string `json:"id"`
// 	Email        string `json:"code"`
// 	token        string
// 	isDeleted    bool      `json:"deleted"`
// 	IsCollector  bool      `json:"collector"`
// 	RewardPoints int       `json:"reward_points"`
// 	CreatedDate  time.Time `json:"create_date"`
// 	UpdatedDate  time.Time `json:"updated_date"`
// }

// type UserInfoResponse struct defines the data format
// sent to the REST API client.
type UserInfoResponse struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
	UserInfo
}

//  Session defines fields to store session for a particular user
// type Session struct {
// 	uuid               string
// 	userId             string
// 	token              string
// 	isCollector        bool
// 	email              string
// 	rewardPoints       int
// 	updatedDate        time.Time
// 	sessionCreatedTime int64
// }
