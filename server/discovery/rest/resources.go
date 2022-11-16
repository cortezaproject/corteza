package rest

import (
	"context"
	"github.com/cortezaproject/corteza/server/discovery/rest/internal/documents"
	"github.com/cortezaproject/corteza/server/discovery/rest/request"
)

type (
	resources struct {
		sys interface {
			Users(ctx context.Context, limit uint, cur string, userID uint64, deleted uint) (*documents.Response, error)
		}

		cmp interface {
			Namespaces(ctx context.Context, limit uint, cur string, namespaceID uint64, deleted uint) (*documents.Response, error)
			Modules(ctx context.Context, namespaceID uint64, limit uint, cur string, moduleID uint64, deleted uint) (*documents.Response, error)
			Records(ctx context.Context, namespaceID, moduleID uint64, limit uint, cur string, recordID uint64, deleted uint) (*documents.Response, error)
		}
	}
)

func Resources() *resources {
	return &resources{
		sys: documents.SystemResources(),
		cmp: documents.ComposeResources(),
	}
}

func (ctrl resources) SystemUsers(ctx context.Context, r *request.ResourcesSystemUsers) (interface{}, error) {
	return ctrl.sys.Users(ctx, r.Limit, r.PageCursor, r.UserID, r.Deleted)
}

func (ctrl resources) ComposeNamespaces(ctx context.Context, r *request.ResourcesComposeNamespaces) (interface{}, error) {
	return ctrl.cmp.Namespaces(ctx, r.Limit, r.PageCursor, r.NamespaceID, r.Deleted)
}

func (ctrl resources) ComposeModules(ctx context.Context, r *request.ResourcesComposeModules) (interface{}, error) {
	return ctrl.cmp.Modules(ctx, r.NamespaceID, r.Limit, r.PageCursor, r.ModuleID, r.Deleted)
}

func (ctrl resources) ComposeRecords(ctx context.Context, r *request.ResourcesComposeRecords) (interface{}, error) {
	return ctrl.cmp.Records(ctx, r.NamespaceID, r.ModuleID, r.Limit, r.PageCursor, r.RecordID, r.Deleted)
}
