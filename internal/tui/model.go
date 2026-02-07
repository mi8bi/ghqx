package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/domain"
)

// StatusModel は status TUI の Bubble Tea モデル
type StatusModel struct {
	app          *app.App
	projects     []ProjectRow
	cursor       int
	viewState    ViewState
	message      *Message
	err          error
	width        int
	height       int
	showWorktree bool
	showDetail   bool // 詳細表示モード
}

// NewStatusModel は新しい StatusModel を作成する
func NewStatusModel(application *app.App, showWorktree bool) StatusModel {
	return StatusModel{
		app:          application,
		projects:     []ProjectRow{},
		cursor:       0,
		viewState:    ViewStateLoading,
		showWorktree: showWorktree,
	}
}

// Init は Bubble Tea の初期化処理
func (m StatusModel) Init() tea.Cmd {
	return m.loadProjects()
}

// Update は Bubble Tea のイベント処理
func (m StatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case projectsLoadedMsg:
		m.projects = msg.projects
		m.viewState = ViewStateList
		m.message = &Message{
			Text: fmt.Sprintf("%d 個のプロジェクトを読み込みました", len(m.projects)),
			Type: MessageTypeInfo,
		}
		return m, nil

	case errorMsg:
		m.viewState = ViewStateError
		m.err = msg.err
		// ユーザー向けメッセージを抽出
		if ghqxErr, ok := msg.err.(*domain.GhqxError); ok {
			m.message = &Message{
				Text: ghqxErr.Message,
				Type: MessageTypeError,
				Hint: ghqxErr.Hint,
			}
		} else {
			m.message = &Message{
				Text: "エラーが発生しました",
				Type: MessageTypeError,
			}
		}
		return m, nil

	case operationSuccessMsg:
		m.message = &Message{
			Text: msg.message,
			Type: MessageTypeSuccess,
		}
		// 成功後は再読み込み
		return m, m.loadProjects()

	case operationErrorMsg:
		// エラーメッセージを表示するが状態は維持
		if ghqxErr, ok := msg.err.(*domain.GhqxError); ok {
			m.message = &Message{
				Text: ghqxErr.Message,
				Type: MessageTypeError,
				Hint: ghqxErr.Hint,
			}
		} else {
			m.message = &Message{
				Text: msg.err.Error(),
				Type: MessageTypeError,
			}
		}
		return m, nil
	}

	return m, nil
}

// handleKeyPress はキーボード入力を処理する
func (m StatusModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.message = nil // メッセージをクリア
		}
		return m, nil

	case "down", "j":
		if m.cursor < len(m.projects)-1 {
			m.cursor++
			m.message = nil
		}
		return m, nil

	case "d":
		// 詳細表示トグル
		m.showDetail = !m.showDetail
		return m, nil

	case "enter":
		// プロモート実行
		if m.viewState == ViewStateList && len(m.projects) > 0 {
			row := m.projects[m.cursor]
			if row.CanPromote {
				return m, m.promoteProject(row.Project)
			} else {
				m.message = &Message{
					Text: "プロモートできません",
					Type: MessageTypeWarning,
					Hint: row.PromoteHint,
				}
				return m, nil
			}
		}

	case "u":
		// undo 実行
		return m, m.undoPromote()

	case "r":
		// 再読み込み
		m.viewState = ViewStateLoading
		m.message = &Message{
			Text: "再読み込み中...",
			Type: MessageTypeInfo,
		}
		return m, m.loadProjects()
	}

	return m, nil
}

// View は Bubble Tea の画面描画
func (m StatusModel) View() string {
	if m.viewState == ViewStateLoading {
		return styleTitle.Render("ghqx status - 読み込み中...") + "\n\n"
	}

	if m.viewState == ViewStateError {
		s := styleTitle.Render("ghqx status - エラー") + "\n\n"
		if m.message != nil {
			s += styleError.Render("✗ "+m.message.Text) + "\n"
			if m.message.Hint != "" {
				s += styleHelp.Render("ヒント: "+m.message.Hint) + "\n"
			}
		}
		s += "\n" + styleHelp.Render("q: 終了 | r: 再試行")
		return s
	}

	return m.renderList()
}

// renderList はプロジェクトリストを描画する
func (m StatusModel) renderList() string {
	s := styleTitle.Render("ghqx status - プロジェクト一覧") + "\n\n"

	// 詳細表示モードの場合
	if m.showDetail && len(m.projects) > 0 {
		return m.renderDetailView()
	}

	// ヘッダー
	header := fmt.Sprintf("%-30s %-10s %-5s %-8s", "名前", "ゾーン", "Git", "状態")
	if m.showWorktree {
		header += fmt.Sprintf(" %-10s", "Worktree")
	}
	s += styleHeader.Render(header) + "\n"

	// プロジェクト行
	for i, row := range m.projects {
		line := m.renderProjectRow(row, i == m.cursor)
		s += line + "\n"
	}

	// フッター: メッセージ
	s += "\n"
	if m.message != nil {
		msgStyle := getMessageStyle(m.message.Type)
		s += msgStyle.Render(m.message.Text)
		if m.message.Hint != "" {
			s += "\n" + styleHelp.Render("ヒント: "+m.message.Hint)
		}
		s += "\n"
	}

	// ヘルプ
	s += styleFooter.Render(m.renderHelp())

	return s
}

// renderDetailView は選択中のプロジェクトの詳細を表示する
func (m StatusModel) renderDetailView() string {
	row := m.projects[m.cursor]
	proj := row.Project

	s := styleTitle.Render("ghqx status - プロジェクト詳細") + "\n\n"

	// 基本情報
	s += styleSandbox.Render("■ 基本情報") + "\n"
	s += fmt.Sprintf("  名前:     %s\n", proj.Name)
	s += fmt.Sprintf("  パス:     %s\n", proj.Path)
	s += fmt.Sprintf("  ゾーン:   %s\n", getZoneStyle(string(proj.Zone)).Render(string(proj.Zone)))
	s += fmt.Sprintf("  ルート:   %s\n", proj.Root)
	s += "\n"

	// Git 情報
	s += styleSandbox.Render("■ Git 情報") + "\n"
	if proj.HasGit {
		s += "  Git管理:  yes\n"

		status := "clean"
		statusStyle := styleClean
		if proj.Dirty {
			status = "dirty"
			statusStyle = styleDirty
		}
		s += fmt.Sprintf("  状態:     %s\n", statusStyle.Render(status))

		if proj.Branch != "" {
			s += fmt.Sprintf("  ブランチ: %s\n", proj.Branch)
		}

		if m.showWorktree {
			s += fmt.Sprintf("  Worktree: %d\n", proj.WorktreeCount)
		}
	} else {
		s += "  Git管理:  no\n"
	}
	s += "\n"

	// プロモート情報
	s += styleSandbox.Render("■ プロモート") + "\n"
	if row.CanPromote {
		s += styleSuccess.Render("  プロモート可能") + "\n"
		s += fmt.Sprintf("  実行先: %s\n", m.app.Config.Promote.To)
	} else {
		s += styleWarning.Render("  プロモート不可") + "\n"
		if row.PromoteHint != "" {
			s += fmt.Sprintf("  理由: %s\n", row.PromoteHint)
		}
	}
	s += "\n"

	// メッセージ
	if m.message != nil {
		msgStyle := getMessageStyle(m.message.Type)
		s += msgStyle.Render(m.message.Text)
		if m.message.Hint != "" {
			s += "\n" + styleHelp.Render("ヒント: "+m.message.Hint)
		}
		s += "\n\n"
	}

	// ヘルプ
	s += styleFooter.Render("d: リスト表示 | Enter: プロモート | u: undo | r: 再読み込み | q: 終了")

	return s
}

// renderProjectRow はプロジェクト行を描画する
func (m StatusModel) renderProjectRow(row ProjectRow, selected bool) string {
	proj := row.Project

	// 基本情報
	name := proj.Name
	if len(name) > 28 {
		name = name[:25] + "..."
	}

	zone := string(proj.Zone)
	zoneStyled := getZoneStyle(zone).Render(fmt.Sprintf("%-10s", zone))

	gitStatus := "no"
	if proj.HasGit {
		gitStatus = "yes"
	}

	status := "clean"
	statusStyle := styleClean
	if proj.Dirty {
		status = "dirty"
		statusStyle = styleDirty
	}

	line := fmt.Sprintf("%-30s %s %-5s %s",
		name,
		zoneStyled,
		gitStatus,
		statusStyle.Render(fmt.Sprintf("%-8s", status)),
	)

	if m.showWorktree && proj.HasGit {
		line += fmt.Sprintf(" %-10d", proj.WorktreeCount)
	}

	// 選択行はハイライト
	if selected {
		return styleSelectedRow.Render("> " + line)
	}
	return styleRow.Render("  " + line)
}

// renderHelp はヘルプテキストを描画する
func (m StatusModel) renderHelp() string {
	help := "↑↓/jk: 移動 | d: 詳細 | Enter: プロモート | u: undo | r: 再読み込み | q: 終了"
	return help
}
