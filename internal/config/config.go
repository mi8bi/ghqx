package config

import (
	"os"
	"path/filepath"

	"github.com/mi8bi/ghqx/internal/domain"
)

// Config represents the ghqx application configuration.
// It defines the workspace roots and default settings.
type Config struct {
	// Roots maps root names to their filesystem paths
	// Example: {"dev": "/home/user/ghqx/dev", "sandbox": "/home/user/ghqx/sandbox"}
	Roots map[string]string `toml:"roots"`
	// Default specifies default settings like which root to use
	Default DefaultConfig `toml:"default"`
}

// DefaultConfig represents default application settings.
type DefaultConfig struct {
	// Root specifies the default workspace root name
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

// GetRoot returns the filesystem path for the given root name.
// It returns empty string and false if the root name doesn't exist.
func (c *Config) GetRoot(name string) (string, bool) {
	path, exists := c.Roots[name]
	return path, exists
}

// GetDefaultRoot returns the default root name.
// If no explicit default is configured, it returns the first available root.
// Returns empty string if no roots are configured.
func (c *Config) GetDefaultRoot() string {
	// Use explicit default if set
	if c.Default.Root != "" {
		return c.Default.Root
	}
	// Fallback to first available root
	for name := range c.Roots {
		return name
	}
	return ""
}

// NewDefaultConfig creates a default configuration with standard workspace roots.
// Creates three roots: sandbox, dev, and release under $HOME/ghqx with sandbox as default.
// This is the single source of truth for default configuration values.
func NewDefaultConfig() *Config {
	// Determine base path using user's home directory
	homeDir, _ := os.UserHomeDir()
	var basePath string

	if homeDir != "" {
		basePath = filepath.Join(homeDir, "ghqx")
	} else {
		// Fallback for systems where home directory cannot be determined
		basePath = filepath.Join(".", "ghqx")
	}

	return &Config{
		Roots: map[string]string{
			"sandbox": filepath.Join(basePath, "sandbox"), // For experimental/temp projects
			"dev":     filepath.Join(basePath, "dev"),     // For development projects
			"release": filepath.Join(basePath, "release"), // For stable/release projects
		},
		Default: DefaultConfig{
			Root: "sandbox", // Default to sandbox for new repositories
		},
	}
}
