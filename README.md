# ghqx - ghq-compatible workspace manager

ghqx extends ghq by managing multiple workspaces (dev/release/sandbox).

## Features

- **Project status** across all workspaces
- **Shell integration** for `cd` command
- **Configuration management** (interactive init, viewing, and TUI editor)
- **Terminal UI (TUI)** for interactive project listing and navigation
- **Zone-aware cloning** with `ghqx get`

## Installation

```bash
go install github.com/mi8bi/ghqx@latest
```

After installation, create the initial configuration:
```bash
ghqx config init
```

## Commands

### `ghqx status`
Show the state of all projects across all roots.

```bash
# Compact view (default)
ghqx status

# Verbose view with full paths
ghqx status -v

# Launch TUI mode
ghqx status --tui
```

Output includes:
- Project name
- Zone (sandbox/dev/release)
- Git managed status
- Clean/dirty status

### `ghqx tui`
Launch the interactive Terminal UI. This provides a visual project list with keyboard navigation.

```bash
ghqx tui
```

**Keybindings:**
- **↑↓** or **j/k** - Navigate through projects
- **d** - Toggle detail view for selected project
- **r** - Refresh project list
- **q** or **Ctrl+C** - Quit

### `ghqx cd` (Shell Integration)

The `ghqx cd` command launches an interactive TUI to select a project and then prints the selected project's directory path to standard output. To actually change directories, you need to use shell integration.

**1. Bash / Zsh:**

Add the following function to your `.bashrc` or `.zshrc` file:
```bash
ghqx-cd() {
  local path
  path=$(ghqx cd)
  if [ -n "$path" ]; then
    cd "$path"
  fi
}
```
Usage:
```bash
# This will open the TUI to select a project
ghqx-cd
```

**2. PowerShell:**

Add the following function to your PowerShell profile (usually `$PROFILE`):
```powershell
function ghqx-cd {
  $path = (ghqx cd)
  if ($path) {
    Set-Location $path
  }
}
```
Usage:
```powershell
# This will open the TUI to select a project
ghqx-cd
```

### `ghqx get <repository>`
Clones a repository into a specified workspace zone using `ghq`.

The repository can be specified as:
- Full URL: `https://github.com/user/repo`
- Short form: `github.com/user/repo`
- User/repo: `user/repo` (assumes github.com)

By default, repositories are cloned to the `sandbox` zone.

```bash
# Clone to sandbox (default)
ghqx get user/repo

# Clone to dev zone
ghqx get user/repo --zone dev
```

### `ghqx config`
Manages the `ghqx` configuration.

**`ghqx config init`**
Creates a new configuration file interactively.
- Prompts for each setting with defaults in `[brackets]`.
- Automatically creates configured root directories.
- Use `--yes` for non-interactive setup with all default values.

**`ghqx config show`**
Displays the current configuration.

**`ghqx config edit`**
Launches an interactive TUI to edit the configuration file.

### `ghqx clean`
Resets `ghqx` to its initial state by deleting all configuration and managed repositories.

**This is a destructive operation.** It will:
1. Delete the `ghqx` configuration file.
2. Delete all configured root directories (`sandbox`, `dev`, `release`) and all the repositories within them.

The command will ask for explicit confirmation before proceeding.

```bash
ghqx clean
```

### `ghqx doctor`
Checks if the `ghqx` environment is set up correctly, verifying:
- Configuration file existence and validity.
- `ghq` command availability.
- `git` command availability.

## Configuration

The configuration file is located at `~/.config/ghqx/config.toml` by default.

Example `config.toml`:

```toml
[roots]
sandbox = "C:/Users/YourUser/ghqx/sandbox"
dev = "C:/Users/YourUser/ghqx/dev"
release = "C:/Users/YourUser/ghqx/release"

[default]
root = "sandbox"
language = "ja"
```

- **`[roots]`**: Defines the paths for your different workspaces (zones).
- **`[default]`**:
  - `root`: The default root to use for certain operations.
  - `language`: The display language (`en` or `ja`).

## Architecture

```
ghqx/
├── cmd/ghqx/          # Thin CLI layer
│   ├── root.go
│   ├── status.go
│   ├── cd.go
│   ├── config.go
│   ├── get.go
│   ├── tui.go
│   └── clean.go
├── internal/
│   ├── app/           # Application orchestration
│   ├── config/        # Config loading & validation
│   ├── domain/        # Core models & errors
│   ├── fs/            # Filesystem operations
│   ├── git/           # Git operations
│   ├── ghq/           # ghq command client
│   ├── i18n/          # Internationalization
│   ├── selector/      # TUI project selector
│   ├── status/        # Status scanning logic
│   ├── tui/           # Main TUI components
│   └── ui/            # CLI output formatting
├── go.mod
├── Makefile
└── README.md
```

## License

[MIT](LICENSE)