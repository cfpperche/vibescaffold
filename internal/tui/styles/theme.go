package styles

import "github.com/charmbracelet/lipgloss"

var (
	Green   = lipgloss.Color("#39ff14")
	Black   = lipgloss.Color("#060608")
	Surface = lipgloss.Color("#0a0a0a")
	Border  = lipgloss.Color("#1a1a1a")
	Muted   = lipgloss.Color("#555555")
	White   = lipgloss.Color("#e0e0e0")
	Red     = lipgloss.Color("#ff4444")
	Yellow  = lipgloss.Color("#fbbf24")

	Logo = lipgloss.NewStyle().
		Foreground(Green).
		Bold(true)

	Title = lipgloss.NewStyle().
		Foreground(White).
		Bold(true)

	Subtle = lipgloss.NewStyle().
		Foreground(Muted)

	Success = lipgloss.NewStyle().
		Foreground(Green)

	Error = lipgloss.NewStyle().
		Foreground(Red)

	Warning = lipgloss.NewStyle().
		Foreground(Yellow)

	Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Border).
		Padding(1, 2)

	ActiveBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Green).
		Padding(1, 2)

	Footer = lipgloss.NewStyle().
		Foreground(Muted).
		MarginTop(1)
)

const LogoASCII = ` ██╗   ██╗███████╗
 ██║   ██║██╔════╝
 ██║   ██║███████╗
 ╚██╗ ██╔╝╚════██║
  ╚████╔╝ ███████║
   ╚═══╝  ╚══════╝`
