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

type Category struct {
	Name  string
	Files []string
}

var Categories = []Category{
	{"Produto", []string{"docs/PRODUCT_BRIEF.md", "docs/PERSONA.md", "docs/VIRAL_LOOP.md"}},
	{"Raiz do projeto", []string{"README.md", "CHANGELOG.md", "LICENSE", "CODE_OF_CONDUCT.md", ".editorconfig"}},
	{"Claude Code", []string{"CLAUDE.md", ".claude/settings.json", ".claude/hooks/", ".claude/commands/"}},
	{"Documentacao", []string{"docs/CONTEXT.md", "docs/ROADMAP.md", "docs/ARCHITECTURE.md", "docs/GLOSSARY.md", "docs/TESTING.md", "docs/DEPLOYMENT.md"}},
	{"Requisitos", []string{"docs/requirements/SRS.md", "docs/requirements/RF.md", "docs/requirements/RNF.md", "docs/requirements/USER_STORIES.md", "docs/requirements/USE_CASES.md"}},
	{"Decisoes", []string{"docs/adr/0001-stack.md"}},
	{"GitHub", []string{".github/workflows/ci.yml", ".github/workflows/release.yml", ".github/dependabot.yml", ".github/PULL_REQUEST_TEMPLATE.md", ".github/ISSUE_TEMPLATE/bug_report.md", ".github/ISSUE_TEMPLATE/feature_request.md"}},
	{"Qualidade", []string{".pre-commit-config.yaml", "CONTRIBUTING.md", "SECURITY.md"}},
	{"Scripts", []string{"scripts/setup.sh"}},
}

// FileMap provides fast lookup by path.
var FileMap map[string]MDFile

func init() {
	FileMap = make(map[string]MDFile, len(Files))
	for _, f := range Files {
		FileMap[f.Path] = f
	}
}

var Files = []MDFile{
	// --- Produto ---
	{
		Path:        "docs/PRODUCT_BRIEF.md",
		FillLevel:   FilledByScaffold,
		Description: "Brief completo do produto — hook, loop, share trigger, viral loop, monetizacao",
		ScaffoldFills: []string{
			"Hook: por que alguem tenta na primeira vez",
			"Loop: Day 1 / Day 7 / Day 30 / Day 90",
			"Share Trigger: momento exato que gera compartilhamento",
			"Mecanicas de engajamento selecionadas",
			"Modelo de monetizacao",
			"Teste 'mostre pra sua mae'",
			"Riscos e mitigacoes",
		},
		AgentFills: []string{
			"Refinamento das mecanicas apos validacao com usuarios reais",
			"Atualizacao do loop conforme metricas reais de retencao",
			"Novos riscos descobertos durante o desenvolvimento",
			"Ajustes de monetizacao baseados em feedback",
		},
	},
	{
		Path:        "docs/PERSONA.md",
		FillLevel:   FilledPartial,
		Description: "Personas do produto — quem usa, como, e por que",
		ScaffoldFills: []string{
			"Publico primario com comportamentos e motivacoes",
			"Publico crossover",
			"Anti-persona (quem NAO e o usuario)",
			"Jobs-to-be-done por persona",
		},
		AgentFills: []string{
			"Personas refinadas apos entrevistas com usuarios reais",
			"Novos segmentos descobertos organicamente",
			"Comportamentos inesperados observados",
		},
	},
	{
		Path:        "docs/VIRAL_LOOP.md",
		FillLevel:   FilledPartial,
		Description: "Diagrama e analise do loop viral do produto",
		ScaffoldFills: []string{
			"Diagrama ASCII do loop viral",
			"Coeficiente viral estimado (low/medium/high)",
			"Formatos de share planejados (tweet, story, embed)",
			"Anti-patterns a evitar",
		},
		AgentFills: []string{
			"K-factor real medido apos launch",
			"Otimizacoes do loop baseadas em dados",
			"Novos vetores de share descobertos",
		},
	},

	// --- Raiz do projeto ---
	{
		Path:        "README.md",
		FillLevel:   FilledByScaffold,
		Description: "Porta de entrada do projeto — primeiro arquivo lido por qualquer pessoa ou agente",
		ScaffoldFills: []string{
			"Nome, descricao e tagline do projeto",
			"Tabela da stack tecnica",
			"Instrucoes de setup (bash scripts/setup.sh)",
			"Tabela de portas dos servicos",
			"Comandos uteis (make up, make db, etc)",
			"Tabela do roadmap resumido",
			"Licenca",
		},
		AgentFills: []string{
			"Badges de CI, cobertura e versao conforme configurados",
			"Screenshots e demos quando o projeto tiver interface",
			"Exemplos de uso da API quando implementada",
			"Documentacao de variaveis de ambiente reais",
		},
	},
	{
		Path:        "CHANGELOG.md",
		FillLevel:   FilledPartial,
		Description: "Historico de versoes seguindo Keep a Changelog + Semantic Versioning",
		ScaffoldFills: []string{
			"Secao Unreleased com itens do scaffold inicial",
			"Versao 0.1.0-alpha com data",
			"Estrutura de secoes (Added, Changed, Fixed, Removed)",
		},
		AgentFills: []string{
			"Entrada a cada release com o que foi adicionado",
			"Breaking changes documentados",
			"Migracao de versoes quando necessario",
		},
	},
	{
		Path:        "LICENSE",
		FillLevel:   FilledByScaffold,
		Description: "Licenca do projeto",
		ScaffoldFills: []string{
			"Licenca proprietaria ou MIT baseado na escolha no wizard",
			"Copyright com nome do autor e ano",
		},
		AgentFills: []string{},
	},
	{
		Path:        "CODE_OF_CONDUCT.md",
		FillLevel:   FilledByScaffold,
		Description: "Codigo de conduta para contribuidores",
		ScaffoldFills: []string{
			"Contributor Covenant 2.1",
			"Contato para reporte de violacoes",
		},
		AgentFills: []string{},
	},
	{
		Path:        ".editorconfig",
		FillLevel:   FilledByScaffold,
		Description: "Configuracao de indentacao e encoding para todos os editores e agentes",
		ScaffoldFills: []string{
			"Regras globais: LF, UTF-8, trim whitespace",
			"Regras por linguagem: Go (tabs), TS/JS (2 spaces), Python/C++ (4 spaces)",
			"Makefile e scripts: tabs obrigatorio",
		},
		AgentFills: []string{
			"Regras adicionais conforme novas linguagens forem adicionadas",
		},
	},

	// --- Claude Code ---
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

	// --- Documentacao ---
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
		Path:        "docs/GLOSSARY.md",
		FillLevel:   FilledPartial,
		Description: "Glossario de termos tecnicos e de dominio do projeto",
		ScaffoldFills: []string{
			"Estrutura do arquivo com secoes por letra",
			"Termos tecnicos da stack escolhida",
		},
		AgentFills: []string{
			"Termos de dominio especificos do negocio",
			"Acronimos e siglas usados no codebase",
			"Definicoes de conceitos arquiteturais do projeto",
			"Linguagem ubiqua do DDD se selecionado",
		},
	},
	{
		Path:        "docs/TESTING.md",
		FillLevel:   FilledPartial,
		Description: "Estrategia de testes do projeto",
		ScaffoldFills: []string{
			"Piramide de testes: unit, integration, e2e",
			"Cobertura minima configurada (80%)",
			"Comandos para rodar cada tipo de teste",
			"Test runner baseado na stack escolhida",
		},
		AgentFills: []string{
			"Estrategias especificas para cada modulo",
			"Fixtures e mocks documentados",
			"Testes de carga e performance quando implementados",
			"Como testar componentes criticos",
		},
	},
	{
		Path:        "docs/DEPLOYMENT.md",
		FillLevel:   FilledPartial,
		Description: "Guia de deploy em producao",
		ScaffoldFills: []string{
			"Requisitos de servidor",
			"Variaveis de ambiente obrigatorias",
			"Comandos de deploy",
			"Como fazer rollback",
		},
		AgentFills: []string{
			"Configuracao real da infraestrutura quando provisionada",
			"IaC (Terraform/Pulumi) documentado",
			"Monitoramento e alertas quando configurados",
			"Runbooks de incidentes",
		},
	},

	// --- Requisitos ---
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
		Path:        "docs/requirements/RF.md",
		FillLevel:   FilledPartial,
		Description: "Requisitos Funcionais detalhados — o que o sistema deve fazer",
		ScaffoldFills: []string{
			"Template padrao por requisito (prioridade, status, descricao, entrada, saida)",
			"RF-001 placeholder para o primeiro requisito",
		},
		AgentFills: []string{
			"Requisitos funcionais completos baseados nas features implementadas",
			"Status atualizado (Planejado → Em desenvolvimento → Concluido)",
			"Regras de negocio especificas descobertas",
			"Relacionamentos entre requisitos",
		},
	},
	{
		Path:        "docs/requirements/RNF.md",
		FillLevel:   FilledPartial,
		Description: "Requisitos Nao-Funcionais — como o sistema opera",
		ScaffoldFills: []string{
			"Secoes: Performance, Seguranca, Escalabilidade, Confiabilidade",
			"Valores padrao (p95 < 200ms, uptime 99.5%, etc)",
		},
		AgentFills: []string{
			"Metricas reais baseadas em testes de carga",
			"Requisitos de seguranca especificos do dominio",
			"SLAs definidos com o produto",
		},
	},
	{
		Path:        "docs/requirements/USER_STORIES.md",
		FillLevel:   FilledByAgent,
		Description: "Historias de usuario na perspectiva do usuario final",
		ScaffoldFills: []string{
			"Template: Como [usuario], quero [acao], para [beneficio]",
			"Estrutura de criterios de aceitacao (Given/When/Then)",
			"Personas identificadas no wizard",
		},
		AgentFills: []string{
			"User stories completas baseadas nas features implementadas",
			"Criterios de aceitacao detalhados",
			"Historias de edge cases descobertos",
			"Historias de usuarios admin e moderador",
		},
	},
	{
		Path:        "docs/requirements/USE_CASES.md",
		FillLevel:   FilledByAgent,
		Description: "Casos de uso — interacoes entre atores e o sistema",
		ScaffoldFills: []string{
			"Template padrao por caso de uso",
			"Lista de atores identificados",
		},
		AgentFills: []string{
			"UC-001 a UC-N completos conforme features implementadas",
			"Fluxos alternativos e excecoes",
			"Diagramas textuais de sequencia",
		},
	},

	// --- Decisoes ---
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

	// --- GitHub ---
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
		Path:        ".github/workflows/release.yml",
		FillLevel:   FilledByScaffold,
		Description: "Release automatico quando tag v*.*.* e criada",
		ScaffoldFills: []string{
			"Trigger em tags v*.*.*",
			"Geracao automatica de release notes do CHANGELOG",
			"GitHub Release criado automaticamente",
		},
		AgentFills: []string{
			"Build de binarios para multiplas plataformas quando necessario",
			"Upload de assets da release",
			"Notificacoes de release",
		},
	},
	{
		Path:        ".github/dependabot.yml",
		FillLevel:   FilledByScaffold,
		Description: "Atualizacoes automaticas de dependencias",
		ScaffoldFills: []string{
			"Schedule semanal para Go modules, npm, GitHub Actions",
			"Limite de PRs abertos simultaneos",
		},
		AgentFills: []string{},
	},
	{
		Path:        ".github/PULL_REQUEST_TEMPLATE.md",
		FillLevel:   FilledByScaffold,
		Description: "Template padrao para todos os PRs do projeto",
		ScaffoldFills: []string{
			"Tipo de mudanca (feat/fix/refactor/docs/chore)",
			"Secao de descricao",
			"Como testar",
			"Checklist: CI verde, testes, CHANGELOG, sem secrets",
		},
		AgentFills: []string{
			"Itens de checklist especificos do projeto conforme evolui",
		},
	},
	{
		Path:        ".github/ISSUE_TEMPLATE/bug_report.md",
		FillLevel:   FilledByScaffold,
		Description: "Template para reportar bugs",
		ScaffoldFills: []string{
			"Descricao, passos para reproduzir",
			"Comportamento esperado vs atual",
			"Logs e ambiente (OS, GPU, versao)",
		},
		AgentFills: []string{},
	},
	{
		Path:        ".github/ISSUE_TEMPLATE/feature_request.md",
		FillLevel:   FilledByScaffold,
		Description: "Template para solicitar novas features",
		ScaffoldFills: []string{
			"Problema que resolve",
			"Solucao proposta",
			"Fase do roadmap que pertence",
		},
		AgentFills: []string{},
	},

	// --- Qualidade ---
	{
		Path:        ".pre-commit-config.yaml",
		FillLevel:   FilledByScaffold,
		Description: "Hooks que rodam antes de cada commit garantindo qualidade",
		ScaffoldFills: []string{
			"trailing-whitespace, end-of-file-fixer, check-yaml",
			"detect-private-key — nunca commita secrets",
			"Lint automatico da linguagem principal",
		},
		AgentFills: []string{
			"Hooks especificos do projeto conforme crescem",
			"Validacoes de schema de banco",
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

	// --- Scripts ---
	{
		Path:        "scripts/setup.sh",
		FillLevel:   FilledByScaffold,
		Description: "Setup do ambiente do zero em um comando",
		ScaffoldFills: []string{
			"Copia .env.example → .env",
			"Instala pre-commit hooks",
			"Sobe docker compose se presente",
			"Mensagem final com proximos passos",
		},
		AgentFills: []string{
			"Instalacao de dependencias especificas descobertas",
			"Configuracoes de ambiente adicionais",
			"Verificacoes de saude do ambiente",
		},
	},
}
