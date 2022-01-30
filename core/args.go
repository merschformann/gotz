package core

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

// request defines per invokation options.
type Request struct {
	// Time to display
	Time *time.Time
}

// parseArgs parses the command line arguments and applies them to the given configuration.
func ParseFlags(startConfig Config) (Config, Request, bool, error) {
	// Check for any changes
	var changed bool
	// Define configuration flags
	var symbols, timezones, markers, stretch, colorize string
	flag.StringVar(
		&symbols,
		"symbols",
		"",
		"symbols to use for timezone markers (one of: "+
			SymbolModeRectangles+", "+
			SymbolModeSunMoon+", "+
			SymbolModeClocks+")",
	)
	flag.StringVar(
		&timezones,
		"timezones",
		"",
		"timezones to display, comma-separated (for example: 'America/New_York,Europe/London,Asia/Shanghai' or named 'Office:America/New_York,Home:Europe/London' "+
			" - for TZ names see TZ database name in https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)",
	)
	flag.StringVar(
		&markers,
		"markers",
		"",
		"indicates whether to use markers on the time axis (one of: true, false)",
	)
	flag.StringVar(
		&stretch,
		"stretch",
		"",
		"indicates whether to stretch across the terminal width at cost of accuracy (one of: true, false)",
	)
	flag.StringVar(
		&colorize,
		"colorize",
		"",
		"indicates whether to colorize the symbols (one of: true, false)",
	)

	// Define direct flags
	var requestTime string
	flag.StringVar(
		&requestTime,
		"time",
		"",
		"time to display (e.g. 20:00 or 2000 or 20 or 8pm)",
	)

	// Parse flags
	flag.Parse()

	// Handle configuration
	if symbols != "" {
		changed = true
		startConfig.Symbols = symbols
		if !checkSymbolMode(startConfig.Symbols) {
			return startConfig, Request{}, false, fmt.Errorf("invalid symbol mode: %s", symbols)
		}
	}
	if timezones != "" {
		changed = true
		tzs, err := parseTimezones(timezones)
		if err != nil {
			return startConfig, Request{}, changed, err
		}
		startConfig.Timezones = tzs
	}
	if markers != "" {
		changed = true
		if strings.ToLower(markers) == "true" {
			startConfig.Markers = true
		} else if strings.ToLower(markers) == "false" {
			startConfig.Markers = false
		} else {
			return startConfig, Request{}, changed, fmt.Errorf("invalid value for markers: %s", markers)
		}
	}
	if stretch != "" {
		changed = true
		if strings.ToLower(stretch) == "true" {
			startConfig.Stretch = true
		} else if strings.ToLower(stretch) == "false" {
			startConfig.Stretch = false
		} else {
			return startConfig, Request{}, changed, fmt.Errorf("invalid value for stretch: %s", stretch)
		}
	}
	if colorize != "" {
		changed = true
		if strings.ToLower(colorize) == "true" {
			startConfig.Colorize = true
		} else if strings.ToLower(colorize) == "false" {
			startConfig.Colorize = false
		} else {
			return startConfig, Request{}, changed, fmt.Errorf("invalid value for colorize: %s", colorize)
		}
	}

	// Handle direct flags
	var request Request
	if requestTime != "" {
		changed = true
		// Parse time
		t, err := parseTime(requestTime)
		if err != nil {
			return startConfig, Request{}, changed, err
		}
		request.Time = &t
	}

	// Handle last argument as time, if it starts with a digit
	if flag.NArg() > 0 {
		// Get last argument
		lastArg := flag.Arg(flag.NArg() - 1)
		// If last argument is a time, parse it
		if len(lastArg) > 0 && lastArg[0] >= '0' && lastArg[0] <= '9' {
			// Parse time
			t, err := parseTime(lastArg)
			if err != nil {
				return startConfig, Request{}, changed, err
			}
			request.Time = &t
		}
	}

	return startConfig, request, changed, nil
}

// parseTimezones parses a comma-separated list of timezones.
func parseTimezones(timezones string) ([]Location, error) {
	var timezoneList []Location
	for _, timezone := range strings.Split(timezones, ",") {
		// Skip empty timezones
		if timezone == "" {
			continue
		}

		if strings.Contains(timezone, ":") {
			// Handle named timezones
			parts := strings.Split(timezone, ":")
			if len(parts) != 2 {
				return timezoneList, fmt.Errorf("invalid timezone: %s", timezone)
			}
			if !checkTimezoneLocation(parts[1]) {
				return timezoneList, fmt.Errorf("invalid timezone: %s", timezone)
			}
			timezoneList = append(timezoneList, Location{
				Name: parts[0],
				TZ:   parts[1],
			})
		} else {
			// Handle simple timezones
			if !checkTimezoneLocation(timezone) {
				return timezoneList, fmt.Errorf("invalid timezone: %s", timezone)
			}
			timezoneList = append(timezoneList, Location{
				Name: timezone,
				TZ:   timezone,
			})
		}
	}
	return timezoneList, nil
}

// checkTimezoneLocation checks if a timezone name is valid.
func checkTimezoneLocation(timezone string) bool {
	_, err := time.LoadLocation(timezone)
	return err == nil
}

// parseTime parses a time string.
func parseTime(t string) (time.Time, error) {
	// Try all supported formats
	for _, format := range []string{
		"15",
		"15:04",
		"15:04:05",
		"3:04pm",
		"3:04:05pm",
		"3pm",
		"1504",
		"150405",
	} {
		if t, err := time.Parse(format, t); err == nil {
			n := time.Now()
			t = time.Date(n.Year(), n.Month(), n.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time: %s", t)
}
