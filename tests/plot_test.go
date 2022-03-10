package core_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
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

func TestTableStatic(t *testing.T) {
	// Get all test configurations
	testConfigurations, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Fatal(err)
	}

	// Set local time to UTC for reproducibility
	time.Local = time.UTC

	// Specify test time
	loc, _ := time.LoadLocation("Europe/Berlin")
	testTime := time.Date(1985, 8, 24, 16, 0, 0, 0, loc)

	// Setup plotter (collect output in stringbuilder for comparison)
	sb := strings.Builder{}
	plotter := core.Plotter{
		Now:           true,
		TerminalWidth: 72,
		PlotLine: func(t core.ContextType, line ...interface{}) {
			_ = t
			sb.WriteString(fmt.Sprint(line...) + "\n")
		},
		PlotString: func(t core.ContextType, msg string) {
			_ = t
			sb.WriteString(msg)
		},
	}

	// Run all tests
	for _, configFile := range testConfigurations {
		t.Run(strings.Replace(configFile, ".json", "", -1), func(t *testing.T) {
			// Reset stringbuilder
			sb.Reset()
			// Read configuration file
			var config core.Config
			data, err := ioutil.ReadFile(configFile)
			if err != nil {
				t.Fatal(err)
			}
			// Unmarshal
			err = json.Unmarshal(data, &config)
			if err != nil {
				t.Fatal(err)
			}
			// Get expected output
			goldenFile := strings.Replace(configFile, ".json", ".golden", -1)
			expected, err := readExpectation(goldenFile)
			if err != nil && !*update {
				t.Fatal(err)
			}
			// Create plot
			err = core.PlotTime(plotter, config, testTime)
			if err != nil {
				t.Errorf("error plotting time: %s", err)
			}
			// Get actual output
			actual := sb.String()
			// Update golden file
			if *update {
				if err := ioutil.WriteFile(goldenFile, []byte(actual), 0644); err != nil {
					t.Fatal(err)
				}
			} else {
				// Compare actual output with expected output
				if actual != expected {
					t.Errorf("\nExpected: %s\nActual:   %s", expected, actual)
				}
			}
		})
	}
}
