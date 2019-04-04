package db

import (
	"encoding/json"
	"os"
	"testing"
)

type sample struct {
	name string
}

var samples []sample

func TestBefore(t *testing.T) {
	samples = append(samples, sample{name: "test"})
	samples = append(samples, sample{name: "test1"})
	samples = append(samples, sample{name: "test2"})
}

func TestGenerateSecretKey(t *testing.T) {
	if err := os.RemoveAll("./../../db/"); err != nil {
		t.Fatal(err.Error())
	}
	LoadSecretKey()
	GenerateSecretKey()

	LoadSecretKey()
	GenerateSecretKey()

	if nl := GetProxies(); nl != nil {
		t.Fatal("nl must be nil")
	}
}

func TestInsertProxy(t *testing.T) {
	jData, _ := json.Marshal(samples)
	InsertProxy(jData)
	InsertProxy(jData)

	SecretKey = ""
	InsertProxy(jData)
}

func TestGetProxies(t *testing.T) {
	GetProxies()
	LoadSecretKey()
	GetProxies()
}
