package rest

import (
	"encoding/json"
	"github.com/ah8ad3/gateway/pkg/routes"
	"net/http"
)

func GetServices(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(routes.Service)
	_, _ = w.Write(data)
}
