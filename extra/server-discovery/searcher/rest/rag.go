package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/extra/server-discovery/searcher/rest/request"
	"github.com/tmc/langchaingo/llms"
)

type (
	searchAPI interface {
		SearchResources(context.Context, *request.SearchResources) (interface{}, error)
	}

	ConversationalRAGChain struct {
		searchAPI searchAPI
		llm       llms.Model
		memory    []llms.ChatMessage
	}
)

func Rag(searchAPI searchAPI, llm llms.Model) *ConversationalRAGChain {
	return &ConversationalRAGChain{
		searchAPI: searchAPI,
		llm:       llm,
		memory:    make([]llms.ChatMessage, 0),
	}
}

func (crc *ConversationalRAGChain) Query(ctx context.Context, question string, topK int) (string, error) {
	searchReq := request.NewSearchListResources()
	searchReq.Q = question

	results, err := crc.searchAPI.SearchResources(ctx, searchReq)
	if err != nil {
		return "", err
	}

	prompt, err := crc.buildRAGPrompt(results, question)
	if err != nil {
		return "", err
	}

	response, err := llms.GenerateFromSinglePrompt(ctx, crc.llm, prompt)
	if err != nil {
		return "", err
	}

	// update the memory
	crc.memory = append(crc.memory,
		llms.HumanChatMessage{Content: question},
		llms.AIChatMessage{Content: response},
	)

	return response, nil
}

func (crc *ConversationalRAGChain) buildRAGPrompt(results interface{}, question string) (string, error) {
	const maxTokens = 12000

	jsonData, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	var searchResults *cdResults
	if err := json.Unmarshal(jsonData, &searchResults); err != nil {
		return "", fmt.Errorf("failed to unmarshal results: %w", err)
	}

	// Remove unnecessary fields
	for i := range searchResults.Hits {
		if searchResults.Hits[i].Value != nil {
			delete(searchResults.Hits[i].Value, "vectorsValue")
			delete(searchResults.Hits[i].Value, "catch_all")
			delete(searchResults.Hits[i].Value, "security")
			delete(searchResults.Hits[i].Value, "@id")
		}
	}

	contextJSON, err := json.MarshalIndent(searchResults, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	// If still too large, reduce number of hits
	maxChars := maxTokens * 4
	for len(contextJSON) > maxChars && len(searchResults.Hits) > 1 {
		searchResults.Hits = searchResults.Hits[:len(searchResults.Hits)-1]
		contextJSON, _ = json.MarshalIndent(searchResults, "", "  ")
	}

	history := crc.buildHistory()

	prompt := fmt.Sprintf(`You are a Corteza AI assistant with expertise in the Corteza Low-Code Platform.

## Your Role
- Help users understand and work with Corteza Compose records, modules, and namespaces
- Provide accurate information based on the available context
- Guide users on Corteza-specific concepts and terminology

## Conversation History
%s

## Retrieved Context from Corteza
%s

## Guidelines
1. **Answer based on evidence**: Use only information from the conversation history and retrieved context
2. **Be specific to Corteza**: Reference specific modules, fields, record IDs, and namespaces when relevant
3. **Admit limitations**: If information isn't in the context, clearly state: "I don't have that information in the current context"
4. **Cite your sources**: When referencing specific records or data, mention the module/namespace it came from
5. **Structured responses**: For lists of records, use clear formatting (tables or bullet points)
6. **Suggest actions**: When appropriate, suggest what Corteza operations might help (e.g., "You could filter by..." or "Check the X module for...")

## Current Question
%s

## Your Response
`, history, string(contextJSON), question)

	return prompt, nil
}

func (crc *ConversationalRAGChain) buildHistory() string {
	var history strings.Builder

	for _, msg := range crc.memory {
		history.WriteString(fmt.Sprintf("%s: %s\n", msg.GetType(), msg.GetContent()))
	}

	if history.Len() == 0 {
		return "No previous conversation."
	}

	return history.String()
}
