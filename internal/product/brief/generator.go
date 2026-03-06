package brief

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Generate creates all product documentation files from the brief.
func Generate(b *Brief, projectDir string) ([]string, error) {
	type genFile struct {
		path    string
		content func(*Brief) string
	}

	files := []genFile{
		{"docs/PRODUCT_BRIEF.md", generateProductBrief},
		{"docs/PERSONA.md", generatePersona},
		{"docs/VIRAL_LOOP.md", generateViralLoop},
	}

	var created []string
	for _, f := range files {
		content := f.content(b)
		p := filepath.Join(projectDir, f.path)
		dir := filepath.Dir(p)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return created, fmt.Errorf("mkdir %s: %w", dir, err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			return created, fmt.Errorf("write %s: %w", f.path, err)
		}
		created = append(created, f.path)
	}
	return created, nil
}

func generateProductBrief(b *Brief) string {
	var s strings.Builder

	fmt.Fprintf(&s, "# %s\n\n", b.Name)
	fmt.Fprintf(&s, "**Tagline:** %s\n\n", b.Tagline)
	fmt.Fprintf(&s, "**Category:** %s\n\n", ScaleLabel(b.Category))

	if len(b.Comparables) > 0 {
		s.WriteString("**Comparable(s):** ")
		var names []string
		for _, c := range b.Comparables {
			names = append(names, c.Name)
		}
		s.WriteString(strings.Join(names, " meets "))
		s.WriteString("\n\n")
	}

	// Problem & Solution
	s.WriteString("## The Problem\n\n")
	fmt.Fprintf(&s, "%s\n\n", b.Problem)

	if b.Solution != "" {
		s.WriteString("## The Solution\n\n")
		fmt.Fprintf(&s, "%s\n\n", b.Solution)
	}

	// Audience
	s.WriteString("## Audience\n\n")
	for _, a := range b.Audience {
		label := a
		if l, ok := AudienceLabels[a]; ok {
			label = l
		}
		fmt.Fprintf(&s, "- %s\n", label)
	}
	s.WriteString("\n")

	// Hook
	s.WriteString("## The Hook\n\n")
	fmt.Fprintf(&s, "> %s\n\n", b.Hook)

	// Loop
	s.WriteString("## The Loop\n\n")
	fmt.Fprintf(&s, "| Dia | O que acontece |\n")
	fmt.Fprintf(&s, "|-----|----------------|\n")
	fmt.Fprintf(&s, "| Day 1 | %s |\n", b.Loop.Day1)
	fmt.Fprintf(&s, "| Day 7 | %s |\n", b.Loop.Day7)
	fmt.Fprintf(&s, "| Day 30 | %s |\n", b.Loop.Day30)
	fmt.Fprintf(&s, "| Day 90 | %s |\n\n", b.Loop.Day90)

	// Share Trigger
	s.WriteString("## The Share Trigger\n\n")
	fmt.Fprintf(&s, "> %s\n\n", b.ShareTrigger)

	// Mechanics
	s.WriteString("## Mechanics Breakdown\n\n")
	s.WriteString("| Mecanica | Layer | Descricao |\n")
	s.WriteString("|----------|-------|-----------|\n")
	for _, m := range b.Mechanics {
		fmt.Fprintf(&s, "| %s | %s | %s |\n", m.Name, m.Layer, m.Description)
	}
	s.WriteString("\n")

	// Viral Loop
	if b.ViralLoop != "" {
		s.WriteString("## Viral Loop\n\n")
		s.WriteString("```\n")
		fmt.Fprintf(&s, "%s\n", b.ViralLoop)
		s.WriteString("```\n\n")
		fmt.Fprintf(&s, "**Coeficiente viral estimado:** %s\n\n", b.ViralCoef)
	}

	// Monetization
	s.WriteString("## Monetization\n\n")
	if b.FreeForever {
		s.WriteString("Free forever — cresce por viral, monetiza depois.\n\n")
	}
	if len(b.Monetization) > 0 {
		s.WriteString("| Modelo | Preco | Racional |\n")
		s.WriteString("|--------|-------|----------|\n")
		for _, r := range b.Monetization {
			fmt.Fprintf(&s, "| %s | %s | %s |\n", r.Name, r.PricePoint, r.Rationale)
		}
		s.WriteString("\n")
	}

	// Risks
	if len(b.Risks) > 0 {
		s.WriteString("## What Could Kill This\n\n")
		s.WriteString("| Risco | Severidade | Mitigacao |\n")
		s.WriteString("|-------|------------|-----------|\n")
		for _, r := range b.Risks {
			fmt.Fprintf(&s, "| %s | %s | %s |\n", r.Description, r.Severity, r.Mitigation)
		}
		s.WriteString("\n")
	}

	// Show Your Mom Test
	s.WriteString("## The \"Show Your Mom\" Test\n\n")
	fmt.Fprintf(&s, "1. %s\n", b.ShowYourMomTest[0])
	fmt.Fprintf(&s, "2. %s\n", b.ShowYourMomTest[1])

	return s.String()
}

func generatePersona(b *Brief) string {
	var s strings.Builder

	fmt.Fprintf(&s, "# %s — Personas\n\n", b.Name)

	s.WriteString("## Publico Primario\n\n")
	for _, a := range b.Audience {
		label := a
		if l, ok := AudienceLabels[a]; ok {
			label = l
		}
		fmt.Fprintf(&s, "### %s\n\n", label)
		s.WriteString("- **Comportamento:** TODO\n")
		s.WriteString("- **Motivacao:** TODO\n")
		s.WriteString("- **Frustracao atual:** TODO\n")
		fmt.Fprintf(&s, "- **Job-to-be-done:** %s\n\n", b.Problem)
	}

	s.WriteString("## Anti-Persona\n\n")
	s.WriteString("Quem NAO e o usuario:\n")
	s.WriteString("- TODO: definir quem nao deve usar o produto\n\n")

	s.WriteString("## Jobs-to-be-done\n\n")
	fmt.Fprintf(&s, "1. **Funcional:** %s\n", b.Problem)
	s.WriteString("2. **Emocional:** TODO\n")
	s.WriteString("3. **Social:** TODO\n")

	return s.String()
}

func generateViralLoop(b *Brief) string {
	var s strings.Builder

	fmt.Fprintf(&s, "# %s — Viral Loop\n\n", b.Name)

	s.WriteString("## Diagrama\n\n")
	s.WriteString("```\n")
	if b.ViralLoop != "" {
		fmt.Fprintf(&s, "%s\n", b.ViralLoop)
	} else {
		s.WriteString("Usuario descobre → Testa (Hook) → Usa (Value) → Compartilha (Trigger) → Amigo descobre\n")
		s.WriteString("     ↑                                                                        |\n")
		s.WriteString("     └────────────────────────────────────────────────────────────────────────┘\n")
	}
	s.WriteString("```\n\n")

	fmt.Fprintf(&s, "## Coeficiente Viral Estimado: %s\n\n", b.ViralCoef)

	switch b.ViralCoef {
	case ViralLow:
		s.WriteString("K < 0.5 — crescimento depende de aquisicao paga ou organica.\n\n")
	case ViralMedium:
		s.WriteString("K 0.5-1.0 — viral ajuda mas nao sustenta sozinho.\n\n")
	case ViralHigh:
		s.WriteString("K > 1.0 — crescimento exponencial possivel.\n\n")
	}

	s.WriteString("## Share Trigger\n\n")
	fmt.Fprintf(&s, "> %s\n\n", b.ShareTrigger)

	s.WriteString("## Formatos de Share\n\n")
	s.WriteString("- [ ] Tweet / post curto\n")
	s.WriteString("- [ ] Screenshot / imagem compartilhavel\n")
	s.WriteString("- [ ] Story (Instagram/TikTok)\n")
	s.WriteString("- [ ] Embed interativo\n")
	s.WriteString("- [ ] Link direto com preview (OG tags)\n\n")

	s.WriteString("## Anti-patterns a Evitar\n\n")
	s.WriteString("- Forcar compartilhamento antes do valor ser entregue\n")
	s.WriteString("- Share gates (bloquear funcionalidade ate compartilhar)\n")
	s.WriteString("- Spam de notificacoes para amigos do usuario\n")
	s.WriteString("- Compartilhamento que nao gera curiosidade no receptor\n\n")

	s.WriteString("## Metricas\n\n")
	s.WriteString("| Metrica | Target | Como medir |\n")
	s.WriteString("|---------|--------|------------|\n")
	s.WriteString("| K-factor | > 0.5 | invites_sent * conversion_rate |\n")
	s.WriteString("| Share rate | > 10% | users_who_share / active_users |\n")
	s.WriteString("| Invite conversion | > 20% | signups_from_invite / invites_clicked |\n")
	s.WriteString("| Time to share | < 7 dias | median_days_to_first_share |\n")

	return s.String()
}
