package core

import "sort"

// Define sorting modes
const (
	// SortingModeNone keeps the order of the timezones as they are defined.
	SortingModeNone = "none"
	// SortingModeOffset sorts the timezones by their offset.
	SortingModeOffset = "offset"
	// SortingModeName sorts the timezones by their name.
	SortingModeName = "name"
	// SortingModeDefault is the default sorting mode.
	SortingModeDefault = SortingModeNone
)

// isValidSortingMode checks if the given sorting mode is defined and valid.
func isValidSortingMode(mode string) bool {
	switch mode {
	case SortingModeNone, SortingModeOffset, SortingModeName:
		return true
	default:
		return false
	}
}

// sortByOffset sorts the given locations by their offset.
func sortByOffset(locations []Location) {
	sort.Slice(locations, func(i, j int) bool {
		return locations[i].
	})
}
