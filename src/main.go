package main

import (
	"log"
	"net/http"
	"encoding/json"
	"errors"
	"time"
	"strconv"
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
	// Enable CORS
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "*")
	
	// Accept a POST endpoint
	if (req.Method == "POST") {
		loginPostHandler(rw, req)
	} else if (req.Method == "OPTIONS") {
		return
	} else {
		http.Error(rw, errors.New("Method Not Supported").Error(), http.StatusMethodNotAllowed)
		return
	}
	
}

func validateCode (code string, offset int) (err error) {
	// If the code was generated within {offset} minutes before
	// 	or after now, it is a valid code.

	// If offset is 1, the current hour and minute must match
	//	the code that was generated.
	if (offset < 1) {
		return errors.New("Invalid Offset")
	}
	str_hour := code[0:2]
	str_min := code [2:4]

	int_hour, err := strconv.Atoi(str_hour)
	if err != nil {
		errors.New("Invalid Code")
	}
	int_min, err := strconv.Atoi(str_min)
	if err != nil {
		errors.New("Invalid Code")
	}

	now := time.Now()
	submitted := time.Date(now.Year(), now.Month(), now.Day(), int_hour, int_min, now.Second(), now.Nanosecond(), time.UTC)
	max_accept := now.Add(time.Minute * time.Duration(offset))
	min_accept := now.Add(time.Minute * -1 * time.Duration(offset))

	if ((submitted.Before(max_accept) && submitted.After(min_accept)) || submitted.Equal(now)) {
		return nil
	}

	return errors.New("Invalid Code")
}

func loginPostHandler (rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

    var loginRequest LoginRequest

	// If this was a production environment, we'd hash + salt the password
	// 	and store in a database so that these passwords aren't stored in
	//	plaintext. But for now we'll just pretend these are okay.

	validEmail := "c137@onecause.com"
	validPassword := "#th@nH@rm#y#r!$100%D0p#"
	submissionOffeset := 2

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

	err = validateCode(loginRequest.Code, submissionOffeset)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Success: Redirect to http://onecause.com to signify movement to the next page
	//	This actually just redirects the fetch request, not the user's browser
	//	Send an OK response instead and redirect client side
	// http.Redirect(rw, req, "http://onecause.com", http.StatusMovedPermanently)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/login", loginHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}