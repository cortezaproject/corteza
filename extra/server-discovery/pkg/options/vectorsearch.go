package options

import (
	"github.com/cortezaproject/corteza/server/pkg/options"
)

const (
	TraditionalSearchMode string = "traditional"
	SemanticSearchMode    string = "semantic"
	HybridSearchMode      string = "hybrid"
)

type (
	VectorSearchOpt struct {
		EmbeddingsAPIUrl string `env:"EMBEDDINGS_API_URL"`
		EmbeddingsAPIKey string `env:"EMBEDDINGS_API_KEY"`
		EmbeddingsModel  string `env:"EMBEDDINGS_MODEL"`
		SearchMode       string `env:"SEARCH_MODE"`
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
		o.EmbeddingsModel = options.EnvString("EMBEDDINGS_MODEL", "paraphrase-multilingual-MiniLM-L12-v2")
		o.LLMBaseURL = options.EnvString("LLM_BASE_URL", "")
		o.LLMAPIKey = options.EnvString("LLM_API_KEY", "")
		o.LLMModel = options.EnvString("LLM_MODEL", "mistral-medium-latest")
		o.SearchMode = options.EnvString("SEARCH_MODE", TraditionalSearchMode)

		return nil
	}()
}
