package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/diff"
	"github.com/opencode-ai/opencode/internal/history"
	"github.com/opencode-ai/opencode/internal/logging"
	"github.com/opencode-ai/opencode/internal/lsp"
	"github.com/opencode-ai/opencode/internal/permission"
)

type WriteParams struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

type WritePermissionsParams struct {
	FilePath string `json:"file_path"`
	Diff     string `json:"diff"`
}

type writeTool struct {
	lspClients  map[string]*lsp.Client
	permissions permission.Service
	files       history.Service
}

type WriteResponseMetadata struct {
	Diff      string `json:"diff"`
	Additions int    `json:"additions"`
	Removals  int    `json:"removals"`
}

const (
	WriteToolName    = "write"
	writeDescription = `File writing tool that creates or updates files in the filesystem, allowing you to save or modify text content.

WHEN TO USE THIS TOOL:
- Use when you need to create a new file
- Helpful for updating existing files with modified content
- Perfect for saving generated code, configurations, or text data

HOW TO USE:
- Provide the path to the file you want to write
- Include the content to be written to the file
- The tool will create any necessary parent directories

FEATURES:
- Can create new files or overwrite existing ones
- Creates parent directories automatically if they don't exist
- Checks if the file has been modified since last read for safety
- Avoids unnecessary writes when content hasn't changed

LIMITATIONS:
- You should read a file before writing to it to avoid conflicts
- Cannot append to files (rewrites the entire file)


TIPS:
- Use the View tool first to examine existing files before modifying them
- Use the LS tool to verify the correct location when creating new files
- Combine with Glob and Grep tools to find and modify multiple files
- Always include descriptive comments when making changes to existing code`
)

func NewWriteTool(lspClients map[string]*lsp.Client, permissions permission.Service, files history.Service) BaseTool {
	return &writeTool{
		lspClients:  lspClients,
		permissions: permissions,
		files:       files,
	}
}

func (w *writeTool) Info() ToolInfo {
	return ToolInfo{
		Name:        WriteToolName,
		Description: writeDescription,
		Parameters: map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "The path to the file to write",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "The content to write to the file",
			},
		},
		Required: []string{"file_path", "content"},
	}
}

func (w *writeTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	logging.Info("Write tool starting execution", "input", call.Input[:min(100, len(call.Input))])
	
	var params WriteParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		logging.Info("Error parsing write tool parameters", "error", err)
		return NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}
	logging.Info("Write tool parameters parsed successfully", "file_path", params.FilePath, "content_length", len(params.Content))

	if params.FilePath == "" {
		logging.Info("Write tool error: file_path is required")
		return NewTextErrorResponse("file_path is required"), nil
	}

	if params.Content == "" {
		logging.Info("Write tool error: content is required")
		return NewTextErrorResponse("content is required"), nil
	}

	filePath := params.FilePath
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(config.WorkingDirectory(), filePath)
	}
	logging.Info("Write tool resolved file path", "original_path", params.FilePath, "resolved_path", filePath)

	fileInfo, err := os.Stat(filePath)
	logging.Info("Write tool file stat result", "path", filePath, "exists", err == nil, "error", err)
	if err == nil {
		if fileInfo.IsDir() {
			logging.Info("Write tool error: path is directory", "path", filePath)
			return NewTextErrorResponse(fmt.Sprintf("Path is a directory, not a file: %s", filePath)), nil
		}

		modTime := fileInfo.ModTime()
		lastRead := getLastReadTime(filePath)
		logging.Info("Write tool checking file modification time", "mod_time", modTime, "last_read", lastRead, "modified_since_read", modTime.After(lastRead))
		if modTime.After(lastRead) {
			logging.Info("Write tool error: file modified since last read")
			return NewTextErrorResponse(fmt.Sprintf("File %s has been modified since it was last read.\nLast modification: %s\nLast read: %s\n\nPlease read the file again before modifying it.",
				filePath, modTime.Format(time.RFC3339), lastRead.Format(time.RFC3339))), nil
		}

		oldContent, readErr := os.ReadFile(filePath)
		logging.Info("Write tool read existing content", "read_error", readErr, "content_length", len(oldContent), "content_matches", readErr == nil && string(oldContent) == params.Content)
		if readErr == nil && string(oldContent) == params.Content {
			logging.Info("Write tool: content unchanged, no write needed")
			return NewTextErrorResponse(fmt.Sprintf("File %s already contains the exact content. No changes made.", filePath)), nil
		}
	} else if !os.IsNotExist(err) {
		logging.Info("Write tool error checking file", "error", err)
		return ToolResponse{}, fmt.Errorf("error checking file: %w", err)
	}

	dir := filepath.Dir(filePath)
	logging.Info("Write tool creating directory", "dir", dir)
	if err = os.MkdirAll(dir, 0o755); err != nil {
		logging.Info("Write tool error creating directory", "dir", dir, "error", err)
		return ToolResponse{}, fmt.Errorf("error creating directory: %w", err)
	}

	oldContent := ""
	if fileInfo != nil && !fileInfo.IsDir() {
		oldBytes, readErr := os.ReadFile(filePath)
		if readErr == nil {
			oldContent = string(oldBytes)
		}
		logging.Info("Write tool loaded old content", "content_length", len(oldContent))
	}

	sessionID, messageID := GetContextValues(ctx)
	logging.Info("Write tool retrieved context values", "session_id", sessionID, "message_id", messageID)
	if sessionID == "" || messageID == "" {
		logging.Info("Write tool error: missing context values")
		return ToolResponse{}, fmt.Errorf("session_id and message_id are required")
	}

	diff, additions, removals := diff.GenerateDiff(
		oldContent,
		params.Content,
		filePath,
	)
	logging.Info("Write tool generated diff", "additions", additions, "removals", removals)

	rootDir := config.WorkingDirectory()
	permissionPath := filepath.Dir(filePath)
	if strings.HasPrefix(filePath, rootDir) {
		permissionPath = rootDir
	}
	logging.Info("Write tool requesting permissions", "permission_path", permissionPath, "file_path", filePath)
	p := w.permissions.Request(
		permission.CreatePermissionRequest{
			SessionID:   sessionID,
			Path:        permissionPath,
			ToolName:    WriteToolName,
			Action:      "write",
			Description: fmt.Sprintf("Create file %s", filePath),
			Params: WritePermissionsParams{
				FilePath: filePath,
				Diff:     diff,
			},
		},
	)
	logging.Info("Write tool permission result", "granted", p)
	if !p {
		logging.Info("Write tool error: permission denied")
		return ToolResponse{}, permission.ErrorPermissionDenied
	}

	logging.Info("Write tool writing file to disk", "path", filePath, "content_length", len(params.Content))
	err = os.WriteFile(filePath, []byte(params.Content), 0o644)
	if err != nil {
		logging.Info("Write tool error writing file", "error", err)
		return ToolResponse{}, fmt.Errorf("error writing file: %w", err)
	}
	logging.Info("Write tool file written successfully")

	// Check if file exists in history
	logging.Info("Write tool checking file history")
	file, err := w.files.GetByPathAndSession(ctx, filePath, sessionID)
	if err != nil {
		logging.Info("Write tool creating new file in history", "error", err)
		_, err = w.files.Create(ctx, sessionID, filePath, oldContent)
		if err != nil {
			// Log error but don't fail the operation
			logging.Info("Write tool error creating file history", "error", err)
			return ToolResponse{}, fmt.Errorf("error creating file history: %w", err)
		}
	}
	if file.Content != oldContent {
		// User Manually changed the content store an intermediate version
		logging.Info("Write tool creating intermediate history version")
		_, err = w.files.CreateVersion(ctx, sessionID, filePath, oldContent)
		if err != nil {
			logging.Debug("Error creating file history version", "error", err)
		}
	}
	// Store the new version
	logging.Info("Write tool creating new history version")
	_, err = w.files.CreateVersion(ctx, sessionID, filePath, params.Content)
	if err != nil {
		logging.Debug("Error creating file history version", "error", err)
	}

	logging.Info("Write tool recording file operations")
	recordFileWrite(filePath)
	RecordFileRead(filePath)
	
	logging.Info("Write tool waiting for LSP diagnostics")
	waitForLspDiagnostics(ctx, filePath, w.lspClients)

	logging.Info("Write tool generating final result")
	result := fmt.Sprintf("File successfully written: %s", filePath)
	result = fmt.Sprintf("<result>\n%s\n</result>", result)
	result += getDiagnostics(filePath, w.lspClients)
	
	logging.Info("Write tool completed successfully", "result_length", len(result))
	return WithResponseMetadata(NewTextResponse(result),
		WriteResponseMetadata{
			Diff:      diff,
			Additions: additions,
			Removals:  removals,
		},
	), nil
}
