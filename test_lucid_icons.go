package main

import (
	"context"

	"github.com/opencode-ai/opencode/internal/llm/tools"
	"github.com/opencode-ai/opencode/internal/logging"
)

var svg_repo_tool = &tools.SVGRepoTool{}
var ctx = context.Background()
var toolcall = tools.ToolCall{
	ID:    "test_svg_repo",
	Name:  "svg_repo",
	Input: `{"assets": ["hospital", "person", "dog"]}`,
}

func testSVGRepo() {
	var ToolResponse, err = svg_repo_tool.Run(ctx, toolcall)
	logging.Info("Tool Response", "response", ToolResponse, "error", err)

}

/*
response is this
{
  "Type": "text",
  "Content": "Downloaded 9 SVG icons to frontend/public directory",
  "Metadata": {
    "filepaths": [
      "frontend/public/hospital_1.svg"
      "frontend/public/hospital_2.svg",
      "frontend/public/hospital_3.svg",
      "frontend/public/person_1.svg",
      "frontend/public/person_2.svg",
      "frontend/public/person_3.svg",
      "frontend/public/dog_1.svg",
      "frontend/public/dog_2.svg",
      "frontend/public/dog_3.svg"
    ]
  },
  "IsError": false,
  "error": null
}
*/
