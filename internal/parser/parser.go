package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/williamtcastro/shortcuts-tui/internal/config"
	"github.com/williamtcastro/shortcuts-tui/internal/tui"
)

var aliasRegex = regexp.MustCompile(`^alias\s+([^=]+)="([^"]+)"\s*(?:#\s*(.*))?`)

func LoadItems(cfg config.Config) []list.Item {
	items := []list.Item{}

	// 1. Load Individual Aliases from multiple ZSH directories
	for _, scriptDir := range cfg.ScriptsDirs {
		files, err := os.ReadDir(scriptDir)
		if err != nil {
			continue // Skip if directory doesn't exist
		}

		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".zsh") {
				path := filepath.Join(scriptDir, f.Name())
				file, err := os.Open(path)
				if err != nil {
					continue
				}

				// Extract category from filename (e.g., "ai.zsh" -> "Ai")
				category := strings.TrimSuffix(f.Name(), ".zsh")
				category = strings.Title(category)

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())
					matches := aliasRegex.FindStringSubmatch(line)
					if len(matches) >= 3 {
						name := matches[1]
						cmd := matches[2]
						desc := ""
						if len(matches) > 3 {
							desc = matches[3]
						}
						if desc == "" {
							desc = cmd
						}

						items = append(items, tui.Item{
							ItemTitle:   name,
							ItemDesc:    desc,
							ItemContent: fmt.Sprintf("# %s Alias: %s\n\n**Command:**\n`%s`\n\n**Description:**\n%s", category, name, cmd, desc),
							Category:    category,
							IsAlias:     true,
							Command:     cmd,
						})
					}
				}
				file.Close()

				// Special case for local functions file
				if f.Name() == "functions.zsh" {
					if data, err := os.ReadFile(path); err == nil {
						items = append(items, tui.Item{
							ItemTitle:   "Shell Functions",
							ItemDesc:    "Raw functions file",
							ItemContent: "```zsh\n" + string(data) + "\n```",
							Category:    "Functions",
						})
					}
				}
			}
		}
	}

	// 2. Load Markdown Guides from docs directories
	for _, docsDir := range cfg.DocsDirs {
		docFiles, err := os.ReadDir(docsDir)
		if err != nil {
			continue
		}

		for _, f := range docFiles {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
				path := filepath.Join(docsDir, f.Name())
				if data, err := os.ReadFile(path); err == nil {
					name := strings.TrimSuffix(f.Name(), ".md")
					name = strings.Title(strings.ReplaceAll(name, "_", " "))
					items = append(items, tui.Item{
						ItemTitle:   name,
						ItemDesc:    "Markdown Guide",
						ItemContent: string(data),
						Category:    "Docs",
					})
				}
			}
		}
	}

	// 3. Always include Neovim guide if it exists
	home, _ := os.UserHomeDir()
	nvimGuide := filepath.Join(home, ".config", "nvim", "NVIM_GUIDE.md")
	if data, err := os.ReadFile(nvimGuide); err == nil {
		items = append(items, tui.Item{
			ItemTitle:   "Neovim Guide",
			ItemDesc:    "Markdown Guide",
			ItemContent: string(data),
			Category:    "Docs",
		})
	}

	return items
}
