package wizard

import (
	"strings"

	"github.com/cfpperche/vibeforge/internal/product/brief"
)

// ReadinessLevel indicates how ready the user is.
type ReadinessLevel string

const (
	ReadyClear   ReadinessLevel = "clear"
	ReadyRefine  ReadinessLevel = "refine"
	ReadyProblem ReadinessLevel = "problem"
	ReadyExplore ReadinessLevel = "explore"
)

// WizardData collects all raw form data before converting to Brief.
type WizardData struct {
	// Round 0
	Readiness string

	// Round 1
	Name    string
	Tagline string
	Scale   string

	// Round 2
	Problem  string
	Audience []string

	// Round 3
	Hook         string
	ShareTrigger string
	Loop         string

	// Round 4
	Mechanics []string

	// Round 5
	Monetization string
	MomTest      string
}

// ToBrief converts the raw wizard data into a structured Brief.
func (d *WizardData) ToBrief() *brief.Brief {
	b := &brief.Brief{
		Name:     d.Name,
		Tagline:  d.Tagline,
		Category: brief.BuildScale(d.Scale),
		Problem:  d.Problem,
		Audience: d.Audience,
		Hook:     d.Hook,
		ShareTrigger: d.ShareTrigger,
	}

	// Loop — split into day estimates based on the single text
	b.Loop = brief.LoopDays{
		Day1:  d.Loop,
		Day7:  "TODO: refinar apos validacao",
		Day30: "TODO: refinar apos validacao",
		Day90: "TODO: refinar apos validacao",
	}

	// Mechanics
	for _, key := range d.Mechanics {
		if m, ok := brief.MechanicCatalog[key]; ok {
			b.Mechanics = append(b.Mechanics, m)
		}
	}

	// Monetization
	b.Monetization = []brief.RevenueStream{
		{
			Name:      brief.MonetizationLabels[d.Monetization],
			Rationale: "Definido no wizard",
		},
	}
	b.FreeForever = d.Monetization == "free"

	// Viral coefficient estimate based on mechanics
	b.ViralCoef = estimateViralCoef(d.Mechanics)

	// Viral loop diagram
	b.ViralLoop = generateViralDiagram(b)

	// Mom test
	lines := strings.SplitN(d.MomTest, "\n", 2)
	if len(lines) >= 1 {
		b.ShowYourMomTest[0] = strings.TrimSpace(lines[0])
	}
	if len(lines) >= 2 {
		b.ShowYourMomTest[1] = strings.TrimSpace(lines[1])
	} else if len(lines) == 1 {
		b.ShowYourMomTest[1] = "TODO"
	}

	// Default risks
	b.Risks = defaultRisks(b)

	return b
}

func estimateViralCoef(mechanics []string) brief.ViralPotential {
	score := 0
	viralMechanics := map[string]int{
		"identity":      3,
		"buildinpublic": 3,
		"leagues":       2,
		"seasonal":      2,
		"collection":    1,
		"streak":        1,
		"pet":           2,
		"copresence":    1,
		"spatial":       1,
		"idle":          0,
	}
	for _, m := range mechanics {
		score += viralMechanics[m]
	}
	if score >= 5 {
		return brief.ViralHigh
	}
	if score >= 3 {
		return brief.ViralMedium
	}
	return brief.ViralLow
}

func generateViralDiagram(b *brief.Brief) string {
	var s strings.Builder
	s.WriteString("Usuario descobre\n")
	s.WriteString("      |\n")
	s.WriteString("      v\n")
	s.WriteString("   [HOOK] " + truncate(b.Hook, 40) + "\n")
	s.WriteString("      |\n")
	s.WriteString("      v\n")
	s.WriteString("   [VALUE] Usa o produto\n")
	s.WriteString("      |\n")
	s.WriteString("      v\n")
	s.WriteString("   [TRIGGER] " + truncate(b.ShareTrigger, 40) + "\n")
	s.WriteString("      |\n")
	s.WriteString("      v\n")
	s.WriteString("   [SHARE] Amigo ve → volta ao topo\n")
	return s.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func defaultRisks(b *brief.Brief) []brief.Risk {
	risks := []brief.Risk{
		{
			Description: "Produto nao resolve dor real",
			Severity:    "high",
			Mitigation:  "Validar com 5 usuarios antes de buildar",
		},
		{
			Description: "Loop viral nao funciona organicamente",
			Severity:    "medium",
			Mitigation:  "Testar share trigger com 10 usuarios",
		},
	}
	if b.Category == brief.ScaleWeekend {
		risks = append(risks, brief.Risk{
			Description: "Scope creep — weekend vira mes",
			Severity:    "medium",
			Mitigation:  "Definir hard deadline de 3 dias",
		})
	}
	return risks
}
