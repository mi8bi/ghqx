package app

import (
	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/status"
)

// App represents the application and its dependencies.
type App struct {
	Config  *config.Config
	Status  *status.Service
}

// New creates a new application instance.
func New(cfg *config.Config) *App {
	return &App{
		Config:  cfg,
		Status:  status.NewService(cfg),
	}
}

// NewFromConfigPath creates an app instance by loading config from a path.
func NewFromConfigPath(configPath string) (*App, error) {
	loader := config.NewLoader()
	cfg, err := loader.Load(configPath)
	if err != nil {
		return nil, err
	}

	return New(cfg), nil
}
