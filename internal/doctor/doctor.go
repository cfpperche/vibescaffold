package doctor

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cfpperche/vibeforge/internal/i18n"
)

type Check struct {
	Name   string
	Status string // "ok", "warn", "fail"
	Detail string
}

func Run() []Check {
	checks := []Check{
		fileCheck("CLAUDE.md", true),
		fileCheck("docs/CONTEXT.md", true),
		fileCheck("docs/ROADMAP.md", false),
		fileCheck(".claude/settings.json", false),
		dirCheck(".claude/hooks", false),
		ciCheck(),
		linterCheck(),
		testRunnerCheck(),
	}
	return checks
}

func Score(checks []Check) (int, int) {
	total := len(checks)
	ok := 0
	for _, c := range checks {
		if c.Status == "ok" {
			ok++
		}
	}
	return ok, total
}

func fileCheck(path string, required bool) Check {
	name := path
	if _, err := os.Stat(path); err == nil {
		return Check{Name: name, Status: "ok", Detail: i18n.T("doctor.found")}
	}
	if required {
		return Check{Name: name, Status: "fail", Detail: i18n.T("doctor.not_found")}
	}
	return Check{Name: name, Status: "warn", Detail: i18n.T("doctor.not_found")}
}

func dirCheck(path string, required bool) Check {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		entries, _ := os.ReadDir(path)
		if len(entries) > 0 {
			return Check{Name: path, Status: "ok", Detail: i18n.T("doctor.configured")}
		}
		return Check{Name: path, Status: "warn", Detail: i18n.T("doctor.empty")}
	}
	if required {
		return Check{Name: path, Status: "fail", Detail: i18n.T("doctor.not_found")}
	}
	return Check{Name: path, Status: "warn", Detail: i18n.T("doctor.not_found")}
}

func ciCheck() Check {
	matches, _ := filepath.Glob(".github/workflows/*.yml")
	yaml, _ := filepath.Glob(".github/workflows/*.yaml")
	matches = append(matches, yaml...)
	if len(matches) > 0 {
		return Check{Name: ".github/workflows/", Status: "ok", Detail: i18n.T("doctor.ci_configured")}
	}
	return Check{Name: ".github/workflows/", Status: "warn", Detail: i18n.T("doctor.no_ci")}
}

func linterCheck() Check {
	linters := map[string]string{
		"biome.json":           "biome",
		".eslintrc.json":       "eslint",
		".eslintrc.js":         "eslint",
		"eslint.config.js":     "eslint",
		"eslint.config.mjs":    "eslint",
		".golangci.yml":        "golangci-lint",
		".golangci.yaml":       "golangci-lint",
		"ruff.toml":            "ruff",
		"pyproject.toml":       "ruff/black",
	}
	for file, name := range linters {
		if _, err := os.Stat(file); err == nil {
			return Check{Name: "Linter", Status: "ok", Detail: name}
		}
	}
	return Check{Name: "Linter", Status: "warn", Detail: i18n.T("doctor.no_linter")}
}

func testRunnerCheck() Check {
	runners := []struct {
		file string
		name string
	}{
		{"go.mod", "go test"},
		{"vitest.config.ts", "vitest"},
		{"jest.config.ts", "jest"},
		{"jest.config.js", "jest"},
		{"pytest.ini", "pytest"},
		{"pyproject.toml", "pytest"},
	}

	for _, r := range runners {
		if _, err := os.Stat(r.file); err == nil {
			if r.file == "pyproject.toml" {
				data, _ := os.ReadFile(r.file)
				if !strings.Contains(string(data), "pytest") {
					continue
				}
			}
			return Check{Name: "Test runner", Status: "ok", Detail: r.name}
		}
	}

	// Check package.json for test script
	if data, err := os.ReadFile("package.json"); err == nil {
		s := string(data)
		if strings.Contains(s, "\"test\"") {
			if strings.Contains(s, "vitest") {
				return Check{Name: "Test runner", Status: "ok", Detail: "vitest"}
			}
			if strings.Contains(s, "jest") {
				return Check{Name: "Test runner", Status: "ok", Detail: "jest"}
			}
			if strings.Contains(s, "bun test") {
				return Check{Name: "Test runner", Status: "ok", Detail: "bun test"}
			}
			return Check{Name: "Test runner", Status: "ok", Detail: i18n.T("doctor.test_found")}
		}
	}

	return Check{Name: "Test runner", Status: "warn", Detail: i18n.T("doctor.no_test")}
}
