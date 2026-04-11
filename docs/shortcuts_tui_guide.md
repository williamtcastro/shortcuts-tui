# Shortcuts TUI Guide 🚀

The **Shortcuts TUI** is your central hub for managing your development environment's efficiency.

## Interactive Mode

- **Search:** Press `/` and type anything (name, command, or description). The list filters in real-time.
- **Run:** Press **Enter** on any alias to execute it. It runs in a subshell, so you don't lose your place.
- **View:** For non-alias entries (like this guide), **Enter** opens the full view.

## Keybindings Reference

### List Navigation
- `j` / `k` (or arrows): Select item
- `/`: Start filtering
- `Esc`: Clear filter
- `Enter`: Execute/Open
- `q`: Quit

### Document Navigation (Viewport)
- `j` / `k`: Scroll 1 line
- `d` / `u`: Scroll half page
- `f` / `b`: Scroll full page
- `g` / `G`: Top / Bottom
- `q` / `Esc`: Back to list

## How Parsing Works

The TUI uses a regex to find all lines starting with `alias` in your configured `.zsh` files. It looks for:
1. Alias name (after the space)
2. Command (between the quotes)
3. Description (after the `#` symbol)

Keep your aliases tidy with comments to make them searchable in this TUI!
