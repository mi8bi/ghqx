# ghqx - ghq-compatible workspace manager

ghqx extends ghq by managing multiple workspaces (dev/release/sandbox).

## Features

- **Project status** across all workspaces
- **Shell integration** for `cd` command
- **Configuration management** (interactive init, viewing, and TUI editor)
- **Zone-aware cloning** with `ghqx get`
- **Default workspace mode selection** with `ghqx mode`

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
```

Output includes:
- Project name
- Zone (sandbox/dev/release)
- Git managed status
- Clean/dirty status
- Non-git managed directories are also shown.

### `ghqx cd` (Shell Integration)

`ghqx cd` launches an interactive Terminal UI to select a project or directory and then prints its full path to standard output. This command cannot directly change your shell's current directory. To do that, you need to use shell integration as described below.

**Keybindings:**
- **↑↓** or **j/k** - Navigate through projects
- **/** - Start searching
- **Enter** - Select project and exit
- **Esc** or **Ctrl+C** - Quit without selecting

**1. Bash / Zsh:**

Add the following function to your `.bashrc` or `.zshrc` file:
```bash
ghqxc() {
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
ghqxc
```

**2. PowerShell:**

Add the following function to your PowerShell profile (usually `$PROFILE`):
```powershell
function ghqxc {
  $path = (ghqx cd)
  if ($path) {
    Set-Location $path
  }
}
```
Usage:
```powershell
# This will open the TUI to select a project
ghqxc
```

### `ghqx mode`
Select and set the default workspace mode.

This command provides an interactive TUI to choose the default root for certain `ghqx` operations (e.g., `ghqx get` without specifying `--zone`).

**Keybindings:**
- **↑↓** or **j/k** - Navigate through options
- **Enter** - Select mode and exit
- **Esc** or **Ctrl+C** - Quit without selecting

### `ghqx get <repository>`
Clones a repository into a specified workspace zone using `ghq`.

The repository can be specified as:
- Full URL: `https://github.com/user/repo`
- Short form: `github.com/user/repo`
- User/repo: `user/repo` (assumes github.com)

By default, repositories are cloned to the configured default zone (`ghqx mode` can change this).

```bash
# Clone to default zone (e.g., dev if configured)
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
Launches an interactive TUI to edit the configuration file. The default root is selected via a TUI.

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
```

- **`[roots]`**: Defines the paths for your different workspaces (zones).
- **`[default]`**:
  - `root`: The default root to use for certain operations.

## Architecture

```
ghqx/
├── cmd/ghqx/          # Thin CLI layer
│   ├── root.go
│   ├── status.go
│   ├── cd.go
│   ├── config.go
│   ├── get.go
│   ├── clean.go
│   └── mode.go
├── internal/
│   ├── app/           # Application orchestration
│   ├── config/        # Config loading & validation
│   ├── domain/        # Core models & errors
│   ├── fs/            # Filesystem operations
│   ├── git/           # Git operations
│   ├── ghq/           # ghq command client
│   ├── i18n/          # Internationalization
│   ├── selector/      # TUI project selector (used by ghqx cd)
│   ├── status/        # Status scanning logic
│   ├── tui/           # Main TUI components (used by ghqx status --tui)
│   └── ui/            # CLI output formatting
├── go.mod
├── Makefile
└── README.md
```

## License

[MIT](LICENSE)
