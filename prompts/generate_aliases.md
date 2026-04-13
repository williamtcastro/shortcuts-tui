# 🤖 Shortcuts TUI: Alias Generation Prompt

Use this prompt with your preferred AI (Claude, Gemini, etc.) to convert your command history or a list of tasks into high-quality aliases for **Shortcuts TUI**.

---

## 📝 The Prompt

> **System Prompt:**
> You are an expert shell automation assistant specializing in `zsh` and `bash`. Your goal is to help me generate a list of clean, descriptive, and high-performance aliases compatible with **Shortcuts TUI**.
>
> **Task:**
> Analyze the provided command history or list of tasks and output a set of aliases in exactly this format:
> `alias <name>="<command>" # <brief description>`
>
> **Guidelines:**
> 1. **Mnemonic Names:** Use short, memorable names (e.g., `gs` for `git status`, `dcu` for `docker-compose up`).
> 2. **Descriptive Comments:** The text after the `#` is what **Shortcuts TUI** uses for searching. Make it clear and include keywords.
> 3. **Portability:** Use environment variables (like `$HOME`) where appropriate.
> 4. **No Boilerplate:** Output ONLY the alias lines, no introductory text or markdown code blocks unless requested.
> 5. **Categorization:** Group them by purpose (e.g., Git, Docker, System, AI).
>
> **Input:**
> [PASTE YOUR COMMAND HISTORY OR TASK LIST HERE]

---

## ⚡️ Quick Automation (via AI CLIs)

You can automate this directly from your terminal using any of these popular AI CLI tools:

### Using Gemini CLI (Recommended)
```zsh
# Generate aliases from your last 50 commands
history -n -50 | gemini -p "$(cat prompts/generate_aliases.md)" >> ~/.dotfiles/scripts/local/generated.zsh
```

### Using Claude Code
```zsh
# Send history to Claude and save its output
history -n -50 | claude -p "$(cat prompts/generate_aliases.md)" >> ~/.dotfiles/scripts/local/generated.zsh
```

### Using OpenAI CLI
```zsh
# Use gpt-4o to process your command history
history -n -50 | openai api chat.completions.create -m gpt-4o -g user "$(cat prompts/generate_aliases.md)" >> ~/.dotfiles/scripts/local/generated.zsh
```

### Using a Generic `ask` function
If you have a custom AI wrapper (like the one in `ai.zsh`):
```zsh
history -n -50 | ask "$(cat prompts/generate_aliases.md)" >> ~/.dotfiles/scripts/local/generated.zsh
```

---

## 🎨 Example Output

```zsh
alias gs="git status" # Show current git status and staged changes
alias dcu="docker-compose up -d" # Start docker containers in detached mode
alias gcap="git commit -a --amend --no-edit && git push --force" # Amend last commit and force push
alias killp="fuser -k" # Kill process running on a specific port
```
