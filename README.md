# Autoassigner

A CLI tool for automatically assigning tasks to team members based on various selection strategies and availability checks.

## Features

- Multiple selection strategies:
  - Round Robin: Cycles through team members in order
  - Random: Randomly selects a team member
  - Least Assigned: Selects the team member with the fewest assignments
- Availability checking:
  - In/Out status: Checks external API for member availability
  - Always Available: Simple implementation that always returns available
- Configuration via YAML files
- Assignment tracking and history
- Group management and validation
- Dry run mode for testing assignments
- Assignment count tracking and reset
- Extensible component system for custom implementations

## Installation

```bash
go install github.com/justinthomas/autoassigner@latest
```

## Usage

```bash
# Basic usage
autoassigner [groupname]

# List all available groups (both commands do the same thing)
autoassigner --list-groups
autoassigner -l

# Show assignment counts for a group
autoassigner [groupname] --show-counts

# Reset assignment counts for a group
autoassigner [groupname] --reset-counts

# Simulate assignment without updating logs or counts
autoassigner [groupname] --dry-run

# Use a custom configuration file (both commands do the same thing)
autoassigner --config /path/to/config.json [groupname]
autoassigner -c /path/to/config.json [groupname]

# Display version information (both commands do the same thing)
autoassigner --version
autoassigner -v
```

## Configuration

1. Create a `config.json` file:
```json
{
    "storage": {
        "data_dir": "var/data",
        "conf_dir": "etc"
    },
    "availability": {
        "inout_api_url_prefix": "https://api.example.com/status/",
        "inout_unavailable_statuses": ["OOO", "AWAY"]
    }
}
```

2. Create group configuration files in the `etc` directory:
```yaml
strategy: round_robin
availability_checker: inout
users:
  - user1
  - user2
  - user3
```

## Extending the System

The system is designed to be extensible through a component-based architecture. You can implement custom versions of any component by implementing the appropriate interface:

### Assignment Strategies

```go
type CustomStrategy struct{}

func (s *CustomStrategy) SelectNext(users []string, lastIndex int, counts map[string]int) (int, error) {
    // Custom selection logic
}
```

### Availability Checkers

```go
type CustomChecker struct{}

func (c *CustomChecker) IsAvailable(username string) (bool, error) {
    // Custom availability logic
}
```

## Error Handling

The tool provides clear error messages for common issues:

- Invalid group names:
  ```
  Error: group team-alpa does not exist
  Use --list-groups to see available groups
  ```

- Missing configuration:
  ```
  Error: configuration file not found: config.json
  Please create a config.json file or specify a different path with --config
  ```

- Invalid configuration:
  ```
  Error: invalid configuration: data_dir is required in storage configuration
  Please check your config file format and required fields
  ```

## Data Storage

The tool maintains several types of data files:

- `var/data/<group>/assignments.log`: Assignment history
- `var/data/<group>/counts.json`: Assignment counts
- `var/data/<group>/index.log`: Assignment indices

## Development

1. Clone the repository
2. Install dependencies: `go mod download`
3. Build: `go build`
4. Run tests: `go test ./...`

## Testing

The codebase includes comprehensive tests for:
- Selection strategies
- Availability checking
- Configuration loading
- Error handling
- Group validation
- Component interfaces

Run tests with:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request 