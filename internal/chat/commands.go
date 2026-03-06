package chat

import (
	"fmt"
	"strings"

	"github.com/cfpperche/vibescaffold/internal/agent"
	"github.com/cfpperche/vibescaffold/internal/doctor"
)

type CommandResult struct {
	Output string
	Quit   bool
}

// HandleCommand processes a slash command and returns the output.
func HandleCommand(session *Session, input string) CommandResult {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return CommandResult{Output: "comando vazio"}
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "/switch":
		return handleSwitch(session, args)
	case "/doctor":
		return handleDoctor(session)
	case "/status":
		return handleStatus(session)
	case "/context":
		return handleContext(session)
	case "/clear":
		return CommandResult{Output: "__clear__"}
	case "/exit":
		return CommandResult{Quit: true}
	case "/help":
		return handleHelp()
	default:
		return CommandResult{Output: fmt.Sprintf("comando desconhecido: %s\nDigite /help para ver comandos disponíveis", cmd)}
	}
}

func handleSwitch(session *Session, args []string) CommandResult {
	if len(args) == 0 {
		var lines []string
		lines = append(lines, "Uso: /switch <agente> [modelo]")
		lines = append(lines, "")
		detected := agent.DetectAll()
		for _, a := range detected {
			status := "✗"
			if a.Installed {
				status = "✓"
			}
			active := " "
			if a.Key == session.Agent.Key {
				active = "●"
			}
			lines = append(lines, fmt.Sprintf("  %s %s %-12s %s", active, status, a.Key, a.Name))
		}
		return CommandResult{Output: strings.Join(lines, "\n")}
	}

	key := args[0]
	ollamaModel := ""
	if key == "ollama" && len(args) > 1 {
		ollamaModel = args[1]
	}

	if session.SwitchAgent(key, ollamaModel) {
		name := session.Agent.Name
		if ollamaModel != "" {
			name += " (" + ollamaModel + ")"
		}
		return CommandResult{Output: fmt.Sprintf("● Trocando para %s...\n● Contexto mantido — injetando CONTEXT.md\n● %s ativo", name, name)}
	}

	return CommandResult{Output: fmt.Sprintf("agente desconhecido: %s", key)}
}

func handleDoctor(session *Session) CommandResult {
	checks := doctor.Run()
	ok, total := doctor.Score(checks)
	pct := 0
	if total > 0 {
		pct = ok * 100 / total
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("Doctor — %s", session.ProjectName))
	lines = append(lines, "")
	for _, c := range checks {
		icon := "✓"
		switch c.Status {
		case "warn":
			icon = "⚠"
		case "fail":
			icon = "✗"
		}
		lines = append(lines, fmt.Sprintf("  %s %-24s %s", icon, c.Name, c.Detail))
	}
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  Score: %d/%d (%d%%)", ok, total, pct))

	return CommandResult{Output: strings.Join(lines, "\n")}
}

func handleStatus(session *Session) CommandResult {
	ctx, err := BuildContext(session.ProjectDir)
	if err != nil {
		return CommandResult{Output: "Nenhum arquivo de roadmap encontrado"}
	}
	// Just show first 40 lines
	lines := strings.Split(ctx, "\n")
	if len(lines) > 40 {
		lines = lines[:40]
		lines = append(lines, "... (truncado)")
	}
	return CommandResult{Output: strings.Join(lines, "\n")}
}

func handleContext(session *Session) CommandResult {
	files := ContextFiles(session.ProjectDir)
	if len(files) == 0 {
		return CommandResult{Output: "Nenhum arquivo de contexto encontrado"}
	}

	var lines []string
	lines = append(lines, "Contexto carregado:")
	for _, f := range files {
		lines = append(lines, fmt.Sprintf("  ● %s", f))
	}
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Total: %d bytes de contexto", len(session.Context)))
	return CommandResult{Output: strings.Join(lines, "\n")}
}

func handleHelp() CommandResult {
	return CommandResult{Output: `Comandos disponíveis:
  /switch <agente>    Troca de agente LLM (claude, codex, gemini, ollama, aider)
  /doctor             Health check do projeto
  /status             Mostra contexto e roadmap
  /context            Mostra arquivos de contexto carregados
  /clear              Limpa o histórico visual
  /exit               Sai do VibeScaffold
  /help               Mostra esta ajuda`}
}
