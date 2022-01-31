package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

// Config is the configuration struct.
type Config struct {
	// Day segmentation
	DaySegments DaySegmentation `json:"day_segments"`
	// All timezones to display
	Timezones []Location `json:"timezones"`

	// Symbol mode
	Symbols string `json:"symbols"`
	// Indicates whether to plot tics on the time axis
	Tics bool `json:"tics"`
	// Indicates whether to stretch across the terminal width at cost of accuracy
	Stretch bool `json:"stretch"`
	// Indicates whether to colorize the symbols
	Colorize bool `json:"colorize"`
	// Indicates whether to use the 24-hour clock
	Hours12 bool `json:"hours12"`
}

// Location describes a timezone the user wants to display.
type Location struct {
	// Descriptive name of the timezone
	Name string
	// Machine-readable timezone name
	TZ string
}

// DaySegmentation defines how to segment the day.
type DaySegmentation struct {
	// MorningHour is the hour at which the morning starts.
	MorningHour int `json:"morning"`
	// MorningColor is the color to use for the morning segment.
	MorningColor string `json:"morning_color"`
	// DayHour is the hour at which the day starts (basically business hours).
	DayHour int `json:"day"`
	// DayColor is the color to use for the day segment.
	DayColor string `json:"day_color"`
	// EveningHour is the hour at which the evening starts.
	EveningHour int `json:"evening"`
	// EveningColor is the color to use for the evening segment.
	EveningColor string `json:"evening_color"`
	// NightHour is the hour at which the night starts.
	NightHour int `json:"night"`
	// NightColor is the color to use for the night segment.
	NightColor string `json:"night_color"`
}

// DefaultConfig configuration generator.
func DefaultConfig() Config {
	tzs := []Location{}
	// Add some default locations
	ny, _ := time.LoadLocation("America/New_York")
	tzs = append(tzs, Location{"New York", ny.String()})
	london, _ := time.LoadLocation("Europe/London")
	tzs = append(tzs, Location{"London", london.String()})
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	tzs = append(tzs, Location{"Shanghai", shanghai.String()})
	sydney, _ := time.LoadLocation("Australia/Sydney")
	tzs = append(tzs, Location{"Sydney", sydney.String()})
	// Return default configuration
	return Config{
		Timezones: tzs,
		DaySegments: DaySegmentation{
			MorningHour:  6,
			MorningColor: string(ColorRed),
			DayHour:      8,
			DayColor:     string(ColorYellow),
			EveningHour:  18,
			EveningColor: string(ColorRed),
			NightHour:    22,
			NightColor:   string(ColorBlue),
		},
		Tics:    false,
		Stretch: true,
	}
}

// defaultConfigFile is the path of the default configuration file.
func defaultConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Could not get user home directory: %s", err))
	}
	return path.Join(home, ".gotz.config.json")
}

// Load configuration from file.
func Load() (Config, error) {
	// If no configuration file exists, create one
	if _, err := os.Stat(defaultConfigFile()); os.IsNotExist(err) {
		return saveDefault()
	}
	// Read configuration file
	var config Config
	data, err := ioutil.ReadFile(defaultConfigFile())
	if err != nil {
		fmt.Println("Error reading config file (replacing with default config):", err)
		return saveDefault()
	}
	// Unmarshal
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshaling config file (replacing with default config):", err)
		return saveDefault()
	}
	// Validate (replace invalid values with defaults)
	config = config.validate()
	return config, nil
}

// saveDefault creates a default config and immediately saves it.
func saveDefault() (Config, error) {
	c := DefaultConfig()
	return c, c.Save()
}

// Save configuration to file.
func (c *Config) Save() error {
	// Marshal
	data, err := json.Marshal(c)
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
	if !checkSymbolMode(c.Symbols) {
		fmt.Printf("Warning - invalid symbols (using default): %s\n", c.Symbols)
		c.Symbols = SymbolModeDefault
	}
	return c
}
