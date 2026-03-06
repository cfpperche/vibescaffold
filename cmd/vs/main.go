package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cfpperche/vibescaffold/internal/config"
	"github.com/cfpperche/vibescaffold/internal/tui"
)

func main() {
	var m tea.Model

	// If in a VS project (has CLAUDE.md), go straight to chat
	if config.DetectProject() {
		m = tui.NewWithChat()
	} else {
		m = tui.New()
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
