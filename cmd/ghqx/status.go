package main

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/mi8bi/ghqx/internal/i18n"
	"github.com/mi8bi/ghqx/internal/status"
	"github.com/mi8bi/ghqx/internal/tui"
	"github.com/spf13/cobra"
)

var (
	statusVerbose bool
	statusTUI     bool
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "", // Will be set in root.go init() after locale is determined
	Long:  "", // Will be set in root.go init() after locale is determined
	RunE:  runStatus,
}

func init() {
	statusCmd.Flags().BoolVarP(&statusVerbose, "verbose", "v", false, i18n.T("status.flag.verbose"))
	statusCmd.Flags().BoolVar(&statusTUI, "tui", false, i18n.T("status.flag.tui"))
}

func runStatus(cmd *cobra.Command, args []string) error {
	if err := loadApp(); err != nil {
		return err
	}

	// TUI mode
	if statusTUI {
		return tui.RunStatus(application)
	}

	// CLI mode
	opts := status.Options{
		CheckDirty: true,
		LoadBranch: false,
	}

	rawProjects, err := application.Status.GetAll(opts)
	if err != nil {
		return err
	}

	// Convert domain.Project to ProjectDisplay
	displayProjects := make([]status.ProjectDisplay, len(rawProjects))
	for i, p := range rawProjects {
		displayProjects[i] = status.NewProjectDisplay(p)
	}

	if statusVerbose {
		return outputVerboseTable(displayProjects)
	}

	return outputCompactTable(displayProjects)
}

func outputCompactTable(projects []status.ProjectDisplay) error {
	// Use i18n keys for headers
	headerName := i18n.T("status.header.name")
	headerWorkspace := i18n.T("status.header.workspace")
	headerGitManaged := i18n.T("status.header.gitManaged")
	headerStatus := i18n.T("status.header.status")

	// Calculate minimum column widths based on headers
	minNameWidth := runewidth.StringWidth(headerName)
	minWorkspaceWidth := runewidth.StringWidth(headerWorkspace)
	minGitManagedWidth := runewidth.StringWidth(headerGitManaged)
	minStatusWidth := runewidth.StringWidth(headerStatus)

	// Calculate max width of content in each column
	for _, proj := range projects {
		nameLen := runewidth.StringWidth(proj.Repo)
		if nameLen > minNameWidth {
			minNameWidth = nameLen
		}

		workspaceLen := runewidth.StringWidth(proj.Workspace)
		if workspaceLen > minWorkspaceWidth {
			minWorkspaceWidth = workspaceLen
		}

		gitManagedLen := runewidth.StringWidth(proj.GitManaged)
		if gitManagedLen > minGitManagedWidth {
			minGitManagedWidth = gitManagedLen
		}

		statusLen := runewidth.StringWidth(proj.Status)
		if statusLen > minStatusWidth {
			minStatusWidth = statusLen
		}
	}

	// Print header
	fmt.Printf("%s  %s  %s  %s\n",
		padRight(headerName, minNameWidth),
		padRight(headerWorkspace, minWorkspaceWidth),
		padRight(headerGitManaged, minGitManagedWidth),
		padRight(headerStatus, minStatusWidth),
	)

	// Print separator line
	fmt.Printf("%s  %s  %s  %s\n",
		strings.Repeat("-", minNameWidth),
		strings.Repeat("-", minWorkspaceWidth),
		strings.Repeat("-", minGitManagedWidth),
		strings.Repeat("-", minStatusWidth),
	)

	// Print project data
	for _, proj := range projects {
		fmt.Printf("%s  %s  %s  %s\n",
			padRight(proj.Repo, minNameWidth),
			padRight(proj.Workspace, minWorkspaceWidth),
			padRight(proj.GitManaged, minGitManagedWidth),
			padRight(proj.Status, minStatusWidth),
		)
	}

	return nil
}

func outputVerboseTable(projects []status.ProjectDisplay) error {
	// Use i18n keys for headers
	headerName := i18n.T("status.header.name")
	headerWorkspace := i18n.T("status.header.workspace")
	headerRoot := i18n.T("status.header.root")
	headerPath := i18n.T("status.header.path")
	headerGitManaged := i18n.T("status.header.gitManaged")
	headerStatus := i18n.T("status.header.status")

	// Calculate minimum widths
	minNameWidth := max(runewidth.StringWidth(headerName), 20)
	minWorkspaceWidth := max(runewidth.StringWidth(headerWorkspace), 10)
	minRootWidth := max(runewidth.StringWidth(headerRoot), 8)
	minPathWidth := max(runewidth.StringWidth(headerPath), 30)
	minGitManagedWidth := max(runewidth.StringWidth(headerGitManaged), 10)
	minStatusWidth := max(runewidth.StringWidth(headerStatus), 8)

	// Calculate max content width
	for _, proj := range projects {
		nameLen := runewidth.StringWidth(proj.Repo)
		if nameLen > minNameWidth && nameLen < 50 {
			minNameWidth = nameLen
		}

		workspaceLen := runewidth.StringWidth(proj.Workspace)
		if workspaceLen > minWorkspaceWidth {
			minWorkspaceWidth = workspaceLen
		}

		rootLen := runewidth.StringWidth(string(proj.RawProject.Root))
		if rootLen > minRootWidth {
			minRootWidth = rootLen
		}

		pathLen := runewidth.StringWidth(proj.FullPath)
		if pathLen > minPathWidth && pathLen < 80 {
			minPathWidth = pathLen
		}

		gitManagedLen := runewidth.StringWidth(proj.GitManaged)
		if gitManagedLen > minGitManagedWidth {
			minGitManagedWidth = gitManagedLen
		}

		statusLen := runewidth.StringWidth(proj.Status)
		if statusLen > minStatusWidth {
			minStatusWidth = statusLen
		}
	}

	// Print header
	fmt.Printf("%s  %s  %s  %s  %s  %s\n",
		padRight(headerName, minNameWidth),
		padRight(headerWorkspace, minWorkspaceWidth),
		padRight(headerRoot, minRootWidth),
		padRight(headerPath, minPathWidth),
		padRight(headerGitManaged, minGitManagedWidth),
		padRight(headerStatus, minStatusWidth),
	)

	// Print separator
	fmt.Printf("%s  %s  %s  %s  %s  %s\n",
		strings.Repeat("-", minNameWidth),
		strings.Repeat("-", minWorkspaceWidth),
		strings.Repeat("-", minRootWidth),
		strings.Repeat("-", minPathWidth),
		strings.Repeat("-", minGitManagedWidth),
		strings.Repeat("-", minStatusWidth),
	)

	// Print data
	for _, proj := range projects {
		displayName := truncateString(proj.Repo, minNameWidth)
		displayPath := truncateString(proj.FullPath, minPathWidth)

		fmt.Printf("%s  %s  %s  %s  %s  %s\n",
			padRight(displayName, minNameWidth),
			padRight(proj.Workspace, minWorkspaceWidth),
			padRight(string(proj.RawProject.Root), minRootWidth),
			padRight(displayPath, minPathWidth),
			padRight(proj.GitManaged, minGitManagedWidth),
			padRight(proj.Status, minStatusWidth),
		)
	}

	return nil
}

// padRight pads a string to the right with spaces to match the specified width
func padRight(s string, width int) string {
	return runewidth.FillRight(s, width)
}

// truncateString truncates a string to a given length, adding "..." if truncated
func truncateString(s string, length int) string {
	if runewidth.StringWidth(s) <= length {
		return s
	}
	if length <= 3 {
		return s[:length]
	}
	return runewidth.Truncate(s, length, "...")
}
