package onboarding

type FillLevel int

const (
	FilledByScaffold FillLevel = iota // scaffold preenche tudo
	FilledPartial                     // scaffold preenche estrutura, agente completa
	FilledByAgent                     // scaffold cria vazio, agente preenche
)

type MDFile struct {
	Path          string
	FillLevel     FillLevel
	Description   string
	ScaffoldFills []string
	AgentFills    []string
}

var Files = []MDFile{
	{
		Path:        "CLAUDE.md",
		FillLevel:   FilledByScaffold,
		Description: "Instrucoes principais do agente — lido automaticamente pelo Claude Code",
		ScaffoldFills: []string{
			"Regras absolutas (nunca commitar sem testes, sempre push, etc)",
			"Principios ativos (TDD, Clean Arch, SOLID — baseado na sua escolha)",
			"Workflow padrao (le → implementa → testa → commit → push)",
			"Comandos disponiveis (/test /review /ship /commit)",
			"Stack do projeto e convencoes de codigo",
		},
		AgentFills: []string{
			"Decisoes tecnicas especificas descobertas durante o desenvolvimento",
			"Padroes do codebase que o agente deve seguir",
			"Gotchas e armadilhas especificas do projeto",
			"Notas sobre dependencias e integracoes",
		},
	},
	{
		Path:        "docs/CONTEXT.md",
		FillLevel:   FilledByScaffold,
		Description: "Briefing completo do projeto — contexto injetado em cada sessao do agente",
		ScaffoldFills: []string{
			"Nome, descricao e objetivo do projeto",
			"Stack tecnica completa",
			"Principios de desenvolvimento",
			"Estado inicial (scaffold criado, proximos passos)",
			"Autor e links do repo",
		},
		AgentFills: []string{
			"Decisoes arquiteturais tomadas durante o desenvolvimento",
			"Integracoes implementadas e como funcionam",
			"Estado atualizado das fases",
			"Aprendizados e descobertas do projeto",
		},
	},
	{
		Path:        "docs/ROADMAP.md",
		FillLevel:   FilledPartial,
		Description: "Fases e marcos do projeto",
		ScaffoldFills: []string{
			"Estrutura de fases (Fase 0 Setup, Fase 1 MVP, etc)",
			"Fase 0 marcada como concluida",
			"Proximas fases com itens genericos",
		},
		AgentFills: []string{
			"Itens especificos de cada fase baseados no projeto real",
			"Atualizacao de status conforme features sao implementadas",
			"Novas fases descobertas durante o desenvolvimento",
		},
	},
	{
		Path:        "docs/ARCHITECTURE.md",
		FillLevel:   FilledByAgent,
		Description: "Diagrama e descricao da arquitetura do sistema",
		ScaffoldFills: []string{
			"Secoes vazias com headers (Visao Geral, Diagrama, Fluxos)",
			"Nota: 'a definir conforme o desenvolvimento'",
		},
		AgentFills: []string{
			"Diagrama ASCII da arquitetura real apos primeiras implementacoes",
			"Descricao de cada componente e suas responsabilidades",
			"Fluxos de dados entre componentes",
			"Decisoes de design e trade-offs",
		},
	},
	{
		Path:        "docs/adr/0001-stack.md",
		FillLevel:   FilledPartial,
		Description: "Architecture Decision Record — decisao da stack tecnologica",
		ScaffoldFills: []string{
			"Status: Aceito",
			"Stack escolhida no wizard",
			"Data da decisao",
		},
		AgentFills: []string{
			"Contexto detalhado do porque da escolha",
			"Alternativas consideradas e por que foram descartadas",
			"Consequencias e trade-offs da decisao",
			"ADRs subsequentes para outras decisoes",
		},
	},
	{
		Path:        "docs/requirements/SRS.md",
		FillLevel:   FilledByAgent,
		Description: "Software Requirements Specification",
		ScaffoldFills: []string{
			"Estrutura de secoes (Introducao, Descricao, RF, RNF)",
			"Campos a preencher com placeholders",
		},
		AgentFills: []string{
			"Requisitos funcionais baseados no que foi implementado",
			"Requisitos nao-funcionais (performance, seguranca, escala)",
			"User stories derivadas das features reais",
		},
	},
	{
		Path:        ".claude/settings.json",
		FillLevel:   FilledByScaffold,
		Description: "Permissoes e configuracoes do Claude Code",
		ScaffoldFills: []string{
			"Lista de comandos permitidos (git, make, docker, npm...)",
			"Lista de comandos bloqueados",
			"Variaveis de ambiente do projeto",
		},
		AgentFills: []string{
			"Permissoes adicionais conforme o projeto necessita",
		},
	},
	{
		Path:        ".claude/hooks/",
		FillLevel:   FilledByScaffold,
		Description: "Executa antes de cada tool call do agente",
		ScaffoldFills: []string{
			"Bloqueio de comandos perigosos",
			"Log de tool calls para auditoria",
		},
		AgentFills: []string{
			"Validacoes especificas do projeto",
			"Checks de qualidade automaticos",
		},
	},
	{
		Path:        ".claude/commands/",
		FillLevel:   FilledByScaffold,
		Description: "Comandos customizados (/test /review /ship /commit)",
		ScaffoldFills: []string{
			"/test  — roda testes e reporta cobertura",
			"/review — code review do diff atual",
			"/ship  — testa + commit + push",
			"/commit — commit com conventional commits",
		},
		AgentFills: []string{
			"Comandos especificos do projeto conforme necessidade",
		},
	},
	{
		Path:        ".github/workflows/ci.yml",
		FillLevel:   FilledByScaffold,
		Description: "CI/CD — roda em cada push e PR",
		ScaffoldFills: []string{
			"Jobs de lint, test e build",
			"Setup do runtime correto (Go, Bun, Python...)",
			"Cache de dependencias",
		},
		AgentFills: []string{
			"Steps especificos conforme o projeto cresce",
			"Deploy automatico quando implementado",
			"Testes de integracao e e2e",
		},
	},
	{
		Path:        "CONTRIBUTING.md",
		FillLevel:   FilledByScaffold,
		Description: "Guia de contribuicao e convencoes",
		ScaffoldFills: []string{
			"Convencoes de commits (conventional commits)",
			"Processo de PR",
			"Como rodar localmente",
			"Regras de TDD se selecionado",
		},
		AgentFills: []string{
			"Detalhes especificos do setup do projeto",
			"Gotchas de desenvolvimento",
		},
	},
	{
		Path:        "SECURITY.md",
		FillLevel:   FilledPartial,
		Description: "Politica de seguranca e como reportar vulnerabilidades",
		ScaffoldFills: []string{
			"Contato para reportar vulnerabilidades",
			"SLA de resposta",
			"Escopo basico",
		},
		AgentFills: []string{
			"Escopo detalhado baseado nas features implementadas",
			"Threats especificas do projeto",
		},
	},
}
