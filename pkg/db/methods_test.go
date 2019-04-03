package db

import (
	"os"
	"testing"
)

func TestGenerateSecretKey(t *testing.T) {
	if err := os.RemoveAll("./../../db/secret.bin"); err != nil {
		t.Fatal(err.Error())
	}
	LoadSecretKey()

	GenerateSecretKey()
	GenerateSecretKey()
}

func TestGetProxies(t *testing.T) {
	GetProxies()
}
