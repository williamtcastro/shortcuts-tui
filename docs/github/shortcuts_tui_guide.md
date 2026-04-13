# Shortcuts TUI Guide 🚀

The **Shortcuts TUI** is a centralized hub for managing and accessing your terminal efficiency tools. This guide explains how to get the most out of it.

---

## 🔍 Searching and Filtering
The search feature is designed to be fast and intuitive. 

1. Press **`/`** to enter search mode.
2. Type any part of a command, its name, or its description.
3. The list filters instantly. 
4. Press **Enter** to execute the selected item or **Esc** to clear the filter.

---

## ⚡️ Executing Shortcuts
When you select an alias and press **Enter**, the TUI:
1. Identifies the shell command.
2. Runs it in a subshell using your default `$SHELL`.
3. Returns you to the TUI (unless `auto_exit` is set to `true` in your config).

*Tip: Use the `auto_clear` setting in your `config.yaml` if you want a clean terminal before each execution.*

---

## 📂 Managing Tabs (Views)
You can organize your tools into multiple tabs. Switch between them using **`Tab`** / **`Shift+Tab`** or **`l`** / **`h`**.

To add a new tab, update your `config.yaml`:

```yaml
views:
  - name: "Git"
    type: "alias"
    dirs: ["~/dotfiles/git"]
  - name: "Guides"
    type: "doc"
    dirs: ["~/docs/cheatsheets"]
```

---

## 🛠 Adding New Aliases
The TUI scans your files for standard ZSH/Bash alias definitions. To include a description:

```bash
alias gc="git commit -v" # Commit changes with verbose output
```
*The text after the `#` will be displayed as the description in the TUI.*

---

## 🗂 Automatic Subdivisions
You can group shortcuts within a single tab by using subdirectories. This is ideal for organizing your tools into `Local`, `Work`, or `GitHub` projects.

1. Create a subfolder in your scripts directory: `~/dotfiles/scripts/work/`.
2. Move your `.zsh` files there: `work/git.zsh`.
3. The TUI will display these items with a subdivision prefix: **`[Work > Git]`**.

*Tip: Subdivisions are indexed for search. Type `work` to instantly filter all work-related tools.*

---

## ⌨️ Navigation Reference

### List View
- `j` / `k`: Navigation
- `Tab` / `l`: Next Tab
- `Shift+Tab` / `h`: Previous Tab
- `/`: Search (Titles, Descriptions, Subdivisions)
- `Enter`: Run/View
- `y`: Copy to Clipboard
- `q`: Quit

### Document View
- `j` / `k`: Scroll 1 line
- `d` / `u`: Half page
- `f` / `b`: Full page
- `g` / `G`: Top / Bottom
- `q` / `Esc`: Back to list
