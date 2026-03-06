package chat

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BuildContext reads project context files and returns a formatted string.
func BuildContext(projectDir string) (string, error) {
	var parts []string

	files := []struct {
		path  string
		label string
	}{
		{"CLAUDE.md", "CLAUDE.md"},
		{"docs/CONTEXT.md", "docs/CONTEXT.md"},
		{"docs/ROADMAP.md", "docs/ROADMAP.md (ultimas fases)"},
	}

	for _, f := range files {
		p := filepath.Join(projectDir, f.path)
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		content := strings.TrimSpace(string(data))
		if content != "" {
			parts = append(parts, fmt.Sprintf("# %s\n\n%s", f.label, content))
		}
	}

	if len(parts) == 0 {
		return "", fmt.Errorf("no context files found")
	}

	return strings.Join(parts, "\n\n---\n\n"), nil
}

// ContextFiles returns which context files exist in the project.
func ContextFiles(projectDir string) []string {
	candidates := []string{
		"CLAUDE.md",
		"docs/CONTEXT.md",
		"docs/ROADMAP.md",
		"docs/SRS.md",
	}
	var found []string
	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(projectDir, c)); err == nil {
			found = append(found, c)
		}
	}
	return found
}
