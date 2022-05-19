package main

import (
	"assignment3/limoBookingApp"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

var errLog *os.File
var userLog *os.File
var adminLog *os.File

var ErrorLogger *log.Logger
var UserLogger *log.Logger
var AdminLogger *log.Logger

func init() {
	var err error

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

	router.HandleFunc("/", limoBookingApp.Index)
	router.HandleFunc("/login", limoBookingApp.Login)
	router.HandleFunc("/logout", limoBookingApp.Logout)
	router.HandleFunc("/signup", limoBookingApp.Signup)
	router.HandleFunc("/admin_login", limoBookingApp.AdminLogin)
	router.HandleFunc("/admin_index", limoBookingApp.AdminIndex)
	router.HandleFunc("/admin_delete_users", limoBookingApp.DeleteUsers)
	router.HandleFunc("/admin_delete_sessions", limoBookingApp.DeleteSessions)
	router.HandleFunc("/admin_view_delete_bookings", limoBookingApp.AdminViewDeleteBookings)
	router.HandleFunc("/admin_delete_booking_confirmed", limoBookingApp.AdminDeleteBookingConfirmed)
	router.HandleFunc("/new_booking", limoBookingApp.NewBookingPage)
	router.HandleFunc("/booking_confirmed", limoBookingApp.BookingConfirmed)
	router.HandleFunc("/view_all_bookings", limoBookingApp.ViewAllBookings)
	router.HandleFunc("/change_booking_page", limoBookingApp.ChangeBookingPage)
	router.HandleFunc("/get_changes", limoBookingApp.GetChanges)
	router.HandleFunc("/print_changed_booking", limoBookingApp.PrintChangedBooking)
	router.HandleFunc("/delete_booking_page", limoBookingApp.DeleteBookingPage)
	router.HandleFunc("/delete_confirmed", limoBookingApp.DeleteBooking)
	router.Handle("/favicon.ico", http.NotFoundHandler())

	err := http.ListenAndServeTLS(":5221", "ssl/cert.pem", "ssl/key.pem", router)
	//err := http.ListenAndServe("localhost:5221", nil)
	if err != nil {
		ErrorLogger.Fatalf("error starting server: %v", err)
	}
}
