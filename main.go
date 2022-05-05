package main

import (
	"errors"
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login) //TODO login page
	http.HandleFunc("/signup", signUp)
	//http.HandleFunc("/logout", logout) //TODO logout page
	http.Handle("/favicon.ico", http.NotFoundHandler())
	err := http.ListenAndServe("localhost:5221", nil)
	if err != nil {
		panic(errors.New("error starting server"))
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	currentUser := getUser(r)
	err := tpl.ExecuteTemplate(w, "index.html", currentUser)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func login(w http.ResponseWriter, r *http.Request) {

}

func signUp(w http.ResponseWriter, r *http.Request) {

	//check if already logged in
	if getUser(r).Username != "" {
		//if already logged in redirect to index
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		//if not logged in createUser
		createUser(w, r)
		//redirect back to main after createUser
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//execute template
	err := tpl.ExecuteTemplate(w, "signup.html", nil)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

//func alreadyLoggedIn(r *http.Request) bool {
//	sessionCookie, err := r.Cookie("sessionId")
//	if err != nil {
//		return false
//	}
//	userName := mapSessions[sessionCookie.Value]
//	_, ok := mapUsers[userName]
//	return ok
//}
