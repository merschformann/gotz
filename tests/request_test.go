package core_test

import (
	"testing"
	"time"

	"github.com/merschformann/gotz/core"
)

func TestParseRequest(t *testing.T) {
	defaultConfig := core.DefaultConfig()
	londonTZ, _ := time.LoadLocation("Europe/London")
	berlinTZ, _ := time.LoadLocation("Europe/Berlin")
	// Define test cases
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "London",
			input:    "1100@Europe/London",
			expected: time.Date(2023, 10, 1, 11, 0, 0, 0, londonTZ),
		},
		{
			name:     "Berlin",
			input:    "7pm@Europe/Berlin",
			expected: time.Date(2023, 10, 1, 19, 0, 0, 0, berlinTZ),
		},
		{
			name:     "BerlinIndexed",
			input:    "7pm@2",
			expected: time.Date(2023, 10, 1, 19, 0, 0, 0, berlinTZ),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Parse the request
			parsedTime, err := core.ParseRequestTime(defaultConfig, test.input)
			if err != nil {
				t.Fatalf("Error parsing request: %v", err)
			}
			// Check if the parsed time matches the expected time (ignore the
			// date)
			if parsedTime.Hour() != test.expected.Hour() || parsedTime.Minute() != test.expected.Minute() {
				t.Errorf("Expected %v, got %v", test.expected, t)
			}
			// Check if the timezone matches the expected timezone
			if parsedTime.Location().String() != test.expected.Location().String() {
				t.Errorf("Expected timezone %v, got %v", test.expected.Location(), parsedTime.Location())
			}
		})
	}
}
