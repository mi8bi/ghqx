package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mi8bi/ghqx/internal/git"
	"github.com/mi8bi/ghqx/internal/ui"
	"github.com/spf13/cobra"
)

var worktreeCmd = &cobra.Command{
	Use:   "worktree <project-name>",
	Short: "List git worktrees for a project",
	Long: `Worktree lists all git worktrees for a given project.

Git worktrees allow you to check out multiple branches simultaneously
in different directories. This command shows all worktrees for a project.`,
	Args: cobra.ExactArgs(1),
	RunE: runWorktree,
}

func runWorktree(cmd *cobra.Command, args []string) error {
	app, err := loadApp()
	if err != nil {
		return err
	}

	projectName := args[0]

	// Find the project
	project, err := app.Status.FindProject(projectName)
	if err != nil {
		return err
	}

	// Check if it's a git repository
	if !project.HasGit {
		fmt.Fprintf(os.Stderr, "%s%s is not a git repository\n",
			ui.FormatWarning(""), project.Name)
		return nil
	}

	// List worktrees
	gitClient := git.NewClient()
	worktrees, err := gitClient.ListWorktrees(project.Path)
	if err != nil {
		return err
	}

	if len(worktrees) == 0 {
		fmt.Printf("No worktrees found for %s\n", project.Name)
		return nil
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(worktrees)
	}

	fmt.Printf("Worktrees for %s (%d total):\n\n", project.Name, len(worktrees))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "PATH\tBRANCH\tSTATUS")
	fmt.Fprintln(w, "----\t------\t------")

	for _, wt := range worktrees {
		status := ""
		if wt.Bare {
			status = "bare"
		} else if wt.Locked {
			status = "locked"
		} else {
			status = "active"
		}

		branch := wt.Branch
		if branch == "" {
			branch = "(detached)"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\n", wt.Path, branch, status)
	}

	return nil
}
