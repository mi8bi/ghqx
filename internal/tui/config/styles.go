package configtui

import "github.com/charmbracelet/lipgloss"

var (
	// タイトルスタイル
	styleTitleBar = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	// フィールドリストスタイル
	styleFieldName = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("81")).
			Width(20)

	styleFieldValue = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	styleFieldDescription = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Italic(true)

	// 選択行スタイル
	styleSelectedField = lipgloss.NewStyle().
				Background(lipgloss.Color("240")).
				Padding(0, 1)

	styleNormalField = lipgloss.NewStyle().
				Padding(0, 1)

	// 編集モードスタイル
	styleEditInput = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			Bold(true)

	// ステータスバースタイル
	styleStatusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	styleModified = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	// メッセージスタイル
	styleInfoMessage = lipgloss.NewStyle().
				Foreground(lipgloss.Color("86"))

	styleSuccessMessage = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")).
				Bold(true)

	styleErrorMessage = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Bold(true)

	// ヘルプスタイル
	styleHelp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			MarginTop(1).
			PaddingTop(1)
)
