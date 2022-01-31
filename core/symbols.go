package core

import "fmt"

const (
	ColorBlack   string = "\u001b[30m"
	ColorWhite   string = "\u001b[37m"
	ColorRed     string = "\u001b[31m"
	ColorYellow  string = "\u001b[33m"
	ColorMagenta string = "\u001b[35m"
	ColorGreen   string = "\u001b[32m"
	ColorCyan    string = "\u001b[36m"
	ColorBlue    string = "\u001b[34m"
	ColorReset   string = "\u001b[0m"
)

// namedColors defines all terminal colors supported by name.
var namedColors = map[string]string{
	"black":   ColorBlack,
	"white":   ColorWhite,
	"red":     ColorRed,
	"yellow":  ColorYellow,
	"magenta": ColorMagenta,
	"green":   ColorGreen,
	"blue":    ColorBlue,
	"cyan":    ColorCyan,
}

// colorize the given string with the given color.
func colorize(color string, message string) string {
	if c, ok := namedColors[color]; ok {
		return fmt.Sprint(string(c), message, ColorReset)
	}
	return fmt.Sprint(string(color), message, string(ColorReset))
}

// Define symbol modes
const (
	// SymbolModeRectangles uses different kinds of rectangles to represent the hours.
	SymbolModeRectangles = "rectangles"
	// SymbolModeSunMoon uses the sun and moon symbols to represent the hours.
	SymbolModeSunMoon = "sun-moon"
	// SymbolModeMono uses a single character to represent the hours (use
	// coloring instead).
	SymbolModeMono = "mono"
	// SymbolModeDefault is the default symbol mode.
	SymbolModeDefault = SymbolModeRectangles
)

// checkSymbolMode checks whether the given symbol mode is valid (true if valid).
func checkSymbolMode(mode string) bool {
	switch mode {
	case SymbolModeRectangles, SymbolModeSunMoon, SymbolModeMono:
		return true
	default:
		return false
	}
}

const (
	DaySegmentNight   = "night"
	DaySegmentMorning = "morning"
	DaySegmentDay     = "day"
	DaySegmentEvening = "evening"
)

var (
	// SunMoonSymbols is a map of day segment to sun/moon symbol.
	SunMoonSymbols = map[string]string{
		DaySegmentNight:   "☾",
		DaySegmentMorning: "☼",
		DaySegmentDay:     "☀",
		DaySegmentEvening: "☼",
	}
	// RectangleSymbols is a map of day segment to rectangle symbol.
	RectangleSymbols = map[string]string{
		DaySegmentNight:   " ",
		DaySegmentMorning: "▒",
		DaySegmentDay:     "█",
		DaySegmentEvening: "▒",
	}
)

// getDaySegment returns the day segment for the given hour.
func getDaySegment(seg DaySegmentation, hour int) string {
	switch {
	case hour < seg.MorningHour || hour >= seg.NightHour:
		return DaySegmentNight
	case hour < seg.DayHour:
		return DaySegmentMorning
	case hour < seg.EveningHour:
		return DaySegmentDay
	case hour < seg.NightHour:
		return DaySegmentEvening
	default:
		panic(fmt.Sprintf("invalid hour: %d", hour))
	}
}

// GetHourSymbol returns a symbol representing the hour in a day.
func GetHourSymbol(mode string, seg DaySegmentation, color bool, hour int) string {
	// Try https://en.wikipedia.org/wiki/Geometric_Shapes_(Unicode)
	// Small sanity check
	if hour < 0 || hour > 23 {
		panic(fmt.Sprintf("invalid hour: %d", hour))
	}
	// Find matching symbol
	var s string
	switch mode {
	case SymbolModeSunMoon:
		s = SunMoonSymbols[getDaySegment(seg, hour)]
	case SymbolModeMono:
		s = "#"
	default:
		s = RectangleSymbols[getDaySegment(seg, hour)]
	}
	// Colorize, if requested
	if color {
		switch getDaySegment(seg, hour) {
		case DaySegmentMorning:
			s = colorize(seg.MorningColor, s)
		case DaySegmentDay:
			s = colorize(seg.DayColor, s)
		case DaySegmentEvening:
			s = colorize(seg.EveningColor, s)
		case DaySegmentNight:
			s = colorize(seg.NightColor, s)
		default:
		}
	}
	return s
}
