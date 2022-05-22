package core

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
)

type ContextType string

const (
	ContextNormal  ContextType = "normal"
	ContextMorning ContextType = "morning"
	ContextDay     ContextType = "day"
	ContextEvening ContextType = "evening"
	ContextNight   ContextType = "night"
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

// NamedStaticColors defines all terminal colors supported by name.
var NamedStaticColors = map[string]string{
	"black":   ColorBlack,
	"white":   ColorWhite,
	"red":     ColorRed,
	"yellow":  ColorYellow,
	"magenta": ColorMagenta,
	"green":   ColorGreen,
	"blue":    ColorBlue,
	"cyan":    ColorCyan,
}

// getDynamicColorMap returns a map of dynamic colors for the given style
// configuration.
func getDynamicColorMap(sty PlotColors) map[ContextType]tcell.Style {
	// Define lookup function
	getColor := func(colorValue string) tcell.Color {
		// Check if color is hex color
		if strings.HasPrefix(colorValue, "#") {
			return tcell.GetColor(strings.ToLower(colorValue))
		}
		// Check if color is a named color
		if c, ok := tcell.ColorNames[strings.ToLower(colorValue)]; ok {
			return c
		}
		// Use default color
		return tcell.ColorDefault
	}
	// Get default foreground / background
	fg, bg, _ := tcell.StyleDefault.Decompose()
	if sty.DynamicColorForeground != "" {
		fg = getColor(sty.DynamicColorForeground)
	}
	if sty.DynamicColorBackground != "" {
		bg = getColor(sty.DynamicColorBackground)
	}
	baseStyle := tcell.StyleDefault.Background(bg).Foreground(fg)
	// Create dynamic color map
	dynamicColorMap := make(map[ContextType]tcell.Style)
	dynamicColorMap[ContextNormal] = baseStyle
	dynamicColorMap[ContextMorning] = baseStyle.Foreground(getColor(sty.DynamicColorMorning))
	dynamicColorMap[ContextDay] = baseStyle.Foreground(getColor(sty.DynamicColorDay))
	dynamicColorMap[ContextEvening] = baseStyle.Foreground(getColor(sty.DynamicColorEvening))
	dynamicColorMap[ContextNight] = baseStyle.Foreground(getColor(sty.DynamicColorNight))
	return dynamicColorMap
}

// getStaticColorMap returns a map of static colors for the given style
// configuration.
func getStaticColorMap(sty PlotColors) map[ContextType]string {
	// Define lookup function
	getColor := func(colorValue string) string {
		// Check if color is a named color
		if color, ok := NamedStaticColors[colorValue]; ok {
			return color
		}
		// Check if color is hex color
		if strings.HasPrefix(colorValue, "#") {
			r, g, b, err := convertHexToRgb(strings.ToLower(colorValue))
			if err != nil {
				panic(err)
			}
			return fmt.Sprintf("\u001b[38;2;%d;%d;%dm", r, g, b)
		}
		// At this point color must be a valid terminal color code
		return colorValue
	}
	// Create static color map
	staticColorMap := make(map[ContextType]string)
	staticColorMap[ContextNormal] = getColor(sty.StaticColorForeground) // Override of background not supported in static mode
	staticColorMap[ContextMorning] = getColor(sty.StaticColorMorning)
	staticColorMap[ContextDay] = getColor(sty.StaticColorDay)
	staticColorMap[ContextEvening] = getColor(sty.StaticColorEvening)
	staticColorMap[ContextNight] = getColor(sty.StaticColorNight)
	return staticColorMap
}

// Define symbol modes
const (
	// SymbolModeRectangles uses different kinds of rectangles to represent the
	// hours
	SymbolModeRectangles = "rectangles"
	// SymbolModeSunMoon uses the sun and moon symbols to represent the hours.
	SymbolModeSunMoon = "sun-moon"
	// SymbolModeMono uses a single character to represent the hours (use
	// coloring instead).
	SymbolModeMono = "mono"
	// SymbolModeBlocks uses all blocks to represent the hours.
	SymbolModeBlocks = "blocks"
	// SymbolModeCustom uses a custom user-defined symbols to represent the
	// hours.
	SymbolModeCustom = "custom"
	// SymbolModeDefault is the default symbol mode.
	SymbolModeDefault = SymbolModeRectangles
)

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

// checkSymbolMode checks if the given symbol mode is valid.
func checkSymbolMode(mode string) error {
	switch mode {
	case SymbolModeRectangles, SymbolModeSunMoon, SymbolModeMono, SymbolModeBlocks, SymbolModeCustom:
		return nil
	default:
		return fmt.Errorf("invalid symbols: %s", mode)
	}
}

// checkSymbolMode does a small sanity check on the symbol definition.
func checkSymbolConfig(sty Style) error {
	if sty.Symbols == SymbolModeCustom {
		if len(sty.CustomSymbols) <= 0 {
			return fmt.Errorf("custom symbols not defined")
		}
		seenHours := map[int]bool{}
		for _, s := range sty.CustomSymbols {
			if utf8.RuneCountInString(s.Symbol) != 1 {
				return fmt.Errorf("custom symbol %s is not a single character", s.Symbol)
			}
			if _, ok := seenHours[s.Start]; ok {
				return fmt.Errorf("duplicate custom symbol for hour %d", s.Start)
			}
			seenHours[s.Start] = true
		}
	}
	return nil
}

func GetSymbols(sty Style) []string {
	symbols := make([]string, 24)
	switch sty.Symbols {
	default:
		fallthrough
	case SymbolModeRectangles:
		for h := range symbols {
			symbols[h] = RectangleSymbols[getDaySegment(sty.DaySegmentation, h)]
		}
	case SymbolModeSunMoon:
		for h := range symbols {
			symbols[h] = SunMoonSymbols[getDaySegment(sty.DaySegmentation, h)]
		}
	case SymbolModeMono:
		for h := range symbols {
			symbols[h] = "#"
		}
	case SymbolModeBlocks:
		for h := range symbols {
			symbols[h] = "█"
		}
	case SymbolModeCustom:
		// Sort custom symbols by hour
		customSymbols := make([]TimeSymbol, len(sty.CustomSymbols))
		copy(customSymbols, sty.CustomSymbols)
		sort.Slice(customSymbols, func(i, j int) bool {
			return customSymbols[i].Start < customSymbols[j].Start
		})
		// Start with the symbol the previous day ends with
		currentSym, currentIdx := customSymbols[len(customSymbols)-1].Symbol, -1
		for h := range symbols {
			// Find the next custom symbol
			if currentIdx < len(customSymbols)-1 && h == customSymbols[currentIdx+1].Start {
				currentSym, currentIdx = customSymbols[currentIdx+1].Symbol, currentIdx+1
			}
			symbols[h] = currentSym
		}
	}
	return symbols
}
