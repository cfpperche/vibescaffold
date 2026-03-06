package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cfpperche/vibeforge/internal/onboarding"
	"github.com/cfpperche/vibeforge/internal/tui/components"
	"github.com/cfpperche/vibeforge/internal/tui/styles"
)

const onboardingListHeight = 20 // visible rows in the file list

// listEntry represents either a category header or a file in the flat list.
type listEntry struct {
	isCategory bool
	catIdx     int    // index into Categories
	filePath   string // path into FileMap
}

type OnboardingModel struct {
	width     int
	height    int
	cursor    int
	scroll    int // scroll offset for the file list
	entries   []listEntry
	expanded  map[int]bool // which categories are expanded (by catIdx)
}

func NewOnboarding() OnboardingModel {
	expanded := make(map[int]bool)
	for i := range onboarding.Categories {
		expanded[i] = true // all expanded by default
	}

	m := OnboardingModel{expanded: expanded}
	m.rebuildEntries()
	return m
}

func (m *OnboardingModel) rebuildEntries() {
	m.entries = nil
	for i, cat := range onboarding.Categories {
		m.entries = append(m.entries, listEntry{isCategory: true, catIdx: i})
		if m.expanded[i] {
			for _, path := range cat.Files {
				m.entries = append(m.entries, listEntry{filePath: path})
			}
		}
	}
	if m.cursor >= len(m.entries) {
		m.cursor = len(m.entries) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
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
				m.ensureVisible()
			}
		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
				m.ensureVisible()
			}
		case "enter", " ":
			if m.cursor < len(m.entries) {
				e := m.entries[m.cursor]
				if e.isCategory {
					m.expanded[e.catIdx] = !m.expanded[e.catIdx]
					m.rebuildEntries()
					m.ensureVisible()
				}
			}
		case "esc", "q":
			return m, func() tea.Msg { return NavigateMsg{Target: "home"} }
		}
	}
	return m, nil
}

func (m OnboardingModel) View() string {
	var b strings.Builder

	b.WriteString(components.Header())
	b.WriteString(styles.Title.Render("  $ help"))
	b.WriteString(styles.Subtle.Render("  — arquivos e estrutura\n\n"))
	b.WriteString(styles.Subtle.Render("  Cada arquivo tem um papel. O scaffold cria, o agente completa.\n\n"))

	// Split: left = file list, right = detail panel
	leftW := 40
	rightW := 46

	// Build left panel (file list with categories)
	var listLines []string
	for i, e := range m.entries {
		isCursor := i == m.cursor

		if e.isCategory {
			cat := onboarding.Categories[e.catIdx]
			arrow := "▸"
			if m.expanded[e.catIdx] {
				arrow = "▾"
			}
			prefix := "  "
			catStyle := styles.Subtle
			if isCursor {
				prefix = styles.Success.Render("> ")
				catStyle = styles.Title
			}
			count := len(cat.Files)
			line := fmt.Sprintf("%s%s %s %s",
				prefix,
				catStyle.Render(arrow),
				catStyle.Bold(true).Render(cat.Name),
				styles.Subtle.Render(fmt.Sprintf("(%d)", count)),
			)
			listLines = append(listLines, line)
		} else {
			f, ok := onboarding.FileMap[e.filePath]
			if !ok {
				continue
			}
			prefix := "    "
			nameStyle := styles.Subtle
			if isCursor {
				prefix = styles.Success.Render(">") + "   "
				nameStyle = styles.Title
			}
			fill := fillIndicator(f.FillLevel)
			short := f.Path
			if len(short) > 22 {
				short = "…" + short[len(short)-21:]
			}
			line := fmt.Sprintf("%s%s %s",
				prefix,
				nameStyle.Render(fmt.Sprintf("%-22s", short)),
				fill,
			)
			listLines = append(listLines, line)
		}
	}

	// Apply scroll window
	h := m.listHeight()
	end := m.scroll + h
	if end > len(listLines) {
		end = len(listLines)
	}
	visibleLines := listLines[m.scroll:end]

	// Scroll indicator
	if m.scroll > 0 {
		visibleLines = append([]string{styles.Subtle.Render("  ↑ mais...")}, visibleLines...)
	}
	if end < len(listLines) {
		visibleLines = append(visibleLines, styles.Subtle.Render("  ↓ mais..."))
	}

	listContent := strings.Join(visibleLines, "\n")
	listBox := styles.Box.Width(leftW).Render(listContent)

	// Build right panel (detail for selected)
	var detailBox string
	if m.cursor >= 0 && m.cursor < len(m.entries) {
		e := m.entries[m.cursor]
		if e.isCategory {
			cat := onboarding.Categories[e.catIdx]
			detailBox = m.renderCategoryDetail(cat, rightW)
		} else if f, ok := onboarding.FileMap[e.filePath]; ok {
			detailBox = m.renderFileDetail(f, rightW)
		}
	}

	// Render side by side
	leftLines := strings.Split(listBox, "\n")
	rightLines := strings.Split(detailBox, "\n")
	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	for i := 0; i < maxLines; i++ {
		left := ""
		right := ""
		if i < len(leftLines) {
			left = leftLines[i]
		}
		if i < len(rightLines) {
			right = rightLines[i]
		}
		// Pad left to fixed width using visible width
		leftVisible := lipgloss.Width(left)
		padding := leftW + 4 - leftVisible
		if padding < 1 {
			padding = 1
		}
		b.WriteString(left)
		b.WriteString(strings.Repeat(" ", padding))
		b.WriteString(right)
		b.WriteString("\n")
	}

	// Legend
	b.WriteString("\n")
	b.WriteString(styles.Subtle.Render("  "))
	b.WriteString(styles.Success.Render("●●●●"))
	b.WriteString(styles.Subtle.Render(" scaffold  "))
	b.WriteString(styles.Warning.Render("●●○○"))
	b.WriteString(styles.Subtle.Render(" parcial  "))
	b.WriteString(agentColor.Render("○○○○"))
	b.WriteString(styles.Subtle.Render(" agente\n"))

	b.WriteString(components.Footer("  [↑↓] navegar  [enter] expandir/colapsar  [q] voltar"))

	return b.String()
}

var agentColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#60a5fa"))

func (m OnboardingModel) renderCategoryDetail(cat onboarding.Category, w int) string {
	var lines []string
	lines = append(lines, styles.Title.Render(cat.Name))
	lines = append(lines, "")

	// Count fill levels
	scaffold, partial, agent := 0, 0, 0
	for _, path := range cat.Files {
		if f, ok := onboarding.FileMap[path]; ok {
			switch f.FillLevel {
			case onboarding.FilledByScaffold:
				scaffold++
			case onboarding.FilledPartial:
				partial++
			case onboarding.FilledByAgent:
				agent++
			}
		}
	}

	lines = append(lines, fmt.Sprintf("%s %d arquivos",
		styles.Subtle.Render("Total:"),
		len(cat.Files),
	))
	if scaffold > 0 {
		lines = append(lines, fmt.Sprintf("  %s %d preenchidos pelo scaffold",
			styles.Success.Render("●●●●"),
			scaffold,
		))
	}
	if partial > 0 {
		lines = append(lines, fmt.Sprintf("  %s %d parciais",
			styles.Warning.Render("●●○○"),
			partial,
		))
	}
	if agent > 0 {
		lines = append(lines, fmt.Sprintf("  %s %d preenchidos pelo agente",
			agentColor.Render("○○○○"),
			agent,
		))
	}

	lines = append(lines, "")
	lines = append(lines, styles.Subtle.Render("Arquivos:"))
	for _, path := range cat.Files {
		lines = append(lines, styles.Subtle.Render("  "+path))
	}

	content := strings.Join(lines, "\n")
	return styles.ActiveBox.Width(w).Render(content)
}

func (m OnboardingModel) renderFileDetail(f onboarding.MDFile, w int) string {
	var lines []string

	lines = append(lines, fmt.Sprintf("%s  %s",
		styles.Title.Render(f.Path),
		fillIndicator(f.FillLevel),
	))
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
		lines = append(lines, agentColor.Render("Preenchido pelo agente LLM:"))
		for _, s := range f.AgentFills {
			lines = append(lines, fmt.Sprintf("  %s %s",
				agentColor.Render("◈"),
				styles.Subtle.Render(s),
			))
		}
	}

	content := strings.Join(lines, "\n")
	return styles.ActiveBox.Width(w).Render(content)
}

func fillIndicator(level onboarding.FillLevel) string {
	switch level {
	case onboarding.FilledByScaffold:
		return styles.Success.Render("●●●●")
	case onboarding.FilledPartial:
		return styles.Warning.Render("●●○○")
	case onboarding.FilledByAgent:
		return agentColor.Render("○○○○")
	default:
		return "    "
	}
}

func (m *OnboardingModel) listHeight() int {
	// Use terminal height minus header/footer/legend/title overhead (~10 lines)
	h := m.height - 10
	if h < onboardingListHeight {
		h = onboardingListHeight
	}
	return h
}

func (m *OnboardingModel) ensureVisible() {
	h := m.listHeight()
	if m.cursor < m.scroll {
		m.scroll = m.cursor
	}
	if m.cursor >= m.scroll+h {
		m.scroll = m.cursor - h + 1
	}
}

func (m *OnboardingModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.ensureVisible()
}
