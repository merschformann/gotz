package core

import "fmt"

type ContextType string

const (
	ContextForeground ContextType = "foreground"
	ContextBackground ContextType = "background"
	ContextMorning    ContextType = "morning"
	ContextDay        ContextType = "day"
	ContextEvening    ContextType = "evening"
	ContextNight      ContextType = "night"
)

// getDaySegment returns the day segment for the given hour.
func getDaySegment(seg DaySegmentation, hour int) ContextType {
	switch {
	case hour < seg.MorningHour || hour >= seg.NightHour:
		return ContextNight
	case hour < seg.DayHour:
		return ContextMorning
	case hour < seg.EveningHour:
		return ContextDay
	case hour < seg.NightHour:
		return ContextEvening
	default:
		panic(fmt.Sprintf("invalid hour: %d", hour))
	}
}

// Terminal color codes.
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

// NamedColors defines all terminal colors supported by name.
var NamedColors = map[string]string{
	"black":   ColorBlack,
	"white":   ColorWhite,
	"red":     ColorRed,
	"yellow":  ColorYellow,
	"magenta": ColorMagenta,
	"green":   ColorGreen,
	"blue":    ColorBlue,
	"cyan":    ColorCyan,
}

func getStaticColorMap(sty PlotColors) map[ContextType]string {
	// Define lookup function
	getColor := func(colorValue string) string {
		// Check if color is a named color
		if color, ok := NamedColors[colorValue]; ok {
			return color
		}
		// At this point color must be a valid terminal color code
		return colorValue
	}
	// Create static color map
	staticColorMap := make(map[ContextType]string)
	staticColorMap[ContextForeground] = getColor(sty.StaticColorForeground)
	staticColorMap[ContextBackground] = "" // No coloring of background
	staticColorMap[ContextMorning] = getColor(sty.StaticColorMorning)
	staticColorMap[ContextDay] = getColor(sty.StaticColorDay)
	staticColorMap[ContextEvening] = getColor(sty.StaticColorEvening)
	staticColorMap[ContextNight] = getColor(sty.StaticColorNight)
	return staticColorMap
}

// colorizeStatic colorizes the given string with the given color. Uses terminal
// color codes or named colors reflecting the same.
func colorizeStatic(style Style, hour int, message string) string {
	// Define coloring function using terminal color codes
	colorize := func(color string) string {
		if c, ok := NamedColors[color]; ok {
			return fmt.Sprint(string(c), message, ColorReset)
		}
		return fmt.Sprint(string(color), message, string(ColorReset))
	}
	// Colorize depending on segment in day
	segment := getDaySegment(style.DaySegmentation, hour)
	switch segment {
	case ContextMorning:
		return colorize(style.Coloring.StaticColorMorning)
	case ContextDay:
		return colorize(style.Coloring.StaticColorDay)
	case ContextEvening:
		return colorize(style.Coloring.StaticColorEvening)
	case ContextNight:
		return colorize(style.Coloring.StaticColorNight)
	default:
		panic(fmt.Sprintf("invalid segment: %s", segment))
	}
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

var (
	// SunMoonSymbols is a map of day segment to sun/moon symbol.
	SunMoonSymbols = map[ContextType]string{
		ContextNight:   "☾",
		ContextMorning: "☼",
		ContextDay:     "☀",
		ContextEvening: "☼",
	}
	// RectangleSymbols is a map of day segment to rectangle symbol.
	RectangleSymbols = map[ContextType]string{
		ContextNight:   " ",
		ContextMorning: "▒",
		ContextDay:     "█",
		ContextEvening: "▒",
	}
)

// GetHourSymbol returns a symbol representing the hour in a day.
func GetHourSymbol(sty Style, hour int) string {
	// Small sanity check
	if hour < 0 || hour > 23 {
		panic(fmt.Sprintf("invalid hour: %d", hour))
	}
	// Get symbol depending on symbol mode
	switch sty.Symbols {
	case SymbolModeRectangles:
		return RectangleSymbols[getDaySegment(sty.DaySegmentation, hour)]
	case SymbolModeSunMoon:
		return SunMoonSymbols[getDaySegment(sty.DaySegmentation, hour)]
	case SymbolModeMono:
		return "#"
	default:
		panic(fmt.Sprintf("invalid symbol mode: %s", sty.Symbols))
	}
}
