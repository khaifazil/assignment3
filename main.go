package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"sync"
)

var tpl *template.Template
var booking *BookingInfoNode
var wg sync.WaitGroup

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

	router := mux.NewRouter()

	router.HandleFunc("/", index)
	router.HandleFunc("/login", login)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/signup", signup)
	router.HandleFunc("/admin_login", adminLogin)
	router.HandleFunc("/admin_index", adminIndex)
	router.HandleFunc("/admin_delete_users", deleteUsers)
	router.HandleFunc("/admin_delete_sessions", deleteSessions)
	router.HandleFunc("/admin_view_delete_bookings", adminViewDeleteBookings)
	router.HandleFunc("/admin_delete_booking_confirmed", adminDeleteBookingConfirmed)
	router.HandleFunc("/new_booking", newBookingPage)
	router.HandleFunc("/booking_confirmed", bookingConfirmed)
	router.HandleFunc("/view_all_bookings", viewAllBookings)
	router.HandleFunc("/change_booking_page", changeBookingPage)
	router.HandleFunc("/get_changes", getChanges)
	router.HandleFunc("/print_changed_booking", printChangedBooking)
	router.HandleFunc("/delete_booking_page", deleteBookingPage)
	router.HandleFunc("/delete_confirmed", deleteBooking)
	router.Handle("/favicon.ico", http.NotFoundHandler())

	err := http.ListenAndServeTLS(":5221", "ssl/cert.pem", "ssl/key.pem", router)
	//err := http.ListenAndServe("localhost:5221", nil)
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
			http.Error(w, "Username and/or password is not valid", http.StatusUnauthorized)
			return
		}
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			http.Error(w, "Username and/or password is not valid", http.StatusUnauthorized)
			return
		}

		setSessionIDCookie(w, username) //TODO set cookie expiry
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
