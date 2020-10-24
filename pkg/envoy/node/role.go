package node

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	// role represents a Role
	Role struct {
		Res *types.Role
	}
)

func (n Role) GetResource() *types.Role {
	return n.Res
}

func (n Role) Identifiers() envoy.NodeIdentifiers {
	return makeIdentifiers(n.Res.ID, n.Res.Handle)
}

func (n Role) Matches(resource string, identifiers ...string) bool {
	return resource == n.Resource() && n.Identifiers().HasAny(identifiers...)
}

func (n Role) Resource() string {
	return "role"
}

func (n Role) Relations() envoy.NodeRelationships {
	return nil
}
