package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
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
		myUser := mapUsers[username]
		for _, v := range myUser.UserBookings {
			deleteFromCarsArr(v)
			bookings.deleteBookingNode(v)
		}

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

func newBookingPage(w http.ResponseWriter, r *http.Request) {
	if getUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	var err error
	if r.Method == http.MethodPost {
		car := r.FormValue("cars")
		if car == "none" {
			fmt.Fprintln(w, "A car was not selected, go back to select car")
			return
		}
		date := r.FormValue("date")
		bookingTime, _ := strconv.Atoi(r.FormValue("bookingTime"))
		if bookingTime == 0 {
			fmt.Fprintln(w, "A time was not selected, go back to select time")
			return
		}
		userName := getUser(r).Username
		pickUp := r.FormValue("pickUp")
		dropOff := r.FormValue("dropOff")
		contact, _ := strconv.Atoi(r.FormValue("contact"))
		remarks := r.FormValue("remarks")

		booking, err = bookings.makeNewBooking(car, date, bookingTime, userName, pickUp, dropOff, contact, remarks)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			booking = nil
			return
		}

		myUser := mapUsers[userName]
		myUser.UserBookings = change(myUser.UserBookings, booking)
		myUser.UserBookings = sortBookingsByTime(myUser.UserBookings, len(myUser.UserBookings))
		myUser.UserBookings = sortBookingsByDate(myUser.UserBookings, len(myUser.UserBookings))
		mapUsers[userName] = myUser

		http.Redirect(w, r, "/booking_confirmed", http.StatusSeeOther)
		return
	}
	err = tpl.ExecuteTemplate(w, "newBooking.html", nil)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func bookingConfirmed(w http.ResponseWriter, r *http.Request) {
	if getUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	err := tpl.ExecuteTemplate(w, "bookingConfirmed.html", booking)
	if err != nil {
		panic(errors.New("error executing template"))
	}
	booking = nil
}

func viewAllBookings(w http.ResponseWriter, r *http.Request) {
	if getUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	userBookings := getUser(r).UserBookings

	err := tpl.ExecuteTemplate(w, "viewAllBookings.html", userBookings)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func changeBookingPage(w http.ResponseWriter, r *http.Request) {
	//validate login
	if getUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		var err error
		//get booking number
		bookingId := r.FormValue("bookingId")
		//iterate through slice to get bookingNode
		myUser := getUser(r)
		booking, err = searchId(myUser.UserBookings, bookingId)
		//fmt.Println(booking)
		if err != nil {
			fmt.Fprintf(w, "there are no bookings with that Booking ID, go back to re-enter ID")
			return
		}
	}
	err := tpl.ExecuteTemplate(w, "changeBookingPage.html", booking)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func getChanges(w http.ResponseWriter, r *http.Request) { //TODO check if new car array is nil
	//myUser := getUser(r)
	if getUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost {
		//collect old car data
		oldCarArr := getCarArr(booking.Car)
		oldDate := convertDate(booking.Date)
		oldTime := convertTime(booking.BookingTime)

		car := r.FormValue("cars")
		date := r.FormValue("date")
		bookingTime, _ := strconv.Atoi(r.FormValue("bookingTime"))
		pickUp := r.FormValue("pickUp")
		dropOff := r.FormValue("dropOff")
		contact, _ := strconv.Atoi(r.FormValue("contact"))
		remarks := r.FormValue("remarks")

		if car == "none" && date == "" && bookingTime == 0 {
			if pickUp != "" {
				booking.PickUp = pickUp
			}
			if dropOff != "" {
				booking.DropOff = dropOff
			}
			if contact != 0 {
				booking.ContactInfo = contact
			}
			if remarks != "" {
				booking.Remarks = remarks
			}

			http.Redirect(w, r, "/print_changed_booking", http.StatusSeeOther)
			return
		}
		if car == "none" {
			car = booking.Car
		}
		//collect new car data
		newCarArr := getCarArr(car)
		newDate := convertDate(date)
		newTime := convertTime(bookingTime)

		if newCarArr[newDate][newTime] != nil { //check for empty timeslot
			err := errors.New("there is already a booking at that time and date")
			fmt.Fprintf(w, "Error: %v , go back to select a new slot", err)
			return
		}
		booking.Car = car
		if date != "" {
			booking.Date = date
		}
		if bookingTime != 0 {
			booking.BookingTime = bookingTime
		}
		if pickUp != "" {
			booking.PickUp = pickUp
		}
		if dropOff != "" {
			booking.DropOff = dropOff
		}
		if contact != 0 {
			booking.ContactInfo = contact
		}
		if remarks != "" {
			booking.Remarks = remarks
		}

		//if car, date or time is changed, booking is moved and old booking is deleted
		newCarArr[newDate][newTime] = oldCarArr[oldDate][oldTime]
		oldCarArr[oldDate][oldTime] = nil
		//sort userBookings slice
		myUser := mapUsers[getUser(r).Username]
		myUser.UserBookings = sortBookingsByTime(myUser.UserBookings, len(myUser.UserBookings))
		myUser.UserBookings = sortBookingsByDate(myUser.UserBookings, len(myUser.UserBookings))
		mapUsers[getUser(r).Username] = myUser
		http.Redirect(w, r, "/print_changed_booking", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(w, "changeBooking.html", booking)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func printChangedBooking(w http.ResponseWriter, r *http.Request) {
	if getUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	err := tpl.ExecuteTemplate(w, "printChangedBooking.html", booking)
	if err != nil {
		panic(errors.New("error executing template"))
	}
	booking = nil
}

func deleteBookingPage(w http.ResponseWriter, r *http.Request) {
	//validate login
	if getUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		var err error
		//get booking number
		bookingId := r.FormValue("bookingId")
		//iterate through slice to get bookingNode
		myUser := getUser(r)
		booking, err = searchId(myUser.UserBookings, bookingId)
		//fmt.Println(booking)
		if err != nil {
			fmt.Fprintf(w, "there are no bookings with that Booking ID, go back to re-enter ID")
			return
		}
	}
	err := tpl.ExecuteTemplate(w, "deleteBookingPage.html", booking)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func deleteBooking(w http.ResponseWriter, r *http.Request) {
	if getUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	myUser := getUser(r)
	deleteFromCarsArr(booking)
	if err := deleteBookingUserArr(myUser, booking); err != nil {
		_ = fmt.Errorf("error: %s", err)
	}
	if err := bookings.deleteBookingNode(booking); err != nil {
		_ = fmt.Errorf("error: %s", err)
	}
	booking = nil

	err := tpl.ExecuteTemplate(w, "deleteConfirmed.html", nil)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}
