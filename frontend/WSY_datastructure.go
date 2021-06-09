package main

import (
	"regexp"
)

// changed by Sook on June 8 integration
type user struct {
	UserName string
	// Password string
	Key                string
	userId             string
	userName           string
	token              string
	isCollector        bool
	email              string
	sessionCreatedTime int64
}

var (
	// tpl           *template.Template
	emailRegex = regexp.MustCompile("^[\\w!#$%&'*+/=?`{|}~^-]+(?:\\.[\\w!#$%&'*+/=?`{|}~^-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,6}$") // regular expression
	// mapUsers      = map[string]user{"r@l.com": user{"recycle", "278d0e77-76c2-4447-bbfb-6fb032f57414"}}                                 //**temporary use data
	mapUsers      = map[string]user{} //**temporary use data
	mapSessions   = map[string]string{}
	matchPassword = map[string]string{"r@l.com": "password"} //**need to get from Database
)
