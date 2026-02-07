package main

import (
	"fmt"

	"github.com/mi8bi/ghqx/internal/promote"
	"github.com/mi8bi/ghqx/internal/ui"
	"github.com/spf13/cobra"
)

var (
	promoteFrom       string
	promoteTo         string
	promoteForce      bool
	promoteAutoGit    bool
	promoteAutoCommit bool
)

var promoteCmd = &cobra.Command{
	Use:   "promote <project-name>",
	Short: "Move a project from one root to another",
	Long: `Promote moves a project from one workspace to another.

By default, promotes from sandbox to dev (configurable in config.toml).

Safety checks:
  - Refuses to promote dirty git repositories (unless --force)
  - Checks for conflicts at destination
  - Records operation in history for undo`,
	Args: cobra.ExactArgs(1),
	RunE: runPromote,
}

func init() {
	promoteCmd.Flags().StringVar(&promoteFrom, "from", "", "source root (default: from config)")
	promoteCmd.Flags().StringVar(&promoteTo, "to", "", "destination root (default: from config)")
	promoteCmd.Flags().BoolVar(&promoteForce, "force", false, "promote even if repository is dirty")
	promoteCmd.Flags().BoolVar(&promoteAutoGit, "git-init", false, "initialize git if not present (default: from config)")
	promoteCmd.Flags().BoolVar(&promoteAutoCommit, "auto-commit", false, "commit after promote (default: from config)")
}

func runPromote(cmd *cobra.Command, args []string) error {
	app, err := loadApp()
	if err != nil {
		return err
	}

	projectName := args[0]

	// Use config defaults if not specified
	from := promoteFrom
	if from == "" {
		from = app.Config.Promote.From
	}

	to := promoteTo
	if to == "" {
		to = app.Config.Promote.To
	}

	autoGit := promoteAutoGit
	if !cmd.Flags().Changed("git-init") {
		autoGit = app.Config.Promote.AutoGitInit
	}

	autoCommit := promoteAutoCommit
	if !cmd.Flags().Changed("auto-commit") {
		autoCommit = app.Config.Promote.AutoCommit
	}

	opts := promote.Options{
		ProjectName: projectName,
		FromRoot:    from,
		ToRoot:      to,
		Force:       promoteForce,
		DryRun:      dryRun,
		AutoGitInit: autoGit,
		AutoCommit:  autoCommit,
	}

	record, err := app.Promote.Promote(opts)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Println(ui.FormatInfo("DRY RUN: Would promote:"))
		fmt.Printf("  From: %s (%s)\n", record.FromRoot, record.FromPath)
		fmt.Printf("  To:   %s (%s)\n", record.ToRoot, record.ToPath)
		return nil
	}

	fmt.Print(ui.FormatSuccess(fmt.Sprintf("Promoted %s from %s to %s",
		projectName, from, to)))
	fmt.Printf("  New location: %s\n", record.ToPath)

	return nil
}
