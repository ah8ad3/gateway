package rest

import (
	"encoding/json"
	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/ah8ad3/gateway/pkg/proxy"
	"net/http"
	"time"
)


func getServices() []byte {
	var jData []byte
	jData, _ = json.Marshal(proxy.Services)

	return jData
}

// GETServices get all services in admin mode
func GETServices(w http.ResponseWriter, r *http.Request) {
	_ = r
	_, _ = w.Write(getServices())
}

func PostService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var service proxy.Service
	err := decoder.Decode(&service)
	if err != nil {
		logger.SetUserLog(logger.UserLog{Time: time.Now(), IP: r.RemoteAddr, RequestURL: r.URL.Path, Log: logger.Log{
			Event: "critical", Description: err.Error()}})
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error": "Input not in Service struct, check you're input'"}`))
		return
	}
	proxy.Services = append(proxy.Services, service)

	_, _ = w.Write([]byte(`{"status": "ok"}`))
	return
}
