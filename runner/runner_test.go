package runner

import (
	"autoassigner/config"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestAssign(t *testing.T) {
	// Create temporary test directory
	testDir, err := os.MkdirTemp("", "autoassigner-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Set up configuration
	config.Settings.Storage.ConfDir = testDir
	config.Settings.Storage.DataDir = filepath.Join(testDir, "data")

	// Create test group config
	groupConfig := AssigneeGroupConfig{
		Strategy:            "round_robin",
		AvailabilityChecker: "always_available",
		Users:               []string{"user1", "user2", "user3"},
	}

	configData, err := yaml.Marshal(groupConfig)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	configPath := filepath.Join(testDir, "test-group.yaml")
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test cases
	tests := []struct {
		name    string
		group   string
		dryRun  bool
		wantErr bool
	}{
		{
			name:    "valid group",
			group:   "test-group",
			dryRun:  false,
			wantErr: false,
		},
		{
			name:    "valid group dry run",
			group:   "test-group",
			dryRun:  true,
			wantErr: false,
		},
		{
			name:    "non-existent group",
			group:   "non-existent",
			dryRun:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Assign(tt.group, tt.dryRun)
			if (err != nil) != tt.wantErr {
				t.Errorf("Assign() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
