package routes

import (
	"testing"
)

func TestV1(t *testing.T) {
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