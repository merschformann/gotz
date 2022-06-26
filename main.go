package main

import (
	"fmt"
	"os"

	"github.com/merschformann/gotz/core"
)

func main() {
	// Get configuration
	config, err := core.Load()
	// If there was an error loading the config, offer the user the option to
	// reset it (or simply exit).
	if err != nil {
		fmt.Println("error loading configuration:", err)
		// Ask the user if they want to reset the config
		if ok, in_err := core.AskUser("Reset configuration?"); in_err != nil {
			fmt.Println("error asking user:", in_err)
			os.Exit(1)
		} else if ok {
			// Reset config
			config, in_err = core.SaveDefault()
			if in_err != nil {
				fmt.Println("error resetting configuration:", in_err)
				os.Exit(1)
			}
		} else {
			// Exit
			os.Exit(0)
		}
	}
	// Parse flags
	config, rt, changed, err := core.ParseFlags(config, Version)
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
	err = core.Plot(config, rt)
	if err != nil {
		fmt.Println("error plotting time:", err)
		os.Exit(1)
	}
}
