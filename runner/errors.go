package runner

import "fmt"

type ConfigError struct {
	Group string
	Err   error
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("configuration error for group %s: %v", e.Group, e.Err)
}

type SelectionError struct {
	Group string
	Err   error
}

func (e *SelectionError) Error() string {
	return fmt.Sprintf("selection error for group %s: %v", e.Group, e.Err)
}

type AvailabilityError struct {
	User string
	Err  error
}

func (e *AvailabilityError) Error() string {
	return fmt.Sprintf("availability check error for user %s: %v", e.User, e.Err)
}

type NoAvailableAssigneeError struct {
	Group string
}

func (e *NoAvailableAssigneeError) Error() string {
	return fmt.Sprintf("no available assignee found for group %s", e.Group)
}

type InvalidGroupError struct {
	Group string
}

func (e *InvalidGroupError) Error() string {
	return fmt.Sprintf("group %s does not exist", e.Group)
}
