package i18n

func loadEnglishMessages() {
	RegisterMessages(LocaleEN, map[string]string{
		// Doctor Command
		"doctor.command.short": "Diagnose ghqx environment",
		"doctor.command.long": `doctor checks the presence of necessary settings and commands for ghqx to function correctly.

It diagnoses the following items:
- Config file (~/.config/ghqx/config.toml)
- ghq command
- git command`,
		"doctor.check.config.name":      "config",
		"doctor.check.config.ok":        "Config file loaded successfully",
		"doctor.check.config.fail":      "Config file not found or invalid",
		"doctor.check.ghq.name":         "ghq",
		"doctor.check.ghq.ok":           "ghq found at: %s",
		"doctor.check.ghq.fail.found":   "ghq not found",
		"doctor.check.ghq.fail.exec":    "Failed to execute ghq --version",
		"doctor.check.ghq.hint.install": "Install ghq: https://github.com/x-motemen/ghq",
		"doctor.check.git.name":         "git",
		"doctor.check.git.ok":           "git found at: %s",
		"doctor.check.git.fail.found":   "git not found",
		"doctor.check.git.fail.exec":    "Failed to execute git --version",
		"doctor.check.git.hint.install": "Install git",

		// cd Command
		"cd.command.short": "Select a project or directory and output its path",
		"cd.command.long": `cd displays an interactive TUI to select a project or directory and outputs its full path.
This command cannot directly change your shell's current directory. To do that, you need to use shell integration.`,

		// version Command
		"version.command.short": "Show application version",
		"version.command.long":  "version displays the version information of ghqx.",

		// Errors (Messages and Hints)
		"error.config.notFoundAny.message":        "No configuration file found",
		"error.config.notFoundAny.hint":           "Run 'ghqx config init' to create a config file",
		"error.config.notFoundAt.message":         "Config file not found at specified path",
		"error.config.notFoundAt.hint":            "Check the path provided with --config flag",
		"error.config.invalidToml.message":        "Failed to parse config file",
		"error.config.invalidToml.hint":           "Check the TOML syntax in your config file",
		"error.config.noRoots.message":            "No roots defined in configuration",
		"error.config.noRoots.hint":               "Add at least one root in the [roots] section",
		"error.config.invalidDefaultRoot.message": "Default root does not exist in roots",
		"error.config.invalidDefaultRoot.hint":    "Set default.root to one of the defined roots",

		"error.root.notFound.message":    "Root not found: %s",
		"error.root.notFound.hint":       "Check your config.toml for available roots",
		"error.root.dirNotExist.message": "Root directory does not exist: %s",
		"error.root.dirNotExist.hint":    "Create the directory or update config.toml",

		"error.project.notFound.message":    "Project not found: %s",
		"error.project.notFound.hint":       "Use 'ghqx status' to see all available projects",
		"error.project.nameInvalid.message": "Invalid project name",
		"error.project.nameInvalid.hint":    "Project name contains forbidden characters",

		"error.argument.required": "Argument required",

		"error.git.dirtyRepo.message":     "Repository has uncommitted changes",
		"error.git.dirtyRepo.hint":        "Commit or stash changes, or use --force",
		"error.git.timeout.message":       "Git operation timed out: %s",
		"error.git.commandFailed.message": "Git operation failed: %s",

		"error.fs.readDir.message":   "Failed to read directory",
		"error.fs.createDir.message": "Failed to create directory",
		"error.fs.scanRoot.message":  "Failed to scan root directory",

		// UI Formatter
		"ui.error.prefix":          "Error",
		"ui.error.hintPrefix":      "Hint",
		"ui.error.debugInfoPrefix": "Debug Information",
		"ui.error.internalPrefix":  "Internal",
		"ui.error.causePrefix":     "Cause",
		"ui.success.prefix":        "✓",
		"ui.warning.prefix":        "⚠",
		"ui.info.prefix":           "•",

		// Status display strings
		"status.git.managed":   "Managed",
		"status.git.unmanaged": "Unmanaged",
		"status.repo.clean":    "clean",
		"status.repo.dirty":    "dirty",

		// Status table headers
		"status.header.name":       "Repo",
		"status.header.workspace":  "Workspace", // Renamed from status.header.zone
		"status.header.gitManaged": "GitManaged",
		"status.header.status":     "Status",
		"status.header.root":       "Root",
		"status.header.path":       "Path",

		// Status messages
		"status.message.projectsLoaded": "%d projects loaded",
		"status.message.errorOccurred":  "An error occurred",
		"status.message.reloading":      "Reloading...",

		// TUI Titles
		"status.title.loading": "ghqx status - Loading...",
		"status.title.error":   "ghqx status - Error",
		"status.title.list":    "ghqx status - Project List",
		"status.title.detail":  "ghqx status - Project Detail",

		// TUI Detail View
		"status.detail.basicInfo":  "■ Basic Info",
		"status.detail.name":       "Name",
		"status.detail.path":       "Path",
		"status.detail.workspace":  "Workspace", // Renamed from status.detail.zone
		"status.detail.root":       "Root",
		"status.detail.gitInfo":    "■ Git Info",
		"status.detail.gitManaged": "Git Managed",
		"status.detail.status":     "Status",
		"status.detail.branch":     "Branch",

		// TUI Help
		"status.help.error": "q: Quit | r: Retry",
		"status.help.main":  "↑↓/jk: Move | d: Detail | r: Reload | q: Quit",

		// Selector
		"selector.title":              "Select a project",
		"selector.search.placeholder": "Filter projects...",
		"selector.search.label":       "Search:",
		"selector.search.noMatches":   "No matching projects found.",
		// "selector.help":                 "↑↓: Move | Enter: Select | Esc/q: Quit", Removed this line
		"selector.helpWithPecoSearch": "↑↓: Move | Enter: Select | Esc: Quit", // New key

		"doctor.result.ok":   "[OK]",
		"doctor.result.ng":   "[NG]",
		"doctor.result.hint": "Hint",

		// Config Command
		"config.init.useDefault":         "Using default configuration",
		"config.init.creatingDirs":       "Creating root directories...",
		"config.init.fileCreated":        "Config file created",
		"config.init.summaryHeader":      "Configuration Summary:",
		"config.show.title":              "ghqx Configuration",
		"config.prompt.intro1":           "Interactively create ghqx configuration",
		"config.prompt.intro2":           "Press Enter to use default values",
		"config.prompt.section.roots":    "■ Workspace Roots",
		"config.prompt.path.dev":         "Path for dev root",
		"config.prompt.path.release":     "Path for release root",
		"config.prompt.path.sandbox":     "Path for sandbox root",
		"config.prompt.section.default":  "■ Default Settings",
		"config.prompt.defaultRoot":      "Default root (dev/release/sandbox)",
		"config.summary.section.roots":   "[Roots]",
		"config.summary.section.default": "[Default]",

		// Get Command
		"get.command.short":    "Clone a repository into a workspace",
		"get.command.long":     "Clones the specified repository into a workspace using ghq. The default workspace can be set with `ghqx mode`.",
		"get.repositoryExists": "Repository already exists in %s workspace", // Updated from zone
		"get.continueFetch":    "Continuing fetch...",
		"get.cloning":          "Cloning %s to %s workspace...",          // Updated from zone
		"get.cloneSuccess":     "Successfully cloned %s to %s workspace", // Updated from zone
		"get.flag.workspace":   "target workspace (sandbox/dev/release)", // Updated from get.flag.zone

		// Root Command
		"root.command.short": "ghqx - ghq-compatible workspace manager",
		"root.command.long":  "ghqx extends ghq by managing multiple workspaces (dev/release/sandbox).",
		"root.flag.config":   "config file path",

		// Status Command
		"status.command.short": "Show the state of all projects across all roots",
		"status.command.long":  "Status quickly visualizes workspace state.\n\nProjects are classified by workspace:\n  sandbox\n  dev\n  release\n\nAdditional information:\n  - Git managed or not\n  - Dirty/clean status", // Updated from zone
		"status.flag.verbose":  "show detailed information including paths",
		"status.flag.tui":      "launch interactive TUI mode",

		// Config Command
		"config.command.short":           "Manage ghqx configuration",
		"config.init.command.short":      "Create a default configuration file",
		"config.init.command.long":       "Initialize a new ghqx configuration file.\n\nInteractive mode (default):\n  Prompts for each configuration value.\n  Press Enter to use default values shown in [brackets].\n\nNon-interactive mode (--yes):\n  Creates config with default values immediately.\n\nThe config file will be created at:\n  ~/.config/ghqx/config.toml (Linux/macOS)\n  %USERPROFILE%\\config\\ghqx\\config.toml (Windows)\n\nUse --config to specify a different location.",
		"config.init.flag.yes":           "non-interactive mode: use all defaults",
		"config.show.command.short":      "Show current configuration",
		"config.show.command.long":       "Display the current ghqx configuration in human-readable format.\n\nShows:\n  - All configured roots\n  - Default settings",
		"config.edit.command.short":      "Edit configuration interactively (TUI)",
		"config.edit.command.long":       "Launch an interactive TUI editor for ghqx configuration.\n\nFeatures:\n  - Visual field editor with descriptions\n  - Real-time validation\n\nKeybindings:\n  ↑↓ or j/k  - Navigate fields\n  Enter       - Edit selected field\n  Esc         - Cancel edit\n  Ctrl+S      - Save configuration\n  q           - Quit (warns if unsaved)\n  Ctrl+Q      - Force quit without saving",
		"config.error.fileAlreadyExists": "Config file already exists: %s",

		// Clean Command
		"clean.command.short":          "Reset ghqx configuration and managed information",
		"clean.command.long":           "Resets ghqx to its initial state. Deletes configuration files and all managed repositories.",
		"clean.warning.title":          "Reset ghqx",
		"clean.warning.description":    "This operation is destructive. It will delete ghqx configuration files and all repositories within managed root directories.",
		"clean.warning.targetRoots":    "The following root directories will be deleted:",
		"clean.warning.noConfigFound":  "No configuration file found, so no root directories will be deleted.",
		"clean.warning.confirm":        "Type 'yes' to continue:",
		"clean.aborted":                "Clean up aborted.",
		"clean.deleting.roots":         "Deleting root directories...",
		"clean.deleting.success":       "Deleted",
		"clean.deleting.config":        "Deleting configuration file...",
		"clean.deleting.noConfigFound": "Configuration file not found. Skipping deletion.",
		"clean.deleting.noConfigPath":  "Configuration file path unknown. Skipping deletion.",
		"clean.complete":               "ghqx clean up complete.",

		// Mode Command
		"mode.command.short":  "Switch default workspace mode (dev/release/sandbox)",
		"mode.command.long":   "Interactively selects and sets the default workspace mode (dev, release, or sandbox) for ghqx operations.",
		"mode.selector.title": "Select default workspace mode",
		"mode.selector.help":  "↑↓: Move | Enter: Select | Esc/q: Quit",
		"mode.error.noRoots":  "No roots defined in configuration. Cannot select a mode.",
		"mode.noChange":       "Default mode is already set to the selected one. No change made.",
		"mode.success":        "Default mode set to: ",
		"mode.aborted":        "Mode selection aborted.",
	})
}
