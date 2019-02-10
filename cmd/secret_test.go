package cmd

import (
	"context"
	"testing"
)

func TestNewSecretKeyCmd(t *testing.T) {
	a := NewSecretKeyCmd(context.Background())
	err := a.Execute()
	if err != nil {
		t.Fatal(err.Error())
	}
}
