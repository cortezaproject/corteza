package node

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	// application represents a Application
	Application struct {
		Res *types.Application
	}
)

func (n Application) GetResource() *types.Application {
	return n.Res
}

func (n Application) Identifiers() envoy.NodeIdentifiers {
	return makeIdentifiers(n.Res.ID)
}

func (n Application) Matches(resource string, identifiers ...string) bool {
	return resource == n.Resource() && n.Identifiers().HasAny(identifiers...)
}

func (n Application) Resource() string {
	return types.ApplicationRBACResource.String()
}

func (n Application) Relations() envoy.NodeRelationships {
	return nil
}
