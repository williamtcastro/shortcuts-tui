# Shortcuts TUI 🚀

![Shortcuts TUI Demo](docs/demo.png)

A high-performance Terminal User Interface (TUI) for exploring, searching, and executing shell aliases and documentation guides.
 Built with Go and the Charm Bracelet [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

## 🚀 Installation

Install via Homebrew (macOS/Linux):

```bash
brew tap williamtcastro/tap
brew install shortcuts-tui
```

---

> **Note:** This project was developed and documented with the assistance of AI (**Gemini CLI**).

## Features

- **🔍 Deep Search:** Press `/` to search through titles, descriptions, **full commands**, and document content.
- **⚡️ Direct Execution:** Press **Enter** on any alias to execute it immediately in your shell.
- **📋 Copy to Clipboard:** Press **y** to copy any alias command to your system clipboard.
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

<details>
<summary><b>Click to expand Configuration Options</b></summary>

### View Configuration
Each view represents a tab in the TUI.
- `name`: The display name of the tab.
- `type`: Either `alias` (for .zsh files) or `doc` (for .md files).
- `dirs`: List of directories to scan.

### Behavioral Settings
- `pagination`: `"numeric"` (1/3) or `"dots"` (•●•).
- `auto_clear`: `true/false` - Clear the terminal before running a shortcut.
- `auto_exit`: `true/false` - Close the TUI immediately after running or copying.

### Theme Colors (Hex codes)
- `primary`: Tabs and headers.
- `secondary`: Dimmed text and borders.
- `text`: Main content text.
- `accent`: Pointer and active selection border.
- `mauve`: Item descriptions.
- `flamingo`: Search bar prompt and cursor.

</details>

**Example `config.yaml`:**
```yaml
# Define your custom tabs/views
views:
  - name: "Aliases"
    type: "alias"
    dirs: 
      - "$HOME/dotfiles/scripts"
  - name: "Docs"
    type: "doc"
    dirs:
      - "$HOME/dotfiles/docs"

pagination: "dots"
auto_clear: false
auto_exit: false

# Catppuccin Mocha Defaults
theme:
  primary: "#a6e3a1"
  secondary: "#6c7086"
  text: "#cdd6f4"
  accent: "#f9e2af"
  mauve: "#cba6f7"
  flamingo: "#f2cdcd"
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
| `y` | Copy alias command to clipboard |
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
