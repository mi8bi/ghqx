package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/promote"
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

// operationSuccessMsg は操作成功メッセージ
type operationSuccessMsg struct {
	message string
}

// operationErrorMsg は操作エラーメッセージ
type operationErrorMsg struct {
	err error
}

// loadProjects はプロジェクトを非同期で読み込む
func (m StatusModel) loadProjects() tea.Cmd {
	return func() tea.Msg {
		opts := status.Options{
			CheckDirty:     true,
			LoadBranch:     false,
			CountWorktrees: m.showWorktree,
		}

		projects, err := m.app.Status.GetAll(opts)
		if err != nil {
			return errorMsg{err: err}
		}

		// ProjectRow に変換
		rows := make([]ProjectRow, len(projects))
		for i, proj := range projects {
			rows[i] = NewProjectRow(proj)
		}

		return projectsLoadedMsg{projects: rows}
	}
}

// promoteProject はプロジェクトをプロモートする
func (m StatusModel) promoteProject(project domain.Project) tea.Cmd {
	return func() tea.Msg {
		// デフォルトの promote 設定を使用
		opts := promote.Options{
			ProjectName: project.Name,
			FromRoot:    string(project.Root),
			ToRoot:      m.app.Config.Promote.To,
			Force:       false,
			DryRun:      false,
			AutoGitInit: m.app.Config.Promote.AutoGitInit,
			AutoCommit:  m.app.Config.Promote.AutoCommit,
		}

		record, err := m.app.Promote.Promote(opts)
		if err != nil {
			return operationErrorMsg{err: err}
		}

		msg := "プロモート成功: " + record.ProjectName + " (" +
			string(record.FromRoot) + " → " + string(record.ToRoot) + ")"
		return operationSuccessMsg{message: msg}
	}
}

// undoPromote は直前のプロモートを取り消す
func (m StatusModel) undoPromote() tea.Cmd {
	return func() tea.Msg {
		record, err := m.app.Promote.Undo(false)
		if err != nil {
			return operationErrorMsg{err: err}
		}

		msg := "undo 成功: " + record.ProjectName + " を元に戻しました"
		return operationSuccessMsg{message: msg}
	}
}
