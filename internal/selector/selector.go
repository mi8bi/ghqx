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

// UI Configuration constants
const (
	// Input field configuration
	searchInputMaxChars = 100
	searchInputWidth    = 50

	// Display configuration
	separatorLength = 60

	// Color constants
	colorTitle     = "36"  // Cyan
	colorSeparator = "240" // Dark gray
	colorLabel     = "246" // Medium gray
	colorMessage   = "241" // Light gray
	colorWarning   = "208" // Orange
	colorCursor    = "24"  // Dark blue
	colorSelected  = "255" // White
)

// Model is the Bubble Tea model for the interactive project selector.
// It manages the state for searching, filtering, and selecting projects.
type Model struct {
	// projects holds all available projects (unfiltered)
	projects []status.ProjectDisplay

	// filteredProjects holds the current search results
	filteredProjects []status.ProjectDisplay

	// textInput is the search box for filtering projects
	textInput textinput.Model

	// cursor tracks the currently highlighted project index
	cursor int

	// selected stores the FullPath of the selected project
	selected string

	// quitting indicates whether the user has exited the selector
	quitting bool
}

// NewModel creates a new selector model with the given projects.
// Input field is pre-configured with localized placeholder text and focused.
func NewModel(projects []status.ProjectDisplay) Model {
	ti := textinput.New()
	ti.Placeholder = i18n.T("selector.search.placeholder")
	ti.CharLimit = searchInputMaxChars
	ti.Width = searchInputWidth
	ti.Prompt = "" // Hide default "> " prompt for cleaner display
	ti.Focus()     // Start with focus on search box for immediate typing

	return Model{
		projects:         projects,
		filteredProjects: projects, // Initially show all projects
		textInput:        ti,
		cursor:           0,
		quitting:         false,
	}
}

// Init implements tea.Model interface. Returns command for blinking text cursor.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model interface. Handles all user input and state changes.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyInput(msg)
	case tea.WindowSizeMsg:
		// Ignore window size changes; use fixed dimensions
		return m, nil
	}

	return m, cmd
}

// handleKeyInput processes keyboard input and updates the model state accordingly.
func (m Model) handleKeyInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := msg.String()

	// Global quit commands
	if keyStr == "ctrl+c" {
		m.quitting = true
		return m, tea.Quit
	}

	if keyStr == "esc" {
		m.quitting = true
		return m, tea.Quit
	}

	// Selection (Enter key)
	if keyStr == "enter" {
		if m.cursor >= 0 && m.cursor < len(m.filteredProjects) {
			m.selected = m.filteredProjects[m.cursor].FullPath
		}
		return m, tea.Quit
	}

	// Delegate to input handler for navigation and typing
	return m.handleInput(msg)
}

// handleInput processes navigation (arrow keys) and search input (text).
// Navigation is limited to arrow keys (up/down) for peco-like behavior.
// All other input is forwarded to the text input field.
func (m Model) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := msg.String()

	// Handle vertical navigation with arrow keys
	switch keyStr {
	case "up":
		m.moveCursorUp()
		return m, nil

	case "down":
		m.moveCursorDown()
		return m, nil
	}

	// All other input goes to the search box
	return m.updateSearchInput(msg)
}

// moveCursorUp moves the cursor to the previous project, wrapping around to the end if at the top.
func (m *Model) moveCursorUp() {
	if m.cursor > 0 {
		m.cursor--
	} else if len(m.filteredProjects) > 0 {
		// Wrap around to the end
		m.cursor = len(m.filteredProjects) - 1
	}
}

// moveCursorDown moves the cursor to the next project, wrapping around to the start if at the end.
func (m *Model) moveCursorDown() {
	if m.cursor < len(m.filteredProjects)-1 {
		m.cursor++
	} else if len(m.filteredProjects) > 0 {
		// Wrap around to the start
		m.cursor = 0
	}
}

// updateSearchInput updates the text input field and applies filtering if the text changed.
func (m Model) updateSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	oldValue := m.textInput.Value()

	// Process the key input through the text input component
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	newValue := m.textInput.Value()

	// Re-apply filter only if the search text changed
	if newValue != oldValue {
		m.applyFilter()
	}

	return m, cmd
}

// applyFilter filters projects based on the current search query.
// Empty query shows all projects. Results are sorted with cursor at the first match.
func (m *Model) applyFilter() {
	query := strings.ToLower(strings.TrimSpace(m.textInput.Value()))

	// Empty query: show all projects
	if query == "" {
		m.filteredProjects = m.projects
		m.cursor = 0
		return
	}

	// Filter projects based on query
	var filtered []status.ProjectDisplay
	for _, p := range m.projects {
		if m.matchesQuery(p, query) {
			filtered = append(filtered, p)
		}
	}

	m.filteredProjects = filtered
	m.cursor = 0 // Reset to first result
}

// matchesQuery checks if a project matches the search query.
// For simple tokens (no slash), it prioritizes matching the repository name.
// For complex queries (with slash), it matches across repo, path, and workspace.
func (m *Model) matchesQuery(p status.ProjectDisplay, query string) bool {
	repoLower := strings.ToLower(p.Repo)
	pathLower := strings.ToLower(p.FullPath)
	workspaceLower := strings.ToLower(p.Workspace)

	// Simple query (no slash): match against repo name and owner
	// e.g., "ghqx" matches "mi8bi/ghqx" by repository name
	if !strings.Contains(query, "/") {
		return m.matchesSimpleQuery(repoLower, query)
	}

	// Complex query (with slash): match across all fields
	// e.g., "mi8bi/gh" matches against full path, repo, and workspace
	return strings.Contains(repoLower, query) ||
		strings.Contains(pathLower, query) ||
		strings.Contains(workspaceLower, query)
}

// matchesSimpleQuery matches a simple query against the repository basename and owner.
// Returns true if either the repo name (after '/') or owner (before '/') contains the query.
func (m *Model) matchesSimpleQuery(repoLower, query string) bool {
	// Extract repo name (after last '/') and owner (before last '/')
	repoName := repoLower
	ownerName := repoLower

	if idx := strings.LastIndex(repoLower, "/"); idx != -1 {
		repoName = repoLower[idx+1:]
		ownerName = repoLower[:idx]
	}

	return strings.Contains(repoName, query) || strings.Contains(ownerName, query)
}

// View implements tea.Model interface. Renders the UI to the terminal.
func (m Model) View() string {
	if m.quitting {
		return "" // Return empty when quitting to avoid final frame
	}

	var s strings.Builder

	// Render each section of the UI
	m.renderHeader(&s)
	m.renderSearchInput(&s)
	m.renderProjectList(&s)
	m.renderFooter(&s)

	return s.String()
}

// renderHeader renders the title and separator line.
func (m *Model) renderHeader(s *strings.Builder) {
	// Title with emphasis
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorTitle)).
		Render(i18n.T("selector.title"))
	s.WriteString(title + "\n")

	// Separator for visual clarity
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorSeparator)).
		Render(strings.Repeat("─", separatorLength))
	s.WriteString(separator + "\n\n")
}

// renderSearchInput renders the search input field with match count.
func (m *Model) renderSearchInput(s *strings.Builder) {
	// Search label
	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorLabel)).
		Render(i18n.T("selector.search.label") + " ")
	s.WriteString(label)

	// Input box
	s.WriteString(m.textInput.View())

	// Match count indicator
	if len(m.projects) > 0 {
		count := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorLabel)).
			Render(fmt.Sprintf(" [%d/%d]", len(m.filteredProjects), len(m.projects)))
		s.WriteString(count)
	}

	s.WriteString("\n\n")
}

// renderProjectList renders the filtered list of projects with cursor highlighting.
func (m *Model) renderProjectList(s *strings.Builder) {
	// Handle no matches case
	if len(m.filteredProjects) == 0 {
		if m.textInput.Value() != "" {
			noMatches := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorWarning)).
				Render(i18n.T("selector.search.noMatches"))
			s.WriteString("  " + noMatches + "\n")
		}
		return
	}

	// Render each project item
	for i, p := range m.filteredProjects {
		m.renderProjectItem(s, i, p)
	}
}

// renderProjectItem renders a single project line with optional highlighting.
func (m *Model) renderProjectItem(s *strings.Builder, index int, project status.ProjectDisplay) {
	// Cursor indicator
	cursorChar := "  "
	if m.cursor == index {
		cursorChar = "❯ " // Pointing indicator for current selection
	}

	// Format: [cursor] [repo name]  [workspace]
	workspaceStr := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorSeparator)).
		Render(project.Workspace)

	line := fmt.Sprintf("%s%-40s  %s", cursorChar, project.Repo, workspaceStr)

	// Highlight selected item
	if m.cursor == index {
		line = lipgloss.NewStyle().
			Background(lipgloss.Color(colorCursor)).
			Foreground(lipgloss.Color(colorSelected)).
			Render(line)
	}

	s.WriteString(line + "\n")
}

// renderFooter renders the help text with keybinding instructions.
func (m *Model) renderFooter(s *strings.Builder) {
	s.WriteString("\n")

	helpText := i18n.T("selector.helpWithPecoSearch")
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMessage)).
		Render(helpText)

	s.WriteString(help)
}

// Run displays the interactive selector in an alternate screen buffer.
// Returns the full path of the selected project, or empty string if canceled.
func Run(projects []status.ProjectDisplay) (string, error) {
	// Early exit if no projects available
	if len(projects) == 0 {
		return "", nil
	}

	// Initialize and run the TUI
	model := NewModel(projects)
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	// Extract the selected project path from the final model
	if m, ok := finalModel.(Model); ok && !m.quitting {
		return m.selected, nil
	}

	// User canceled or no selection made
	return "", nil
}
