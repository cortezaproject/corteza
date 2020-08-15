package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/compose_namespaces.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	composeNamespaceCreate struct {
		Done chan struct{}
		res  *types.Namespace
		err  error
	}

	composeNamespaceUpdate struct {
		Done chan struct{}
		res  *types.Namespace
		err  error
	}

	composeNamespaceRemove struct {
		Done chan struct{}
		res  *types.Namespace
		err  error
	}
)

// CreateComposeNamespace creates a new ComposeNamespace
// create job that can be pushed to store's transaction handler
func CreateComposeNamespace(res *types.Namespace) *composeNamespaceCreate {
	return &composeNamespaceCreate{res: res}
}

// Do Executes composeNamespaceCreate job
func (j *composeNamespaceCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateComposeNamespace(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateComposeNamespace creates a new ComposeNamespace
// update job that can be pushed to store's transaction handler
func UpdateComposeNamespace(res *types.Namespace) *composeNamespaceUpdate {
	return &composeNamespaceUpdate{res: res}
}

// Do Executes composeNamespaceUpdate job
func (j *composeNamespaceUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateComposeNamespace(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveComposeNamespace creates a new ComposeNamespace
// remove job that can be pushed to store's transaction handler
func RemoveComposeNamespace(res *types.Namespace) *composeNamespaceRemove {
	return &composeNamespaceRemove{res: res}
}

// Do Executes composeNamespaceRemove job
func (j *composeNamespaceRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveComposeNamespace(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}
