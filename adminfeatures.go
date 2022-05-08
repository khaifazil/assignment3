package main

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

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
	err := tpl.ExecuteTemplate(w, "adminLogin.gohtml", nil)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func adminIndex(w http.ResponseWriter, r *http.Request) {
	currentAdmin := getAdmin(r)
	err := tpl.ExecuteTemplate(w, "adminIndex.gohtml", currentAdmin)
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
	err := tpl.ExecuteTemplate(w, "deleteUsers.gohtml", mapUsers)
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
	err := tpl.ExecuteTemplate(w, "deleteSessions.gohtml", mapSessions)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func adminViewDeleteBookings(w http.ResponseWriter, r *http.Request) {

	if getAdmin(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	toDisplay, _ := bookings.appendAllToSlice()

	if r.Method == http.MethodPost {
		id := r.FormValue("bookingId")
		//fmt.Println(id)
		if toDelete, err := searchId(toDisplay, id); err != nil {
			fmt.Fprintf(w, "searchId error: %v", err)
		} else {
			myUser := mapUsers[toDelete.UserName]
			wg.Add(1)
			go deleteFromCarsArr(toDelete)
			if err := deleteBookingUserArr(myUser, toDelete); err != nil {
				fmt.Errorf("error: %s", err)
			}
			if err := bookings.deleteBookingNode(toDelete); err != nil {
				fmt.Errorf("error: %s", err)
			}
			wg.Done()
			wg.Wait()
			http.Redirect(w, r, "/admin_delete_booking_confirmed", http.StatusSeeOther)
			return
		}
	}
	err := tpl.ExecuteTemplate(w, "adminViewBookings.gohtml", toDisplay)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}

func adminDeleteBookingConfirmed(w http.ResponseWriter, r *http.Request) {
	if getAdmin(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	err := tpl.ExecuteTemplate(w, "adminDeleteConfirmed.gohtml", nil)
	if err != nil {
		panic(errors.New("error executing template"))
	}
}
