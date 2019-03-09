package routes

import (
	"github.com/ah8ad3/gateway/pkg/db"
	"github.com/ah8ad3/gateway/pkg/proxy"
	"testing"
)

func TestV1(t *testing.T) {
	db.GenerateSecretKey()
	proxy.LoadServices(true, "./../../services.json")

	routes := V1()

	if routes == nil {
		t.Log("RouteV1 not correct")
	}
}

func TestV2(t *testing.T) {
	routes := V2()

	if routes == nil {
		t.Log("Routes V2 not correct")
	}
}