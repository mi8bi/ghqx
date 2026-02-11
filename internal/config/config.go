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
}

// DefaultConfig represents default settings.
type DefaultConfig struct {
	Root string `toml:"root"`
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
	// Determine OS-specific default paths using $HOME/ghqx
	homeDir, _ := os.UserHomeDir()
	var basePath string

	if homeDir != "" {
		basePath = filepath.Join(homeDir, "ghqx")
	} else {
		// Fallback for systems where home dir cannot be determined
		basePath = filepath.Join(".", "ghqx")
	}

	return &Config{
		Roots: map[string]string{
			"sandbox": filepath.Join(basePath, "sandbox"),
			"dev":     filepath.Join(basePath, "dev"),
			"release": filepath.Join(basePath, "release"),
		},
		Default: DefaultConfig{
			Root: "sandbox",
		},
	}
}
