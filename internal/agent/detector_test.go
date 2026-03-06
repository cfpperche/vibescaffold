package agent_test

import (
	"testing"

	"github.com/cfpperche/vibescaffold/internal/agent"
)

func TestDetectAllReturnsAllAgents(t *testing.T) {
	detected := agent.DetectAll()
	if len(detected) != 5 {
		t.Errorf("expected 5 agents, got %d", len(detected))
	}

	keys := map[string]bool{}
	for _, d := range detected {
		keys[d.Key] = true
	}

	expected := []string{"claude", "codex", "gemini", "ollama", "aider"}
	for _, k := range expected {
		if !keys[k] {
			t.Errorf("missing agent: %s", k)
		}
	}
}

func TestDefaultAgents(t *testing.T) {
	agents := agent.DefaultAgents()
	if len(agents) == 0 {
		t.Fatal("expected default agents")
	}
	for _, a := range agents {
		if a.Key == "" || a.Name == "" || a.Command == "" {
			t.Errorf("agent has empty fields: %+v", a)
		}
	}
}

func TestInstallHint(t *testing.T) {
	tests := []struct {
		key  string
		want string
	}{
		{"claude", "npm install"},
		{"codex", "npm install"},
		{"gemini", "npm install"},
		{"ollama", "curl"},
		{"aider", "pip install"},
		{"unknown", ""},
	}
	for _, tt := range tests {
		hint := agent.InstallHint(tt.key)
		if tt.want != "" && !contains(hint, tt.want) {
			t.Errorf("InstallHint(%s) = %q, want to contain %q", tt.key, hint, tt.want)
		}
		if tt.want == "" && hint != "" {
			t.Errorf("InstallHint(%s) = %q, want empty", tt.key, hint)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && searchStr(s, sub)
}

func searchStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
