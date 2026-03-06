package views

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/cfpperche/vibeforge/internal/config"
	"github.com/cfpperche/vibeforge/internal/i18n"
	"github.com/cfpperche/vibeforge/internal/scaffold"
	"github.com/cfpperche/vibeforge/internal/tui/components"
	"github.com/cfpperche/vibeforge/internal/tui/styles"
)

type initPhase int

const (
	phaseForm initPhase = iota
	phaseGenerating
	phaseDone
)

type scaffoldDoneMsg struct {
	files      []string
	err        error
	projectDir string
}

type InitModel struct {
	width      int
	height     int
	phase      initPhase
	form       *huh.Form
	cfg        config.Config
	files      []string
	err        error
	projectDir string
}

func NewInit() InitModel {
	cfg := config.Config{
		Features: []string{"claude-md", "context-docs"},
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(i18n.T("init.field.name")).
				Placeholder(i18n.T("init.field.name.placeholder")).
				Value(&cfg.Name),
			huh.NewInput().
				Title(i18n.T("init.field.desc")).
				Placeholder(i18n.T("init.field.desc.placeholder")).
				Value(&cfg.Desc),
			huh.NewInput().
				Title(i18n.T("init.field.author")).
				Placeholder(i18n.T("init.field.author.placeholder")).
				Value(&cfg.Author),
			huh.NewConfirm().
				Title(i18n.T("init.field.repo")).
				Value(&cfg.Repo),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(i18n.T("init.field.type")).
				Options(
					huh.NewOption("SaaS Web", "saas"),
					huh.NewOption("API Backend", "api"),
					huh.NewOption("Game", "game"),
					huh.NewOption("CLI Tool", "cli"),
					huh.NewOption("Mobile", "mobile"),
					huh.NewOption("Library", "lib"),
				).
				Value(&cfg.Type),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(i18n.T("init.field.stack")).
				Options(
					huh.NewOption("Go + Bubble Tea", "go-bubbletea"),
					huh.NewOption("Go + Chi", "go-chi"),
					huh.NewOption("Go + Gin", "go-gin"),
					huh.NewOption("TypeScript + React + Vite", "ts-react-vite"),
					huh.NewOption("TypeScript + Next.js", "ts-nextjs"),
					huh.NewOption("TypeScript + Hono", "ts-hono"),
					huh.NewOption("TypeScript + Elysia (Bun)", "ts-elysia"),
					huh.NewOption("Python + FastAPI", "py-fastapi"),
					huh.NewOption("Python + Django", "py-django"),
					huh.NewOption("Rust + Axum", "rs-axum"),
				).
				Value(&cfg.Stack),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(i18n.T("init.field.principles")).
				Options(
					huh.NewOption("TDD", "tdd"),
					huh.NewOption("Clean Architecture", "clean-arch"),
					huh.NewOption("SOLID", "solid"),
					huh.NewOption("12-Factor", "12-factor"),
					huh.NewOption("DDD", "ddd"),
					huh.NewOption("CQRS", "cqrs"),
					huh.NewOption("Event Sourcing", "event-sourcing"),
					huh.NewOption("Hexagonal", "hexagonal"),
				).
				Value(&cfg.Principles),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(i18n.T("init.field.features")).
				Options(
					huh.NewOption("CLAUDE.md (regras)", "claude-md"),
					huh.NewOption("docs/CONTEXT.md", "context-docs"),
					huh.NewOption("docs/ROADMAP.md", "roadmap"),
					huh.NewOption("docs/SRS.md", "srs"),
					huh.NewOption("ADRs (Architecture Decision Records)", "adrs"),
					huh.NewOption(".claude/hooks", "hooks"),
					huh.NewOption(".claude/commands", "commands"),
					huh.NewOption("GitHub Actions CI", "ci"),
				).
				Value(&cfg.Features),
		),
	).WithTheme(huh.ThemeCharm())

	return InitModel{
		phase: phaseForm,
		form:  form,
		cfg:   cfg,
	}
}

func (m InitModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m InitModel) Update(msg tea.Msg) (InitModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && m.phase != phaseGenerating {
			if m.phase == phaseDone {
				return m, func() tea.Msg { return NavigateMsg{Target: "home"} }
			}
			return m, func() tea.Msg { return NavigateMsg{Target: "home"} }
		}
		if msg.String() == "enter" && m.phase == phaseDone && m.err == nil {
			// Transition to chat
			dir := m.projectDir
			name := m.cfg.Name
			summary := i18n.TF("init.summary", name, len(m.files))
			return m, func() tea.Msg {
				return EnterChatMsg{
					ProjectDir:  dir,
					ProjectName: name,
					Summary:     summary,
				}
			}
		}
	case scaffoldDoneMsg:
		m.phase = phaseDone
		m.files = msg.files
		m.err = msg.err
		m.projectDir = msg.projectDir
		return m, nil
	}

	if m.phase == phaseForm {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		if m.form.State == huh.StateCompleted {
			m.phase = phaseGenerating
			cfg := m.cfg
			return m, func() tea.Msg {
				cwd, _ := os.Getwd()
				files, err := scaffold.Scaffold(cfg, cwd)
				projectDir := filepath.Join(cwd, cfg.Name)
				return scaffoldDoneMsg{files: files, err: err, projectDir: projectDir}
			}
		}
		return m, cmd
	}

	return m, nil
}

func (m InitModel) View() string {
	var b strings.Builder

	b.WriteString(components.Header())
	b.WriteString(styles.Title.Render("  " + i18n.T("init.title")))
	b.WriteString(styles.Subtle.Render("  " + i18n.T("init.subtitle") + "\n\n"))

	switch m.phase {
	case phaseForm:
		b.WriteString(m.form.View())
	case phaseGenerating:
		sp := components.NewSpinner()
		b.WriteString(i18n.TF("init.generating", sp.View()) + "\n")
	case phaseDone:
		if m.err != nil {
			b.WriteString(styles.Error.Render(i18n.TF("init.error", m.err)))
			b.WriteString("\n")
		} else {
			b.WriteString(styles.Success.Render(i18n.TF("init.success", m.cfg.Name)))
			b.WriteString("\n\n")
			b.WriteString(styles.Subtle.Render(i18n.T("init.files_generated") + "\n"))
			for _, f := range m.files {
				b.WriteString(styles.Success.Render("    ✓ "))
				b.WriteString(styles.Subtle.Render(f))
				b.WriteString("\n")
			}
			b.WriteString("\n")
			b.WriteString(styles.Success.Render("  ▸ "))
			b.WriteString(styles.Title.Render(i18n.T("init.press_enter_chat")))
			b.WriteString("\n")
		}
	}

	b.WriteString(components.Footer(i18n.T("init.footer")))
	return b.String()
}

func (m *InitModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
