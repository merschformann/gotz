package main

import (
	"fmt"
	"os"

	"github.com/merschformann/gotz/core"
)

func main() {
	// Get configuration
	config, err := core.Load()
	if err != nil {
		fmt.Println("Error handling configuration:", err)
		os.Exit(1)
	}
	// Parse flags
	config, request, changed, err := core.ParseFlags(config)
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		os.Exit(1)
	}
	// Update config, if necessary
	if changed {
		err = config.Save()
		if err != nil {
			fmt.Println("Error saving configuration update:", err)
			os.Exit(1)
		}
	}
	// Plot time
	err = config.PlotTime(request)
	if err != nil {
		fmt.Println("Error plotting time:", err)
		os.Exit(1)
	}
}
