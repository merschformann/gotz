package core

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

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

// AskUser asks the user a yes/no question and returns true if the user answers
// yes.
func AskUser(question string) (bool, error) {
	// Ask the user
	fmt.Printf("%s (y/N): ", question)
	// Read user input
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	// Normalize input
	input = strings.ToLower(input)
	input = strings.TrimSpace(input)
	// Check input
	return input == "y" || input == "yes", nil
}
