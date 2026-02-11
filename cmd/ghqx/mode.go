package main

import (
	"fmt"
	"sort"

	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/i18n"
	"github.com/spf13/cobra"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var modeCmd = &cobra.Command{
	Use:   "mode",
	Short: i18n.T("mode.command.short"),
	Long:  i18n.T("mode.command.long"),
	RunE: runMode,
}

func init() {
}

// ModeSelectorModel is the Bubble Tea model for the mode selector.
type ModeSelectorModel struct {
	workspaceNames []string // Renamed from choices
	cursor  int      // which choice is selected
	selected string   // the selected choice
	quitting bool
}

// Init implements tea.Model.
func (m ModeSelectorModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m ModeSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.workspaceNames)-1 { // Updated to workspaceNames
				m.cursor++
			}
		case "enter":
			if m.cursor >= 0 && m.cursor < len(m.workspaceNames) { // Updated to workspaceNames
				m.selected = m.workspaceNames[m.cursor] // Updated to workspaceNames
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m ModeSelectorModel) View() string {
	if m.quitting {
		return ""
	}

	s := lipgloss.NewStyle().Bold(true).Render(i18n.T("mode.selector.title")) + "\n\n"

	for i, choice := range m.workspaceNames { // Updated to workspaceNames
		cursor := "  " // no cursor
		if m.cursor == i {
			cursor = "> " // cursor!
		}

		style := lipgloss.NewStyle()
		if m.cursor == i {
			style = style.Background(lipgloss.Color("240")).Foreground(lipgloss.Color("230"))
		}

		s += style.Render(fmt.Sprintf("%s%s", cursor, choice)) + "\n"
	}
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1).Render(i18n.T("mode.selector.help"))
	return s
}

func runMode(cmd *cobra.Command, args []string) error {
	if err := loadApp(); err != nil {
		return err
	}

	// Get available root names from config
	var rootNames []string
	for name := range application.Config.Roots {
		rootNames = append(rootNames, name)
	}
	sort.Strings(rootNames) // Sort for consistent display

	if len(rootNames) == 0 {
		return fmt.Errorf(i18n.T("mode.error.noRoots"))
	}

	// Initialize the TUI model
	model := ModeSelectorModel{
		workspaceNames: rootNames, // Updated to workspaceNames
		cursor:  0,
	}

	// If a default root is already set, try to pre-select it
	currentDefaultRoot := application.Config.GetDefaultRoot()
	for i, name := range rootNames {
		if name == currentDefaultRoot {
			model.cursor = i
			break
		}
	}

	// Run the TUI
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	// Check selected value
	if m, ok := finalModel.(ModeSelectorModel); ok && !m.quitting && m.selected != "" {
		if m.selected == currentDefaultRoot {
			fmt.Println(i18n.T("mode.noChange"))
			return nil
		}

		// Update config
		application.Config.Default.Root = m.selected
		loader := config.NewLoader()
		
		// Determine config path for saving
		savePath := configPath
		if savePath == "" {
			savePath, err = config.GetDefaultConfigPath()
			if err != nil {
				return err
			}
		}

		if err := loader.Save(application.Config, savePath); err != nil {
			return err
		}
		// Updated message to use workspace terminology
		fmt.Print(i18n.T("mode.success") + m.selected)
	} else {
		fmt.Println(i18n.T("mode.aborted"))
	}

	return nil
}
