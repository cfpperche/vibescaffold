package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cfpperche/vibeforge/internal/config"
	"github.com/cfpperche/vibeforge/internal/tui"
)

func main() {
	// Skip terminal background color query — prevents escape sequence
	// leaking into text inputs on Windows Terminal + WSL.
	lipgloss.SetHasDarkBackground(true)
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
