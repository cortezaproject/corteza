package llm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cortezaproject/corteza/server/pkg/api"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/go-chi/chi/v5"
)

func MountRoutes(svc *Service) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", searchHandler(svc))
		r.Post("/", createHandler(svc))
		r.Get("/{providerID}", lookupHandler(svc))
		r.Put("/{providerID}", updateHandler(svc))
		r.Delete("/{providerID}", deleteHandler(svc))
		r.Post("/{providerID}/prompt", promptHandler(svc))
	}
}

func createHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			sysTypes.LlmProvider
			APIKey string `json:"apiKey,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			api.Send(w, r, err)
			return
		}

		result, err := svc.Create(r.Context(), &payload.LlmProvider, payload.APIKey)
		api.Send(w, r, err, result)
	}
}

func lookupHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerID, err := paramID(r, "providerID")
		if err != nil {
			api.Send(w, r, err)
			return
		}

		result, err := svc.LookupByID(r.Context(), providerID)
		api.Send(w, r, err, result)
	}
}

func updateHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerID, err := paramID(r, "providerID")
		if err != nil {
			api.Send(w, r, err)
			return
		}

		var p sysTypes.LlmProvider
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			api.Send(w, r, err)
			return
		}

		p.ID = providerID
		result, err := svc.Update(r.Context(), &p)
		api.Send(w, r, err, result)
	}
}

func deleteHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerID, err := paramID(r, "providerID")
		if err != nil {
			api.Send(w, r, err)
			return
		}

		api.Send(w, r, svc.Delete(r.Context(), providerID, 0))
	}
}

func searchHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := sysTypes.LlmProviderFilter{
			Provider: r.URL.Query().Get("provider"),
			Status:   r.URL.Query().Get("status"),
		}

		result, err := svc.Search(r.Context(), f)
		api.Send(w, r, err, result)
	}
}

func promptHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerID, err := paramID(r, "providerID")
		if err != nil {
			api.Send(w, r, err)
			return
		}

		var payload struct {
			Messages []Message `json:"messages"`
			Tools    []Tool    `json:"tools,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			api.Send(w, r, err)
			return
		}

		result, err := svc.Prompt(r.Context(), providerID, payload.Messages, payload.Tools)
		api.Send(w, r, err, result)
	}
}

func paramID(r *http.Request, key string) (uint64, error) {
	v := chi.URLParam(r, key)
	id, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return id, nil
}
