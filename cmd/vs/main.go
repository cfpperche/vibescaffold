package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cfpperche/vibescaffold/internal/tui"
)

func main() {
	m := tui.New()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
