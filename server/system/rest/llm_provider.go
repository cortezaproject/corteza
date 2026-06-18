package rest

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	LlmProvider struct {
		svc llmProviderService
		ac  llmProviderAccessController
	}

	llmProviderPayload struct {
		*types.LlmProvider

		CanGrant              bool `json:"canGrant"`
		CanUpdateLlmProvider  bool `json:"canUpdateLlmProvider"`
		CanDeleteLlmProvider  bool `json:"canDeleteLlmProvider"`
	}

	llmProviderSetPayload struct {
		Set []*llmProviderPayload `json:"set"`
	}

	llmProviderService interface {
		Create(ctx context.Context, p *types.LlmProvider, apiKey string) (*types.LlmProvider, error)
		FindByID(ctx context.Context, id uint64) (*types.LlmProvider, error)
		Update(ctx context.Context, p *types.LlmProvider, apiKey string) (*types.LlmProvider, error)
		Delete(ctx context.Context, id uint64, deletedBy uint64) error
		Search(ctx context.Context, f types.LlmProviderFilter) (types.LlmProviderSet, error)
		ListModels(ctx context.Context, providerID uint64) ([]string, error)
		Validate(ctx context.Context, providerID uint64) error
	}

	llmProviderAccessController interface {
		CanGrant(ctx context.Context) bool
		CanReadLlmProvider(ctx context.Context, p *types.LlmProvider) bool
		CanUpdateLlmProvider(ctx context.Context, p *types.LlmProvider) bool
		CanDeleteLlmProvider(ctx context.Context, p *types.LlmProvider) bool
	}
)

func (LlmProvider) New() LlmProvider {
	return LlmProvider{svc: service.DefaultLlmService, ac: service.DefaultAccessControl}
}

func (ctrl LlmProvider) List(ctx context.Context, r *request.LlmProviderList) (interface{}, error) {
	set, err := ctrl.svc.Search(ctx, types.LlmProviderFilter{
		Provider: r.Provider,
		Status:   r.Status,
	})
	return ctrl.makeFilterPayload(ctx, set, err)
}

func (ctrl LlmProvider) Create(ctx context.Context, r *request.LlmProviderCreate) (interface{}, error) {
	p := &types.LlmProvider{
		Handle:   r.Handle,
		Provider: r.Provider,
		Status:   r.Status,
		Meta:     r.Meta,
		Config:   r.Config,
	}
	res, err := ctrl.svc.Create(ctx, p, r.ApiKey)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl LlmProvider) Read(ctx context.Context, r *request.LlmProviderRead) (interface{}, error) {
	p, err := ctrl.svc.FindByID(ctx, r.LlmProviderID)
	if err != nil {
		return nil, err
	}
	if !ctrl.ac.CanReadLlmProvider(ctx, p) {
		return nil, fmt.Errorf("not allowed to read LLM provider")
	}
	return ctrl.makePayload(ctx, p, nil)
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
	res, err := ctrl.svc.Update(ctx, p, r.ApiKey)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl LlmProvider) Delete(ctx context.Context, r *request.LlmProviderDelete) (interface{}, error) {
	return nil, ctrl.svc.Delete(ctx, r.LlmProviderID, auth.GetIdentityFromContext(ctx).Identity())
}

func (ctrl LlmProvider) Models(ctx context.Context, r *request.LlmProviderModels) (interface{}, error) {
	return ctrl.svc.ListModels(ctx, r.LlmProviderID)
}

func (ctrl LlmProvider) Validate(ctx context.Context, r *request.LlmProviderValidate) (interface{}, error) {
	if err := ctrl.svc.Validate(ctx, r.LlmProviderID); err != nil {
		return nil, err
	}
	return true, nil
}

func (ctrl LlmProvider) makePayload(ctx context.Context, p *types.LlmProvider, err error) (*llmProviderPayload, error) {
	if err != nil || p == nil {
		return nil, err
	}

	return &llmProviderPayload{
		LlmProvider: p,

		CanGrant:             ctrl.ac.CanGrant(ctx),
		CanUpdateLlmProvider: ctrl.ac.CanUpdateLlmProvider(ctx, p),
		CanDeleteLlmProvider: ctrl.ac.CanDeleteLlmProvider(ctx, p),
	}, nil
}

func (ctrl LlmProvider) makeFilterPayload(ctx context.Context, set types.LlmProviderSet, err error) (*llmProviderSetPayload, error) {
	if err != nil {
		return nil, err
	}

	pp := &llmProviderSetPayload{Set: make([]*llmProviderPayload, len(set))}
	for i := range set {
		pp.Set[i], _ = ctrl.makePayload(ctx, set[i], nil)
	}

	return pp, nil
}
