package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mi8bi/ghqx/internal/app"
)

// RunStatus は status TUI を起動する
func RunStatus(application *app.App, showWorktree bool) error {
	model := NewStatusModel(application, showWorktree)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // 別画面モード
		tea.WithMouseCellMotion(), // マウスサポート（オプショナル）
	)

	_, err := p.Run()
	return err
}
