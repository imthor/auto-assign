// Package selector provides different strategies for selecting team members for task assignment.
package selector

import (
	"fmt"
)

// RoundRobin implements the Selector interface to choose team members
// in a circular order. This strategy ensures that each team member
// is selected in turn, providing a fair distribution of tasks.
type RoundRobin struct{}

// SelectNext chooses the next team member to assign a task to by
// cycling through the list of team members in order. It uses the
// lastIndex to determine the next member in the sequence.
//
// Parameters:
//   - users: List of available team members
//   - lastIndex: Index of the last assigned team member
//   - counts: Map of assignment counts for each team member (not used in this strategy)
//
// Returns:
//   - int: Index of the selected team member
//   - error: Any error that occurred during selection
//
// Example:
//
//	users := []string{"alice", "bob", "charlie"}
//	lastIndex := 1
//	index, err := roundRobin.SelectNext(users, lastIndex, nil)
//	// index will be 2 (charlie) as it's the next in the sequence
func (r *RoundRobin) SelectNext(users []string, lastIndex int, counts map[string]int) (int, error) {
	if len(users) == 0 {
		return -1, fmt.Errorf("empty users list")
	}
	return (lastIndex + 1) % len(users), nil
}
