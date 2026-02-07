package main

import (
	"fmt"

	"github.com/mi8bi/ghqx/internal/ui"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Revert the most recent promote operation",
	Long: `Undo reverts the last promote operation by moving the project back
to its original location.

Only the most recent promote can be undone. The operation will fail if:
  - The promoted project no longer exists at its destination
  - The original location is now occupied
  - History tracking is disabled`,
	RunE: runUndo,
}

func runUndo(cmd *cobra.Command, args []string) error {
	app, err := loadApp()
	if err != nil {
		return err
	}

	record, err := app.Promote.Undo(dryRun)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Println(ui.FormatInfo("DRY RUN: Would undo:"))
		fmt.Printf("  Project: %s\n", record.ProjectName)
		fmt.Printf("  From:    %s (%s)\n", record.ToRoot, record.ToPath)
		fmt.Printf("  To:      %s (%s)\n", record.FromRoot, record.FromPath)
		return nil
	}

	fmt.Print(ui.FormatSuccess(fmt.Sprintf("Undone promote of %s", record.ProjectName)))
	fmt.Printf("  Restored to: %s\n", record.FromPath)

	return nil
}
