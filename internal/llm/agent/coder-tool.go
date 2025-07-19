package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/llm/tools"
	"github.com/opencode-ai/opencode/internal/lsp"
	"github.com/opencode-ai/opencode/internal/message"
	"github.com/opencode-ai/opencode/internal/session"
)

type coderAgentTool struct {
	sessions   session.Service
	messages   message.Service
	lspClients map[string]*lsp.Client
}

const (
	CoderAgentToolName = "coder"
)

type CoderAgentParams struct {
	Prompt string `json:"prompt"`
}

func (c *coderAgentTool) Info() tools.ToolInfo {
	return tools.ToolInfo{
		Name:        CoderAgentToolName,
		Description: "Launch a coder agent specifically designed for writing Motion Canvas animations. This agent has no access to additional tools and focuses solely on generating TypeScript code for Motion Canvas projects. The coder agent follows Motion Canvas developer style guidelines and creates reactive, well-structured animation code.",
		Parameters: map[string]any{
			"prompt": map[string]any{
				"type":        "string",
				"description": "The coding task for the coder agent to perform, such as creating a specific Motion Canvas animation or component",
			},
		},
		Required: []string{"prompt"},
	}
}

func (c *coderAgentTool) Run(ctx context.Context, call tools.ToolCall) (tools.ToolResponse, error) {
	var params CoderAgentParams
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

	done, err := agent.Run(ctx, session.ID, params.Prompt)
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

func NewCoderAgentTool(
	Sessions session.Service,
	Messages message.Service,
	LspClients map[string]*lsp.Client,
) tools.BaseTool {
	return &coderAgentTool{
		sessions:   Sessions,
		messages:   Messages,
		lspClients: LspClients,
	}
}