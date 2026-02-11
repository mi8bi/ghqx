package selector

import (
	"fmt"
	"strings"

	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mi8bi/ghqx/internal/i18n"
	"github.com/mi8bi/ghqx/internal/status"
)

// Model is the Bubble Tea model for the interactive selector.
type Model struct {
	projects         []status.ProjectDisplay // All projects (original unfiltered list)
	filteredProjects []status.ProjectDisplay // Projects after filtering by text input
	textInput        textinput.Model         // Search input
	cursor           int
	selected         string // FullPath of the selected project
	quitting         bool
}

// NewModel creates a new model for the selector.
func NewModel(projects []status.ProjectDisplay) Model {
	ti := textinput.New()
	ti.Placeholder = i18n.T("selector.search.placeholder")
	ti.CharLimit = 100
	ti.Width = 20
	// ti.Focus() // Moved to Init() to ensure it's called after model setup

	return Model{
		projects:         projects,
		filteredProjects: projects, // Initially all projects are filtered
		textInput:        ti,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return m.textInput.Focus() // Return the command to focus the text input
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			} else if len(m.filteredProjects) > 0 {
				m.cursor = len(m.filteredProjects) - 1
			}
			// Navigation keys do not update text input or trigger filterProjects directly.
			// Cursor adjustment is enough.
		case "down":
			if m.cursor < len(m.filteredProjects)-1 {
				m.cursor++
			} else if len(m.filteredProjects) > 0 {
				m.cursor = 0
			}
			// Navigation keys do not update text input or trigger filterProjects directly.
			// Cursor adjustment is enough.
		case "enter":
			if m.cursor >= 0 && m.cursor < len(m.filteredProjects) {
				m.selected = m.filteredProjects[m.cursor].FullPath
			}
			return m, tea.Quit
		default: // Regular character input for filtering
			m.textInput, cmd = m.textInput.Update(msg)
			m.filterProjects()
		}
	case tea.WindowSizeMsg:
		// Handle window size changes if necessary
		// For now, no specific action is needed.
	}

	// Adjust cursor if filtered projects list is shorter than cursor position
	// This applies after any potential filtering or initial load.
	if m.cursor >= len(m.filteredProjects) {
		if len(m.filteredProjects) > 0 {
			m.cursor = len(m.filteredProjects) - 1
		} else {
			m.cursor = 0 // No projects, cursor at 0
		}
	}

	return m, cmd
}

// filterProjects filters the projects based on the text input value.
func (m *Model) filterProjects() {
	query := strings.ToLower(m.textInput.Value())
	if query == "" {
		m.filteredProjects = m.projects
		m.cursor = 0 // Reset cursor when filter is cleared
		return
	}

	var newFiltered []status.ProjectDisplay
	for _, p := range m.projects {
		if strings.Contains(strings.ToLower(p.Repo), query) ||
			strings.Contains(strings.ToLower(p.FullPath), query) {
			newFiltered = append(newFiltered, p)
		}
	}
	m.filteredProjects = newFiltered
	m.cursor = 0 // Reset cursor to top of filtered list
}

// View implements tea.Model.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	s := strings.Builder{}
	s.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("36")).Render("> ") +
		lipgloss.NewStyle().Bold(true).Render(i18n.T("selector.title")) + "\n")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("─", 50)) + "\n\n")

	// Render search input with better formatting
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Render("Filter: "))
	s.WriteString(m.textInput.View())

	// Show match count
	matchCount := len(m.filteredProjects)
	if matchCount > 0 || m.textInput.Value() == "" {
		s.WriteString("  ")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Render(
			fmt.Sprintf("[%d/%d]", matchCount, len(m.projects))))
	}
	s.WriteString("\n\n")

	if len(m.filteredProjects) == 0 && m.textInput.Value() != "" {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("  No matches found\n"))
	}

	// Render projects
	for i, p := range m.filteredProjects {
		cursor := " "
		if m.cursor == i {
			cursor = "❯"
		}

		// Create highlighted line with better formatting
		line := fmt.Sprintf("%s %s", cursor, p.Repo)

		if m.cursor == i {
			// Highlight selected item with background
			s.WriteString(lipgloss.NewStyle().
				Background(lipgloss.Color("24")).
				Foreground(lipgloss.Color("255")).
				Render(line))
		} else {
			s.WriteString(line)
		}
		s.WriteString("\n")
	}

	// Help text
	if len(m.filteredProjects) > 0 {
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Render("  ↑/↓ to navigate • Enter to select • Esc to cancel"))
	}

	return s.String()
}

// Run displays the interactive selector and returns the full path of the selected project.
func Run(projects []status.ProjectDisplay) (string, error) {
	if len(projects) == 0 {
		return "", nil // No projects to select
	}

	model := NewModel(projects)                     // Get the initial model
	p := tea.NewProgram(model, tea.WithAltScreen()) // Pass model only, Init() handles textinput.Focus()
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	if m, ok := finalModel.(Model); ok && !m.quitting {
		return m.selected, nil
	}
	return "", nil // User quit or nothing selected
}
