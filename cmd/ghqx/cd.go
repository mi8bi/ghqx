package main

import (
	"fmt"

	"github.com/mi8bi/ghqx/internal/selector"
	"github.com/mi8bi/ghqx/internal/status"
	"github.com/spf13/cobra"
)

var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "", // Will be set in root.go init() after locale is determined
	Long:  "", // Will be set in root.go init() after locale is determined
	RunE:  runCD,
}

// runCD launches an interactive TUI to select a project and outputs its path.
// The path is printed to stdout and can be used with shell integration
// to change the current directory.
func runCD(cmd *cobra.Command, args []string) error {
	if err := loadApp(); err != nil {
		return err
	}

	// Load projects from the default workspace
	projects, err := loadProjectsForSelection()
	if err != nil {
		return err
	}

	// Display interactive selector
	selectedPath, err := selector.Run(projects)
	if err != nil {
		return err
	}

	// Output selected path to stdout
	if selectedPath != "" {
		fmt.Println(selectedPath)
	}

	return nil
}

// loadProjectsForSelection loads all projects from the default root
// and converts them to display format for the selector.
func loadProjectsForSelection() ([]status.ProjectDisplay, error) {
	// Check if application is initialized
	if application == nil {
		return nil, fmt.Errorf("application not initialized")
	}

	opts := status.Options{
		CheckDirty: false, // Not needed for cd operation
		LoadBranch: false, // Not needed for cd operation
	}

	// Filter by default root to reduce clutter
	defaultRoot := application.Config.GetDefaultRoot()
	rawProjects, err := application.Status.GetAll(opts, defaultRoot)
	if err != nil {
		return nil, err
	}

	// Convert to display format
	displayProjects := make([]status.ProjectDisplay, len(rawProjects))
	for i, p := range rawProjects {
		displayProjects[i] = status.NewProjectDisplay(p)
	}

	return displayProjects, nil
}
