package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	helloHandler := func (w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello world!\n")
	}
	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServeTLS(":8080", "/app/cert.pem", "/app/key.pem", nil))
}