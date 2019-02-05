package admin

import (
	"encoding/json"
	"github.com/ah8ad3/gateway/pkg/proxy"
	"net/http"
)

func getServices() []byte {
	var jData []byte
	jData, _ = json.Marshal(proxy.Services)

	return jData
}

// Welcome just an sample welcome
func Welcome(w http.ResponseWriter, r *http.Request) {
	_ = r
	_, _ = w.Write([]byte("Welcome to gateway"))
}

// GETServices get all services in admin mode
func GETServices(w http.ResponseWriter, r *http.Request) {
	_ = r
	_, _ = w.Write(getServices())
}
