package config

import (
	"os"
	"path/filepath"
)

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

func DetectProject() bool {
	_, err := os.Stat("CLAUDE.md")
	return err == nil
}

func ProjectName() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return filepath.Base(dir)
}
