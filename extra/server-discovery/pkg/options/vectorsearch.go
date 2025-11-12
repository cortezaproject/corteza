package options

import (
	"github.com/cortezaproject/corteza/server/pkg/options"
)

type (
	VectorSearchOpt struct {
		EmbeddingsAPIUrl string `env:"EMBEDDINGS_API_URL"`
		EmbeddingsAPIKey string `env:"EMBEDDINGS_API_KEY"`
		EmbeddingsModel  string `env:"EMBEDDINGS_MODEL"`
		LLMBaseURL       string `env:"LLM_BASE_URL"`
		LLMAPIKey        string `env:"LLM_API_KEY"`
		LLMModel         string `env:"LLM_MODEL"`
	}
)

func VectorSearch() (o *VectorSearchOpt, err error) {
	o = &VectorSearchOpt{}

	return o, func() error {
		o.EmbeddingsAPIKey = options.EnvString("EMBEDDINGS_API_KEY", "")
		o.EmbeddingsAPIUrl = options.EnvString("EMBEDDINGS_API_URL", "")
		o.EmbeddingsModel = options.EnvString("EMBEDDINGS_MODEL", "mistral-embed")
		o.LLMBaseURL = options.EnvString("LLM_BASE_URL", "")
		o.LLMAPIKey = options.EnvString("LLM_API_KEY", "")
		o.LLMModel = options.EnvString("LLM_MODEL", "mistral-medium-latest")

		return nil
	}()
}
