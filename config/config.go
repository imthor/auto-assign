// Package config provides configuration management for the autoassigner.
// It handles loading and parsing of configuration files, including:
// - Storage configuration (data directory and config directory)
// - Availability configuration (API endpoints and status settings)
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// StorageConfig defines the storage-related configuration settings.
type StorageConfig struct {
	DataDir string `json:"data_dir"` // Base directory for all data files
	ConfDir string `json:"conf_dir"` // Directory for group configuration files
}

// AvailabilityConfig defines the availability-related configuration settings.
type AvailabilityConfig struct {
	InOutApiUrlPrefix        string   `json:"inout_api_url_prefix"`       // Base URL for the In/Out API
	InOutUnavailableStatuses []string `json:"inout_unavailable_statuses"` // List of statuses indicating unavailability
}

// Config represents the complete configuration for the autoassigner.
type Config struct {
	Storage      StorageConfig      `json:"storage"`      // Storage-related settings
	Availability AvailabilityConfig `json:"availability"` // Availability-related settings
}

// Settings holds the global configuration settings.
var Settings Config

// GetGroupDataDir returns the data directory for a specific group.
// It creates the directory if it doesn't exist.
func GetGroupDataDir(group string) (string, error) {
	dir := filepath.Join(Settings.Storage.DataDir, group)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// ListGroups returns a list of all valid group names from the config directory.
// A valid group is one that has a .yaml configuration file.
func ListGroups() ([]string, error) {
	entries, err := os.ReadDir(Settings.Storage.ConfDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	var groups []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			// Remove .yaml extension to get group name
			groupName := strings.TrimSuffix(entry.Name(), ".yaml")
			groups = append(groups, groupName)
		}
	}

	return groups, nil
}

// LoadConfig loads the configuration from the specified config file.
// It reads the file, parses the JSON content, and populates the Settings variable.
// Returns an error if the file cannot be read or parsed.
func LoadConfig(configPath string) error {
	// Validate config file path
	if configPath == "" {
		return fmt.Errorf("config file path cannot be empty")
	}

	// Check if file exists
	info, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file does not exist: %s", configPath)
		}
		return fmt.Errorf("failed to access config file: %w", err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return fmt.Errorf("config path is not a regular file: %s", configPath)
	}

	// Check if file is readable
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Parse JSON content
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&Settings); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if err := validateConfig(&Settings); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	return nil
}

// validateConfig checks if the configuration has all required fields.
func validateConfig(cfg *Config) error {
	if cfg.Storage.DataDir == "" {
		return fmt.Errorf("data_dir is required in storage configuration")
	}
	if cfg.Storage.ConfDir == "" {
		return fmt.Errorf("conf_dir is required in storage configuration")
	}
	if cfg.Availability.InOutApiUrlPrefix == "" {
		return fmt.Errorf("inout_api_url_prefix is required in availability configuration")
	}
	if len(cfg.Availability.InOutUnavailableStatuses) == 0 {
		return fmt.Errorf("inout_unavailable_statuses is required in availability configuration")
	}
	return nil
}
