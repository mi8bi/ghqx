package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/mi8bi/ghqx/internal/domain"
)

// Loader handles configuration file discovery and loading.
type Loader struct{}

// NewLoader creates a new config loader.
func NewLoader() *Loader {
	return &Loader{}
}

// Load finds and loads the configuration file.
// Search order:
// 1. configPath argument (if provided)
// 2. GHQX_CONFIG environment variable
// 3. $XDG_CONFIG_HOME/ghqx/config.toml
// 4. ~/.config/ghqx/config.toml
// 5. ~/.ghqx.toml
func (l *Loader) Load(configPath string) (*Config, error) {
	path, err := l.findConfigPath(configPath)
	if err != nil {
		return nil, err
	}

	return l.loadFromPath(path)
}

// Save writes configuration to the specified path.
// If path is empty, uses the default config location.
func (l *Loader) Save(cfg *Config, configPath string) error {
	path := configPath
	if path == "" {
		var err error
		path, err = GetDefaultConfigPath()
		if err != nil {
			return err
		}
	}

	// Validate before saving
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return domain.ErrFSCreateDir(err)
	}

	// Write config file
	f, err := os.Create(path)
	if err != nil {
		return domain.NewErrorWithCause(
			domain.ErrCodeFSError,
			"Failed to create config file",
			err,
		)
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	if err := enc.Encode(cfg); err != nil {
		return domain.NewErrorWithCause(
			domain.ErrCodeConfigInvalid,
			"Failed to write config",
			err,
		)
	}

	return nil
}

// GetDefaultConfigPath returns the default config file path.
// This is a package-level function accessible from other parts of the config package.
func GetDefaultConfigPath() (string, error) {
	// Prefer XDG_CONFIG_HOME
	if xdgHome := os.Getenv("XDG_CONFIG_HOME"); xdgHome != "" {
		return filepath.Join(xdgHome, "ghqx", "config.toml"), nil
	}

	// Fallback to ~/.config/ghqx/config.toml
	home, err := os.UserHomeDir()
	if err != nil {
		return "", domain.NewErrorWithCause(
			domain.ErrCodeConfigInvalid,
			"Cannot determine home directory",
			err,
		)
	}

	return filepath.Join(home, ".config", "ghqx", "config.toml"), nil
}

// findConfigPath returns the first existing config file path.
func (l *Loader) findConfigPath(explicitPath string) (string, error) {
	// 1. Explicit path via flag
	if explicitPath != "" {
		if _, err := os.Stat(explicitPath); err == nil {
			return explicitPath, nil
		}
		return "", domain.ErrConfigNotFoundAt(explicitPath)
	}

	// 2. GHQX_CONFIG environment variable
	if envPath := os.Getenv("GHQX_CONFIG"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}
	}

	// 3. XDG_CONFIG_HOME/ghqx/config.toml
	if xdgHome := os.Getenv("XDG_CONFIG_HOME"); xdgHome != "" {
		path := filepath.Join(xdgHome, "ghqx", "config.toml")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// 4. ~/.config/ghqx/config.toml
	home, err := os.UserHomeDir()
	if err == nil {
		path := filepath.Join(home, ".config", "ghqx", "config.toml")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// 5. ~/.ghqx.toml
	if err == nil {
		path := filepath.Join(home, ".ghqx.toml")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", domain.ErrConfigNotFoundAny
}

// loadFromPath reads and parses a TOML config file.
func (l *Loader) loadFromPath(path string) (*Config, error) {
	var cfg Config

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, domain.ErrConfigInvalidTOML(err).WithInternal("path: " + path)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
