package core

import (
	"sort"
	"time"
)

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

// locationContainer is a container for a location with additional information.
type locationContainer struct {
	location    *time.Location
	description string
	offset      int
}

// sortLocations sorts the given locations based on the given sorting mode.
func sortLocations(locations []locationContainer, sortingMode string, localTop bool) {
	sort.Slice(locations, func(i, j int) bool {
		// If the local timezone should be kept at the top, check if one of the
		// locations is the local timezone.
		if localTop {
			if locations[i].location == time.Local {
				return true
			} else if locations[j].location == time.Local {
				return false
			}
		}
		// Sort based on the sorting mode
		switch sortingMode {
		case SortingModeOffset:
			return locations[i].offset < locations[j].offset
		case SortingModeName:
			return locations[i].description < locations[j].description
		default:
			return i < j
		}
	})
}
