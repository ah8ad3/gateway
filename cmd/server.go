package cmd

import (
	"context"
	"os"

	"github.com/ah8ad3/gateway/pkg"
	"github.com/spf13/cobra"
)

var test int

func init() {
	if os.Getenv("TEST") == "1" {
		test = 1
	} else {
		test = 0
	}
}


// NewServerCmd creates a new version command
func NewServerCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "run",
		Short:   "Run gateway server",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			ip, port, route := "0.0.0.0", "3000", "v2"
			pkg.RUN(ip, port, route, test)
		},
	}
}
