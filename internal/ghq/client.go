package ghq

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/domain"
)

// Client はghq コマンドを実行するクライアント
type Client struct {
	cfg     *config.Config
	timeout time.Duration
}

// NewClient は新しい ghq Client を作成する
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:     cfg,
		timeout: 30 * time.Second, // ghq get は時間がかかる可能性があるため長めに設定
	}
}

// GetOptions は ghq get コマンドのオプション
type GetOptions struct {
	Repository string // リポジトリURL または短縮形
	Workspace  string // 取得先のワークスペース (sandbox/dev/release) - Renamed from Zone
}

// Get は ghq get コマンドを実行してリポジトリを取得する
func (c *Client) Get(opts GetOptions) error {
	// ghq コマンドが利用可能か確認
	if !c.hasGhq() {
		return domain.NewError(
			domain.ErrCodeGitError,
			"ghq command not found",
		).WithHint("Install ghq: https://github.com/x-motemen/ghq")
	}

	// ワークスペースに対応する root を取得 - Renamed from zone
	rootPath, exists := c.cfg.GetRoot(opts.Workspace) // Updated opts.Zone to opts.Workspace
	if !exists {
		return domain.ErrRootNotFound(opts.Workspace) // Updated opts.Zone to opts.Workspace
	}

	// コンテキストとタイムアウト設定
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// ghq get コマンドを構築
	cmd := exec.CommandContext(ctx, "ghq", "get", opts.Repository)
	
	// GHQ_ROOT 環境変数を設定してクローン先を指定
	cmd.Env = append(os.Environ(), "GHQ_ROOT="+rootPath)
	
	// 標準出力・標準エラー出力を親プロセスに接続
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// コマンド実行
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return domain.NewError(
				domain.ErrCodeGitError,
				"ghq get operation timed out",
			).WithHint("Repository may be too large or network is slow")
		}
		return domain.NewErrorWithCause(
			domain.ErrCodeGitError,
			"ghq get failed",
			err,
		).WithHint("Check repository URL and network connection")
	}

	return nil
}

// hasGhq は ghq コマンドが利用可能かチェックする
func (c *Client) hasGhq() bool {
	cmd := exec.Command("ghq", "--version")
	return cmd.Run() == nil
}