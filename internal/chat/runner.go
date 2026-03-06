package chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// StreamToken represents a chunk of output from the agent.
type StreamToken struct {
	Content string
	Done    bool
	Err     error
}

// RunAgent executes a CLI agent (claude, codex, gemini, aider) as a subprocess
// and streams its output line by line.
func RunAgent(session *Session, input string) (<-chan StreamToken, error) {
	key := session.Agent.Key

	if key == "ollama" {
		model := session.AppCfg.OllamaModel
		if model == "" {
			model = "llama3.2"
		}
		return RunOllama(session, input, model)
	}

	// For CLI agents, we pipe input and stream output
	var cmd *exec.Cmd
	switch key {
	case "claude":
		// claude --print sends to stdout without interactive mode
		args := []string{"--print"}
		if session.Context != "" {
			args = append(args, "--system-prompt", session.Context)
		}
		args = append(args, input)
		cmd = exec.Command("claude", args...)
	case "codex":
		args := []string{"--print"}
		if session.Context != "" {
			args = append(args, "--system-prompt", session.Context)
		}
		args = append(args, input)
		cmd = exec.Command("codex", args...)
	case "gemini":
		args := []string{"--print"}
		if session.Context != "" {
			args = append(args, "--system-prompt", session.Context)
		}
		args = append(args, input)
		cmd = exec.Command("gemini", args...)
	case "aider":
		args := []string{"--message", input, "--no-git", "--yes"}
		cmd = exec.Command("aider", args...)
	default:
		return nil, fmt.Errorf("unknown agent: %s", key)
	}

	cmd.Dir = session.ProjectDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start %s: %w", key, err)
	}

	ch := make(chan StreamToken, 64)
	go func() {
		defer close(ch)

		// Read stdout
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
		for scanner.Scan() {
			ch <- StreamToken{Content: scanner.Text() + "\n"}
		}

		// Read any stderr
		errOut, _ := io.ReadAll(stderr)
		if len(errOut) > 0 {
			errStr := strings.TrimSpace(string(errOut))
			if errStr != "" {
				ch <- StreamToken{Content: errStr + "\n"}
			}
		}

		err := cmd.Wait()
		if err != nil {
			ch <- StreamToken{Err: err}
		}
		ch <- StreamToken{Done: true}
	}()

	return ch, nil
}

// RunOllama streams a chat completion from the Ollama API.
func RunOllama(session *Session, input string, model string) (<-chan StreamToken, error) {
	// Build messages
	messages := []map[string]string{}
	if session.Context != "" {
		messages = append(messages, map[string]string{
			"role": "system", "content": session.Context,
		})
	}
	// Add history (last 20 messages for context window)
	start := 0
	if len(session.History) > 20 {
		start = len(session.History) - 20
	}
	for _, m := range session.History[start:] {
		role := "user"
		if m.Role == "agent" {
			role = "assistant"
		} else if m.Role == "system" {
			continue
		}
		messages = append(messages, map[string]string{
			"role": role, "content": m.Content,
		})
	}
	messages = append(messages, map[string]string{
		"role": "user", "content": input,
	})

	body := map[string]interface{}{
		"model":    model,
		"messages": messages,
		"stream":   true,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Post("http://localhost:11434/api/chat", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("ollama API: %w", err)
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("ollama API returned %d", resp.StatusCode)
	}

	ch := make(chan StreamToken, 64)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			var chunk struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
				Done bool `json:"done"`
			}
			if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
				continue
			}
			if chunk.Message.Content != "" {
				ch <- StreamToken{Content: chunk.Message.Content}
			}
			if chunk.Done {
				ch <- StreamToken{Done: true}
				return
			}
		}
		ch <- StreamToken{Done: true}
	}()

	return ch, nil
}
