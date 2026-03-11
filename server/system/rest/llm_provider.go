package rest

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	LlmProvider struct {
		svc llmProviderService
		ac  llmProviderAccessController
	}

	llmProviderService interface {
		Create(ctx context.Context, p *types.LlmProvider, apiKey string) (*types.LlmProvider, error)
		LookupByID(ctx context.Context, id uint64) (*types.LlmProvider, error)
		Update(ctx context.Context, p *types.LlmProvider) (*types.LlmProvider, error)
		Delete(ctx context.Context, id uint64, deletedBy uint64) error
		Search(ctx context.Context, f types.LlmProviderFilter) (types.LlmProviderSet, error)
		ListModels(ctx context.Context, providerID uint64) ([]string, error)
	}

	llmProviderAccessController interface {
		CanReadLlmProvider(ctx context.Context, p *types.LlmProvider) bool
	}
)

func (LlmProvider) New() LlmProvider {
	return LlmProvider{svc: service.DefaultLlmService, ac: service.DefaultAccessControl}
}

func (ctrl LlmProvider) List(ctx context.Context, r *request.LlmProviderList) (interface{}, error) {
	return ctrl.svc.Search(ctx, types.LlmProviderFilter{
		Provider: r.Provider,
		Status:   r.Status,
	})
}

func (ctrl LlmProvider) Create(ctx context.Context, r *request.LlmProviderCreate) (interface{}, error) {
	p := &types.LlmProvider{
		Handle:   r.Handle,
		Provider: r.Provider,
		Status:   r.Status,
		Meta:     r.Meta,
		Config:   r.Config,
	}
	return ctrl.svc.Create(ctx, p, r.ApiKey)
}

func (ctrl LlmProvider) Read(ctx context.Context, r *request.LlmProviderRead) (interface{}, error) {
	p, err := ctrl.svc.LookupByID(ctx, r.LlmProviderID)
	if err != nil {
		return nil, err
	}
	if !ctrl.ac.CanReadLlmProvider(ctx, p) {
		return nil, fmt.Errorf("not allowed to read LLM provider")
	}
	return p, nil
}

func (ctrl LlmProvider) Update(ctx context.Context, r *request.LlmProviderUpdate) (interface{}, error) {
	p := &types.LlmProvider{
		ID:       r.LlmProviderID,
		Handle:   r.Handle,
		Provider: r.Provider,
		Status:   r.Status,
		Meta:     r.Meta,
		Config:   r.Config,
	}
	return ctrl.svc.Update(ctx, p)
}

func (ctrl LlmProvider) Delete(ctx context.Context, r *request.LlmProviderDelete) (interface{}, error) {
	return nil, ctrl.svc.Delete(ctx, r.LlmProviderID, auth.GetIdentityFromContext(ctx).Identity())
}

func (ctrl LlmProvider) Models(ctx context.Context, r *request.LlmProviderModels) (interface{}, error) {
	return ctrl.svc.ListModels(ctx, r.LlmProviderID)
}
