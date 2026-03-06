package tui

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cfpperche/vibescaffold/internal/chat"
	"github.com/cfpperche/vibescaffold/internal/tui/views"
)

type view int

const (
	viewHome view = iota
	viewInit
	viewDoctor
	viewStatus
	viewAgent
	viewChat
)

type Model struct {
	currentView view
	width       int
	height      int
	home        views.HomeModel
	init        views.InitModel
	doctor      views.DoctorModel
	status      views.StatusModel
	agent       views.AgentModel
	chat        views.ChatModel
}

func New() Model {
	return Model{
		currentView: viewHome,
		home:        views.NewHome(),
		init:        views.NewInit(),
		doctor:      views.NewDoctor(),
		status:      views.NewStatus(),
		agent:       views.NewAgent(),
	}
}

// NewWithChat creates the TUI and goes directly to chat mode.
func NewWithChat() Model {
	cwd, _ := os.Getwd()
	projectName := filepath.Base(cwd)
	session := chat.NewSession(cwd, projectName)
	chatModel := views.NewChat(session)

	return Model{
		currentView: viewChat,
		home:        views.NewHome(),
		chat:        chatModel,
	}
}

func (m Model) Init() tea.Cmd {
	if m.currentView == viewChat {
		return m.chat.Init()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.home.SetSize(msg.Width, msg.Height)
		m.init.SetSize(msg.Width, msg.Height)
		m.doctor.SetSize(msg.Width, msg.Height)
		m.status.SetSize(msg.Width, msg.Height)
		m.agent.SetSize(msg.Width, msg.Height)
		m.chat.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" && m.currentView == viewHome {
			return m, tea.Quit
		}
		if msg.String() == "ctrl+c" && m.currentView != viewChat {
			return m, tea.Quit
		}

	case views.NavigateMsg:
		switch msg.Target {
		case "home":
			m.currentView = viewHome
			return m, nil
		case "init":
			m.init = views.NewInit()
			m.init.SetSize(m.width, m.height)
			m.currentView = viewInit
			return m, m.init.Init()
		case "doctor":
			m.doctor = views.NewDoctor()
			m.doctor.SetSize(m.width, m.height)
			m.currentView = viewDoctor
			return m, nil
		case "status":
			m.status = views.NewStatus()
			m.status.SetSize(m.width, m.height)
			m.currentView = viewStatus
			return m, nil
		case "agent":
			m.agent = views.NewAgent()
			m.agent.SetSize(m.width, m.height)
			m.currentView = viewAgent
			return m, nil
		case "chat":
			return m, nil // chat is entered via EnterChatMsg
		}

	case views.EnterChatMsg:
		session := chat.NewSession(msg.ProjectDir, msg.ProjectName)
		session.AddMessage("system", msg.Summary)
		m.chat = views.NewChat(session)
		m.chat.SetSize(m.width, m.height)
		m.currentView = viewChat
		return m, m.chat.Init()

	case views.LaunchAgentMsg:
		// Suspend TUI, launch agent, then resume
		return m, tea.ExecProcess(
			launchAgentCmd(msg),
			func(err error) tea.Msg {
				return views.NavigateMsg{Target: "home"}
			},
		)
	}

	var cmd tea.Cmd
	switch m.currentView {
	case viewHome:
		m.home, cmd = m.home.Update(msg)
	case viewInit:
		m.init, cmd = m.init.Update(msg)
	case viewDoctor:
		m.doctor, cmd = m.doctor.Update(msg)
	case viewStatus:
		m.status, cmd = m.status.Update(msg)
	case viewAgent:
		m.agent, cmd = m.agent.Update(msg)
	case viewChat:
		m.chat, cmd = m.chat.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.currentView {
	case viewInit:
		return m.init.View()
	case viewDoctor:
		return m.doctor.View()
	case viewStatus:
		return m.status.View()
	case viewAgent:
		return m.agent.View()
	case viewChat:
		return m.chat.View()
	default:
		return m.home.View()
	}
}
