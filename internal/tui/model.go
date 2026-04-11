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
	tabHeight      = 2 // Text + Border
	footerHeight   = 1
	verticalMargin = 2 // Padding around the app
)

// --- Styles ---
type Styles struct {
	App          lipgloss.Style
	Header       lipgloss.Style
	Footer       lipgloss.Style
	ActiveTab    lipgloss.Style
	InactiveTab  lipgloss.Style
	TabBorder    lipgloss.Style
	SelectionBar lipgloss.Style
	Title        lipgloss.Style
	Desc         lipgloss.Style
	Dim          lipgloss.Style
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
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(primary).
			Padding(0, 2).
			Bold(true),
		InactiveTab: lipgloss.NewStyle().
			Foreground(secondary).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("0")). // Bottom border matches background
			Padding(0, 2),
		TabBorder: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("238")),
		SelectionBar: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(primary).
			PaddingLeft(1),
		Title: lipgloss.NewStyle().Width(25).Bold(true),
		Desc:  lipgloss.NewStyle().Foreground(secondary),
		Dim:   lipgloss.NewStyle().Foreground(secondary).Faint(true),
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
	if !ok {
		return
	}

	title := i.Title()
	desc := i.Description()

	if index == m.Index() {
		t := d.styles.Title.Foreground(d.active).Render("󰄾 " + title)
		dsc := d.styles.Desc.Render(desc)
		line := lipgloss.JoinHorizontal(lipgloss.Center, t, " ", dsc)
		fmt.Fprint(w, d.styles.SelectionBar.Render(line))
	} else {
		t := d.styles.Title.PaddingLeft(2).Render(title)
		dsc := d.styles.Dim.Render(desc)
		line := lipgloss.JoinHorizontal(lipgloss.Center, t, " ", dsc)
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
	if len(m.config.Views) == 0 {
		return
	}
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

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate available space
		h, v := m.styles.App.GetFrameSize()
		// height - padding - tabs - footer - spacing
		listHeight := msg.Height - v - tabHeight - footerHeight - 2
		
		m.list.SetSize(msg.Width-h, listHeight)
		m.viewport = viewport.New(msg.Width-h, msg.Height-v-4)
		m.viewport.YPosition = 4
		m.ready = true

	case tea.KeyMsg:
		// Global Navigation
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

		// Viewport Mode
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

		// List Mode
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
	return m, cmd
}

func runCommand(command string) tea.Cmd {
	shell := os.Getenv("SHELL")
	if shell == "" { shell = "zsh" }
	c := exec.Command(shell, "-c", command+"; echo ''; echo 'Press Enter to return...'; read")
	return tea.ExecProcess(c, func(err error) tea.Msg { return nil })
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	if m.showViewport {
		i := m.list.SelectedItem().(Item)
		header := m.styles.Header.Render(" " + i.Title() + " ")
		help := m.styles.Dim.Render(" %3.f%% • q/esc: back • j/k: scroll")
		footer := fmt.Sprintf(help, m.viewport.ScrollPercent()*100)
		return m.styles.App.Render(lipgloss.JoinVertical(lipgloss.Left, header, "\n", m.viewport.View(), "\n", footer))
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
	
	// Add a full-width bottom border to the tab bar for a professional look
	width, _ := m.styles.App.GetFrameSize()
	fullWidth := m.width - width
	
	// Pad the tab row to fill width and add border
	tabBar := m.styles.TabBorder.Width(fullWidth).Render(tabRow)

	// 2. Help/Footer
	help := m.styles.Footer.Render("enter: run/view • x: exec • tab: switch • /: filter • q: quit")
	
	// 3. Assemble
	content := lipgloss.JoinVertical(lipgloss.Left, tabBar, "\n", m.list.View(), help)
	
	return m.styles.App.Render(content)
}
