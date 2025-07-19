package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/opencode-ai/opencode/internal/config"
)

// InitSystemPromptLogging creates the system_prompts_used directory and cleans it
func InitSystemPromptLogging() error {
	// Initialize all log directories
	logDirs := []string{
		"logs/system_prompts_used",
		"logs/input",
		"logs/output",
	}
	
	for _, logsDir := range logDirs {
		// Remove existing directory and recreate it (empty on app restart)
		if err := os.RemoveAll(logsDir); err != nil {
			Warn("Failed to remove existing log directory", "error", err, "dir", logsDir)
		}
		
		// Create the directory
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory %s: %w", logsDir, err)
		}
		
		Info("Log directory initialized", "path", logsDir)
	}
	
	return nil
}

// LogSystemPrompt writes the system prompt to a file with agent name and timestamp
func LogSystemPrompt(agentName config.AgentName, systemPrompt string) {
	if systemPrompt == "" {
		return
	}
	
	// Create timestamp for filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.txt", string(agentName), timestamp)
	filepath := filepath.Join("logs", "system_prompts_used", filename)
	
	// Write system prompt to file
	if err := os.WriteFile(filepath, []byte(systemPrompt), 0644); err != nil {
		Warn("Failed to write system prompt to file", "error", err, "filepath", filepath, "agent", agentName)
		return
	}
	
	Debug("System prompt logged", "filepath", filepath, "agent", agentName, "prompt_length", len(systemPrompt))
}

// LogInput writes the input content to a file with agent name and timestamp
func LogInput(agentName config.AgentName, inputContent string) {
	if inputContent == "" {
		return
	}
	
	// Create timestamp for filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.txt", string(agentName), timestamp)
	filepath := filepath.Join("logs", "input", filename)
	
	// Write input content to file
	if err := os.WriteFile(filepath, []byte(inputContent), 0644); err != nil {
		Warn("Failed to write input to file", "error", err, "filepath", filepath, "agent", agentName)
		return
	}
	
	Debug("Input logged", "filepath", filepath, "agent", agentName, "content_length", len(inputContent))
}

// LogOutput writes the output content to a file with agent name and timestamp
// For tool calls, it marshals the data to JSON with indentation
func LogOutput(agentName config.AgentName, content string, toolCalls interface{}) {
	if content == "" && toolCalls == nil {
		return
	}
	
	// Create timestamp for filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.txt", string(agentName), timestamp)
	filepath := filepath.Join("logs", "output", filename)
	
	var outputContent string
	
	// If there are tool calls, format as JSON
	if toolCalls != nil {
		jsonData, err := json.MarshalIndent(toolCalls, "", "  ")
		if err != nil {
			Warn("Failed to marshal tool calls to JSON", "error", err, "agent", agentName)
			outputContent = content // fallback to text content
		} else {
			if content != "" {
				outputContent = fmt.Sprintf("Content:\n%s\n\nTool Calls:\n%s", content, string(jsonData))
			} else {
				outputContent = fmt.Sprintf("Tool Calls:\n%s", string(jsonData))
			}
		}
	} else {
		outputContent = content
	}
	
	// Write output content to file
	if err := os.WriteFile(filepath, []byte(outputContent), 0644); err != nil {
		Warn("Failed to write output to file", "error", err, "filepath", filepath, "agent", agentName)
		return
	}
	
	Debug("Output logged", "filepath", filepath, "agent", agentName, "content_length", len(outputContent))
}