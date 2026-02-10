package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mi8bi/ghqx/internal/status"
)

// メッセージ型定義

// projectsLoadedMsg はプロジェクト読み込み完了メッセージ
type projectsLoadedMsg struct {
	projects []ProjectRow
}

// errorMsg はエラーメッセージ
type errorMsg struct {
	err error
}

// loadProjects はプロジェクトを非同期で読み込む
func (m StatusModel) loadProjects() tea.Cmd {
	return func() tea.Msg {
		opts := status.Options{
			CheckDirty: true,
			LoadBranch: false,
		}

		projects, err := m.app.Status.GetAll(opts)
		if err != nil {
			return errorMsg{err: err}
		}

		// ProjectRow に変換
		rows := make([]ProjectRow, len(projects))
		for i, proj := range projects {
			rows[i] = NewProjectRow(status.NewProjectDisplay(proj))
		}

		return projectsLoadedMsg{projects: rows}
	}
}
