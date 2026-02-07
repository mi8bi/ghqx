package main

import (
	"github.com/mi8bi/ghqx/internal/app"
	"github.com/spf13/cobra"
)

var (
	configPath string
	jsonOutput bool
	dryRun     bool
	rootFlag   string
)

var rootCmd = &cobra.Command{
	Use:   "ghqx",
	Short: "ghqx - ghq-compatible workspace lifecycle manager",
	Long: `ghqx extends ghq by managing multiple workspaces (dev/release/sandbox)
and supporting lifecycle operations such as promote, undo, and status.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file path")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "show what would be done without doing it")
	rootCmd.PersistentFlags().StringVar(&rootFlag, "root", "", "target specific root")

	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(promoteCmd)
	rootCmd.AddCommand(undoCmd)
	rootCmd.AddCommand(cdCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(worktreeCmd)
	rootCmd.AddCommand(tuiCmd)
}

// loadApp is a helper to load the app with config.
func loadApp() (*app.App, error) {
	application, err := app.NewFromConfigPath(configPath)
	if err != nil {
		return nil, err
	}
	return application, nil
}
