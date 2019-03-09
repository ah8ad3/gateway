package proxy

import (
	"testing"

	"github.com/ah8ad3/gateway/pkg/db"
)

func TestLoadServices(t *testing.T) {
	db.GenerateSecretKey()
	LoadServices(true, "./../../services.json")

	t.Log("Test if can make Services full")

	if len(Services) == 0 {
		t.Errorf("Services cant load")
	}
}

func TestCheckServices(t *testing.T) {
	CheckServices(false)
	CheckServices(true)
}
