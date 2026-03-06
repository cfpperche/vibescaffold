package agent

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// DetectedAgent holds the detection result for a single agent.
type DetectedAgent struct {
	Agent
	Installed bool
	Version   string
	Models    []OllamaModel // only for ollama
}

// OllamaModel represents a model available in Ollama.
type OllamaModel struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type ollamaTagsResponse struct {
	Models []struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
	} `json:"models"`
}

// DetectAll checks which agents are installed on the system.
func DetectAll() []DetectedAgent {
	agents := DefaultAgents()
	results := make([]DetectedAgent, len(agents))

	for i, a := range agents {
		d := DetectedAgent{Agent: a}
		if path, err := exec.LookPath(a.Command); err == nil && path != "" {
			d.Installed = true
			d.Version = getVersion(a.Command)
		}
		if a.Key == "ollama" && d.Installed {
			d.Models = getOllamaModels()
			if !isOllamaRunning() {
				d.Version = "instalado (parado)"
			}
		}
		results[i] = d
	}
	return results
}

func getVersion(command string) string {
	cmd := exec.Command(command, "--version")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	v := strings.TrimSpace(string(out))
	// Take first line only
	if idx := strings.IndexByte(v, '\n'); idx != -1 {
		v = v[:idx]
	}
	// Truncate long versions
	if len(v) > 40 {
		v = v[:40] + "..."
	}
	return v
}

func isOllamaRunning() bool {
	client := &http.Client{Timeout: time.Second}
	resp, err := client.Get("http://localhost:11434")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}

func getOllamaModels() []OllamaModel {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var tags ollamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil
	}

	models := make([]OllamaModel, len(tags.Models))
	for i, m := range tags.Models {
		models[i] = OllamaModel{Name: m.Name, Size: m.Size}
	}
	return models
}
