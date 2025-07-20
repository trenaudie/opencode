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

type LucidIconsTool struct{}

type LucidIconsParams struct {
	Assets []string `json:"assets"`
}

type LucidIconsResponseMetadata struct {
	FilesPaths []string `json:"filepaths"`
}

const (
	LucidIconsToolName = "lucid_icons"
	lucidIconsDescription = `Icon asset fetching tool that downloads SVG icons from SVG Repo based on asset names.

WHEN TO USE THIS TOOL:
- Use when you need to download icon assets for your frontend application
- Perfect for getting high-quality SVG icons from SVG Repo
- Use when building UI components that need visual icons

HOW TO USE:
- Provide a list of asset names (e.g., "hospital", "person", "dog")
- The tool will search SVG Repo and fetch corresponding SVG files
- Icons are saved to the frontend/public directory
- Returns relative file paths to the downloaded icons

FEATURES:
- Downloads SVG icons from SVG Repo (svgrepo.com)
- Saves icons to frontend/public directory for web usage
- Returns relative file paths for easy integration
- Handles multiple assets in a single request
- Fetches up to 3 variations per asset name

LIMITATIONS:
- Requires internet connection for fetching icons
- Limited to assets available on SVG Repo
- May not find exact matches for very specific queries

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
	
	// Create frontend/public directory if it doesn't exist
	publicDir := "frontend/public"
	if err := os.MkdirAll(publicDir, 0755); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to create directory: %s", err)), nil
	}
	
	var allFilepaths []string
	totalDownloaded := 0
	
	for _, asset := range params.Assets {
		logging.Info("Scraping SVGs for asset", "asset", asset)
		
		// Scrape up to 3 SVGs per asset
		svgs, err := scrapesvgs.ScrapeSVG(asset, 3)
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
	
	if totalDownloaded == 0 {
		return NewTextErrorResponse("No SVG icons were successfully downloaded"), nil
	}
	
	return WithResponseMetadata(
		NewTextResponse(fmt.Sprintf("Downloaded %d SVG icons to %s directory", totalDownloaded, publicDir)),
		LucidIconsResponseMetadata{
			FilesPaths: allFilepaths,
		},
	), nil
}
