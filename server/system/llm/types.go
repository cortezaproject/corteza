package llm

type (
	Message struct {
		Role       string     `json:"role"`
		Content    string     `json:"content"`
		ToolCalls  []ToolCall `json:"toolCalls,omitempty"`
		ToolCallID string     `json:"toolCallID,omitempty"`
	}

	ToolCall struct {
		ID        string         `json:"id"`
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}

	Response struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finishReason"`
		Usage        Usage   `json:"usage"`
	}

	Usage struct {
		PromptTokens     int `json:"promptTokens"`
		CompletionTokens int `json:"completionTokens"`
		TotalTokens      int `json:"totalTokens"`
	}

	Tool struct {
		Type     string       `json:"type"`
		Function ToolFunction `json:"function"`
	}

	ToolFunction struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  any    `json:"parameters"`
	}
)
