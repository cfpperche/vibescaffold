import { Link } from '@tanstack/react-router';

const features = [
  { tag: 'CLAUDE.md', desc: 'Regras e contexto gerados para o Claude Code entender seu projeto.' },
  { tag: 'Git Hooks', desc: 'Pre-commit com Biome lint e TypeScript check automaticos.' },
  { tag: 'CI/CD', desc: 'GitHub Actions pipeline com build, lint e testes.' },
  { tag: 'Stack', desc: 'Frontend, backend, banco de dados configurados em um script.' },
  { tag: 'CONTEXT.md', desc: 'Documentacao de contexto para onboarding do AI coder.' },
  { tag: 'Docker', desc: 'Dockerfile e compose opcionais para deploy.' },
];

export function Landing() {
  return (
    <div className="min-h-screen">
      <header className="flex flex-col items-center justify-center px-6 pt-28 pb-20 text-center">
        <div className="mb-6 inline-block rounded-full border border-neon/20 bg-neon/5 px-4 py-1 text-xs text-neon">
          scaffold para vibecoding
        </div>
        <h1 className="mb-6 text-5xl font-bold tracking-tight text-white md:text-7xl">
          vibe<span className="text-neon">scaffold</span>
        </h1>
        <p className="mb-10 max-w-lg text-sm leading-relaxed text-neutral-500">
          Gere um projeto completo com CLAUDE.md, hooks, CI/CD e docs em 5 passos. Cole o script no
          terminal e comece a vibecoding.
        </p>
        <Link
          to="/wizard"
          className="rounded-lg border border-neon/50 bg-neon/10 px-8 py-3 text-sm font-semibold text-neon transition hover:bg-neon/20 hover:shadow-[0_0_30px_#39ff1425]"
        >
          $ criar projeto
        </Link>
      </header>

      <section className="mx-auto max-w-4xl px-6 py-16">
        <h2 className="mb-10 text-center text-xs font-semibold uppercase tracking-widest text-neon/50">
          o que voce ganha
        </h2>
        <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
          {features.map((f) => (
            <div
              key={f.tag}
              className="rounded-xl border border-surface-border bg-surface-light p-5 transition hover:border-neon/20"
            >
              <span className="mb-2 inline-block rounded-md bg-neon/10 px-2 py-0.5 text-[10px] font-bold uppercase text-neon/70">
                {f.tag}
              </span>
              <p className="text-xs leading-relaxed text-neutral-500">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-2xl px-6 py-16 text-center">
        <h2 className="mb-10 text-xs font-semibold uppercase tracking-widest text-neon/50">
          como funciona
        </h2>
        <div className="flex flex-col gap-6 md:flex-row md:justify-between">
          {[
            '$ configure o projeto',
            '$ escolha tipo & stack',
            '$ copie o script',
          ].map((s, i) => (
            <div key={s} className="flex flex-col items-center">
              <div className="mb-3 flex h-10 w-10 items-center justify-center rounded-full border border-neon/30 text-sm font-bold text-neon">
                {i + 1}
              </div>
              <p className="text-xs text-neutral-500">{s}</p>
            </div>
          ))}
        </div>
      </section>

      <footer className="border-t border-surface-border py-6 text-center text-[10px] uppercase tracking-widest text-neutral-700">
        vibescaffold — scaffold para vibecoding com claude code
      </footer>
    </div>
  );
}
