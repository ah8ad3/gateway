package cmd

import (
	"context"
	"github.com/ah8ad3/gateway/pkg"
	"github.com/spf13/cobra"
)

// NewServerCmd creates a new version command
func NewServerCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "run",
		Short:   "Run gateway server",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			pkg.RUN("", "")
		},
	}
}
