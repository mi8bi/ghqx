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
	Short: i18n.T("doctor.command.short"),
	Long:  i18n.T("doctor.command.long"),
	RunE: runDoctor,
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
