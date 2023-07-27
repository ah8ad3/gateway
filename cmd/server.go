package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ah8ad3/gateway/pkg"
	"github.com/spf13/cobra"
)

func init() {

}

const (
	ServerStop    string = "ServerStop"
	ServerStart   string = "ServerStart"
	ServerRestart string = "ServerRestart"
)

// NewServerCmd creates a new version command
func NewServerCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "run",
		Short:   "Run gateway server",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			serverCommand := make(chan string)
			var sw sync.WaitGroup
			sw.Add(1)
			go run(serverCommand)
			serverCommand <- ServerStart
			time.Sleep(7 * time.Second)
			// serverCommand <- ServerStart
			sw.Wait()
		},
	}
}

func run(serverCommand chan string) {
	ip, port, route := "0.0.0.0", "3000", "v2"
	stopSignal := make(chan bool)

	for {
		select {
		case command := <-serverCommand:
			switch command {
			case ServerStart:
				fmt.Println("StartCommand")
				go pkg.RUN(ip, port, route, stopSignal)
			case ServerStop:
				stopSignal <- true
			}

			// case ServerStop <- serverCommand:
			// 	stopSignal <- true

		}
	}

}
