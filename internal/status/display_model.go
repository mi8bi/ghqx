package status

import (
	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/i18n"
)

// ProjectDisplay は表示用のプロジェクト情報を保持します
type ProjectDisplay struct {
	Repo        string // 短縮されたリポジトリ名 (例: user/repo)
	Zone        string // ゾーン (例: sandbox)
	GitManaged  string // Git管理状態 (例: 管理 / 未管理)
	Status      string // リポジトリの状態 (例: clean / dirty)
	FullPath    string // プロジェクトのフルパス
	RawProject  domain.Project // 元のプロジェクトデータ
}

// NewProjectDisplay は domain.Project から ProjectDisplay を作成します
func NewProjectDisplay(p domain.Project) ProjectDisplay {
	gitManaged := i18n.T("status.git.managed")
	if !p.HasGit {
		gitManaged = i18n.T("status.git.unmanaged")
	}

	status := i18n.T("status.repo.clean")
	if p.Dirty {
		status = i18n.T("status.repo.dirty")
	}
	// For non-git repos, status is not applicable
	if !p.HasGit {
		status = "-"
	}

	return ProjectDisplay{
		Repo:        p.DisplayName,
		Zone:        string(p.Zone),
		GitManaged:  gitManaged,
		Status:      status,
		FullPath:    p.Path,
		RawProject:  p,
	}
}
