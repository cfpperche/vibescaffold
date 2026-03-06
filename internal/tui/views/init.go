package views

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/cfpperche/vibeforge/internal/config"
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
				Title("Nome do projeto").
				Placeholder("meu-projeto").
				Value(&cfg.Name),
			huh.NewInput().
				Title("Descricao").
				Placeholder("Descreva seu projeto...").
				Value(&cfg.Desc),
			huh.NewInput().
				Title("Autor").
				Placeholder("seu-nome").
				Value(&cfg.Author),
			huh.NewConfirm().
				Title("Criar repo privado no GitHub?").
				Value(&cfg.Repo),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Tipo de projeto").
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
				Title("Stack").
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
				Title("Principios").
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
				Title("Ferramentas Claude").
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
			summary := fmt.Sprintf("✓ Projeto '%s' scaffolado — %d arquivos criados", name, len(m.files))
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
	b.WriteString(styles.Title.Render("  $ init"))
	b.WriteString(styles.Subtle.Render("  — scaffold de projeto\n\n"))

	switch m.phase {
	case phaseForm:
		b.WriteString(m.form.View())
	case phaseGenerating:
		sp := components.NewSpinner()
		b.WriteString(fmt.Sprintf("\n  %s Gerando scaffold...\n", sp.View()))
	case phaseDone:
		if m.err != nil {
			b.WriteString(styles.Error.Render(fmt.Sprintf("\n  ✗ Erro: %s\n", m.err)))
		} else {
			b.WriteString(styles.Success.Render(fmt.Sprintf("\n  ✓ Projeto '%s' criado!\n\n", m.cfg.Name)))
			b.WriteString(styles.Subtle.Render("  Arquivos gerados:\n"))
			for _, f := range m.files {
				b.WriteString(styles.Success.Render("    ✓ "))
				b.WriteString(styles.Subtle.Render(f))
				b.WriteString("\n")
			}
			b.WriteString("\n")
			b.WriteString(styles.Success.Render("  ▸ "))
			b.WriteString(styles.Title.Render("Pressione [enter] para abrir o chat"))
			b.WriteString("\n")
		}
	}

	b.WriteString(components.Footer("\n  [enter] abrir chat  [esc] voltar"))
	return b.String()
}

func (m *InitModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
