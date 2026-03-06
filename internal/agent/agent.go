package agent

// Agent represents an LLM agent configuration.
type Agent struct {
	Key     string
	Name    string
	Command string
	Args    []string
}

// DefaultAgents returns the list of supported agents in display order.
func DefaultAgents() []Agent {
	return []Agent{
		{Key: "claude", Name: "Claude Code", Command: "claude"},
		{Key: "codex", Name: "Codex CLI", Command: "codex"},
		{Key: "gemini", Name: "Gemini CLI", Command: "gemini"},
		{Key: "ollama", Name: "Ollama", Command: "ollama", Args: []string{"run"}},
		{Key: "aider", Name: "Aider", Command: "aider", Args: []string{"--model", "gpt-4o"}},
	}
}

// InstallHint returns the install command for an agent.
func InstallHint(key string) string {
	switch key {
	case "claude":
		return "npm install -g @anthropic-ai/claude-code"
	case "codex":
		return "npm install -g @openai/codex"
	case "gemini":
		return "npm install -g @google/gemini-cli"
	case "ollama":
		return "curl -fsSL https://ollama.com/install.sh | sh"
	case "aider":
		return "pip install aider-chat"
	default:
		return ""
	}
}
