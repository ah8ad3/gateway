package cmd

import (
	"context"
	"os"
	"testing"
)

func TestNewLoadCmd(t *testing.T) {
	err := NewLoadCmd(context.Background()).Execute()

	if err != nil {
		t.Fatal(err.Error())
	}

	os.Setenv("TEST", "0")

	err = NewLoadCmd(context.Background()).Execute()

	if err != nil {
		t.Fatal(err.Error())
	}
	os.Setenv("TEST", "1")

}
