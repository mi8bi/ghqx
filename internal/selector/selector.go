package selector

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mi8bi/ghqx/internal/i18n"
	"github.com/mi8bi/ghqx/internal/status"
)

// Model is the Bubble Tea model for the interactive selector.
type Model struct {
	projects []status.ProjectDisplay
	cursor   int
	selected string // FullPath of the selected project
	quitting bool
}

// NewModel creates a new model for the selector.
func NewModel(projects []status.ProjectDisplay) Model {
	return Model{projects: projects}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.projects)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor >= 0 && m.cursor < len(m.projects) {
				m.selected = m.projects[m.cursor].FullPath
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	s := strings.Builder{}
	s.WriteString(lipgloss.NewStyle().Bold(true).Render(i18n.T("selector.title")) + "\n\n")

	for i, p := range m.projects {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		
		// Render with styles (simplified for now, actual width/align will be adjusted later)
		line := fmt.Sprintf("%s %-30s %-10s %s",
			cursor,
			p.Repo,
			p.Zone,
			p.FullPath,
		)

		if m.cursor == i {
			s.WriteString(lipgloss.NewStyle().Background(lipgloss.Color("240")).Foreground(lipgloss.Color("230")).Render(line))
		} else {
			s.WriteString(line)
		}
		s.WriteString("\n")
	}
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1).Render(i18n.T("selector.help")))
	return s.String()
}

// Run displays the interactive selector and returns the full path of the selected project.
func Run(projects []status.ProjectDisplay) (string, error) {
	if len(projects) == 0 {
		return "", nil // No projects to select
	}

	p := tea.NewProgram(NewModel(projects))
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	if m, ok := finalModel.(Model); ok && !m.quitting {
		return m.selected, nil
	}
	return "", nil // User quit or nothing selected
}
