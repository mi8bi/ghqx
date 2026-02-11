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
	ti.Width = 50
	ti.Prompt = "" // Remove the default "> " prompt
	ti.Focus()     // Focus the text input by default

	return Model{
		projects:         projects,
		filteredProjects: projects, // Initially all projects are shown
		textInput:        ti,
		cursor:           0,
		quitting:         false,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyStr := msg.String()

		// Handle global keys first
		if keyStr == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		if keyStr == "esc" {
			m.quitting = true
			return m, tea.Quit
		}

		// Selection
		if keyStr == "enter" {
			if m.cursor >= 0 && m.cursor < len(m.filteredProjects) {
				m.selected = m.filteredProjects[m.cursor].FullPath
			}
			return m, tea.Quit
		}

		// Navigation and search handling
		return m.handleInput(msg)

	case tea.WindowSizeMsg:
		return m, nil
	}

	return m, cmd
}

// handleInput handles all input in peco-like mode
// Always in "search mode" with navigation via arrow keys
func (m Model) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := msg.String()

	// Handle navigation keys
	switch keyStr {
	case "up":
		if m.cursor > 0 {
			m.cursor--
		} else if len(m.filteredProjects) > 0 {
			m.cursor = len(m.filteredProjects) - 1
		}
		return m, nil

	case "down":
		if m.cursor < len(m.filteredProjects)-1 {
			m.cursor++
		} else if len(m.filteredProjects) > 0 {
			m.cursor = 0
		}
		return m, nil
	}

	// All other input goes to the search box
	oldValue := m.textInput.Value()

	// Update the text input
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	newValue := m.textInput.Value()

	// If the value changed, apply the filter
	if newValue != oldValue {
		m.applyFilter()
	}

	return m, cmd
}

// applyFilter filters the projects based on the current search query
func (m *Model) applyFilter() {
	query := strings.ToLower(strings.TrimSpace(m.textInput.Value()))

	// If query is empty, show all projects
	if query == "" {
		m.filteredProjects = m.projects
		m.cursor = 0
		return
	}

	// Filter projects
	var filtered []status.ProjectDisplay
	for _, p := range m.projects {
		// Check multiple fields for match
		if m.matchesQuery(p, query) {
			filtered = append(filtered, p)
		}
	}

	m.filteredProjects = filtered

	// Reset cursor to first item
	m.cursor = 0
}

// matchesQuery checks if a project matches the search query
func (m *Model) matchesQuery(p status.ProjectDisplay, query string) bool {
	repoLower := strings.ToLower(p.Repo)
	pathLower := strings.ToLower(p.FullPath)
	workspaceLower := strings.ToLower(p.Workspace)

	// When the user types a simple token (no slash), prefer matching
	// against the repository basename (the part after the last '/').
	// e.g., typing "ghqx" should match "mi8bi/ghqx" by its repo name.
	if !strings.Contains(query, "/") {
		// Match either the repo basename (after '/') or the owner (before '/')
		repoName := repoLower
		ownerName := repoLower
		if idx := strings.LastIndex(repoLower, "/"); idx != -1 {
			repoName = repoLower[idx+1:]
			ownerName = repoLower[:idx]
		}

		return strings.Contains(repoName, query) || strings.Contains(ownerName, query)
	}

	// For queries with a slash or more specific tokens, fall back to
	// broader matching across repo, full path, and workspace.
	return strings.Contains(repoLower, query) ||
		strings.Contains(pathLower, query) ||
		strings.Contains(workspaceLower, query)
}

// View implements tea.Model.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("36")).
		Render(i18n.T("selector.title"))
	s.WriteString(title + "\n")

	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(strings.Repeat("─", 60))
	s.WriteString(separator + "\n\n")

	// Search input
	m.renderSearchInput(&s)

	// Project list
	m.renderProjectList(&s)

	// Help text
	m.renderHelp(&s)

	return s.String()
}

// renderSearchInput renders the search input field
func (m *Model) renderSearchInput(s *strings.Builder) {
	filterLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("246")).
		Render(i18n.T("selector.search.label") + " ")
	s.WriteString(filterLabel)

	s.WriteString(m.textInput.View())

	// Match count
	if len(m.projects) > 0 {
		matchInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")).
			Render(fmt.Sprintf(" [%d/%d]", len(m.filteredProjects), len(m.projects)))
		s.WriteString(matchInfo)
	}
	s.WriteString("\n\n")
}

// renderProjectList renders the filtered project list
func (m *Model) renderProjectList(s *strings.Builder) {
	// No matches message
	if len(m.filteredProjects) == 0 {
		if m.textInput.Value() != "" {
			noMatches := lipgloss.NewStyle().
				Foreground(lipgloss.Color("208")).
				Render(i18n.T("selector.search.noMatches"))
			s.WriteString("  " + noMatches + "\n")
		}
		return
	}

	// Project items
	for i, p := range m.filteredProjects {
		cursor := "  "
		if m.cursor == i {
			cursor = "❯ "
		}

		line := fmt.Sprintf("%s%-40s  %s", cursor, p.Repo,
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(p.Workspace))

		if m.cursor == i {
			// Highlight selected item
			highlighted := lipgloss.NewStyle().
				Background(lipgloss.Color("24")).
				Foreground(lipgloss.Color("255")).
				Render(line)
			s.WriteString(highlighted)
		} else {
			s.WriteString(line)
		}
		s.WriteString("\n")
	}
}

// renderHelp renders the help text
func (m *Model) renderHelp(s *strings.Builder) {
	s.WriteString("\n")

	helpText := i18n.T("selector.helpWithPecoSearch")

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(helpText)
	s.WriteString(help)
}

// Run displays the interactive selector and returns the full path of the selected project.
func Run(projects []status.ProjectDisplay) (string, error) {
	if len(projects) == 0 {
		return "", nil // No projects to select
	}

	model := NewModel(projects)
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	if m, ok := finalModel.(Model); ok && !m.quitting {
		return m.selected, nil
	}
	return "", nil // User quit or nothing selected
}
