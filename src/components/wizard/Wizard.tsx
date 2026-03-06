import { useState } from 'react';
import { generateScript } from '../../lib/generateScript';
import { type ProjectType, useWizardStore } from '../../stores/wizard';
import { Link } from '@tanstack/react-router';

const STEP_TITLES = ['Projeto', 'Tipo & Stack', 'Principios', 'Ferramentas Claude', 'Script'];

function StepIndicator({ current }: { current: number }) {
  return (
    <div className="mb-8 flex items-center justify-center gap-1">
      {STEP_TITLES.map((title, i) => (
        <div key={title} className="flex items-center gap-1">
          <button
            type="button"
            onClick={() => {
              if (i < current) useWizardStore.getState().setStep(i);
            }}
            className={`flex h-8 w-8 items-center justify-center rounded-full text-xs font-bold transition
              ${i === current ? 'bg-neon text-black' : i < current ? 'bg-neon/20 text-neon cursor-pointer hover:bg-neon/30' : 'bg-surface-lighter text-neutral-600'}`}
          >
            {i < current ? '\u2713' : i + 1}
          </button>
          {i < STEP_TITLES.length - 1 && (
            <div className={`h-px w-6 ${i < current ? 'bg-neon/40' : 'bg-surface-lighter'}`} />
          )}
        </div>
      ))}
    </div>
  );
}

function StepProject() {
  const { projectName, projectDescription, setField } = useWizardStore();
  return (
    <div className="space-y-5">
      <div>
        <label htmlFor="pname" className="mb-2 block text-xs font-semibold uppercase tracking-wider text-neon/70">
          Nome do Projeto
        </label>
        <input
          id="pname"
          value={projectName}
          onChange={(e) => setField('projectName', e.target.value)}
          placeholder="meu-projeto"
          className="w-full rounded-lg border border-surface-border bg-surface px-4 py-3 text-sm text-white placeholder-neutral-600 outline-none transition focus:border-neon/50 focus:shadow-[0_0_12px_#39ff1420]"
        />
      </div>
      <div>
        <label htmlFor="pdesc" className="mb-2 block text-xs font-semibold uppercase tracking-wider text-neon/70">
          Descricao
        </label>
        <textarea
          id="pdesc"
          value={projectDescription}
          onChange={(e) => setField('projectDescription', e.target.value)}
          placeholder="Descreva seu projeto..."
          rows={3}
          className="w-full rounded-lg border border-surface-border bg-surface px-4 py-3 text-sm text-white placeholder-neutral-600 outline-none transition focus:border-neon/50 focus:shadow-[0_0_12px_#39ff1420]"
        />
      </div>
    </div>
  );
}

const PROJECT_TYPES: { id: ProjectType; label: string; icon: string }[] = [
  { id: 'saas', label: 'SaaS', icon: '>' },
  { id: 'api', label: 'API', icon: '{' },
  { id: 'game', label: 'Game', icon: '#' },
  { id: 'cli', label: 'CLI', icon: '$' },
  { id: 'mobile', label: 'Mobile', icon: '@' },
  { id: 'lib', label: 'Library', icon: '*' },
];

const STACK_OPTIONS: Record<string, { label: string; options: string[] }> = {
  frontend: { label: 'Frontend', options: ['react', 'vue', 'svelte', 'next', 'none'] },
  backend: { label: 'Backend', options: ['node', 'hono', 'fastify', 'none'] },
  database: { label: 'Database', options: ['postgres', 'sqlite', 'mysql', 'none'] },
};

function StepTypeAndStack() {
  const { projectType, frontend, backend, database, setField } = useWizardStore();
  return (
    <div className="space-y-6">
      <div>
        <span className="mb-3 block text-xs font-semibold uppercase tracking-wider text-neon/70">
          Tipo de Projeto
        </span>
        <div className="grid grid-cols-3 gap-2">
          {PROJECT_TYPES.map((t) => (
            <button
              key={t.id}
              type="button"
              onClick={() => setField('projectType', t.id)}
              className={`rounded-lg border px-3 py-3 text-left text-sm transition
                ${projectType === t.id
                  ? 'border-neon bg-neon/10 text-neon shadow-[0_0_12px_#39ff1415]'
                  : 'border-surface-border bg-surface text-neutral-400 hover:border-neutral-600'}`}
            >
              <span className="mr-2 text-xs opacity-50">{t.icon}</span>
              {t.label}
            </button>
          ))}
        </div>
      </div>

      <div className="border-t border-surface-border pt-5">
        <span className="mb-3 block text-xs font-semibold uppercase tracking-wider text-neon/70">Stack</span>
        <div className="space-y-4">
          {Object.entries(STACK_OPTIONS).map(([key, { label, options }]) => {
            const current = key === 'frontend' ? frontend : key === 'backend' ? backend : database;
            return (
              <div key={key}>
                <span className="mb-2 block text-xs text-neutral-500">{label}</span>
                <div className="flex flex-wrap gap-2">
                  {options.map((opt) => (
                    <button
                      key={opt}
                      type="button"
                      onClick={() => setField(key as 'frontend' | 'backend' | 'database', opt)}
                      className={`rounded-md border px-3 py-1.5 text-xs font-medium capitalize transition
                        ${current === opt
                          ? 'border-neon/50 bg-neon/10 text-neon'
                          : 'border-surface-border text-neutral-500 hover:border-neutral-600 hover:text-neutral-300'}`}
                    >
                      {opt}
                    </button>
                  ))}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

const PRINCIPLE_OPTIONS = [
  { id: 'never-commit-without-build', label: 'Nunca commitar sem build passando' },
  { id: 'never-expose-secrets', label: 'Nunca expor secrets (.env, keys)' },
  { id: 'always-push-after-commit', label: 'Sempre push apos commit' },
  { id: 'tests-before-merge', label: 'Nunca merge sem testes' },
  { id: 'small-commits', label: 'Commits pequenos e frequentes' },
  { id: 'no-console-log', label: 'Sem console.log em producao' },
  { id: 'type-safety', label: 'Strict TypeScript (sem any)' },
];

function ToggleList({
  label,
  selected,
  options,
  onChange,
}: {
  label: string;
  selected: string[];
  options: { id: string; label: string }[];
  onChange: (v: string[]) => void;
}) {
  const toggle = (id: string) => {
    onChange(selected.includes(id) ? selected.filter((s) => s !== id) : [...selected, id]);
  };
  return (
    <div>
      <span className="mb-3 block text-xs font-semibold uppercase tracking-wider text-neon/70">{label}</span>
      <div className="space-y-1.5">
        {options.map((opt) => {
          const active = selected.includes(opt.id);
          return (
            <button
              key={opt.id}
              type="button"
              onClick={() => toggle(opt.id)}
              className={`flex w-full items-center gap-3 rounded-lg border px-4 py-2.5 text-left text-sm transition
                ${active
                  ? 'border-neon/30 bg-neon/5 text-neon'
                  : 'border-surface-border text-neutral-500 hover:border-neutral-600'}`}
            >
              <span
                className={`flex h-4 w-4 shrink-0 items-center justify-center rounded-sm border text-[10px] transition
                  ${active ? 'border-neon bg-neon text-black' : 'border-neutral-600'}`}
              >
                {active ? '\u2713' : ''}
              </span>
              {opt.label}
            </button>
          );
        })}
      </div>
    </div>
  );
}

function StepPrinciples() {
  const { principles, setField } = useWizardStore();
  return (
    <ToggleList
      label="Principios do Projeto"
      selected={principles}
      options={PRINCIPLE_OPTIONS}
      onChange={(v) => setField('principles', v)}
    />
  );
}

const CLAUDE_TOOL_OPTIONS = [
  { id: 'claude-md', label: 'CLAUDE.md — regras para o AI' },
  { id: 'context-docs', label: 'docs/CONTEXT.md — contexto do projeto' },
  { id: 'git-hooks', label: 'Git Hooks — pre-commit com lint + typecheck' },
  { id: 'github-actions', label: 'GitHub Actions — CI/CD pipeline' },
  { id: 'biome', label: 'Biome — lint + format config' },
  { id: 'docker', label: 'Docker — Dockerfile + compose' },
];

function StepClaudeTools() {
  const { claudeTools, setField } = useWizardStore();
  return (
    <ToggleList
      label="Ferramentas Claude Code"
      selected={claudeTools}
      options={CLAUDE_TOOL_OPTIONS}
      onChange={(v) => setField('claudeTools', v)}
    />
  );
}

function StepScript() {
  const store = useWizardStore();
  const [copied, setCopied] = useState(false);
  const script = generateScript(store);

  const copy = () => {
    navigator.clipboard.writeText(script);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <span className="text-xs font-semibold uppercase tracking-wider text-neon/70">
          Script Gerado
        </span>
        <button
          type="button"
          onClick={copy}
          className={`rounded-md border px-4 py-1.5 text-xs font-semibold transition
            ${copied
              ? 'border-neon bg-neon text-black'
              : 'border-neon/50 text-neon hover:bg-neon/10'}`}
        >
          {copied ? 'Copiado!' : 'Copiar'}
        </button>
      </div>
      <pre className="max-h-[400px] overflow-auto rounded-lg border border-surface-border bg-black p-4 text-xs leading-relaxed text-neon/80">
        {script}
      </pre>
      <p className="mt-4 text-xs text-neutral-600">
        Cole no terminal: <code className="text-neon/50">bash scaffold.sh</code>
      </p>
    </div>
  );
}

const STEPS = [StepProject, StepTypeAndStack, StepPrinciples, StepClaudeTools, StepScript];

export function Wizard() {
  const { step, nextStep, prevStep, projectName } = useWizardStore();
  const CurrentStep = STEPS[step];
  const canNext = step === 0 ? projectName.trim().length > 0 : true;

  return (
    <div className="mx-auto min-h-screen max-w-2xl px-6 py-10">
      <div className="mb-6 flex items-center justify-between">
        <Link to="/" className="text-xs text-neutral-600 transition hover:text-neon/50">
          &lt;- voltar
        </Link>
        <span className="text-xs text-neutral-600">
          {step + 1}/{STEP_TITLES.length}
        </span>
      </div>

      <h1 className="mb-1 text-2xl font-bold text-white">
        <span className="text-neon">$</span> {STEP_TITLES[step]}
      </h1>
      <p className="mb-6 text-xs text-neutral-600">passo {step + 1} de {STEP_TITLES.length}</p>

      <StepIndicator current={step} />

      <div className="rounded-xl border border-surface-border bg-surface-light p-6">
        <CurrentStep />
      </div>

      <div className="mt-6 flex justify-between">
        <button
          type="button"
          onClick={prevStep}
          disabled={step === 0}
          className="rounded-lg border border-surface-border px-5 py-2.5 text-sm text-neutral-500 transition hover:border-neutral-600 hover:text-neutral-300 disabled:opacity-20 disabled:hover:border-surface-border disabled:hover:text-neutral-500"
        >
          Voltar
        </button>
        {step < 4 && (
          <button
            type="button"
            onClick={nextStep}
            disabled={!canNext}
            className="rounded-lg border border-neon/50 bg-neon/10 px-5 py-2.5 text-sm font-semibold text-neon transition hover:bg-neon/20 hover:shadow-[0_0_20px_#39ff1420] disabled:opacity-20"
          >
            Proximo
          </button>
        )}
      </div>
    </div>
  );
}
