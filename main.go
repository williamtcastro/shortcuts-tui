package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#545454")).
			Render
)

type item struct {
	title, desc, content string
	isAlias              bool
	command              string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title + " " + i.desc }

type model struct {
	list         list.Model
	viewport     viewport.Model
	ready        bool
	width        int
	height       int
	renderer     *glamour.TermRenderer
	showViewport bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		m.viewport = viewport.New(msg.Width-h, msg.Height-v-4)
		m.viewport.YPosition = 4

		if !m.ready {
			m.ready = true
		}

	case tea.KeyMsg:
		if m.showViewport {
			switch msg.String() {
			case "esc", "q":
				m.showViewport = false
				return m, nil
			case "j":
				m.viewport.LineDown(1)
				return m, nil
			case "k":
				m.viewport.LineUp(1)
				return m, nil
			case "d":
				m.viewport.HalfPageDown()
				return m, nil
			case "u":
				m.viewport.HalfPageUp()
				return m, nil
			case "f":
				m.viewport.PageDown()
				return m, nil
			case "b":
				m.viewport.PageUp()
				return m, nil
			case "g":
				m.viewport.GotoTop()
				return m, nil
			case "G":
				m.viewport.GotoBottom()
				return m, nil
			}
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			if !m.list.SettingFilter() {
				return m, tea.Quit
			}
		case "x":
			if !m.list.SettingFilter() {
				i, ok := m.list.SelectedItem().(item)
				if ok && i.isAlias {
					return m, runCommand(i.command)
				}
			}
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				if i.isAlias && !m.list.SettingFilter() {
					return m, runCommand(i.command)
				}
				out, _ := m.renderer.Render(i.content)
				m.viewport.SetContent(out)
				m.viewport.GotoTop()
				m.showViewport = true
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func runCommand(command string) tea.Cmd {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "zsh"
	}
	c := exec.Command(shell, "-c", command+"; echo ''; echo 'Press Enter to return...'; read")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return nil
	})
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	if m.showViewport {
		header := titleStyle.Render(m.list.SelectedItem().(item).Title())
		footer := infoStyle(fmt.Sprintf("%3.f%% (q/esc to back, j/k to scroll)", m.viewport.ScrollPercent()*100))
		return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left, header, m.viewport.View(), footer))
	}

	return appStyle.Render(m.list.View())
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func loadItems() []list.Item {
	items := []list.Item{}

	home, _ := os.UserHomeDir()
	
	// Configurable paths via environment variables
	scriptsDir := getEnv("SHORTCUTS_SCRIPTS_DIR", filepath.Join(home, "dotfiles", "scripts"))
	docsDir := getEnv("SHORTCUTS_DOCS_DIR", "./docs")
	
	// 1. Load Individual Aliases from ZSH files
	// Look for any .zsh file in the scripts directory
	files, _ := os.ReadDir(scriptsDir)
	aliasRegex := regexp.MustCompile(`^alias\s+([^=]+)="([^"]+)"\s*(?:#\s*(.*))?`)

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".zsh") {
			path := filepath.Join(scriptsDir, f.Name())
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

					items = append(items, item{
						title:   name,
						desc:    desc,
						content: fmt.Sprintf("# Alias: %s\n\n**Command:**\n`%s`\n\n**Description:**\n%s", name, cmd, desc),
						isAlias: true,
						command: cmd,
					})
				}
			}
			file.Close()
		}
	}

	// 2. Load Markdown Guides from docs directory
	docFiles, _ := os.ReadDir(docsDir)
	for _, f := range docFiles {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			path := filepath.Join(docsDir, f.Name())
			if data, err := os.ReadFile(path); err == nil {
				name := strings.TrimSuffix(f.Name(), ".md")
				name = strings.Title(strings.ReplaceAll(name, "_", " "))
				items = append(items, item{
					title:   name,
					desc:    "Markdown Guide",
					content: string(data),
				})
			}
		}
	}

	// 3. Special case for local functions.zsh if not already parsed as aliases
	funcPath := filepath.Join(scriptsDir, "functions.zsh")
	if data, err := os.ReadFile(funcPath); err == nil {
		items = append(items, item{
			title:   "Shell Functions",
			desc:    "Raw functions file",
			content: "```zsh\n" + string(data) + "\n```",
		})
	}

	return items
}

func main() {
	items := loadItems()

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Shortcuts Explorer"
	l.Styles.Title = titleStyle

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	m := model{
		list:     l,
		renderer: renderer,
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
