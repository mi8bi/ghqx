package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

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

func outputCompactTable(projects []status.ProjectDisplay) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "%-30s  %-10s  %-10s  %-8s\n",
		i18n.T("status.header.name"),
		i18n.T("status.header.zone"),
		i18n.T("status.header.gitManaged"),
		i18n.T("status.header.status"),
	)
	fmt.Fprintf(w, "%-30s  %-10s  %-10s  %-8s\n",
		strings.Repeat("-", 30),
		strings.Repeat("-", 10),
		strings.Repeat("-", 10),
		strings.Repeat("-", 8),
	)

	for _, proj := range projects {
		repo := truncateString(proj.Repo, 30)
		fmt.Fprintf(w, "%-30s  %-10s  %-10s  %-8s\n",
			repo, proj.Zone, proj.GitManaged, proj.Status)
	}

	return nil
}

func outputVerboseTable(projects []status.ProjectDisplay) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "%-30s  %-10s  %-10s  %-50s  %-10s  %-8s\n",
		i18n.T("status.header.name"),
		i18n.T("status.header.zone"),
		i18n.T("status.header.root"),
		i18n.T("status.header.path"),
		i18n.T("status.header.gitManaged"),
		i18n.T("status.header.status"),
	)
	fmt.Fprintf(w, "%-30s  %-10s  %-10s  %-50s  %-10s  %-8s\n",
		strings.Repeat("-", 30),
		strings.Repeat("-", 10),
		strings.Repeat("-", 10),
		strings.Repeat("-", 50),
		strings.Repeat("-", 10),
		strings.Repeat("-", 8),
	)

	for _, proj := range projects {
		repo := truncateString(proj.Repo, 30)
		fullPath := truncateString(proj.FullPath, 50)
		fmt.Fprintf(w, "%-30s  %-10s  %-10s  %-50s  %-10s  %-8s\n",
			repo, proj.Zone, proj.RawProject.Root, fullPath, proj.GitManaged, proj.Status)
	}

	return nil
}

// truncateString truncates a string to a given length, adding "..." if truncated.
func truncateString(s string, length int) string {
	if len(s) > length && length > 3 {
		return s[:length-3] + "..."
	}
	if len(s) > length {
		return s[:length]
	}
	return s
}
