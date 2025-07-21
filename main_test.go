package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/opencode-ai/opencode/internal/app"
	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/db"
	"github.com/opencode-ai/opencode/internal/llm/agent"
	"github.com/opencode-ai/opencode/internal/logging"
)

func TestMain(m *testing.M) {
	logging.InitGlobalLogging("test.log")
	
	if err := logging.InitSystemPromptLogging(); err != nil {
		log.Fatal("Failed to initialize system prompt logging:", err)
	}
	
	os.Exit(m.Run())
}

func TestIntegrationAgentSession(t *testing.T) {
	if len(os.Args) < 2 {
		t.Skip("No prompt provided as command line argument. Usage: go test -run TestIntegrationAgentSession -args \"your prompt here\"")
	}
	
	prompt := os.Args[len(os.Args)-1]
	if prompt == "" {
		t.Fatal("Empty prompt provided")
	}
	
	fmt.Printf("Running integration test with prompt: %s\n", truncateString(prompt, 100))
	
	ctx := context.Background()
	
	appInstance, err := initializeTestApp(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize app: %v", err)
	}
	defer appInstance.Shutdown()
	
	server := &ChatServer{
		app: appInstance,
	}
	
	err = server.runAgentSession(ctx, prompt)
	if err != nil {
		t.Fatalf("Agent session failed: %v", err)
	}
}

func initializeTestApp(ctx context.Context) (*app.App, error) {
	cwd, err := os.Getwd()
	if err != nil {
		logging.Error("Failed to get current working directory", "error", err)
		return nil, err
	}
	logging.Info("Current working directory obtained", "cwd", cwd)

	_, err = config.Load(cwd, false)
	if err != nil {
		logging.Error("Failed to load configuration", "error", err, "cwd", cwd)
		return nil, err
	}
	
	dbConn, err := db.Connect()
	if err != nil {
		logging.Error("Failed to connect to database", "error", err)
		return nil, err
	}
	logging.Info("Database connection established successfully")

	appInstance, err := app.New(ctx, dbConn)
	if err != nil {
		logging.Error("Failed to create app instance", "error", err)
		return nil, err
	}
	logging.Info("App instance created successfully")

	return appInstance, nil
}

func (s *ChatServer) runAgentSession(ctx context.Context, prompt string) error {
	logging.Debug("Starting agent session integration test")
	
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		logging.Info("Cancelling agent session context")
		cancel()
	}()

	if prompt == "" {
		return fmt.Errorf("prompt is required")
	}

	logging.Debug("Received prompt", "prompt_length", len(prompt), "prompt_preview", truncateString(prompt, 100))

	logging.Debug("Creating new session for Motion Canvas Scene Generation")
	session, err := s.app.Sessions.Create(ctx, "Motion Canvas Scene Generation")
	if err != nil {
		logging.Error("Failed to create session", "error", err)
		return fmt.Errorf("failed to create session: %w", err)
	}

	fmt.Printf("Session created successfully: %s\n", session.ID)
	
	s.app.Permissions.AutoApproveSession(session.ID)
	logging.Debug("Auto-approved permissions for session", "session_id", session.ID)

	logging.Debug("Starting CoderAgent for session", "session_id", session.ID, "prompt", truncateString(prompt, 200))
	done, err := s.app.CoderAgent.Run(ctx, session.ID, prompt)
	if err != nil {
		logging.Error("Failed to start CoderAgent", "error", err, "session_id", session.ID)
		return fmt.Errorf("failed to start agent: %w", err)
	}
	logging.Debug("CoderAgent started successfully", "session_id", session.ID)

	logging.Debug("Subscribing to agent events", "session_id", session.ID)
	eventChan := s.app.CoderAgent.Subscribe(ctx)
	logging.Debug("Successfully subscribed to agent events", "session_id", session.ID)

	fmt.Printf("Processing agent events for session: %s\n", session.ID)
	
	logging.Debug("Starting event loop for session", "session_id", session.ID)
	for {
		select {
		case event := <-eventChan:
			logging.Debug("Received agent event", "event_type", event.Type, "session_id", event.Payload.SessionID, "target_session", session.ID)
			if event.Payload.SessionID == session.ID {
				logging.Debug("Processing agent event for our session", "event_type", event.Type, "session_id", session.ID)
				s.handleAgentEventConsole(event.Payload)
			} else {
				logging.Debug("Ignoring event for different session", "event_session", event.Payload.SessionID, "our_session", session.ID)
			}
		case result := <-done:
			logging.Debug("Agent processing completed", "session_id", session.ID, "has_error", result.Error != nil)
			if result.Error != nil {
				logging.Error("Agent completed with error", "error", result.Error, "session_id", session.ID)
				fmt.Printf("❌ Agent error: %s\n", result.Error.Error())
				return result.Error
			} else {
				logging.Debug("Agent completed successfully", "session_id", session.ID, "content_length", len(result.Message.Content().String()))
				fmt.Printf("✅ Agent completed successfully\n")
				fmt.Printf("📄 Final response:\n%s\n", result.Message.Content().String())
				fmt.Printf("🎯 Session completed: %s\n", session.ID)
			}
			logging.Debug("Agent session integration test completed", "session_id", session.ID)
			return nil
		case <-ctx.Done():
			logging.Debug("Context cancelled, ending agent session", "session_id", session.ID)
			return ctx.Err()
		case <-time.After(5 * time.Minute):
			logging.Warn("Agent session timed out after 5 minutes", "session_id", session.ID)
			return fmt.Errorf("agent session timed out")
		}
	}
}

func (s *ChatServer) handleAgentEventConsole(event agent.AgentEvent) {
	logging.Debug("Handling agent event", "event_type", event.Type, "session_id", event.SessionID)
	switch event.Type {
	case agent.AgentEventTypeResponse:
		logging.Info("Processing agent response event", "session_id", event.SessionID, "content_length", len(event.Message.Content().String()))
		fmt.Printf("🤖 Agent response: %s\n", truncateString(event.Message.Content().String(), 200))
	case agent.AgentEventTypeError:
		logging.Error("Processing agent error event", "error", event.Error, "session_id", event.SessionID)
		fmt.Printf("❌ Agent error: %s\n", event.Error.Error())
	default:
		logging.Warn("Unknown agent event type", "event_type", event.Type, "session_id", event.SessionID)
		fmt.Printf("❓ Unknown event type: %s\n", event.Type)
	}
}