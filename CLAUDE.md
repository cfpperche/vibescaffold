# VibeScaffold — Claude Code

## O que e
TUI em Go com Bubble Tea para scaffold de projetos de vibecoding.
Gerenciador de contexto para Claude Code, Codex CLI, Gemini CLI.

## Stack
Go + Bubble Tea + Bubbles + Lip Gloss + Huh + Cobra

## Comandos
- make dev        → roda a TUI
- make build      → compila para dist/vs
- make install    → instala globalmente
- make test       → roda todos os testes (teatest + unit)
- make test-v     → testes verbose
- make test-update → atualiza golden files
- make demo       → grava demo GIF com VHS
- make demo-quick → grava demo curta com VHS

## Testes
- teatest: testes interativos da TUI (internal/tui/app_test.go)
- unit: testes de doctor e scaffold
- golden files: testdata/*.golden (atualizar com -update)
- VHS tapes: demos/*.tape (geram GIFs para documentacao)

## Regras
1. NUNCA commite sem go build passando
2. Cada view e um sub-modelo Bubble Tea independente
3. Estilos SEMPRE em internal/tui/styles/theme.go
4. Logica de negocio NUNCA no pacote tui — vai em internal/scaffold ou internal/doctor
5. SEMPRE git push apos commit

## Arquitetura Bubble Tea
msg → Update() → cmd → View()
Cada view implementa: Init() / Update() / View()

## Fluxo principal
vs new  → wizard produto (6 rounds) → gera PRODUCT_BRIEF.md + PERSONA.md + VIRAL_LOOP.md
       → scaffold tecnico (le brief, adapta todos os docs)
       → chat com agente (contexto completo: produto + tecnico)

vs init → scaffold tecnico direto (sem wizard de produto)

## Pacotes de produto
- internal/product/brief/brief.go      → estrutura de dados do brief
- internal/product/brief/generator.go  → gera os arquivos MD
- internal/product/wizard/wizard.go    → wizard data + conversao para Brief
- internal/tui/views/product.go        → wizard interativo com huh (6 rounds)

## Arquitetura do Chat

O chat e o coracao do produto. Apos init, o terminal vira
um ambiente persistente que vive durante todo o desenvolvimento.

Fluxo: input do usuario → IsCommand? → HandleCommand()
                                     → Send() → RunAgent() → stream → viewport

Pacotes:
- internal/chat/session.go   → estado da sessao
- internal/chat/runner.go    → executa agente como processo filho ou API (Ollama)
- internal/chat/context.go   → build e injecao de contexto
- internal/chat/commands.go  → /switch /doctor /status /context /clear /exit /help
- internal/tui/views/chat.go → UI do chat com streaming

Deteccao automatica: se CLAUDE.md existe no cwd, vs abre chat direto.

## Estrutura
- cmd/vs/main.go — entry point (detecta projeto → chat ou home)
- internal/tui/ — modelo principal + views + components + styles
- internal/tui/views/product.go — wizard de produto (vs new)
- internal/tui/views/chat.go — chat persistente com streaming
- internal/tui/views/agent.go — seletor de agente LLM
- internal/chat/ — sessao, runner, contexto, comandos
- internal/onboarding/ — dados dos arquivos MD, controle de first-run
- internal/tui/views/onboarding.go — tela de onboarding navegavel
- internal/agent/detector.go — detecta binarios instalados (claude, codex, gemini, ollama, aider)
- internal/agent/launcher.go — lanca agente com contexto injetado
- internal/scaffold/ — logica de geracao de projetos
- internal/doctor/ — health check
- internal/config/ — config persistida em ~/.vibescaffold/ + deteccao de projeto
