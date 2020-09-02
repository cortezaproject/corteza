package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/rest/request"
)

type (
	namespacePayload struct {
		*types.Namespace

		CanGrant           bool `json:"canGrant"`
		CanUpdateNamespace bool `json:"canUpdateNamespace"`
		CanDeleteNamespace bool `json:"canDeleteNamespace"`
		CanManageNamespace bool `json:"canManageNamespace"`
		CanCreateModule    bool `json:"canCreateModule"`
		CanCreateChart     bool `json:"canCreateChart"`
		CanCreatePage      bool `json:"canCreatePage"`
	}

	payload struct{}
	Foobar  struct{}

	accessController interface {
		CanGrant(context.Context) bool
	}
)

func (Foobar) New() *Foobar {
	return &Foobar{}
}

func (ctrl Foobar) Foobar(ctx context.Context, r *request.FoobarFoobar) (interface{}, error) {
	return &struct{}{}, nil
}
