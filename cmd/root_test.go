package cmd

import (
	"testing"
)

func TestNewRootCmd(t *testing.T) {
	if err := NewRootCmd().Execute(); err != nil {
		t.Fatal(err.Error())
	}
}
