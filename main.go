package main

import (
	"fmt"
	"os"
	"path"

	"github.com/merschformann/gotz/core"
)

func main() {
	// TODO: remove this
	// Delete configuration file if it exists
	home, _ := os.UserHomeDir()
	configFile := path.Join(home, ".gotz.config.json")
	if _, err := os.Stat(configFile); err == nil {
		os.Remove(configFile)
	}
	// Get configuration
	config, err := core.Load()
	if err != nil {
		fmt.Println("Error handling configuration:", err)
		os.Exit(1)
	}
	// Parse flags
	config, request, err := core.ParseFlags(config)
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		os.Exit(1)
	}
	// Plot time
	err = config.PlotTime(request)
	if err != nil {
		fmt.Println("Error plotting time:", err)
		os.Exit(1)
	}
}
