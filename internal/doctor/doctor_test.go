package doctor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cfpperche/vibescaffold/internal/doctor"
)

func TestRunInEmptyDir(t *testing.T) {
	// Run in temp dir with nothing
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(orig)

	checks := doctor.Run()
	if len(checks) == 0 {
		t.Fatal("expected checks to be non-empty")
	}

	ok, total := doctor.Score(checks)
	if ok >= total {
		t.Errorf("expected score < total in empty dir, got %d/%d", ok, total)
	}
}

func TestRunWithClaudeMd(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(orig)

	// Create CLAUDE.md
	os.WriteFile("CLAUDE.md", []byte("# Test"), 0o644)

	checks := doctor.Run()
	found := false
	for _, c := range checks {
		if c.Name == "CLAUDE.md" && c.Status == "ok" {
			found = true
		}
	}
	if !found {
		t.Error("expected CLAUDE.md check to pass")
	}
}

func TestRunWithFullProject(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(orig)

	// Simulate a full project
	os.WriteFile("CLAUDE.md", []byte("# Test"), 0o644)
	os.MkdirAll("docs", 0o755)
	os.WriteFile("docs/CONTEXT.md", []byte("# Ctx"), 0o644)
	os.WriteFile("docs/ROADMAP.md", []byte("# Road"), 0o644)
	os.MkdirAll(".claude/hooks", 0o755)
	os.WriteFile(".claude/settings.json", []byte("{}"), 0o644)
	os.WriteFile(".claude/hooks/pre-commit", []byte("#!/bin/sh"), 0o644)
	os.MkdirAll(".github/workflows", 0o755)
	os.WriteFile(".github/workflows/ci.yml", []byte("name: CI"), 0o644)
	os.WriteFile("biome.json", []byte("{}"), 0o644)
	os.WriteFile("go.mod", []byte("module test"), 0o644)

	checks := doctor.Run()
	ok, total := doctor.Score(checks)
	if ok < total-1 {
		t.Errorf("expected nearly perfect score, got %d/%d", ok, total)
		for _, c := range checks {
			t.Logf("  %s: %s (%s)", c.Name, c.Status, c.Detail)
		}
	}
}

func TestScoreCalculation(t *testing.T) {
	checks := []doctor.Check{
		{Name: "a", Status: "ok"},
		{Name: "b", Status: "ok"},
		{Name: "c", Status: "warn"},
		{Name: "d", Status: "fail"},
	}
	ok, total := doctor.Score(checks)
	if ok != 2 || total != 4 {
		t.Errorf("expected 2/4, got %d/%d", ok, total)
	}
}

func TestLinterDetection(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(orig)

	os.WriteFile("biome.json", []byte("{}"), 0o644)

	checks := doctor.Run()
	for _, c := range checks {
		if c.Name == "Linter" {
			if c.Status != "ok" || c.Detail != "biome" {
				t.Errorf("expected Linter ok/biome, got %s/%s", c.Status, c.Detail)
			}
			return
		}
	}
	t.Error("Linter check not found")
}

func TestCIDetection(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(orig)

	os.MkdirAll(filepath.Join(tmp, ".github/workflows"), 0o755)
	os.WriteFile(filepath.Join(tmp, ".github/workflows/ci.yml"), []byte("name: CI"), 0o644)

	checks := doctor.Run()
	for _, c := range checks {
		if c.Name == ".github/workflows/" {
			if c.Status != "ok" {
				t.Errorf("expected CI ok, got %s", c.Status)
			}
			return
		}
	}
	t.Error("CI check not found")
}
