package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Build-time variables set via -ldflags
// These are populated during the build process with version info.
// Example: go build -ldflags "-X main.version=v1.0.0 -X main.commit=abc123 -X main.buildTime=2026-02-11T04:12:00Z"
var (
	version   = "dev"
	commit    = "none"
	buildTime = "unknown"
)

var (
	versionVerbose bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "", // Will be set in root.go init() after locale is determined
	Long:  "", // Will be set in root.go init() after locale is determined
	Run:   runVersion,
}

func init() {
	versionCmd.Flags().BoolVarP(&versionVerbose, "verbose", "v", false, "Show detailed version information")
}

// runVersion outputs the application version.
// Default output: ghqx vX.Y.Z (single line, human-friendly)
// With --verbose: extended information including commit hash, build time, and Go version
func runVersion(cmd *cobra.Command, args []string) {
	if versionVerbose {
		// Detailed output format
		fmt.Printf("ghqx %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built at: %s\n", buildTime)
		fmt.Printf("go version: %s\n", runtime.Version())
	} else {
		// Simple output format (default)
		fmt.Printf("ghqx %s\n", version)
	}
}
