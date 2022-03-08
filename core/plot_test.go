package core_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/merschformann/gotz/core"
)

type ConsolePrinter struct{}

func (cp *ConsolePrinter) Print(value string) {
	fmt.Printf("this is value: %s", value)
}

func TestPlotStatic(t *testing.T) {
	// Get configuration
	config := core.DefaultConfig()

	// Create test plotter
	plotter := core.Plotter{}

	// Specify time
	loc, _ := time.LoadLocation("Europe/Berlin")
	testTime := time.Date(1985, 8, 24, 16, 0, 0, 0, loc)

	// Plot time
	err := core.PlotTime(plotter, config, testTime)
	if err != nil {
		t.Errorf("error plotting time: %s", err)
	}
}
