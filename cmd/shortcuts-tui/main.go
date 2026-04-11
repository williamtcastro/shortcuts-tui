package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/williamtcastro/shortcuts-tui/internal/config"
	"github.com/williamtcastro/shortcuts-tui/internal/parser"
	"github.com/williamtcastro/shortcuts-tui/internal/tui"
)

func main() {
	cfg := config.LoadConfig()
	items := parser.LoadItems(cfg)
	m := tui.New(items, cfg)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
