package handlers

import (
	"fmt"
	"html/template"

	"net/http"

	browser "github.com/heartziq/recycle-lah/server/utility"
	uuid "github.com/satori/go.uuid"
)

var tpl *template.Template
var mapSessions = map[string]string{} //to record total user that have loged in

//User
func NewUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST": //creat new user

		// Send data to gohtml file
		Data := struct {
			PageName  string
			UserName  string
			MsgToUser string
		}{PageName: "New User Registration", UserName: ""}
		tpl.ExecuteTemplate(w, "NewUser.gohtml", Data)

		//this struct is to store into DB
		type User struct {
			UserName string
			email    string
			password string
			uuid     string
		}

		// var currentUser user
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			email := r.FormValue("email")
			password := r.FormValue("password")
			confirmpassword := r.FormValue("confirmpassword")
			// _, nameFound := ctms.ExistingCustomer(username)
			emailOk := browser.IsEmailValid(email)

			// check user input condition
			if /*!nameFound &&*/ password == confirmpassword && emailOk {
				id := uuid.NewV4().String()
				browser.SetCookie(w, id)
				ToDB := User{username, email, password, id}
				// 	UserName: username,
				// 	email:    email,
				// 	password: password,
				// 	uuid:     id,
				// }
				fmt.Println(ToDB)
				// ctms.AddCustomer(username, password, id) //** send data back to api server
				mapSessions[id] = username
				Data.MsgToUser = "New User Registration Done! You may process to log in."
				defer fmt.Fprintf(w, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			} else if !emailOk {
				Data.MsgToUser = "Please enter correct email!"
				defer fmt.Fprintf(w, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			} else if password != confirmpassword {
				Data.MsgToUser = "Confirm Password is not same!"
				defer fmt.Fprintf(w, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)
			} /*} else if nameFound {
			// fmt.Println("Htmlmain.newUser - name in existing data")
			fmt.Scanf(Data.MsgToUser, "Please use other user name! '%v' has been taken!", username)
			defer fmt.Fprintf(w, "<br><script>document.getElementById('MsgToUser').innerHTML = '%v';</script>", Data.MsgToUser)*/
		}

	}

}
