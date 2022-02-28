package core

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

type timeslot struct {
	Time time.Time
}

type plotter struct {
	plotLine      func(msgs ...interface{}) // func for plotting a line (with line-break)
	plotString    func(msg string)          // func for plotting simple strings
	terminalWidth int                       // Terminal width
}

func (c Config) Plot(request Request) error {
	if c.Live {
		// --> Plot time using tcell
		// Initialize screen
		// defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
		// boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)
		s, err := tcell.NewScreen()
		if err != nil {
			return fmt.Errorf("error creating tcell screen for live view - %+v", err)
		}
		if err := s.Init(); err != nil {
			return fmt.Errorf("error initializing tcell screen for live view - %+v", err)
		}
		quit := func() {
			s.Fini()
			os.Exit(0)
		}

		// s.SetStyle(defStyle)
		// s.EnableMouse()
		// s.EnablePaste()

		// Define plotting functions for tcell
		x, y := 0, 0
		plotLine := func(msgs ...interface{}) {
			for _, msg := range msgs {
				fmt.Println(msg)
				for _, r := range fmt.Sprint(msg) {
					s.SetContent(x, y, r, nil, tcell.StyleDefault)
					x++
				}
			}
			x = 0
			y++
		}
		plotString := func(msg string) {
			fmt.Println(msg)
			for i, r := range fmt.Sprint(msg) {
				s.SetContent(x+i, y, r, nil, tcell.StyleDefault)
				x++
			}
		}

		// Track terminal size
		width, height := 0, 0

		// Enter main loop
		for {
			// Refresh if size changed (or initializing)
			w, h := s.Size()
			if w != width || h != height {
				// Clear screen
				s.Clear()
				// Update plot info
				width = w
				x, y = 0, 0
				plt := plotter{
					terminalWidth: width,
					plotLine:      plotLine,
					plotString:    plotString,
				}
				// Redraw
				err := c.plotTime(request, plt)
				if err != nil {
					return err
				}
				s.Sync()
			}

			// Update screen
			s.Show()

			// Poll event
			ev := s.PollEvent()

			// Process event
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					quit()
				} else if ev.Key() == tcell.KeyCtrlL {
					s.Sync()
				}
			}
		}
	} else {
		// --> Plot time using fmt
		err := c.plotTime(request, plotter{
			terminalWidth: getTerminalWidth(),
			plotLine:      func(line ...interface{}) { fmt.Println(line...) },
			plotString:    func(msg string) { fmt.Print(msg) },
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c Config) plotTime(request Request, plt plotter) error {
	// Get terminal width
	width := plt.terminalWidth
	// Set hours to plot
	hours := 24
	// Get terminal width
	if !c.Stretch {
		width = width / 24 * 24
	}
	// Get current time
	requestedTime := true
	if request.Time == nil {
		requestedTime = false
		t := time.Now()
		request.Time = &t
	}
	// Determine time slot basics
	timeSlots := make([]timeslot, width)
	nowSlot := width / 2
	slotMinutes := hours * 60 / width
	offsetMinutes := slotMinutes * width / 2
	// Print header
	nowDescription := "now"
	if requestedTime {
		nowDescription = "time"
	}
	plt.plotLine(strings.Repeat(" ",
		nowSlot-(len(nowDescription)+1)) +
		nowDescription + " v " +
		formatTime(c.Hours12, *request.Time))
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
			formatDay(c.Hours12, (*request.Time).In(timezones[i])),
			formatTime(c.Hours12, (*request.Time).In(timezones[i])))
		if len(desc)-1 < nowSlot {
			desc = desc + strings.Repeat(" ", nowSlot-len(desc)) + "|"
		}
		plt.plotLine(desc)
		for j := 0; j < width; j++ {
			// Convert to tz time
			tzTime := timeSlots[j].Time.In(timezones[i])
			// Get symbol of slot
			symbol := GetHourSymbol(c.Symbols, c.DaySegments, c.Colorize, tzTime.Hour())
			if j == nowSlot {
				symbol = "|"
			}
			plt.plotString(symbol)
		}
		plt.plotLine()
	}

	// Print tics
	if c.Tics {
		printTics(c.Hours12, timeSlots, width, plt)
	}

	return nil
}

// printTics prints the tics on the plot.
func printTics(twelve bool, timeSlots []timeslot, width int, plt plotter) {
	// Prepare tics
	tics := make([]string, width)
	currentHour := -1
	for i := 0; i < width; i++ {
		// Get hour of slot
		hour := timeSlots[i].Time.Truncate(time.Hour)
		if hour.Hour()%3 == 0 && hour.Hour() != currentHour {
			if twelve {
				tics[i] = hour.Format("3PM")
			} else {
				tics[i] = fmt.Sprintf("%d", hour.Hour())
			}
			currentHour = hour.Hour()
		}
	}
	// Print tics
	for i := 0; i < width; i++ {
		if tics[i] != "" {
			plt.plotString("^")
		} else {
			plt.plotString(" ")
		}
	}
	plt.plotLine()
	// Print tics
	for i := 0; i < width; i++ {
		if tics[i] != "" && i+len(tics[i]) < width {
			plt.plotString(tics[i])
			i += len(tics[i]) - 1
		} else {
			plt.plotString(" ")
		}
	}
	plt.plotLine()
}

// maxStringLength returns the length of the longest string in the given slice.
func maxStringLength(s []string) int {
	length := 0
	for _, str := range s {
		if len(str) > length {
			length = len(str)
		}
	}
	return length
}

// formatTime formats the time in the default way (distinguishing 12/24 hours
// though).
func formatTime(twelve bool, t time.Time) string {
	if twelve {
		return t.Format("3:04PM")
	} else {
		return t.Format("15:04")
	}
}

// formatDay formats the day in the default way.
func formatDay(twelve bool, t time.Time) string {
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
