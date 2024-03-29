package proxy

import (
	"os"
	"testing"

	"github.com/ah8ad3/gateway/plugins"

	"github.com/ah8ad3/gateway/pkg/db"
)

func TestBefore(t *testing.T) {
	plugins.RegisterPlugins()
}

func TestLoadServices(t *testing.T) {
	db.GenerateSecretKey()
	LoadServices(true, "./../../services.json")
	err := LoadServices(true, "./../../ser.json")
	LoadServices(true, "./../../integrates.json")
	if err.Message == "" {
		t.Fatal("File not found error")
	}
	LoadServices(false, "")

	os.Remove("./../../db/proxy.bin")
	LoadServices(false, "")

	LoadServices(true, "./../../services.json")

	if len(Services) == 0 {
		t.Fatal("Services cant load")
	}
}

func TestCheckServices(t *testing.T) {
	CheckServices(false)
	CheckServices(true)
}

func TestHealthCheck(t *testing.T) {
	go HealthCheck()
}

func TestSyncPlugins(t *testing.T) {
	for _, val := range Services {
		SyncPlugins(val.Name)
	}
}

func TestAddPlugin(t *testing.T) {
	AddPlugin("service1", 1, "rateLimiter", nil)
	AddPlugin("service1", 1, "rate", nil)
	AddPlugin("service2", 1, "rateLimiter", nil)
	AddPlugin("service2", 1, "rateLimiter", map[string]interface{}{})
	AddPlugin("service2", 2, "rateLimiter", map[string]interface{}{})
}

func TestRemovePlugin(t *testing.T) {
	RemovePlugin("service1", 1, "rate")
	RemovePlugin("service1", 2, "rate")
	RemovePlugin("service1", 1, "rateLimiter")
}
