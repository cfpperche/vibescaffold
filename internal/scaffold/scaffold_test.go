package scaffold_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cfpperche/vibescaffold/internal/config"
	"github.com/cfpperche/vibescaffold/internal/scaffold"
)

func TestScaffoldCreatesFiles(t *testing.T) {
	tmp := t.TempDir()

	cfg := config.Config{
		Name:       "test-project",
		Desc:       "A test project",
		Author:     "tester",
		Type:       "cli",
		Stack:      "go-bubbletea",
		Principles: []string{"tdd", "solid"},
		Features:   []string{"claude-md", "context-docs", "roadmap", "ci"},
		Repo:       false,
	}

	files, err := scaffold.Scaffold(cfg, tmp)
	if err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("expected files to be created")
	}

	// Check key files exist
	expectedFiles := []string{
		"CLAUDE.md",
		"docs/CONTEXT.md",
		"docs/ROADMAP.md",
		".claude/settings.json",
		".github/workflows/ci.yml",
		".gitignore",
	}

	root := filepath.Join(tmp, "test-project")
	for _, f := range expectedFiles {
		path := filepath.Join(root, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", f)
		}
	}

	// Check CLAUDE.md content
	data, err := os.ReadFile(filepath.Join(root, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	content := string(data)
	if !contains(content, "test-project") {
		t.Error("CLAUDE.md should contain project name")
	}
	if !contains(content, "TDD") && !contains(content, "Test-Driven") {
		t.Error("CLAUDE.md should contain TDD principle")
	}
}

func TestScaffoldMinimalFeatures(t *testing.T) {
	tmp := t.TempDir()

	cfg := config.Config{
		Name:     "minimal",
		Desc:     "Minimal project",
		Type:     "lib",
		Stack:    "go-chi",
		Features: []string{"claude-md"},
		Repo:     false,
	}

	files, err := scaffold.Scaffold(cfg, tmp)
	if err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	root := filepath.Join(tmp, "minimal")

	// CLAUDE.md should exist
	if _, err := os.Stat(filepath.Join(root, "CLAUDE.md")); os.IsNotExist(err) {
		t.Error("CLAUDE.md should exist")
	}

	// docs/CONTEXT.md should NOT exist (not in features)
	if _, err := os.Stat(filepath.Join(root, "docs/CONTEXT.md")); err == nil {
		t.Error("docs/CONTEXT.md should not exist with minimal features")
	}

	// CI should NOT exist
	if _, err := os.Stat(filepath.Join(root, ".github/workflows/ci.yml")); err == nil {
		t.Error("CI should not exist with minimal features")
	}

	_ = files
}

func TestScaffoldWithADRs(t *testing.T) {
	tmp := t.TempDir()

	cfg := config.Config{
		Name:     "adr-project",
		Desc:     "Project with ADRs",
		Type:     "saas",
		Stack:    "ts-react-vite",
		Features: []string{"claude-md", "adrs"},
		Repo:     false,
	}

	_, err := scaffold.Scaffold(cfg, tmp)
	if err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	adrPath := filepath.Join(tmp, "adr-project", "docs/adr/001-stack.md")
	if _, err := os.Stat(adrPath); os.IsNotExist(err) {
		t.Error("ADR file should exist")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
