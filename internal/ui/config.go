package ui

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds persistent user preferences
type Config struct {
	Theme     string `json:"theme"`
	StereoSep int    `json:"stereo_separation"`
	LastUsed  string `json:"last_file,omitempty"`
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		Theme:     "default",
		StereoSep: 50,
	}
}

// configPath returns the path to the config file
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".modtui.json"), nil
}

// LoadConfig loads configuration from ~/.modtui.json
func LoadConfig() (Config, error) {
	path, err := configPath()
	if err != nil {
		return DefaultConfig(), err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// File doesn't exist, return defaults
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	return cfg, nil
}

// SaveConfig saves configuration to ~/.modtui.json
func SaveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
