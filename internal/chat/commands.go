package chat

import (
	"fmt"
	"strings"

	"github.com/cfpperche/vibeforge/internal/agent"
	"github.com/cfpperche/vibeforge/internal/doctor"
	"github.com/cfpperche/vibeforge/internal/i18n"
)

type CommandResult struct {
	Output string
	Quit   bool
}

// HandleCommand processes a slash command and returns the output.
func HandleCommand(session *Session, input string) CommandResult {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return CommandResult{Output: i18n.T("commands.empty")}
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
		return CommandResult{Output: i18n.TF("commands.unknown", cmd)}
	}
}

func handleSwitch(session *Session, args []string) CommandResult {
	if len(args) == 0 {
		var lines []string
		lines = append(lines, i18n.T("commands.switch_usage"))
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
		return CommandResult{Output: i18n.TF("commands.switch_success", name, name)}
	}

	return CommandResult{Output: i18n.TF("commands.switch_unknown", key)}
}

func handleDoctor(session *Session) CommandResult {
	checks := doctor.Run()
	ok, total := doctor.Score(checks)
	pct := 0
	if total > 0 {
		pct = ok * 100 / total
	}

	var lines []string
	lines = append(lines, i18n.TF("commands.doctor_title", session.ProjectName))
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
	lines = append(lines, i18n.TF("commands.score", ok, total, pct))

	return CommandResult{Output: strings.Join(lines, "\n")}
}

func handleStatus(session *Session) CommandResult {
	ctx, err := BuildContext(session.ProjectDir)
	if err != nil {
		return CommandResult{Output: i18n.T("commands.no_roadmap")}
	}
	// Just show first 40 lines
	lines := strings.Split(ctx, "\n")
	if len(lines) > 40 {
		lines = lines[:40]
		lines = append(lines, i18n.T("chat.truncated"))
	}
	return CommandResult{Output: strings.Join(lines, "\n")}
}

func handleContext(session *Session) CommandResult {
	files := ContextFiles(session.ProjectDir)
	if len(files) == 0 {
		return CommandResult{Output: i18n.T("commands.no_context")}
	}

	var lines []string
	lines = append(lines, i18n.T("commands.context_loaded"))
	for _, f := range files {
		lines = append(lines, i18n.TF("commands.context_file", f))
	}
	lines = append(lines, "")
	lines = append(lines, i18n.TF("commands.context_total", len(session.Context)))
	return CommandResult{Output: strings.Join(lines, "\n")}
}

func handleHelp() CommandResult {
	return CommandResult{Output: i18n.T("commands.help")}
}
