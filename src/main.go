package main

import (
	// "io"
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"errors"
)

// Woah, why is code a string? code is a four digit number that contains
// 	the two digit hour and minute the request was generated. What happens
//	when the hour has a leading zero, ex: 0908 for 9:08 AM? If made a
//	number, the leading zero would be removed. So, let's make it a string
type LoginRequest struct {
	Email string
	Password string
	Code string
}

func loginHandler (rw http.ResponseWriter, req *http.Request) {
	// Accept a POST endpoint
	if (req.Method != "POST") {
		http.Error(rw, errors.New("Method Not Supported").Error(), http.StatusMethodNotAllowed)
		return
	}
	loginPostHandler(rw, req)
}

func loginPostHandler (rw http.ResponseWriter, req *http.Request) {
    var loginRequest LoginRequest
	
	// If this was a production environment, we'd hash + salt the password
	// 	and store in a database so that these passwords aren't stored in
	//	plaintext. But for now we'll just pretend these are okay.

	validEmail := "c137@onecause.com"
	validPassword := "#th@nH@rm#y#r!$100%D0p#"
	// submissionOffeset := 0

    err := json.NewDecoder(req.Body).Decode(&loginRequest)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

	// Accept a POST endpoint with a JSON payload with username, password, and token
	if (len(loginRequest.Email) == 0) {
		http.Error(rw, errors.New("Required field: Email").Error(), http.StatusBadRequest)
		return
	}

	if (len(loginRequest.Password) == 0) {
		http.Error(rw, errors.New("Required field: Password").Error(), http.StatusBadRequest)
		return
	}

	if (len(loginRequest.Code) == 0) {
		http.Error(rw, errors.New("Required field: Code").Error(), http.StatusBadRequest)
		return
	}

	if (loginRequest.Email != validEmail || loginRequest.Password != validPassword) {
		http.Error(rw, errors.New("Invalid username/password combination").Error(), http.StatusBadRequest)
		return
	}

	// Success: Redirect to http://onecause.com to signify movement to the next page
	http.Redirect(rw, req, "http://onecause.com", http.StatusMovedPermanently)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/login", loginHandler)

	log.Fatal(http.ListenAndServeTLS(":8080", "/app/cert.pem", "/app/key.pem", mux))
}