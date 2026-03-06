package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cfpperche/vibescaffold/internal/config"
	"github.com/cfpperche/vibescaffold/internal/doctor"
	"github.com/cfpperche/vibescaffold/internal/tui/components"
	"github.com/cfpperche/vibescaffold/internal/tui/styles"
)

type DoctorModel struct {
	width  int
	height int
	checks []doctor.Check
}

func NewDoctor() DoctorModel {
	return DoctorModel{
		checks: doctor.Run(),
	}
}

func (m DoctorModel) Init() tea.Cmd {
	return nil
}

func (m DoctorModel) Update(msg tea.Msg) (DoctorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return NavigateMsg{Target: "home"} }
		case "r":
			m.checks = doctor.Run()
			return m, nil
		}
	}
	return m, nil
}

func (m DoctorModel) View() string {
	var b strings.Builder

	b.WriteString(components.Header())
	b.WriteString(styles.Title.Render(fmt.Sprintf("  $ doctor — %s", config.ProjectName())))
	b.WriteString("\n\n")

	// Checks
	var lines []string
	for _, c := range m.checks {
		var icon string
		var style func(strs ...string) string
		switch c.Status {
		case "ok":
			icon = "✓"
			style = styles.Success.Render
		case "warn":
			icon = "⚠"
			style = styles.Warning.Render
		case "fail":
			icon = "✗"
			style = styles.Error.Render
		}
		line := fmt.Sprintf("  %s %-24s %s",
			style(icon),
			styles.Subtle.Render(c.Name),
			style(c.Detail),
		)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	b.WriteString(styles.Box.Width(50).Render(content))
	b.WriteString("\n\n")

	// Score
	ok, total := doctor.Score(m.checks)
	pct := 0
	if total > 0 {
		pct = ok * 100 / total
	}
	barFull := ok * 10 / total
	barEmpty := 10 - barFull
	bar := strings.Repeat("█", barFull) + strings.Repeat("░", barEmpty)

	b.WriteString(fmt.Sprintf("  Score: %d/%d %s  %d%%\n",
		ok, total,
		styles.Success.Render(bar),
		pct,
	))

	b.WriteString(components.Footer("  [r] refresh  [q] voltar"))

	return b.String()
}

func (m *DoctorModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
