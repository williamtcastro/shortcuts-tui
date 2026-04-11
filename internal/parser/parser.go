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

var blocklist = map[string]bool{
	"..":    true,
	"...":   true,
	"....":  true,
	"-":     true,
	"ls":    true,
	"ll":    true,
	"la":    true,
	"tree":  true,
	"cat":   true,
	"v":     true,
	"vim":   true,
}

func LoadItems(cfg config.Config) []list.Item {
	items := []list.Item{}

	for _, v := range cfg.Views {
		for _, dir := range v.Dirs {
			files, err := os.ReadDir(dir)
			if err != nil {
				continue
			}

			for _, f := range files {
				if f.IsDir() {
					continue
				}

				path := filepath.Join(dir, f.Name())
				
				// Process based on view type
				if v.Type == "alias" && strings.HasSuffix(f.Name(), ".zsh") {
					file, err := os.Open(path)
					if err != nil {
						continue
					}
					
					category := strings.TrimSuffix(f.Name(), ".zsh")
					category = strings.Title(category)
					
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						line := strings.TrimSpace(scanner.Text())
						matches := aliasRegex.FindStringSubmatch(line)
						if len(matches) >= 3 {
							name := matches[1]
							
							// Filter out blocklisted aliases
							if blocklist[name] {
								continue
							}

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
								ViewName:    v.Name,
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
								ViewName:    v.Name,
							})
						}
					}
				} else if v.Type == "doc" && strings.HasSuffix(f.Name(), ".md") {
					if data, err := os.ReadFile(path); err == nil {
						name := strings.TrimSuffix(f.Name(), ".md")
						name = strings.Title(strings.ReplaceAll(name, "_", " "))
						items = append(items, tui.Item{
							ItemTitle:   name,
							ItemDesc:    "Markdown Guide",
							ItemContent: string(data),
							Category:    "Docs",
							ViewName:    v.Name,
						})
					}
				}
			}
		}
	}

	// Global Neovim guide check (add to any Doc view or first view)
	home, _ := os.UserHomeDir()
	nvimGuide := filepath.Join(home, ".config", "nvim", "NVIM_GUIDE.md")
	if data, err := os.ReadFile(nvimGuide); err == nil {
		targetView := "Docs"
		if len(cfg.Views) > 0 {
			targetView = cfg.Views[0].Name
			for _, v := range cfg.Views {
				if v.Type == "doc" {
					targetView = v.Name
					break
				}
			}
		}

		items = append(items, tui.Item{
			ItemTitle:   "Neovim Guide",
			ItemDesc:    "Markdown Guide",
			ItemContent: string(data),
			Category:    "Docs",
			ViewName:    targetView,
		})
	}

	return items
}
