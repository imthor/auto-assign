package runner

import (
	"autoassigner/availability"
	"autoassigner/selector"
	"fmt"
)

// ComponentFactory creates components for the runner
type ComponentFactory struct {
	configLoader     ConfigLoader
	storageManager   StorageManager
	countManager     CountManager
	assignmentLogger AssignmentLogger
}

// NewComponentFactory creates a new component factory
func NewComponentFactory(
	configLoader ConfigLoader,
	storageManager StorageManager,
	countManager CountManager,
	assignmentLogger AssignmentLogger,
) *ComponentFactory {
	return &ComponentFactory{
		configLoader:     configLoader,
		storageManager:   storageManager,
		countManager:     countManager,
		assignmentLogger: assignmentLogger,
	}
}

// CreateAssignmentStrategy creates an assignment strategy based on the strategy name
func (f *ComponentFactory) CreateAssignmentStrategy(strategy string) (AssignmentStrategy, error) {
	switch strategy {
	case "random":
		return &selector.Random{}, nil
	case "least_assigned":
		return &selector.LeastAssigned{}, nil
	case "round_robin":
		return &selector.RoundRobin{}, nil
	default:
		return nil, fmt.Errorf("unknown strategy: %s", strategy)
	}
}

// CreateAvailabilityChecker creates an availability checker based on the checker name
func (f *ComponentFactory) CreateAvailabilityChecker(checker string) (AvailabilityChecker, error) {
	switch checker {
	case "inout":
		return &availability.InOutChecker{}, nil
	case "always_available":
		return &availability.AlwaysAvailable{}, nil
	default:
		return nil, fmt.Errorf("unknown availability checker: %s", checker)
	}
}

// GetConfigLoader returns the config loader
func (f *ComponentFactory) GetConfigLoader() ConfigLoader {
	return f.configLoader
}

// GetStorageManager returns the storage manager
func (f *ComponentFactory) GetStorageManager() StorageManager {
	return f.storageManager
}

// GetCountManager returns the count manager
func (f *ComponentFactory) GetCountManager() CountManager {
	return f.countManager
}

// GetAssignmentLogger returns the assignment logger
func (f *ComponentFactory) GetAssignmentLogger() AssignmentLogger {
	return f.assignmentLogger
}
