package tui_test

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/cfpperche/vibeforge/internal/tui"
)

func TestHomeViewRenders(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		tui.New(),
		teatest.WithInitialTermSize(80, 24),
	)

	// Wait for the home view to render with the logo
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return strings.Contains(string(bts), "vibeforge")
		},
		teatest.WithDuration(3*time.Second),
		teatest.WithCheckInterval(100*time.Millisecond),
	)

	// Quit
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

func TestNavigateToDoctor(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		tui.New(),
		teatest.WithInitialTermSize(80, 24),
	)

	// Wait for home to render
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return strings.Contains(string(bts), "init")
		},
		teatest.WithDuration(3*time.Second),
		teatest.WithCheckInterval(100*time.Millisecond),
	)

	// Press 3 to go to doctor
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})

	// Wait for doctor view
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return strings.Contains(string(bts), "doctor")
		},
		teatest.WithDuration(3*time.Second),
		teatest.WithCheckInterval(100*time.Millisecond),
	)

	// Go back and quit
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			// Back at home, or doctor still showing - either way we can quit
			return strings.Contains(string(bts), "selecionar") || strings.Contains(string(bts), "voltar")
		},
		teatest.WithDuration(3*time.Second),
		teatest.WithCheckInterval(100*time.Millisecond),
	)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

func TestNavigateToStatus(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		tui.New(),
		teatest.WithInitialTermSize(80, 24),
	)

	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return strings.Contains(string(bts), "vibeforge")
		},
		teatest.WithDuration(3*time.Second),
		teatest.WithCheckInterval(100*time.Millisecond),
	)

	// Press 4 for status
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})

	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return strings.Contains(string(bts), "status")
		},
		teatest.WithDuration(3*time.Second),
		teatest.WithCheckInterval(100*time.Millisecond),
	)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return strings.Contains(string(bts), "selecionar") || strings.Contains(string(bts), "vibeforge")
		},
		teatest.WithDuration(3*time.Second),
		teatest.WithCheckInterval(100*time.Millisecond),
	)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

func TestArrowNavigation(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		tui.New(),
		teatest.WithInitialTermSize(80, 24),
	)

	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return strings.Contains(string(bts), "init")
		},
		teatest.WithDuration(3*time.Second),
		teatest.WithCheckInterval(100*time.Millisecond),
	)

	// Navigate down twice to doctor (index 2)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})

	// Press enter to go to doctor
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return strings.Contains(string(bts), "doctor")
		},
		teatest.WithDuration(3*time.Second),
		teatest.WithCheckInterval(100*time.Millisecond),
	)

	// Quit back
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return strings.Contains(string(bts), "selecionar") || strings.Contains(string(bts), "vibeforge")
		},
		teatest.WithDuration(3*time.Second),
		teatest.WithCheckInterval(100*time.Millisecond),
	)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

func TestCtrlCQuits(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		tui.New(),
		teatest.WithInitialTermSize(80, 24),
	)

	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return strings.Contains(string(bts), "vibeforge")
		},
		teatest.WithDuration(3*time.Second),
		teatest.WithCheckInterval(100*time.Millisecond),
	)

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}
