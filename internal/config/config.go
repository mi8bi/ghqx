package config

import (
	"os"
	"path/filepath"

	"github.com/mi8bi/ghqx/internal/domain"
)

// Config represents the ghqx configuration.
type Config struct {
	Roots   map[string]string `toml:"roots"`
	Default DefaultConfig     `toml:"default"`
	Promote PromoteConfig     `toml:"promote"`
	History HistoryConfig     `toml:"history"`
}

// DefaultConfig represents default settings.
type DefaultConfig struct {
	Root string `toml:"root"`
}

// PromoteConfig represents promote behavior settings.
type PromoteConfig struct {
	From        string `toml:"from"`
	To          string `toml:"to"`
	AutoGitInit bool   `toml:"auto_git_init"`
	AutoCommit  bool   `toml:"auto_commit"`
}

// HistoryConfig represents history tracking settings.
type HistoryConfig struct {
	Enabled bool `toml:"enabled"`
	Max     int  `toml:"max"`
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if len(c.Roots) == 0 {
		return domain.ErrConfigNoRoots
	}

	if c.Default.Root != "" {
		if _, exists := c.Roots[c.Default.Root]; !exists {
			return domain.ErrConfigInvalidDefaultRoot
		}
	}

	if c.Promote.From != "" {
		if _, exists := c.Roots[c.Promote.From]; !exists {
			return domain.NewError(domain.ErrCodeConfigInvalid, "Promote 'from' root does not exist").
				WithHint("Set promote.from to one of the defined roots")
		}
	}

	if c.Promote.To != "" {
		if _, exists := c.Roots[c.Promote.To]; !exists {
			return domain.NewError(domain.ErrCodeConfigInvalid, "Promote 'to' root does not exist").
				WithHint("Set promote.to to one of the defined roots")
		}
	}

	return nil
}

// GetRoot returns the path for a given root name.
func (c *Config) GetRoot(name string) (string, bool) {
	path, exists := c.Roots[name]
	return path, exists
}

// GetDefaultRoot returns the default root name.
func (c *Config) GetDefaultRoot() string {
	if c.Default.Root != "" {
		return c.Default.Root
	}
	// Fallback to first root if no default is set
	for name := range c.Roots {
		return name
	}
	return ""
}

// NewDefaultConfig creates a default configuration.
// This is the single source of truth for default values.
func NewDefaultConfig() *Config {
	// Determine OS-specific default paths
	homeDir, _ := os.UserHomeDir()
	var basePath string

	if homeDir != "" {
		basePath = filepath.Join(homeDir, "src")
	} else {
		basePath = filepath.Join("C:", "src")
	}

	return &Config{
		Roots: map[string]string{
			"dev":     filepath.Join(basePath, "ghq-dev"),
			"release": filepath.Join(basePath, "ghq-release"),
			"sandbox": filepath.Join(basePath, "sandbox"),
		},
		Default: DefaultConfig{
			Root: "dev",
		},
		Promote: PromoteConfig{
			From:        "sandbox",
			To:          "dev",
			AutoGitInit: true,
			AutoCommit:  false,
		},
		History: HistoryConfig{
			Enabled: true,
			Max:     50,
		},
	}
}
