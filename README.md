# Opsy - Terminal-based SOP Executor

Opsy is a terminal-based TUI tool inspired by `lazygit` and `lazydocker`. It allows users to browse, execute, and log Standard Operating Procedure (SOP) documents written in Markdown. Commands inside SOPs (like `curl` or shell commands) can be executed interactively, with outputs captured and logged for audit and review.

## Features

### Core Functionality
- 📁 Interactive navigation of SOP directories and files
- 🔍 Detection of executable command blocks in SOPs
- ▶️ Execute, edit, or skip commands directly from the interface
- 📝 Comprehensive logging of all executions with timestamps and outputs
- 📊 Human-readable log format saved in a structured directory

### Modern TUI
- 🎨 Beautiful, professional interface with color-coded status badges
- 📊 Visual progress bar showing SOP completion
- 📦 Bordered command and output blocks for clarity
- ⚡ Immediate output display after command execution
- 🔄 Smooth scrolling with accurate position tracking
- 📏 Full-width layout utilizing all available space
- ⌨️ Vim-style navigation (j/k keys)

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd opsy

# Build the binary
go build .

# Run the application
./opsy run    # Launch the TUI for executing SOPs
./opsy list   # List all available SOPs
```

## Usage

### SOP File Format

SOP files are written in Markdown and can contain executable command blocks:

```markdown
# Deploy Nginx

This SOP deploys nginx to the server.

## Step 1: Check if nginx is running
Check the status of nginx service.

```bash
curl -I localhost
```

## Step 2: Install nginx if not present
Install nginx if it's not already installed.

```bash
which nginx || echo "nginx not found"
```
```

### Key Bindings

#### Browse Mode
- `↑`/`↓` - Navigate files and directories
- `Enter` - Open file or enter directory
- `g` - Go to parent directory
- `H` - Go to home directory
- `?` - Toggle help
- `q` - Quit application

#### Execute Mode
- `↑`/`↓` or `j`/`k` - Navigate between steps
- `Ctrl+u`/`Ctrl+d` - Scroll up/down half page
- `Enter` - Execute current step
- `e` - Edit command
- `s` - Skip step
- `l` - Save execution log
- `?` - Toggle help
- `q` - Return to browser

#### Edit Mode
- `Enter` - Save changes
- `Esc` - Cancel editing

### Default Directory Structure

```
~/.opsy/
├── sops/
│   ├── infra/
│   │   ├── deploy-nginx.md
│   │   └── restart-services.md
│   └── db/
│       └── backup-db.md
└── logs/
    └── <date>/
        └── <sop-folder>/
            └── <sop-name>_<HH-MM-SS>.log.md
```

SOP files are stored in `~/.opsy/sops/` by default, and execution logs are saved in `~/.opsy/logs/` with a date-based directory structure.

## Example SOPs

Example SOP files have been created in the default location for demonstration:
- `~/.opsy/sops/infra/deploy-nginx.md` - Example deployment SOP
- `~/.opsy/sops/db/backup-db.md` - Example backup SOP

## Development

The project follows a modular architecture:
- `internal/parser` - Parses Markdown SOP files
- `internal/executor` - Executes commands safely
- `internal/logger` - Logs execution results
- `internal/tui` - Terminal user interface
- `internal/config` - Configuration management

## License

MIT