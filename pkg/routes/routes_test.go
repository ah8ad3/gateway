package routes

import (
	"github.com/ah8ad3/gateway/pkg/db"
	"github.com/ah8ad3/gateway/pkg/integrate"
	"github.com/ah8ad3/gateway/pkg/proxy"
	"github.com/ah8ad3/gateway/plugins"
	"net/http"
	"testing"
)

func TestBefore(t *testing.T) {
	db.GenerateSecretKey()
	proxy.LoadServices(true, "./../../services.json")
	plugins.RegisterPlugins()
	proxy.LoadServices(false, "")

	integrate.LoadIntegration("./../../integrates.json")
}

func TestV1(t *testing.T) {
	routes := V1()

	if routes == nil {
		t.Log("RouteV1 not correct")
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = req

	//rr := httptest.NewRecorder()
	//handler := http.HandlerFunc(HealthCheckHandler)
}

func TestV2(t *testing.T) {
	routes := V2()

	if routes == nil {
		t.Log("Routes V2 not correct")
	}
}
