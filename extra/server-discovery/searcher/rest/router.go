package rest

import (
	"github.com/cortezaproject/corteza/extra/server-discovery/pkg/auth"
	"github.com/cortezaproject/corteza/extra/server-discovery/pkg/options"
	"github.com/cortezaproject/corteza/extra/server-discovery/pkg/vectorsearch"
	"github.com/cortezaproject/corteza/extra/server-discovery/searcher/rest/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/tmc/langchaingo/llms/openai"
	"go.uber.org/zap"
)

func MountRoutes(log *zap.Logger, opt options.VectorSearchOpt) func(r chi.Router) {
	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			embeddingSvc, err := vectorsearch.Embedder(log, opt)
			if err != nil {
				log.Error(err.Error())
			}

			search := Search(embeddingSvc)

			r.Use(auth.HttpTokenValidator("discovery"))

			llm, err := openai.New(
				openai.WithBaseURL(opt.LLMBaseURL),
				openai.WithToken(opt.LLMAPIKey),
				openai.WithModel(opt.LLMModel),
			)

			if err != nil {
				log.Error(err.Error())
			}

			handlers.NewSearch(search).MountRoutes(r)
			handlers.NewRagQuery(Rag(search, llm, log)).MountRoutes(r)
		})
	}
}
