package main

import (
	"fmt"

	"github.com/mi8bi/ghqx/internal/ghq"
	"github.com/mi8bi/ghqx/internal/i18n"
	"github.com/mi8bi/ghqx/internal/status"
	"github.com/mi8bi/ghqx/internal/ui"
	"github.com/spf13/cobra"
)

var (
	getZone string
)

var getCmd = &cobra.Command{
	Use:   "get <repository>",
	Short: i18n.T("get.command.short"),
	Long: i18n.T("get.command.long"),
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func init() {
	getCmd.Flags().StringVar(&getZone, "zone", "sandbox", i18n.T("get.flag.zone"))
}

func runGet(cmd *cobra.Command, args []string) error {
	if err := loadApp(); err != nil {
		return err
	}

	repository := args[0]

	// 既に同じリポジトリが他のzoneに存在するかチェック
	if existingZone := checkRepositoryExists(repository); existingZone != "" {
		fmt.Print(ui.FormatWarning(fmt.Sprintf(
			i18n.T("get.repositoryExists"),
			existingZone,
		)))
		fmt.Println(i18n.T("get.continueFetch"))
	}

	// ghq client を使用してクローン
	ghqClient := ghq.NewClient(application.Config)
	
	opts := ghq.GetOptions{
		Repository: repository,
		Zone:       getZone,
	}

	fmt.Printf(i18n.T("get.cloning")+"\n", repository, getZone)
	
	if err := ghqClient.Get(opts); err != nil {
		return err
	}

	fmt.Print(ui.FormatSuccess(fmt.Sprintf(
		i18n.T("get.cloneSuccess"),
		repository, getZone,
	)))

	return nil
}

// checkRepositoryExists は指定したリポジトリが既に存在するかチェックする
// 存在する場合はそのzoneを返し、存在しない場合は空文字列を返す
func checkRepositoryExists(repository string) string {
	// global application should be loaded by PersistentPreRunE
	projects, err := application.Status.GetAll(status.Options{})
	if err != nil {
		return "" // エラーは無視して続行
	}

	// リポジトリ名を抽出（簡易的な実装）
	// 例: "github.com/user/repo" → "repo"
	// 例: "user/repo" → "repo"
	repoName := extractRepoName(repository)

	for _, proj := range projects {
		if contains(proj.Name, repoName) {
			return string(proj.Zone)
		}
	}

	return ""
}

// extractRepoName はリポジトリURLから名前を抽出する（簡易版）
func extractRepoName(repository string) string {
	// 最後の / 以降を取得
	for i := len(repository) - 1; i >= 0; i-- {
		if repository[i] == '/' {
			return repository[i+1:]
		}
	}
	return repository
}

// contains はaにbが含まれるかチェックする
func contains(a, b string) bool {
	return len(a) >= len(b) && (a == b || len(a) > len(b) && 
		(a[len(a)-len(b)-1] == '/' && a[len(a)-len(b):] == b))
}
