package cmd

import (
	"context"
	"github.com/ah8ad3/gateway/pkg"
	"github.com/spf13/cobra"
)

var test int

func init() {

}


// NewServerCmd creates a new version command
func NewServerCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "run",
		Short:   "Run gateway server",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			ip, port, route := "0.0.0.0", "3000", "v2"

			pkg.RUN(ip, port, route)
		},
	}
}
