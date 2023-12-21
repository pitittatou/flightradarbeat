package main

import (
	"os"

	"github.com/pitittatou/flightradarbeat/cmd"

	_ "github.com/pitittatou/flightradarbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
