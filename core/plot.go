package core

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/term"
)

type Timeslot struct {
	Time time.Time
}

func (c Config) PlotTime() error {
	// Set hours to plot
	hours := 24
	// Get terminal width
	width := GetTerminalWidth()
	if !c.Stretch {
		width = width / 24 * 24
	}
	// Get current time
	t := time.Now()
	// Determine time slot basics
	timeSlots := make([]Timeslot, width)
	nowSlot := width / 2
	slotMinutes := hours * 60 / width
	offsetMinutes := slotMinutes * width / 2
	// Print header
	fmt.Println(strings.Repeat(" ", nowSlot-4) + "now v " + t.Format("15:04:05"))
	// Prepare slots
	for i := 0; i < width; i++ {
		// Get time of slot
		slotTime := t.Add(time.Duration(i*slotMinutes-offsetMinutes) * time.Minute)
		// Store timeslot info
		timeSlots[i] = Timeslot{
			Time: slotTime,
		}
	}

	// Prepare timezones to plot
	timezones := make([]*time.Location, len(c.Timezones)+1)
	descriptions := make([]string, len(c.Timezones)+1)
	timezones[0] = time.Local
	descriptions[0] = "Local"
	for i, tz := range c.Timezones {
		// Get timezone
		loc, err := time.LoadLocation(tz.TZ)
		if err != nil {
			return fmt.Errorf("error loading timezone %s: %s", tz.TZ, err)
		}
		// Store timezone
		timezones[i+1] = loc
		descriptions[i+1] = tz.Name
	}
	descriptionLength := MaxStringLength(descriptions)

	// Plot all timezones
	for i := range timezones {
		// Print header
		desc := fmt.Sprintf("%-*s", descriptionLength, descriptions[i])
		desc = fmt.Sprintf("%s: %s %s", desc, FormatDay(t.In(timezones[i])), FormatTime(t.In(timezones[i])))
		if len(desc)-1 < nowSlot {
			desc = desc + strings.Repeat(" ", nowSlot-len(desc)) + "|"
		}
		fmt.Println(desc)
		for j := 0; j < width; j++ {
			// Convert to tz time
			tzTime := timeSlots[j].Time.In(timezones[i])
			// Get symbol of slot
			symbol := GetHourSymbol(tzTime.Hour())
			if j == nowSlot {
				symbol = "|"
			}
			fmt.Print(symbol)
		}
		fmt.Println()
	}

	// Print markers
	PrintMarkers(timeSlots, width)

	return nil
}

func PrintMarkers(timeSlots []Timeslot, width int) {
	// Prepare tics
	tics := make([]string, width)
	currentHour := timeSlots[0].Time.Hour()
	for i := 0; i < width; i++ {
		hour := timeSlots[i].Time.Hour()
		if hour%3 == 0 && hour != currentHour {
			tics[i] = fmt.Sprint(hour)
			currentHour = hour
		}
	}
	// Print markers
	for i := 0; i < width; i++ {
		if tics[i] != "" {
			fmt.Print("^")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	// Print tics
	for i := 0; i < width; i++ {
		if tics[i] != "" {
			fmt.Print(tics[i])
			i += len(tics[i]) - 1
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}

func MaxStringLength(s []string) int {
	length := 0
	for _, str := range s {
		if len(str) > length {
			length = len(str)
		}
	}
	return length
}

func FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}

func FormatDay(t time.Time) string {
	return t.Format("Mon 02 Jan 2006")
}

// GetTerminalWidth returns the width of the terminal.
func GetTerminalWidth() int {
	width, _, err := term.GetSize(0)
	if err != nil || width < 24 {
		return 72
	}
	return width
}

// GetHourSymbol returns a symbol representing the hour in a day.
func GetHourSymbol(hour int) string {
	switch {
	case hour >= 6 && hour < 12:
		return "☀"
	case hour >= 12 && hour < 18:
		return "☼"
	default:
		return "☾"
	}
}
