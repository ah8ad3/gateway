package cmd

import (
	"context"
	"fmt"
	"github.com/ah8ad3/gateway/pkg/db"
	"github.com/spf13/cobra"
)

// NewLoadCmd creates a new version command
func NewSecretKeyCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "secret",
		Short:   "generate secret key",
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			db.GenerateSecretKey()
			fmt.Printf("Secret key generated \n")
		},
	}
}
