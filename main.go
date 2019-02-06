package main

import (
	"log"
	"os"

	"github.com/ah8ad3/gateway/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Cmd cant start, look at cmd package")
		os.Exit(1)
	}
}
