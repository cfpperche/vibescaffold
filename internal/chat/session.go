package chat

import (
	"time"

	"github.com/cfpperche/vibeforge/internal/agent"
	"github.com/cfpperche/vibeforge/internal/config"
)

type Message struct {
	Role    string // "user", "agent", "system"
	Content string
	Time    time.Time
}

type Session struct {
	ProjectName string
	ProjectDir  string
	Agent       agent.Agent
	AppCfg      config.AppConfig
	Context     string
	History     []Message
}

func NewSession(projectDir, projectName string) *Session {
	appCfg := config.LoadAppConfig()

	// Find active agent
	var activeAgent agent.Agent
	for _, a := range agent.DefaultAgents() {
		if a.Key == appCfg.ActiveAgent {
			activeAgent = a
			break
		}
	}
	if activeAgent.Key == "" {
		activeAgent = agent.DefaultAgents()[0]
	}

	ctx, _ := agent.InjectContext(projectDir)

	return &Session{
		ProjectName: projectName,
		ProjectDir:  projectDir,
		Agent:       activeAgent,
		AppCfg:      appCfg,
		Context:     ctx,
		History:     []Message{},
	}
}

func (s *Session) AddMessage(role, content string) {
	s.History = append(s.History, Message{
		Role:    role,
		Content: content,
		Time:    time.Now(),
	})
}

func (s *Session) SwitchAgent(key string, ollamaModel string) bool {
	for _, a := range agent.DefaultAgents() {
		if a.Key == key {
			s.Agent = a
			s.AppCfg.ActiveAgent = key
			if ollamaModel != "" {
				s.AppCfg.OllamaModel = ollamaModel
			}
			config.SaveAppConfig(s.AppCfg)
			return true
		}
	}
	return false
}

func (s *Session) IsCommand(input string) bool {
	return len(input) > 0 && input[0] == '/'
}
