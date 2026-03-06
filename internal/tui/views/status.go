package views

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cfpperche/vibeforge/internal/config"
	"github.com/cfpperche/vibeforge/internal/tui/components"
	"github.com/cfpperche/vibeforge/internal/tui/styles"
)

type phase struct {
	name    string
	total   int
	done    int
}

type StatusModel struct {
	width  int
	height int
	phases []phase
	found  bool
}

func NewStatus() StatusModel {
	m := StatusModel{}
	m.loadRoadmap()
	return m
}

func (m *StatusModel) loadRoadmap() {
	data, err := os.ReadFile("docs/ROADMAP.md")
	if err != nil {
		m.found = false
		return
	}
	m.found = true
	m.phases = parseRoadmap(string(data))
}

func parseRoadmap(content string) []phase {
	var phases []phase
	phaseRe := regexp.MustCompile(`(?m)^##\s+(.+)$`)
	taskDone := regexp.MustCompile(`(?m)^- \[x\]`)
	taskTodo := regexp.MustCompile(`(?m)^- \[ \]`)

	matches := phaseRe.FindAllStringIndex(content, -1)
	names := phaseRe.FindAllStringSubmatch(content, -1)

	for i, name := range names {
		start := matches[i][1]
		end := len(content)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		section := content[start:end]
		done := len(taskDone.FindAllString(section, -1))
		todo := len(taskTodo.FindAllString(section, -1))
		total := done + todo
		if total == 0 {
			total = 1 // avoid div by zero
		}
		phases = append(phases, phase{name: name[1], total: total, done: done})
	}
	return phases
}

func (m StatusModel) Init() tea.Cmd {
	return nil
}

func (m StatusModel) Update(msg tea.Msg) (StatusModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return NavigateMsg{Target: "home"} }
		case "r":
			m.loadRoadmap()
			return m, nil
		}
	}
	return m, nil
}

func (m StatusModel) View() string {
	var b strings.Builder

	b.WriteString(components.Header())
	b.WriteString(styles.Title.Render(fmt.Sprintf("  $ status — %s", config.ProjectName())))
	b.WriteString("\n\n")

	if !m.found {
		b.WriteString(styles.Warning.Render("  ⚠ docs/ROADMAP.md nao encontrado\n"))
		b.WriteString(styles.Subtle.Render("  Crie um ROADMAP.md com ## Fase e - [x]/- [ ] tarefas\n"))
		b.WriteString(components.Footer("  [q] voltar"))
		return b.String()
	}

	var lines []string
	totalDone := 0
	totalAll := 0
	for _, p := range m.phases {
		pct := p.done * 100 / p.total
		barFull := p.done * 12 / p.total
		barEmpty := 12 - barFull
		bar := strings.Repeat("█", barFull) + strings.Repeat("░", barEmpty)

		var barStyle func(strs ...string) string
		if pct == 100 {
			barStyle = styles.Success.Render
		} else if pct > 0 {
			barStyle = styles.Warning.Render
		} else {
			barStyle = styles.Subtle.Render
		}

		line := fmt.Sprintf("  %-28s %s %3d%%",
			styles.Subtle.Render(p.name),
			barStyle(bar),
			pct,
		)
		lines = append(lines, line)
		totalDone += p.done
		totalAll += p.total
	}

	content := strings.Join(lines, "\n")
	b.WriteString(styles.Box.Width(54).Render(content))
	b.WriteString("\n\n")

	// Overall
	if totalAll > 0 {
		overallPct := totalDone * 100 / totalAll
		overallFull := totalDone * 10 / totalAll
		overallEmpty := 10 - overallFull
		overallBar := strings.Repeat("█", overallFull) + strings.Repeat("░", overallEmpty)
		b.WriteString(fmt.Sprintf("  Progresso geral: %s  %d%%\n",
			styles.Success.Render(overallBar),
			overallPct,
		))
	}

	b.WriteString(styles.Subtle.Render(fmt.Sprintf("\n  Ultima atualizacao: %s\n", time.Now().Format("2006-01-02"))))
	b.WriteString(components.Footer("  [r] refresh  [q] voltar"))

	return b.String()
}

func (m *StatusModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
