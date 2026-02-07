package configtui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mi8bi/ghqx/internal/config"
)

// メッセージ型

// saveSuccessMsg は保存成功メッセージ
type saveSuccessMsg struct{}

// saveErrorMsg は保存エラーメッセージ
type saveErrorMsg struct {
	err error
}

// saveConfig は設定を保存する
func (m Model) saveConfig() tea.Cmd {
	return func() tea.Msg {
		// 変更を Config に反映
		m.editor.ApplyChanges()

		// バリデーション
		if err := m.editor.Config.Validate(); err != nil {
			return saveErrorMsg{err: err}
		}

		// 保存
		loader := config.NewLoader()
		if err := loader.Save(m.editor.Config, m.editor.ConfigPath); err != nil {
			return saveErrorMsg{err: err}
		}

		return saveSuccessMsg{}
	}
}
