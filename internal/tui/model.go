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

// --- Layout Constants ---
const (
	headerHeight = 1
	tabHeight    = 2 // Tabs line + Horizontal separator
	footerHeight = 1
)

// --- Styles ---
type Styles struct {
	App          lipgloss.Style
	Header       lipgloss.Style
	Footer       lipgloss.Style
	ActiveTab    lipgloss.Style
	InactiveTab  lipgloss.Style
	TabSeparator lipgloss.Style
	SelectionBar lipgloss.Style
	Title        lipgloss.Style
	Desc         lipgloss.Style
	Dim          lipgloss.Style
	DocViewport  lipgloss.Style
	DocHeader    lipgloss.Style
}

func DefaultStyles(cfg config.Config) Styles {
	primary := lipgloss.Color(cfg.Theme.PrimaryColor)
	secondary := lipgloss.Color(cfg.Theme.SecondaryColor)
	text := lipgloss.Color(cfg.Theme.TextColor)

	return Styles{
		App: lipgloss.NewStyle().Padding(1, 2),
		Header: lipgloss.NewStyle().
			Foreground(text).
			Background(primary).
			Padding(0, 1).
			Bold(true),
		Footer: lipgloss.NewStyle().
			Foreground(secondary).
			MarginTop(1),
		ActiveTab: lipgloss.NewStyle().
			Foreground(primary).
			Background(lipgloss.Color("236")). // Subtle block highlight
			Padding(0, 2).
			Bold(true),
		InactiveTab: lipgloss.NewStyle().
			Foreground(secondary).
			Padding(0, 2),
		TabSeparator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("238")).
			Border(lipgloss.NormalBorder(), false, false, true, false), // Single clean line
		SelectionBar: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(primary).
			PaddingLeft(1),
		Title: lipgloss.NewStyle().Width(25).Bold(true),
		Desc:  lipgloss.NewStyle().Foreground(secondary),
		Dim:   lipgloss.NewStyle().Foreground(secondary).Faint(true),
		DocViewport: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondary).
			Padding(0, 1),
		DocHeader: lipgloss.NewStyle().
			Foreground(primary).
			Bold(true).
			MarginBottom(1),
	}
}

// --- Delegate ---
type itemDelegate struct {
	styles Styles
	active lipgloss.Color
}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok { return }

	titleStr := i.Title()
	descStr := i.Description()

	if index == m.Index() {
		t := d.styles.Title.Foreground(d.active).Render("󰄾 " + titleStr)
		dsc := d.styles.Desc.Render(descStr)
		line := lipgloss.JoinHorizontal(lipgloss.Top, t, " ", dsc)
		fmt.Fprint(w, d.styles.SelectionBar.Render(line))
	} else {
		t := d.styles.Title.PaddingLeft(2).Render(titleStr)
		dsc := d.styles.Dim.Render(descStr)
		line := lipgloss.JoinHorizontal(lipgloss.Top, t, " ", dsc)
		fmt.Fprint(w, lipgloss.NewStyle().PaddingLeft(1).Render(line))
	}
}

// --- Model ---
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
	styles         Styles
}

func New(items []list.Item, cfg config.Config) Model {
	s := DefaultStyles(cfg)
	delegate := itemDelegate{
		styles: s,
		active: lipgloss.Color(cfg.Theme.PrimaryColor),
	}

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	m := Model{
		list:           l,
		renderer:       renderer,
		styles:         s,
		activeTabIndex: 0,
		allItem:        items,
		config:         cfg,
	}
	
	m.updateListForTab()
	return m
}

func (m *Model) updateListForTab() {
	if len(m.config.Views) == 0 { return }
	activeView := m.config.Views[m.activeTabIndex]
	var filtered []list.Item
	for _, item := range m.allItem {
		if i, ok := item.(Item); ok && i.ViewName == activeView.Name {
			filtered = append(filtered, item)
		}
	}
	m.list.SetItems(filtered)
	m.list.ResetSelected()
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		h, v := m.styles.App.GetFrameSize()
		// Overhead: v (2 padding) + tabs (1) + separator (1) + help (1) + 2 (JoinVertical newlines)
		listHeight := msg.Height - v - 5
		if listHeight < 0 { listHeight = 0 }
		
		m.list.SetSize(msg.Width-h, listHeight)
		m.viewport = viewport.New(msg.Width-h-4, msg.Height-v-8)
		m.ready = true

	case tea.KeyMsg:
		if !m.list.SettingFilter() && len(m.config.Views) > 0 {
			switch msg.String() {
			case "tab", "l", "right":
				if m.activeTabIndex < len(m.config.Views)-1 {
					m.activeTabIndex++
					m.updateListForTab()
					m.showViewport = false
					return m, nil
				}
			case "shift+tab", "h", "left":
				if m.activeTabIndex > 0 {
					m.activeTabIndex--
					m.updateListForTab()
					m.showViewport = false
					return m, nil
				}
			}
		}

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

		switch msg.String() {
		case "ctrl+c", "q":
			if !m.list.SettingFilter() { return m, tea.Quit }
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
	return m, cmd
}

func runCommand(command string) tea.Cmd {
	shell := os.Getenv("SHELL")
	if shell == "" { shell = "zsh" }
	c := exec.Command(shell, "-c", command+"; echo ''; echo 'Press Enter to return...'; read")
	return tea.ExecProcess(c, func(err error) tea.Msg { return nil })
}

func (m Model) View() string {
	if !m.ready { return "\n  Initializing..." }

	if m.showViewport {
		i := m.list.SelectedItem().(Item)
		header := m.styles.DocHeader.Render("󰧮 " + i.Title())
		view := m.styles.DocViewport.Render(m.viewport.View())
		helpText := m.styles.Dim.Render(fmt.Sprintf(" %3.f%% • q/esc: back • j/k: scroll", m.viewport.ScrollPercent()*100))
		footer := m.styles.Footer.Render(helpText)
		return m.styles.App.Render(lipgloss.JoinVertical(lipgloss.Left, header, view, footer))
	}

	// 1. Render Tab Bar
	var tabs []string
	for i, v := range m.config.Views {
		label := strings.ToUpper(v.Name)
		if i == m.activeTabIndex {
			tabs = append(tabs, m.styles.ActiveTab.Render(label))
		} else {
			tabs = append(tabs, m.styles.InactiveTab.Render(label))
		}
	}
	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	
	// 2. Full-width separator line (bulletproof approach)
	h, _ := m.styles.App.GetFrameSize()
	fullWidth := m.width - h
	if fullWidth < 0 { fullWidth = 0 }
	
	// Create a simple solid line for the separator
	tabSeparator := m.styles.Dim.Render(strings.Repeat("─", fullWidth))

	// 3. Help Footer
	help := m.styles.Footer.Render("enter: run/view • x: exec • tab: switch • /: filter • q: quit")
	
	// 4. Final Assembly (JoinVertical adds 1 newline between each element)
	content := lipgloss.JoinVertical(lipgloss.Left, tabRow, tabSeparator, m.list.View(), help)
	
	return m.styles.App.Render(content)
}
