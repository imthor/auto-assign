// Package selector provides different strategies for selecting team members for task assignment.
package selector

import (
	"fmt"
)

// LeastAssigned implements the Selector interface to choose team members
// who have been assigned the fewest tasks. This strategy helps maintain
// a balanced workload across the team by prioritizing members with fewer
// assignments.
type LeastAssigned struct{}

// SelectNext chooses the next team member to assign a task to based on
// the number of previous assignments. It selects the team member with
// the lowest assignment count.
//
// Parameters:
//   - users: List of available team members
//   - lastIndex: Index of the last assigned team member (not used in this strategy)
//   - counts: Map of assignment counts for each team member
//
// Returns:
//   - int: Index of the selected team member
//   - error: Any error that occurred during selection
//
// Example:
//
//	users := []string{"alice", "bob", "charlie"}
//	counts := map[string]int{"alice": 2, "bob": 1, "charlie": 3}
//	index, err := leastAssigned.SelectNext(users, -1, counts)
//	// index will be 1 (bob) as they have the fewest assignments
func (l *LeastAssigned) SelectNext(users []string, lastIndex int, counts map[string]int) (int, error) {
	if len(users) == 0 {
		return -1, fmt.Errorf("empty users list")
	}
	min := 1<<31 - 1 // Initialize with maximum possible integer value
	index := 0
	for i, u := range users {
		if counts[u] < min {
			min = counts[u]
			index = i
		}
	}
	return index, nil
}
