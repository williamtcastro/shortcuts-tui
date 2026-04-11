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

func LoadItems(cfg config.Config) []list.Item {
	items := []list.Item{}

	// 1. Load Individual Aliases from ZSH files
	files, _ := os.ReadDir(cfg.ScriptsDir)
	aliasRegex := regexp.MustCompile(`^alias\s+([^=]+)="([^"]+)"\s*(?:#\s*(.*))?`)

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".zsh") {
			path := filepath.Join(cfg.ScriptsDir, f.Name())
			file, err := os.Open(path)
			if err != nil {
				continue
			}
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
						ItemContent: fmt.Sprintf("# Alias: %s\n\n**Command:**\n`%s`\n\n**Description:**\n%s", name, cmd, desc),
						IsAlias:     true,
						Command:     cmd,
					})
				}
			}
			file.Close()
		}
	}

	// 2. Load Markdown Guides from docs directory
	docFiles, _ := os.ReadDir(cfg.DocsDir)
	for _, f := range docFiles {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			path := filepath.Join(cfg.DocsDir, f.Name())
			if data, err := os.ReadFile(path); err == nil {
				name := strings.TrimSuffix(f.Name(), ".md")
				name = strings.Title(strings.ReplaceAll(name, "_", " "))
				items = append(items, tui.Item{
					ItemTitle:   name,
					ItemDesc:    "Markdown Guide",
					ItemContent: string(data),
				})
			}
		}
	}

	// 3. Special case for local functions.zsh
	funcPath := filepath.Join(cfg.ScriptsDir, "functions.zsh")
	if data, err := os.ReadFile(funcPath); err == nil {
		items = append(items, tui.Item{
			ItemTitle:   "Shell Functions",
			ItemDesc:    "Raw functions file",
			ItemContent: "```zsh\n" + string(data) + "\n```",
		})
	}

	return items
}
