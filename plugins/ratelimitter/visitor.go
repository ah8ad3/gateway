package ratelimitter

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ah8ad3/gateway/plugins/ip"
	"golang.org/x/time/rate"
)

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
	allow    bool
}

// Change the the map to hold values of the type visitor.
var visitors = make(map[string]*visitor)
var mtx sync.Mutex

func addVisitor(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(5, 5)
	mtx.Lock()
	// Include the current time when creating a new visitor.
	visitors[ip] = &visitor{limiter, time.Now(), true}
	mtx.Unlock()
	return limiter
}

func getVisitor(ip string) *rate.Limiter {
	mtx.Lock()
	v, exists := visitors[ip]
	if !exists {
		mtx.Unlock()
		return addVisitor(ip)
	}

	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	mtx.Unlock()
	return v.limiter
}

// CleanupVisitors Every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
func CleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mtx.Lock()
		for id, v := range visitors {
			if time.Now().Sub(v.lastSeen) > 3*time.Minute {
				delete(visitors, id)
			}
		}
		mtx.Unlock()
	}
}

// LimitMiddleware to check for the too many request every too many requests
// will ban for 1 minutes
func LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := getVisitor(r.RemoteAddr)
		if limiter.Allow() == false {
			splitRoute := strings.Split(r.URL.Path, "/")
			// extract server path from url
			path := splitRoute[1]
			ip.AddBlockList(r.RemoteAddr, path, time.Duration(time.Minute*1), false)
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func TestMiddle(a time.Duration) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := getVisitor(r.RemoteAddr)
			if limiter.Allow() == false {
				splitRoute := strings.Split(r.URL.Path, "/")
				// extract server path from url
				path := splitRoute[1]
				ip.AddBlockList(r.RemoteAddr, path, time.Duration(time.Second*a), false)
				http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
