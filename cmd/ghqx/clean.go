package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/i18n"
	"github.com/mi8bi/ghqx/internal/ui"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "", // Will be set in root.go init() after locale is determined
	Long:  "", // Will be set in root.go init() after locale is determined
	RunE:  runClean,
}

func runClean(cmd *cobra.Command, args []string) error {
	fmt.Println(ui.FormatWarning(i18n.T("clean.warning.title")))
	fmt.Println(i18n.T("clean.warning.description"))

	// Load the app to get config, but handle errors gracefully
	// as the config file might not even exist.
	loadedApp, err := app.NewFromConfigPath(configPath)
	if err == nil {
		fmt.Println("\n" + i18n.T("clean.warning.targetRoots"))
		for name, path := range loadedApp.Config.Roots {
			fmt.Printf("- %s (%s)\n", name, path)
		}
	} else {
		fmt.Println(i18n.T("clean.warning.noConfigFound"))
	}

	fmt.Printf("\n%s ", i18n.T("clean.warning.confirm"))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	if input != "yes" {
		fmt.Println(i18n.T("clean.aborted"))
		return nil
	}

	fmt.Println()

	// --- Deletion Phase ---

	// 1. Delete root directories
	if loadedApp != nil {
		fmt.Println(i18n.T("clean.deleting.roots"))
		for name, path := range loadedApp.Config.Roots {
			fmt.Printf("  - %s (%s)... ", name, path)
			if err := os.RemoveAll(path); err != nil {
				fmt.Println(ui.FormatError(err))
			} else {
				fmt.Println(i18n.T("clean.deleting.success"))
			}
		}
	}

	// 2. Delete config file
	fmt.Println(i18n.T("clean.deleting.config"))

	// Try to find the config path again to be sure.
	// loader := config.NewLoader() // This line is removed.
	// We'll rely on the global `configPath` flag or the default path.
	cfgPathToDelete := configPath
	if cfgPathToDelete == "" {
		// Attempt to find default path
		// We suppress the error because the file might not exist, which is fine.
		defaultPath, _ := config.GetDefaultConfigPath()
		cfgPathToDelete = defaultPath
	}

	// Double check we have a path to delete.
	if cfgPathToDelete != "" {
		if _, err := os.Stat(cfgPathToDelete); err == nil {
			fmt.Printf("  - %s... ", cfgPathToDelete)
			if err := os.Remove(cfgPathToDelete); err != nil {
				fmt.Println(ui.FormatError(err))
			} else {
				fmt.Println(i18n.T("clean.deleting.success"))
			}
		} else {
			fmt.Println(i18n.T("clean.deleting.noConfigFound"))
		}
	} else {
		fmt.Println(i18n.T("clean.deleting.noConfigPath"))
	}

	fmt.Println(ui.FormatSuccess(i18n.T("clean.complete")))
	return nil
}
