package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// AgentCallLogEntry represents a structured log entry for agent tool calls
type AgentCallLogEntry struct {
	Timestamp float64     `json:"timestamp"`
	Agent     string      `json:"agent"`
	Tool      string      `json:"tool"`
	Input     interface{} `json:"input"`
	Output    interface{} `json:"output"`
}

var agentCallLogFile *os.File
var agentCallLogMutex sync.Mutex

// InitSystemPromptLogging creates the system_prompts_used directory and cleans it
func InitSystemPromptLogging() error {
	// Initialize all log directories
	logDirs := []string{
		"logs/system_prompts_used",
		"logs/input",
		"logs/output",
		"logs/agent_calls",
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
	
	// Initialize agent call log file
	if err := initAgentCallLog(); err != nil {
		return fmt.Errorf("failed to initialize agent call log: %w", err)
	}
	
	return nil
}

// initAgentCallLog creates and opens the agent call log file for writing
func initAgentCallLog() error {
	agentCallLogMutex.Lock()
	defer agentCallLogMutex.Unlock()

	// Close existing log file if open
	if agentCallLogFile != nil {
		agentCallLogFile.Close()
	}

	// Create agent_calls.jsonl file
	logPath := filepath.Join("logs", "agent_calls", "agent_calls.jsonl")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create agent call log file: %w", err)
	}

	agentCallLogFile = file
	Info("Agent call log initialized", "path", logPath)
	return nil
}

// LogSystemPrompt writes the system prompt to a file with agent name and timestamp
func LogSystemPrompt(agentName string, systemPrompt string) {
	if systemPrompt == "" {
		return
	}
	
	// Create timestamp for filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.txt", agentName, timestamp)
	filepath := filepath.Join("logs", "system_prompts_used", filename)
	
	// Write system prompt to file
	if err := os.WriteFile(filepath, []byte(systemPrompt), 0644); err != nil {
		Warn("Failed to write system prompt to file", "error", err, "filepath", filepath, "agent", agentName)
		return
	}
	
	Debug("System prompt logged", "filepath", filepath, "agent", agentName, "prompt_length", len(systemPrompt))
}

// LogInput writes the input content to a file with agent name and timestamp
func LogInput(agentName string, inputContent string, systemPrompt string) {
	if inputContent == "" {
		return
	}
	
	// Create timestamp for filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.txt", agentName, timestamp)
	filepath := filepath.Join("logs", "input", filename)
	
	// Combine input content and system prompt
	var combinedContent string
	if systemPrompt != "" {
		combinedContent = fmt.Sprintf("System Prompt:\n%s\n\nInput Content:\n%s", systemPrompt, inputContent)
	} else {
		combinedContent = inputContent
	}
	
	// Write combined content to file
	if err := os.WriteFile(filepath, []byte(combinedContent), 0644); err != nil {
		Warn("Failed to write input to file", "error", err, "filepath", filepath, "agent", agentName)
		return
	}
	
	Debug("Input logged", "filepath", filepath, "agent", agentName, "content_length", len(combinedContent))
}

// LogOutput writes the output content to a file with agent name and timestamp
// For tool calls, it marshals the data to JSON with indentation
func LogOutput(agentName string, content string, toolCalls interface{}) {
	if content == "" && toolCalls == nil {
		return
	}
	
	// Create timestamp for filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.txt", agentName, timestamp)
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

// LogAgentCall logs an agent tool call invocation to the JSONL log file
func LogAgentCall(agentName string, toolName string, inputData interface{}, outputData interface{}) {
	if agentCallLogFile == nil {
		Warn("Agent call log file not initialized, skipping log entry")
		return
	}

	agentCallLogMutex.Lock()
	defer agentCallLogMutex.Unlock()

	// Create log entry
	entry := AgentCallLogEntry{
		Timestamp: float64(time.Now().UnixMilli()) / 1000.0, // Convert to seconds with decimal precision
		Agent:     agentName,
		Tool:      toolName,
		Input:     inputData,
		Output:    outputData,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		Warn("Failed to marshal agent call log entry", "error", err, "agent", agentName, "tool", toolName)
		return
	}

	// Write to JSONL file (one JSON object per line)
	if _, err := agentCallLogFile.WriteString(string(jsonData) + "\n"); err != nil {
		Warn("Failed to write agent call log entry", "error", err, "agent", agentName, "tool", toolName)
		return
	}

	// Ensure data is written to disk
	if err := agentCallLogFile.Sync(); err != nil {
		Warn("Failed to sync agent call log file", "error", err)
	}

	Debug("Agent call logged", "agent", agentName, "tool", toolName, "timestamp", entry.Timestamp)
}

// CloseAgentCallLog closes the agent call log file (should be called on shutdown)
func CloseAgentCallLog() {
	agentCallLogMutex.Lock()
	defer agentCallLogMutex.Unlock()

	if agentCallLogFile != nil {
		agentCallLogFile.Close()
		agentCallLogFile = nil
	}
}