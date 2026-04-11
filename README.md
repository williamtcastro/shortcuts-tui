# Shortcuts TUI 🚀

![Shortcuts TUI Demo](docs/demo.png)

**Shortcuts TUI** is a high-performance terminal tool built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea). It provides a beautiful, searchable interface for your shell aliases, functions, and documentation guides, allowing you to explore and execute them with minimal keystrokes.

---

## ⚡️ Quick Start

1. **Install via Homebrew:**
   ```bash
   brew tap williamtcastro/tap
   brew install shortcuts-tui
   ```

2. **Initialize Configuration:**
   Create a config file at `~/.config/shortcuts/config.yaml`. See the [Configuration](#-configuration) section for details.

3. **Launch:**
   Simply run `shortcuts-tui`.

---

## ✨ Features

- **🔍 Deep Search:** Press `/` to instantly filter through titles, descriptions, and full commands.
- **⚡️ Direct Execution:** Hit **Enter** on any alias to run it immediately in your shell.
- **📋 Clipboard Integration:** Press **y** to copy any command to your system clipboard.
- **📖 Markdown Documentation:** Renders `.md` files with beautiful formatting for your guides and cheatsheets.
- **⌨️ Vim-First Navigation:** Full support for `j/k`, `d/u`, and `f/b` motions.
- **📂 Multi-Tab Interface:** Organize your workflows into logical tabs (e.g., "Dev", "Ops", "Guides").
- **🎨 Catppuccin Theme:** Built-in support for the high-contrast Catppuccin Mocha palette.

---

## 🛠 Installation

### Homebrew (Recommended)
```bash
brew tap williamtcastro/tap
brew install shortcuts-tui
```

### From Source
Requires [Go](https://go.dev/doc/install) 1.21+.
```bash
git clone https://github.com/williamtcastro/shortcuts-tui.git
cd shortcuts-tui
go build -o shortcuts-tui ./cmd/shortcuts-tui
mv shortcuts-tui /usr/local/bin/ # Or any directory in your $PATH
```

---

## ⚙️ Configuration

By default, **Shortcuts TUI** looks for `~/.config/shortcuts/config.yaml`.

### Example `config.yaml`
```yaml
# Define your tabs (Views)
views:
  - name: "Aliases"
    type: "alias" # Scans for 'alias name="cmd" # description'
    dirs: 
      - "$HOME/.zsh/aliases"
  - name: "Guides"
    type: "doc" # Renders markdown files
    dirs:
      - "$HOME/docs/cheatsheets"

# Behavioral Settings
pagination: "dots"   # "numeric" (1/3) or "dots" (•●•)
auto_clear: false    # Clear terminal before running a shortcut
auto_exit: false     # Close TUI after running/copying

# Theme (Catppuccin Mocha)
theme:
  primary: "#a6e3a1"   # Tab & Header text
  secondary: "#6c7086" # Dimmed text & Borders
  text: "#cdd6f4"      # Main content
  accent: "#f9e2af"    # Pointer & Active selection
  mauve: "#cba6f7"     # Item descriptions
  flamingo: "#f2cdcd"  # Search prompt & Cursor
```

### Alias Structure
To make your aliases searchable and descriptive, use this format in your `.zsh` or `.bash` files:
```zsh
alias gs="git status" # Show current git status
```
*The TUI parses the text after the `#` as the description.*

---

## ⌨️ Keybindings

### Main List View
| Key | Action |
|-----|--------|
| `j` / `k` | Move selection down / up |
| `Tab` / `l` | Next tab |
| `Shift+Tab` / `h` | Previous tab |
| `/` | Start searching |
| `Enter` | **Run Alias** or **View Document** |
| `y` | Copy command to clipboard |
| `q` / `Ctrl+C` | Quit |

### Document Viewer
| Key | Action |
|-----|--------|
| `j` / `k` | Scroll 1 line |
| `d` / `u` | Scroll half page |
| `f` / `b` | Scroll full page |
| `g` / `G` | Jump to top / bottom |
| `q` / `Esc` | Return to list |

---

## 📝 About
This project was developed and documented with the assistance of **Gemini CLI**. It is designed to be privacy-conscious and respects your local environment:
- No hardcoded home directory paths (uses `$HOME`).
- Executes commands using your default `$SHELL`.
- Operates entirely offline with no telemetry.

## 📄 License
MIT
