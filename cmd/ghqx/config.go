package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/domain"
	configtui "github.com/mi8bi/ghqx/internal/tui/config"
	"github.com/mi8bi/ghqx/internal/ui"
	"github.com/spf13/cobra"
)

var (
	configInitYes bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage ghqx configuration",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a default configuration file",
	Long: `Initialize a new ghqx configuration file.

Interactive mode (default):
  Prompts for each configuration value.
  Press Enter to use default values shown in [brackets].

Non-interactive mode (--yes):
  Creates config with default values immediately.

The config file will be created at:
  ~/.config/ghqx/config.toml (Linux/macOS)
  %USERPROFILE%\.config\ghqx\config.toml (Windows)

Use --config to specify a different location.`,
	RunE: runConfigInit,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long: `Display the current ghqx configuration in human-readable format.

Shows:
  - All configured roots
  - Default settings
  - Promote behavior
  - History settings`,
	RunE: runConfigShow,
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration interactively (TUI)",
	Long: `Launch an interactive TUI editor for ghqx configuration.

Features:
  - Visual field editor with descriptions
  - Real-time validation
  - Undo changes before saving
  - Japanese error messages

Keybindings:
  ↑↓ or j/k  - Navigate fields
  Enter       - Edit selected field
  Esc         - Cancel editing
  Ctrl+S      - Save configuration
  q           - Quit (warns if unsaved)
  Ctrl+Q      - Force quit without saving`,
	RunE: runConfigEdit,
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configEditCmd)

	configInitCmd.Flags().BoolVar(&configInitYes, "yes", false, "non-interactive mode: use all defaults")
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
			"Config file already exists: "+path,
		).WithHint("Delete the existing file or use a different path with --config")
	}

	var cfg *config.Config

	if configInitYes {
		// Non-interactive mode
		cfg = config.NewDefaultConfig()
		fmt.Println(ui.FormatInfo("Using default configuration"))
	} else {
		// Interactive mode
		var err error
		cfg, err = promptForConfig()
		if err != nil {
			return err
		}
	}

	// Save config
	if err := loader.Save(cfg, path); err != nil {
		return err
	}

	fmt.Print(ui.FormatSuccess("Config file created: " + path))
	fmt.Println("\n設定内容:")
	printConfigSummary(cfg)

	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	app, err := loadApp()
	if err != nil {
		return err
	}

	fmt.Println("ghqx 設定")
	fmt.Println("==================")
	printConfigSummary(app.Config)

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
	fmt.Println("ghqx 設定を対話的に作成します")
	fmt.Println("各項目でEnterを押すとデフォルト値を使用します")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	defaults := config.NewDefaultConfig()

	cfg := &config.Config{
		Roots:   make(map[string]string),
		Default: config.DefaultConfig{},
		Promote: config.PromoteConfig{},
		History: config.HistoryConfig{},
	}

	// Roots
	fmt.Println("■ ワークスペースルート")

	devPath := promptWithDefault(reader, "dev ルートのパス", defaults.Roots["dev"])
	cfg.Roots["dev"] = devPath

	releasePath := promptWithDefault(reader, "release ルートのパス", defaults.Roots["release"])
	cfg.Roots["release"] = releasePath

	sandboxPath := promptWithDefault(reader, "sandbox ルートのパス", defaults.Roots["sandbox"])
	cfg.Roots["sandbox"] = sandboxPath

	fmt.Println()

	// Default root
	fmt.Println("■ デフォルト設定")
	defaultRoot := promptWithDefault(reader, "デフォルトルート (dev/release/sandbox)", defaults.Default.Root)
	cfg.Default.Root = defaultRoot

	fmt.Println()

	// Promote settings
	fmt.Println("■ プロモート設定")
	promoteFrom := promptWithDefault(reader, "プロモート元", defaults.Promote.From)
	cfg.Promote.From = promoteFrom

	promoteTo := promptWithDefault(reader, "プロモート先", defaults.Promote.To)
	cfg.Promote.To = promoteTo

	autoGitInit := promptYesNo(reader, "Git 未管理プロジェクトに git init を実行", defaults.Promote.AutoGitInit)
	cfg.Promote.AutoGitInit = autoGitInit

	autoCommit := promptYesNo(reader, "プロモート後に自動コミット", defaults.Promote.AutoCommit)
	cfg.Promote.AutoCommit = autoCommit

	fmt.Println()

	// History settings
	fmt.Println("■ 履歴設定")
	historyEnabled := promptYesNo(reader, "undo 履歴を有効化", defaults.History.Enabled)
	cfg.History.Enabled = historyEnabled

	if historyEnabled {
		maxStr := promptWithDefault(reader, "最大履歴数", fmt.Sprintf("%d", defaults.History.Max))
		var max int
		fmt.Sscanf(maxStr, "%d", &max)
		if max <= 0 {
			max = defaults.History.Max
		}
		cfg.History.Max = max
	} else {
		cfg.History.Max = 0
	}

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
	fmt.Println("\n[Roots]")
	for name, path := range cfg.Roots {
		fmt.Printf("  %-10s = %s\n", name, path)
	}

	fmt.Println("\n[Default]")
	fmt.Printf("  root       = %s\n", cfg.Default.Root)

	fmt.Println("\n[Promote]")
	fmt.Printf("  from       = %s\n", cfg.Promote.From)
	fmt.Printf("  to         = %s\n", cfg.Promote.To)
	fmt.Printf("  auto_init  = %t\n", cfg.Promote.AutoGitInit)
	fmt.Printf("  auto_commit= %t\n", cfg.Promote.AutoCommit)

	fmt.Println("\n[History]")
	fmt.Printf("  enabled    = %t\n", cfg.History.Enabled)
	if cfg.History.Enabled {
		fmt.Printf("  max        = %d\n", cfg.History.Max)
	}
}
