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

const CoderAgentDescription = `Coder Agent specifically designed to generate and return Motion Canvas types This agent has no access to additional tools and focuses solely on generating TypeScript code for Motion Canvas projects.
You will be provided a SPEC FORMAT for the MotionCanvas scenes, in a json format. 
If you are making a tool call to the coder agent, your tool response for the coder agent can be one of two types:
- WRITE : If the user wants to WRITE a new scene, you MUST provide a JSON specification, and provide a very detailed json output to the next agent, with optional comments. The next agent should have all the specifications it needs to generate the MotionCanvas scene from scratch. Only provide input code to the coder agent if it is the output of the VIEW tool, ie it is not code that you  . 
eg. 
   "{\n  \"prompt\": \"Create a Motion Canvas scene where a triangle-shaped pyramid (an equilateral triangle, currently upside-down) is smoothly flipped so its apex points upward, right side up. Scene style guidelines: \\n- Use only Rect, Node, and Path (NOT Layout) for containers and positioning.\\n- The triangle's points and orientation should be computed reactively with createSignal and createComputed, referencing the parent rect's dimensions for positioning and rotation—avoid hardcoded pixel values whenever possible.\\n- Background should be black via view.fill('#000').\\n- All imports from '@motion-canvas/2d' and '@motion-canvas/core' only.\\n- Animate the triangle flipping over by rotating or morphing the points from the upside-down position to the right-side-up position, using a tween with an ease-in-out timing.\\n- Overwrite the entire content of example.tsx.\\n\\nFull SPEC:\\n{\\n  \\\"title\\\": \\\"Flip Triangle Pyramid Upright\\\",\\n  \\\"description\\\": \\\"An equilateral triangle is shown upside-down in the center. It animates by flipping to point upward, using a rotation or direct vertex morph.\\\",\\n  \\\"sceneMetadata\\\": {\\n    \\\"background\\\": \\\"#000000\\\",\\n    \\\"canvasDefaults\\\": null,\\n    \\\"viewport\\\": null,\\n    \\\"canvasResolution\\\": null,\\n    \\\"other\\\": null\\n  },\\n  \\\"imports\\\": [\\n    ...",
The SPEC FORMAT that you will be given gives you many fields, but these are optional. Only add a field in the spec that you provide to the coder agent if you have a non-null string value for it that you specifically want. 

- UPDATE : If the user wants to UPDATE an existing scene, you MUST provide a simple instruction to the coder agent WITHOUT the spec. In that case, you MUST also provide the existing typescript code in the example.tsx file, so that the coder agent can understand what to update.
You can invent names of functions or attributes if you do not know them, but in that case add a comment to the spec or instruction to the coder agent, so that it can understand that you are unsure about the real name of the function or attribute.
`

func (c *coderAgentTool) Info() tools.ToolInfo {
	return tools.ToolInfo{
		Name:        CoderAgentToolName,
		Description: CoderAgentDescription,
		Parameters: map[string]any{
			"prompt": map[string]any{
				"type":        "string",
				"description": "The coding task for the coder agent to perform, such as creating a specific Motion Canvas animation or component",
			},
			"type": map[string]any{
				"type":        "string",
				"description": "<WRITE|UPDATE> - Specify whether this is a write request to create new code or an update request to modify existing code",
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
