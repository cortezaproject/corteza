package envoy

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"github.com/cortezaproject/corteza/server/system/types"
	"gopkg.in/yaml.v3"
)

func (e YamlEncoder) encodeAuthClientSecurityC(ctx context.Context, p envoyx.EncodeParams, tt envoyx.Traverser, n *envoyx.Node, ac *types.AuthClient, sec *types.AuthClientSecurity) (_ any, err error) {
	sqPermittedRoles, err := e.encodeRoleSlice(n, tt, "Security.PermittedRoles", sec.PermittedRoles)
	if err != nil {
		return
	}

	sqProhibitedRoles, err := e.encodeRoleSlice(n, tt, "Security.ProhibitedRoles", sec.ProhibitedRoles)
	if err != nil {
		return
	}

	sqForcedRoles, err := e.encodeRoleSlice(n, tt, "Security.ForcedRoles", sec.ForcedRoles)
	if err != nil {
		return
	}

	var impersonateUser string
	if _, ok := n.References["Security.ImpersonateUser.UserID"]; ok {
		node := tt.ParentForRef(n, n.References["Security.ImpersonateUser.UserID"])
		if node == nil {
			err = fmt.Errorf("node not found @todo error")
			return
		}
		impersonateUser = n.Identifiers.FriendlyIdentifier()
	}

	return y7s.MakeMap(
		"impersonateUser", impersonateUser,
		"permittedRoles", sqPermittedRoles,
		"prohibitedRoles", sqProhibitedRoles,
		"forcedRoles", sqForcedRoles,
	)
}

func (e YamlEncoder) encodeRoleSlice(n *envoyx.Node, tt envoyx.Traverser, k string, rr []string) (out *yaml.Node, err error) {
	sq, _ := y7s.MakeSeq()

	for i := range rr {
		node := tt.ParentForRef(n, n.References[fmt.Sprintf("%s.%d.RoleID", k, i)])
		if node == nil {
			err = fmt.Errorf("node not found @todo error")
			return
		}

		sq, err = y7s.AddSeq(sq, node.Identifiers.FriendlyIdentifier())
		if err != nil {
			return
		}
	}

	return sq, nil
}
