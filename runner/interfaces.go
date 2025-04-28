package runner

// AssignmentStrategy defines how tasks are assigned to team members
type AssignmentStrategy interface {
	// SelectNext chooses the next team member to assign a task to
	SelectNext(users []string, lastIndex int, counts map[string]int) (int, error)
}

// AvailabilityChecker defines how to check if a team member is available
type AvailabilityChecker interface {
	// IsAvailable checks if a team member is available for assignment
	IsAvailable(username string) (bool, error)
}

// AssignmentLogger defines how assignments are logged
type AssignmentLogger interface {
	// LogAssignment records an assignment in the log
	LogAssignment(group, user, strategy string, lastIndex, nextIndex int, counts map[string]int) error
}

// CountManager defines how assignment counts are managed
type CountManager interface {
	// GetCounts retrieves the current assignment counts for a group
	GetCounts(group string) (map[string]int, error)
	// IncrementCount increments the assignment count for a user
	IncrementCount(group, user string) error
	// ResetCounts resets the assignment counts for a group
	ResetCounts(group string) error
}

// ConfigLoader defines how group configurations are loaded
type ConfigLoader interface {
	// LoadConfig loads the configuration for a group
	LoadConfig(group string) (*AssigneeGroupConfig, error)
}

// StorageManager defines how data is stored and retrieved
type StorageManager interface {
	// GetGroupDataDir returns the data directory for a group
	GetGroupDataDir(group string) (string, error)
	// ReadLastIndex reads the last assigned index for a group
	ReadLastIndex(group string) (int, error)
	// WriteLastIndex writes the last assigned index for a group
	WriteLastIndex(group string, index int) error
}
