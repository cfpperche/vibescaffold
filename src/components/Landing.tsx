import { Link } from '@tanstack/react-router'

const features = [
  { title: 'CLAUDE.md', desc: 'Regras e contexto gerados automaticamente para o Claude Code.' },
  { title: 'Git Hooks', desc: 'Pre-commit com lint e typecheck configurados no script.' },
  { title: 'CI/CD', desc: 'GitHub Actions pipeline pronto para build e testes.' },
  { title: 'Stack Completa', desc: 'Frontend, backend, banco de dados e testes em um comando.' },
  { title: 'docs/CONTEXT.md', desc: 'Documentacao de contexto para onboarding do AI.' },
  { title: 'Pronto pra Vibes', desc: 'Cole o script no terminal e comece a vibecoding.' },
]

export function Landing() {
  return (
    <div className="min-h-screen">
      {/* Hero */}
      <header className="flex flex-col items-center justify-center px-6 pt-24 pb-16 text-center">
        <div className="mb-4 inline-block rounded-full bg-primary/20 px-4 py-1 text-sm text-accent">
          Scaffold para Vibecoding
        </div>
        <h1 className="mb-6 text-5xl font-bold tracking-tight text-white md:text-7xl">
          Vibe<span className="text-primary">Scaffold</span>
        </h1>
        <p className="mb-10 max-w-2xl text-lg text-slate-400">
          Gere um projeto completo com CLAUDE.md, hooks, CI/CD e docs em 5 passos.
          Cole o script bash no terminal e comece a codar com o Claude Code.
        </p>
        <Link
          to="/wizard"
          className="rounded-xl bg-primary px-8 py-4 text-lg font-semibold text-white shadow-lg shadow-primary/25 transition hover:bg-primary-dark hover:shadow-xl"
        >
          Criar Projeto
        </Link>
      </header>

      {/* Features */}
      <section className="mx-auto max-w-5xl px-6 py-16">
        <h2 className="mb-12 text-center text-3xl font-bold text-white">O que voce ganha</h2>
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {features.map((f) => (
            <div
              key={f.title}
              className="rounded-2xl border border-surface-lighter bg-surface-light p-6 transition hover:border-primary/50"
            >
              <h3 className="mb-2 text-lg font-semibold text-white">{f.title}</h3>
              <p className="text-sm text-slate-400">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* How it works */}
      <section className="mx-auto max-w-3xl px-6 py-16 text-center">
        <h2 className="mb-12 text-3xl font-bold text-white">Como funciona</h2>
        <div className="flex flex-col gap-8 md:flex-row md:justify-between">
          {['Configure o projeto', 'Escolha a stack', 'Copie o script'].map((s, i) => (
            <div key={s} className="flex flex-col items-center">
              <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-primary text-xl font-bold text-white">
                {i + 1}
              </div>
              <p className="text-slate-300">{s}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-surface-lighter py-8 text-center text-sm text-slate-500">
        VibeScaffold — Scaffold para vibecoding com Claude Code
      </footer>
    </div>
  )
}
