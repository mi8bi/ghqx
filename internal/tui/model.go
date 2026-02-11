package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/i18n"
)

// StatusModel は status TUI の Bubble Tea モデル

type StatusModel struct {

	app          *app.App

	projects     []ProjectRow // ProjectDisplay を含む

	cursor       int

	viewState    ViewState

	message      *Message

	err          error

	width        int

	height       int

	showDetail   bool // 詳細表示モード

}



// NewStatusModel は新しい StatusModel を作成する

func NewStatusModel(application *app.App) StatusModel {

	return StatusModel{

		app:          application,

		projects:     []ProjectRow{},

		cursor:       0,

		viewState:    ViewStateLoading,

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

			Text: fmt.Sprintf(i18n.T("status.message.projectsLoaded"), len(m.projects)),

			Type: MessageTypeInfo,

		}

		return m, nil



	case errorMsg:

		m.viewState = ViewStateError

		m.err = msg.err

		if ghqxErr, ok := msg.err.(*domain.GhqxError); ok {

			m.message = &Message{

				Text: ghqxErr.Message,

				Type: MessageTypeError,

				Hint: ghqxErr.Hint,

			}

		} else {

			m.message = &Message{

				Text: i18n.T("status.message.errorOccurred"),

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

			m.message = nil

		}

		return m, nil



	case "down", "j":

		if m.cursor < len(m.projects)-1 {

			m.cursor++

			m.message = nil

		}

		return m, nil



	case "d":

		m.showDetail = !m.showDetail

		return m, nil



	case "r":

		m.viewState = ViewStateLoading

		m.message = &Message{

			Text: i18n.T("status.message.reloading"),

			Type: MessageTypeInfo,

		}

		return m, m.loadProjects()

	}



	return m, nil

}



// View は Bubble Tea の画面描画

func (m StatusModel) View() string {

	if m.viewState == ViewStateLoading {

		return styleTitle.Render(i18n.T("status.title.loading")) + "\n\n"

	}



	if m.viewState == ViewStateError {

		s := styleTitle.Render(i18n.T("status.title.error")) + "\n\n"

		if m.message != nil {

			s += styleError.Render("✗ "+m.message.Text) + "\n"

			if m.message.Hint != "" {

				s += styleHelp.Render(i18n.T("ui.error.hintPrefix")+": "+m.message.Hint) + "\n"

			}

		}

		s += "\n" + styleHelp.Render(i18n.T("status.help.error"))

		return s

	}



	return m.renderList()

}



// renderList はプロジェクトリストを描画する

func (m StatusModel) renderList() string {

	s := styleTitle.Render(i18n.T("status.title.list")) + "\n\n"



	if m.showDetail && len(m.projects) > 0 {

		return m.renderDetailView()

	}



	// ヘッダー

	repoHeader := lipgloss.NewStyle().Width(30).Align(lipgloss.Left).Render(i18n.T("status.header.name"))

	workspaceHeader := lipgloss.NewStyle().Width(10).Align(lipgloss.Left).Render(i18n.T("status.header.workspace")) // Renamed from zoneHeader and updated i18n key

	gitManagedHeader := lipgloss.NewStyle().Width(10).Align(lipgloss.Left).Render(i18n.T("status.header.gitManaged"))

	statusHeader := lipgloss.NewStyle().Width(8).Align(lipgloss.Left).Render(i18n.T("status.header.status"))



	header := fmt.Sprintf("%s %s %s %s", repoHeader, workspaceHeader, gitManagedHeader, statusHeader)

	s += styleHeader.Render(header) + "\n"



	// プロジェクト行

	for i, row := range m.projects {

		s += m.renderProjectRow(row, i == m.cursor) + "\n"

	}



	// フッター: メッセージ

	s += "\n"

	if m.message != nil {

		msgStyle := getMessageStyle(m.message.Type)

		s += msgStyle.Render(m.message.Text)

		if m.message.Hint != "" {

			s += "\n" + styleHelp.Render(i18n.T("ui.error.hintPrefix")+": "+m.message.Hint)

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

	proj := row.RawProject



	s := styleTitle.Render(i18n.T("status.title.detail")) + "\n\n"



	// 基本情報

	s += styleSandbox.Render(i18n.T("status.detail.basicInfo")) + "\n"

	s += fmt.Sprintf("  %s:     %s\n", i18n.T("status.detail.name"), row.Repo)

	s += fmt.Sprintf("  %s:     %s\n", i18n.T("status.detail.path"), row.FullPath)

	s += fmt.Sprintf("  %s:   %s\n", i18n.T("status.detail.workspace"), getWorkspaceStyle(row.Workspace).Render(row.Workspace)) // Updated to row.Workspace and getWorkspaceStyle

	s += fmt.Sprintf("  %s:   %s\n", i18n.T("status.detail.root"), proj.Root)

	s += "\n"



	// Git 情報

	s += styleSandbox.Render(i18n.T("status.detail.gitInfo")) + "\n"

	s += fmt.Sprintf("  %s:  %s\n", i18n.T("status.detail.gitManaged"), row.GitManaged)



	s += fmt.Sprintf("  %s:     %s\n", i18n.T("status.detail.status"), getStatusStyle(row.Status).Render(row.Status))



	if proj.Branch != "" {

		s += fmt.Sprintf("  %s: %s\n", i18n.T("status.detail.branch"), proj.Branch)

	}

	s += "\n"



	// メッセージ

	if m.message != nil {

		msgStyle := getMessageStyle(m.message.Type)

		s += msgStyle.Render(m.message.Text)

		if m.message.Hint != "" {

			s += "\n" + styleHelp.Render(i18n.T("ui.error.hintPrefix")+": "+m.message.Hint)

		}

		s += "\n\n"

	}



	// ヘルプ

	s += styleFooter.Render(m.renderHelp())



	return s

}



// renderProjectRow はプロジェクト行を描画する

func (m StatusModel) renderProjectRow(row ProjectRow, selected bool) string {

	repoCell := lipgloss.NewStyle().Width(30).Align(lipgloss.Left).Render(row.Repo)

	workspaceCell := lipgloss.NewStyle().Width(10).Align(lipgloss.Left).Render(row.Workspace) // Renamed from zoneCell and updated to row.Workspace

	gitManagedCell := lipgloss.NewStyle().Width(10).Align(lipgloss.Left).Render(row.GitManaged)

	statusCell := lipgloss.NewStyle().Width(8).Align(lipgloss.Left).Render(row.Status)



	line := fmt.Sprintf("%s %s %s %s", repoCell, workspaceCell, gitManagedCell, statusCell)

	// 選択行はハイライト

	if selected {

		return styleSelectedRow.Render("> " + line)

	}

	return styleRow.Render("  " + line)

}



// renderHelp はヘルプテキストを描画する

func (m StatusModel) renderHelp() string {

	help := i18n.T("status.help.main")

	return help

}
