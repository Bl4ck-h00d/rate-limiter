package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	// Router
	mux := http.NewServeMux()

	mux.HandleFunc("/limited", handleLimited)
	mux.HandleFunc("/unlimited", handleDefault)

	log.Print("Listening...")
	http.ListenAndServe(":3000", mux)
}

func handleLimited(w http.ResponseWriter, req *http.Request) {
	// TODO: rate limiter
	w.Header().Set("Status", "200")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Limited, don't over use me!")
}

func handleDefault(w http.ResponseWriter, req *http.Request) {
	// TODO: rate limiter
	w.Header().Set("Status", "200")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Unlimited! Let's Go!")
}
