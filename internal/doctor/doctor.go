package doctor

import (
	"fmt"
	"os/exec"

	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/i18n"
)

// CheckResult は診断結果を保持します
type CheckResult struct {
	Name    string
	OK      bool
	Message string
	Hint    string
}

// Service は環境診断サービスです
type Service struct {
	configLoader *config.Loader
	configPath   string // Add this field
}

// NewService は新しい Service を作成します
func NewService() *Service {
	return &Service{
		configLoader: config.NewLoader(),
		configPath:   "", // Empty means use default search
	}
}

// NewServiceWithConfigPath は指定された設定パスで Service を作成します
func NewServiceWithConfigPath(configPath string) *Service {
	return &Service{
		configLoader: config.NewLoader(),
		configPath:   configPath,
	}
}

// RunChecks はすべての診断を実行します
func (s *Service) RunChecks() []CheckResult {
	return []CheckResult{
		s.CheckConfig(),
		s.CheckGhq(),
		s.CheckGit(),
	}
}

// CheckConfig は設定ファイルを診断します
func (s *Service) CheckConfig() CheckResult {
	_, err := s.configLoader.Load(s.configPath)
	if err != nil {
		return CheckResult{
			Name:    i18n.T("doctor.check.config.name"),
			OK:      false,
			Message: i18n.T("doctor.check.config.fail"),
			Hint:    err.Error(),
		}
	}
	return CheckResult{
		Name:    i18n.T("doctor.check.config.name"),
		OK:      true,
		Message: i18n.T("doctor.check.config.ok"),
	}
}

// CheckGhq は ghq コマンドを診断します
func (s *Service) CheckGhq() CheckResult {
	path, err := exec.LookPath("ghq")
	if err != nil {
		return CheckResult{
			Name:    i18n.T("doctor.check.ghq.name"),
			OK:      false,
			Message: i18n.T("doctor.check.ghq.fail.found"),
			Hint:    i18n.T("doctor.check.ghq.hint.install"),
		}
	}

	cmd := exec.Command(path, "--version")
	if _, err := cmd.Output(); err != nil {
		return CheckResult{
			Name:    i18n.T("doctor.check.ghq.name"),
			OK:      false,
			Message: i18n.T("doctor.check.ghq.fail.exec"),
			Hint:    err.Error(),
		}
	}

	return CheckResult{
		Name:    i18n.T("doctor.check.ghq.name"),
		OK:      true,
		Message: fmt.Sprintf(i18n.T("doctor.check.ghq.ok"), path),
	}
}

// CheckGit は git コマンドを診断します
func (s *Service) CheckGit() CheckResult {
	path, err := exec.LookPath("git")
	if err != nil {
		return CheckResult{
			Name:    i18n.T("doctor.check.git.name"),
			OK:      false,
			Message: i18n.T("doctor.check.git.fail.found"),
			Hint:    i18n.T("doctor.check.git.hint.install"),
		}
	}

	cmd := exec.Command(path, "--version")
	if _, err := cmd.Output(); err != nil {
		return CheckResult{
			Name:    i18n.T("doctor.check.git.name"),
			OK:      false,
			Message: i18n.T("doctor.check.git.fail.exec"),
			Hint:    err.Error(),
		}
	}

	return CheckResult{
		Name:    i18n.T("doctor.check.git.name"),
		OK:      true,
		Message: fmt.Sprintf(i18n.T("doctor.check.git.ok"), path),
	}
}
