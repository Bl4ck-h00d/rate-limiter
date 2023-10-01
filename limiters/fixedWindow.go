package limiters

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

var userRequestHistory = make(map[string]int)
var mutex sync.Mutex

func consume(ipAddress, timestamp string) bool {
	userId := ipAddress + "@" + timestamp
	mutex.Lock()
	if userRequestHistory[userId] == 10 {
		mutex.Unlock()
		return false
	}
	userRequestHistory[userId] = int(math.Min(10, float64(userRequestHistory[userId]+1)))
	log.Printf("user = %s, counter = %d", userId, userRequestHistory[userId])
	mutex.Unlock()
	return true
}

func discardStaleUserIds() {
	for k, _ := range userRequestHistory {
		currentTime := time.Now()
		timestamp, error := time.ParseInLocation("02-Jan-2006 15:04:00", strings.Split(k, "@")[1], currentTime.Location())
		if error != nil {
			log.Print(error)
			return
		}

		diff := currentTime.Sub(timestamp)
		log.Print(diff.Minutes())
		if diff.Minutes() >= 1 {
			mutex.Lock()
			delete(userRequestHistory, k)
			mutex.Unlock()
			log.Printf("Discarded user = %s", k)
		}
	}
}

func Ticker() {
	for range time.Tick(time.Second * 30) {
		discardStaleUserIds()
	}
}

func FixedWindowMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := r.RemoteAddr
		if reqAccepted := consume(ipAddress, time.Now().Format("02-Jan-2006 15:04:00")); !reqAccepted {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode("Too many requests !!")
			log.Printf("Too many request, %s", ipAddress)
			return
		}
		next.ServeHTTP(w, r)
	})
}
