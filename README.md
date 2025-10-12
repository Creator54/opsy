# opsy

TUI for managing executable SOPs.

## Installation

```bash
go build .
./opsy        # Launch TUI
```

## Quick Start

Create Markdown SOPs in `~/.opsy/sops/`:

```
~/.opsy/sops/
├── infra/deploy-nginx.md
└── database/backup-postgres.md
```

Opsy executes bash/shell code blocks from your Markdown files. Example:

```markdown
# Deploy Nginx

## Check Status
​```bash
curl -I localhost
​```

## Install
​```bash
which nginx || sudo apt install nginx -y
​```

...
```

## Key Bindings

### Browse Mode
- `↑` `↓` - Navigate
- `Enter` - Open file/directory
- `←` - Back to parent
- `h` - Home directory
- `l` - View logs
- `q` - Quit

### Execute Mode
- `↑` `↓` - Navigate steps
- `Enter` - Execute current step
- `e` - Edit command before execution
- `s` - Skip current step
- `l` - View logs
- `q` - Back to browser

### Logs Mode
- `↑` `↓` - Navigate logs
- `Enter` - View selected log
- `←` - Back
- `q` - Exit logs

## License

MIT
