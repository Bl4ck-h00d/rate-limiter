package limiters

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"
)

type Bucket struct {
	timestamp time.Time
	tokens    int `default:"10"`
}

var userBuckets = make(map[string]*Bucket)

func (b *Bucket) consume(ipAddress string) bool {

	if b.tokens == 0 {
		return false
	}
	b.tokens--
	log.Printf("user %s, tokens = %d", ipAddress, b.tokens)
	return true
}

func (b *Bucket) refill() {
	if b.tokens == 10 {
		return
	}

	diff := time.Since(b.timestamp)
	if diff.Seconds() >= 1 {
		b.tokens = int(math.Min(10, float64(b.tokens+1)))
		b.timestamp = time.Now()
	}
}

func TokenBucketMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := r.RemoteAddr

		bucket, ok := userBuckets[ipAddress]
		if ok {
			if reqAccepted := bucket.consume(ipAddress); !reqAccepted {
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode("Too many requests !!")
				log.Printf("Too many request, %s", ipAddress)
				return
			}
		} else {
			// Add new user
			userBuckets[ipAddress] = &Bucket{timestamp: time.Now(), tokens: 10}
		}
		next.ServeHTTP(w, r)
	})
}

func RefillTokenBuckets() {
	for range time.Tick(time.Second * 1) {
		for _, bucket := range userBuckets {
			bucket.refill()
		}
	}
}
