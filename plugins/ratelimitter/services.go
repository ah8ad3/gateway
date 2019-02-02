package ratelimitter

import (
	"golang.org/x/time/rate"
)

var limiter = make(map[string]*rate.Limiter)

// SetLimmiter to set new limit for server
func SetLimmiter(server string, r rate.Limit, bursts int) *rate.Limiter {
	var limit *rate.Limiter
	mtx.Lock()
	limiter[server] = rate.NewLimiter(r, bursts)
	limit = limiter[server]
	mtx.Unlock()

	return limit

}

// DeleteLimiter for delete limitation of service
func DeleteLimiter(server string) {
	delete(limiter, server)
}
