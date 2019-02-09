package cmd

import (
	"context"
	"fmt"
	"github.com/ah8ad3/gateway/pkg/db"
	"github.com/ah8ad3/gateway/pkg/proxy"

	"github.com/spf13/cobra"
)

// NewLoadCmd creates a new version command
func NewLoadCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "load",
		Short:   "load services from services.json",
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			str, _ :=db.LoadSecretKey()
			db.SecretKey = str
			proxy.LoadServices(true)
			fmt.Printf("Services loaded successfully \n")
		},
	}
}
