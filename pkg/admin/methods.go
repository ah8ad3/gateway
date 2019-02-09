package admin

import (
	"encoding/json"
	"net/http"

	"github.com/ah8ad3/gateway/pkg/proxy"
)

func getServices() []byte {
	var jData []byte
	jData, _ = json.Marshal(proxy.Services)

	return jData
}

// Welcome just an sample welcome
func Welcome(w http.ResponseWriter, r *http.Request) {
	_ = r

	//str, _ := proxy.AddPlugin("service1", "rateLimiter", nil)
	_, _ = w.Write([]byte("Welcome To Gateway"))
}

// GETServices get all services in admin mode
func GETServices(w http.ResponseWriter, r *http.Request) {
	_ = r
	_, _ = w.Write(getServices())
}
