package cmd

import (
	"context"
	"testing"
)

func TestNewServerCmd(t *testing.T) {
	// i will test this in server file later
	err := NewServerCmd(context.Background()).Execute()

	if err != nil {
		t.Error("Server error in test")
	}

}
