package db

import (
	"os"
	"testing"
)

func TestGenerateSecretKey(t *testing.T) {
	LoadSecretKey()

	GenerateSecretKey()

	if err := os.RemoveAll("./../../db/secret.bin"); err != nil {
		t.Fatal(err.Error())
	}
	LoadSecretKey()
	GenerateSecretKey()
}

func TestGetProxies(t *testing.T) {
	//GetProxies()
}
