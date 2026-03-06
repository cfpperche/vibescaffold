# VibeScaffold — Claude Code

## O que é
Plataforma web que gera scaffold de projetos para vibecoding.
Wizard 5 passos → bash script com CLAUDE.md, hooks, docs, CI/CD.

## Stack
Vite + React 19 + TypeScript + TanStack Router + Zustand + Tailwind v4
Hono.js em Cloudflare Workers para server-side

## Regras
1. NUNCA commite sem build passando (npm run build)
2. NUNCA exponha secrets
3. SEMPRE git push após git commit
4. Componentes em src/components/
5. Estado global em src/stores/ (Zustand)
6. Rotas em src/routes/ (TanStack Router file-based)
7. Lógica de negócio em src/lib/

## Workflow
Tarefa → lê arquivos → implementa → npm run build → commit → push

## Referência de design
/mnt/user-data/outputs/vibescaffold-landing.jsx
