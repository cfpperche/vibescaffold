package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cfpperche/vibescaffold/internal/onboarding"
	"github.com/cfpperche/vibescaffold/internal/tui/components"
	"github.com/cfpperche/vibescaffold/internal/tui/styles"
)

type OnboardingModel struct {
	width  int
	height int
	cursor int
	files  []onboarding.MDFile
}

func NewOnboarding() OnboardingModel {
	return OnboardingModel{
		files: onboarding.Files,
	}
}

func (m OnboardingModel) Init() tea.Cmd {
	return nil
}

func (m OnboardingModel) Update(msg tea.Msg) (OnboardingModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}
		case "esc", "q":
			onboarding.MarkSeen()
			return m, func() tea.Msg { return NavigateMsg{Target: "home"} }
		}
	}
	return m, nil
}

func (m OnboardingModel) View() string {
	var b strings.Builder

	b.WriteString(components.Header())
	b.WriteString(styles.Title.Render("  $ onboarding"))
	b.WriteString(styles.Subtle.Render("  — como funciona\n\n"))

	b.WriteString(styles.Subtle.Render("  O scaffold cria a estrutura. O agente LLM preenche.\n\n"))

	// File list
	var listLines []string
	for i, f := range m.files {
		cursor := "  "
		nameStyle := styles.Subtle
		if i == m.cursor {
			cursor = styles.Success.Render("> ")
			nameStyle = styles.Title
		}

		fill := fillIndicator(f.FillLevel)
		label := fillLabel(f.FillLevel)

		line := fmt.Sprintf("%s%-28s %s %s",
			cursor,
			nameStyle.Render(f.Path),
			fill,
			styles.Subtle.Render(label),
		)
		listLines = append(listLines, line)
	}

	listContent := strings.Join(listLines, "\n")
	listBox := styles.Box.Width(62).Render(listContent)
	b.WriteString(listBox)
	b.WriteString("\n\n")

	// Detail panel for selected file
	if m.cursor >= 0 && m.cursor < len(m.files) {
		selected := m.files[m.cursor]
		b.WriteString(m.renderDetail(selected))
	}

	// Legend
	b.WriteString("\n")
	b.WriteString(styles.Subtle.Render("  "))
	b.WriteString(styles.Success.Render("●●●●"))
	b.WriteString(styles.Subtle.Render(" scaffold  "))
	b.WriteString(styles.Warning.Render("●●○○"))
	b.WriteString(styles.Subtle.Render(" parcial  "))
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#60a5fa")).Render("○○○○"))
	b.WriteString(styles.Subtle.Render(" agente\n"))

	b.WriteString(components.Footer("  [↑↓] navegar  [q] voltar"))

	return b.String()
}

func (m OnboardingModel) renderDetail(f onboarding.MDFile) string {
	var lines []string

	lines = append(lines, styles.Title.Render(f.Path))
	lines = append(lines, styles.Subtle.Render(f.Description))
	lines = append(lines, "")

	if len(f.ScaffoldFills) > 0 {
		lines = append(lines, styles.Success.Render("Criado pelo scaffold:"))
		for _, s := range f.ScaffoldFills {
			lines = append(lines, fmt.Sprintf("  %s %s",
				styles.Success.Render("✓"),
				styles.Subtle.Render(s),
			))
		}
		lines = append(lines, "")
	}

	if len(f.AgentFills) > 0 {
		agentColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#60a5fa"))
		lines = append(lines, agentColor.Render("Preenchido pelo agente LLM:"))
		for _, s := range f.AgentFills {
			lines = append(lines, fmt.Sprintf("  %s %s",
				agentColor.Render("◈"),
				styles.Subtle.Render(s),
			))
		}
	}

	content := strings.Join(lines, "\n")
	detailBox := styles.ActiveBox.Width(62).Render(content)
	return detailBox
}

func fillIndicator(level onboarding.FillLevel) string {
	switch level {
	case onboarding.FilledByScaffold:
		return styles.Success.Render("●●●●")
	case onboarding.FilledPartial:
		return styles.Warning.Render("●●○○")
	case onboarding.FilledByAgent:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#60a5fa")).Render("○○○○")
	default:
		return "    "
	}
}

func fillLabel(level onboarding.FillLevel) string {
	switch level {
	case onboarding.FilledByScaffold:
		return "preenchido"
	case onboarding.FilledPartial:
		return "parcial"
	case onboarding.FilledByAgent:
		return "preenchido pelo agente"
	default:
		return ""
	}
}

func (m *OnboardingModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
