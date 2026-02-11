package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/i18n" // Added missing import
	"github.com/mi8bi/ghqx/internal/ui"
	"github.com/spf13/cobra"
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
	Use:           "ghqx",
	Short:         "ghqx", // Will be set in init() after locale is determined
	Long:          "ghqx", // Will be set in init() after locale is determined
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

		// Determine locale precedence: Env Var > OS Language > Default (Japanese)
		targetLocale := getOSLanguageLocale()

		// Override with environment variable if set
		if lang := os.Getenv("GHQX_LANG"); lang != "" {
			switch lang {
			case "en", "en_US":
				targetLocale = i18n.LocaleEN
			case "ja", "ja_JP":
				targetLocale = i18n.LocaleJA
			}
		}
		i18n.SetLocale(targetLocale)
		return nil
	},
}

func init() {
	// Initialize locale BEFORE setting command descriptions
	initLocale()

	// Now set descriptions using the correct locale
	rootCmd.Short = i18n.T("root.command.short")
	rootCmd.Long = i18n.T("root.command.long")

	// Set all subcommand descriptions after locale is initialized
	cdCmd.Short = i18n.T("cd.command.short")
	cdCmd.Long = i18n.T("cd.command.long")

	statusCmd.Short = i18n.T("status.command.short")
	statusCmd.Long = i18n.T("status.command.long")

	configCmd.Short = i18n.T("config.command.short")
	configInitCmd.Short = i18n.T("config.init.command.short")
	configInitCmd.Long = i18n.T("config.init.command.long")
	configShowCmd.Short = i18n.T("config.show.command.short")
	configShowCmd.Long = i18n.T("config.show.command.long")
	configEditCmd.Short = i18n.T("config.edit.command.short")
	configEditCmd.Long = i18n.T("config.edit.command.long")

	getCmd.Short = i18n.T("get.command.short")
	getCmd.Long = i18n.T("get.command.long")

	doctorCmd.Short = i18n.T("doctor.command.short")
	doctorCmd.Long = i18n.T("doctor.command.long")

	cleanCmd.Short = i18n.T("clean.command.short")
	cleanCmd.Long = i18n.T("clean.command.long")

	modeCmd.Short = i18n.T("mode.command.short")
	modeCmd.Long = i18n.T("mode.command.long")

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", i18n.T("root.flag.config"))

	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(cdCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(modeCmd)
}

// initLocale initializes the locale before any command descriptions are rendered
func initLocale() {
	// Determine locale precedence: Env Var > OS Language > Default (Japanese)
	targetLocale := getOSLanguageLocale()

	// Override with environment variable if set
	if lang := os.Getenv("GHQX_LANG"); lang != "" {
		switch lang {
		case "en", "en_US":
			targetLocale = i18n.LocaleEN
		case "ja", "ja_JP":
			targetLocale = i18n.LocaleJA
		}
	}
	i18n.SetLocale(targetLocale)
}

// getOSLanguageLocale determines the locale based on OS environment.
func getOSLanguageLocale() i18n.Locale {
	// Check environment variables in priority order

	// 1. Check LC_ALL (highest priority)
	if lang := os.Getenv("LC_ALL"); lang != "" {
		if strings.Contains(strings.ToLower(lang), "ja") {
			return i18n.LocaleJA
		} else if strings.Contains(strings.ToLower(lang), "en") {
			return i18n.LocaleEN
		}
	}

	// 2. Check LANG (Unix/Linux/macOS standard)
	if lang := os.Getenv("LANG"); lang != "" {
		if strings.Contains(strings.ToLower(lang), "ja") {
			return i18n.LocaleJA
		} else if strings.Contains(strings.ToLower(lang), "en") {
			return i18n.LocaleEN
		}
	}

	// 3. Check LANGUAGE (Linux)
	if lang := os.Getenv("LANGUAGE"); lang != "" {
		// LANGUAGE can contain colon-separated list
		if strings.Contains(strings.ToLower(lang), "ja") {
			return i18n.LocaleJA
		} else if strings.Contains(strings.ToLower(lang), "en") {
			return i18n.LocaleEN
		}
	}

	// Default to Japanese if none found
	return i18n.LocaleJA
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
