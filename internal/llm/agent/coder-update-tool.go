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
	"github.com/opencode-ai/opencode/internal/session"
)

type coderUpdateAgentTool struct {
	sessions   session.Service
	messages   message.Service
	lspClients map[string]*lsp.Client
}

const (
	CoderUpdateAgentToolName = "coder"
)

type CoderUpdateAgentParams struct {
	Prompt string `json:"prompt"`
}

const CoderUpdateAgentDescription = `Coder Update Agent specifically designed to EDIT Motion Canvas scene code, written in Typescript. This agent has no access to additional tools and focuses solely on generating TypeScript code for an existing Motion Canvas scene.
WHEN TO USE THIS TOOL:
If the user wants to UPDATE an existing scene, you MUST provide a simple instruction to the coder update agent, telling it specifically what to change, what to remove, and/or what to add. 
You may give code to the update agent, if you have previously used a view tool, and have knowledge of the scene (example.tsx)
In your instructions can invent names of functions or attributes if you do not know them, but in that case add a comment to the spec or instruction to the coder agent, so that it can understand that you are unsure about the real name of the function or attribute.
HOW TO USE THIS TOOL:
- provide a simple prompt to the Coder Update Agent, regarding what to change. Explain it in concise and precise text. 
`

func (c *coderUpdateAgentTool) Info() tools.ToolInfo {
	return tools.ToolInfo{
		Name:        CoderUpdateAgentToolName,
		Description: CoderUpdateAgentDescription,
		Parameters: map[string]any{
			"prompt": map[string]any{
				"type":        "string",
				"description": "The coding task for the coder agent to perform, such as creating a specific Motion Canvas animation or component",
			},
			"code_to_modify": map[string]any{
				"type":        "string",
				"description": "The Motion Canvas scene, written in typescript code, to edit. Most of the code may be working, only an edit of some of the code's lines is necessary.",
			},
		},
		Required: []string{"prompt"},
	}
}

func (c *coderUpdateAgentTool) Run(ctx context.Context, call tools.ToolCall) (tools.ToolResponse, error) {
	var params CoderUpdateAgentParams
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

	// Create coder agent with no tools (empty tools slice)
	agent, err := NewAgent(config.AgentCoder, c.sessions, c.messages, []tools.BaseTool{})
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error creating coder agent: %s", err)
	}

	session, err := c.sessions.CreateTaskSession(ctx, call.ID, sessionID, "Coder Agent Session")
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
		return tools.ToolResponse{}, fmt.Errorf("error running coder agent: %s", err)
	}
	result := <-done
	if result.Error != nil {
		return tools.ToolResponse{}, fmt.Errorf("error from coder agent: %s", result.Error)
	}

	response := result.Message
	if response.Role != message.Assistant {
		return tools.NewTextErrorResponse("no response from coder agent"), nil
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

func NewCoderUpdateAgentTool(
	Sessions session.Service,
	Messages message.Service,
	LspClients map[string]*lsp.Client,
) tools.BaseTool {
	return &coderUpdateAgentTool{
		sessions:   Sessions,
		messages:   Messages,
		lspClients: LspClients,
	}
}
