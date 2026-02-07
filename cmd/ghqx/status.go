package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mi8bi/ghqx/internal/status"
	"github.com/mi8bi/ghqx/internal/tui"
	"github.com/spf13/cobra"

	"github.com/mi8bi/ghqx/internal/domain"
)

var (
	statusVerbose   bool
	statusWorktrees bool
	statusTUI       bool
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the state of all projects across all roots",
	Long: `Status quickly visualizes workspace state.

Projects are classified by zone:
  sandbox  - flat directory in sandbox root
  dev      - repository in dev root
  release  - repository in release root

Additional information:
  - Git managed or not
  - Dirty/clean status
  - Worktree count (with --worktrees)

TUI mode (Terminal UI):
  --tui flag launches an interactive terminal interface
  with keyboard navigation and operations.`,
	RunE: runStatus,
}

func init() {
	statusCmd.Flags().BoolVarP(&statusVerbose, "verbose", "v", false, "show detailed information including paths")
	statusCmd.Flags().BoolVarP(&statusWorktrees, "worktrees", "w", false, "count git worktrees")
	statusCmd.Flags().BoolVar(&statusTUI, "tui", false, "launch interactive TUI mode")
}

func runStatus(cmd *cobra.Command, args []string) error {
	app, err := loadApp()
	if err != nil {
		return err
	}

	// TUI モードの場合
	if statusTUI {
		return tui.RunStatus(app, statusWorktrees)
	}

	// CLI モード（既存の実装）
	opts := status.Options{
		CheckDirty:     true,
		LoadBranch:     false, // Fast mode for CLI
		CountWorktrees: statusWorktrees,
		RootFilter:     rootFlag,
	}

	projects, err := app.Status.GetAll(opts)
	if err != nil {
		return err
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(projects)
	}

	// Convert []domain.Project to []any for table output functions
	projectsAny := make([]any, len(projects))
	for i, p := range projects {
		projectsAny[i] = p
	}

	if statusVerbose {
		return outputVerboseTable(projectsAny)
	}

	return outputCompactTable(projectsAny)
}

func outputCompactTable(projects []any) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "NAME\tZONE\tGIT\tSTATUS\tWORKTREES")
	fmt.Fprintln(w, "----\t----\t---\t------\t---------")

	for _, p := range projects {
		proj, ok := p.(domain.Project)
		if !ok {
			// This should not happen if `projectsAny` is populated correctly.
			// Handle error or skip if necessary.
			continue
		}

		hasGit := "no"
		if proj.HasGit {
			hasGit = "yes"
		}

		gitStatus := "clean"
		if proj.Dirty {
			gitStatus = "dirty"
		}
		if !proj.HasGit {
			gitStatus = "-"
		}

		worktrees := "-"
		if proj.HasGit && proj.WorktreeCount > 0 {
			worktrees = fmt.Sprintf("%d", proj.WorktreeCount)
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			proj.Name, proj.Zone, hasGit, gitStatus, worktrees)
	}

	return nil
}

func outputVerboseTable(projects []any) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "NAME\tZONE\tROOT\tPATH\tGIT\tSTATUS\tWORKTREES")
	fmt.Fprintln(w, "----\t----\t----\t----\t---\t------\t---------")

	for _, p := range projects {
		proj, ok := p.(domain.Project)
		if !ok {
			// This should not happen if `projectsAny` is populated correctly.
			// Handle error or skip if necessary.
			continue
		}

		hasGit := "no"
		if proj.HasGit {
			hasGit = "yes"
		}

		gitStatus := "clean"
		if proj.Dirty {
			gitStatus = "dirty"
		}
		if !proj.HasGit {
			gitStatus = "-"
		}

		worktrees := "-"
		if proj.HasGit && proj.WorktreeCount > 0 {
			worktrees = fmt.Sprintf("%d", proj.WorktreeCount)
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			proj.Name, proj.Zone, proj.Root, proj.Path, hasGit, gitStatus, worktrees)
	}

	return nil
}
