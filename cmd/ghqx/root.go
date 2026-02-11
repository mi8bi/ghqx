package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/i18n" // Added missing import
	"github.com/mi8bi/ghqx/internal/ui"
)

var (
	configPath string

	// Global application instance
	application *app.App
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, ui.FormatDetailedError(err))
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ghqx",
	Short: i18n.T("root.command.short"),
	Long:  i18n.T("root.command.long"),
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// If the command is 'config init', we don't try to load existing config.
		// We set a default locale for its output.
		if cmd == configInitCmd {
			i18n.SetLocale(i18n.LocaleJA) // Default to Japanese for init output
			return nil
		}

		// For all other commands, load the app and determine locale.
		if err := loadApp(); err != nil {
			return err
		}

		// Determine locale precedence: Env Var > Config > Default
		targetLocale := i18n.LocaleJA // Default fallback

		if lang := os.Getenv("GHQX_LANG"); lang != "" {
			switch lang {
			case "en", "en_US":
				targetLocale = i18n.LocaleEN
			case "ja", "ja_JP":
				targetLocale = i18n.LocaleJA
			}
		} else if application.Config != nil && application.Config.Default.Language != "" { // Check application.Config for nil
			switch application.Config.Default.Language {
			case "en":
				targetLocale = i18n.LocaleEN
			case "ja":
				targetLocale = i18n.LocaleJA
			}
		}
		i18n.SetLocale(targetLocale)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", i18n.T("root.flag.config"))

	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(cdCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(modeCmd)
}

// loadApp is a helper to load the app with config and set the global application variable.
func loadApp() error {
	var err error
	application, err = app.NewFromConfigPath(configPath)
	if err != nil {
		return err
	}
	return nil
}
