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

		if !IsAlphabetic(username) {
			err := errors.New("username includes invalid characters")
			AdminLogger.Printf("LOGIN UNSUCCESSFUL: '%s', %v", username, err)
			http.Error(w, "Username and/or password is not valid", http.StatusUnauthorized)
			return
		}

		if !IsAlphanumeric(password) {
			err := errors.New("password includes invalid characters")
			AdminLogger.Printf("LOGIN UNSUCCESSFUL: %v tried to login, %v", username, err)
			http.Error(w, "Username and/or password is not valid", http.StatusUnauthorized)
			return
		}

		myAdminUser, ok := mapAdmins[username]
		if !ok {
			err := errors.New("username not found")
			AdminLogger.Printf("LOGIN UNSUCCESSFUL: '%s', %v", username, err)
			http.Error(w, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}
		err := bcrypt.CompareHashAndPassword(myAdminUser.Password, []byte(password))
		if err != nil {
			err := errors.New("wrong password")
			AdminLogger.Printf("LOGIN UNSUCCESSFUL: %v tried to login, %v", username, err)
			http.Error(w, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}

		setSessionIDCookie(w, username)
		AdminLogger.Printf("LOGIN SUCCESSFUL: %s logged in", username)
		http.Redirect(w, r, "/admin_index", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(w, "adminLogin.html", nil)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}

func adminIndex(w http.ResponseWriter, r *http.Request) {
	currentAdmin := getAdmin(r)
	err := tpl.ExecuteTemplate(w, "adminIndex.html", currentAdmin)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
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
			err := bookings.deleteBookingNode(v)
			if err != nil {
				ErrorLogger.Println(err)
			}
		}

		delete(mapUsers, username)
	}
	err := tpl.ExecuteTemplate(w, "deleteUsers.html", mapUsers)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
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
		ErrorLogger.Panicf("error executing template: %v", err)
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
			ErrorLogger.Printf("searchId error: %v", err)
			_, err := fmt.Fprintf(w, "searchId error: %v", err)
			if err != nil {
				ErrorLogger.Println(err)
			}
		} else {
			myUser := mapUsers[toDelete.UserName]
			wg.Add(1)
			go deleteFromCarsArr(toDelete)
			if err := deleteBookingUserArr(myUser, toDelete); err != nil {
				ErrorLogger.Println(err)
			}
			if err := bookings.deleteBookingNode(toDelete); err != nil {
				ErrorLogger.Println(err)
			}
			wg.Done()
			wg.Wait()
			http.Redirect(w, r, "/admin_delete_booking_confirmed", http.StatusSeeOther)
			return
		}
	}
	err := tpl.ExecuteTemplate(w, "adminViewBookings.html", toDisplay)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}

func adminDeleteBookingConfirmed(w http.ResponseWriter, r *http.Request) {
	if getAdmin(r).Username == "" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	err := tpl.ExecuteTemplate(w, "adminDeleteConfirmed.html", nil)
	if err != nil {
		ErrorLogger.Panicf("error executing template: %v", err)
	}
}
