package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

var tpl *template.Template
var booking *BookingInfoNode
var wg sync.WaitGroup

var errLog *os.File
var userLog *os.File
var adminLog *os.File

var ErrorLogger *log.Logger
var UserLogger *log.Logger
var AdminLogger *log.Logger

var funcMap = template.FuncMap{
	"add": add,
}

func init() {
	var err error
	tpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*"))
	errLog, err = os.OpenFile("logs/errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	userLog, err = os.OpenFile("logs/userLoginAndLogout.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open user log file:", err)
	}
	adminLog, err = os.OpenFile("logs/adminLoginAndLogout.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open admin log file:", err)
	}

	flags := log.LstdFlags | log.Lshortfile
	ErrorLogger = log.New(io.MultiWriter(errLog, os.Stderr), "ERROR: ", flags)
	UserLogger = log.New(io.MultiWriter(userLog, os.Stderr), "USER LOG: ", flags)
	AdminLogger = log.New(io.MultiWriter(adminLog, os.Stderr), "ADMIN LOG: ", flags)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()
	defer func(errLog *os.File) {
		err := errLog.Close()
		if err != nil {
			ErrorLogger.Fatalln("unable to close error.log")
		}
	}(errLog)
	defer func(userLog *os.File) {
		err := userLog.Close()
		if err != nil {
			ErrorLogger.Fatalln("unable to close userLoginAndLogout.log")
		}
	}(userLog)

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
		ErrorLogger.Fatalf("error starting server: %v", err)
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

		if !IsAlphabetic(username) {
			err := errors.New("username includes invalid characters")
			UserLogger.Printf("LOGIN UNSUCCESSFUL: '%s', %v", username, err)
			http.Error(w, "Username and/or password is not valid", http.StatusUnauthorized)
			return
		}

		if !IsAlphanumeric(password) {
			err := errors.New("password includes invalid characters")
			UserLogger.Printf("LOGIN UNSUCCESSFUL: %v tried to login, %v", username, err)
			http.Error(w, "Username and/or password is not valid", http.StatusUnauthorized)
			return
		}

		myUser, ok := mapUsers[username]
		if !ok {
			err := errors.New("username not found")
			UserLogger.Printf("LOGIN UNSUCCESSFUL: '%s', %v", username, err)
			http.Error(w, "Username and/or password is not valid", http.StatusUnauthorized)
			return
		}
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			err := errors.New("wrong password")
			UserLogger.Printf("LOGIN UNSUCCESSFUL: %v tried to login, %v", username, err)
			http.Error(w, "Username and/or password is not valid", http.StatusUnauthorized)
			return
		}

		setSessionIDCookie(w, username) //TODO set cookie expiry
		UserLogger.Printf("LOGIN SUCCESSFUL: %s logged in", username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
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
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {

	UserLogger.Printf("USER LOGOUT: %v has logged out", getUser(r).Username)
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
