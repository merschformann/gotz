package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type Config struct {
	// All timezones to display
	Timezones []Location
}

type Location struct {
	Name string
	TZ   string
}

// Default configuration generator.
func Default() Config {
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
	}
}

// DefaultConfigFile is the path of the default configuration file.
func DefaultConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Could not get user home directory: %s", err))
	}
	return path.Join(home, ".gotz.config.json")
}

// Load configuration from file.
func Load() (Config, error) {
	// If no configuration file exists, create one
	if _, err := os.Stat(DefaultConfigFile()); os.IsNotExist(err) {
		c := Default()
		c.Save()
		return c, nil
	}
	// Read configuration file
	var config Config
	data, err := ioutil.ReadFile(DefaultConfigFile())
	if err != nil {
		return config, err
	}
	// Unmarshal
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// Save configuration to file.
func (c *Config) Save() error {
	// Marshal
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	// Write file
	err = ioutil.WriteFile(DefaultConfigFile(), data, 0644)
	if err != nil {
		return err
	}
	return nil
}
