package limoBookingApp

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type admin struct {
	Username string
	Password []byte
}

var mapAdmins = make(map[string]admin)

func init() {
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	mapAdmins["admin"] = admin{Username: "admin", Password: bPassword}
}

//GetAdmin gets the Admin user from the cookie sessionID and returns the admin user
func GetAdmin(r *http.Request) admin {
	// get current session cookie
	sessionCookie, err := r.Cookie("sessionId")
	if err != nil { //if no cookie, just return empty user
		return admin{}
	}

	//if cookie exists, continue on
	var myAdmin admin
	if userName, ok := mapSessions[sessionCookie.Value]; ok {
		myAdmin = mapAdmins[userName]
	}
	return myAdmin
}
