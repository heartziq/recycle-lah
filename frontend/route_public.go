package main

import (
	"fmt"
	"net/http"
)

//  index() - home page
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("making changes at the remote end")
	Data := struct {
		PageName string
		UserName string
	}{PageName: "Recyle Lah! Home Page"}
	executeTemplate(w, "index.gohtml", Data)
}

//  contact()  provides contact detailss
func contact(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "contact.gohtml", nil)
}

//  unauthorized()  shows error
func unauthorized(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>unauthorized access<h1>")
	// executeTemplate(w, "contact.gohtml", nil)
}

//  signUpSuccess()  shows user account successfully created and link to login
func signUpSuccess(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "signupsuccess.gohtml", nil)
}
