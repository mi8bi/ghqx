package main

import (
	"fmt"
	"os"

	"github.com/mi8bi/ghqx/internal/doctor"
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
	doctorService := doctor.NewService()
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
		// NGがあった場合は終了コード 1 で終了
		os.Exit(1)
	}

	return nil
}
