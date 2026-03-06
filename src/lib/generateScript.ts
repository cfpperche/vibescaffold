import type { WizardState } from '../stores/wizard';

type Config = Pick<
  WizardState,
  | 'projectName'
  | 'projectDescription'
  | 'projectType'
  | 'frontend'
  | 'backend'
  | 'database'
  | 'principles'
  | 'claudeTools'
>;

const FRONTEND_DEPS: Record<string, string> = {
  react: 'react react-dom @vitejs/plugin-react vite typescript',
  vue: 'vue @vitejs/plugin-vue vite typescript',
  svelte: 'svelte @sveltejs/vite-plugin-svelte vite typescript',
  next: 'next react react-dom typescript @types/react',
  none: '',
};

const BACKEND_DEPS: Record<string, string> = {
  node: 'express @types/express tsx typescript',
  hono: 'hono @hono/node-server typescript tsx',
  fastify: 'fastify @fastify/cors tsx typescript',
  none: '',
};

const DB_DEPS: Record<string, string> = {
  postgres: 'pg @types/pg drizzle-orm drizzle-kit',
  sqlite: 'better-sqlite3 @types/better-sqlite3 drizzle-orm drizzle-kit',
  mysql: 'mysql2 drizzle-orm drizzle-kit',
  none: '',
};

const PRINCIPLE_LABELS: Record<string, string> = {
  'never-commit-without-build': 'NUNCA commite sem build passando',
  'never-expose-secrets': 'NUNCA exponha secrets (.env, keys, tokens)',
  'always-push-after-commit': 'SEMPRE git push apos git commit',
  'tests-before-merge': 'NUNCA merge PR sem testes passando',
  'small-commits': 'Commits pequenos e frequentes com mensagens claras',
  'no-console-log': 'Remover console.log antes de commitar',
  'type-safety': 'Manter strict TypeScript sem any',
};

function generateClaudeMd(config: Config): string {
  const rules = config.principles
    .map((r, i) => `${i + 1}. ${PRINCIPLE_LABELS[r] ?? r}`)
    .join('\n');

  const tools: string[] = [];
  if (config.claudeTools.includes('claude-md')) tools.push('- CLAUDE.md com regras do projeto');
  if (config.claudeTools.includes('context-docs')) tools.push('- docs/CONTEXT.md com contexto');
  if (config.claudeTools.includes('git-hooks')) tools.push('- Git hooks (pre-commit)');
  if (config.claudeTools.includes('github-actions')) tools.push('- GitHub Actions CI/CD');
  if (config.claudeTools.includes('biome')) tools.push('- Biome (lint + format)');
  if (config.claudeTools.includes('docker')) tools.push('- Dockerfile + docker-compose');

  return `# ${config.projectName}

## O que e
${config.projectDescription}

## Tipo
${config.projectType}

## Stack
- Frontend: ${config.frontend}
- Backend: ${config.backend}
- Database: ${config.database}

## Regras
${rules}

## Ferramentas
${tools.join('\n')}

## Workflow
Tarefa -> le arquivos -> implementa -> bun run build -> commit -> push
`;
}

function generateHooks(config: Config): string {
  if (!config.claudeTools.includes('git-hooks')) return '';
  return `
# Git hooks
cat > .git/hooks/pre-commit << 'HOOKEOF'
#!/bin/sh
set -e
bunx biome check src/
tsc --noEmit
HOOKEOF
chmod +x .git/hooks/pre-commit
`;
}

function generateCicd(config: Config): string {
  if (!config.claudeTools.includes('github-actions')) return '';
  return `
# GitHub Actions
mkdir -p .github/workflows
cat > .github/workflows/ci.yml << 'CIEOF'
name: CI
on:
  push: { branches: [main] }
  pull_request: { branches: [main] }
jobs:
  ci:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: oven-sh/setup-bun@v2
        with: { bun-version: latest }
      - run: bun install --frozen-lockfile
      - run: bunx biome check .
      - run: tsc --noEmit
      - run: bun run build
CIEOF
`;
}

function generateDocker(config: Config): string {
  if (!config.claudeTools.includes('docker')) return '';
  return `
# Docker
cat > Dockerfile << 'DOCKEREOF'
FROM oven/bun:1 AS base
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile
COPY . .
RUN bun run build
EXPOSE 3000
CMD ["bun", "run", "start"]
DOCKEREOF

cat > docker-compose.yml << 'COMPOSEEOF'
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
COMPOSEEOF
`;
}

function generateBiome(config: Config): string {
  if (!config.claudeTools.includes('biome')) return '';
  return `
# Biome config
cat > biome.json << 'BIOMEEOF'
{
  "$schema": "https://biomejs.dev/schemas/2.0.0/schema.json",
  "linter": { "enabled": true, "rules": { "recommended": true } },
  "formatter": { "enabled": true, "indentStyle": "space", "indentWidth": 2 },
  "javascript": { "formatter": { "quoteStyle": "single", "semicolons": "always" } }
}
BIOMEEOF
bun add -d @biomejs/biome
`;
}

export function generateScript(config: Config): string {
  const frontendDeps = FRONTEND_DEPS[config.frontend] ?? '';
  const backendDeps = BACKEND_DEPS[config.backend] ?? '';
  const dbDeps = DB_DEPS[config.database] ?? '';
  const claudeMd = generateClaudeMd(config);
  const allDeps = [frontendDeps, backendDeps, dbDeps].filter(Boolean).join(' ');

  return `#!/bin/bash
set -e

PROJECT_NAME="${config.projectName}"

echo ">>> Creating project: $PROJECT_NAME"
mkdir -p "$PROJECT_NAME"
cd "$PROJECT_NAME"

# Initialize
git init
bun init -y

# Install dependencies
${allDeps ? `bun add ${allDeps}` : '# No dependencies selected'}
bun add -d typescript @types/node

# CLAUDE.md
cat > CLAUDE.md << 'CLAUDEEOF'
${claudeMd}CLAUDEEOF

# Directory structure
mkdir -p src/components src/lib src/stores docs

# docs/CONTEXT.md
${config.claudeTools.includes('context-docs') ? `cat > docs/CONTEXT.md << 'DOCSEOF'
# ${config.projectName} — Contexto

## Produto
${config.projectDescription}

## Tipo
${config.projectType}

## Stack
- Frontend: ${config.frontend}
- Backend: ${config.backend}
- Database: ${config.database}
DOCSEOF` : '# docs/CONTEXT.md skipped'}

# .gitignore
cat > .gitignore << 'GIEOF'
node_modules/
dist/
.env
.env.local
*.log
.DS_Store
.wrangler/
GIEOF
${generateBiome(config)}${generateHooks(config)}${generateCicd(config)}${generateDocker(config)}
# Initial commit
git add -A
git commit -m "Initial scaffold by VibeScaffold"

echo ""
echo ">>> Project '$PROJECT_NAME' scaffolded!"
echo ">>> cd $PROJECT_NAME && claude"
`;
}
