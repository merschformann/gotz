package core

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type timeslot struct {
	Time time.Time
}

// Plotter compiles functionality & configuration for plotting.
type Plotter struct {
	// func for plotting a line (with line-break)
	PlotLine func(t ContextType, msgs ...interface{})
	// func for plotting simple strings
	PlotString func(t ContextType, msg string)
	// Terminal width
	TerminalWidth int
	// Whether to plot the current time
	Now bool
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

// updateTimeNeeded indicates whether the time shown should be updated.
func updateTimeNeeded(shown, now time.Time) bool {
	// Check if time has changed
	return shown.Second() != now.Second() ||
		shown.Minute() != now.Minute() ||
		shown.Hour() != now.Hour() ||
		shown.Day() != now.Day() ||
		shown.Month() != now.Month() ||
		shown.Year() != now.Year()
}

// formatDay formats the day in the default way.
func formatDay(twelve bool, t time.Time) string {
	return t.Format("Mon 02 Jan 2006")
}

// Plot is the main plotting function. It either plots to the terminal in a
// conventional way or uses tcell for providing a continuously updating plot.
func Plot(c Config, t time.Time) error {
	if c.Live && t.IsZero() /* Only enter live mode if no time was requested */ {
		// --> Plot time using tcell
		// Initialize styles
		styles := map[ContextType]tcell.Style{}
		if c.Style.Colorize {
			styles = getDynamicColorMap(c.Style.Coloring)
		}

		// Initialize screen
		s, err := tcell.NewScreen()
		if err != nil {
			return fmt.Errorf("error creating tcell screen for live view - %+v", err)
		}
		if err := s.Init(); err != nil {
			return fmt.Errorf("error initializing tcell screen for live view - %+v", err)
		}
		exit := func() {
			s.Fini()
			os.Exit(0)
		}

		// Track update events
		width, height := s.Size()
		now := time.Time{} // Requested time is discarded in live mode; first 'now' is set to trigger refresh

		// Define plotting functions for tcell
		x, y := 0, 0
		plotLine := func(t ContextType, msgs ...interface{}) {
			// Get style
			style := tcell.StyleDefault
			if c.Style.Colorize {
				style = styles[t]
			}
			// Print message
			for _, msg := range msgs {
				for _, r := range fmt.Sprint(msg) {
					s.SetContent(x, y, r, nil, style)
					x++
				}
			}
			// Fill previous line to the end
			for i := x; i < width; i++ {
				s.SetContent(i, y, ' ', nil, style)
			}
			// Move cursor to next line
			x = 0
			y++
		}
		plotString := func(t ContextType, msg string) {
			// Get style
			style := tcell.StyleDefault
			if c.Style.Colorize {
				style = styles[t]
			}
			// Print message
			for i, r := range fmt.Sprint(msg) {
				s.SetContent(x+i, y, r, nil, style)
				x++
			}
		}

		// Prepare plotter
		plt := Plotter{PlotLine: plotLine, PlotString: plotString, Now: true}

		// Refresh time periodically
		updateTimeout := time.Duration(40) * time.Millisecond

		// Enter main loop
		for {
			// Check whether to refresh the plot (due to time or resizing)
			w, h := s.Size()
			t = time.Now()
			if w != width || h != height || updateTimeNeeded(now, t) {
				// Update dynamic plot information
				width, height = w, h
				now = t
				x, y = 0, 0
				plt.TerminalWidth = w
				// Refresh time (pass zero time to indicate now should be used)
				s.Clear()
				err := PlotTime(plt, c, now)
				if err != nil {
					return err
				}
				// Fill remaining lines
				for i := y; i < h; i++ {
					for j := 0; j < w; j++ {
						s.SetContent(j, i, ' ', nil, styles[ContextNormal])
					}
				}
				// Update screen
				s.Sync()
			}

			// Poll event or simply wait
			if s.HasPendingEvent() {
				// Grab the event
				ev := s.PollEvent()
				// Process event
				switch ev := ev.(type) {
				case *tcell.EventResize:
					s.Sync()
				case *tcell.EventKey:
					if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
						exit()
					}
				}
			} else {
				// Just sleep before redrawing
				time.Sleep(updateTimeout)
			}
		}
	} else {
		// --> Plot time using fmt
		// Prepare plotter
		colorMap := getStaticColorMap(c.Style.Coloring)
		plt := Plotter{
			Now:           t.IsZero(),
			TerminalWidth: getTerminalWidth(),
			PlotLine: func(t ContextType, line ...interface{}) {
				if ch, ok := colorMap[t]; ok && ch != "" && c.Style.Colorize {
					fmt.Println(ch + fmt.Sprint(line) + ColorReset)
				} else {
					fmt.Println(line...)
				}
			},
			PlotString: func(t ContextType, msg string) {
				if ch, ok := colorMap[t]; ok && ch != "" && c.Style.Colorize {
					fmt.Print(ch + fmt.Sprint(msg) + ColorReset)
				} else {
					fmt.Print(msg)
				}
			},
		}
		// Get current time, if no specific time was requested
		if plt.Now {
			t = time.Now()
		}
		// Plot
		err := PlotTime(plt, c, t)
		if err != nil {
			return err
		}
	}

	return nil
}

// PlotTime plots the time on the terminal.
func PlotTime(plt Plotter, cfg Config, t time.Time) error {
	// Get terminal width
	width := plt.TerminalWidth
	// Set hours to plot
	hours := 24
	// Get terminal width
	if !cfg.Stretch {
		width = width / 24 * 24
	}
	// Determine time slot basics
	timeSlots := make([]timeslot, width)
	nowSlot := width / 2
	slotMinutes := hours * 60 / width
	offsetMinutes := slotMinutes * width / 2
	// Plot header
	nowDescription := "now"
	if !plt.Now {
		nowDescription = "time"
	}
	plt.PlotLine(
		ContextNormal,
		strings.Repeat(" ",
			nowSlot-(len(nowDescription)+1))+
			nowDescription+" v "+
			formatTime(cfg.Hours12, t))
	// Prepare slots
	for i := 0; i < width; i++ {
		// Get time of slot
		slotTime := t.Add(time.Duration(i*slotMinutes-offsetMinutes) * time.Minute)
		// Store timeslot info
		timeSlots[i] = timeslot{
			Time: slotTime,
		}
	}

	// Prepare timezones to plot
	timezones := make([]*time.Location, len(cfg.Timezones)+1)
	descriptions := make([]string, len(cfg.Timezones)+1)
	timezones[0] = time.Local
	descriptions[0] = "Local"
	for i, tz := range cfg.Timezones {
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
		// --> Plot header
		desc := fmt.Sprintf("%-*s", descriptionLength, descriptions[i])
		desc = fmt.Sprintf(
			"%s: %s %s",
			desc,
			formatDay(cfg.Hours12, t.In(timezones[i])),
			formatTime(cfg.Hours12, t.In(timezones[i])))
		if len(desc)-1 < nowSlot {
			desc = desc + strings.Repeat(" ", nowSlot-len(desc)) + "|"
		}
		plt.PlotLine(ContextNormal, desc)
		// --> Plot timeslots
		for j := 0; j < width; j++ {
			// Convert to tz time
			tzTime := timeSlots[j].Time.In(timezones[i])
			// Get symbol of slot
			s := GetHourSymbol(cfg.Style, tzTime.Hour())
			// Get segment type of slot
			seg := getDaySegment(cfg.Style.DaySegmentation, tzTime.Hour())
			if j == nowSlot {
				s = "|"
			}
			plt.PlotString(seg, s)
		}
		plt.PlotLine(ContextNormal)
	}

	// Plot tics
	if cfg.Tics {
		plotTics(plt, cfg.Hours12, timeSlots, width)
	}

	return nil
}

// plotTics adds tics to the plot.
func plotTics(plt Plotter, hours12 bool, timeSlots []timeslot, width int) {
	// Prepare tics
	tics := make([]string, width)
	currentHour := -1
	for i := 0; i < width; i++ {
		// Get hour of slot
		hour := timeSlots[i].Time.Truncate(time.Hour)
		if hour.Hour()%3 == 0 && hour.Hour() != currentHour {
			if hours12 {
				tics[i] = hour.Format("3PM")
			} else {
				tics[i] = fmt.Sprintf("%d", hour.Hour())
			}
			currentHour = hour.Hour()
		}
	}
	// Plot tics
	for i := 0; i < width; i++ {
		if tics[i] != "" {
			plt.PlotString(ContextNormal, "^")
		} else {
			plt.PlotString(ContextNormal, " ")
		}
	}
	plt.PlotLine(ContextNormal)
	// Plot tics
	for i := 0; i < width; i++ {
		if tics[i] != "" && i+len(tics[i]) < width {
			plt.PlotString(ContextNormal, tics[i])
			i += len(tics[i]) - 1
		} else {
			plt.PlotString(ContextNormal, " ")
		}
	}
	plt.PlotLine(ContextNormal)
}
