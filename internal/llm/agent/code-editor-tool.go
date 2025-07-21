package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/llm/tools"
	"github.com/opencode-ai/opencode/internal/lsp"
	"github.com/opencode-ai/opencode/internal/message"
	"github.com/opencode-ai/opencode/internal/permission"
	"github.com/opencode-ai/opencode/internal/session"
)

type codeEditorAgentTool struct {
	sessions   session.Service
	messages   message.Service
	lspClients map[string]*lsp.Client
}

const (
	CodeEditorAgentToolName = "code_editor"
)

type CodeEditorAgentParams struct {
	Prompt string `json:"prompt"`
}

const CodeEditorAgentDescription = `Code Editor Agent specifically designed to EDIT Motion Canvas scene code IN-PLACE, written in Typescript. This is the EDITOR agent that modifies existing files by providing exact text replacements.
WHEN TO USE THIS TOOL:
Use this agent when the user wants to edit, modify, or update an existing Motion Canvas scene. This agent will output old_string (exact text to be replaced) and new_string (exact replacement text) that will be passed to a deterministic edit tool.
CRITICAL REQUIREMENTS - THE AGENT MUST OUTPUT EXACT STRINGS:
The orchestrator will use the agent's output with a deterministic edit tool that performs exact string replacement. Therefore:
- The old_string MUST match the existing file content EXACTLY (including all whitespace, indentation, line breaks)
- The new_string MUST be the exact replacement text
- The old_string must be UNIQUE in the file - include enough context (3-5 lines before/after) to ensure uniqueness
- If old_string appears multiple times, provide more surrounding context to make it unique
ACCURACY IS CRITICAL - if the strings don't match exactly, the edit will fail.
HOW TO USE THIS TOOL:
- Provide clear instructions about what specific code sections need to be changed
- The agent will analyze the current file and provide exact old_string/new_string pairs for replacement
- You can mention function/attribute names even if unsure - add comments so the agent understands uncertainty
`

func (c *codeEditorAgentTool) Info() tools.ToolInfo {
	return tools.ToolInfo{
		Name:        CodeEditorAgentToolName,
		Description: CodeEditorAgentDescription,
		Parameters: map[string]any{
			"prompt": map[string]any{
				"type":        "string",
				"description": "The coding task for the code editor agent to perform, such as removing this or replacing that line of the example.tsx file",
			},
			// "code_to_modify": map[string]any{
			// 	"type":        "string",
			// 	"description": "The Motion Canvas scene, written in typescript code, to edit. Most of the code may be working, only an edit of some of the code's lines is necessary.",
			// }, do not add this here because it would create unnecessary generations
		},
		Required: []string{"prompt"},
	}
}

func (c *codeEditorAgentTool) Run(ctx context.Context, call tools.ToolCall) (tools.ToolResponse, error) {
	var params CodeEditorAgentParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return tools.NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}
	if params.Prompt == "" {
		return tools.NewTextErrorResponse("prompt is required"), nil
	}

	sessionID, messageID := tools.GetContextValues(ctx)
	if sessionID == "" || messageID == "" {
		return tools.ToolResponse{}, fmt.Errorf("session_id and message_id are required")
	}

	// Create code-editor agent with edit tools
	permissions := permission.NewPermissionService()
	// For the history service, we'll pass nil since these agents run in isolated sessions
	// and don't need persistent history tracking
	agent, err := NewAgent(config.AgentCodeEditor, c.sessions, c.messages, CodeEditorAgentTools(permissions, nil, c.lspClients))
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error creating code-editor agent: %s", err)
	}

	session, err := c.sessions.CreateTaskSession(ctx, call.ID, sessionID, "CodeEditor Agent Session")
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error creating session: %s", err)
	}

	filePath := filepath.Join(config.WorkingDirectory(), "frontend/src/scenes/example.tsx")
	currentScene, err := os.ReadFile(filePath)
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error reading file: %s", err)
	}
	done, err := agent.Run(ctx, session.ID, params.Prompt+"\n\nCurrent file (called example.tsx):\n"+string(currentScene))
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error running code-editor agent: %s", err)
	}
	result := <-done
	if result.Error != nil {
		return tools.ToolResponse{}, fmt.Errorf("error from code-editor agent: %s", result.Error)
	}

	response := result.Message
	if response.Role != message.Assistant {
		return tools.NewTextErrorResponse("no response from code-editor agent"), nil
	}

	updatedSession, err := c.sessions.Get(ctx, session.ID)
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error getting session: %s", err)
	}
	parentSession, err := c.sessions.Get(ctx, sessionID)
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error getting parent session: %s", err)
	}

	parentSession.Cost += updatedSession.Cost

	_, err = c.sessions.Save(ctx, parentSession)
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error saving parent session: %s", err)
	}
	return tools.NewTextResponse(response.Content().String()), nil
}

func NewCodeEditorAgentTool(
	Sessions session.Service,
	Messages message.Service,
	LspClients map[string]*lsp.Client,
) tools.BaseTool {
	return &codeEditorAgentTool{
		sessions:   Sessions,
		messages:   Messages,
		lspClients: LspClients,
	}
}
