package main

import (
	"github.com/spf13/cobra"
	"github.com/mi8bi/ghqx/internal/i18n"
	"github.com/mi8bi/ghqx/internal/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: i18n.T("tui.command.short"),
	Long: i18n.T("tui.command.long"),
	RunE: runTUI,
}

func init() {
}

func runTUI(cmd *cobra.Command, args []string) error {
	if err := loadApp(); err != nil {
		return err
	}

	return tui.RunStatus(application)
}
