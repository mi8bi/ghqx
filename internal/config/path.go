package config

import (
	"os"

	"github.com/mi8bi/ghqx/internal/domain"
)

// EnsureRootDirectories は設定内の全ルートディレクトリを作成する
// ディレクトリが既に存在する場合は何もしない
func EnsureRootDirectories(cfg *Config) error {
	for rootName, rootPath := range cfg.Roots {
		if err := ensureDirectory(rootPath); err != nil {
			return domain.NewErrorWithCause(
				domain.ErrCodeFSError,
				"Failed to create root directory: "+rootName,
				err,
			).WithHint("Check directory permissions and path validity").
				WithInternal("path: " + rootPath)
		}
	}
	return nil
}

// ensureDirectory はディレクトリを作成する（既に存在する場合は何もしない）
func ensureDirectory(path string) error {
	// ディレクトリが既に存在するか確認
	info, err := os.Stat(path)
	if err == nil {
		// 存在する場合、ディレクトリかどうか確認
		if !info.IsDir() {
			return domain.NewError(
				domain.ErrCodeFSError,
				"Path exists but is not a directory: "+path,
			).WithHint("Remove the file or choose a different path")
		}
		// ディレクトリとして既に存在する場合は成功
		return nil
	}

	// ディレクトリが存在しない場合は作成
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
		return nil
	}

	// その他のエラー（権限不足など）
	return err
}
