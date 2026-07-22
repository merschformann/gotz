package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseTimezones parses a comma-separated list of timezones.
func ParseTimezones(timezones string) ([]Location, error) {
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

// inputTimeFormat defines accepted time formats.
type inputTimeFormat struct {
	// The format string.
	Format string
	// Indicates whether the input declared a date too.
	Date bool
	// Indicates whether the input declared a timezone too.
	TZInfo bool
}

// ParseRequestTime parses a requested time in various formats. Furthermore, it
// reads an optional timezone index and uses its timezone instead of local.
func ParseRequestTime(config Config, t string) (time.Time, error) {
	tzSeparator := "@"
	tz := time.Local
	// Check whether a different time zone than the local one was specified.
	if strings.Contains(t, tzSeparator) {
		// Split time and timezone
		parts := strings.Split(t, tzSeparator)
		if len(parts) != 2 {
			return time.Time{}, fmt.Errorf("invalid time format: %s (should be <time>@<timezone-index> or <time>@<timezone-identifier>)", t)
		}
		if checkTimezoneLocation(parts[1]) {
			// >>> Handle the explicitly given timezone
			explicitTZ, err := time.LoadLocation(parts[1])
			if err == nil {
				// Use explicit timezone
				tz = explicitTZ
			}
		} else {
			// >>> Handle the timezone index referring to the configured TZs
			tzIndex, err := strconv.Atoi(parts[1])
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid time format: %s (should be <time>@<timezone-index> or <time>@<timezone-identifier>)", t)
			}
			if tzIndex < 0 || tzIndex > len(config.Timezones) {
				return time.Time{}, fmt.Errorf("invalid time format: %s (timezone-index out of range)", t)
			}
			// Get timezone at index (offset by one to account for 0 as local timezone)
			if tzIndex > 0 {
				tz, err = time.LoadLocation(config.Timezones[tzIndex-1].TZ)
				if err != nil {
					return time.Time{}, fmt.Errorf("invalid timezone: %s (given index %d)",
						config.Timezones[tzIndex-1].TZ,
						tzIndex)
				}
			}
		}
		// Keep time element
		t = parts[0]
	}
	// Parse time
	rt, err := parseTime(t, tz)
	if err != nil {
		return time.Time{}, err
	}
	return rt, nil
}

// parseTime parses a time string in various formats.
func parseTime(t string, tz *time.Location) (time.Time, error) {
	// Try all supported formats
	for _, format := range []inputTimeFormat{
		{"15", false, false},
		{"15:04", false, false},
		{"15:04:05", false, false},
		{"3:04pm", false, false},
		{"3:04:05pm", false, false},
		{"3pm", false, false},
		{"1504", false, false},
		{"150405", false, false},
		{"2006-01-02T15:04:05", true, false},
		{"2006-01-02T15:04:05Z07:00", true, true},
	} {
		if t, err := time.Parse(format.Format, t); err == nil {
			n := time.Now()
			if !format.TZInfo {
				if format.Date {
					t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, tz)
				} else {
					t = time.Date(n.Year(), n.Month(), n.Day(), t.Hour(), t.Minute(), t.Second(), 0, tz)
				}
			}
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time: %s", t)
}
