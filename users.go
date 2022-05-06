package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type user struct {
	Username     string
	Password     []byte
	First        string
	Last         string
	userBookings []*bookingInfoNode
}

var mapUsers = make(map[string]user)
var mapSessions = make(map[string]string)

func init() {
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	mapUsers["khai"] = user{"khai", bPassword, "khai", "fazil", []*bookingInfoNode{}}
	mapUsers["joseph"] = user{"joseph", bPassword, "joseph", "seow", []*bookingInfoNode{}}
	mapUsers["doug"] = user{"doug", bPassword, "doug", "choo", []*bookingInfoNode{}}
	mapUsers["iza"] = user{"iza", bPassword, "iza", "zainuddin", []*bookingInfoNode{}}
	fmt.Println(mapUsers)
}

func getUser(r *http.Request) user {
	// get current session cookie
	sessionCookie, err := r.Cookie("sessionId")
	if err != nil { //if no cookie, just return empty user
		return user{}
	}

	//if cookie exists, continue on
	var myUser user
	if userName, ok := mapSessions[sessionCookie.Value]; ok {
		myUser = mapUsers[userName]
	}
	return myUser
}

func createUser(w http.ResponseWriter, r *http.Request) {
	//get inputs (username, password, firstname, lastname)
	username := r.FormValue("username")
	password := r.FormValue("password")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	//check if inputs are empty
	if username == "" || password == "" || firstName == "" || lastName == "" {
		http.Error(w, "One or more inputs are empty", http.StatusForbidden)
		return
	}
	//check if username is taken
	if _, ok := mapUsers[username]; ok {
		http.Error(w, "Username  already taken", http.StatusForbidden)
		return
	}
	//convert password to hash form
	bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	//if reach here, means that passed all prior checks
	//Save data to struct & mapUsers
	newUser := user{
		Username: username,
		Password: bPassword,
		First:    firstName,
		Last:     lastName,
	}
	mapUsers[username] = newUser
	//generate & set sessionID with function
	setSessionIDCookie(w, username)
}

func setSessionIDCookie(w http.ResponseWriter, username string) {
	//generate new UUID
	id := uuid.NewV4()
	//create new cookie with name and UUID
	sessionCookie := &http.Cookie{
		Name:  "sessionId",
		Value: id.String(),
	}
	//set cookie
	http.SetCookie(w, sessionCookie)
	//mapSession
	mapSessions[sessionCookie.Value] = username
}
