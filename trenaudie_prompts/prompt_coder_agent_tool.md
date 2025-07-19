
  For my Motion Canvas agent LLM project, I needed to restructure the tool architecture to implement a proper
  orchestrator-coder agent separation. The existing setup had the coder agent with all tools, but the architecture
  requires:

  1. Orchestrator Agent - Should have ALL tools for coordination and file operations
  2. Coder Agent - Should have NO tools, focusing purely on code generation based on context and instructions

  Why This Architecture?

  Separation of Concerns:
  - The orchestrator agent handles workflow coordination, file operations, viewing files, editing/writing code, calling
  external APIs, and delegating to other agents
  - The coder agent focuses exclusively on generating Motion Canvas TypeScript code based on the context and specific
  instructions it receives from the orchestrator

  Tool Distribution Rationale:
  - Orchestrator gets all tools because it needs to:
    - View existing code structure (ViewTool, LsTool, GrepTool)
    - Edit and write files (EditTool, WriteTool, PatchTool)
    - Execute commands (BashTool)
    - Call external services (FetchTool)
    - Call other agents including the coder (AgentTool)
    - Handle diagnostics (DiagnosticsTool)
  - Coder gets no tools because it should:
    - Focus purely on code generation without distractions
    - Rely on rich context provided by the orchestrator
    - Follow Motion Canvas patterns and conventions from examples
    - Generate clean, self-contained TypeScript code

  Implementation Details

  Files Modified:

  - @internal/llm/agent/tools.go - Restructured tool distribution functions
  - @internal/config/config.go - Added orchestrator agent to all provider defaults

  Changes Made:

  1. Tool Function Restructuring in tools.go:
  // NEW: OrchestratorAgentTools - Gets ALL the tools that were previously in CoderAgentTools
  func OrchestratorAgentTools(
      permissions permission.Service,
      sessions session.Service, 
      messages message.Service,
      history history.Service,
      lspClients map[string]*lsp.Client,
  ) []tools.BaseTool {
      // Contains: BashTool, EditTool, FetchTool, GlobTool, GrepTool, 
      // LsTool, SourcegraphTool, ViewTool, PatchTool, WriteTool, 
      // LucidIconsTool, AgentTool, plus MCP tools and diagnostics
  }

  // UPDATED: CoderAgentTools - Now returns empty slice
  func CoderAgentTools(
      permissions permission.Service,
      sessions session.Service,
      messages message.Service, 
      history history.Service,
      lspClients map[string]*lsp.Client,
  ) []tools.BaseTool {
      return []tools.BaseTool{} // No tools - pure code generation
  }

  2. Configuration Updates in config.go:
  Added agents.orchestrator.model defaults to all provider configurations:
  - Copilot: models.CopilotGPT4o
  - Anthropic: models.Claude4Sonnet
  - OpenAI: models.GPT41
  - Gemini: models.Gemini25
  - Groq: models.QWENQwq
  - OpenRouter: models.OpenRouterClaude37Sonnet
  - XAI: models.XAIGrok3Beta
  - AWS Bedrock: models.BedrockClaude37Sonnet
  - Azure: models.AzureGPT41
  - VertexAI: models.VertexAIGemini25

  Orchestrator Prompt Integration

  The orchestrator agent uses the prompt defined in @internal/llm/prompt/orchestrator.go, which emphasizes:
  - 4-step workflow: Planning → Viewing → Implementation → Completion
  - Tool delegation: Always calls coder agent for code generation
  - File operations: Uses view tool after coder response to understand how to write code to example.tsx
  - Motion Canvas specificity: Focus on frontend/src/scenes/example.tsx as the target file

  Coder Prompt Integration

  The coder agent uses the Motion Canvas-specific prompt that includes:
  - Detailed Motion Canvas coding guidelines
  - Import patterns from @motion-canvas/2d and @motion-canvas/core
  - Extensive examples of proper Motion Canvas syntax
  - No tool access - relies entirely on context and instructions

  Expected Workflow

  1. User Request → Orchestrator Agent
  2. Orchestrator analyzes request and creates plan
  3. Orchestrator calls Coder Agent with specific instructions
  4. Coder Agent generates Motion Canvas code (no tools, pure generation)
  5. Orchestrator receives code and uses view tool to understand current file
  6. Orchestrator uses edit or write tools to implement the generated code
  7. Orchestrator verifies and completes the animation

  This architecture ensures clean separation of concerns while maintaining the ability to generate sophisticated Motion
  Canvas animations through coordinated agent interaction.