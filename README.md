# Shortcuts TUI 🚀

A high-performance Terminal User Interface (TUI) for exploring, searching, and executing shell aliases and documentation guides. Built with Go and the Charm Bracelet [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

---

> **Note:** This project was developed and documented with the assistance of AI (**Gemini CLI**).

## Features

- **🔍 Live Search:** Press `/` to search through all aliases, commands, and descriptions.
- **⚡️ Direct Execution:** Press **Enter** on any alias to execute it immediately in your shell.
- **📖 Markdown Support:** Renders documentation guides with beautiful formatting.
- **⌨️ Vim Keybindings:** Full support for Vim motions for navigation.
- **📂 Multi-Tab Interface:** Separate views for Aliases/Functions and Documentation Guides.
- **🔄 Infinite Navigation:** Tabs and lists wrap around seamlessly when you reach the end.

## How to Use

### 1. Structure your Aliases
Shortcuts TUI looks for ZSH alias definitions with optional trailing comments for descriptions:

```zsh
alias gs="git status" # Show current git status
```

### 2. Configuration (`config.yaml`)
By default, the TUI looks for a config file at `~/.config/shortcuts/config.yaml`. This allows you to customize your views (tabs), script paths, and UI colors.

**Example `config.yaml`:**
```yaml
# Define your custom tabs/views
views:
  - name: "Aliases"
    type: "alias"
    dirs: 
      - "/Users/youruser/dotfiles/scripts"
  - name: "Docs"
    type: "doc"
    dirs:
      - "/Users/youruser/dotfiles/docs"
      - "./docs"
  - name: "Work"
    type: "alias"
    dirs:
      - "/Users/youruser/work/scripts"

# Customize your TUI colors!
theme:
  primary: "#25A065"
  secondary: "#545454"
  text: "#FFFDF5"
```

You can still use environment variables (`SHORTCUTS_SCRIPTS_DIRS`, etc.) if preferred.

### 3. Launch the TUI
Simply run `shortcuts-tui` from your terminal.

## Keybindings

### List View (Main Menu)
| Key | Action |
|-----|--------|
| `j` / `k` | Move selection down / up |
| `Tab` / `l` | Switch to next tab |
| `Shift+Tab` / `h` | Switch to previous tab |
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
go build -o shortcuts-tui ./cmd/shortcuts-tui
mv shortcuts-tui ~/.local/bin/
```

## OSS Publishing & Privacy
This version has been refactored for Open Source use:
- No hardcoded home directory paths.
- Uses `os.UserHomeDir()` and environment variables for configuration.
- Supports any shell defined in `$SHELL` for command execution.

## License
MIT
