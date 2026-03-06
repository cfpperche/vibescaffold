package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cfpperche/vibeforge/internal/config"
	"github.com/cfpperche/vibeforge/internal/i18n"
	"github.com/cfpperche/vibeforge/internal/tui/components"
	"github.com/cfpperche/vibeforge/internal/tui/styles"

	"github.com/cfpperche/vibeforge/internal/agent"
)

type menuItem struct {
	key  string
	name string
	desc func() string
}

var menuItems = []menuItem{
	{"1", "new", func() string { return i18n.T("home.menu.new.desc") }},
	{"2", "init", func() string { return i18n.T("home.menu.init.desc") }},
	{"3", "doctor", func() string { return i18n.T("home.menu.doctor.desc") }},
	{"4", "status", func() string { return i18n.T("home.menu.status.desc") }},
	{"5", "agent", func() string { return i18n.T("home.menu.agent.desc") }},
	{"6", "help", func() string { return i18n.T("home.menu.help.desc") }},
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

// EnterChatMsg signals the app to enter chat mode for a project.
type EnterChatMsg struct {
	ProjectDir  string
	ProjectName string
	Summary     string // scaffold summary
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
			return m, func() tea.Msg { return NavigateMsg{Target: "new"} }
		case "2":
			return m, func() tea.Msg { return NavigateMsg{Target: "init"} }
		case "3":
			return m, func() tea.Msg { return NavigateMsg{Target: "doctor"} }
		case "4":
			return m, func() tea.Msg { return NavigateMsg{Target: "status"} }
		case "5":
			return m, func() tea.Msg { return NavigateMsg{Target: "agent"} }
		case "6":
			return m, func() tea.Msg { return NavigateMsg{Target: "help"} }
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
			styles.Subtle.Render(item.desc()),
		)
		menuLines = append(menuLines, line)
	}

	menu := styles.Box.Width(52).Render(strings.Join(menuLines, "\n"))
	b.WriteString(menu)
	b.WriteString("\n\n")

	// Project detection
	if config.DetectProject() {
		b.WriteString(styles.Success.Render("  ✓ "))
		b.WriteString(styles.Subtle.Render(i18n.TF("home.project_detected", config.ProjectName())))
	} else {
		b.WriteString(styles.Warning.Render("  ⚠ "))
		b.WriteString(styles.Subtle.Render(i18n.T("home.no_project")))
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
	b.WriteString(styles.Subtle.Render("  " + i18n.T("home.agent_label")))
	b.WriteString(styles.Success.Render("● " + agentName))
	b.WriteString("\n")

	b.WriteString(components.Footer(i18n.T("home.footer")))

	return b.String()
}

func (m *HomeModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
