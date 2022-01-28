package main

import (
	"fmt"
	"os"

	"github.com/merschformann/gotz/core"
)

func main() {
	// TODO: remove this
	// Delete configuration file if it exists
	if _, err := os.Stat(".gotz.config"); err == nil {
		os.Remove(".gotz.config")
	}
	// Get configuration
	config, err := core.Load()
	if err != nil {
		fmt.Println("Error handling configuration:", err)
		os.Exit(1)
	}
	// Plot time
	err = config.PlotTime()
	if err != nil {
		fmt.Println("Error plotting time:", err)
		os.Exit(1)
	}
}
