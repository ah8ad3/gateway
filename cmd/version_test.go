package cmd

import (
	"context"
	"testing"
)

func TestNewVersionCmd(t *testing.T) {
	err := NewVersionCmd(context.Background()).Execute()

	if err != nil {
		t.Fatal(err.Error())
	}
}
