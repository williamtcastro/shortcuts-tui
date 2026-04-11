package tui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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

type Model struct {
	list         list.Model
	viewport     viewport.Model
	ready        bool
	width        int
	height       int
	renderer     *glamour.TermRenderer
	showViewport bool
}

func New(items []list.Item) Model {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Shortcuts Explorer"
	l.Styles.Title = titleStyle

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	return Model{
		list:     l,
		renderer: renderer,
	}
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
		header := titleStyle.Render(i.Title())
		footer := infoStyle(fmt.Sprintf("%3.f%% (q/esc to back, j/k to scroll)", m.viewport.ScrollPercent()*100))
		return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left, header, m.viewport.View(), footer))
	}

	return appStyle.Render(m.list.View())
}
