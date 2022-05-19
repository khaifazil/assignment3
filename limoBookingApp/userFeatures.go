package limoBookingApp

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

//Index is the handler for the user index page
func Index(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	err := tpl.ExecuteTemplate(w, "index.html", currentUser)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}

//Login is the handler for the user login page
func Login(w http.ResponseWriter, r *http.Request) {
	if GetUser(r).Username != "" {
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

		SetSessionIDCookie(w, username)
		UserLogger.Printf("LOGIN SUCCESSFUL: %s logged in", username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}

//Signup is the handler for the user signup page.
func Signup(w http.ResponseWriter, r *http.Request) {

	//check if already logged in
	if GetUser(r).Username != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		//if not logged in createUser
		CreateUser(w, r)

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

//Logout is the handler for the admin and user /logout path.
func Logout(w http.ResponseWriter, r *http.Request) {

	UserLogger.Printf("USER LOGOUT: %v has logged out", GetUser(r).Username)
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

//NewBookingPage is the handler for the user make new booking page.
func NewBookingPage(w http.ResponseWriter, r *http.Request) {
	if GetUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	var err error
	if r.Method == http.MethodPost {
		car := r.FormValue("cars")
		if car == "none" {
			_, err := fmt.Fprintln(w, "A car was not selected, go back to select car")
			if err != nil {
				ErrorLogger.Println(err)
			}
			return
		}
		date := r.FormValue("date")
		date = StripHtmlRegex(date)
		if err := CheckDate(date); err != nil {
			_, err := fmt.Fprintf(w, "%v, go back to change date", err)
			if err != nil {
				ErrorLogger.Println(err)
			}
			return
		}
		bookingTime, _ := strconv.Atoi(r.FormValue("bookingTime"))
		if bookingTime == 0 {
			_, err := fmt.Fprintln(w, "A time was not selected, go back to select time")
			if err != nil {
				ErrorLogger.Println(err)
			}
			return
		}
		userName := GetUser(r).Username
		pickUp := r.FormValue("pickUp")
		pickUp = StripHtmlRegex(pickUp)
		dropOff := r.FormValue("dropOff")
		dropOff = StripHtmlRegex(dropOff)
		contact, err := strconv.Atoi(r.FormValue("contact"))
		if err != nil {
			_, err := fmt.Fprintln(w, "Invalid contact number, go back to input new contact number")
			if err != nil {
				ErrorLogger.Println(err)
			}
			return
		}
		if !CheckContactLen(contact) {
			_, err := fmt.Fprintln(w, "Invalid contact number, go back to input new contact number")
			if err != nil {
				ErrorLogger.Println(err)
			}
			return
		}
		remarks := r.FormValue("remarks")
		remarks = StripHtmlRegex(remarks)

		booking, err = bookings.MakeNewBooking(car, date, bookingTime, userName, pickUp, dropOff, contact, remarks)
		if err != nil {
			_, err := fmt.Fprintf(w, "%v", err)
			if err != nil {
				ErrorLogger.Println(err)
			}
			booking = nil
			return
		}

		myUser := mapUsers[userName]
		myUser.UserBookings = AppendNodeToSlice(myUser.UserBookings, booking)
		myUser.UserBookings = SortBookingsByTime(myUser.UserBookings, len(myUser.UserBookings))
		myUser.UserBookings = SortBookingsByDate(myUser.UserBookings, len(myUser.UserBookings))
		mapUsers[userName] = myUser

		http.Redirect(w, r, "/booking_confirmed", http.StatusSeeOther)
		return
	}
	err = tpl.ExecuteTemplate(w, "newBooking.html", nil)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}

//BookingConfirmed is the handler for the user booking confirmed page.
func BookingConfirmed(w http.ResponseWriter, r *http.Request) {
	if GetUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	err := tpl.ExecuteTemplate(w, "bookingConfirmed.html", booking)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
	booking = nil
}

//ViewAllBookings is the handler for the user view all bookings page.
func ViewAllBookings(w http.ResponseWriter, r *http.Request) {
	if GetUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	userBookings := GetUser(r).UserBookings

	err := tpl.ExecuteTemplate(w, "viewAllBookings.html", userBookings)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}

//ChangeBookingPage is the handler for the user change bookings page.
func ChangeBookingPage(w http.ResponseWriter, r *http.Request) {
	//validate login
	if GetUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		var err error
		//get booking number
		bookingId := r.FormValue("bookingId")
		bookingId = StripHtmlRegex(bookingId)
		//iterate through slice to get bookingNode
		myUser := GetUser(r)
		booking, err = SearchId(myUser.UserBookings, bookingId)
		if err != nil {
			_, err := fmt.Fprintf(w, "there are no bookings with that Booking ID, go back to re-enter ID")
			if err != nil {
				ErrorLogger.Println(err)
			}
			return
		}
	}
	err := tpl.ExecuteTemplate(w, "changeBookingPage.html", booking)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}

//GetChanges is the handler for the user getChanges path
func GetChanges(w http.ResponseWriter, r *http.Request) {
	//myUser := GetUser(r)
	if GetUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost {
		//collect old car data
		oldCarArr := GetCarArr(booking.Car)
		oldDate := ConvertDate(booking.Date)
		oldTime := ConvertTime(booking.BookingTime)

		car := r.FormValue("cars")
		date := r.FormValue("date")
		date = StripHtmlRegex(date)
		if err := CheckDate(date); err != nil {
			_, err := fmt.Fprintf(w, "%v, go back to change date", err)
			if err != nil {
				ErrorLogger.Println(err)
			}
			return
		}
		bookingTime, _ := strconv.Atoi(r.FormValue("bookingTime"))
		pickUp := r.FormValue("pickUp")
		pickUp = StripHtmlRegex(pickUp)
		dropOff := r.FormValue("dropOff")
		dropOff = StripHtmlRegex(dropOff)
		contact, err := strconv.Atoi(r.FormValue("contact"))
		if err != nil {
			_, err := fmt.Fprintln(w, "Invalid contact number, go back to input new contact number")
			if err != nil {
				ErrorLogger.Println(err)
			}
			return
		}
		if !CheckContactLen(contact) {
			_, err := fmt.Fprintln(w, "Invalid contact number, go back to input new contact number")
			if err != nil {
				ErrorLogger.Println(err)
			}
			return
		}
		remarks := r.FormValue("remarks")
		remarks = StripHtmlRegex(remarks)

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
		newCarArr := GetCarArr(car)
		newDate := ConvertDate(date)
		newTime := ConvertTime(bookingTime)

		if newCarArr[newDate][newTime] != nil { //check for empty timeslot
			_, err := fmt.Fprintf(w, "Error: %v , go back to select a new slot", errors.New("there is already a booking at that time and date"))
			if err != nil {
				ErrorLogger.Println(err)
			}
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
		myUser := mapUsers[GetUser(r).Username]
		myUser.UserBookings = SortBookingsByTime(myUser.UserBookings, len(myUser.UserBookings))
		myUser.UserBookings = SortBookingsByDate(myUser.UserBookings, len(myUser.UserBookings))
		mapUsers[GetUser(r).Username] = myUser
		http.Redirect(w, r, "/print_changed_booking", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(w, "changeBooking.html", booking)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}

//PrintChangedBooking is the handler for the Print changed booking page.
func PrintChangedBooking(w http.ResponseWriter, r *http.Request) {
	if GetUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	err := tpl.ExecuteTemplate(w, "printChangedBooking.html", booking)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
	booking = nil
}

//DeleteBookingPage is the handler for the user delete booking page.
func DeleteBookingPage(w http.ResponseWriter, r *http.Request) {
	//validate login
	if GetUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		var err error
		//get booking number
		bookingId := r.FormValue("bookingId")
		bookingId = StripHtmlRegex(bookingId)
		//iterate through slice to get bookingNode
		myUser := GetUser(r)
		booking, err = SearchId(myUser.UserBookings, bookingId)
		if err != nil {
			_, err := fmt.Fprintf(w, "there are no bookings with that Booking ID, go back to re-enter ID")
			if err != nil {
				ErrorLogger.Println(err)
			}
			return
		}
	}
	err := tpl.ExecuteTemplate(w, "deleteBookingPage.html", booking)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}

//DeleteBooking is the handler for the delete booking path.
func DeleteBooking(w http.ResponseWriter, r *http.Request) {
	if GetUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	myUser := GetUser(r)
	DeleteFromCarsArr(booking)
	if err := DeleteBookingUserArr(myUser, booking); err != nil {
		ErrorLogger.Printf("error: %s", err)
	}
	if err := bookings.DeleteBookingNode(booking); err != nil {
		ErrorLogger.Printf("error: %s", err)
	}
	booking = nil

	err := tpl.ExecuteTemplate(w, "deleteConfirmed.html", nil)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}
