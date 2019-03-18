package cmd

import (
	"context"
	"fmt"
	"github.com/ah8ad3/gateway/pkg/db"
	"github.com/ah8ad3/gateway/pkg/proxy"
	"os"

	"github.com/spf13/cobra"
)

var serLocation string

func init() {

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

			if os.Getenv("TEST") == "1" {
				fmt.Println("Run On Test Mode!")
				serLocation = "./../services.json"
			} else {
				serLocation = "services.json"
			}

			proxy.LoadServices(true, serLocation)
			fmt.Printf("Services loaded successfully \n")
		},
	}
}
