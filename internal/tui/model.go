package tui

import (
	"fmt"
	"io"
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

type Tab int

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
	
	activeTabIndex int
	config         config.Config
	allItem        []list.Item
	
	// Styles
	titleStyle       lipgloss.Style
	infoStyle        lipgloss.Style
	activeTabStyle   lipgloss.Style
	inactiveTabStyle lipgloss.Style
	helpStyle        lipgloss.Style
	selectionStyle   lipgloss.Style
}

type itemDelegate struct {
	activeColor    lipgloss.Color
	inactiveColor  lipgloss.Color
	selectionStyle lipgloss.Style
}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	titleStr := i.Title()
	descStr := i.Description()

	// Base styles
	titleStyle := lipgloss.NewStyle().Width(25)
	descStyle := lipgloss.NewStyle()

	if index == m.Index() {
		// Active item styling
		title := titleStyle.Foreground(d.activeColor).Bold(true).Render("󰄾 " + titleStr)
		desc := descStyle.Foreground(d.inactiveColor).Render(descStr)
		
		// Join them and apply selection style to the whole line
		line := lipgloss.JoinHorizontal(lipgloss.Center, title, " ", desc)
		fmt.Fprint(w, d.selectionStyle.Render(line))
	} else {
		// Inactive item styling
		title := titleStyle.PaddingLeft(2).Render(titleStr)
		desc := descStyle.Foreground(d.inactiveColor).Faint(true).Render(descStr)
		
		line := lipgloss.JoinHorizontal(lipgloss.Center, title, " ", desc)
		fmt.Fprint(w, lipgloss.NewStyle().PaddingLeft(1).Render(line))
	}
}

func New(items []list.Item, cfg config.Config) Model {
	primary := lipgloss.Color(cfg.Theme.PrimaryColor)
	secondary := lipgloss.Color(cfg.Theme.SecondaryColor)
	text := lipgloss.Color(cfg.Theme.TextColor)

	titleStyle := lipgloss.NewStyle().
		Foreground(text).
		Background(primary).
		Padding(0, 1).
		Bold(true)

	infoStyle := lipgloss.NewStyle().
		Foreground(secondary)

	activeTabStyle := lipgloss.NewStyle().
		Foreground(primary).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(primary).
		Padding(0, 2).
		MarginRight(1).
		Bold(true)

	inactiveTabStyle := lipgloss.NewStyle().
		Foreground(secondary).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("none")). // Same height, invisible border
		Padding(0, 2).
		MarginRight(1)

	selectionStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(primary).
		PaddingLeft(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(secondary)

	delegate := itemDelegate{
		activeColor:    primary,
		inactiveColor:  secondary,
		selectionStyle: selectionStyle,
	}

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowTitle(false) // Hide internal title to avoid "double title"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	m := Model{
		list:             l,
		renderer:         renderer,
		titleStyle:       titleStyle,
		infoStyle:        infoStyle,
		activeTabStyle:   activeTabStyle,
		inactiveTabStyle: inactiveTabStyle,
		helpStyle:        helpStyle,
		selectionStyle:   selectionStyle,
		activeTabIndex:   0,
		allItem:          items,
		config:           cfg,
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
		// Deduct 5: 1 (tab row) + 1 (empty line) + 1 (help footer) + 2 (extra padding safety)
		m.list.SetSize(msg.Width-h, msg.Height-v-5)

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
			case "j": m.viewport.LineDown(1); return m, nil
			case "k": m.viewport.LineUp(1); return m, nil
			case "d": m.viewport.HalfPageDown(); return m, nil
			case "u": m.viewport.HalfPageUp(); return m, nil
			case "f": m.viewport.PageDown(); return m, nil
			case "b": m.viewport.PageUp(); return m, nil
			case "g": m.viewport.GotoTop(); return m, nil
			case "G": m.viewport.GotoBottom(); return m, nil
			}
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

		if !m.list.SettingFilter() && len(m.config.Views) > 0 {
			switch msg.String() {
			case "tab", "l", "right":
				if m.activeTabIndex < len(m.config.Views)-1 {
					m.activeTabIndex++
					m.updateListForTab()
					return m, nil
				}
			case "shift+tab", "h", "left":
				if m.activeTabIndex > 0 {
					m.activeTabIndex--
					m.updateListForTab()
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
		header := m.titleStyle.Render(" " + i.Title() + " ")
		footer := m.infoStyle.Render(fmt.Sprintf(" %3.f%% • q/esc: back • j/k: scroll", m.viewport.ScrollPercent()*100))
		return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left, header, "\n", m.viewport.View(), "\n", footer))
	}

	// 1. Tab Row
	var tabs []string
	for i, v := range m.config.Views {
		label := strings.ToUpper(v.Name)
		if i == m.activeTabIndex {
			tabs = append(tabs, m.activeTabStyle.Render(label))
		} else {
			tabs = append(tabs, m.inactiveTabStyle.Render(label))
		}
	}
	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	
	// 2. Help/Footer
	help := m.helpStyle.Render("enter: run/view • x: execute • tab: switch tab • /: search • q: quit")
	
	// Total deduction in WindowSizeMsg was 5 lines.
	// tabRow now has a border (height 1) + 1 line separator + list + 1 line separator + help
	content := lipgloss.JoinVertical(lipgloss.Left, tabRow, "\n", m.list.View(), "\n", help)
	
	return appStyle.Render(content)
}
