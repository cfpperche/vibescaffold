package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cfpperche/vibeforge/internal/config"
)

var principleLabels = map[string]string{
	"tdd":              "Test-Driven Development — escreva testes antes do codigo",
	"clean-arch":       "Clean Architecture — separacao de camadas",
	"solid":            "SOLID — principios de design OO",
	"12-factor":        "12-Factor App — boas praticas de deploy",
	"ddd":              "Domain-Driven Design — modelagem pelo dominio",
	"cqrs":             "CQRS — separacao de leitura e escrita",
	"event-sourcing":   "Event Sourcing — historico de eventos como fonte de verdade",
	"hexagonal":        "Hexagonal Architecture — ports and adapters",
}

func Scaffold(cfg config.Config, targetDir string) ([]string, error) {
	root := filepath.Join(targetDir, cfg.Name)
	var created []string

	dirs := []string{
		"src", "docs", ".claude",
	}
	for _, d := range dirs {
		p := filepath.Join(root, d)
		if err := os.MkdirAll(p, 0o755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %w", d, err)
		}
	}

	// CLAUDE.md
	claudeMd := generateClaudeMd(cfg)
	if err := writeFile(root, "CLAUDE.md", claudeMd); err != nil {
		return nil, err
	}
	created = append(created, "CLAUDE.md")

	// docs/CONTEXT.md
	if hasFeature(cfg, "context-docs") {
		ctx := generateContext(cfg)
		if err := writeFile(root, "docs/CONTEXT.md", ctx); err != nil {
			return nil, err
		}
		created = append(created, "docs/CONTEXT.md")
	}

	// docs/ROADMAP.md
	if hasFeature(cfg, "roadmap") {
		roadmap := fmt.Sprintf("# %s — Roadmap\n\n## Fase 0 — Setup\n- [x] Scaffold inicial\n\n## Fase 1 — MVP\n- [ ] TODO\n", cfg.Name)
		if err := writeFile(root, "docs/ROADMAP.md", roadmap); err != nil {
			return nil, err
		}
		created = append(created, "docs/ROADMAP.md")
	}

	// docs/SRS.md
	if hasFeature(cfg, "srs") {
		srs := fmt.Sprintf("# %s — Software Requirements Specification\n\n## 1. Introducao\n%s\n\n## 2. Requisitos Funcionais\n- [ ] TODO\n\n## 3. Requisitos Nao-Funcionais\n- [ ] TODO\n", cfg.Name, cfg.Desc)
		if err := writeFile(root, "docs/SRS.md", srs); err != nil {
			return nil, err
		}
		created = append(created, "docs/SRS.md")
	}

	// ADR template
	if hasFeature(cfg, "adrs") {
		if err := os.MkdirAll(filepath.Join(root, "docs/adr"), 0o755); err != nil {
			return nil, err
		}
		adr := "# ADR-001: Escolha de stack\n\n## Status\nAceito\n\n## Contexto\nPrecisamos definir a stack do projeto.\n\n## Decisao\nUsaremos " + cfg.Stack + ".\n\n## Consequencias\nTime precisa conhecer a stack escolhida.\n"
		if err := writeFile(root, "docs/adr/001-stack.md", adr); err != nil {
			return nil, err
		}
		created = append(created, "docs/adr/001-stack.md")
	}

	// .claude/settings.json
	settings := `{
  "permissions": {
    "allow": ["Read", "Write", "Edit", "Bash", "Glob", "Grep"]
  }
}
`
	if err := writeFile(root, ".claude/settings.json", settings); err != nil {
		return nil, err
	}
	created = append(created, ".claude/settings.json")

	// .claude/commands
	if hasFeature(cfg, "commands") {
		if err := os.MkdirAll(filepath.Join(root, ".claude/commands"), 0o755); err != nil {
			return nil, err
		}
		commitCmd := "Review all staged changes, then create a commit with a descriptive message following conventional commits."
		if err := writeFile(root, ".claude/commands/commit.md", commitCmd); err != nil {
			return nil, err
		}
		created = append(created, ".claude/commands/commit.md")
	}

	// Hooks
	if hasFeature(cfg, "hooks") {
		if err := os.MkdirAll(filepath.Join(root, ".claude/hooks"), 0o755); err != nil {
			return nil, err
		}
		hook := `{
  "hooks": [
    {
      "event": "pre-commit",
      "command": "go build ./... && go vet ./..."
    }
  ]
}
`
		if err := writeFile(root, ".claude/hooks/hooks.json", hook); err != nil {
			return nil, err
		}
		created = append(created, ".claude/hooks/hooks.json")
	}

	// CI
	if hasFeature(cfg, "ci") {
		if err := os.MkdirAll(filepath.Join(root, ".github/workflows"), 0o755); err != nil {
			return nil, err
		}
		ci := generateCI(cfg)
		if err := writeFile(root, ".github/workflows/ci.yml", ci); err != nil {
			return nil, err
		}
		created = append(created, ".github/workflows/ci.yml")
	}

	// .gitignore
	gitignore := "node_modules/\ndist/\n.env\n.env.local\n*.log\n.DS_Store\n/tmp/\n"
	if err := writeFile(root, ".gitignore", gitignore); err != nil {
		return nil, err
	}
	created = append(created, ".gitignore")

	// git init + commit
	if err := runInDir(root, "git", "init"); err != nil {
		return created, fmt.Errorf("git init: %w", err)
	}
	if err := runInDir(root, "git", "add", "-A"); err != nil {
		return created, fmt.Errorf("git add: %w", err)
	}
	if err := runInDir(root, "git", "commit", "-m", "Initial scaffold by VibeForge"); err != nil {
		return created, fmt.Errorf("git commit: %w", err)
	}

	// gh repo create
	if cfg.Repo {
		_ = runInDir(root, "gh", "repo", "create", cfg.Name, "--private", "--source", ".", "--push")
	}

	return created, nil
}

func generateClaudeMd(cfg config.Config) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", cfg.Name)
	fmt.Fprintf(&b, "## O que e\n%s\n\n", cfg.Desc)
	fmt.Fprintf(&b, "## Tipo\n%s\n\n", cfg.Type)
	fmt.Fprintf(&b, "## Stack\n%s\n\n", cfg.Stack)

	if len(cfg.Principles) > 0 {
		b.WriteString("## Principios\n")
		for i, p := range cfg.Principles {
			label, ok := principleLabels[p]
			if !ok {
				label = p
			}
			fmt.Fprintf(&b, "%d. %s\n", i+1, label)
		}
		b.WriteString("\n")
	}

	b.WriteString("## Regras\n")
	b.WriteString("1. NUNCA commite sem build passando\n")
	b.WriteString("2. NUNCA exponha secrets\n")
	b.WriteString("3. SEMPRE git push apos commit\n\n")
	b.WriteString("## Workflow\nTarefa -> le arquivos -> implementa -> build -> commit -> push\n")

	return b.String()
}

func generateContext(cfg config.Config) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s — Contexto\n\n", cfg.Name)
	fmt.Fprintf(&b, "## Produto\n%s\n\n", cfg.Desc)
	fmt.Fprintf(&b, "## Tipo\n%s\n\n", cfg.Type)
	fmt.Fprintf(&b, "## Stack\n%s\n\n", cfg.Stack)
	if cfg.Author != "" {
		fmt.Fprintf(&b, "## Autor\n%s\n", cfg.Author)
	}
	return b.String()
}

func generateCI(cfg config.Config) string {
	return `name: CI
on:
  push: { branches: [main] }
  pull_request: { branches: [main] }
jobs:
  ci:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: go build ./...
      - run: go test ./...
      - run: go vet ./...
`
}

func hasFeature(cfg config.Config, f string) bool {
	for _, feat := range cfg.Features {
		if feat == f {
			return true
		}
	}
	return false
}

func writeFile(root, rel, content string) error {
	p := filepath.Join(root, rel)
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, []byte(content), 0o644)
}

func runInDir(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
