# Shortcuts TUI 🚀

A high-performance Terminal User Interface (TUI) for exploring, searching, and executing shell aliases and documentation guides. Built with Go and the Charm Bracelet [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

---

> **Note:** This project was developed and documented with the assistance of AI (**Gemini CLI**).

## Features

- **🔍 Deep Search:** Press `/` to search through titles, descriptions, **full commands**, and document content.
- **⚡️ Direct Execution:** Press **Enter** on any alias to execute it immediately in your shell.
- **📖 Markdown Support:** Renders documentation guides with beautiful formatting.
- **⌨️ Vim Keybindings:** Full support for Vim motions for navigation.
- **📂 Multi-Tab Interface:** Separate views for Aliases/Functions and Documentation Guides via `config.yaml`.
- **🔄 Infinite Navigation:** Tabs and lists wrap around seamlessly when you reach the end.
- **📱 Compact UI:** High-density listing with title and description on a single line.
- **🎨 Catppuccin Mocha:** Beautiful, high-contrast theme out of the box.

## How to Use

### 1. Structure your Aliases
Shortcuts TUI looks for ZSH alias definitions with optional trailing comments for descriptions:

```zsh
alias gs="git status" # Show current git status
```

### 2. Configuration (`config.yaml`)
By default, the TUI looks for a config file at `~/.config/shortcuts/config.yaml`.

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

# Customize your TUI colors (Catppuccin Mocha defaults shown)
theme:
  primary: "#a6e3a1"   # Green (Tabs/Headers)
  secondary: "#6c7086" # Overlay0 (Dimmed text)
  text: "#cdd6f4"      # Text (Default text)
  accent: "#f9e2af"    # Yellow (Active cursor)
  mauve: "#cba6f7"     # Mauve (Descriptions)
  flamingo: "#f2cdcd"  # Flamingo (Search bar)
```

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
| `x` | Execute alias (if not in search mode) |
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
