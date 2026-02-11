package app

import (
	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/status"
)

// App represents the main application instance and holds all necessary dependencies.
// It provides access to application configuration and project status management.
type App struct {
	// Config holds the application configuration including workspace roots
	Config *config.Config
	// Status provides project scanning and status management services
	Status *status.Service
}

// New creates a new application instance with the given configuration.
// This is the primary constructor for App and initializes all services.
func New(cfg *config.Config) *App {
	return &App{
		Config: cfg,
		Status: status.NewService(cfg),
	}
}

// NewFromConfigPath creates an app instance by loading configuration from a file path.
// This is a convenience constructor that handles loading the config file.
// If configPath is empty, it will use the default config location.
func NewFromConfigPath(configPath string) (*App, error) {
	loader := config.NewLoader()
	cfg, err := loader.Load(configPath)
	if err != nil {
		return nil, err
	}

	return New(cfg), nil
}
