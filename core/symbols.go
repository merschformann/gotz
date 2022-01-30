package core

import "fmt"

type Color string

const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed    Color = "\u001b[31m"
	ColorGreen  Color = "\u001b[32m"
	ColorYellow Color = "\u001b[33m"
	ColorBlue   Color = "\u001b[34m"
	ColorReset  Color = "\u001b[0m"
)

// colorize the given string with the given color.
func colorize(color Color, message string) string {
	return fmt.Sprint(string(color), message, string(ColorReset))
}

// Define symbol modes
const (
	// SymbolModeRectangles uses different kinds of rectangles to represent the hours.
	SymbolModeRectangles = "rectangles"
	// SymbolModeSunMoon uses the sun and moon symbols to represent the hours.
	SymbolModeSunMoon = "sun-moon"
	// SymbolModeClocks uses the 12-hour clock to represent the hours.
	SymbolModeClocks = "clocks"
	// SymbolModeDefault is the default symbol mode.
	SymbolModeDefault = SymbolModeRectangles
)

var (
	// ClockSymbols is a map of hour to clock symbol.
	ClockSymbols = []string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚"}
)

// GetHourSymbol returns a symbol representing the hour in a day.
func GetHourSymbol(mode string, color bool, hour int) string {
	// Try https://en.wikipedia.org/wiki/Geometric_Shapes_(Unicode)
	// Small sanity check
	if hour < 0 || hour > 23 {
		panic(fmt.Sprintf("invalid hour: %d", hour))
	}
	// Find matching symbol
	var s string
	switch mode {
	case SymbolModeSunMoon:
		switch {
		case hour >= 6 && hour < 12:
			s = "☀"
		case hour >= 12 && hour < 18:
			s = "☼"
		default:
			s = "☾"
		}
	case SymbolModeClocks:
		s = ClockSymbols[hour]
	default:
		switch hour {
		case 18, 19, 20, 21, 22, 23, 0, 1, 2, 3, 4, 5:
			s = "▲"
		case 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17:
			s = "▼"
		}
	}
	// Colorize, if requested
	if color {
		switch {
		case hour >= 6 && hour < 12:
			s = colorize(ColorGreen, s)
		case hour >= 12 && hour < 18:
			s = colorize(ColorYellow, s)
		default:
			s = colorize(ColorRed, s)
		}
	}
	return s
}
