package core_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/merschformann/gotz/core"
)

// update indicates whether to update the golden files instead of using them to
// compare the output.
var update = flag.Bool("update", false, "Update golden files")

// readExpectation reads the golden file and returns its content.
func readExpectation(goldenFile string) (string, error) {
	// Read expected output string from file
	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		return "", err
	}
	return string(expected), nil
}

func TestPlotStatic(t *testing.T) {
	// Define expected output
	goldenFile := "testdata/plot_static.golden"
	expected, err := readExpectation(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	// Get configuration
	config := core.DefaultConfig()

	// Create test plotter, collect output in stringbuilder for comparison
	sb := strings.Builder{}
	plotter := core.Plotter{
		Now:           true,
		TerminalWidth: 80,
		PlotLine: func(t core.ContextType, line ...interface{}) {
			_ = t
			sb.WriteString(fmt.Sprint(line...) + "\n")
		},
		PlotString: func(t core.ContextType, msg string) {
			_ = t
			sb.WriteString(msg)
		},
	}

	// Specify time
	loc, _ := time.LoadLocation("Europe/Berlin")
	testTime := time.Date(1985, 8, 24, 16, 0, 0, 0, loc)

	// Plot time
	err = core.PlotTime(plotter, config, testTime)
	if err != nil {
		t.Errorf("error plotting time: %s", err)
	}

	// Collect output
	output := sb.String()

	if *update {
		// Update golden file
		err = ioutil.WriteFile(goldenFile, []byte(output), 0644)
		if err != nil {
			t.Errorf("error updating golden file: %s", err)
		}
	} else {
		// Compare output with golden file
		if output != expected {
			t.Errorf("expected output:\n%s\nbut got:\n%s", expected, output)
		}
	}
}
