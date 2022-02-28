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
		fmt.Println("error handling configuration:", err)
		os.Exit(1)
	}
	// Parse flags
	config, request, changed, err := core.ParseFlags(config)
	if err != nil {
		fmt.Println("error parsing flags:", err)
		os.Exit(1)
	}
	// Update config, if necessary
	if changed {
		err = config.Save()
		if err != nil {
			fmt.Println("error saving configuration update:", err)
			os.Exit(1)
		}
	}
	// Plot time
	err = config.Plot(request)
	if err != nil {
		fmt.Println("error plotting time:", err)
		os.Exit(1)
	}
}
