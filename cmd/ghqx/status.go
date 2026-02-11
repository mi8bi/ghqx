package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

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
	Short: i18n.T("status.command.short"),
	Long:  i18n.T("status.command.long"),
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

// getDisplayWidth is a helper to get the display width of an i18n string
func getDisplayWidth(key string) int {
	return runewidth.StringWidth(i18n.T(key))
}

func outputCompactTable(projects []status.ProjectDisplay) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Use i18n keys for headers
	headerName := i18n.T("status.header.name")
	headerWorkspace := i18n.T("status.header.workspace")
	headerGitManaged := i18n.T("status.header.gitManaged")
	headerStatus := i18n.T("status.header.status")

	// Calculate initial minimum widths based on header text (display width)
	minNameWidth := runewidth.StringWidth(headerName)
	minWorkspaceWidth := runewidth.StringWidth(headerWorkspace)
	minGitManagedWidth := runewidth.StringWidth(headerGitManaged)
	minStatusWidth := runewidth.StringWidth(headerStatus)

	// Additionally, calculate max width of content in each column
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
	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s\n",
		minNameWidth, headerName,
		minWorkspaceWidth, headerWorkspace,
		minGitManagedWidth, headerGitManaged,
		minStatusWidth, headerStatus,
	)
	// Print separator line
	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s\n",
		minNameWidth, strings.Repeat("-", minNameWidth),
		minWorkspaceWidth, strings.Repeat("-", minWorkspaceWidth),
		minGitManagedWidth, strings.Repeat("-", minGitManagedWidth),
		minStatusWidth, strings.Repeat("-", minStatusWidth),
	)

	// Print project data
	for _, proj := range projects {
		// Use FillRight to ensure correct padding for multi-byte characters
		fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s\n",
			minNameWidth, runewidth.FillRight(proj.Repo, minNameWidth),
			minWorkspaceWidth, runewidth.FillRight(proj.Workspace, minWorkspaceWidth),
			minGitManagedWidth, runewidth.FillRight(proj.GitManaged, minGitManagedWidth),
			minStatusWidth, runewidth.FillRight(proj.Status, minStatusWidth),
		)
	}

	return nil
}

func outputVerboseTable(projects []status.ProjectDisplay) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Use i18n keys for headers
	headerName := i18n.T("status.header.name")
	headerWorkspace := i18n.T("status.header.workspace")
	headerRoot := i18n.T("status.header.root")
	headerPath := i18n.T("status.header.path")
	headerGitManaged := i18n.T("status.header.gitManaged")
	headerStatus := i18n.T("status.header.status")

	// Calculate initial minimum widths based on header text (display width)
	minNameWidth := runewidth.StringWidth(headerName)
	minWorkspaceWidth := runewidth.StringWidth(headerWorkspace)
	minRootWidth := runewidth.StringWidth(headerRoot)
	minPathWidth := runewidth.StringWidth(headerPath)
	minGitManagedWidth := runewidth.StringWidth(headerGitManaged)
	minStatusWidth := runewidth.StringWidth(headerStatus)

	// Additionally, calculate max width of content in each column
	for _, proj := range projects {
		nameLen := runewidth.StringWidth(proj.Repo)
		if nameLen > minNameWidth {
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
		if pathLen > minPathWidth {
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
	
	// Limit Path width to a reasonable maximum to avoid overly wide lines
	if minPathWidth > 80 {
		minPathWidth = 80
	}
	// Limit Repo name width to a reasonable maximum
	if minNameWidth > 50 {
		minNameWidth = 50
	}


	// Print header
	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s  %-*s  %-*s\n",
		minNameWidth, headerName,
		minWorkspaceWidth, headerWorkspace,
		minRootWidth, headerRoot,
		minPathWidth, headerPath,
		minGitManagedWidth, headerGitManaged,
		minStatusWidth, headerStatus,
	)
	// Print separator line
	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s  %-*s  %-*s\n",
		minNameWidth, strings.Repeat("-", minNameWidth),
		minWorkspaceWidth, strings.Repeat("-", minWorkspaceWidth),
		minRootWidth, strings.Repeat("-", minRootWidth),
		minPathWidth, strings.Repeat("-", minPathWidth),
		minGitManagedWidth, strings.Repeat("-", minGitManagedWidth),
		minStatusWidth, strings.Repeat("-", minStatusWidth),
	)

	// Print project data
	for _, proj := range projects {
		// Use FillRight for correct padding with multi-byte characters
		// Truncate long paths
		displayPath := runewidth.FillRight(proj.FullPath, minPathWidth)
		if runewidth.StringWidth(displayPath) > minPathWidth {
			displayPath = runewidth.Truncate(proj.FullPath, minPathWidth, "...")
		}
		
		displayName := runewidth.FillRight(proj.Repo, minNameWidth)
		if runewidth.StringWidth(displayName) > minNameWidth {
			displayName = runewidth.Truncate(proj.Repo, minNameWidth, "...")
		}


		fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s  %-*s  %-*s\n",
			minNameWidth, displayName,
			minWorkspaceWidth, runewidth.FillRight(proj.Workspace, minWorkspaceWidth),
			minRootWidth, runewidth.FillRight(string(proj.RawProject.Root), minRootWidth),
			minPathWidth, displayPath,
			minGitManagedWidth, runewidth.FillRight(proj.GitManaged, minGitManagedWidth),
			minStatusWidth, runewidth.FillRight(proj.Status, minStatusWidth),
		)
	}

	return nil
}

// truncateString truncates a string to a given length, adding "..." if truncated.
// This function is no longer used directly as runewidth.Truncate is now used.
// Keeping it for now in case of other uses or a fallback.
func truncateString(s string, length int) string {
	if len(s) > length && length > 3 {
		return s[:length-3] + "..."
	}
	if len(s) > length {
		return s[:length]
	}
	return s
}
