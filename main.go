package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", landingHandler)
	mux.HandleFunc("/rsvp", rsvpHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/admin/login", adminLoginHandler)
	mux.HandleFunc("/admin", adminHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf(
		"Wedding RSVP running on http://localhost:%s",
		port,
	)

	log.Fatal(
		http.ListenAndServe(":"+port, mux),
	)
}
