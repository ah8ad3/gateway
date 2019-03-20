package routes

import (
	"testing"
)

func TestGetService(t *testing.T) {
	GetService("google.com", "/about", "")
	GetService("google.com", "/about", "?q=there")
	GetService("localhost:8000", "/about", "?q=there")
	findService("/foo")
	findService("foo")
}

func TestPostService(t *testing.T) {
	PostService("localhost:8000", "/test", nil)
}