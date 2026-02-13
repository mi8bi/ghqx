package configtui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/domain"
)

// Model は config editor の Bubble Tea モデル
type Model struct {
	editor      *ConfigEditor
	cursor      int
	state       EditState
	editValue   string // 編集中の値
	message     string
	messageType MessageType
	width       int
	height      int
}

// MessageType はメッセージの種類
type MessageType int

const (
	MessageTypeNone MessageType = iota
	MessageTypeInfo
	MessageTypeSuccess
	MessageTypeError
)

// NewModel は新しい Model を作成する
func NewModel(cfg *config.Config, configPath string) Model {
	editor := NewConfigEditor(cfg, configPath)
	return Model{
		editor:      editor,
		cursor:      0,
		state:       EditStateList,
		messageType: MessageTypeNone,
	}
}

// Init は Bubble Tea の初期化
func (m Model) Init() tea.Cmd {
	return nil
}

// Update は Bubble Tea のイベント処理
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case saveSuccessMsg:
		m.state = EditStateList
		m.message = "設定を保存しました"
		m.messageType = MessageTypeSuccess
		m.editor.Modified = false
		return m, nil

	case saveErrorMsg:
		m.state = EditStateList
		if ghqxErr, ok := msg.err.(*domain.GhqxError); ok {
			m.message = "保存失敗: " + ghqxErr.Message
			if ghqxErr.Hint != "" {
				m.message += " (" + ghqxErr.Hint + ")"
			}
		} else {
			m.message = "保存失敗: " + msg.err.Error()
		}
		m.messageType = MessageTypeError
		return m, nil
	}

	return m, nil
}

// handleKeyPress はキーボード入力を処理する
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case EditStateList:
		return m.handleListKeys(msg)
	case EditStateEdit:
		return m.handleEditKeys(msg)
	default:
		return m, nil
	}
}

// handleListKeys はリスト表示時のキー処理
func (m Model) handleListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if m.editor.Modified {
			m.message = "未保存の変更があります。Ctrl+S で保存するか、Ctrl+Q で破棄して終了"
			m.messageType = MessageTypeInfo
			return m, nil
		}
		return m, tea.Quit

	case "ctrl+q":
		// 強制終了
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.message = ""
			m.messageType = MessageTypeNone
		}
		return m, nil

	case "down", "j":
		if m.cursor < len(m.editor.Fields)-1 {
			m.cursor++
			m.message = ""
			m.messageType = MessageTypeNone
		}
		return m, nil

	case "enter":
		// 編集モードに入る
		m.state = EditStateEdit
		m.editValue = m.editor.Fields[m.cursor].Value
		m.message = ""
		m.messageType = MessageTypeNone
		return m, nil

	case "ctrl+s":
		// 保存
		return m, m.saveConfig()
	}

	return m, nil
}

// handleEditKeys は編集モード時のキー処理
func (m Model) handleEditKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	field := m.editor.Fields[m.cursor]

	switch msg.String() {
	case "esc", "q": // Added 'q' for canceling selection
		// キャンセル
		m.state = EditStateList
		m.editValue = ""
		return m, nil

	case "enter":
		// 確定
		switch field.Type {
		case FieldTypeBool:
			// bool 型は toggle
			if m.editValue == "true" {
				m.editValue = "false"
			} else {
				m.editValue = "true"
			}
		case FieldTypeSelection:
			// Selection 型は Enter で確定し、編集モードを終了
			m.editor.UpdateField(m.cursor, m.editValue)
			m.state = EditStateList
			m.editValue = ""
			return m, nil
		}
		// For other types, "enter" confirms and exits edit mode
		m.editor.UpdateField(m.cursor, m.editValue)
		m.state = EditStateList
		m.editValue = ""
		return m, nil

	case "left", "right", " ": // Left/Right arrows and space to cycle options
		if field.Type == FieldTypeSelection {
			currentIdx := -1
			for i, opt := range field.Options {
				if opt == m.editValue {
					currentIdx = i
					break
				}
			}

			if msg.String() == "left" {
				if currentIdx <= 0 { // Cycle backward
					m.editValue = field.Options[len(field.Options)-1]
				} else {
					m.editValue = field.Options[currentIdx-1]
				}
			} else { // "right" or "space"
				if currentIdx == -1 || currentIdx >= len(field.Options)-1 { // Cycle forward
					m.editValue = field.Options[0]
				} else {
					m.editValue = field.Options[currentIdx+1]
				}
			}
			return m, nil
		}
		// Fall through for other types, if 'space' was meant for them
		if field.Type == FieldTypeString && msg.String() == " " {
			m.editValue += " "
		}
		return m, nil

	case "backspace":
		if field.Type == FieldTypeString && len(m.editValue) > 0 { // Only allow backspace for string type
			m.editValue = m.editValue[:len(m.editValue)-1]
		}
		return m, nil

	default:
		// 文字入力 (FieldTypeString のみ)
		if field.Type == FieldTypeString && len(msg.String()) == 1 {
			m.editValue += msg.String()
		}
		return m, nil
	}
}

// View は Bubble Tea の画面描画
func (m Model) View() string {
	if m.state == EditStateSaving {
		return styleTitleBar.Render("ghqx config edit - 保存中...") + "\n\n"
	}

	s := styleTitleBar.Render("ghqx config edit") + "\n\n"

	// フィールドリスト
	for i, field := range m.editor.Fields {
		s += m.renderField(field, i == m.cursor) + "\n"
	}

	// メッセージ
	if m.message != "" {
		s += "\n"
		switch m.messageType {
		case MessageTypeSuccess:
			s += styleSuccessMessage.Render("✓ " + m.message)
		case MessageTypeError:
			s += styleErrorMessage.Render("✗ " + m.message)
		default:
			s += styleInfoMessage.Render("• " + m.message)
		}
		s += "\n"
	}

	// ステータスバー
	s += "\n"
	status := ""
	if m.editor.Modified {
		status = styleModified.Render("[変更あり]") + " "
	}
	status += fmt.Sprintf("設定: %s", m.editor.ConfigPath)
	s += styleStatusBar.Render(status) + "\n"

	// ヘルプ
	s += m.renderHelp()

	return s
}

// renderField はフィールドを描画する
func (m Model) renderField(field Field, selected bool) string {
	name := styleFieldName.Render(field.Name)

	var value string
	if m.state == EditStateEdit && selected {
		// 編集中
		if field.Type == FieldTypeSelection {
			value = styleEditInput.Render(fmt.Sprintf("< %s >", m.editValue)) // Visual indicator for selection
		} else {
			value = styleEditInput.Render("> " + m.editValue + "_")
		}
	} else {
		value = styleFieldValue.Render(field.Value)
	}

	desc := styleFieldDescription.Render("  " + field.Description)

	line := fmt.Sprintf("%s  %s\n%s", name, value, desc)

	if selected {
		return styleSelectedField.Render(line)
	}
	return styleNormalField.Render(line)
}

// renderHelp はヘルプテキストを描画する
func (m Model) renderHelp() string {
	var help string

	switch m.state {
	case EditStateEdit:
		field := m.editor.Fields[m.cursor]
		switch field.Type {
		case FieldTypeBool:
			help = "Enter/Space: 切替 | Esc/q: キャンセル"
		case FieldTypeSelection:
			help = "←→: 移動 | Enter: 確定 | Esc/q: キャンセル" // Updated help for selection
		default:
			help = "Enter: 確定 | Esc/q: キャンセル | Backspace: 削除"
		}
	default:
		help = "↑↓/jk: 移動 | Enter: 編集 | Ctrl+S: 保存 | q: 終了"
		if m.editor.Modified {
			help += " | Ctrl+Q: 破棄して終了"
		}
	}

	return styleHelp.Render(help)
}
