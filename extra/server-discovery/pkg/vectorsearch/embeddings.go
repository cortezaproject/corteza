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
		Input []string `json:"input"` // TODO:: replace this with "input"
		Model string   `json:"model,omitempty"`
	}

	EmbeddingResponse struct {
		Embeddings [][]float32 `json:"embeddings"`
		Dimension  int         `json:"dimension"`
		Count      int         `json:"count"`
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

// consolidate this it should be on both
func (embedder *embedder) GenerateEmbeddings(input string) ([]float32, error) {
	requestBody := EmbeddingRequest{
		Input: []string{input},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// TODO: replace this with the env configurable option
	// Add the API KEY section
	url := embedder.vectorSearchOpt.EmbeddingsAPIUrl
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

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

	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return result.Embeddings[0], nil
}

func (embedder *embedder) ValidateEmbeddingsAPI(dimension int) (bool, error) {
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
