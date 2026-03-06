package agent

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Launch starts the selected agent in the given project directory.
// It suspends the TUI (returns exec.Cmd to be run with os/exec).
func Launch(agentKey string, projectDir string, ollamaModel string) error {
	switch agentKey {
	case "claude":
		return runInteractive(projectDir, "claude")

	case "codex":
		return runInteractive(projectDir, "codex")

	case "gemini":
		return runInteractive(projectDir, "gemini")

	case "ollama":
		model := ollamaModel
		if model == "" {
			model = "llama3.2"
		}
		ctx, err := InjectContext(projectDir)
		if err == nil && ctx != "" {
			// Use system prompt via ollama run --system
			return runInteractive(projectDir, "ollama", "run", model, "--system", ctx)
		}
		return runInteractive(projectDir, "ollama", "run", model)

	case "aider":
		args := []string{}
		contextPath := filepath.Join(projectDir, "docs/CONTEXT.md")
		claudePath := filepath.Join(projectDir, "CLAUDE.md")
		if _, err := os.Stat(contextPath); err == nil {
			args = append(args, "--read", contextPath)
		}
		if _, err := os.Stat(claudePath); err == nil {
			args = append(args, "--read", claudePath)
		}
		return runInteractive(projectDir, "aider", args...)

	default:
		return fmt.Errorf("unknown agent: %s", agentKey)
	}
}

// InjectContext reads CLAUDE.md and docs/CONTEXT.md and returns a combined context string.
func InjectContext(projectDir string) (string, error) {
	var parts []string

	claudePath := filepath.Join(projectDir, "CLAUDE.md")
	if data, err := os.ReadFile(claudePath); err == nil {
		parts = append(parts, string(data))
	}

	contextPath := filepath.Join(projectDir, "docs/CONTEXT.md")
	if data, err := os.ReadFile(contextPath); err == nil {
		parts = append(parts, string(data))
	}

	if len(parts) == 0 {
		return "", fmt.Errorf("no context files found in %s", projectDir)
	}

	return strings.Join(parts, "\n\n---\n\n"), nil
}

func runInteractive(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
