package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cfpperche/vibescaffold/internal/agent"
	"github.com/cfpperche/vibescaffold/internal/config"
	"github.com/cfpperche/vibescaffold/internal/tui/components"
	"github.com/cfpperche/vibescaffold/internal/tui/styles"
)

type agentSubView int

const (
	agentList agentSubView = iota
	agentInstallHint
	agentOllamaModels
)

// LaunchAgentMsg signals the app to suspend TUI and launch agent.
type LaunchAgentMsg struct {
	AgentKey    string
	OllamaModel string
}

type AgentModel struct {
	width       int
	height      int
	subView     agentSubView
	cursor      int
	agents      []agent.DetectedAgent
	appCfg      config.AppConfig
	ollamaCursor int
}

func NewAgent() AgentModel {
	appCfg := config.LoadAppConfig()
	agents := agent.DetectAll()

	// Set cursor to active agent
	cursor := 0
	for i, a := range agents {
		if a.Key == appCfg.ActiveAgent {
			cursor = i
			break
		}
	}

	return AgentModel{
		agents:  agents,
		appCfg:  appCfg,
		cursor:  cursor,
		subView: agentList,
	}
}

func (m AgentModel) Init() tea.Cmd {
	return nil
}

func (m AgentModel) Update(msg tea.Msg) (AgentModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.subView {
		case agentList:
			return m.updateList(msg)
		case agentInstallHint:
			return m.updateInstallHint(msg)
		case agentOllamaModels:
			return m.updateOllamaModels(msg)
		}
	}
	return m, nil
}

func (m AgentModel) updateList(msg tea.KeyMsg) (AgentModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.agents)-1 {
			m.cursor++
		}
	case "enter":
		selected := m.agents[m.cursor]
		if !selected.Installed {
			m.subView = agentInstallHint
			return m, nil
		}
		// If Ollama with models, show model selector
		if selected.Key == "ollama" && len(selected.Models) > 0 {
			m.subView = agentOllamaModels
			m.ollamaCursor = 0
			return m, nil
		}
		// Set as active and launch
		m.appCfg.ActiveAgent = selected.Key
		config.SaveAppConfig(m.appCfg)
		return m, func() tea.Msg {
			return LaunchAgentMsg{AgentKey: selected.Key}
		}
	case "i":
		selected := m.agents[m.cursor]
		if !selected.Installed {
			m.subView = agentInstallHint
		}
	case "r":
		m.agents = agent.DetectAll()
	case "c":
		// Set as active without launching
		selected := m.agents[m.cursor]
		if selected.Installed {
			m.appCfg.ActiveAgent = selected.Key
			config.SaveAppConfig(m.appCfg)
		}
	case "esc", "q":
		return m, func() tea.Msg { return NavigateMsg{Target: "home"} }
	}
	return m, nil
}

func (m AgentModel) updateInstallHint(msg tea.KeyMsg) (AgentModel, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.subView = agentList
	case "r":
		m.agents = agent.DetectAll()
		m.subView = agentList
	}
	return m, nil
}

func (m AgentModel) updateOllamaModels(msg tea.KeyMsg) (AgentModel, tea.Cmd) {
	selected := m.agents[m.cursor]
	switch msg.String() {
	case "up", "k":
		if m.ollamaCursor > 0 {
			m.ollamaCursor--
		}
	case "down", "j":
		if m.ollamaCursor < len(selected.Models)-1 {
			m.ollamaCursor++
		}
	case "enter":
		model := selected.Models[m.ollamaCursor].Name
		m.appCfg.ActiveAgent = "ollama"
		m.appCfg.OllamaModel = model
		config.SaveAppConfig(m.appCfg)
		return m, func() tea.Msg {
			return LaunchAgentMsg{AgentKey: "ollama", OllamaModel: model}
		}
	case "esc", "q":
		m.subView = agentList
	}
	return m, nil
}

func (m AgentModel) View() string {
	switch m.subView {
	case agentInstallHint:
		return m.viewInstallHint()
	case agentOllamaModels:
		return m.viewOllamaModels()
	default:
		return m.viewList()
	}
}

func (m AgentModel) viewList() string {
	var b strings.Builder

	b.WriteString(components.Header())
	b.WriteString(styles.Title.Render("  $ agent"))
	b.WriteString(styles.Subtle.Render("  — selecionar agente LLM\n\n"))

	var lines []string
	for i, a := range m.agents {
		cursor := "  "
		bullet := "○"
		if a.Key == m.appCfg.ActiveAgent {
			bullet = "●"
		}
		nameStyle := styles.Subtle
		if i == m.cursor {
			cursor = styles.Success.Render("> ")
			nameStyle = styles.Title
		}

		var statusStr string
		if a.Installed {
			ver := a.Version
			if ver == "" {
				ver = "instalado"
			}
			statusStr = styles.Success.Render("✓ ") + styles.Subtle.Render(ver)
		} else {
			statusStr = styles.Error.Render("✗ nao encontrado")
		}

		bulletStyle := styles.Subtle
		if a.Key == m.appCfg.ActiveAgent {
			bulletStyle = styles.Success
		}

		line := fmt.Sprintf("%s%s %-16s %s",
			cursor,
			bulletStyle.Render(bullet),
			nameStyle.Render(a.Name),
			statusStr,
		)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	b.WriteString(styles.Box.Width(52).Render(content))
	b.WriteString("\n\n")

	// Active agent status
	b.WriteString(styles.Subtle.Render("  Agente ativo: "))
	b.WriteString(styles.Success.Render(m.activeAgentName()))
	b.WriteString("\n")

	b.WriteString(components.Footer("  [↑↓] selecionar  [enter] usar  [i] instalar  [c] ativar  [r] re-detectar  [q] voltar"))

	return b.String()
}

func (m AgentModel) viewInstallHint() string {
	var b strings.Builder
	selected := m.agents[m.cursor]
	hint := agent.InstallHint(selected.Key)

	b.WriteString(components.Header())
	b.WriteString(styles.Title.Render(fmt.Sprintf("  $ instalar %s", selected.Name)))
	b.WriteString("\n\n")

	var lines []string
	lines = append(lines, "")
	lines = append(lines, styles.Success.Render("  "+hint))
	lines = append(lines, "")
	lines = append(lines, styles.Subtle.Render("  Apos instalar, pressione [r] para detectar"))
	lines = append(lines, styles.Subtle.Render("  novamente."))
	lines = append(lines, "")

	b.WriteString(styles.ActiveBox.Width(52).Render(strings.Join(lines, "\n")))
	b.WriteString("\n")

	b.WriteString(components.Footer("  [r] re-detectar  [q] voltar"))

	return b.String()
}

func (m AgentModel) viewOllamaModels() string {
	var b strings.Builder
	selected := m.agents[m.cursor]

	b.WriteString(components.Header())
	b.WriteString(styles.Title.Render("  $ ollama — selecionar modelo"))
	b.WriteString("\n\n")

	if len(selected.Models) == 0 {
		b.WriteString(styles.Warning.Render("  ⚠ Nenhum modelo encontrado\n"))
		b.WriteString(styles.Subtle.Render("  Execute: ollama pull llama3.2\n"))
		b.WriteString(components.Footer("  [q] voltar"))
		return b.String()
	}

	var lines []string
	for i, model := range selected.Models {
		cursor := "  "
		bullet := "○"
		nameStyle := styles.Subtle
		if i == m.ollamaCursor {
			cursor = styles.Success.Render("> ")
			bullet = "●"
			nameStyle = styles.Title
		}
		if model.Name == m.appCfg.OllamaModel {
			bullet = "●"
		}

		sizeStr := formatSize(model.Size)
		line := fmt.Sprintf("%s%s %-24s %s",
			cursor,
			styles.Success.Render(bullet),
			nameStyle.Render(model.Name),
			styles.Subtle.Render(sizeStr),
		)
		lines = append(lines, line)
	}

	b.WriteString(styles.Box.Width(52).Render(strings.Join(lines, "\n")))
	b.WriteString("\n")
	b.WriteString(components.Footer("  [↑↓] selecionar  [enter] confirmar  [q] voltar"))

	return b.String()
}

func (m AgentModel) activeAgentName() string {
	for _, a := range m.agents {
		if a.Key == m.appCfg.ActiveAgent {
			if a.Key == "ollama" && m.appCfg.OllamaModel != "" {
				return a.Name + " (" + m.appCfg.OllamaModel + ")"
			}
			return a.Name
		}
	}
	return m.appCfg.ActiveAgent
}

func formatSize(bytes int64) string {
	if bytes == 0 {
		return ""
	}
	gb := float64(bytes) / (1024 * 1024 * 1024)
	if gb >= 1 {
		return fmt.Sprintf("%.1fGB", gb)
	}
	mb := float64(bytes) / (1024 * 1024)
	return fmt.Sprintf("%.0fMB", mb)
}

func (m *AgentModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
