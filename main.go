package main

import (
	"os"

	"github.com/pitittatou/flightradarbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
