package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cfpperche/vibescaffold/internal/tui/views"
)

type view int

const (
	viewHome view = iota
	viewInit
	viewDoctor
	viewStatus
)

type Model struct {
	currentView view
	width       int
	height      int
	home        views.HomeModel
	init        views.InitModel
	doctor      views.DoctorModel
	status      views.StatusModel
}

func New() Model {
	return Model{
		currentView: viewHome,
		home:        views.NewHome(),
		init:        views.NewInit(),
		doctor:      views.NewDoctor(),
		status:      views.NewStatus(),
	}
}

func (m Model) Init() tea.Cmd {
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
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" && m.currentView == viewHome {
			return m, tea.Quit
		}
		if msg.String() == "ctrl+c" {
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
		case "context":
			// context is same as doctor for now
			m.currentView = viewDoctor
			return m, nil
		}
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
	default:
		return m.home.View()
	}
}
