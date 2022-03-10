package core

import (
	"strconv"
	"strings"

	"golang.org/x/term"
)

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

// getTerminalWidth returns the width of the terminal.
func getTerminalWidth() int {
	width, _, err := term.GetSize(0)
	if err != nil || width < 24 {
		return 72
	}
	return width
}

// convertHexToRgb converts hex color code to rgb.
func convertHexToRgb(hex string) (r, g, b uint8, err error) {
	hex = strings.TrimPrefix(hex, "#")
	rgb, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, 0, 0, err
	}
	r = uint8(rgb >> 16)
	g = uint8((rgb & 0x00ff00) >> 8)
	b = uint8(rgb & 0x0000ff)
	return r, g, b, nil
}
