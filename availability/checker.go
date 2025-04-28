// Package availability provides different implementations for checking team member availability.
// It includes:
// - In/Out status checker: Checks external API for member availability
// - Always Available: Simple implementation that always returns available
package availability

// Checker defines the interface for checking team member availability.
// Each implementation must provide a way to determine if a team member is available.
type Checker interface {
	// IsAvailable checks if a team member is available for assignment.
	// Parameters:
	//   - username: The username of the team member to check
	// Returns:
	//   - bool: True if the team member is available, false otherwise
	//   - error: Any error that occurred during the check
	IsAvailable(username string) (bool, error)
}
