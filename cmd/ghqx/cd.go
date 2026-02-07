package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cdCmd = &cobra.Command{
	Use:   "cd <project-name>",
	Short: "Print the path to a project for shell cd integration",
	Long: `CD outputs the full path to a project, designed for shell integration.

Usage with shell function (add to your .bashrc or .zshrc):

  ghqx-cd() {
    local path=$(ghqx cd "$1")
    if [ -n "$path" ]; then
      cd "$path"
    fi
  }

Then use:
  ghqx-cd myproject`,
	Args: cobra.ExactArgs(1),
	RunE: runCD,
}

func runCD(cmd *cobra.Command, args []string) error {
	app, err := loadApp()
	if err != nil {
		return err
	}

	projectName := args[0]

	project, err := app.Status.FindProject(projectName)
	if err != nil {
		return err
	}

	// Output only the path for shell consumption
	fmt.Println(project.Path)

	return nil
}
