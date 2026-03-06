package brief

type ViralPotential string
type BuildScale string

const (
	ViralLow    ViralPotential = "low"    // < 0.5 coeficiente
	ViralMedium ViralPotential = "medium" // 0.5-1.0
	ViralHigh   ViralPotential = "high"   // > 1.0

	ScaleWeekend BuildScale = "weekend" // 2-3 dias
	ScaleSide    BuildScale = "side"    // 2-4 semanas
	ScaleSerious BuildScale = "serious" // 1-3 meses
)

type Brief struct {
	// Identidade
	Name     string
	Tagline  string // max 10 palavras
	Category BuildScale

	// O produto
	Problem  string   // qual dor resolve
	Solution string   // como resolve
	Audience []string // publico primario + crossover

	// Viral mechanics
	Hook         string         // por que alguem tenta na primeira vez
	Loop         LoopDays       // por que voltam
	ShareTrigger string         // momento exato que faz compartilhar
	ViralLoop    string         // diagrama textual do loop
	ViralCoef    ViralPotential

	// Mecanicas escolhidas
	Mechanics []Mechanic

	// Modelo de negocio
	Monetization []RevenueStream
	FreeForever  bool

	// Referencias
	Comparables []Comparable // "X meets Y"

	// Testes de sanidade
	ShowYourMomTest [2]string // 2 frases para nao-devs
	Risks           []Risk
}

type LoopDays struct {
	Day1  string
	Day7  string
	Day30 string
	Day90 string
}

type Mechanic struct {
	Name        string
	Layer       string // "core" | "social" | "identity"
	Description string
}

type RevenueStream struct {
	Name       string
	PricePoint string
	Rationale  string
}

type Comparable struct {
	Name           string
	URL            string
	WhatWeBorrow   string
}

type Risk struct {
	Description string
	Severity    string // "high" | "medium" | "low"
	Mitigation  string
}

// ScaleLabel returns the human-readable label for the build scale.
func ScaleLabel(s BuildScale) string {
	switch s {
	case ScaleWeekend:
		return "Weekend vibe (2-3 dias)"
	case ScaleSide:
		return "Side project (2-4 semanas)"
	case ScaleSerious:
		return "Produto serio (1-3 meses)"
	default:
		return string(s)
	}
}

// AudienceLabels maps audience keys to human labels.
var AudienceLabels = map[string]string{
	"dev_core":   "Developers (core)",
	"non_tech":   "Builders nao-tecnicos",
	"influencer": "Tech Twitter / influenciadores",
	"general":    "Publico geral",
}

// MechanicCatalog is the full list of engagement mechanics.
var MechanicCatalog = map[string]Mechanic{
	"identity":      {Name: "Identity Visualization", Layer: "identity", Description: "Atividade vira artefato visual compartilhavel"},
	"idle":          {Name: "Idle / Incremental", Layer: "core", Description: "Progresso acontece enquanto voce esta fora"},
	"streak":        {Name: "Streak + Loss Aversion", Layer: "core", Description: "Sequencia diaria com punicao por quebrar"},
	"pet":           {Name: "Virtual Pet / Companion", Layer: "social", Description: "Criatura que reflete sua atividade"},
	"collection":    {Name: "Collection / Completionism", Layer: "core", Description: "Itens que acumulam, gotta catch 'em all"},
	"leagues":       {Name: "Competitive Tiers / Leagues", Layer: "social", Description: "Competicao temporal com promocao/rebaixamento"},
	"copresence":    {Name: "Body Doubling / Co-presence", Layer: "social", Description: "Trabalha ao lado de outros sem colaborar"},
	"spatial":       {Name: "Spatial / World Metaphor", Layer: "identity", Description: "Conceitos mapeados em espaco navegavel"},
	"seasonal":      {Name: "Seasonal / Event-driven", Layer: "social", Description: "Eventos temporais com FOMO"},
	"buildinpublic": {Name: "Build-in-public", Layer: "identity", Description: "Processo de criacao visivel e compartilhavel"},
}

// MonetizationLabels maps monetization keys to labels.
var MonetizationLabels = map[string]string{
	"free":         "Free forever",
	"freemium":     "Freemium",
	"subscription": "Subscription",
	"onetime":      "One-time",
	"b2b":          "B2B / Teams",
}
