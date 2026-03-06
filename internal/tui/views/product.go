package views

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/cfpperche/vibescaffold/internal/product/brief"
	"github.com/cfpperche/vibescaffold/internal/product/wizard"
	"github.com/cfpperche/vibescaffold/internal/scaffold"
	"github.com/cfpperche/vibescaffold/internal/config"
	"github.com/cfpperche/vibescaffold/internal/tui/components"
	"github.com/cfpperche/vibescaffold/internal/tui/styles"
)

type productPhase int

const (
	productWizard productPhase = iota
	productConfirm
	productGenerating
	productDone
)

type productBriefDoneMsg struct {
	briefFiles    []string
	scaffoldFiles []string
	err           error
	projectDir    string
}

type ProductModel struct {
	width  int
	height int
	phase  productPhase
	round  int
	data   wizard.WizardData
	forms  []*huh.Form
	brief  *brief.Brief
	files  []string
	err    error
	projectDir string
}

func NewProduct() ProductModel {
	m := ProductModel{
		phase: productWizard,
		round: 0,
	}
	m.forms = m.buildForms()
	return m
}

func (m *ProductModel) buildForms() []*huh.Form {
	forms := make([]*huh.Form, 6)

	// Round 0 — Warmup
	forms[0] = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Voce ja sabe o que quer construir?").
				Description("Antes de codar, vamos entender o produto.\nO maior erro do vibecoding e pular essa parte.").
				Options(
					huh.NewOption("Tenho ideia clara — so quero estruturar", "clear"),
					huh.NewOption("Tenho uma direcao — preciso refinar", "refine"),
					huh.NewOption("Tenho um problema — nao sei a solucao ainda", "problem"),
					huh.NewOption("Quero explorar — me mostre possibilidades", "explore"),
				).
				Value(&m.data.Readiness),
		),
	).WithTheme(huh.ThemeCharm())

	// Round 1 — Identity
	forms[1] = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Nome do projeto").
				Placeholder("vibescaffold").
				Value(&m.data.Name),
			huh.NewInput().
				Title("Tagline — maximo 10 palavras").
				Description("Deve fazer um nao-dev ficar curioso").
				Placeholder("Scaffold que pensa no produto antes do codigo").
				Value(&m.data.Tagline),
			huh.NewSelect[string]().
				Title("Escala de ambicao").
				Options(
					huh.NewOption("Weekend vibe — 2-3 dias, ship e veja se cola", "weekend"),
					huh.NewOption("Side project — 2-4 semanas, MVP polido", "side"),
					huh.NewOption("Produto serio — 1-3 meses, monetizacao real", "serious"),
				).
				Value(&m.data.Scale),
		),
	).WithTheme(huh.ThemeCharm())

	// Round 2 — Problem & Audience
	forms[2] = huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Qual dor esse produto resolve?").
				Description("Seja especifico. 'Developers perdem tempo com X porque Y'").
				CharLimit(300).
				Value(&m.data.Problem),
			huh.NewMultiSelect[string]().
				Title("Quem deve usar isso?").
				Options(
					huh.NewOption("Developers que vibecoding diariamente (core)", "dev_core"),
					huh.NewOption("Builders nao-tecnicos usando IA (crossover)", "non_tech"),
					huh.NewOption("Tech Twitter / influenciadores dev (amplificacao)", "influencer"),
					huh.NewOption("Publico geral — viral alem da bolha dev", "general"),
				).
				Value(&m.data.Audience),
		),
	).WithTheme(huh.ThemeCharm())

	// Round 3 — Hook, Share Trigger, Loop
	forms[3] = huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("O HOOK — por que alguem tenta na primeira vez?").
				Description("Deve caber em um tweet. Se nao cabe, e complexo demais.\n\nExemplos:\n  OK: 'Transforma seus commits em uma cidade pixel art'\n  OK: 'Seu AI agent tem um pet que morre se voce parar de codar'\n  Ruim: 'Plataforma de desenvolvimento com IA integrada'").
				CharLimit(140).
				Value(&m.data.Hook),
			huh.NewText().
				Title("O SHARE TRIGGER — qual momento exato faz compartilhar?").
				Description("Screenshot? Comparacao? Resultado absurdo? Conquista?\n\nSeja especifico sobre o MOMENTO:\n  OK: 'Quando veem o predio deles maior que o do amigo'\n  Ruim: 'Quando tem uma boa experiencia'").
				CharLimit(200).
				Value(&m.data.ShareTrigger),
			huh.NewText().
				Title("O LOOP — por que voltam amanha?").
				Description("O que muda entre hoje e amanha?").
				CharLimit(200).
				Value(&m.data.Loop),
		),
	).WithTheme(huh.ThemeCharm())

	// Round 4 — Mechanics
	forms[4] = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Quais mecanicas de engajamento fazem sentido?").
				Description("Escolha 2-4 que se complementam").
				Options(
					huh.NewOption("Identity Visualization — atividade vira artefato visual", "identity"),
					huh.NewOption("Idle / Incremental — progresso enquanto voce esta fora", "idle"),
					huh.NewOption("Streak + Loss Aversion — sequencia diaria com punicao", "streak"),
					huh.NewOption("Virtual Pet / Companion — criatura que reflete atividade", "pet"),
					huh.NewOption("Collection / Completionism — gotta catch 'em all", "collection"),
					huh.NewOption("Competitive Tiers / Leagues — promocao/rebaixamento", "leagues"),
					huh.NewOption("Body Doubling / Co-presence — trabalha ao lado de outros", "copresence"),
					huh.NewOption("Spatial / World Metaphor — espaco navegavel", "spatial"),
					huh.NewOption("Seasonal / Event-driven — FOMO temporal", "seasonal"),
					huh.NewOption("Build-in-public — processo visivel e compartilhavel", "buildinpublic"),
				).
				Value(&m.data.Mechanics),
		),
	).WithTheme(huh.ThemeCharm())

	// Round 5 — Monetization & Mom Test
	forms[5] = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Modelo de monetizacao inicial").
				Options(
					huh.NewOption("Free forever — cresce por viral, monetiza depois", "free"),
					huh.NewOption("Freemium — core gratis, features premium pagas", "freemium"),
					huh.NewOption("Subscription — valor recorrente desde o inicio", "subscription"),
					huh.NewOption("One-time — paga uma vez, usa sempre", "onetime"),
					huh.NewOption("B2B / Teams — vende para empresas", "b2b"),
				).
				Value(&m.data.Monetization),
			huh.NewText().
				Title("Teste 'Mostre pra sua mae'").
				Description("Explique em 2 frases para alguem nao-tecnico.\nSe nao consegue, o conceito precisa simplificacao.").
				Placeholder("Frase 1: O que e.\nFrase 2: Por que e legal.").
				CharLimit(300).
				Value(&m.data.MomTest),
		),
	).WithTheme(huh.ThemeCharm())

	return forms
}

func (m ProductModel) Init() tea.Cmd {
	if len(m.forms) > 0 {
		return m.forms[0].Init()
	}
	return nil
}

func (m ProductModel) Update(msg tea.Msg) (ProductModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			switch m.phase {
			case productConfirm:
				// Go back to last form round
				m.phase = productWizard
				m.round = len(m.forms) - 1
				m.forms = m.buildForms() // rebuild to reset state
				return m, m.forms[m.round].Init()
			case productDone:
				if m.err == nil {
					return m, func() tea.Msg { return NavigateMsg{Target: "home"} }
				}
				return m, func() tea.Msg { return NavigateMsg{Target: "home"} }
			case productWizard:
				return m, func() tea.Msg { return NavigateMsg{Target: "home"} }
			}
		}

		if msg.String() == "enter" && m.phase == productConfirm {
			m.phase = productGenerating
			data := m.data
			b := data.ToBrief()
			return m, func() tea.Msg {
				return doGenerate(b, &data)
			}
		}

		if msg.String() == "enter" && m.phase == productDone && m.err == nil {
			dir := m.projectDir
			name := m.data.Name
			summary := fmt.Sprintf("Projeto '%s' criado — %d arquivos", name, len(m.files))
			return m, func() tea.Msg {
				return EnterChatMsg{
					ProjectDir:  dir,
					ProjectName: name,
					Summary:     summary,
				}
			}
		}

	case productBriefDoneMsg:
		m.phase = productDone
		m.err = msg.err
		m.projectDir = msg.projectDir
		m.files = append(msg.briefFiles, msg.scaffoldFiles...)
		return m, nil
	}

	if m.phase == productWizard && m.round < len(m.forms) {
		form, cmd := m.forms[m.round].Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.forms[m.round] = f
		}

		if m.forms[m.round].State == huh.StateCompleted {
			// Check for "explore" exit
			if m.round == 0 && m.data.Readiness == "explore" {
				m.phase = productDone
				m.err = fmt.Errorf("explore")
				return m, nil
			}

			m.round++
			if m.round >= len(m.forms) {
				// All rounds done — show confirmation
				m.brief = m.data.ToBrief()
				m.phase = productConfirm
				return m, nil
			}
			return m, m.forms[m.round].Init()
		}
		return m, cmd
	}

	return m, nil
}

func doGenerate(b *brief.Brief, data *wizard.WizardData) productBriefDoneMsg {
	cwd, _ := os.Getwd()
	projectDir := filepath.Join(cwd, data.Name)

	// Generate product brief files
	briefFiles, err := brief.Generate(b, projectDir)
	if err != nil {
		return productBriefDoneMsg{err: err, projectDir: projectDir}
	}

	// Run scaffold
	cfg := config.Config{
		Name:     data.Name,
		Desc:     b.Problem,
		Type:     "saas",
		Stack:    "go-bubbletea",
		Features: []string{"claude-md", "context-docs", "roadmap"},
	}

	scaffoldFiles, err := scaffold.Scaffold(cfg, cwd)
	if err != nil {
		return productBriefDoneMsg{
			briefFiles: briefFiles,
			err:        err,
			projectDir: projectDir,
		}
	}

	return productBriefDoneMsg{
		briefFiles:    briefFiles,
		scaffoldFiles: scaffoldFiles,
		projectDir:    projectDir,
	}
}

func (m ProductModel) View() string {
	var b strings.Builder

	b.WriteString(components.Header())

	switch m.phase {
	case productWizard:
		b.WriteString(m.viewWizard())
	case productConfirm:
		b.WriteString(m.viewConfirm())
	case productGenerating:
		b.WriteString(styles.Title.Render("  $ new"))
		b.WriteString(styles.Subtle.Render("  — gerando produto\n\n"))
		sp := components.NewSpinner()
		b.WriteString(fmt.Sprintf("\n  %s Gerando brief e scaffold...\n", sp.View()))
	case productDone:
		b.WriteString(m.viewDone())
	}

	return b.String()
}

func (m ProductModel) viewWizard() string {
	var b strings.Builder

	roundLabels := []string{
		"Aquecimento",
		"Identidade do produto",
		"Problema e audiencia",
		"Hook, Loop e Share Trigger",
		"Mecanicas de engajamento",
		"Monetizacao e riscos",
	}

	b.WriteString(styles.Title.Render("  $ new"))
	b.WriteString(styles.Subtle.Render(fmt.Sprintf("  — %s", roundLabels[m.round])))
	b.WriteString(styles.Subtle.Render(fmt.Sprintf("  (%d/%d)\n\n", m.round+1, len(m.forms))))

	// Progress bar
	progress := ""
	for i := range m.forms {
		if i < m.round {
			progress += styles.Success.Render("●")
		} else if i == m.round {
			progress += styles.Warning.Render("●")
		} else {
			progress += styles.Subtle.Render("○")
		}
		if i < len(m.forms)-1 {
			progress += styles.Subtle.Render("─")
		}
	}
	b.WriteString("  " + progress + "\n\n")

	if m.round < len(m.forms) {
		b.WriteString(m.forms[m.round].View())
	}

	b.WriteString(components.Footer("  [enter] proximo  [esc] voltar"))

	return b.String()
}

func (m ProductModel) viewConfirm() string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("  $ new"))
	b.WriteString(styles.Subtle.Render("  — confirmar brief\n\n"))

	if m.brief == nil {
		return b.String()
	}

	br := m.brief

	// Build the confirmation box
	var lines []string

	lines = append(lines, styles.Title.Bold(true).Render(fmt.Sprintf("Product Brief — %s", br.Name)))
	lines = append(lines, "")
	lines = append(lines, styles.Subtle.Render(fmt.Sprintf("\"%s\"", br.Tagline)))
	lines = append(lines, "")

	lines = append(lines, fmt.Sprintf("  %s  %s",
		styles.Subtle.Render("Escala:"),
		styles.Title.Render(brief.ScaleLabel(br.Category)),
	))

	// Audience
	var audLabels []string
	for _, a := range br.Audience {
		if l, ok := brief.AudienceLabels[a]; ok {
			audLabels = append(audLabels, l)
		}
	}
	lines = append(lines, fmt.Sprintf("  %s  %s",
		styles.Subtle.Render("Audiencia:"),
		styles.Title.Render(strings.Join(audLabels, " + ")),
	))
	lines = append(lines, "")

	lines = append(lines, fmt.Sprintf("  %s  %s",
		styles.Success.Render("Hook:"),
		truncateView(br.Hook, 50),
	))
	lines = append(lines, fmt.Sprintf("  %s  %s",
		styles.Warning.Render("Trigger:"),
		truncateView(br.ShareTrigger, 50),
	))
	lines = append(lines, fmt.Sprintf("  %s  %s",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#60a5fa")).Render("Loop:"),
		truncateView(br.Loop.Day1, 50),
	))
	lines = append(lines, "")

	// Mechanics
	lines = append(lines, styles.Subtle.Render("  Mecanicas:"))
	for _, mech := range br.Mechanics {
		lines = append(lines, fmt.Sprintf("    %s %s",
			styles.Success.Render("●"),
			styles.Title.Render(mech.Name),
		))
	}
	lines = append(lines, "")

	// Monetization
	if len(br.Monetization) > 0 {
		lines = append(lines, fmt.Sprintf("  %s  %s",
			styles.Subtle.Render("Monetizacao:"),
			styles.Title.Render(br.Monetization[0].Name),
		))
	}

	// Viral coef
	viralStyle := styles.Subtle
	switch br.ViralCoef {
	case brief.ViralHigh:
		viralStyle = styles.Success
	case brief.ViralMedium:
		viralStyle = styles.Warning
	}
	lines = append(lines, fmt.Sprintf("  %s  %s",
		styles.Subtle.Render("Viral:"),
		viralStyle.Render(string(br.ViralCoef)),
	))

	content := strings.Join(lines, "\n")
	box := styles.ActiveBox.Width(60).Render(content)
	b.WriteString(box)
	b.WriteString("\n\n")

	b.WriteString(components.Footer("  [enter] gerar brief e scaffold  [esc] ajustar"))

	return b.String()
}

func (m ProductModel) viewDone() string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("  $ new"))

	if m.err != nil {
		if m.err.Error() == "explore" {
			b.WriteString(styles.Subtle.Render("  — explorar\n\n"))
			b.WriteString(styles.Warning.Render("  Para exploracao livre, use o Vibe Forge:\n"))
			b.WriteString(styles.Subtle.Render("  /spock-vibe-forge no Claude Code\n"))
			b.WriteString(styles.Subtle.Render("  Ele pesquisa tendencias e gera conceitos virais.\n\n"))
			b.WriteString(styles.Subtle.Render("  Voltando com uma ideia? Rode: vs new\n"))
			b.WriteString(components.Footer("  [esc] voltar"))
			return b.String()
		}
		b.WriteString(styles.Subtle.Render("  — erro\n\n"))
		b.WriteString(styles.Error.Render(fmt.Sprintf("\n  Erro: %s\n", m.err)))
		b.WriteString(components.Footer("  [esc] voltar"))
		return b.String()
	}

	b.WriteString(styles.Subtle.Render("  — produto criado\n\n"))
	b.WriteString(styles.Success.Render(fmt.Sprintf("  Projeto '%s' criado!\n\n", m.data.Name)))

	b.WriteString(styles.Subtle.Render("  Arquivos gerados:\n"))
	for _, f := range m.files {
		b.WriteString(styles.Success.Render("    + "))
		b.WriteString(styles.Subtle.Render(f))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	b.WriteString(styles.Success.Render("  > "))
	b.WriteString(styles.Title.Render("Pressione [enter] para abrir o chat"))
	b.WriteString("\n")

	b.WriteString(components.Footer("  [enter] abrir chat  [esc] voltar"))

	return b.String()
}

func truncateView(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func (m *ProductModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
