// Package cmd provides the command-line interface for the autoassigner.
// It implements the main command structure and argument handling.
package cmd

import (
	"autoassigner/config"
	"autoassigner/runner"
	"autoassigner/version"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var (
	dryRun      bool
	showCounts  bool
	resetCounts bool
	configFile  string
	listGroups  bool
	showVersion bool
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "autoassigner [groupname]",
	Short: "Auto assigner CLI",
	Long: `Autoassigner is a tool for automatically assigning tasks to team members.
It uses various selection strategies and availability checks to determine the next assignee.

Example:
  autoassigner team-alpha`,
	Args: func(cmd *cobra.Command, args []string) error {
		if listGroups || showVersion {
			return nil
		}
		if len(args) != 1 {
			return fmt.Errorf("requires exactly one argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle version flag
		if showVersion {
			fmt.Println(version.String())
			return nil
		}

		if err := config.LoadConfig(configFile); err != nil {
			// Provide more user-friendly error messages for common config issues
			errMsg := err.Error()
			if strings.Contains(errMsg, "does not exist") {
				return fmt.Errorf("configuration file not found: %s\nPlease create a config.json file or specify a different path with --config", configFile)
			}
			if strings.Contains(errMsg, "invalid config") {
				return fmt.Errorf("invalid configuration: %s\nPlease check your config file format and required fields", errMsg)
			}
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Handle list-groups flag
		if listGroups {
			groups, err := config.ListGroups()
			if err != nil {
				return fmt.Errorf("failed to list groups: %w", err)
			}
			if len(groups) == 0 {
				fmt.Println("No groups found in config directory")
				return nil
			}
			sort.Strings(groups)
			fmt.Println("Available groups:")
			for _, group := range groups {
				fmt.Printf("  %s\n", group)
			}
			return nil
		}

		groupName := args[0]

		// Handle show-counts flag
		if showCounts {
			counts, orderedUsers, err := runner.GetCounts(groupName)
			if err != nil {
				if _, ok := err.(*runner.InvalidGroupError); ok {
					return fmt.Errorf("%v\nUse --list-groups to see available groups", err)
				}
				return fmt.Errorf("failed to get counts: %w", err)
			}
			fmt.Printf("Assignment counts for group %s:\n", groupName)
			for _, user := range orderedUsers {
				fmt.Printf("  %s: %d\n", user, counts[user])
			}
			return nil
		}

		// Handle reset-counts flag
		if resetCounts {
			if err := runner.ResetCounts(groupName); err != nil {
				if _, ok := err.(*runner.InvalidGroupError); ok {
					return fmt.Errorf("%v\nUse --list-groups to see available groups", err)
				}
				return fmt.Errorf("failed to reset counts: %w", err)
			}
			fmt.Printf("Successfully reset assignment counts for group %s\n", groupName)
			return nil
		}

		// Normal assignment with optional dry-run
		if err := runner.Assign(groupName, dryRun); err != nil {
			switch e := err.(type) {
			case *runner.InvalidGroupError:
				return fmt.Errorf("%v\nUse --list-groups to see available groups", e)
			case *runner.ConfigError:
				return fmt.Errorf("configuration error: %w", e)
			case *runner.SelectionError:
				return fmt.Errorf("selection error: %w", e)
			case *runner.AvailabilityError:
				return fmt.Errorf("availability error: %w", e)
			case *runner.NoAvailableAssigneeError:
				return fmt.Errorf("no available assignee: %w", e)
			default:
				return fmt.Errorf("unexpected error: %w", err)
			}
		}
		return nil
	},
	SilenceUsage:  true, // Don't show usage on error
	SilenceErrors: true, // Don't show errors (we'll handle them)
}

func init() {
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate assignment without updating logs or counts")
	rootCmd.Flags().BoolVar(&showCounts, "show-counts", false, "Display current assignment counts for the group")
	rootCmd.Flags().BoolVar(&resetCounts, "reset-counts", false, "Reset assignment counts for the group")
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "config.json", "Path to the configuration file")
	rootCmd.Flags().BoolVarP(&listGroups, "list-groups", "l", false, "List all available groups")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Display version information")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
