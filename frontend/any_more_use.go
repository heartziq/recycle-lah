package main

import "net/http"

//Log Out
func logOut(res http.ResponseWriter, req *http.Request) {

	Data := struct {
		PageName string
		UserName string
	}{PageName: "Log Out", UserName: "bye-bye"}
	Cookie, err := req.Cookie("RecycleLah")
	if err == nil {
		Cookie.MaxAge = -1
		delete(mapUsers, mapSessions[Cookie.Value])
		delete(mapSessions, Cookie.Value)
		http.SetCookie(res, Cookie)
		// fmt.Println("Cookie deleted")
	} else {
		// fmt.Println("No Cookie found and deleted")
		http.Redirect(res, req, "/logIn", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(res, "LogOut.gohtml", Data)
}
