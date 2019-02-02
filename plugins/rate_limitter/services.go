package rate_limitter

import (
	"golang.org/x/time/rate"
)


var limiter = make(map[string] *rate.Limiter)


func SetLimmiter(server string, r rate.Limit, bursts int) *rate.Limiter {
	var limit *rate.Limiter
	mtx.Lock()
	limiter[server] = rate.NewLimiter(r, bursts)
	limit = limiter[server]
	mtx.Unlock()

	return limit

}

func DeleteLimiter(server string) {
	delete(limiter, server)
}
