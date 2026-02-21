<p align="center">
  <a href="https://github.com/mi8bi/ghqx">
    <img src="./assets/ogp.png" alt="OGP Image" style="max-width:640px;width:100%;height:auto;">
  </a>
</p>

# ghqx - ghq-compatible workspace manager

ghqx extends ghq by managing multiple workspaces (dev/release/sandbox).

[![Build Status](https://github.com/mi8bi/ghqx/actions/workflows/codecov.yml/badge.svg)](https://github.com/mi8bi/ghqx/actions/workflows/codecov.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mi8bi/ghqx)](https://goreportcard.com/report/github.com/mi8bi/ghqx)
[![codecov](https://codecov.io/gh/mi8bi/ghqx/branch/main/graph/badge.svg)](https://codecov.io/gh/mi8bi/ghqx)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest Release](https://img.shields.io/github/v/release/mi8bi/ghqx)](https://github.com/mi8bi/ghqx/releases/latest)
[![GitHub Stars](https://img.shields.io/github/stars/mi8bi/ghqx?style=social)](https://github.com/mi8bi/ghqx/stargazers)


## Features

- **Project status** across all workspaces
- **Shell integration** for quick project navigation with fzf/peco
- **Configuration management** (interactive init, viewing, and TUI editor)
- **Workspace-aware cloning** with `ghqx get`
- **Default workspace mode selection** with `ghqx mode`
- **Project listing** with `ghqx list` for integration with fuzzy finders

## Installation

### Manual Installation (Recommended)

Download the latest binary from the [GitHub Releases](https://github.com/mi8bi/ghqx/releases) page:

#### macOS

```bash
# For Apple Silicon (ARM64)
curl -L https://github.com/mi8bi/ghqx/releases/latest/download/ghqx_Darwin_arm64.tar.gz | tar xz
sudo mv ghqx /usr/local/bin/

# For Intel (x86_64)
curl -L https://github.com/mi8bi/ghqx/releases/latest/download/ghqx_Darwin_x86_64.tar.gz | tar xz
sudo mv ghqx /usr/local/bin/
```

#### Linux

```bash
# For x86_64
curl -L https://github.com/mi8bi/ghqx/releases/latest/download/ghqx_Linux_x86_64.tar.gz | tar xz
sudo mv ghqx /usr/local/bin/

# For ARM64
curl -L https://github.com/mi8bi/ghqx/releases/latest/download/ghqx_Linux_arm64.tar.gz | tar xz
sudo mv ghqx /usr/local/bin/
```

#### Windows

1. Go to the [Releases](https://github.com/mi8bi/ghqx/releases) page
2. Download `ghqx_Windows_x86_64.zip`
3. Extract the archive using File Explorer's "Extract All..." option
4. Move `ghqx.exe` to a directory in your PATH (e.g., `%USERPROFILE%\bin\`)

#### Verify Installation

```bash
ghqx version
```

After installation, create the initial configuration:

```bash
ghqx config init
```

### For Developers

If you want to build from source or contribute to development:

```bash
go install github.com/mi8bi/ghqx/cmd/ghqx@latest
```

**Note**: When installing via `go install`, the `ghqx version` command will not display the correct version information. Use the manual installation method for the full experience.

## Shell Integration

ghqx works seamlessly with fuzzy finders like `fzf` or `peco` for quick project navigation. The `ghqx list` command outputs projects from your default workspace (set with `ghqx mode`), making it easy to integrate with your favorite fuzzy finder.

### Prerequisites

**Install fzf (recommended) or peco:**

```bash
# macOS
brew install fzf

# Windows (PowerShell)
scoop install fzf
# or
choco install fzf

# Linux
sudo apt install fzf  # Debian/Ubuntu
sudo dnf install fzf  # Fedora
```

### Setup Shell Function

Add the following function to your shell configuration file:

**Bash (~/.bashrc) / Zsh (~/.zshrc):**

```bash
ghqxc() {
  local selected
  selected=$(ghqx list --full-path | fzf --height 40% --reverse --border --prompt="Select project: ")
  if [ -n "$selected" ]; then
    cd "$selected"
  fi
}
```

**PowerShell ($PROFILE):**

```powershell
function ghqxc {
  $selected = ghqx list --full-path | fzf --height 40% --reverse --border --prompt="Select project: "
  if ($selected) {
    Set-Location $selected
  }
}
```

**Fish (~/.config/fish/config.fish):**

```fish
function ghqxc
  set selected (ghqx list --full-path | fzf --height 40% --reverse --border --prompt="Select project: ")
  test -n "$selected"; and cd $selected
end
```

### Usage

```bash
# Interactive project selection and navigation
ghqxc

# Or use directly in one line
cd $(ghqx list --full-path | fzf)

# With peco instead of fzf
cd $(ghqx list --full-path | peco)

# List projects from all workspaces (not just default)
cd $(ghqx list --all --full-path | fzf)

# Browse relative paths (for reference, but cd needs full paths)
ghqx list | fzf
```

### How it works

- **`ghqx list`**: Lists projects from your default workspace (set with `ghqx mode`)
- **`ghqx list --all`**: Lists projects from all workspaces (dev/release/sandbox)
- **`ghqx list --full-path`**: Shows absolute paths instead of relative paths
- **`ghqx mode`**: Changes which workspace is used by default for `ghqx list` and `ghqx get`

## Commands

### `ghqx status`
Show the state of all projects across all roots.

```bash
# Compact view (default)
ghqx status

# Verbose view with full paths
ghqx status -v

# Interactive TUI mode
ghqx status --tui
```

Output includes:
- Project name
- Workspace (sandbox/dev/release)
- Git managed status
- Clean/dirty status
- Non-git managed directories are also shown.

### `ghqx list`

List project paths for integration with fuzzy finders.

```bash
# List projects from default workspace (respects ghqx mode)
ghqx list

# List projects from all workspaces
ghqx list --all
ghqx list -a

# Show full absolute paths
ghqx list --full-path
ghqx list -p
```

**Output format:**
- Default: Relative paths from workspace root (like `github.com/user/repo`)
- With `--full-path`: Absolute paths (like `/home/user/ghqx/dev/github.com/user/repo`)
- Default scope: Only projects in your default workspace
- With `--all`: Projects from all workspaces (dev/release/sandbox)

**Integration examples:**
```bash
# Navigate with fzf (use --full-path for cd)
cd $(ghqx list --full-path | fzf)

# Navigate with peco
cd $(ghqx list --full-path | peco)

# Search all workspaces
cd $(ghqx list --all --full-path | fzf)
```

### `ghqx mode`
Select and set the default workspace mode.

This command provides an interactive TUI to choose the default root for `ghqx list` and `ghqx get` operations.

```bash
ghqx mode
```

**Keybindings:**
- **↑↓** or **j/k** - Navigate through options
- **Enter** - Select mode and exit
- **Esc** or **Ctrl+C** - Quit without selecting

### `ghqx get <repository>`
Clones a repository into a specified workspace using `ghq`.

The repository can be specified as:
- Full URL: `https://github.com/user/repo`
- Short form: `github.com/user/repo`
- User/repo: `user/repo` (assumes github.com)

By default, repositories are cloned to the configured default workspace (`ghqx mode` can change this).

```bash
# Clone to default workspace
ghqx get user/repo

# Clone to specific workspace
ghqx get user/repo --workspace dev
ghqx get user/repo --workspace sandbox
ghqx get user/repo --workspace release
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

### `ghqx version`
Display version information.

```bash
ghqx version

# With detailed build information
ghqx version --verbose
```

## Configuration

The configuration file is located at `~/.config/ghqx/config.toml` by default.

Example `config.toml`:

```toml
[roots]
dev = "/Users/username/ghqx/dev"
release = "/Users/username/ghqx/release"
sandbox = "/Users/username/ghqx/sandbox"

[default]
root = "sandbox"
```

- **`[roots]`**: Defines the paths for your different workspaces.
- **`[default]`**:
  - `root`: The default workspace to use for `ghqx list` and `ghqx get` operations.

## Architecture

```
ghqx/
├── cmd/ghqx/          # Thin CLI layer
│   ├── root.go
│   ├── status.go
│   ├── list.go        # Project listing for fzf/peco integration
│   ├── config.go
│   ├── get.go
│   ├── clean.go
│   ├── mode.go
│   └── version.go
├── internal/
│   ├── app/           # Application orchestration
│   ├── config/        # Config loading & validation
│   ├── domain/        # Core models & errors
│   ├── fs/            # Filesystem operations
│   ├── git/           # Git operations
│   ├── ghq/           # ghq command client
│   ├── i18n/          # Internationalization
│   ├── status/        # Status scanning logic
│   ├── tui/           # Main TUI components (used by ghqx status --tui)
│   └── ui/            # CLI output formatting
├── go.mod
├── Makefile
├── .goreleaser.yaml       # Multi-platform release config
├── .github/workflows/     # CI/CD workflows
│   ├── codecov.yml        # Automated tests with coverage
│   └── release.yml        # Automated releases
├── BUILDING.md            # Build instructions
├── CI_CD.md              # CI/CD pipeline documentation
└── README.md
```

## Building

### Quick Build

```bash
# Development build
make build

# Release build with version information
make build-release VERSION=v0.3.0
```

For detailed build instructions, see [BUILDING.md](BUILDING.md).

### Release Process

ghqx uses **GoReleaser** for automated multi-platform releases to GitHub.

To create a release:

```bash
git tag v0.3.0
git push origin v0.3.0
```

This triggers the GitHub Actions [release workflow](.github/workflows/release.yml), which:
- Builds for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64)
- Creates checksums
- Publishes to GitHub Releases

See [CI_CD.md](CI_CD.md) for detailed CI/CD documentation.

## License

[MIT](LICENSE)
