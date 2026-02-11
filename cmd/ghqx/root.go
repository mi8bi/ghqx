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
		// Special case: config init doesn't require an existing app instance
		if cmd == configInitCmd {
			i18n.SetLocale(i18n.LocaleJA) // Default to Japanese for init output
			return nil
		}

		// Load app configuration for all other commands
		if err := loadApp(); err != nil {
			return err
		}

		// Determine and set the locale based on environment
		locale := determineLocale()
		i18n.SetLocale(locale)
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

	versionCmd.Short = i18n.T("version.command.short")
	versionCmd.Long = i18n.T("version.command.long")

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
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(modeCmd)
}

// initLocale initializes the locale before any command descriptions are rendered.
func initLocale() {
	locale := determineLocale()
	i18n.SetLocale(locale)
}

// determineLocale determines the appropriate locale based on environment variables.
// Priority: GHQX_LANG environment variable > OS environment > Default (Japanese)
func determineLocale() i18n.Locale {
	// 1. Check GHQX_LANG environment variable (explicit app setting)
	if lang := os.Getenv("GHQX_LANG"); lang != "" {
		if locale := parseLocaleFromEnv(lang); locale != "" {
			return locale
		}
	}

	// 2. Determine from OS environment variables
	return getOSLanguageLocale()
}

// parseLocaleFromEnv converts a language string to locale.
func parseLocaleFromEnv(lang string) i18n.Locale {
	switch lang {
	case "en", "en_US":
		return i18n.LocaleEN
	case "ja", "ja_JP":
		return i18n.LocaleJA
	default:
		return ""
	}
}

// getOSLanguageLocale determines the locale based on OS environment variables.
// Checks in standard order: LC_ALL > LANG > LANGUAGE > Default (Japanese)
func getOSLanguageLocale() i18n.Locale {
	// 1. Check LC_ALL (highest priority in Unix/Linux/macOS)
	if lang := os.Getenv("LC_ALL"); lang != "" {
		if locale := matchLocaleString(lang); locale != "" {
			return locale
		}
	}

	// 2. Check LANG (Unix/Linux/macOS standard)
	if lang := os.Getenv("LANG"); lang != "" {
		if locale := matchLocaleString(lang); locale != "" {
			return locale
		}
	}

	// 3. Check LANGUAGE (Linux, can contain colon-separated list)
	if lang := os.Getenv("LANGUAGE"); lang != "" {
		if locale := matchLocaleString(lang); locale != "" {
			return locale
		}
	}

	// Default to Japanese if no OS locale is detected
	return i18n.LocaleJA
}

// matchLocaleString checks if a string contains language hints and returns appropriate locale.
func matchLocaleString(str string) i18n.Locale {
	lowerStr := strings.ToLower(str)
	if strings.Contains(lowerStr, "ja") {
		return i18n.LocaleJA
	}
	if strings.Contains(lowerStr, "en") {
		return i18n.LocaleEN
	}
	return ""
}

// loadApp initializes the global application instance with configuration.
// It loads config from the specified path (or default location if not specified).
func loadApp() error {
	var err error
	application, err = app.NewFromConfigPath(configPath)
	if err != nil {
		return err
	}
	return nil
}
