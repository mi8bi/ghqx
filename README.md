# ghqx - ghq-compatible workspace lifecycle manager

ghqx extends ghq by managing multiple workspaces (dev/release/sandbox) and supporting lifecycle operations.

## Features

### Phase 1: Core CLI
- Project status across all workspaces
- Promote/undo operations
- Shell integration for cd
- Configuration management

### Phase 2: Enhanced Information
- Zone-based classification (sandbox/dev/release)
- Git worktree support
- Centralized error handling with user/internal separation
- Verbose and compact output modes

### Phase 3: Terminal UI (TUI)
- **Interactive terminal interface** powered by Bubble Tea
- Visual project list with keyboard navigation
- Direct promote and undo operations from TUI
- Real-time status updates
- Japanese error messages
- Detail view for selected projects

### Phase 4: Config Management
- **Interactive config creation** with prompts
- **TUI config editor** with real-time validation
- Config viewing and inspection
- Unified config design (single source of truth)
- Auto-validation before save

## Configuration

### Creating Config

**Interactive mode:**
```bash
ghqx config init
```

Prompts for each setting with defaults in [brackets].

**Non-interactive mode:**
```bash
ghqx config init --yes
```

Uses all default values (for scripts).

### Viewing Config

```bash
ghqx config show
```

Displays current configuration in readable format.

### Editing Config

**TUI Editor:**
```bash
ghqx config edit
```

Interactive editor with:
- Field-by-field editing
- Real-time validation
- Unsaved changes warning
- Japanese error messages

**Keybindings:**
- ↑↓ or j/k - Navigate fields
- Enter - Edit field
- Esc - Cancel edit
- Ctrl+S - Save
- q - Quit (warns if unsaved)
- Ctrl+Q - Force quit

## Installation

```bash
make deps
make build
./bin/ghqx config init
```

## TUI Mode (New in Phase 3)

Launch the interactive terminal interface:

```bash
# From status command
ghqx status --tui

# Or standalone TUI command
ghqx tui
ghqx tui -w  # with worktree counts
```

### TUI Keybindings

- **↑↓** or **j/k** - Navigate through projects
- **d** - Toggle detail view for selected project
- **Enter** - Promote selected project (if eligible)
- **u** - Undo last promotion
- **r** - Refresh project list
- **q** or **Ctrl+C** - Quit

### TUI Features

1. **Visual Project List**
   - Color-coded zones (sandbox/dev/release)
   - Git status indicators (clean/dirty)
   - Worktree counts (optional)
   - Selected row highlighting

2. **Detail View**
   - Full project path
   - Git branch information
   - Promote eligibility and hints
   - Zone and root information

3. **Interactive Operations**
   - One-key promote (Enter)
   - One-key undo (u)
   - Instant feedback with Japanese messages
   - Error hints for failed operations

4. **Smart Validation**
   - Prevents promoting dirty repositories
   - Shows why promotion is disabled
   - Validates zone requirements

## CLI Mode (Traditional)

### Enhanced Status Command

The `status` command now provides comprehensive project information:

```bash
# Compact view (default)
ghqx status

# Verbose view with full paths
ghqx status -v

# Include worktree counts
ghqx status -w

# Filter by root
ghqx status --root=dev

# JSON output
ghqx status --json
```

**Output includes:**
- Project name
- Zone (sandbox/dev/release)
- Git managed status
- Clean/dirty status
- Worktree count (with `-w` flag)
- Full path (with `-v` flag)

### Git Worktree Support

List all worktrees for a project:

```bash
ghqx worktree <project-name>
```

This shows:
- Worktree paths
- Associated branches
- Status (active/bare/locked)

### Improved Error Handling

Phase 2 introduces a centralized error system:

**User-facing errors:**
- Clear, actionable messages
- Hints for resolution
- No internal implementation details

**Debug mode:**
```bash
export GHQX_DEBUG=1
ghqx status  # Shows internal error details
```

All errors are defined in `internal/domain/errors.go` for consistency.

## Installation

```bash
make deps
make build
./bin/ghqx config init
```

## Configuration

Example `~/.config/ghqx/config.toml`:

```toml
[roots]
dev     = "C:/src/ghq-dev"
release = "C:/src/ghq-release"
sandbox = "C:/src/sandbox"

[default]
root = "dev"

[promote]
from = "sandbox"
to   = "dev"
auto_git_init = true
auto_commit   = false

[history]
enabled = true
max = 50
```

## Commands

### Status
```bash
ghqx status              # Show all projects
ghqx status -v           # Verbose mode with paths
ghqx status -w           # Include worktree counts
ghqx status --root=dev   # Filter by root
```

### Promote
```bash
ghqx promote myproject             # Promote to default target
ghqx promote myproject --from=sandbox --to=dev
ghqx promote myproject --force     # Ignore dirty state
ghqx promote myproject --dry-run   # Preview changes
```

### Undo
```bash
ghqx undo              # Revert last promote
ghqx undo --dry-run    # Preview undo operation
```

### Worktree
```bash
ghqx worktree myproject  # List all worktrees
ghqx worktree myproject --json
```

### CD (Shell Integration)
```bash
ghqx cd myproject  # Print project path

# Add to .bashrc/.zshrc:
ghqx-cd() {
  local path=$(ghqx cd "$1")
  if [ -n "$path" ]; then
    cd "$path"
  fi
}
```

### Config
```bash
ghqx config init  # Create default config
```

## Architecture

```
ghqx/
├── cmd/ghqx/          # Thin CLI layer
│   ├── main.go
│   ├── root.go
│   ├── status.go
│   ├── promote.go
│   ├── undo.go
│   ├── cd.go
│   ├── config.go
│   └── worktree.go
├── internal/
│   ├── app/           # Application orchestration
│   ├── config/        # Config loading & validation
│   ├── domain/        # Core models & errors
│   │   ├── error.go   # Error type definition
│   │   ├── errors.go  # Centralized error catalog
│   │   └── models.go  # Domain models
│   ├── fs/            # Filesystem operations
│   ├── git/           # Git operations (with worktree support)
│   ├── promote/       # Promote/undo logic
│   ├── status/        # Status scanning
│   └── ui/            # Output formatting
├── go.mod
├── Makefile
└── config.example.toml
```

## Design Principles

1. **No Panic**: All errors are handled gracefully
2. **Thin CLI**: Command handlers delegate to internal packages
3. **Centralized Errors**: User-facing errors defined in one place
4. **Separation of Concerns**: User vs internal error information
5. **Performance**: Parallel scanning, timeouts on git operations
6. **Safety**: Dry-run mode, dirty state checks, undo history

## Error System

Errors are defined in `internal/domain/errors.go` and provide:

- **Code**: Machine-readable error code
- **Message**: User-facing description
- **Hint**: Actionable advice
- **Internal**: Debug information (only shown with GHQX_DEBUG=1)
- **Cause**: Underlying error (for error wrapping)

Example:
```go
return domain.ErrProjectNotFound(name)
// User sees: "Project not found: myproject"
// Hint: "Use 'ghqx status' to see all available projects"
```

## Performance

- Parallel root scanning
- 150ms timeout on git operations
- Lazy loading of branch info (TUI mode only)
- Optional worktree counting

## Development

```bash
# Run tests
make test

# Format code
make fmt

# Build
make build

# Install to $GOPATH/bin
make install
```

## Future Phases

- **Phase 3**: TUI implementation with Bubble Tea
- **Phase 4**: Plugin system
- **Phase 5**: Advanced ghq integration

## License

See LICENSE file for details.
