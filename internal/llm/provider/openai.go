package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/llm/models"
	"github.com/opencode-ai/opencode/internal/llm/tools"
	"github.com/opencode-ai/opencode/internal/logging"
	"github.com/opencode-ai/opencode/internal/message"
)

type openaiOptions struct {
	baseURL         string
	disableCache    bool
	reasoningEffort string
	extraHeaders    map[string]string
}

type OpenAIOption func(*openaiOptions)

type openaiClient struct {
	providerOptions providerClientOptions
	options         openaiOptions
	client          openai.Client
}

type OpenAIClient ProviderClient

func newOpenAIClient(opts providerClientOptions) OpenAIClient {
	openaiOpts := openaiOptions{
		reasoningEffort: "medium",
	}
	for _, o := range opts.openaiOptions {
		o(&openaiOpts)
	}

	openaiClientOptions := []option.RequestOption{}
	if opts.apiKey != "" {
		openaiClientOptions = append(openaiClientOptions, option.WithAPIKey(opts.apiKey))
	}
	if openaiOpts.baseURL != "" {
		openaiClientOptions = append(openaiClientOptions, option.WithBaseURL(openaiOpts.baseURL))
	}

	if openaiOpts.extraHeaders != nil {
		for key, value := range openaiOpts.extraHeaders {
			openaiClientOptions = append(openaiClientOptions, option.WithHeader(key, value))
		}
	}

	client := openai.NewClient(openaiClientOptions...)
	return &openaiClient{
		providerOptions: opts,
		options:         openaiOpts,
		client:          client,
	}
}

func (o *openaiClient) convertMessages(messages []message.Message) (openaiMessages []openai.ChatCompletionMessageParamUnion) {
	// Add system message first
	// logging.Info("System message", "content", o.providerOptions.systemMessage)
	openaiMessages = append(openaiMessages, openai.SystemMessage(o.providerOptions.systemMessage))

	for _, msg := range messages {
		switch msg.Role {
		case message.User:
			var content []openai.ChatCompletionContentPartUnionParam
			textBlock := openai.ChatCompletionContentPartTextParam{Text: msg.Content().String()}
			content = append(content, openai.ChatCompletionContentPartUnionParam{OfText: &textBlock})
			for _, binaryContent := range msg.BinaryContent() {
				imageURL := openai.ChatCompletionContentPartImageImageURLParam{URL: binaryContent.String(models.ProviderOpenAI)}
				imageBlock := openai.ChatCompletionContentPartImageParam{ImageURL: imageURL}

				content = append(content, openai.ChatCompletionContentPartUnionParam{OfImageURL: &imageBlock})
			}

			openaiMessages = append(openaiMessages, openai.UserMessage(content))

		case message.Assistant:
			assistantMsg := openai.ChatCompletionAssistantMessageParam{
				Role: "assistant",
			}

			if msg.Content().String() != "" {
				assistantMsg.Content = openai.ChatCompletionAssistantMessageParamContentUnion{
					OfString: openai.String(msg.Content().String()),
				}
			}

			if len(msg.ToolCalls()) > 0 {
				assistantMsg.ToolCalls = make([]openai.ChatCompletionMessageToolCallParam, len(msg.ToolCalls()))
				for i, call := range msg.ToolCalls() {
					assistantMsg.ToolCalls[i] = openai.ChatCompletionMessageToolCallParam{
						ID:   call.ID,
						Type: "function",
						Function: openai.ChatCompletionMessageToolCallFunctionParam{
							Name:      call.Name,
							Arguments: call.Input,
						},
					}
				}
			}

			openaiMessages = append(openaiMessages, openai.ChatCompletionMessageParamUnion{
				OfAssistant: &assistantMsg,
			})

		case message.Tool:
			for _, result := range msg.ToolResults() {
				openaiMessages = append(openaiMessages,
					openai.ToolMessage(result.Content, result.ToolCallID),
				)
			}
		}
	}

	return
}

func (o *openaiClient) convertTools(tools []tools.BaseTool) []openai.ChatCompletionToolParam {
	openaiTools := make([]openai.ChatCompletionToolParam, len(tools))

	for i, tool := range tools {
		info := tool.Info()
		openaiTools[i] = openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        info.Name,
				Description: openai.String(info.Description),
				Parameters: openai.FunctionParameters{
					"type":       "object",
					"properties": info.Parameters,
					"required":   info.Required,
				},
			},
		}
	}

	return openaiTools
}

func (o *openaiClient) finishReason(reason string) message.FinishReason {
	switch reason {
	case "stop":
		return message.FinishReasonEndTurn
	case "length":
		return message.FinishReasonMaxTokens
	case "tool_calls":
		return message.FinishReasonToolUse
	default:
		return message.FinishReasonUnknown
	}
}

func (o *openaiClient) preparedParams(messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolParam) openai.ChatCompletionNewParams {
	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(o.providerOptions.model.APIModel),
		Messages: messages,
		Tools:    tools,
	}

	if o.providerOptions.model.CanReason == true {
		params.MaxCompletionTokens = openai.Int(o.providerOptions.maxTokens)
		switch o.options.reasoningEffort {
		case "low":
			params.ReasoningEffort = shared.ReasoningEffortLow
		case "medium":
			params.ReasoningEffort = shared.ReasoningEffortMedium
		case "high":
			params.ReasoningEffort = shared.ReasoningEffortHigh
		default:
			params.ReasoningEffort = shared.ReasoningEffortMedium
		}
	} else {
		params.MaxTokens = openai.Int(o.providerOptions.maxTokens)
	}

	return params
}

func (o *openaiClient) send(ctx context.Context, messages []message.Message, tools []tools.BaseTool) (response *ProviderResponse, err error) {
	params := o.preparedParams(o.convertMessages(messages), o.convertTools(tools))
	cfg := config.Get()
	if cfg.Debug {
		jsonData, _ := json.Marshal(params)
		logging.Info("Prepared messages", "messages", string(jsonData))
	}
	attempts := 0
	for {
		attempts++
		// logging.Info("Making OpenAI API call", "model", o.providerOptions.model.APIModel, "system_prompt", o.providerOptions.systemMessage[:1000]+"...."+o.providerOptions.systemMessage[len(o.providerOptions.systemMessage)-1000:], "attempt", attempts, "system_prompt_length", len(o.providerOptions.systemMessage))

		logging.Info("Making OpenAI API call", "model", o.providerOptions.model.APIModel, "attempt", attempts, "system_prompt_length", len(o.providerOptions.systemMessage))
		// Log system prompt and input to files
		agentName := o.extractAgentName()
		logging.LogSystemPrompt(string(agentName), o.providerOptions.systemMessage)

		// Log input content
		inputContent := o.extractInputContent(messages)
		logging.LogInput(string(agentName), inputContent, o.providerOptions.systemMessage)

		// Prepare input data for agent call logging
		inputData := map[string]interface{}{
			"model":                o.providerOptions.model.APIModel,
			"system_prompt_length": len(o.providerOptions.systemMessage),
			"messages_count":       len(messages),
			"tools_count":          len(tools),
			"max_tokens":           o.providerOptions.maxTokens,
			"attempt":              attempts,
		}

		for _, msg := range messages {
			logging.Info("Processing message", "role", msg.Role, "content_length", len(msg.Content().Text), "tool_calls_count", len(msg.ToolCalls()), "tool_results_count", len(msg.ToolResults()))
		}

		openaiResponse, err := o.client.Chat.Completions.New(
			ctx,
			params,
		)
		if err == nil && openaiResponse != nil && len(openaiResponse.Choices) > 0 {
			logging.Info("OpenAI API call completed", "content_length", len(openaiResponse.Choices[0].Message.Content), "finish_reason", openaiResponse.Choices[0].FinishReason, "tool_calls_count", len(openaiResponse.Choices[0].Message.ToolCalls))
		} else {
			logging.Info("OpenAI API call failed", "error", err)
		}
		// If there is an error we are going to see if we can retry the call
		if err != nil {
			retry, after, retryErr := o.shouldRetry(attempts, err)
			if retryErr != nil {
				return nil, retryErr
			}
			if retry {
				logging.WarnPersist(fmt.Sprintf("Retrying due to rate limit... attempt %d of %d", attempts, maxRetries), logging.PersistTimeArg, time.Millisecond*time.Duration(after+100))
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(time.Duration(after) * time.Millisecond):
					continue
				}
			}
			return nil, retryErr
		}

		content := ""
		if openaiResponse.Choices[0].Message.Content != "" {
			content = openaiResponse.Choices[0].Message.Content
		}

		toolCalls := o.toolCalls(*openaiResponse)
		finishReason := o.finishReason(string(openaiResponse.Choices[0].FinishReason))

		if len(toolCalls) > 0 {
			finishReason = message.FinishReasonToolUse
		}

		// Log output content and tool calls
		logging.LogOutput(string(agentName), content, toolCalls)

		// Prepare output data for agent call logging
		outputData := map[string]interface{}{
			"content_length":   len(content),
			"finish_reason":    string(openaiResponse.Choices[0].FinishReason),
			"tool_calls_count": len(toolCalls),
			"usage":            o.usage(*openaiResponse),
		}

		// Log the complete agent call
		logging.LogAgentCall(string(agentName), "openai_api_call", inputData, outputData)

		return &ProviderResponse{
			Content:      content,
			ToolCalls:    toolCalls,
			Usage:        o.usage(*openaiResponse),
			FinishReason: finishReason,
		}, nil
	}
}

func (o *openaiClient) stream(ctx context.Context, messages []message.Message, tools []tools.BaseTool) <-chan ProviderEvent {
	params := o.preparedParams(o.convertMessages(messages), o.convertTools(tools))
	params.StreamOptions = openai.ChatCompletionStreamOptionsParam{
		IncludeUsage: openai.Bool(true),
	}

	cfg := config.Get()
	if cfg.Debug {
		jsonData, _ := json.Marshal(params)
		logging.Debug("Prepared messages", "messages", string(jsonData))
	}

	attempts := 0
	eventChan := make(chan ProviderEvent)

	go func() {
		for {
			attempts++
			// logging.Info("Making OpenAI streaming API call", "model", o.providerOptions.model.APIModel, "attempt", attempts, "system_prompt", o.providerOptions.systemMessage[:1000]+"...."+o.providerOptions.systemMessage[len(o.providerOptions.systemMessage)-1000:], "system_prompt_length", len(o.providerOptions.systemMessage))
			logging.Info("Making OpenAI API call", "model", o.providerOptions.model.APIModel, "attempt", attempts, "system_prompt_length", len(o.providerOptions.systemMessage))
			// Log system prompt and input to files
			agentName := o.extractAgentName()
			logging.LogSystemPrompt(string(agentName), o.providerOptions.systemMessage)

			// Log input content
			inputContent := o.extractInputContent(messages)
			logging.LogInput(string(agentName), inputContent, o.providerOptions.systemMessage)

			// Prepare input data for agent call logging (streaming)
			inputData := map[string]interface{}{
				"model":                o.providerOptions.model.APIModel,
				"system_prompt_length": len(o.providerOptions.systemMessage),
				"messages_count":       len(messages),
				"tools_count":          len(tools),
				"max_tokens":           o.providerOptions.maxTokens,
				"attempt":              attempts,
				"streaming":            true,
			}

			for _, msg := range messages {
				logging.Info("Processing streaming message", "role", msg.Role, "content_length", len(msg.Content().Text), "tool_calls_count", len(msg.ToolCalls()), "tool_results_count", len(msg.ToolResults()))
			}
			openaiStream := o.client.Chat.Completions.NewStreaming(
				ctx,
				params,
			)

			acc := openai.ChatCompletionAccumulator{}
			currentContent := ""
			toolCalls := make([]message.ToolCall, 0)

			for openaiStream.Next() {
				chunk := openaiStream.Current()
				acc.AddChunk(chunk)

				for _, choice := range chunk.Choices {
					if choice.Delta.Content != "" {
						eventChan <- ProviderEvent{
							Type:    EventContentDelta,
							Content: choice.Delta.Content,
						}
						currentContent += choice.Delta.Content
					}
				}
			}
			logging.Info("finished streaming , so here is the content full:", currentContent)
			jsonData, err1 := json.Marshal(acc.ChatCompletion)
			if err1 != nil {
				logging.Error("Failed to marshal JSON", "error", err1)
			} else {
				logging.Info("JSON marshal of the ChatCompletion", "json", string(jsonData))
			}

			err := openaiStream.Err()
			if err == nil || errors.Is(err, io.EOF) {
				// Stream completed successfully
				finishReason := o.finishReason(string(acc.ChatCompletion.Choices[0].FinishReason))
				if len(acc.ChatCompletion.Choices[0].Message.ToolCalls) > 0 {
					toolCalls = append(toolCalls, o.toolCalls(acc.ChatCompletion)...)
				}
				if len(toolCalls) > 0 {
					finishReason = message.FinishReasonToolUse
				}

				logging.Info("OpenAI streaming API call completed", "content_length", len(currentContent), "finish_reason", finishReason, "tool_calls_count", len(toolCalls))

				// Log output content and tool calls
				logging.LogOutput(string(agentName), currentContent, toolCalls)

				// Prepare output data for agent call logging (streaming)
				outputData := map[string]interface{}{
					"content_length":   len(currentContent),
					"finish_reason":    string(acc.ChatCompletion.Choices[0].FinishReason),
					"tool_calls_count": len(toolCalls),
					"usage":            o.usage(acc.ChatCompletion),
					"streaming":        true,
				}

				// Log the complete agent call (streaming)
				logging.LogAgentCall(string(agentName), "openai_streaming_api_call", inputData, outputData)

				eventChan <- ProviderEvent{
					Type: EventComplete,
					Response: &ProviderResponse{
						Content:      currentContent,
						ToolCalls:    toolCalls,
						Usage:        o.usage(acc.ChatCompletion),
						FinishReason: finishReason,
					},
				}
				close(eventChan)
				return
			}

			// If there is an error we are going to see if we can retry the call
			retry, after, retryErr := o.shouldRetry(attempts, err)
			if retryErr != nil {
				eventChan <- ProviderEvent{Type: EventError, Error: retryErr}
				close(eventChan)
				return
			}
			if retry {
				logging.WarnPersist(fmt.Sprintf("Retrying due to rate limit... attempt %d of %d", attempts, maxRetries), logging.PersistTimeArg, time.Millisecond*time.Duration(after+100))
				select {
				case <-ctx.Done():
					// context cancelled
					if ctx.Err() == nil {
						eventChan <- ProviderEvent{Type: EventError, Error: ctx.Err()}
					}
					close(eventChan)
					return
				case <-time.After(time.Duration(after) * time.Millisecond):
					continue
				}
			}
			eventChan <- ProviderEvent{Type: EventError, Error: retryErr}
			close(eventChan)
			return
		}
	}()

	return eventChan
}

func (o *openaiClient) shouldRetry(attempts int, err error) (bool, int64, error) {
	var apierr *openai.Error
	if !errors.As(err, &apierr) {
		return false, 0, err
	}

	if apierr.StatusCode != 429 && apierr.StatusCode != 500 {
		return false, 0, err
	}

	if attempts > maxRetries {
		return false, 0, fmt.Errorf("maximum retry attempts reached for rate limit: %d retries", maxRetries)
	}

	retryMs := 0
	retryAfterValues := apierr.Response.Header.Values("Retry-After")

	backoffMs := 2000 * (1 << (attempts - 1))
	jitterMs := int(float64(backoffMs) * 0.2)
	retryMs = backoffMs + jitterMs
	if len(retryAfterValues) > 0 {
		if _, err := fmt.Sscanf(retryAfterValues[0], "%d", &retryMs); err == nil {
			retryMs = retryMs * 1000
		}
	}
	return true, int64(retryMs), nil
}

func (o *openaiClient) toolCalls(completion openai.ChatCompletion) []message.ToolCall {
	var toolCalls []message.ToolCall

	if len(completion.Choices) > 0 && len(completion.Choices[0].Message.ToolCalls) > 0 {
		for _, call := range completion.Choices[0].Message.ToolCalls {
			toolCall := message.ToolCall{
				ID:       call.ID,
				Name:     call.Function.Name,
				Input:    call.Function.Arguments,
				Type:     "function",
				Finished: true,
			}
			toolCalls = append(toolCalls, toolCall)
		}
	}

	return toolCalls
}

func (o *openaiClient) usage(completion openai.ChatCompletion) TokenUsage {
	cachedTokens := completion.Usage.PromptTokensDetails.CachedTokens
	inputTokens := completion.Usage.PromptTokens - cachedTokens

	return TokenUsage{
		InputTokens:         inputTokens,
		OutputTokens:        completion.Usage.CompletionTokens,
		CacheCreationTokens: 0, // OpenAI doesn't provide this directly
		CacheReadTokens:     cachedTokens,
	}
}

func WithOpenAIBaseURL(baseURL string) OpenAIOption {
	return func(options *openaiOptions) {
		options.baseURL = baseURL
	}
}

func WithOpenAIExtraHeaders(headers map[string]string) OpenAIOption {
	return func(options *openaiOptions) {
		options.extraHeaders = headers
	}
}

func WithOpenAIDisableCache() OpenAIOption {
	return func(options *openaiOptions) {
		options.disableCache = true
	}
}

func WithReasoningEffort(effort string) OpenAIOption {
	return func(options *openaiOptions) {
		defaultReasoningEffort := "medium"
		switch effort {
		case "low", "medium", "high":
			defaultReasoningEffort = effort
		default:
			logging.Warn("Invalid reasoning effort, using default: medium")
		}
		options.reasoningEffort = defaultReasoningEffort
	}
}

// extractAgentName attempts to determine the agent name from the system message
func (o *openaiClient) extractAgentName() config.AgentName {
	sysMsg := strings.ToLower(o.providerOptions.systemMessage)
	// Check for orchestrator first since it may also contain "code" words
	if strings.Contains(sysMsg, "orchestrator") {
		return config.AgentOrchestrator
	}
	if strings.Contains(sysMsg, "coder") || strings.Contains(sysMsg, "code") {
		return config.AgentCoder
	}
	if strings.Contains(sysMsg, "title") {
		return config.AgentTitle
	}
	if strings.Contains(sysMsg, "summariz") {
		return config.AgentSummarizer
	}
	if strings.Contains(sysMsg, "task") {
		return config.AgentTask
	}
	return "unknown"
}

// extractInputContent combines all user messages into a single input string for logging
func (o *openaiClient) extractInputContent(messages []message.Message) string {
	var inputParts []string

	for _, msg := range messages {
		switch msg.Role {
		case message.User:
			if content := msg.Content().String(); content != "" {
				inputParts = append(inputParts, fmt.Sprintf("User: %s", content))
			}
		case message.Tool:
			for _, result := range msg.ToolResults() {
				inputParts = append(inputParts, fmt.Sprintf("Tool Result (%s): %s", result.ToolCallID, result.Content))
			}
		}
	}

	return strings.Join(inputParts, "\n\n")
}
