// Package selector provides different strategies for selecting team members for task assignment.
// It includes implementations for:
// - Round Robin: Cycles through team members in order
// - Random: Randomly selects a team member
// - Least Assigned: Selects the team member with the fewest assignments
package selector

// Selector defines the interface for different selection strategies.
// Each strategy must implement the SelectNext method to determine the next assignee.
type Selector interface {
	// SelectNext chooses the next team member to assign a task to.
	// Parameters:
	//   - users: List of available team members
	//   - lastIndex: Index of the last assigned team member
	//   - counts: Map of assignment counts for each team member
	// Returns:
	//   - int: Index of the selected team member
	//   - error: Any error that occurred during selection
	SelectNext(users []string, lastIndex int, counts map[string]int) (int, error)
}
