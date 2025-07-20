package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opencode-ai/opencode/internal/logging"
	scrapesvgs "github.com/opencode-ai/opencode/scrape_svgs"
)

type SVGRepoTool struct{}

type SVGRepoParams struct {
	Assets []string `json:"assets"`
}

type SVGRepoResponseMetadata struct {
	FilesPaths []string `json:"filepaths"`
}

const (
	SVGRepoToolName    = "svg_repo_fetcher"
	SVGRepoDescription = `Icon asset fetching tool that downloads SVG icons from SVG Repo based on asset names.

WHEN TO USE THIS TOOL:
- Use when you need to download icon assets for your frontend application
- Perfect for getting high-quality SVG icons from SVG Repo
- But, currently, I would like you to use this tool to fetch ANY external asset. This is the only currently available tool for fetching external assets, visual icons or image.

HOW TO USE the Tool:
- Provide a list of asset names (e.g., "hospital", "person", "dog")
- The tool will search SVG Repo and fetch corresponding SVG files
- Icons are saved to the frontend/public directory
- Returns relative file paths to the downloaded icons

HOW TO USE the Tool's output:
- The tool returns a list of file paths, with only one file path.
- You must take this **EXACT FILE PATH** and use it in the call to the Coder Agent tool, be for an update or a write request.

FEATURES:
- Downloads SVG icons from SVG Repo (svgrepo.com)
- Saves icons to frontend/public directory for web usage
- Returns relative file paths for easy integration
- Handles multiple assets in a single request
- Fetches up to 3 variations per asset name
`
)

func NewSVGRepoTool() BaseTool {
	return &SVGRepoTool{}
}

func (l *SVGRepoTool) Info() ToolInfo {
	return ToolInfo{
		Name:        SVGRepoToolName,
		Description: SVGRepoDescription,
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

func (l *SVGRepoTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params SVGRepoParams
	logging.Debug("svg repo tool params", "params", call.Input)
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}

	if len(params.Assets) == 0 {
		return NewTextErrorResponse("assets list is required"), nil
	}

	logging.Info("Received SVGRepoTool call with assets", "assets", params.Assets)

	// Create frontend/public directory if it doesn't exist
	publicDir := "frontend/public"
	if err := os.MkdirAll(publicDir, 0755); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to create directory: %s", err)), nil
	}

	var allFilepaths []string
	totalDownloaded := 0

	for _, asset := range params.Assets {
		logging.Info("Scraping SVGs for asset", "asset", asset)

		// Scrape up to 1 SVGs per asset
		svgs, err := scrapesvgs.ScrapeSVG(asset, 1)
		if err != nil {
			logging.Info("Failed to scrape SVGs for asset", "asset", asset, "error", err)
			continue
		}

		for i, svgContent := range svgs {
			filename := fmt.Sprintf("%s_%d.svg", strings.ReplaceAll(asset, " ", "_"), i+1)
			filepath := filepath.Join(publicDir, filename)

			if err := os.WriteFile(filepath, []byte(svgContent), 0644); err != nil {
				logging.Info("Failed to write SVG file", "filepath", filepath, "error", err)
				continue
			}

			allFilepaths = append(allFilepaths, filepath)
			totalDownloaded++
			logging.Info("Successfully saved SVG", "filepath", filepath)
		}
	}
	var FilepathListDebug []string
	FilepathListDebug = append(FilepathListDebug, allFilepaths[0])
	if totalDownloaded == 0 {
		return NewTextErrorResponse("No SVG icons were successfully downloaded"), nil
	}
	var toolResponse ToolResponse
	toolResponse.Type = ToolResponseTypeText

	responseData := map[string]any{
		"filepaths": FilepathListDebug,
		"text":      fmt.Sprintf("Downloaded %d SVG icons to %s directory", totalDownloaded, publicDir),
	}

	contentBytes, err := json.Marshal(responseData)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("error marshaling response: %s", err)), nil
	}
	toolResponse.Content = string(contentBytes)

	metadataBytes, err := json.Marshal(SVGRepoResponseMetadata{
		FilesPaths: FilepathListDebug,
	})
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("error marshaling metadata: %s", err)), nil
	}
	toolResponse.Metadata = string(metadataBytes)
	toolResponse.IsError = false

	return toolResponse, nil
}
