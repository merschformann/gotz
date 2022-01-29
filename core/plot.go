package core

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/term"
)

type timeslot struct {
	Time time.Time
}

func (c Config) PlotTime(request Request) error {
	// Set hours to plot
	hours := 24
	// Get terminal width
	width := getTerminalWidth()
	if !c.Stretch {
		width = width / 24 * 24
	}
	// Get current time
	if request.Time == nil {
		t := time.Now()
		request.Time = &t
	}
	// Determine time slot basics
	timeSlots := make([]timeslot, width)
	nowSlot := width / 2
	slotMinutes := hours * 60 / width
	offsetMinutes := slotMinutes * width / 2
	// Print header
	fmt.Println(strings.Repeat(" ", nowSlot-4) + "now v " + request.Time.Format("15:04:05"))
	// Prepare slots
	for i := 0; i < width; i++ {
		// Get time of slot
		slotTime := request.Time.Add(time.Duration(i*slotMinutes-offsetMinutes) * time.Minute)
		// Store timeslot info
		timeSlots[i] = timeslot{
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
	descriptionLength := maxStringLength(descriptions)

	// Plot all timezones
	for i := range timezones {
		// Print header
		desc := fmt.Sprintf("%-*s", descriptionLength, descriptions[i])
		desc = fmt.Sprintf(
			"%s: %s %s",
			desc,
			formatDay((*request.Time).In(timezones[i])),
			formatTime((*request.Time).In(timezones[i])))
		if len(desc)-1 < nowSlot {
			desc = desc + strings.Repeat(" ", nowSlot-len(desc)) + "|"
		}
		fmt.Println(desc)
		for j := 0; j < width; j++ {
			// Convert to tz time
			tzTime := timeSlots[j].Time.In(timezones[i])
			// Get symbol of slot
			symbol := GetHourSymbol(c.Symbols, c.Colorize, tzTime.Hour())
			if j == nowSlot {
				symbol = "|"
			}
			fmt.Print(symbol)
		}
		fmt.Println()
	}

	// Print markers
	printMarkers(timeSlots, width)

	return nil
}

func printMarkers(timeSlots []timeslot, width int) {
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

func maxStringLength(s []string) int {
	length := 0
	for _, str := range s {
		if len(str) > length {
			length = len(str)
		}
	}
	return length
}

func formatTime(t time.Time) string {
	return t.Format("15:04:05")
}

func formatDay(t time.Time) string {
	return t.Format("Mon 02 Jan 2006")
}

// getTerminalWidth returns the width of the terminal.
func getTerminalWidth() int {
	width, _, err := term.GetSize(0)
	if err != nil || width < 24 {
		return 72
	}
	return width
}
