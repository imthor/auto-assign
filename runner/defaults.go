package runner

import (
	"autoassigner/config"
	"fmt"
)

// DefaultConfigLoader implements ConfigLoader using YAML files
type DefaultConfigLoader struct{}

func (l *DefaultConfigLoader) LoadConfig(group string) (*AssigneeGroupConfig, error) {
	return loadAssigneeGroupConfig(group)
}

// DefaultStorageManager implements StorageManager using the filesystem
type DefaultStorageManager struct{}

func (m *DefaultStorageManager) GetGroupDataDir(group string) (string, error) {
	return config.GetGroupDataDir(group)
}

func (m *DefaultStorageManager) ReadLastIndex(group string) (int, error) {
	return readLastIndex(group), nil
}

func (m *DefaultStorageManager) WriteLastIndex(group string, index int) error {
	return writeLastIndex(group, index)
}

// DefaultCountManager implements CountManager using JSON files
type DefaultCountManager struct{}

func (m *DefaultCountManager) GetCounts(group string) (map[string]int, error) {
	counts := readCounts(group)
	if len(counts) == 0 {
		return nil, fmt.Errorf("no counts found for group %s", group)
	}
	return counts, nil
}

func (m *DefaultCountManager) IncrementCount(group, user string) error {
	return incrementCount(group, user)
}

func (m *DefaultCountManager) ResetCounts(group string) error {
	return ResetCounts(group)
}

// DefaultAssignmentLogger implements AssignmentLogger using JSON files
type DefaultAssignmentLogger struct{}

func (l *DefaultAssignmentLogger) LogAssignment(group, user, strategy string, lastIndex, nextIndex int, counts map[string]int) error {
	return logAssignment(group, user, strategy, lastIndex, nextIndex, counts)
}
