package core

import "golang.org/x/term"

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
