package node

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	// user represents a User
	User struct {
		Res *types.User
	}
)

func (n User) GetResource() *types.User {
	return n.Res
}

func (n User) Identifiers() envoy.NodeIdentifiers {
	return makeIdentifiers(n.Res.ID, n.Res.Handle, n.Res.Email)
}

func (n User) Matches(resource string, identifiers ...string) bool {
	return resource == n.Resource() && n.Identifiers().HasAny(identifiers...)
}

func (n User) Resource() string {
	return "user"
}

func (n User) Relations() envoy.NodeRelationships {
	return nil
}
