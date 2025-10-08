package options

import (
	"github.com/cortezaproject/corteza/server/pkg/options"
)

type (
	VectorSearchOpt struct {
		EmbeddingsAPIUrl string `env:"EMBEDDINGS_API_URL"`
		EmbeddingsAPIKey string `env:"EMBEDDINGS_API_KEY"`
		EmbeddingsModel  string `env:"EMBEDDINGS_MODEL"`
	}
)

func VectorSearch() (o *VectorSearchOpt, err error) {
	o = &VectorSearchOpt{}

	return o, func() error {
		o.EmbeddingsAPIKey = options.EnvString("EMBEDDINGS_API_KEY", "")
		o.EmbeddingsAPIUrl = options.EnvString("EMBEDDINGS_API_URL", "")
		o.EmbeddingsModel = options.EnvString("EMBEDDINGS_MODEL", "paraphrase-multilingual-MiniLM-L12-v2")

		return nil
	}()
}
