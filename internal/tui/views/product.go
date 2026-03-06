package views

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/cfpperche/vibeforge/internal/config"
	"github.com/cfpperche/vibeforge/internal/i18n"
	"github.com/cfpperche/vibeforge/internal/product/brief"
	"github.com/cfpperche/vibeforge/internal/product/wizard"
	"github.com/cfpperche/vibeforge/internal/scaffold"
	"github.com/cfpperche/vibeforge/internal/tui/components"
	"github.com/cfpperche/vibeforge/internal/tui/styles"
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
				Title(i18n.T("product.readiness.question")).
				Description(i18n.T("product.readiness.desc")).
				Options(
					huh.NewOption(i18n.T("product.readiness.clear"), "clear"),
					huh.NewOption(i18n.T("product.readiness.refine"), "refine"),
					huh.NewOption(i18n.T("product.readiness.problem"), "problem"),
					huh.NewOption(i18n.T("product.readiness.explore"), "explore"),
				).
				Value(&m.data.Readiness),
		),
	).WithTheme(huh.ThemeCharm())

	// Round 1 — Identity
	forms[1] = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(i18n.T("product.name")).
				Placeholder(i18n.T("product.name.placeholder")).
				Value(&m.data.Name),
			huh.NewInput().
				Title(i18n.T("product.tagline")).
				Description(i18n.T("product.tagline.desc")).
				Placeholder(i18n.T("product.tagline.placeholder")).
				Value(&m.data.Tagline),
			huh.NewSelect[string]().
				Title(i18n.T("product.scale")).
				Options(
					huh.NewOption(i18n.T("product.scale.weekend"), "weekend"),
					huh.NewOption(i18n.T("product.scale.side"), "side"),
					huh.NewOption(i18n.T("product.scale.serious"), "serious"),
				).
				Value(&m.data.Scale),
		),
	).WithTheme(huh.ThemeCharm())

	// Round 2 — Problem & Audience
	forms[2] = huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title(i18n.T("product.problem")).
				Description(i18n.T("product.problem.desc")).
				CharLimit(300).
				Value(&m.data.Problem),
			huh.NewMultiSelect[string]().
				Title(i18n.T("product.audience")).
				Options(
					huh.NewOption(i18n.T("product.audience.dev_core"), "dev_core"),
					huh.NewOption(i18n.T("product.audience.non_tech"), "non_tech"),
					huh.NewOption(i18n.T("product.audience.influencer"), "influencer"),
					huh.NewOption(i18n.T("product.audience.general"), "general"),
				).
				Value(&m.data.Audience),
		),
	).WithTheme(huh.ThemeCharm())

	// Round 3 — Hook, Share Trigger, Loop
	forms[3] = huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title(i18n.T("product.hook")).
				Description(i18n.T("product.hook.desc")).
				CharLimit(140).
				Value(&m.data.Hook),
			huh.NewText().
				Title(i18n.T("product.share_trigger")).
				Description(i18n.T("product.share_trigger.desc")).
				CharLimit(200).
				Value(&m.data.ShareTrigger),
			huh.NewText().
				Title(i18n.T("product.loop")).
				Description(i18n.T("product.loop.desc")).
				CharLimit(200).
				Value(&m.data.Loop),
		),
	).WithTheme(huh.ThemeCharm())

	// Round 4 — Mechanics
	forms[4] = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(i18n.T("product.mechanics")).
				Description(i18n.T("product.mechanics.desc")).
				Options(
					huh.NewOption(i18n.T("product.mechanics.identity"), "identity"),
					huh.NewOption(i18n.T("product.mechanics.idle"), "idle"),
					huh.NewOption(i18n.T("product.mechanics.streak"), "streak"),
					huh.NewOption(i18n.T("product.mechanics.pet"), "pet"),
					huh.NewOption(i18n.T("product.mechanics.collection"), "collection"),
					huh.NewOption(i18n.T("product.mechanics.leagues"), "leagues"),
					huh.NewOption(i18n.T("product.mechanics.copresence"), "copresence"),
					huh.NewOption(i18n.T("product.mechanics.spatial"), "spatial"),
					huh.NewOption(i18n.T("product.mechanics.seasonal"), "seasonal"),
					huh.NewOption(i18n.T("product.mechanics.buildinpublic"), "buildinpublic"),
				).
				Value(&m.data.Mechanics),
		),
	).WithTheme(huh.ThemeCharm())

	// Round 5 — Monetization & Mom Test
	forms[5] = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(i18n.T("product.monetization")).
				Options(
					huh.NewOption(i18n.T("product.monetization.free"), "free"),
					huh.NewOption(i18n.T("product.monetization.freemium"), "freemium"),
					huh.NewOption(i18n.T("product.monetization.subscription"), "subscription"),
					huh.NewOption(i18n.T("product.monetization.onetime"), "onetime"),
					huh.NewOption(i18n.T("product.monetization.b2b"), "b2b"),
				).
				Value(&m.data.Monetization),
			huh.NewText().
				Title(i18n.T("product.mom_test")).
				Description(i18n.T("product.mom_test.desc")).
				Placeholder(i18n.T("product.mom_test.placeholder")).
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
			summary := i18n.TF("product.done_summary", name, len(m.files))
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
		b.WriteString(styles.Title.Render("  " + i18n.T("product.title")))
		b.WriteString(styles.Subtle.Render("  " + i18n.T("product.generating_subtitle") + "\n\n"))
		sp := components.NewSpinner()
		b.WriteString(fmt.Sprintf("\n  %s %s\n", sp.View(), i18n.T("product.generating_msg")))
	case productDone:
		b.WriteString(m.viewDone())
	}

	return b.String()
}

func (m ProductModel) viewWizard() string {
	var b strings.Builder

	roundLabels := []string{
		i18n.T("product.warmup_title"),
		i18n.T("product.identity_title"),
		i18n.T("product.problem_title"),
		i18n.T("product.hook_title"),
		i18n.T("product.mechanics_title"),
		i18n.T("product.monetization_title"),
	}

	b.WriteString(styles.Title.Render("  " + i18n.T("product.title")))
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

	b.WriteString(components.Footer(i18n.T("product.wizard_footer")))

	return b.String()
}

func (m ProductModel) viewConfirm() string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("  " + i18n.T("product.title")))
	b.WriteString(styles.Subtle.Render("  " + i18n.T("product.confirm_subtitle") + "\n\n"))

	if m.brief == nil {
		return b.String()
	}

	br := m.brief

	// Build the confirmation box
	var lines []string

	lines = append(lines, styles.Title.Bold(true).Render(i18n.TF("product.brief_title", br.Name)))
	lines = append(lines, "")
	lines = append(lines, styles.Subtle.Render(fmt.Sprintf("\"%s\"", br.Tagline)))
	lines = append(lines, "")

	lines = append(lines, fmt.Sprintf("  %s  %s",
		styles.Subtle.Render(i18n.T("product.scale_label")),
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
		styles.Subtle.Render(i18n.T("product.audience_label")),
		styles.Title.Render(strings.Join(audLabels, " + ")),
	))
	lines = append(lines, "")

	lines = append(lines, fmt.Sprintf("  %s  %s",
		styles.Success.Render(i18n.T("product.hook_label")),
		truncateView(br.Hook, 50),
	))
	lines = append(lines, fmt.Sprintf("  %s  %s",
		styles.Warning.Render(i18n.T("product.trigger_label")),
		truncateView(br.ShareTrigger, 50),
	))
	lines = append(lines, fmt.Sprintf("  %s  %s",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#60a5fa")).Render(i18n.T("product.loop_label")),
		truncateView(br.Loop.Day1, 50),
	))
	lines = append(lines, "")

	// Mechanics
	lines = append(lines, styles.Subtle.Render(i18n.T("product.mechanics_label")))
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
			styles.Subtle.Render(i18n.T("product.monetization_label")),
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
		styles.Subtle.Render(i18n.T("product.viral_label")),
		viralStyle.Render(string(br.ViralCoef)),
	))

	content := strings.Join(lines, "\n")
	box := styles.ActiveBox.Width(60).Render(content)
	b.WriteString(box)
	b.WriteString("\n\n")

	b.WriteString(components.Footer(i18n.T("product.confirm_footer")))

	return b.String()
}

func (m ProductModel) viewDone() string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("  " + i18n.T("product.title")))

	if m.err != nil {
		if m.err.Error() == "explore" {
			b.WriteString(styles.Subtle.Render("  " + i18n.T("product.explore_subtitle") + "\n\n"))
			b.WriteString(styles.Warning.Render(i18n.T("product.explore_instruction") + "\n"))
			b.WriteString(styles.Subtle.Render(i18n.T("product.explore_detail1") + "\n"))
			b.WriteString(styles.Subtle.Render(i18n.T("product.explore_detail2") + "\n\n"))
			b.WriteString(styles.Subtle.Render(i18n.T("product.explore_detail3") + "\n"))
			b.WriteString(components.Footer(i18n.T("product.explore_footer")))
			return b.String()
		}
		b.WriteString(styles.Subtle.Render("  " + i18n.T("product.error_subtitle") + "\n\n"))
		b.WriteString(styles.Error.Render(i18n.TF("product.error_msg", m.err)))
		b.WriteString("\n")
		b.WriteString(components.Footer(i18n.T("product.error_footer")))
		return b.String()
	}

	b.WriteString(styles.Subtle.Render("  " + i18n.T("product.done_subtitle") + "\n\n"))
	b.WriteString(styles.Success.Render(i18n.TF("product.done_success", m.data.Name)))
	b.WriteString("\n\n")

	b.WriteString(styles.Subtle.Render(i18n.T("product.done_files") + "\n"))
	for _, f := range m.files {
		b.WriteString(styles.Success.Render("    + "))
		b.WriteString(styles.Subtle.Render(f))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	b.WriteString(styles.Success.Render("  > "))
	b.WriteString(styles.Title.Render(i18n.T("product.done_press_enter")))
	b.WriteString("\n")

	b.WriteString(components.Footer(i18n.T("product.done_footer")))

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
