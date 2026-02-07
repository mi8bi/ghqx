package configtui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mi8bi/ghqx/internal/config"
)

// Run は config editor TUI を起動する
func Run(cfg *config.Config, configPath string) error {
	model := NewModel(cfg, configPath)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	_, err := p.Run()
	return err
}
