package mapping

import (
	"context"

	"github.com/crusttech/human/server/system/service"
)

type (
	systemAccessControl interface {
		CanSearchUsers(ctx context.Context) bool
		CanSearchRoles(ctx context.Context) bool
		CanSearchApplications(ctx context.Context) bool
		CanSearchAuthClients(ctx context.Context) bool
	}

	systemMapping struct {
		ac systemAccessControl
	}
)

func SystemMapping() *systemMapping {
	return &systemMapping{
		ac: service.DefaultAccessControl,
	}
}

func (svc systemMapping) Users(_ context.Context) ([]*Mapping, error) {
	return []*Mapping{{
		Index: "system-users",
		Mapping: map[string]*property{
			"resourceType": {Type: "keyword"},

			"userID": {Type: "long"},

			"email":  {Type: "keyword", Boost: 2},
			"name":   {Type: "keyword", Boost: 10},
			"handle": {Type: "keyword", Boost: 2},

			"created": change(),
			"updated": change(),
			"deleted": change(),

			"suspendedAt": {Type: "date"},

			"security": security(),

			"catch_all": {
				Type:     "text",
				Analyzer: "standard",
			},
		},
	}}, nil
}
