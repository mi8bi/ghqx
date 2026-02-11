package main

import (
	"fmt"

	"github.com/mi8bi/ghqx/internal/i18n"
	"github.com/mi8bi/ghqx/internal/selector"
	"github.com/mi8bi/ghqx/internal/status"
	"github.com/spf13/cobra"
)

var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: i18n.T("cd.command.short"), // Re-using old i18n key for `cd`
	Long:  i18n.T("cd.command.long"),  // Re-using old i18n key for `cd`
	RunE: runCD,
}

func init() {
}

func runCD(cmd *cobra.Command, args []string) error {
	if err := loadApp(); err != nil {
		return err
	}

	opts := status.Options{}
	// Filter projects by the default root
	projects, err := application.Status.GetAll(opts, application.Config.GetDefaultRoot())
	if err != nil {
		return err
	}

	// Convert domain.Project to ProjectDisplay
	displayProjects := make([]status.ProjectDisplay, len(projects))
	for i, p := range projects {
		displayProjects[i] = status.NewProjectDisplay(p)
	}

	// Always use interactive selection
	selectedPath, err := selector.Run(displayProjects)
	if err != nil {
		return err // Error from Bubble Tea program
	}
	if selectedPath != "" {
		fmt.Println(selectedPath)
	}
	return nil // If nothing selected, just exit without error
}
