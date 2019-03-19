package proxy

import (
	"testing"

	"github.com/ah8ad3/gateway/pkg/db"
)

func TestLoadServices(t *testing.T) {
	db.GenerateSecretKey()
	LoadServices(true, "./../../services.json")
	err := LoadServices(true, "./../../ser.json")
	if err.Message == "" {
		t.Fatal("File not found error")
	}
	LoadServices(false, "")

	if len(Services) == 0 {
		t.Fatal("Services cant load")
	}
}

func TestCheckServices(t *testing.T) {
	CheckServices(false)
	CheckServices(true)
}

func TestHealthCheck(t *testing.T) {
	defer HealthCheck()
}
