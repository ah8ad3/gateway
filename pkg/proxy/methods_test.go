package proxy

import (
	"github.com/ah8ad3/gateway/pkg/db"
	"testing"
)

func TestLoadServices(t *testing.T) {
	db.GenerateSecretKey()
	LoadServices(true, "./../../services.json")

	t.Log("Test if can make Services full")

	if len(Services) == 0 {
		t.Errorf("Services cant load")
	}
}
