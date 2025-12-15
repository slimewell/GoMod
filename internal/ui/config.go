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
	return filepath.Join(home, ".gomod.json"), nil
}

// LoadConfig loads configuration from ~/.gomod.json
func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			defaultCfg := DefaultConfig()
			return &defaultCfg, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves configuration to ~/.gomod.json
func SaveConfig(cfg *Config) error {
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
