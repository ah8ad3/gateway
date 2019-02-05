package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// NewRootCmd creates a new instance of the root command
func NewRootCmd() *cobra.Command {
	ctx := context.Background()

	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "gateway is an API Gateway",
		Long: `
This is a lightweight API Gateway and Management Platform that enables you
to control who accesses your API, when they access it and how they access it.
API Gateway will also record detailed analytics on how your users are interacting
with your API and when things go wrong.`,
	}

	cmd.AddCommand(NewVersionCmd(ctx))
	cmd.AddCommand(NewServerCmd(ctx))

	return cmd
}
