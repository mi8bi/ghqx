package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mi8bi/ghqx/internal/i18n"
)

var (
	// テーブルスタイル
	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true)

	styleRow = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	styleSelectedRow = lipgloss.NewStyle().
				PaddingLeft(1).
				PaddingRight(1).
				Background(lipgloss.Color("240")).
				Foreground(lipgloss.Color("230"))

	// ワークスペーススタイル (Renamed from ゾーンスタイル)
	styleSandbox = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	styleDev = lipgloss.NewStyle().
			Foreground(lipgloss.Color("81"))

	styleRelease = lipgloss.NewStyle().
			Foreground(lipgloss.Color("204"))

	// ステータススタイル (Used by getStatusStyle)
	styleClean = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	styleDirty = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	// メッセージスタイル
	styleInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	styleSuccess = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	styleWarning = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	styleError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	// UI要素スタイル
	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			MarginBottom(1)

	styleHelp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	styleFooter = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			MarginTop(1).
			PaddingTop(1)
)

// getWorkspaceStyle はワークスペースに応じたスタイルを返す (Renamed from getZoneStyle)
func getWorkspaceStyle(workspace string) lipgloss.Style { // Renamed parameter
	switch workspace {
	case "sandbox":
		return styleSandbox
	case "dev":
		return styleDev
	case "release":
		return styleRelease
	default:
		return lipgloss.NewStyle()
	}
}

// getStatusStyle はステータスに応じたスタイルを返す
func getStatusStyle(status string) lipgloss.Style {
	switch status {
	case i18n.T("status.repo.clean"):
		return styleClean
	case i18n.T("status.repo.dirty"):
		return styleDirty
	default: // 未 git 管理, "-" など
		return lipgloss.NewStyle()
	}
}

// getMessageStyle はメッセージタイプに応じたスタイルを返す
func getMessageStyle(msgType MessageType) lipgloss.Style {
	switch msgType {
	case MessageTypeSuccess:
		return styleSuccess
	case MessageTypeWarning:
		return styleWarning
	case MessageTypeError:
		return styleError
	default:
		return styleInfo
	}
}