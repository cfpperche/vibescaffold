package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cfpperche/vibescaffold/internal/config"
	"github.com/cfpperche/vibescaffold/internal/tui/components"
	"github.com/cfpperche/vibescaffold/internal/tui/styles"

	"github.com/cfpperche/vibescaffold/internal/agent"
)

type menuItem struct {
	key  string
	name string
	desc string
}

var menuItems = []menuItem{
	{"1", "init", "scaffold projeto"},
	{"2", "doctor", "health check"},
	{"3", "status", "roadmap"},
	{"4", "agent", "selecionar agente LLM"},
}

type HomeModel struct {
	cursor int
	width  int
	height int
}

func NewHome() HomeModel {
	return HomeModel{}
}

type NavigateMsg struct {
	Target string
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(menuItems)-1 {
				m.cursor++
			}
		case "enter":
			return m, func() tea.Msg {
				return NavigateMsg{Target: menuItems[m.cursor].name}
			}
		case "1":
			return m, func() tea.Msg { return NavigateMsg{Target: "init"} }
		case "2":
			return m, func() tea.Msg { return NavigateMsg{Target: "doctor"} }
		case "3":
			return m, func() tea.Msg { return NavigateMsg{Target: "status"} }
		case "4":
			return m, func() tea.Msg { return NavigateMsg{Target: "agent"} }
		}
	}
	return m, nil
}

func (m HomeModel) View() string {
	var b strings.Builder

	b.WriteString(components.Header())
	b.WriteString("\n")

	// Menu
	var menuLines []string
	for i, item := range menuItems {
		cursor := "  "
		style := styles.Subtle
		if i == m.cursor {
			cursor = styles.Success.Render("> ")
			style = styles.Success
		}
		line := fmt.Sprintf("%s%s  %s    %s",
			cursor,
			style.Render(item.key),
			style.Bold(true).Render(fmt.Sprintf("%-10s", item.name)),
			styles.Subtle.Render(item.desc),
		)
		menuLines = append(menuLines, line)
	}

	menu := styles.Box.Width(42).Render(strings.Join(menuLines, "\n"))
	b.WriteString(menu)
	b.WriteString("\n\n")

	// Project detection
	if config.DetectProject() {
		b.WriteString(styles.Success.Render("  ✓ "))
		b.WriteString(styles.Subtle.Render(fmt.Sprintf("Projeto detectado: %s", config.ProjectName())))
	} else {
		b.WriteString(styles.Warning.Render("  ⚠ "))
		b.WriteString(styles.Subtle.Render("Nenhum projeto detectado nesta pasta"))
	}
	b.WriteString("\n")

	// Active agent
	appCfg := config.LoadAppConfig()
	agentName := appCfg.ActiveAgent
	for _, a := range agent.DefaultAgents() {
		if a.Key == appCfg.ActiveAgent {
			agentName = a.Name
			break
		}
	}
	b.WriteString(styles.Subtle.Render("  Agente: "))
	b.WriteString(styles.Success.Render("● " + agentName))
	b.WriteString("\n")

	b.WriteString(components.Footer("  [1-4] selecionar  [↑↓] navegar  [enter] confirmar  [q] sair"))

	return b.String()
}

func (m *HomeModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
