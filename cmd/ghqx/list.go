package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/status"
	"github.com/spf13/cobra"
)

var (
	listFullPath bool
	listAllRoots bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "", // Will be set in root.go init() after locale is determined
	Long:  "", // Will be set in root.go init() after locale is determined
	RunE:  runList,
}

func init() {
	listCmd.Flags().BoolVarP(&listFullPath, "full-path", "p", false, "show full absolute paths")
	listCmd.Flags().BoolVarP(&listAllRoots, "all", "a", false, "list projects from all roots (default: only default root)")
}

func runList(cmd *cobra.Command, args []string) error {
	// Skip loading if application is already initialized (for testing)
	if application == nil {
		if err := loadApp(); err != nil {
			return err
		}
	}

	opts := status.Options{
		CheckDirty: false,
		LoadBranch: false,
	}

	// Get all projects
	allProjects, err := application.Status.GetAll(opts)
	if err != nil {
		return err
	}

	// Filter projects based on flags
	var projects []domain.Project
	if listAllRoots {
		projects = allProjects
	} else {
		// Filter to only default root
		defaultRootName := application.Config.GetDefaultRoot()
		for _, p := range allProjects {
			if string(p.Root) == defaultRootName {
				projects = append(projects, p)
			}
		}
	}

	// Output project paths, one per line
	for _, p := range projects {
		if listFullPath {
			// Absolute path
			fmt.Println(p.Path)
		} else {
			// Relative path from root (like ghq list format)
			rootPath, exists := application.Config.GetRoot(string(p.Root))
			if exists {
				relPath, err := filepath.Rel(rootPath, p.Path)
				if err == nil && !strings.HasPrefix(relPath, "..") {
					fmt.Println(relPath)
					continue
				}
			}
			// Fallback to project name
			fmt.Println(p.Name)
		}
	}

	return nil
}
