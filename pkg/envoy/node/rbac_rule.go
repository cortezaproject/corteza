package node

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

type (
	// user represents a RbacRule
	RbacRule struct {
		Res     *rbac.Rule
		RefRole string
	}
)

func (n RbacRule) Resource() string { return "rbacRule" }
func (n RbacRule) GetResource() *rbac.Rule {
	return n.Res
}

func (n RbacRule) Identifiers() envoy.NodeIdentifiers { return nil }

func (n RbacRule) Matches(resource string, identifiers ...string) bool {
	return resource == n.Resource() && n.Identifiers().HasAny(identifiers...)
}

func (n RbacRule) Relations() envoy.NodeRelationships {
	ref := make(envoy.NodeRelationships)
	envoy.NodeRelationships{}.Add(Role{}.Resource(), n.RefRole)
	return ref
}

func (n *RbacRule) Update(rr ...envoy.Node) {
	for _, r := range rr {
		switch r := r.(type) {
		case *Role:
			n.updateRole(r)
		}
	}
}

func (n *RbacRule) updateRole(r *Role) {
	if r.Identifiers().HasAny(n.RefRole) {
		n.Res.RoleID = r.Res.ID
		n.RefRole = ""
	}
}
