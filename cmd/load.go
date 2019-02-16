package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ah8ad3/gateway/pkg/db"
	"github.com/ah8ad3/gateway/pkg/proxy"

	"github.com/spf13/cobra"
)

var serLocation string

func init() {
	if os.Getenv("TEST") == "1" {
		serLocation = "./../services.json"
	} else {
		serLocation = "services.json"
	}
}

// NewLoadCmd creates a new version command
func NewLoadCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "load",
		Short:   "load services from services.json",
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			str, _ := db.LoadSecretKey()
			db.SecretKey = str
			proxy.LoadServices(true, serLocation)
			fmt.Printf("Services loaded successfully \n")
		},
	}
}
