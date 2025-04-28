// Package runner provides the core functionality for automatically assigning tasks to team members.
// It implements the main assignment logic, including:
// - Loading and parsing group configurations
// - Selecting assignees based on various strategies
// - Checking member availability
// - Tracking assignment history
package runner

import (
	"autoassigner/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// AssigneeGroupConfig represents the configuration for a group of assignees.
// It specifies the selection strategy, availability checker, and list of users.
type AssigneeGroupConfig struct {
	Strategy            string   `yaml:"strategy"`             // The strategy to use for selecting assignees
	AvailabilityChecker string   `yaml:"availability_checker"` // The type of availability checker to use
	Users               []string `yaml:"users"`                // List of users in the group
}

// AssignmentLog represents a single assignment entry in the log file.
type AssignmentLog struct {
	Timestamp  string `json:"timestamp"`
	Group      string `json:"group"`
	User       string `json:"user"`
	Strategy   string `json:"strategy"`
	LastIndex  int    `json:"last_index"`
	NextIndex  int    `json:"next_index"`
	TotalCount int    `json:"total_count"`
	UserCount  int    `json:"user_count"`
}

// Assign selects an available assignee from the specified group.
// It uses the configured strategy to select a user and checks their availability.
// If dryRun is true, it will simulate the assignment without updating any logs or counts.
// Returns an error if no available assignee is found or if there are configuration issues.
func Assign(group string, dryRun bool) error {
	factory := NewComponentFactory(
		&DefaultConfigLoader{},
		&DefaultStorageManager{},
		&DefaultCountManager{},
		&DefaultAssignmentLogger{},
	)

	// Load group configuration
	groupConf, err := factory.GetConfigLoader().LoadConfig(group)
	if err != nil {
		return &ConfigError{Group: group, Err: err}
	}

	users := groupConf.Users
	if len(users) == 0 {
		return &ConfigError{Group: group, Err: fmt.Errorf("no users found")}
	}

	// Get last index and counts
	lastIndex, err := factory.GetStorageManager().ReadLastIndex(group)
	if err != nil {
		return fmt.Errorf("failed to read last index: %w", err)
	}
	counts, err := factory.GetCountManager().GetCounts(group)
	if err != nil {
		return fmt.Errorf("failed to get counts: %w", err)
	}

	// Create strategy
	strategy, err := factory.CreateAssignmentStrategy(groupConf.Strategy)
	if err != nil {
		return &ConfigError{Group: group, Err: err}
	}

	// Create availability checker
	availChecker, err := factory.CreateAvailabilityChecker(groupConf.AvailabilityChecker)
	if err != nil {
		return &ConfigError{Group: group, Err: err}
	}

	// Select next user
	nextIndex, err := strategy.SelectNext(users, lastIndex, counts)
	if err != nil {
		return &SelectionError{Group: group, Err: err}
	}

	// Try to find an available user
	attempts := 0
	for attempts < len(users) {
		user := users[nextIndex]
		ok, err := availChecker.IsAvailable(user)
		if err != nil {
			return &AvailabilityError{User: user, Err: err}
		}
		if ok {
			if dryRun {
				fmt.Printf("[DRY RUN] Would assign to: %s\n", user)
			} else {
				fmt.Println(user)

				// Update indices and counts
				if err := factory.GetStorageManager().WriteLastIndex(group, nextIndex); err != nil {
					return fmt.Errorf("failed to write last index: %w", err)
				}
				if err := factory.GetCountManager().IncrementCount(group, user); err != nil {
					return fmt.Errorf("failed to increment count: %w", err)
				}

				// Get updated counts
				updatedCounts, err := factory.GetCountManager().GetCounts(group)
				if err != nil {
					return fmt.Errorf("failed to get updated counts: %w", err)
				}

				// Log the assignment
				if err := factory.GetAssignmentLogger().LogAssignment(group, user, groupConf.Strategy, lastIndex, nextIndex, updatedCounts); err != nil {
					return fmt.Errorf("failed to log assignment: %w", err)
				}
			}
			return nil
		}
		nextIndex = (nextIndex + 1) % len(users)
		attempts++
	}

	return &NoAvailableAssigneeError{Group: group}
}

// GetCounts retrieves the current assignment counts for a group.
// Returns the counts in the same order as users are defined in the config file.
func GetCounts(group string) (map[string]int, []string, error) {
	// Validate group exists before proceeding
	if _, err := loadAssigneeGroupConfig(group); err != nil {
		return nil, nil, &InvalidGroupError{Group: group}
	}

	counts := readCounts(group)
	if len(counts) == 0 {
		return nil, nil, fmt.Errorf("no counts found for group %s", group)
	}

	// Get the group configuration to get the user order
	groupConf, err := loadAssigneeGroupConfig(group)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load group config: %w", err)
	}

	// Ensure all users from config have an entry in counts
	for _, user := range groupConf.Users {
		if _, exists := counts[user]; !exists {
			counts[user] = 0
		}
	}

	return counts, groupConf.Users, nil
}

// ResetCounts resets the assignment counts for all users in a group to zero.
func ResetCounts(group string) error {
	groupConf, err := loadAssigneeGroupConfig(group)
	if err != nil {
		return &ConfigError{Group: group, Err: err}
	}

	counts := make(map[string]int)
	for _, user := range groupConf.Users {
		counts[user] = 0
	}

	groupDir, err := config.GetGroupDataDir(group)
	if err != nil {
		return fmt.Errorf("failed to get group data directory: %w", err)
	}

	path := filepath.Join(groupDir, "counts.json")
	data, err := json.MarshalIndent(counts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal counts: %w", err)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write counts file: %w", err)
	}

	return nil
}

// logAssignment creates a log entry for the assignment.
func logAssignment(group, user, strategy string, lastIndex, nextIndex int, counts map[string]int) error {
	// Get the group configuration to get the total number of users
	groupConf, err := loadAssigneeGroupConfig(group)
	if err != nil {
		return fmt.Errorf("failed to load group config for logging: %w", err)
	}

	logEntry := AssignmentLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		Group:      group,
		User:       user,
		Strategy:   strategy,
		LastIndex:  lastIndex,
		NextIndex:  nextIndex,
		TotalCount: len(groupConf.Users), // Use actual number of users in group
		UserCount:  counts[user],
	}

	groupDir, err := config.GetGroupDataDir(group)
	if err != nil {
		return fmt.Errorf("failed to get group data directory: %w", err)
	}

	logPath := filepath.Join(groupDir, "assignments.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}

	return nil
}

// loadAssigneeGroupConfig loads and parses the configuration for a group.
// It reads the YAML file from the configured directory and unmarshals it into an AssigneeGroupConfig.
func loadAssigneeGroupConfig(group string) (*AssigneeGroupConfig, error) {
	confPath := filepath.Join(config.Settings.Storage.ConfDir, group+".yaml")
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var groupConf AssigneeGroupConfig
	if err := yaml.Unmarshal(data, &groupConf); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	return &groupConf, nil
}

// readLastIndex reads the last assigned index for a group from the index file.
// Returns -1 if no previous assignment exists or if there's an error reading the file.
func readLastIndex(group string) int {
	groupDir, err := config.GetGroupDataDir(group)
	if err != nil {
		return -1
	}

	path := filepath.Join(groupDir, "index.log")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return -1
	}
	if len(data) == 0 {
		return -1
	}

	// Read the last line of the file
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return -1
	}

	// Get the last non-empty line
	lastLine := ""
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] != "" {
			lastLine = lines[i]
			break
		}
	}
	if lastLine == "" {
		return -1
	}

	parts := strings.Split(lastLine, "--")
	if len(parts) == 2 {
		idx := strings.TrimSpace(parts[1])
		val, err := strconv.Atoi(idx)
		if err != nil {
			return -1
		}
		return val
	}
	return -1
}

// writeLastIndex writes the last assigned index for a group to the index file.
// The index is written with a timestamp for tracking purposes.
func writeLastIndex(group string, index int) error {
	groupDir, err := config.GetGroupDataDir(group)
	if err != nil {
		return fmt.Errorf("failed to get group data directory: %w", err)
	}

	path := filepath.Join(groupDir, "index.log")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open index file: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s -- %d\n", time.Now().Format(time.RFC3339), index)
	if err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}
	return nil
}

// readCounts reads the assignment counts for all users from the counts file.
// Returns an empty map if the file doesn't exist or if there's an error reading it.
func readCounts(group string) map[string]int {
	counts := map[string]int{}
	groupDir, err := config.GetGroupDataDir(group)
	if err != nil {
		return counts
	}

	path := filepath.Join(groupDir, "counts.json")
	data, err := ioutil.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, &counts); err != nil {
			log.Printf("Warning: failed to parse counts file: %v", err)
		}
	}

	// Initialize counts for all users in the group if they don't exist
	groupConf, err := loadAssigneeGroupConfig(group)
	if err == nil {
		for _, user := range groupConf.Users {
			if _, exists := counts[user]; !exists {
				counts[user] = 0
			}
		}
	}

	return counts
}

// incrementCount increments the assignment count for a user and saves it to the counts file.
// The counts are stored in a JSON file for persistence.
func incrementCount(group, user string) error {
	counts := readCounts(group)

	// Increment the count for the specific user
	counts[user]++

	groupDir, err := config.GetGroupDataDir(group)
	if err != nil {
		return fmt.Errorf("failed to get group data directory: %w", err)
	}

	path := filepath.Join(groupDir, "counts.json")
	data, err := json.MarshalIndent(counts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal counts: %w", err)
	}
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write counts file: %w", err)
	}
	return nil
}

// GetGroupDataDir returns the data directory for a specific group.
// It creates the directory if it doesn't exist.
func GetGroupDataDir(group string) (string, error) {
	// Validate group exists before creating directory
	if _, err := loadAssigneeGroupConfig(group); err != nil {
		return "", &InvalidGroupError{Group: group}
	}

	dir := filepath.Join(config.Settings.Storage.DataDir, group)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}
