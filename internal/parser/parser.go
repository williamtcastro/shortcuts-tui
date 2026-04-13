package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/williamtcastro/shortcuts-tui/internal/config"
	"github.com/williamtcastro/shortcuts-tui/internal/models"
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
			// Expand environment variables in the directory path
			dir = os.ExpandEnv(dir)
			
			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}

				if info.IsDir() {
					return nil
				}

				// Calculate relative path for subdivision
				relPath, _ := filepath.Rel(dir, path)
				relDir := filepath.Dir(relPath)
				subdivision := ""
				if relDir != "." {
					subdivision = strings.Title(strings.ReplaceAll(relDir, "_", " "))
				}

				fileName := info.Name()
				
				// Process based on view type
				if v.Type == "alias" && strings.HasSuffix(fileName, ".zsh") {
					file, err := os.Open(path)
					if err != nil {
						return nil
					}
					defer file.Close()
					
					category := strings.TrimSuffix(fileName, ".zsh")
					category = strings.Title(strings.ReplaceAll(category, "_", " "))
					
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
							
							items = append(items, models.Item{
								ItemTitle:   name,
								ItemDesc:    desc,
								ItemContent: fmt.Sprintf("# %s Alias: %s\n\n**Command:**\n`%s`\n\n**Description:**\n%s", category, name, cmd, desc),
								Category:    category,
								Subdivision: subdivision,
								ViewName:    v.Name,
								IsAlias:     true,
								Command:     cmd,
							})
						}
					}
					
					// Special case for local functions file
					if fileName == "functions.zsh" {
						if data, err := os.ReadFile(path); err == nil {
							items = append(items, models.Item{
								ItemTitle:   "Shell Functions",
								ItemDesc:    "Raw functions file",
								ItemContent: "```zsh\n" + string(data) + "\n```",
								Category:    "Functions",
								Subdivision: subdivision,
								ViewName:    v.Name,
							})
						}
					}
				} else if v.Type == "doc" && strings.HasSuffix(fileName, ".md") {
					if data, err := os.ReadFile(path); err == nil {
						name := strings.TrimSuffix(fileName, ".md")
						name = strings.Title(strings.ReplaceAll(name, "_", " "))
						items = append(items, models.Item{
							ItemTitle:   name,
							ItemDesc:    "Markdown Guide",
							ItemContent: string(data),
							Category:    "Docs",
							Subdivision: subdivision,
							ViewName:    v.Name,
						})
					}
				}
				return nil
			})
			if err != nil {
				continue
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

		items = append(items, models.Item{
			ItemTitle:   "Neovim Guide",
			ItemDesc:    "Markdown Guide",
			ItemContent: string(data),
			Category:    "Docs",
			ViewName:    targetView,
		})
	}

	// Sort items: Subdivision > Category > Title
	sort.Slice(items, func(i, j int) bool {
		a := items[i].(models.Item)
		b := items[j].(models.Item)

		if a.Subdivision != b.Subdivision {
			return a.Subdivision < b.Subdivision
		}
		if a.Category != b.Category {
			return a.Category < b.Category
		}
		return a.ItemTitle < b.ItemTitle
	})

	return items
}
