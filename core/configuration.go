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
	// Symbol mode
	Symbols string `json:"symbols"`
	// All timezones to display
	Timezones []Location `json:"timezones"`
	// Indicates whether to plot markers on the time axis
	Markers bool `json:"tics"`
	// Indicates whether to stretch across the terminal width at cost of accuracy
	Stretch bool `json:"stretch"`
	// Indicates whether to colorize the symbols
	Colorize bool `json:"colorize"`
}

// Location describes a timezone the user wants to display.
type Location struct {
	// Descriptive name of the timezone
	Name string
	// Machine-readable timezone name
	TZ string
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
		Markers:   true,
		Stretch:   true,
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
	if c.Symbols != SymbolModeRectangles && c.Symbols != SymbolModeClocks && c.Symbols != SymbolModeSunMoon {
		fmt.Printf("Warning - invalid symbols (using default): %s\n", c.Symbols)
		c.Symbols = SymbolModeDefault
	}
	return c
}
