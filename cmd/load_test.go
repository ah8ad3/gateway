package cmd

import (
	"context"
	"testing"
)

func TestNewLoadCmd(t *testing.T) {
	err := NewLoadCmd(context.Background()).Execute()

	if err != nil {
		t.Fatal(err.Error())
	}
}
