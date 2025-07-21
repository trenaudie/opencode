package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/db"
	"github.com/opencode-ai/opencode/internal/llm/agent"
	"github.com/opencode-ai/opencode/internal/llm/tools"
	"github.com/opencode-ai/opencode/internal/logging"
	"github.com/opencode-ai/opencode/internal/lsp"
	"github.com/opencode-ai/opencode/internal/message"
	"github.com/opencode-ai/opencode/internal/session"
)

func TestRunAgentTools(t *testing.T) {
	var ctx = context.Background()
	
	// Generate unique session ID for this test run
	sessionID := fmt.Sprintf("test-session-%d", time.Now().UnixNano())
	
	var toolcall = tools.ToolCall{
		ID:    sessionID,
		Name:  "coder",
		Input: `{\"prompt\":\"Remove ALL code related to the person SVG (import, signals, refs, view.add JSX, and associated animation logic). This includes:\\n\\n- Remove the import of person SVG\\n- Remove personRef, personOpacity, personScale\\n- Remove the \u003cSVG ...personRef...\u003e (and comment) JSX block from view.add (lines 66-81 inclusive)\\n- Remove animation blocks: yield* all(...) for personOpacity/personScale, and the for loop that sets fill to white for personRef paths (lines 103-113)\\n\\nEnsure code compiles after removal, don't remove anything sun or sand-related.\"}`,
	}

	// Create temporary test directory
	testDir := filepath.Join(os.TempDir(), "opencode-test")
	if err := os.RemoveAll(testDir); err != nil {
		t.Fatalf("failed to remove existing test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Initialize test configuration with temporary database
	_, err := config.Load(testDir, false)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	
	// Connect to database using the app's setup
	dbConn, err := db.Connect()
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer dbConn.Close()
	
	// Create database querier
	queries := db.New(dbConn)

	// Initialize required services
	service := session.NewService(queries)
	session_obj, err := service.CreateTaskSession(ctx, toolcall.ID, "parent_session_test_id", "Test Session")
	if err != nil {
		t.Errorf("failed to create session: %v", err)
		return
	}
	actualSessionID := session_obj.ID
	messages := message.NewService(queries)
	
	// Create starter message
	starterMessage, err := messages.Create(ctx, actualSessionID, message.CreateMessageParams{
		Role: message.User,
		Parts: []message.ContentPart{
			message.TextContent{Text: "Starting agent task"},
		},
	})
	if err != nil {
		t.Errorf("failed to create starter message: %v", err)
		return
	}
	messageID := starterMessage.ID
	
	// Add IDs to context
	ctx = context.WithValue(ctx, tools.SessionIDContextKey, actualSessionID)
	ctx = context.WithValue(ctx, tools.MessageIDContextKey, messageID)
	
	lspclients := make(map[string]*lsp.Client)
	logging.Info("messages are ", "messages", messages)
	logging.Info("lspclients are ", "lspclients", lspclients)
	logging.Info("actualSessionID", "actualSessionID", actualSessionID)
	logging.Info("messageID", "messageID", messageID)
	coderAgentTool := agent.NewCoderAgentTool(
		service,
		messages,
		lspclients,
	)
	// func (b *agentTool) Run(ctx context.Context, call tools.ToolCall) (tools.ToolResponse, error) {
	response, err := coderAgentTool.Run(ctx, toolcall)
	logging.Info("Tool Response", "response", response, "error", err)

	// 	agent, err := agent.NewAgent("coder", agentTool.sessions, agentTool.messages, TaskCoder(b.lspClients))

	// 	if err != nil {
	// 		log.Fatalf("Error scraping SVGs: %v", err)
	// 	}

	// 	fmt.Printf("Successfully scraped %d SVG files:\n\n", len(svgs))

	//		for i, svg := range svgs {
	//			fmt.Printf("=== SVG %d ===\n", i+1)
	//			fmt.Printf("Length: %d characters\n", len(svg))
	//			fmt.Printf("Preview (first 200 chars): %s...\n\n", truncate(svg, 200))
	//		}
	//	}
}
