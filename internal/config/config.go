package config

import (
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	BaseDirectory string // Base directory for SOP files (default: ~/.opsy/sops/)
	LogDirectory  string // Directory for logs (default: ~/.opsy/logs/)
}

// DefaultBaseDirectory returns the default base directory for SOPs
func DefaultBaseDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".opsy", "sops")
}

// DefaultLogDirectory returns the default directory for logs
func DefaultLogDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".opsy", "logs")
}

// GetConfig returns the application configuration
func GetConfig() *Config {
	return &Config{
		BaseDirectory: DefaultBaseDirectory(),
		LogDirectory:  DefaultLogDirectory(),
	}
}