import type { WizardState } from '../stores/wizard';

type Config = Pick<
  WizardState,
  | 'projectName'
  | 'projectDescription'
  | 'frontend'
  | 'backend'
  | 'database'
  | 'claudeMdRules'
  | 'hooks'
  | 'cicd'
  | 'testing'
>;

const FRONTEND_DEPS: Record<string, string> = {
  react: 'react react-dom @vitejs/plugin-react vite typescript',
  vue: 'vue @vitejs/plugin-vue vite typescript',
  svelte: 'svelte @sveltejs/vite-plugin-svelte vite typescript',
};

const BACKEND_DEPS: Record<string, string> = {
  node: 'express @types/express tsx typescript',
  hono: 'hono @hono/node-server typescript tsx',
  fastify: 'fastify @fastify/cors tsx typescript',
};

const DB_DEPS: Record<string, string> = {
  postgres: 'pg @types/pg drizzle-orm drizzle-kit',
  sqlite: 'better-sqlite3 @types/better-sqlite3 drizzle-orm drizzle-kit',
  mysql: 'mysql2 drizzle-orm drizzle-kit',
};

function generateClaudeMd(config: Config): string {
  const rules = config.claudeMdRules
    .map((r) => {
      switch (r) {
        case 'never-commit-without-build':
          return '1. NUNCA commite sem build passando';
        case 'never-expose-secrets':
          return '2. NUNCA exponha secrets';
        case 'always-push-after-commit':
          return '3. SEMPRE git push apos git commit';
        default:
          return `- ${r}`;
      }
    })
    .join('\n');

  return `# ${config.projectName}

## O que e
${config.projectDescription}

## Stack
- Frontend: ${config.frontend}
- Backend: ${config.backend}
- Database: ${config.database}
- Testing: ${config.testing}
- CI/CD: ${config.cicd}

## Regras
${rules}

## Workflow
Tarefa -> le arquivos -> implementa -> npm run build -> commit -> push
`;
}

function generateHooks(config: Config): string {
  const hookLines: string[] = [];
  for (const hook of config.hooks) {
    if (hook === 'pre-commit-lint') {
      hookLines.push('  npx eslint --max-warnings 0 src/');
    }
    if (hook === 'pre-commit-typecheck') {
      hookLines.push('  npx tsc --noEmit');
    }
  }
  if (hookLines.length === 0) return '';
  return `
cat > .git/hooks/pre-commit << 'HOOKEOF'
#!/bin/sh
set -e
${hookLines.join('\n')}
HOOKEOF
chmod +x .git/hooks/pre-commit
`;
}

function generateCicd(config: Config): string {
  if (config.cicd !== 'github-actions') return '';
  return `
mkdir -p .github/workflows
cat > .github/workflows/ci.yml << 'CIEOF'
name: CI
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm
      - run: npm ci
      - run: npm run build
      - run: npm test
CIEOF
`;
}

export function generateScript(config: Config): string {
  const frontendDeps = FRONTEND_DEPS[config.frontend] ?? '';
  const backendDeps = BACKEND_DEPS[config.backend] ?? '';
  const dbDeps = DB_DEPS[config.database] ?? '';
  const claudeMd = generateClaudeMd(config);

  return `#!/bin/bash
set -e

PROJECT_NAME="${config.projectName}"

echo ">>> Creating project: $PROJECT_NAME"
mkdir -p "$PROJECT_NAME"
cd "$PROJECT_NAME"

# Initialize git & npm
git init
npm init -y

# Install dependencies
npm install ${frontendDeps} ${backendDeps} ${dbDeps}
npm install -D ${config.testing} eslint @types/node

# Create CLAUDE.md
cat > CLAUDE.md << 'CLAUDEEOF'
${claudeMd}CLAUDEEOF

# Create directory structure
mkdir -p src/components src/lib src/stores docs

# Create docs/CONTEXT.md
cat > docs/CONTEXT.md << 'DOCSEOF'
# ${config.projectName} — Contexto

## Produto
${config.projectDescription}

## Stack
- Frontend: ${config.frontend}
- Backend: ${config.backend}
- Database: ${config.database}
- Testing: ${config.testing}
DOCSEOF

# Create .gitignore
cat > .gitignore << 'GIEOF'
node_modules/
dist/
.env
.env.local
*.log
.DS_Store
GIEOF
${generateHooks(config)}${generateCicd(config)}
# Initial commit
git add -A
git commit -m "Initial scaffold by VibeScaffold"

echo ""
echo ">>> Project '$PROJECT_NAME' scaffolded successfully!"
echo ">>> cd $PROJECT_NAME && code ."
`;
}
