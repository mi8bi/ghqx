package main

import (
	"fmt"

	"github.com/mi8bi/ghqx/internal/doctor"
	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/i18n"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "", // Will be set in root.go init() after locale is determined
	Long:  "", // Will be set in root.go init() after locale is determined
	RunE:  runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) error {
	doctorService := doctor.NewServiceWithConfigPath(configPath)
	results := doctorService.RunChecks()

	allOK := true
	for _, res := range results {
		if res.OK {
			fmt.Printf("%s %s\n", i18n.T("doctor.result.ok"), res.Message)
		} else {
			allOK = false
			fmt.Printf("%s %s\n", i18n.T("doctor.result.ng"), res.Message)
			if res.Hint != "" {
				fmt.Printf("     %s: %s\n", i18n.T("doctor.result.hint"), res.Hint)
			}
		}
	}

	if !allOK {
		// Return an error instead of os.Exit(1)
		return domain.NewError(
			domain.ErrCodeUnknown,
			"Environment check failed",
		).WithHint("Fix the issues above and try again")
	}

	return nil
}
