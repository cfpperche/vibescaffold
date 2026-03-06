package tui

import (
	"os"
	"os/exec"

	"github.com/cfpperche/vibeforge/internal/tui/views"
)

func launchAgentCmd(msg views.LaunchAgentMsg) *exec.Cmd {
	cwd, _ := os.Getwd()

	switch msg.AgentKey {
	case "claude":
		cmd := exec.Command("claude")
		cmd.Dir = cwd
		return cmd
	case "codex":
		cmd := exec.Command("codex")
		cmd.Dir = cwd
		return cmd
	case "gemini":
		cmd := exec.Command("gemini")
		cmd.Dir = cwd
		return cmd
	case "ollama":
		model := msg.OllamaModel
		if model == "" {
			model = "llama3.2"
		}
		cmd := exec.Command("ollama", "run", model)
		cmd.Dir = cwd
		return cmd
	case "aider":
		args := []string{}
		if _, err := os.Stat("docs/CONTEXT.md"); err == nil {
			args = append(args, "--read", "docs/CONTEXT.md")
		}
		if _, err := os.Stat("CLAUDE.md"); err == nil {
			args = append(args, "--read", "CLAUDE.md")
		}
		cmd := exec.Command("aider", args...)
		cmd.Dir = cwd
		return cmd
	default:
		return exec.Command("echo", "Unknown agent:", msg.AgentKey)
	}
}
