package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/cfpperche/vibescaffold/internal/chat"
	"github.com/cfpperche/vibescaffold/internal/tui/styles"
)

// ChatMsg types for async operations
type chatStreamMsg struct {
	token chat.StreamToken
}

type chatErrorMsg struct {
	err error
}

type chatEntry struct {
	role    string
	content string
	time    time.Time
}

type ChatModel struct {
	width   int
	height  int
	session *chat.Session

	viewport    viewport.Model
	input       textinput.Model
	spinner     spinner.Model
	entries     []chatEntry
	streaming   bool
	streamBuf   string
	cmdHistory  []string
	cmdIdx      int
	confirmExit bool
}

func NewChat(session *chat.Session) ChatModel {
	ti := textinput.New()
	ti.Placeholder = "mensagem ou /help..."
	ti.Focus()
	ti.CharLimit = 4096
	ti.Width = 80

	vp := viewport.New(80, 20)
	vp.SetContent("")

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = styles.Success

	m := ChatModel{
		session:  session,
		input:    ti,
		viewport: vp,
		spinner:  sp,
	}

	// Add welcome message
	contextFiles := chat.ContextFiles(session.ProjectDir)
	var welcome strings.Builder
	welcome.WriteString(fmt.Sprintf("✓ Projeto: %s\n", session.ProjectName))
	welcome.WriteString(fmt.Sprintf("✓ Agente: %s — pronto\n", session.Agent.Name))
	if len(contextFiles) > 0 {
		welcome.WriteString("\nContexto carregado:\n")
		for _, f := range contextFiles {
			welcome.WriteString(fmt.Sprintf("  ● %s\n", f))
		}
	}
	welcome.WriteString("\nDigite sua mensagem ou /help para comandos.")

	m.entries = append(m.entries, chatEntry{
		role:    "system",
		content: welcome.String(),
		time:    time.Now(),
	})

	return m
}

func (m ChatModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m ChatModel) Update(msg tea.Msg) (ChatModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		return m, nil

	case tea.KeyMsg:
		if m.confirmExit {
			switch msg.String() {
			case "y", "Y":
				return m, tea.Quit
			default:
				m.confirmExit = false
				m.entries = append(m.entries, chatEntry{
					role: "system", content: "Saida cancelada.", time: time.Now(),
				})
				m.updateViewport()
				return m, nil
			}
		}

		switch msg.String() {
		case "ctrl+c":
			if m.streaming {
				// TODO: cancel streaming
				return m, nil
			}
			m.confirmExit = true
			m.entries = append(m.entries, chatEntry{
				role: "system", content: "Deseja sair? (y/n)", time: time.Now(),
			})
			m.updateViewport()
			return m, nil
		case "enter":
			return m.handleSubmit()
		case "up":
			if len(m.cmdHistory) > 0 && m.cmdIdx > 0 {
				m.cmdIdx--
				m.input.SetValue(m.cmdHistory[m.cmdIdx])
				return m, nil
			}
		case "down":
			if m.cmdIdx < len(m.cmdHistory)-1 {
				m.cmdIdx++
				m.input.SetValue(m.cmdHistory[m.cmdIdx])
			} else if m.cmdIdx == len(m.cmdHistory)-1 {
				m.cmdIdx = len(m.cmdHistory)
				m.input.SetValue("")
			}
			return m, nil
		case "pgup", "pgdown":
			var vpCmd tea.Cmd
			m.viewport, vpCmd = m.viewport.Update(msg)
			return m, vpCmd
		}

	case chatStreamMsg:
		if msg.token.Err != nil {
			m.streaming = false
			m.entries = append(m.entries, chatEntry{
				role: "system", content: fmt.Sprintf("Erro: %s", msg.token.Err), time: time.Now(),
			})
			m.updateViewport()
			return m, nil
		}
		if msg.token.Done {
			m.streaming = false
			if m.streamBuf != "" {
				m.session.AddMessage("agent", m.streamBuf)
				m.entries = append(m.entries, chatEntry{
					role: "agent", content: m.streamBuf, time: time.Now(),
				})
				m.streamBuf = ""
			}
			m.updateViewport()
			return m, nil
		}
		m.streamBuf += msg.token.Content
		m.updateViewport()
		// Continue reading from stream
		return m, waitForStream

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update text input
	var tiCmd tea.Cmd
	m.input, tiCmd = m.input.Update(msg)
	cmds = append(cmds, tiCmd)

	return m, tea.Batch(cmds...)
}

func (m ChatModel) handleSubmit() (ChatModel, tea.Cmd) {
	input := strings.TrimSpace(m.input.Value())
	if input == "" {
		return m, nil
	}

	m.input.SetValue("")
	m.cmdHistory = append(m.cmdHistory, input)
	m.cmdIdx = len(m.cmdHistory)

	if m.session.IsCommand(input) {
		result := chat.HandleCommand(m.session, input)

		if result.Quit {
			return m, tea.Quit
		}

		if result.Output == "__clear__" {
			m.entries = m.entries[:0]
			m.entries = append(m.entries, chatEntry{
				role: "system", content: "Historico limpo.", time: time.Now(),
			})
			m.updateViewport()
			return m, nil
		}

		m.entries = append(m.entries, chatEntry{
			role: "user", content: input, time: time.Now(),
		})
		m.entries = append(m.entries, chatEntry{
			role: "system", content: result.Output, time: time.Now(),
		})
		m.updateViewport()
		return m, nil
	}

	// Regular message — send to agent
	m.entries = append(m.entries, chatEntry{
		role: "user", content: input, time: time.Now(),
	})
	m.session.AddMessage("user", input)
	m.streaming = true
	m.streamBuf = ""
	m.updateViewport()

	// Launch agent and start reading stream
	session := m.session
	return m, func() tea.Msg {
		ch, err := chat.RunAgent(session, input)
		if err != nil {
			return chatStreamMsg{token: chat.StreamToken{Err: err}}
		}
		// Store channel globally for subsequent reads
		activeStream = ch
		return waitForStream()
	}
}

var activeStream <-chan chat.StreamToken

func waitForStream() tea.Msg {
	if activeStream == nil {
		return chatStreamMsg{token: chat.StreamToken{Done: true}}
	}
	token, ok := <-activeStream
	if !ok {
		activeStream = nil
		return chatStreamMsg{token: chat.StreamToken{Done: true}}
	}
	return chatStreamMsg{token: token}
}

func (m *ChatModel) updateLayout() {
	headerH := 2
	inputH := 3
	footerH := 2
	vpH := m.height - headerH - inputH - footerH
	if vpH < 5 {
		vpH = 5
	}
	vpW := m.width - 2
	if vpW < 40 {
		vpW = 40
	}
	m.viewport.Width = vpW
	m.viewport.Height = vpH
	m.input.Width = vpW - 4
	m.updateViewport()
}

func (m *ChatModel) updateViewport() {
	var b strings.Builder
	w := m.viewport.Width - 4
	if w < 20 {
		w = 20
	}

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(w),
	)

	for _, e := range m.entries {
		switch e.role {
		case "user":
			b.WriteString(styles.Success.Render("  voce: "))
			b.WriteString(styles.Title.Render(e.content))
			b.WriteString("\n\n")
		case "agent":
			b.WriteString(styles.Success.Render(fmt.Sprintf("  %s: ", m.session.Agent.Name)))
			if renderer != nil {
				rendered, err := renderer.Render(e.content)
				if err == nil {
					b.WriteString(rendered)
				} else {
					b.WriteString(e.content)
				}
			} else {
				b.WriteString(e.content)
			}
			b.WriteString("\n")
		case "system":
			b.WriteString(styles.Subtle.Render("  " + e.content))
			b.WriteString("\n\n")
		}
	}

	// Show streaming buffer
	if m.streaming && m.streamBuf != "" {
		b.WriteString(styles.Success.Render(fmt.Sprintf("  %s: ", m.session.Agent.Name)))
		if renderer != nil {
			rendered, err := renderer.Render(m.streamBuf)
			if err == nil {
				b.WriteString(rendered)
			} else {
				b.WriteString(m.streamBuf)
			}
		} else {
			b.WriteString(m.streamBuf)
		}
		b.WriteString(styles.Subtle.Render("▊"))
		b.WriteString("\n")
	} else if m.streaming {
		b.WriteString(fmt.Sprintf("  %s %s pensando...\n",
			m.spinner.View(),
			styles.Subtle.Render(m.session.Agent.Name),
		))
	}

	m.viewport.SetContent(b.String())
	m.viewport.GotoBottom()
}

func (m ChatModel) View() string {
	if m.width == 0 {
		return "Carregando..."
	}

	var b strings.Builder

	// Header bar
	agentName := m.session.Agent.Name
	if m.session.Agent.Key == "ollama" && m.session.AppCfg.OllamaModel != "" {
		agentName += " (" + m.session.AppCfg.OllamaModel + ")"
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Green).
		Bold(true).
		Padding(0, 1).
		Width(m.width)

	header := fmt.Sprintf(" VibeScaffold  ●  %s  ●  %s  ●  v0.1",
		m.session.ProjectName,
		agentName,
	)
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")

	// Viewport (history)
	b.WriteString(m.viewport.View())
	b.WriteString("\n")

	// Input
	separator := lipgloss.NewStyle().
		Foreground(styles.Border).
		Width(m.width).
		Render(strings.Repeat("─", m.width))
	b.WriteString(separator)
	b.WriteString("\n")

	inputPrefix := styles.Success.Render(" > ")
	b.WriteString(inputPrefix + m.input.View())
	b.WriteString("\n")

	// Footer
	footer := styles.Subtle.Render("  /help  /switch  /doctor  /status  /context  /exit")
	b.WriteString(footer)

	return b.String()
}

func (m *ChatModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.updateLayout()
}
