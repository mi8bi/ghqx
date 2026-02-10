package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/i18n" // Add this import
	configtui "github.com/mi8bi/ghqx/internal/tui/config"
	"github.com/mi8bi/ghqx/internal/ui"
)

var (
	configInitYes bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: i18n.T("config.command.short"),
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: i18n.T("config.init.command.short"),
	Long: i18n.T("config.init.command.long"),
	RunE: runConfigInit,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: i18n.T("config.show.command.short"),
	Long: i18n.T("config.show.command.long"),
	RunE: runConfigShow,
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: i18n.T("config.edit.command.short"),
	Long: i18n.T("config.edit.command.long"),
	RunE: runConfigEdit,
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configEditCmd)
	
	configInitCmd.Flags().BoolVar(&configInitYes, "yes", false, i18n.T("config.init.flag.yes"))
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	loader := config.NewLoader()

	// Determine config path
	path := configPath
	if path == "" {
		defaultPath, err := config.GetDefaultConfigPath()
		if err != nil {
			return err
		}
		path = defaultPath
	}

	// Check if config already exists
	if _, err := os.Stat(path); err == nil {
		return domain.NewError(
			domain.ErrCodeConfigInvalid,
			fmt.Sprintf(i18n.T("config.error.fileAlreadyExists"), path),
		).WithHint("Delete the existing file or use a different path with --config")
	}

	var cfg *config.Config

	if configInitYes {
		// Non-interactive mode
		cfg = config.NewDefaultConfig()
		fmt.Println(ui.FormatInfo(i18n.T("config.init.useDefault")))
	} else {
		// Interactive mode
		var err error
		cfg, err = promptForConfig()
		if err != nil {
			return err
		}
	}

	// Create root directories before saving config
	fmt.Println(ui.FormatInfo(i18n.T("config.init.creatingDirs")))
	if err := config.EnsureRootDirectories(cfg); err != nil {
		return err
	}

	// Save config
	if err := loader.Save(cfg, path); err != nil {
		return err
	}

	fmt.Print(ui.FormatSuccess(i18n.T("config.init.fileCreated") + ": " + path))
	fmt.Println("\n" + i18n.T("config.init.summaryHeader"))
	printConfigSummary(cfg)

	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	if err := loadApp(); err != nil {
		return err
	}

	fmt.Println(i18n.T("config.show.title"))
	fmt.Println("==================")
	printConfigSummary(application.Config)

	return nil
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	loader := config.NewLoader()

	// Load existing config
	cfg, err := loader.Load(configPath)
	if err != nil {
		return err
	}

	// Determine config path for saving
	savePath := configPath
	if savePath == "" {
		savePath, err = config.GetDefaultConfigPath()
		if err != nil {
			return err
		}
	}

	// Launch TUI editor
	return configtui.Run(cfg, savePath)
}

// promptForConfig は対話的に設定を入力する
func promptForConfig() (*config.Config, error) {
	fmt.Println(i18n.T("config.prompt.intro1"))
	fmt.Println(i18n.T("config.prompt.intro2"))
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	defaults := config.NewDefaultConfig()

	cfg := &config.Config{
		Roots:   make(map[string]string),
		Default: config.DefaultConfig{},
	}

	// Roots
	fmt.Println(i18n.T("config.prompt.section.roots"))
	
	devPath := promptWithDefault(reader, i18n.T("config.prompt.path.dev"), defaults.Roots["dev"])
	cfg.Roots["dev"] = devPath

	releasePath := promptWithDefault(reader, i18n.T("config.prompt.path.release"), defaults.Roots["release"])
	cfg.Roots["release"] = releasePath

	sandboxPath := promptWithDefault(reader, i18n.T("config.prompt.path.sandbox"), defaults.Roots["sandbox"])
	cfg.Roots["sandbox"] = sandboxPath

	fmt.Println()

	// Default root
	fmt.Println(i18n.T("config.prompt.section.default"))
	defaultRoot := promptWithDefault(reader, i18n.T("config.prompt.defaultRoot"), defaults.Default.Root)
	cfg.Default.Root = defaultRoot

	defaultLanguage := promptWithDefault(reader, i18n.T("config.prompt.defaultLanguage"), defaults.Default.Language)
	cfg.Default.Language = defaultLanguage

	fmt.Println()

	return cfg, nil
}

// promptWithDefault は入力を促し、空の場合はデフォルト値を返す
func promptWithDefault(reader *bufio.Reader, prompt, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "" {
		return defaultValue
	}
	return input
}

// promptYesNo は yes/no を入力させる
func promptYesNo(reader *bufio.Reader, prompt string, defaultValue bool) bool {
	defaultStr := "n"
	if defaultValue {
		defaultStr = "y"
	}

	fmt.Printf("%s (y/n) [%s]: ", prompt, defaultStr)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	if input == "" {
		return defaultValue
	}

	return input == "y" || input == "yes"
}

// printConfigSummary は設定の要約を表示する
func printConfigSummary(cfg *config.Config) {
	fmt.Println("\n" + i18n.T("config.summary.section.roots"))
	for name, path := range cfg.Roots {
		fmt.Printf("  %-10s = %s\n", name, path)
	}

	fmt.Println("\n" + i18n.T("config.summary.section.default"))
	fmt.Printf("  root       = %s\n", cfg.Default.Root)
	fmt.Printf("  language   = %s\n", cfg.Default.Language)
}
