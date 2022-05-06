package main

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

import (
	"fmt"
	"html/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/admin_login", adminLogin)
	http.HandleFunc("/admin_index", adminIndex)
	http.HandleFunc("/admin_delete_users", deleteUsers)
	http.HandleFunc("/admin_delete_sessions", deleteSessions)
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
	if getUser(r).Username != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "One or more inputs are empty", http.StatusForbidden)
			return
		}

		myUser, ok := mapUsers[username]
		if !ok {
			http.Error(w, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}

		setSessionIDCookie(w, username)
		//fmt.Println(mapSessions)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func signup(w http.ResponseWriter, r *http.Request) {

	//check if already logged in
	if getUser(r).Username != "" {
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

func logout(w http.ResponseWriter, r *http.Request) { //FIXME logout not deleting cookie
	//if getUser(r).Username != "" {
	//	http.Redirect(w, r, "/", http.StatusSeeOther)
	//	return
	//}

	sessionCookie, _ := r.Cookie("sessionId")

	delete(mapSessions, sessionCookie.Value)
	fmt.Println(mapSessions)
	sessionCookie = &http.Cookie{
		Name:   "sessionId",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, sessionCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func adminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "One or more inputs are empty", http.StatusForbidden)
			return
		}

		myAdminUser, ok := mapAdmins[username]
		if !ok {
			http.Error(w, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}
		err := bcrypt.CompareHashAndPassword(myAdminUser.Password, []byte(password))
		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}

		setSessionIDCookie(w, username)
		http.Redirect(w, r, "/admin_index", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(w, "adminLogin.html", nil)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func adminIndex(w http.ResponseWriter, r *http.Request) {
	currentAdmin := getAdmin(r)
	err := tpl.ExecuteTemplate(w, "adminIndex.html", currentAdmin)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func deleteUsers(w http.ResponseWriter, r *http.Request) {
	if getAdmin(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		delete(mapUsers, username)
	}
	err := tpl.ExecuteTemplate(w, "deleteUsers.html", mapUsers)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func deleteSessions(w http.ResponseWriter, r *http.Request) {
	if getAdmin(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		sessionId := r.FormValue("sessionId")
		delete(mapSessions, sessionId)
	}
	err := tpl.ExecuteTemplate(w, "deleteSessions.html", mapSessions)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}
