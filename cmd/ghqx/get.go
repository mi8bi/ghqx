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
	getTargetWorkspace string // Renamed from getZone
)

var getCmd = &cobra.Command{
	Use:   "get <repository>",
	Short: i18n.T("get.command.short"),
	Long: i18n.T("get.command.long"),
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func init() {
	// Renamed from --zone to --workspace, and getZone to getTargetWorkspace
	getCmd.Flags().StringVar(&getTargetWorkspace, "workspace", "sandbox", i18n.T("get.flag.workspace"))
}

func runGet(cmd *cobra.Command, args []string) error {
	if err := loadApp(); err != nil {
		return err
	}

	repository := args[0]

	// Determine the target workspace (Renamed from targetZone)
	targetWorkspace := getTargetWorkspace
	if !cmd.Flags().Changed("workspace") { // Changed flag name
		// If --workspace flag was not explicitly set, use the default root from config
		targetWorkspace = application.Config.GetDefaultRoot()
	}

	// 既に同じリポジトリが他のワークスペースに存在するかチェック (Renamed from zone)
	if existingWorkspace := checkRepositoryExists(repository); existingWorkspace != "" {
		fmt.Print(ui.FormatWarning(fmt.Sprintf(
			i18n.T("get.repositoryExists"),
			existingWorkspace,
		)))
		fmt.Println(i18n.T("get.continueFetch"))
	}

	// ghq client を使用してクローン
	ghqClient := ghq.NewClient(application.Config)
	
	opts := ghq.GetOptions{
		Repository: repository,
		Workspace:  targetWorkspace, // Updated to Workspace
	}

	// Updated message to use targetWorkspace
	fmt.Printf(i18n.T("get.cloning")+"\n", repository, targetWorkspace)
	
	if err := ghqClient.Get(opts); err != nil {
		return err
	}

	// Updated message to use targetWorkspace
	fmt.Print(ui.FormatSuccess(fmt.Sprintf(
		i18n.T("get.cloneSuccess"),
		repository, targetWorkspace,
	)))

	return nil
}

// checkRepositoryExists は指定したリポジリが既に存在するかチェックする (Renamed from zone)
// 存在する場合はそのワークスペースを返し、存在しない場合は空文字列を返す (Renamed from zone)
func checkRepositoryExists(repository string) string {
	// global application should be loaded by PersistentPreRunE
	projects, err := application.Status.GetAll(status.Options{})
	if err != nil {
		return "" // エラーは無視して続行
	}

	// リポジリ名を抽出（簡易的な実装）
	// 例: "github.com/user/repo" → "repo"
	// 例: "user/repo" → "repo"
	repoName := extractRepoName(repository)

	for _, proj := range projects {
		if contains(proj.Name, repoName) {
			// Updated to WorkspaceType
			return string(proj.WorkspaceType)
		}
	}

	return ""
}

// extractRepoName はリポジリURLから名前を抽出する（簡易版）
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