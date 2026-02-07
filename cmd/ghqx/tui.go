package main

import (
	"github.com/mi8bi/ghqx/internal/tui"
	"github.com/spf13/cobra"
)

var (
	tuiWorktrees bool
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive TUI (Terminal UI)",
	Long: `TUI launches an interactive terminal interface for ghqx.

Features:
  - Visual project list with keyboard navigation
  - Direct promote and undo operations
  - Real-time status updates
  - Japanese error messages

Keybindings:
  ↑↓ or j/k  - Navigate projects
  Enter       - Promote selected project
  u           - Undo last promotion
  r           - Refresh project list
  q or Ctrl+C - Quit`,
	RunE: runTUI,
}

func init() {
	tuiCmd.Flags().BoolVarP(&tuiWorktrees, "worktrees", "w", false, "show worktree counts")
}

func runTUI(cmd *cobra.Command, args []string) error {
	app, err := loadApp()
	if err != nil {
		return err
	}

	return tui.RunStatus(app, tuiWorktrees)
}
