package vectorsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cortezaproject/corteza/extra/server-discovery/pkg/options"
	"go.uber.org/zap"
)

type (
	EmbeddingRequest struct {
		Input []string `json:"input"`
		Model string   `json:"model,omitempty"`
	}

	EmbeddingData struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	}

	EmbeddingUsage struct {
		PromptAudioSeconds interface{} `json:"prompt_audio_seconds"`
		PromptTokens       int         `json:"prompt_tokens"`
		TotalTokens        int         `json:"total_tokens"`
		CompletionTokens   int         `json:"completion_tokens"`
		RequestCount       interface{} `json:"request_count"`
		PromptTokenDetails interface{} `json:"prompt_token_details"`
	}

	EmbeddingResponse struct {
		ID      string          `json:"id"`
		Object  string          `json:"object"`
		Data    []EmbeddingData `json:"data"`
		Usage   EmbeddingUsage  `json:"usage"`
		Message string          `json:"message"`
	}

	embedder struct {
		log             *zap.Logger
		vectorSearchOpt options.VectorSearchOpt
	}
)

func Embedder(log *zap.Logger, opt options.VectorSearchOpt) (out *embedder, err error) {
	out = &embedder{log: log, vectorSearchOpt: opt}
	return
}

// GenerateEmbeddings generates embeddings for the given input text
// by calling the external embedding API service.
func (embedder *embedder) GenerateEmbeddings(input string) ([]float64, error) {
	requestBody := EmbeddingRequest{
		Input: []string{input},
		Model: embedder.vectorSearchOpt.EmbeddingsModel,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	url := embedder.vectorSearchOpt.EmbeddingsAPIUrl
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", embedder.vectorSearchOpt.EmbeddingsAPIKey))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call embedding service: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("embedding service error %d: %s", resp.StatusCode, string(respBody))
	}

	var result EmbeddingResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode embedding response: %v", err)
	}

	if result.Object == "error" {
		return nil, fmt.Errorf("embedding service error: %s", result.Message)
	}

	return result.Data[0].Embedding, nil
}

func (embedder *embedder) ValidateEmbeddingDimensions(dimension int) (bool, error) {
	embeddings, err := embedder.GenerateEmbeddings("Corteza")
	if err != nil {
		return false, err
	}

	embeddingsLength := len(embeddings)

	if embeddingsLength == 0 {
		return false, fmt.Errorf("no embeddings returned")
	}

	if embeddingsLength != dimension {
		return false, fmt.Errorf("Embeddings dimensions mismatch: expected %d, got %d", dimension, embeddingsLength)
	}

	return true, nil
}
