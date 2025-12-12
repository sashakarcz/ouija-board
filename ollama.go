package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OllamaClient handles communication with the Ollama API
type OllamaClient struct {
	url       string
	model     string
	timeout   time.Duration
	maxTokens int
	client    *http.Client
}

// OllamaRequest represents the request payload to Ollama API
type OllamaRequest struct {
	Model  string        `json:"model"`
	Prompt string        `json:"prompt"`
	Stream bool          `json:"stream"`
	Options OllamaOptions `json:"options"`
}

// OllamaOptions contains generation options
type OllamaOptions struct {
	NumPredict int `json:"num_predict"`
}

// OllamaResponse represents a single line of the streaming response
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(url, model string, timeout time.Duration, maxTokens int) *OllamaClient {
	return &OllamaClient{
		url:       url,
		model:     model,
		timeout:   timeout,
		maxTokens: maxTokens,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// GenerateAnswer generates an answer using the Ollama API
func (c *OllamaClient) GenerateAnswer(ctx context.Context, question string) (string, error) {
	// Validate input
	if len(question) > 1000 {
		return "", errors.New("question too long")
	}

	// Sanitize question
	question = sanitizeInput(question)

	// Create mystical prompt
	prompt := fmt.Sprintf(
		"Pretend that you are a Ouija board. As a mystical Ouija board, answer the following question in a short answer. "+
			"Respond without using any actions, such as *smiles*, *laughs*, or any text within asterisks. "+
			"If the question is a yes or no question, answer with a yes or a no. "+
			"If the user says goodbye, bye, or farewell, respond with 'Goodbye.' Question: %s",
		question,
	)

	// Create request payload
	reqPayload := OllamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: true,
		Options: OllamaOptions{
			NumPredict: c.maxTokens,
		},
	}

	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return "The spirits cannot answer at this time. Try again later.", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "The spirits cannot answer at this time. Try again later.", nil
	}

	// Process streaming response
	answer := strings.Builder{}
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var ollamaResp OllamaResponse
		if err := json.Unmarshal(line, &ollamaResp); err != nil {
			// Skip malformed lines
			continue
		}

		answer.WriteString(ollamaResp.Response)

		if ollamaResp.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return "The spirits cannot answer at this time. Try again later.", nil
	}

	result := strings.TrimSpace(answer.String())
	if result == "" {
		return "The spirits cannot answer at this time. Try again later.", nil
	}

	return result, nil
}

// sanitizeInput removes potentially dangerous characters from input
func sanitizeInput(input string) string {
	// Remove control characters and trim whitespace
	input = strings.TrimSpace(input)

	// Replace null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	return input
}
