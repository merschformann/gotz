package main

import (
	"fmt"

	"github.com/merschformann/gotz/core"
)

func main() {
	// Print all colors using all basic terminal codes
	for k, v := range core.NamedColors {
		fmt.Println(v + "" + k + core.ColorReset)
	}
}
