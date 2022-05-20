package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/adrg/xdg"
)

// ConfigVersion is the current version of the configuration file.
const ConfigVersion = "1.0"

// Config is the configuration struct.
type Config struct {
	// Version is the version of the configuration file.
	ConfigVersion string `json:"config_version"`

	// All timezones to display.
	Timezones []Location `json:"timezones"`

	// Style defines the style of the timezone plot.
	Style Style `json:"style"`

	// Indicates whether to plot tics on the time axis.
	Tics bool `json:"tics"`
	// Indicates whether to stretch across the terminal width at cost of
	// accuracy.
	Stretch bool `json:"stretch"`
	// Indicates whether to use the 24-hour clock.
	Hours12 bool `json:"hours12"`

	// Indicates whether to continuously update.
	Live bool `json:"live"`
}

// Location describes a timezone the user wants to display.
type Location struct {
	// Descriptive name of the timezone.
	Name string
	// Machine-readable timezone name.
	TZ string
}

type Style struct {
	// Defines the symbols to be used.
	Symbols string `json:"symbols"`
	// Indicates whether to use colors.
	Colorize bool `json:"colorize"`
	// Defines how the day is split up into different ranges.
	DaySegmentation DaySegmentation `json:"day_segments"`
	// Defines the colors to be used in the plot.
	Coloring PlotColors `json:"coloring"`
}

// DaySegmentation defines how to segment the day.
type DaySegmentation struct {
	// MorningHour is the hour at which the morning starts.
	MorningHour int `json:"morning"`
	// DayHour is the hour at which the day starts (basically business hours).
	DayHour int `json:"day"`
	// EveningHour is the hour at which the evening starts.
	EveningHour int `json:"evening"`
	// NightHour is the hour at which the night starts.
	NightHour int `json:"night"`
}

// PlotColors defines the colors to be used in the plot.
type PlotColors struct {
	// StaticColorMorning is the color to use for the morning segment.
	StaticColorMorning string
	// StaticColorDay is the color to use for the day segment.
	StaticColorDay string
	// StaticColorEvening is the color to use for the evening segment.
	StaticColorEvening string
	// StaticColorNight is the color to use for the night segment.
	StaticColorNight string
	// StaticColorForeground is the color to use for the foreground.
	StaticColorForeground string

	// DynamicColorMorning is the color to use for the morning segment (in live mode).
	DynamicColorMorning string
	// DynamicColorDay is the color to use for the morning segment (in live mode).
	DynamicColorDay string
	// DynamicColorEvening is the color to use for the morning segment (in live mode).
	DynamicColorEvening string
	// DynamicColorNight is the color to use for the morning segment (in live mode).
	DynamicColorNight string
	// DynamicColorForeground is the color to use for the foreground (in live mode).
	DynamicColorForeground string
	// DynamicColorBackground is the color to use for the background (in live mode).
	DynamicColorBackground string
}

// DefaultConfig configuration generator.
func DefaultConfig() Config {
	tzs := []Location{}
	// Add some default locations
	ny, _ := time.LoadLocation("America/New_York")
	tzs = append(tzs, Location{"New York", ny.String()})
	london, _ := time.LoadLocation("Europe/Berlin")
	tzs = append(tzs, Location{"Berlin", london.String()})
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	tzs = append(tzs, Location{"Shanghai", shanghai.String()})
	sydney, _ := time.LoadLocation("Australia/Sydney")
	tzs = append(tzs, Location{"Sydney", sydney.String()})
	// Return default configuration
	return Config{
		ConfigVersion: ConfigVersion,
		Timezones:     tzs,
		Style: Style{
			Symbols:  SymbolModeDefault,
			Colorize: false,
			DaySegmentation: DaySegmentation{
				MorningHour: 6,
				DayHour:     8,
				EveningHour: 18,
				NightHour:   22,
			},
			Coloring: PlotColors{
				StaticColorMorning:     "red",
				StaticColorDay:         "yellow",
				StaticColorEvening:     "red",
				StaticColorNight:       "blue",
				StaticColorForeground:  "", // don't override terminal foreground color
				DynamicColorMorning:    "red",
				DynamicColorDay:        "yellow",
				DynamicColorEvening:    "red",
				DynamicColorNight:      "blue",
				DynamicColorForeground: "", // don't override foreground color
				DynamicColorBackground: "", // don't override background color
			},
		},
		Tics:    false,
		Stretch: true,
	}
}

// defaultConfigFile is the path of the default configuration file.
func defaultConfigFile() string {
	configFilePath, err := xdg.ConfigFile("gotz/config.json")
	if err != nil {
		panic(fmt.Sprintf("Could not get user config directory: %s", err))
	}
	return configFilePath
}

// Load configuration from file.
func Load() (Config, error) {
	// If no configuration file exists, create one
	if _, err := os.Stat(defaultConfigFile()); os.IsNotExist(err) {
		return SaveDefault()
	}
	// Read configuration file
	var config Config
	data, err := ioutil.ReadFile(defaultConfigFile())
	if err != nil {
		return config, errors.New("Error reading config file: " + err.Error())
	}
	// Unmarshal
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, errors.New("Error unmarshaling config file: " + err.Error())
	}
	// Check version
	if config.ConfigVersion != ConfigVersion {
		version := config.ConfigVersion
		if version == "" {
			version = "unknown"
		}
		return config, errors.New("Config file version " + version + " is not supported")
	}
	// Validate (replace invalid values with defaults)
	config = config.validate()
	return config, nil
}

// SaveDefault creates a default config and immediately saves it.
func SaveDefault() (Config, error) {
	c := DefaultConfig()
	return c, c.Save()
}

// Save configuration to file.
func (c *Config) Save() error {
	// Marshal and pretty-print
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	// Write file
	err = ioutil.WriteFile(defaultConfigFile(), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// validate validates the configuration.
func (c Config) validate() Config {
	// Check whether symbol mode is known
	if !checkSymbolMode(c.Style.Symbols) {
		fmt.Printf("Warning - invalid symbols (using default): %s\n", c.Style.Symbols)
		c.Style.Symbols = SymbolModeDefault
	}
	return c
}
