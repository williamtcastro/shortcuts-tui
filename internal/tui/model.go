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
	SearchPrompt lipgloss.Style
	SearchCursor lipgloss.Style
	KeyStyle     lipgloss.Style
	StatusStyle  lipgloss.Style
}

func DefaultStyles(cfg config.Config) Styles {
	primary := lipgloss.Color(cfg.Theme.Primary)
	secondary := lipgloss.Color(cfg.Theme.Secondary)
	text := lipgloss.Color(cfg.Theme.Text)
	accent := lipgloss.Color(cfg.Theme.Accent)
	mauve := lipgloss.Color(cfg.Theme.Mauve)
	flamingo := lipgloss.Color(cfg.Theme.Flamingo)

	return Styles{
		App: lipgloss.NewStyle().Padding(1, 2),
		Header: lipgloss.NewStyle().
			Foreground(text).
			Background(primary).
			Padding(0, 1).
			Bold(true),
		Footer: lipgloss.NewStyle().
			MarginTop(1),
		ActiveTab: lipgloss.NewStyle().
			Foreground(primary).
			Background(lipgloss.Color("236")).
			Padding(0, 2).
			Bold(true),
		InactiveTab: lipgloss.NewStyle().
			Foreground(secondary).
			Padding(0, 2),
		TabSeparator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("238")).
			Border(lipgloss.NormalBorder(), false, false, true, false),
		SelectionBar: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(accent).
			PaddingLeft(1),
		Title: lipgloss.NewStyle().Width(25).Bold(true),
		Desc:  lipgloss.NewStyle().Foreground(mauve),
		Dim:   lipgloss.NewStyle().Foreground(secondary).Faint(true),
		DocViewport: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondary).
			Padding(0, 1),
		DocHeader: lipgloss.NewStyle().
			Foreground(primary).
			Bold(true).
			MarginBottom(1),
		SearchPrompt: lipgloss.NewStyle().Foreground(flamingo).Bold(true),
		SearchCursor: lipgloss.NewStyle().Foreground(text).Background(flamingo),
		KeyStyle:     lipgloss.NewStyle().Foreground(flamingo).Bold(true),
		StatusStyle:  lipgloss.NewStyle().Foreground(secondary).MarginLeft(1),
	}
}

// --- Delegate ---
type itemDelegate struct {
	styles Styles
	active lipgloss.Color
	accent lipgloss.Color
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
		t := d.styles.Title.Foreground(d.accent).Render("󰄾 " + titleStr)
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
		active: lipgloss.Color(cfg.Theme.Primary),
		accent: lipgloss.Color(cfg.Theme.Accent),
	}

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	
	l.Styles.FilterPrompt = s.SearchPrompt
	l.Styles.FilterCursor = s.SearchCursor

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
		listHeight := msg.Height - v - 5
		if listHeight < 0 { listHeight = 0 }
		
		m.list.SetSize(msg.Width-h, listHeight)
		m.viewport = viewport.New(msg.Width-h-4, msg.Height-v-8)
		m.ready = true

	case tea.KeyMsg:
		if !m.list.SettingFilter() && len(m.config.Views) > 0 {
			switch msg.String() {
			case "tab", "l", "right":
				m.activeTabIndex = (m.activeTabIndex + 1) % len(m.config.Views)
				m.updateListForTab()
				m.showViewport = false
				return m, nil
			case "shift+tab", "h", "left":
				m.activeTabIndex = (m.activeTabIndex - 1 + len(m.config.Views)) % len(m.config.Views)
				m.updateListForTab()
				m.showViewport = false
				return m, nil
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
			if !m.list.SettingFilter() { return m, tea.Quit }
		case "x":
			if !m.list.SettingFilter() {
				if i, ok := m.list.SelectedItem().(Item); ok && i.IsAlias {
					return m, runCommand(i.Command, m.config.AutoClear)
				}
			}
		case "enter":
			if i, ok := m.list.SelectedItem().(Item); ok {
				if i.IsAlias && !m.list.SettingFilter() {
					return m, runCommand(i.Command, m.config.AutoClear)
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

func runCommand(command string, autoClear bool) tea.Cmd {
	shell := os.Getenv("SHELL")
	if shell == "" { shell = "zsh" }
	
	fullCmd := command
	if autoClear {
		fullCmd = "clear && " + command
	}

	c := exec.Command(shell, "-c", fullCmd+"; echo ''; echo 'Press Enter to return...'; read")
	return tea.ExecProcess(c, func(err error) tea.Msg { return nil })
}

func (m Model) View() string {
	if !m.ready { return "\n  Initializing..." }

	if m.showViewport {
		i := m.list.SelectedItem().(Item)
		header := m.styles.DocHeader.Render("󰧮 " + i.Title())
		view := m.styles.DocViewport.Render(m.viewport.View())
		
		scroll := fmt.Sprintf(" %3.f%% ", m.viewport.ScrollPercent()*100)
		footer := lipgloss.JoinHorizontal(lipgloss.Center,
			m.styles.KeyStyle.Render(" esc "), m.styles.Dim.Render("back • "),
			m.styles.KeyStyle.Render(" j/k "), m.styles.Dim.Render("scroll • "),
			m.styles.StatusStyle.Render(scroll),
		)
		
		return m.styles.App.Render(lipgloss.JoinVertical(lipgloss.Left, header, view, m.styles.Footer.Render(footer)))
	}

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
	
	h, _ := m.styles.App.GetFrameSize()
	fullWidth := m.width - h
	if fullWidth < 0 { fullWidth = 0 }
	tabSeparator := m.styles.Dim.Render(strings.Repeat("─", fullWidth))

	keys := []string{
		m.styles.KeyStyle.Render(" enter "), m.styles.Dim.Render("run • "),
		m.styles.KeyStyle.Render(" x "), m.styles.Dim.Render("exec • "),
		m.styles.KeyStyle.Render(" tab "), m.styles.Dim.Render("switch • "),
		m.styles.KeyStyle.Render(" / "), m.styles.Dim.Render("filter"),
	}
	helpBar := lipgloss.JoinHorizontal(lipgloss.Center, keys...)
// 3. Pagination Info
var pagination string
if m.config.Pagination == "dots" {
	totalPages := m.list.Paginator.TotalPages
	currentPage := m.list.Paginator.Page
	for i := 0; i < totalPages; i++ {
		if i == currentPage {
			pagination += "●"
		} else {
			pagination += "•"
		}
	}
} else {
	pagination = fmt.Sprintf(" %d/%d ", m.list.Paginator.Page+1, m.list.Paginator.TotalPages)
}

status := m.styles.StatusStyle.Render(pagination)

// Create the full width footer with help on left and pagination on right
footerContent := lipgloss.JoinHorizontal(lipgloss.Top, helpBar, lipgloss.NewStyle().Width(fullWidth-lipgloss.Width(helpBar)).Align(lipgloss.Right).Render(status))


	content := lipgloss.JoinVertical(lipgloss.Left, tabRow, tabSeparator, m.list.View(), m.styles.Footer.Render(footerContent))
	
	return m.styles.App.Render(content)
}
