package main

import (
	"fmt"
	"os"

	"github.com/mi8bi/ghqx/internal/ui"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, ui.FormatDetailedError(err))
		os.Exit(1)
	}
}

