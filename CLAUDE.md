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

## Estrutura
- cmd/vs/main.go — entry point
- internal/tui/ — modelo principal + views + components + styles
- internal/tui/views/agent.go — seletor de agente LLM
- internal/agent/detector.go — detecta binarios instalados (claude, codex, gemini, ollama, aider)
- internal/agent/launcher.go — lanca agente com contexto injetado
- internal/scaffold/ — logica de geracao de projetos
- internal/doctor/ — health check
- internal/config/ — config persistida em ~/.vibescaffold/ + deteccao de projeto
