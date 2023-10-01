package main

import (
	"encoding/json"
	"log"
	"net/http"
	"rate-limiter/limiters"
)

func main() {
	// Router
	mux := http.NewServeMux()
	go limiters.RefillTokenBuckets()
	go limiters.Ticker()

	mux.Handle("/limited", limiters.TokenBucketMiddleware(http.HandlerFunc(handleLimited)))
	mux.Handle("/limited-window", limiters.FixedWindowMiddleware(http.HandlerFunc(handleLimited)))

	mux.HandleFunc("/unlimited", handleDefault)

	log.Print("Listening...")
	http.ListenAndServe(":3000", mux)
}

func handleLimited(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Limited use only")
}

func handleDefault(w http.ResponseWriter, req *http.Request) {
	// TODO: rate limiter
	w.Header().Set("Status", "200")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Unlimited! Let's Go!")
}
