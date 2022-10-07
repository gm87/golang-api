package main

import (
	// "io"
	"log"
	"net/http"
	"encoding/json"
)

type LoginRequest struct {
	email string
	password string
	// Woah, why is code a string? code is a four digit number that contains
	// 	the two digit hour and minute the request was generated. What happens
	//	when the hour has a leading zero, ex: 0908 for 9:08 AM? If made a
	//	number, the leading zero would be removed. So, let's make it a string
	code string
}

func setStatus (rw http.ResponseWriter, req *http.Request, status int) {
	switch status {
	case 200:
		rw.WriteHeader(http.StatusOK)
	case 301:
		rw.WriteHeader(http.StatusMovedPermanently)
	case 400:
		rw.WriteHeader(http.StatusBadRequest)
	case 405:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func loginHandler (rw http.ResponseWriter, req *http.Request) {
	if (req.Method == "POST") {
		loginPostHandler(rw, req)
	} else {
		setStatus(rw, req, 405)
	}
}

func loginPostHandler (rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var loginRequest LoginRequest
	err := decoder.Decode(&loginRequest)
	if err != nil {
		log.Panic(err)
	}
	log.Println(loginRequest)
}

func main() {
	http.HandleFunc("/auth/login", loginHandler)

	log.Fatal(http.ListenAndServeTLS(":8080", "/app/cert.pem", "/app/key.pem", nil))
}