package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/williamtcastro/shortcuts-tui/internal/config"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
)

type Model struct {
	list         list.Model
	viewport     viewport.Model
	ready        bool
	width        int
	height       int
	renderer     *glamour.TermRenderer
	showViewport bool
	titleStyle   lipgloss.Style
	infoStyle    func(strings ...string) string
	
	activeTabIndex int
	config         config.Config
	allItem        []list.Item
	
	activeTabStyle   lipgloss.Style
	inactiveTabStyle lipgloss.Style
}

func New(items []list.Item, cfg config.Config) Model {
	primary := lipgloss.Color(cfg.Theme.PrimaryColor)
	secondary := lipgloss.Color(cfg.Theme.SecondaryColor)
	text := lipgloss.Color(cfg.Theme.TextColor)

	titleStyle := lipgloss.NewStyle().
		Foreground(text).
		Background(primary).
		Padding(0, 1)

	infoStyle := lipgloss.NewStyle().
		Foreground(secondary).
		Render

	activeTabStyle := lipgloss.NewStyle().
		Foreground(text).
		Background(primary).
		Padding(0, 2).
		MarginRight(1).
		Bold(true)

	inactiveTabStyle := lipgloss.NewStyle().
		Foreground(primary).
		Padding(0, 2).
		MarginRight(1)

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Shortcuts Explorer"
	l.Styles.Title = titleStyle

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	m := Model{
		list:             l,
		renderer:         renderer,
		titleStyle:       titleStyle,
		infoStyle:        infoStyle,
		activeTabIndex:   0,
		allItem:          items,
		config:           cfg,
		activeTabStyle:   activeTabStyle,
		inactiveTabStyle: inactiveTabStyle,
	}
	
	m.updateListForTab()
	return m
}

func (m *Model) updateListForTab() {
	if len(m.config.Views) == 0 {
		return
	}
	
	activeView := m.config.Views[m.activeTabIndex]
	
	var filtered []list.Item
	for _, item := range m.allItem {
		if i, ok := item.(Item); ok {
			if i.ViewName == activeView.Name {
				filtered = append(filtered, item)
			}
		}
	}
	m.list.SetItems(filtered)
	m.list.ResetSelected()
	m.list.Title = activeView.Name
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-3) // Space for tabs

		m.viewport = viewport.New(msg.Width-h, msg.Height-v-4)
		m.viewport.YPosition = 4

		if !m.ready {
			m.ready = true
		}

	case tea.KeyMsg:
		// Handle Tab switching globally (if not filtering)
		if !m.list.SettingFilter() && len(m.config.Views) > 0 {
			switch msg.String() {
			case "tab", "l", "right":
				m.showViewport = false // Close viewport if open
				m.activeTabIndex = (m.activeTabIndex + 1) % len(m.config.Views)
				m.updateListForTab()
				return m, nil
			case "shift+tab", "h", "left":
				m.showViewport = false // Close viewport if open
				m.activeTabIndex = (m.activeTabIndex - 1 + len(m.config.Views)) % len(m.config.Views)
				m.updateListForTab()
				return m, nil
			}
		}

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

		// Handle infinite list scrolling (wrap around)
		if !m.list.SettingFilter() && len(m.list.Items()) > 0 {
			switch msg.String() {
			case "j", "down":
				if m.list.Index() == len(m.list.Items())-1 {
					m.list.Select(0)
					return m, nil
				}
			case "k", "up":
				if m.list.Index() == 0 {
					m.list.Select(len(m.list.Items()) - 1)
					return m, nil
				}
			}
		}

		switch msg.String() {
		case "ctrl+c", "q":
			if !m.list.SettingFilter() {
				return m, tea.Quit
			}
		case "x":
			if !m.list.SettingFilter() {
				if i, ok := m.list.SelectedItem().(Item); ok && i.IsAlias {
					return m, runCommand(i.Command)
				}
			}
		case "enter":
			if i, ok := m.list.SelectedItem().(Item); ok {
				if i.IsAlias && !m.list.SettingFilter() {
					return m, runCommand(i.Command)
				}
				out, _ := m.renderer.Render(i.ItemContent)
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

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	if m.showViewport {
		i := m.list.SelectedItem().(Item)
		header := m.titleStyle.Render(i.Title())
		footer := m.infoStyle(fmt.Sprintf("%3.f%% (q/esc to back, j/k to scroll)", m.viewport.ScrollPercent()*100))
		return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left, header, m.viewport.View(), footer))
	}

	// Render Dynamic Tabs
	var tabs []string
	for i, v := range m.config.Views {
		label := strings.ToUpper(v.Name)
		if i == m.activeTabIndex {
			tabs = append(tabs, m.activeTabStyle.Render(label))
		} else {
			tabs = append(tabs, m.inactiveTabStyle.Render(label))
		}
	}
	
	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	
	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left, row, "\n", m.list.View()))
}
