package main

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

var tpl *template.Template
var booking *BookingInfoNode

var funcMap = template.FuncMap{
	"add": add,
}

func init() {
	tpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*"))
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/admin_login", adminLogin)
	http.HandleFunc("/admin_index", adminIndex)
	http.HandleFunc("/admin_delete_users", deleteUsers)
	http.HandleFunc("/admin_delete_sessions", deleteSessions)
	http.HandleFunc("/admin_view_delete_bookings", adminViewDeleteBookings)
	http.HandleFunc("/admin_delete_booking_confirmed", adminDeleteBookingConfirmed)
	http.HandleFunc("/new_booking", newBookingPage)
	http.HandleFunc("/booking_confirmed", bookingConfirmed)
	http.HandleFunc("/view_all_bookings", viewAllBookings)
	http.HandleFunc("/change_booking_page", changeBookingPage)
	http.HandleFunc("/get_changes", getChanges)
	http.HandleFunc("/print_changed_booking", printChangedBooking)
	http.HandleFunc("/delete_booking_page", deleteBookingPage)
	http.HandleFunc("/delete_confirmed", deleteBooking)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	err := http.ListenAndServe("localhost:5221", nil)
	if err != nil {
		panic(errors.New("error starting server"))
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

func logout(w http.ResponseWriter, r *http.Request) { //FIXME logout does not work with validation
	//if getUser(r).Username != "" {
	//	http.Redirect(w, r, "/", http.StatusSeeOther)
	//	return
	//}

	sessionCookie, _ := r.Cookie("sessionId")

	delete(mapSessions, sessionCookie.Value)
	sessionCookie = &http.Cookie{
		Name:   "sessionId",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, sessionCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
