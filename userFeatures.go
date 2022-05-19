package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func index(w http.ResponseWriter, r *http.Request) {
	currentUser := getUser(r)
	err := tpl.ExecuteTemplate(w, "index.html", currentUser)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
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
			_, err := fmt.Fprintln(w, "A car was not selected, go back to select car")
			if err != nil {
				ErrorLogger.Println(err)
			}
			return
		}
		date := r.FormValue("date")
		date = StripHtmlRegex(date)
		if err := checkDate(date); err != nil {
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
		userName := getUser(r).Username
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
		if !checkContactLen(contact) {
			_, err := fmt.Fprintln(w, "Invalid contact number, go back to input new contact number")
			if err != nil {
				ErrorLogger.Println(err)
			}
			return
		}
		remarks := r.FormValue("remarks")
		remarks = StripHtmlRegex(remarks)

		booking, err = bookings.makeNewBooking(car, date, bookingTime, userName, pickUp, dropOff, contact, remarks)
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
		myUser.UserBookings = sortBookingsByTime(myUser.UserBookings, len(myUser.UserBookings))
		myUser.UserBookings = sortBookingsByDate(myUser.UserBookings, len(myUser.UserBookings))
		mapUsers[userName] = myUser

		http.Redirect(w, r, "/booking_confirmed", http.StatusSeeOther)
		return
	}
	err = tpl.ExecuteTemplate(w, "newBooking.html", nil)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}

func bookingConfirmed(w http.ResponseWriter, r *http.Request) {
	if getUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	err := tpl.ExecuteTemplate(w, "bookingConfirmed.html", booking)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
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
		ErrorLogger.Panicf("error executing template: %v", err)
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
		bookingId = StripHtmlRegex(bookingId)
		//iterate through slice to get bookingNode
		myUser := getUser(r)
		booking, err = searchId(myUser.UserBookings, bookingId)
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

func getChanges(w http.ResponseWriter, r *http.Request) {
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
		date = StripHtmlRegex(date)
		if err := checkDate(date); err != nil {
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
		if !checkContactLen(contact) {
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
		newCarArr := getCarArr(car)
		newDate := convertDate(date)
		newTime := convertTime(bookingTime)

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
		myUser := mapUsers[getUser(r).Username]
		myUser.UserBookings = sortBookingsByTime(myUser.UserBookings, len(myUser.UserBookings))
		myUser.UserBookings = sortBookingsByDate(myUser.UserBookings, len(myUser.UserBookings))
		mapUsers[getUser(r).Username] = myUser
		http.Redirect(w, r, "/print_changed_booking", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(w, "changeBooking.html", booking)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}

func printChangedBooking(w http.ResponseWriter, r *http.Request) {
	if getUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	err := tpl.ExecuteTemplate(w, "printChangedBooking.html", booking)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
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
		bookingId = StripHtmlRegex(bookingId)
		//iterate through slice to get bookingNode
		myUser := getUser(r)
		booking, err = searchId(myUser.UserBookings, bookingId)
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

func deleteBooking(w http.ResponseWriter, r *http.Request) {
	if getUser(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	myUser := getUser(r)
	deleteFromCarsArr(booking)
	if err := deleteBookingUserArr(myUser, booking); err != nil {
		ErrorLogger.Printf("error: %s", err)
	}
	if err := bookings.deleteBookingNode(booking); err != nil {
		ErrorLogger.Printf("error: %s", err)
	}
	booking = nil

	err := tpl.ExecuteTemplate(w, "deleteConfirmed.html", nil)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}
