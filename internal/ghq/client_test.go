package ghq

import (
	"os"
	"testing"

	"github.com/mi8bi/ghqx/internal/config"
)

func TestGetWithoutGhqReturnsError(t *testing.T) {
	cfg := &config.Config{Roots: map[string]string{"sandbox": "/tmp"}, Default: config.DefaultConfig{Root: "sandbox"}}
	c := NewClient(cfg)

	// Ensure ghq cannot be found by clearing PATH
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", origPath)

	err := c.Get(GetOptions{Repository: "github.com/user/repo", Workspace: "sandbox"})
	if err == nil {
		t.Fatalf("expected error when ghq is not available")
	}
}
