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

// plotter compiles functionality & configuration for plotting.
type plotter struct {
	plotLine      func(t ContextType, msgs ...interface{}) // func for plotting a line (with line-break)
	plotString    func(t ContextType, msg string)          // func for plotting simple strings
	terminalWidth int                                      // Terminal width
	modeStatic    bool                                     // Whether to plot static or live
	now           bool                                     // Whether to plot the current time
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
	if c.Live {
		// --> Plot time using tcell
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

		// Initialize styles
		styles := map[ContextType]tcell.Style{
			ContextBackground: tcell.StyleDefault.Background(tcell.ColorBlack),
			ContextForeground: tcell.StyleDefault.Foreground(tcell.ColorWhite),
			ContextMorning:    tcell.StyleDefault.Foreground(tcell.ColorGreen),
			ContextDay:        tcell.StyleDefault.Foreground(tcell.ColorYellow),
			ContextEvening:    tcell.StyleDefault.Foreground(tcell.ColorRed),
			ContextNight:      tcell.StyleDefault.Foreground(tcell.ColorBlue),
		}

		// Define plotting functions for tcell
		x, y := 0, 0
		plotLine := func(t ContextType, msgs ...interface{}) {
			for _, msg := range msgs {
				fmt.Println(msg)
				for _, r := range fmt.Sprint(msg) {
					s.SetContent(x, y, r, nil, styles[t])
					x++
				}
			}
			x = 0
			y++
		}
		plotString := func(t ContextType, msg string) {
			fmt.Println(msg)
			for i, r := range fmt.Sprint(msg) {
				s.SetContent(x+i, y, r, nil, styles[t])
				x++
			}
		}

		// Prepare plotter
		plt := plotter{plotLine: plotLine, plotString: plotString, now: true}

		// Track update events
		width, height := s.Size()
		now := time.Time{} // Requested time is discarded in live mode; first 'now' is set to trigger refresh

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
				plt.terminalWidth = w
				// Refresh time (pass zero time to indicate now should be used)
				s.Clear()
				err := plotTime(plt, c, now)
				if err != nil {
					return err
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
		plt := plotter{
			now:           t.IsZero(),
			terminalWidth: getTerminalWidth(),
			plotLine: func(t ContextType, line ...interface{}) {
				if c, ok := colorMap[t]; ok && c != "" {
					fmt.Println(c + fmt.Sprint(line) + ColorReset)
				} else {
					fmt.Println(line...)
				}
			},
			plotString: func(t ContextType, msg string) {
				if c, ok := colorMap[t]; ok && c != "" {
					fmt.Print(c + fmt.Sprint(msg) + ColorReset)
				} else {
					fmt.Print(msg)
				}
			},
		}
		// Get current time, if no specific time was requested
		if plt.now {
			t = time.Now()
		}
		// Plot
		err := plotTime(plt, c, t)
		if err != nil {
			return err
		}
	}

	return nil
}

// plotTime plots the time on the terminal.
func plotTime(plt plotter, cfg Config, t time.Time) error {
	// Get terminal width
	width := plt.terminalWidth
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
	// Print header
	nowDescription := "now"
	if !plt.now {
		nowDescription = "time"
	}
	plt.plotLine(
		ContextForeground,
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
		// Print header
		desc := fmt.Sprintf("%-*s", descriptionLength, descriptions[i])
		desc = fmt.Sprintf(
			"%s: %s %s",
			desc,
			formatDay(cfg.Hours12, t.In(timezones[i])),
			formatTime(cfg.Hours12, t.In(timezones[i])))
		if len(desc)-1 < nowSlot {
			desc = desc + strings.Repeat(" ", nowSlot-len(desc)) + "|"
		}
		plt.plotLine(ContextForeground, desc)
		for j := 0; j < width; j++ {
			// Convert to tz time
			tzTime := timeSlots[j].Time.In(timezones[i])
			// Get symbol of slot
			s := GetHourSymbol(cfg.Style, tzTime.Hour())
			// Get segment type of slot
			seg := getDaySegment(cfg.Style.DaySegmentation, tzTime.Hour())
			// Colorize statically if requested
			if cfg.Style.Colorize && plt.modeStatic {
				colorizeStatic(cfg.Style, tzTime.Hour(), s)
			}
			if j == nowSlot {
				s = "|"
			}
			plt.plotString(seg, s)
		}
		plt.plotLine(ContextForeground)
	}

	// Print tics
	if cfg.Tics {
		plotTics(plt, cfg.Hours12, timeSlots, width)
	}

	return nil
}

// plotTics adds tics to the plot.
func plotTics(plt plotter, hours12 bool, timeSlots []timeslot, width int) {
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
	// Print tics
	for i := 0; i < width; i++ {
		if tics[i] != "" {
			plt.plotString(ContextForeground, "^")
		} else {
			plt.plotString(ContextForeground, " ")
		}
	}
	plt.plotLine(ContextForeground)
	// Print tics
	for i := 0; i < width; i++ {
		if tics[i] != "" && i+len(tics[i]) < width {
			plt.plotString(ContextForeground, tics[i])
			i += len(tics[i]) - 1
		} else {
			plt.plotString(ContextForeground, " ")
		}
	}
	plt.plotLine(ContextForeground)
}
