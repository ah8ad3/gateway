package main

import (
	"github.com/ah8ad3/gateway/cmd"
	"log"
	"os"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Cmd cant start, look at cmd package")
		os.Exit(1)
	}
}
