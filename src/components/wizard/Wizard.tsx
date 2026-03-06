import { useState } from 'react';
import { generateScript } from '../../lib/generateScript';
import { useWizardStore } from '../../stores/wizard';

const STEP_TITLES = ['Projeto', 'Stack', 'Claude MD', 'Hooks & CI', 'Gerar Script'];

function StepIndicator({ current }: { current: number }) {
  return (
    <div className="mb-10 flex items-center justify-center gap-2">
      {STEP_TITLES.map((title, i) => (
        <div key={title} className="flex items-center gap-2">
          <div
            className={`flex h-8 w-8 items-center justify-center rounded-full text-sm font-bold transition
            ${i === current ? 'bg-primary text-white' : i < current ? 'bg-accent text-surface' : 'bg-surface-lighter text-slate-400'}`}
          >
            {i < current ? '\u2713' : i + 1}
          </div>
          {i < STEP_TITLES.length - 1 && (
            <div className={`h-0.5 w-8 ${i < current ? 'bg-accent' : 'bg-surface-lighter'}`} />
          )}
        </div>
      ))}
    </div>
  );
}

function StepProject() {
  const { projectName, projectDescription, setField } = useWizardStore();
  return (
    <div className="space-y-6">
      <div>
        <label className="mb-2 block text-sm font-medium text-slate-300">Nome do Projeto</label>
        <input
          value={projectName}
          onChange={(e) => setField('projectName', e.target.value)}
          placeholder="meu-projeto"
          className="w-full rounded-lg border border-surface-lighter bg-surface px-4 py-3 text-white placeholder-slate-500 outline-none focus:border-primary"
        />
      </div>
      <div>
        <label className="mb-2 block text-sm font-medium text-slate-300">Descricao</label>
        <textarea
          value={projectDescription}
          onChange={(e) => setField('projectDescription', e.target.value)}
          placeholder="Descreva seu projeto..."
          rows={3}
          className="w-full rounded-lg border border-surface-lighter bg-surface px-4 py-3 text-white placeholder-slate-500 outline-none focus:border-primary"
        />
      </div>
    </div>
  );
}

const OPTIONS = {
  frontend: ['react', 'vue', 'svelte'],
  backend: ['node', 'hono', 'fastify'],
  database: ['postgres', 'sqlite', 'mysql'],
};

function SelectGroup({
  label,
  value,
  options,
  onChange,
}: {
  label: string;
  value: string;
  options: string[];
  onChange: (v: string) => void;
}) {
  return (
    <div>
      <label className="mb-2 block text-sm font-medium text-slate-300">{label}</label>
      <div className="flex gap-3">
        {options.map((opt) => (
          <button
            key={opt}
            onClick={() => onChange(opt)}
            className={`rounded-lg px-4 py-2 text-sm font-medium capitalize transition
              ${value === opt ? 'bg-primary text-white' : 'bg-surface-lighter text-slate-300 hover:bg-surface-lighter/80'}`}
          >
            {opt}
          </button>
        ))}
      </div>
    </div>
  );
}

function StepStack() {
  const { frontend, backend, database, setField } = useWizardStore();
  return (
    <div className="space-y-6">
      <SelectGroup
        label="Frontend"
        value={frontend}
        options={OPTIONS.frontend}
        onChange={(v) => setField('frontend', v)}
      />
      <SelectGroup
        label="Backend"
        value={backend}
        options={OPTIONS.backend}
        onChange={(v) => setField('backend', v)}
      />
      <SelectGroup
        label="Banco de Dados"
        value={database}
        options={OPTIONS.database}
        onChange={(v) => setField('database', v)}
      />
    </div>
  );
}

const RULE_OPTIONS = [
  { id: 'never-commit-without-build', label: 'Nunca commitar sem build' },
  { id: 'never-expose-secrets', label: 'Nunca expor secrets' },
  { id: 'always-push-after-commit', label: 'Sempre push apos commit' },
];

function CheckboxGroup({
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
      <label className="mb-3 block text-sm font-medium text-slate-300">{label}</label>
      <div className="space-y-2">
        {options.map((opt) => (
          <label
            key={opt.id}
            className="flex cursor-pointer items-center gap-3 rounded-lg bg-surface-lighter px-4 py-3"
          >
            <input
              type="checkbox"
              checked={selected.includes(opt.id)}
              onChange={() => toggle(opt.id)}
              className="accent-primary"
            />
            <span className="text-sm text-slate-300">{opt.label}</span>
          </label>
        ))}
      </div>
    </div>
  );
}

function StepClaudeMd() {
  const { claudeMdRules, setField } = useWizardStore();
  return (
    <CheckboxGroup
      label="Regras do CLAUDE.md"
      selected={claudeMdRules}
      options={RULE_OPTIONS}
      onChange={(v) => setField('claudeMdRules', v)}
    />
  );
}

const HOOK_OPTIONS = [
  { id: 'pre-commit-lint', label: 'Pre-commit: ESLint' },
  { id: 'pre-commit-typecheck', label: 'Pre-commit: TypeScript check' },
];
const CICD_OPTIONS = ['github-actions', 'none'];
const TESTING_OPTIONS = ['vitest', 'jest', 'none'];

function StepHooks() {
  const { hooks, cicd, testing, setField } = useWizardStore();
  return (
    <div className="space-y-6">
      <CheckboxGroup
        label="Git Hooks"
        selected={hooks}
        options={HOOK_OPTIONS}
        onChange={(v) => setField('hooks', v)}
      />
      <SelectGroup
        label="CI/CD"
        value={cicd}
        options={CICD_OPTIONS}
        onChange={(v) => setField('cicd', v)}
      />
      <SelectGroup
        label="Testing"
        value={testing}
        options={TESTING_OPTIONS}
        onChange={(v) => setField('testing', v)}
      />
    </div>
  );
}

function StepGenerate() {
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
        <h3 className="text-lg font-semibold text-white">Seu script esta pronto!</h3>
        <button
          onClick={copy}
          className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-surface transition hover:bg-accent/80"
        >
          {copied ? 'Copiado!' : 'Copiar Script'}
        </button>
      </div>
      <pre className="max-h-96 overflow-auto rounded-xl border border-surface-lighter bg-black p-4 text-sm text-green-400">
        {script}
      </pre>
    </div>
  );
}

const STEPS = [StepProject, StepStack, StepClaudeMd, StepHooks, StepGenerate];

export function Wizard() {
  const { step, nextStep, prevStep, projectName } = useWizardStore();
  const CurrentStep = STEPS[step];

  const canNext = step === 0 ? projectName.trim().length > 0 : true;

  return (
    <div className="mx-auto min-h-screen max-w-2xl px-6 py-12">
      <h1 className="mb-2 text-center text-3xl font-bold text-white">Criar Projeto</h1>
      <p className="mb-8 text-center text-slate-400">
        Passo {step + 1} de {STEP_TITLES.length}: {STEP_TITLES[step]}
      </p>

      <StepIndicator current={step} />

      <div className="rounded-2xl border border-surface-lighter bg-surface-light p-8">
        <CurrentStep />
      </div>

      <div className="mt-8 flex justify-between">
        <button
          onClick={prevStep}
          disabled={step === 0}
          className="rounded-lg bg-surface-lighter px-6 py-3 font-medium text-slate-300 transition hover:bg-surface-lighter/80 disabled:opacity-30"
        >
          Voltar
        </button>
        {step < 4 && (
          <button
            onClick={nextStep}
            disabled={!canNext}
            className="rounded-lg bg-primary px-6 py-3 font-medium text-white transition hover:bg-primary-dark disabled:opacity-30"
          >
            Proximo
          </button>
        )}
      </div>
    </div>
  );
}
