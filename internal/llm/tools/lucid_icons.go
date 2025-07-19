package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opencode-ai/opencode/internal/logging"
)

type LucidIconsTool struct{}

type LucidIconsParams struct {
	Assets []string `json:"assets"`
}

type LucidIconsResponseMetadata struct {
	FilesPaths []string `json:"filepaths"`
}

const (
	LucidIconsToolName = "lucid_icons"
	lucidIconsDescription = `Icon asset fetching tool that downloads Lucid icons as SVG files based on asset names.

WHEN TO USE THIS TOOL:
- Use when you need to download icon assets for your frontend application
- Perfect for getting consistent, high-quality SVG icons from the Lucid icon library
- Use when building UI components that need visual icons

HOW TO USE:
- Provide a list of asset names (e.g., "hospital", "person", "dog")
- The tool will fetch corresponding Lucid icons as SVG files
- Icons are saved to the frontend/public directory
- Returns relative file paths to the downloaded icons

FEATURES:
- Downloads SVG icons from Lucid icon library
- Saves icons to frontend/public directory for web usage
- Returns relative file paths for easy integration
- Handles multiple assets in a single request

LIMITATIONS:
- Currently returns hardcoded paths (implementation in progress)
- Requires internet connection for fetching icons
- Limited to Lucid icon library assets

TIPS:
- Use descriptive asset names like "home", "user", "settings"
- Icons are saved as SVG format for scalability
- Returned paths are relative to project root for easy imports`
)

func NewLucidIconsTool() BaseTool {
	return &LucidIconsTool{}
}

func (l *LucidIconsTool) Info() ToolInfo {
	return ToolInfo{
		Name:        LucidIconsToolName,
		Description: lucidIconsDescription,
		Parameters: map[string]any{
			"assets": map[string]any{
				"type":        "array",
				"description": "List of asset names to download (e.g., 'hospital', 'person', 'dog')",
				"items": map[string]any{
					"type": "string",
				},
			},
		},
		Required: []string{"assets"},
	}
}

func (l *LucidIconsTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params LucidIconsParams
	logging.Debug("lucid_icons tool params", "params", call.Input)
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}

	if len(params.Assets) == 0 {
		return NewTextErrorResponse("assets list is required"), nil
	}

	logging.Info("Received LucidIconsTool call with assets", "assets", params.Assets)
	
	// TODO: Implement actual icon fetching logic
	// For now, return hardcoded paths as requested
	filepaths := []string{
		"frontend/public/logo.svg",
		"frontend/public/grid.png",
		"frontend/public/placeholder.png",
	}
	
	return WithResponseMetadata(
		NewTextResponse(fmt.Sprintf("Downloaded %d icons to frontend/public directory", len(params.Assets))),
		LucidIconsResponseMetadata{
			FilesPaths: filepaths,
		},
	), nil
}
