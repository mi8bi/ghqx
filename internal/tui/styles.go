package tui

import "github.com/charmbracelet/lipgloss"

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

	// ゾーンスタイル
	styleSandbox = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	styleDev = lipgloss.NewStyle().
			Foreground(lipgloss.Color("81"))

	styleRelease = lipgloss.NewStyle().
			Foreground(lipgloss.Color("204"))

	// ステータススタイル
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

// getZoneStyle はゾーンに応じたスタイルを返す
func getZoneStyle(zone string) lipgloss.Style {
	switch zone {
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
