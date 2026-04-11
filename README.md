# Shortcuts TUI 🚀

A high-performance Terminal User Interface (TUI) for exploring, searching, and executing shell aliases and documentation guides. Built with Go and the Charm Bracelet [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

---

> **Note:** This project was developed and documented with the assistance of AI (**Gemini CLI**).

## Features

- **🔍 Live Search:** Press `/` to search through all aliases, commands, and descriptions.
- **⚡️ Direct Execution:** Press **Enter** on any alias to execute it immediately in your shell.
- **📖 Markdown Support:** Renders documentation guides with beautiful formatting.
- **⌨️ Vim Keybindings:** Full support for Vim motions for navigation.
- **📂 Automatic Parsing:** Automatically extracts aliases and their comments from your `.zsh` files.

## How to Use

### 1. Structure your Aliases
Shortcuts TUI looks for ZSH alias definitions with optional trailing comments for descriptions:

```zsh
alias gs="git status" # Show current git status
```

### 2. Configure Paths
Set the following environment variables in your `.zshrc` or `.bashrc` to point the TUI to your files:

```zsh
export SHORTCUTS_SCRIPTS_DIR="$HOME/dotfiles/scripts"
export SHORTCUTS_DOCS_DIR="$HOME/dotfiles/docs"
```

- `SHORTCUTS_SCRIPTS_DIR`: Directory containing your `.zsh` files with aliases.
- `SHORTCUTS_DOCS_DIR`: Directory containing `.md` files you want to browse.

### 3. Launch the TUI
Simply run `shortcuts-tui` from your terminal.

## Keybindings

### List View (Main Menu)
| Key | Action |
|-----|--------|
| `j` / `k` | Move selection down / up |
| `/` | Start searching/filtering |
| `Enter` | **Run Alias** or **View Document** |
| `q` / `Ctrl+C` | Quit |

### Viewport (Document Viewer)
| Key | Action |
|-----|--------|
| `j` / `k` | Scroll down / up (one line) |
| `d` / `u` | Scroll down / up (half page) |
| `f` / `b` | Scroll down / up (full page) |
| `g` / `G` | Jump to top / bottom |
| `q` / `Esc` | Return to list |

## Installation

### From Source
```bash
git clone https://github.com/williamtcastro/shortcuts-tui.git
cd shortcuts-tui
go build -o shortcuts-tui main.go
mv shortcuts-tui ~/.local/bin/
```

## OSS Publishing & Privacy
This version has been refactored for Open Source use:
- No hardcoded home directory paths.
- Uses `os.UserHomeDir()` and environment variables for configuration.
- Supports any shell defined in `$SHELL` for command execution.

## License
MIT
