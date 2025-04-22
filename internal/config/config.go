package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
}

// Load reads the config file and returns the configuration
func Load(configFile string) (*Config, error) {
	// Set default config file path if not provided
	if configFile == "" {
		configFile = "config.json"
	}

	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
