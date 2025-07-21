// package tools

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"time"

// 	"github.com/opencode-ai/opencode/internal/config"
// 	"github.com/opencode-ai/opencode/internal/diff"
// 	"github.com/opencode-ai/opencode/internal/history"
// 	"github.com/opencode-ai/opencode/internal/logging"
// 	"github.com/opencode-ai/opencode/internal/lsp"
// 	"github.com/opencode-ai/opencode/internal/permission"
// )

// type ConfirmWriteParams struct {
// 	FilePath string `json:"file_path"`
// 	// Content  string `json:"content"` I would like the orchestrator to not have to send this
// }

// type ConfirmWritePermissionsParams struct {
// 	FilePath string `json:"file_path"`
// 	Diff     string `json:"diff"`
// }

// type confirmWriteTool struct {
// 	lspClients  map[string]*lsp.Client
// 	permissions permission.Service
// 	files       history.Service
// }

// type ConfirmWriteResponseMetadata struct {
// 	Diff      string `json:"diff"`
// 	Additions int    `json:"additions"`
// 	Removals  int    `json:"removals"`
// }

// const (
// 	ConfirmWriteToolName    = "confirm-write"
// 	confirmWriteDescription = `Tool to confirm the writing of the file to example.tsx.

// WHEN TO USE THIS TOOL:
// - Only use this tool in the case of a WRITE request from the user, ie a request to write a file from scratch.
// - Use when you need to confirm the creation a new file. So, more specifically, use it right after you have called the Coder Agent, and you have received a typescript code from the latter.

// HOW TO USE:
// - No parameter has be used. This tool being called directly writes the exact string content of the generated code to the file "example.tsx".
// `
// )

// func NewConfirmWriteTool(lspClients map[string]*lsp.Client, permissions permission.Service, files history.Service) BaseTool {
// 	return &confirmWriteTool{
// 		lspClients:  lspClients,
// 		permissions: permissions,
// 		files:       files,
// 	}
// }

// func (w *confirmWriteTool) Info() ToolInfo {
// 	return ToolInfo{
// 		Name:        ConfirmWriteToolName,
// 		Description: confirmWriteDescription,
// 		Parameters: map[string]any{
// 			"file_path": map[string]any{
// 				"type":        "string",
// 				"description": "The path to the file to write",
// 			},
// 		},
// 		Required: []string{"file_path", "content"},
// 	}
// }

// func (w *confirmWriteTool) Run(ctx coantext.Context, call ToolCall) (ToolResponse, error) {
// 	var params ConfirmWriteParams
// 	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
// 		return NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
// 	}

// 	if params.FilePath == "" {
// 		return NewTextErrorResponse("file_path is required"), nil
// 	}

// 	filePath := params.FilePath
// 	if !filepath.IsAbs(filePath) {
// 		filePath = filepath.Join(config.WorkingDirectory(), filePath)
// 	}

// 	fileInfo, err := os.Stat(filePath)
// 	if err == nil {
// 		if fileInfo.IsDir() {
// 			return NewTextErrorResponse(fmt.Sprintf("Path is a directory, not a file: %s", filePath)), nil
// 		}

// 		modTime := fileInfo.ModTime()
// 		lastRead := getLastReadTime(filePath)
// 		if modTime.After(lastRead) {
// 			return NewTextErrorResponse(fmt.Sprintf("File %s has been modified since it was last read.\nLast modification: %s\nLast read: %s\n\nPlease read the file again before modifying it.",
// 				filePath, modTime.Format(time.RFC3339), lastRead.Format(time.RFC3339))), nil
// 		}

// 		oldContent, readErr := os.ReadFile(filePath)
// 		if readErr == nil && string(oldContent) == params.Content {
// 			return NewTextErrorResponse(fmt.Sprintf("File %s already contains the exact content. No changes made.", filePath)), nil
// 		}
// 	} else if !os.IsNotExist(err) {
// 		return ToolResponse{}, fmt.Errorf("error checking file: %w", err)
// 	}

// 	dir := filepath.Dir(filePath)
// 	if err = os.MkdirAll(dir, 0o755); err != nil {
// 		return ToolResponse{}, fmt.Errorf("error creating directory: %w", err)
// 	}

// 	oldContent := ""
// 	if fileInfo != nil && !fileInfo.IsDir() {
// 		oldBytes, readErr := os.ReadFile(filePath)
// 		if readErr == nil {
// 			oldContent = string(oldBytes)
// 		}
// 	}

// 	sessionID, messageID := GetContextValues(ctx)
// 	if sessionID == "" || messageID == "" {
// 		return ToolResponse{}, fmt.Errorf("session_id and message_id are required")
// 	}

// 	diff, additions, removals := diff.GenerateDiff(
// 		oldContent,
// 		params.Content,
// 		filePath,
// 	)

// 	rootDir := config.WorkingDirectory()
// 	permissionPath := filepath.Dir(filePath)
// 	if strings.HasPrefix(filePath, rootDir) {
// 		permissionPath = rootDir
// 	}
// 	p := w.permissions.Request(
// 		permission.CreatePermissionRequest{
// 			SessionID:   sessionID,
// 			Path:        permissionPath,
// 			ToolName:    ConfirmWriteToolName,
// 			Action:      "write",
// 			Description: fmt.Sprintf("Create file %s", filePath),
// 			Params: ConfirmWritePermissionsParams{
// 				FilePath: filePath,
// 				Diff:     diff,
// 			},
// 		},
// 	)
// 	if !p {
// 		return ToolResponse{}, permission.ErrorPermissionDenied
// 	}

// 	err = os.WriteFile(filePath, []byte(params.Content), 0o644)
// 	if err != nil {
// 		return ToolResponse{}, fmt.Errorf("error writing file: %w", err)
// 	}

// 	// Check if file exists in history
// 	file, err := w.files.GetByPathAndSession(ctx, filePath, sessionID)
// 	if err != nil {
// 		_, err = w.files.Create(ctx, sessionID, filePath, oldContent)
// 		if err != nil {
// 			// Log error but don't fail the operation
// 			return ToolResponse{}, fmt.Errorf("error creating file history: %w", err)
// 		}
// 	}
// 	if file != nil && file.Content != oldContent {
// 		// User Manually changed the content store an intermediate version
// 		_, err = w.files.CreateVersion(ctx, sessionID, filePath, oldContent)
// 		if err != nil {
// 			logging.Debug("Error creating file history version", "error", err)
// 		}
// 	}
// 	// Store the new version
// 	_, err = w.files.CreateVersion(ctx, sessionID, filePath, params.Content)
// 	if err != nil {
// 		logging.Debug("Error creating file history version", "error", err)
// 	}

// 	recordFileWrite(filePath)
// 	recordFileRead(filePath)
// 	waitForLspDiagnostics(ctx, filePath, w.lspClients)

// 	result := fmt.Sprintf("File successfully written: %s", filePath)
// 	result = fmt.Sprintf("<result>\n%s\n</result>", result)
// 	result += getDiagnostics(filePath, w.lspClients)
// 	return WithResponseMetadata(NewTextResponse(result),
// 		ConfirmWriteResponseMetadata{
// 			Diff:      diff,
// 			Additions: additions,
// 			Removals:  removals,
// 		},
// 	), nil
// }
