package i18n

func loadJapaneseMessages() {
	RegisterMessages(LocaleJA, map[string]string{
		// Doctor Command
		"doctor.command.short": "ghqx の動作環境を診断します",
		"doctor.command.long": `doctor は ghqx が正しく動作するために必要な設定やコマンドの存在をチェックします。

以下の項目を診断します:
- 設定ファイル (~/.config/ghqx/config.toml)
- ghq コマンド
- git コマンド`,
		"doctor.check.config.name":      "config",
		"doctor.check.config.ok":        "設定ファイルを読み込みました",
		"doctor.check.config.fail":      "設定ファイルが見つからないか、不正です",
		"doctor.check.ghq.name":         "ghq",
		"doctor.check.ghq.ok":           "ghq が見つかりました: %s",
		"doctor.check.ghq.fail.found":   "ghq が見つかりません",
		"doctor.check.ghq.fail.exec":    "ghq --version の実行に失敗しました",
		"doctor.check.ghq.hint.install": "ghq をインストールしてください: https://github.com/x-motemen/ghq",
		"doctor.check.git.name":         "git",
		"doctor.check.git.ok":           "git が見つかりました: %s",
		"doctor.check.git.fail.found":   "git が見つかりません",
		"doctor.check.git.fail.exec":    "git --version の実行に失敗しました",
		"doctor.check.git.hint.install": "git をインストールしてください",

		// cd Command
		"cd.command.short": "プロジェクトまたはディレクトリを選択し、そのパスを出力",
		"cd.command.long": `cd は対話的な TUI を表示してプロジェクトまたはディレクトリを選択し、そのフルパスを出力します。
このコマンドは直接シェルのカレントディレクトリを変更することはできません。そのためには、シェル連携を使用する必要があります。`,

		// Errors (Messages and Hints)
		"error.config.notFoundAny.message":        "設定ファイルが見つかりません",
		"error.config.notFoundAny.hint":           "'ghqx config init' を実行して設定ファイルを作成してください",
		"error.config.notFoundAt.message":         "指定されたパスに設定ファイルが見つかりません",
		"error.config.notFoundAt.hint":            "--config フラグで指定したパスを確認してください",
		"error.config.invalidToml.message":        "設定ファイルの解析に失敗しました",
		"error.config.invalidToml.hint":           "設定ファイルの TOML 構文を確認してください",
		"error.config.noRoots.message":            "設定にルートが定義されていません",
		"error.config.noRoots.hint":               "[roots] セクションに少なくとも1つのルートを追加してください",
		"error.config.invalidDefaultRoot.message": "デフォルトルートが [roots] に存在しません",
		"error.config.invalidDefaultRoot.hint":    "default.root を定義済みルートのいずれかに設定してください",

		"error.root.notFound.message":    "ルートが見つかりません: %s",
		"error.root.notFound.hint":       "config.toml で利用可能なルートを確認してください",
		"error.root.dirNotExist.message": "ルートディレクトリが存在しません: %s",
		"error.root.dirNotExist.hint":    "ディレクトリを作成するか、config.toml を更新してください",

		"error.project.notFound.message":    "プロジェクトが見つかりません: %s",
		"error.project.notFound.hint":       "'ghqx status' で利用可能なプロジェクトを確認してください",
		"error.project.nameInvalid.message": "不正なプロジェクト名です",
		"error.project.nameInvalid.hint":    "プロジェクト名に禁止文字が含まれています",

		"error.argument.required": "引数が必要です",

		"error.git.dirtyRepo.message":     "リポジトリにコミットされていない変更があります",
		"error.git.dirtyRepo.hint":        "変更をコミットまたはスタッシュするか、--force を使用してください",
		"error.git.timeout.message":       "Git 操作がタイムアウトしました: %s",
		"error.git.commandFailed.message": "Git 操作に失敗しました: %s",

		"error.fs.readDir.message":   "ディレクトリの読み込みに失敗しました",
		"error.fs.createDir.message": "ディレクトリの作成に失敗しました",
		"error.fs.scanRoot.message":  "ルートディレクトリのスキャンに失敗しました",

		// UI Formatter
		"ui.error.prefix":          "エラー",
		"ui.error.hintPrefix":      "ヒント",
		"ui.error.debugInfoPrefix": "デバッグ情報",
		"ui.error.internalPrefix":  "内部",
		"ui.error.causePrefix":     "原因",
		"ui.success.prefix":        "✓",
		"ui.warning.prefix":        "⚠",
		"ui.info.prefix":           "•",

		// Status display strings
		"status.git.managed":   "管理",
		"status.git.unmanaged": "未管理",
		"status.repo.clean":    "変更なし",
		"status.repo.dirty":    "変更あり",

		// Status table headers
		"status.header.name":       "Repo",
		"status.header.workspace":  "ワークスペース", // Renamed from status.header.zone
		"status.header.gitManaged": "Git管理",
		"status.header.status":     "状態",
		"status.header.root":       "Root",
		"status.header.path":       "Path",

		// Status messages
		"status.message.projectsLoaded": "プロジェクトを %d 個読み込みました",
		"status.message.errorOccurred":  "エラーが発生しました",
		"status.message.reloading":      "再読み込み中...",

		// TUI Titles
		"status.title.loading": "ghqx status - 読み込み中...",
		"status.title.error":   "ghqx status - エラー",
		"status.title.list":    "ghqx status - プロジェクト一覧",
		"status.title.detail":  "ghqx status - プロジェクト詳細",

		// TUI Detail View
		"status.detail.basicInfo":  "■ 基本情報",
		"status.detail.name":       "名前",
		"status.detail.path":       "パス",
		"status.detail.workspace":  "ワークスペース", // Renamed from status.detail.zone
		"status.detail.root":       "ルート",
		"status.detail.gitInfo":    "■ Git 情報",
		"status.detail.gitManaged": "Git管理",
		"status.detail.status":     "状態",
		"status.detail.branch":     "ブランチ",

		// TUI Help
		"status.help.error": "q: 終了 | r: 再試行",
		"status.help.main":  "↑↓/jk: 移動 | d: 詳細 | r: 再読み込み | q: 終了",

		// Selector
		"selector.title":              "プロジェクトを選択してください",
		"selector.search.placeholder": "プロジェクトをフィルタリング...",
		"selector.search.label":       "検索:",
		"selector.search.noMatches":   "一致するプロジェクトは見つかりませんでした。",
		// "selector.help":                 "↑↓: 移動 | Enter: 選択 | Esc/q: 終了", Removed this line
		"selector.helpWithPecoSearch": "↑↓: 移動 | Enter: 選択 | Esc: 終了", // New key

		"doctor.result.ok":   "[OK]",
		"doctor.result.ng":   "[NG]",
		"doctor.result.hint": "ヒント",

		// Config Command
		"config.init.useDefault":         "デフォルト設定を使用しています",
		"config.init.creatingDirs":       "ルートディレクトリを作成中...",
		"config.init.fileCreated":        "設定ファイルを作成しました",
		"config.init.summaryHeader":      "設定内容:",
		"config.show.title":              "ghqx 設定",
		"config.prompt.intro1":           "ghqx 設定を対話的に作成します",
		"config.prompt.intro2":           "各項目でEnterを押すとデフォルト値を使用します",
		"config.prompt.section.roots":    "■ ワークスペースルート",
		"config.prompt.path.dev":         "dev ルートのパス",
		"config.prompt.path.release":     "release ルートのパス",
		"config.prompt.path.sandbox":     "sandbox ルートのパス",
		"config.prompt.section.default":  "■ デフォルト設定",
		"config.prompt.defaultRoot":      "デフォルトルート (dev/release/sandbox)",
		"config.summary.section.roots":   "[Roots]",
		"config.summary.section.default": "[Default]",

		// Get Command
		"get.command.short":    "リポジトリを取得し、ワークスペースにクローンします",
		"get.command.long":     "指定されたリポジトリをghqを使用してワークスペースにクローンします。デフォルトのワークスペースはghqx modeで設定できます。",
		"get.repositoryExists": "リポジトリは既に %s ワークスペースに存在します", // Updated from zone
		"get.continueFetch":    "取得を続行します...",
		"get.cloning":          "リポジトリ %s を %s ワークスペースにクローンしています...", // Updated from zone
		"get.cloneSuccess":     "%s を %s ワークスペースにクローンしました",           // Updated from zone
		"get.flag.workspace":   "ターゲットワークスペース (sandbox/dev/release)", // Updated from get.flag.zone

		// Root Command
		"root.command.short": "ghqx - ghq互換ワークスペースマネージャー",
		"root.command.long":  "ghqx は、複数のワークスペース (dev/release/sandbox) を管理することで ghq を拡張します。",
		"root.flag.config":   "設定ファイルのパス",

		// Status Command
		"status.command.short": "すべてのルートにおける全プロジェクトの状態を表示",
		"status.command.long":  "status はワークスペースの状態を素早く可視化します。\n\nプロジェクトはワークスペースによって分類されます:\n  sandbox\n  dev\n  release\n\n追加情報:\n  - Git管理されているか\n  - Dirty/clean 状態", // Updated from zone
		"status.flag.verbose":  "パスを含む詳細情報を表示",
		"status.flag.tui":      "対話型 TUI モードを起動",

		// Config Command
		"config.command.short":           "ghqx の設定を管理",
		"config.init.command.short":      "デフォルト設定ファイルを作成",
		"config.init.command.long":       "新しい ghqx 設定ファイルを初期化します。\n\n対話モード (デフォルト):\n  各設定値の入力を求めます。\n  [ブラケット] 内に表示されるデフォルト値を使用するには Enter を押します。\n\n非対話モード (--yes):\n  デフォルト値を使用してすぐに設定を作成します。\n\n設定ファイルは以下に作成されます:\n  ~/.config/ghqx/config.toml (Linux/macOS)\n  %USERPROFILE%\\config\\ghqx\\config.toml (Windows)\n\n異なる場所を指定するには --config を使用します。",
		"config.init.flag.yes":           "非対話モード: すべてデフォルト値を使用",
		"config.show.command.short":      "現在の設定を表示",
		"config.show.command.long":       "現在の ghqx 設定を人間が読みやすい形式で表示します。\n\n表示内容:\n  - 設定されているすべてのルート\n  - デフォルト設定",
		"config.edit.command.short":      "設定を対話的に編集 (TUI)",
		"config.edit.command.long":       "ghqx 設定の対話型 TUI エディターを起動します。\n\n機能:\n  - 説明付きの視覚的なフィールドエディター\n  - リアルタイム検証\n\nキーバインド:\n  ↑↓ or j/k  - フィールドをナビゲート\n  Enter       - 選択したフィールドを編集\n  Esc         - 編集をキャンセル\n  Ctrl+S      - 設定を保存\n  q           - 終了 (未保存の場合は警告)\n  Ctrl+Q      - 保存せずに強制終了",
		"config.error.fileAlreadyExists": "設定ファイルが既に存在します: %s",

		// Clean Command
		"clean.command.short":          "ghqx の設定や管理情報をリセット",
		"clean.command.long":           "ghqx を初期状態に戻します。設定ファイルと、管理下の全リポジトリを削除します。",
		"clean.warning.title":          "ghqx のリセット",
		"clean.warning.description":    "この操作は破壊的です。ghqx の設定ファイルと、すべてのルートディレクトリ内のリポジトリが削除されます。",
		"clean.warning.targetRoots":    "以下のルートディレクトリが削除されます:",
		"clean.warning.noConfigFound":  "設定ファイルが見つからないため、ルートディレクトリは削除されません。",
		"clean.warning.confirm":        "続行するには 'yes' と入力してください:",
		"clean.aborted":                "クリーンアップを中止しました。",
		"clean.deleting.roots":         "ルートディレクトリを削除中...",
		"clean.deleting.success":       "削除完了",
		"clean.deleting.config":        "設定ファイルを削除中...",
		"clean.deleting.noConfigFound": "設定ファイルが見つかりません。削除をスキップします。",
		"clean.deleting.noConfigPath":  "設定ファイルのパスが不明です。削除をスキップします。",
		"clean.complete":               "ghqx クリーンアップが完了しました。",

		// Mode Command
		"mode.command.short":  "デフォルトのワークスペースモードを切り替えます (dev/release/sandbox)",
		"mode.command.long":   "ghqx 操作のデフォルトワークスペースモード (dev, release, または sandbox) を対話的に選択し設定します。",
		"mode.selector.title": "デフォルトのワークスペースモードを選択してください",
		"mode.selector.help":  "↑↓: 移動 | Enter: 選択 | Esc/q: 終了",
		"mode.error.noRoots":  "設定にルートが定義されていません。モードを選択できません。",
		"mode.noChange":       "デフォルトモードは既に選択されたモードに設定されています。変更はありません。",
		"mode.success":        "デフォルトモードを次のものに設定しました: ",
		"mode.aborted":        "モード選択は中止されました。",
	})
}
