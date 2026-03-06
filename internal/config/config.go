package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config is the scaffold config (used when creating projects).
type Config struct {
	Name       string   `json:"name"`
	Desc       string   `json:"desc"`
	Author     string   `json:"author"`
	Type       string   `json:"type"`
	Stack      string   `json:"stack"`
	Principles []string `json:"principles"`
	Features   []string `json:"features"`
	Repo       bool     `json:"repo"`
}

// AppConfig is the global persistent config stored in ~/.vibeforge/config.json.
type AppConfig struct {
	ActiveAgent string `json:"active_agent"`
	OllamaModel string `json:"ollama_model"`
}

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".vibeforge"
	}
	return filepath.Join(home, ".vibeforge")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

// LoadAppConfig reads the persistent app config.
func LoadAppConfig() AppConfig {
	cfg := AppConfig{ActiveAgent: "claude"}
	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg
	}
	json.Unmarshal(data, &cfg)
	return cfg
}

// SaveAppConfig writes the persistent app config.
func SaveAppConfig(cfg AppConfig) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0o644)
}

// DetectProject checks if CLAUDE.md exists in cwd.
func DetectProject() bool {
	_, err := os.Stat("CLAUDE.md")
	return err == nil
}

// ProjectName returns the current directory name.
func ProjectName() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return filepath.Base(dir)
}
