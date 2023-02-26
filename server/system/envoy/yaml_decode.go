package envoy

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"github.com/cortezaproject/corteza/server/system/types"
	"gopkg.in/yaml.v3"
)

func (d *auxYamlDoc) unmarshalYAML(k string, n *yaml.Node) (out envoyx.NodeSet, err error) {
	return
}

func (d *auxYamlDoc) unmarshalFiltersExtendedNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	return d.unmarshalApigwFilterNode(dctx, n, meta...)
}

func (d *auxYamlDoc) unmarshalUserRolesNode(r *types.User, n *yaml.Node) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
	refs = make(map[string]envoyx.Ref, len(n.Content))

	i := 0
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		refs[fmt.Sprintf("Roles.%d", i)] = envoyx.Ref{
			ResourceType: types.RoleResourceType,
			Identifiers:  envoyx.MakeIdentifiers(n.Value),
		}

		return nil
	})

	return
}

func (d *auxYamlDoc) unmarshalAuthClientSecurityNode(r *types.AuthClient, n *yaml.Node) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
	refs = make(map[string]envoyx.Ref)

	err = y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "impersonateuser":
			var av string
			err = y7s.DecodeScalar(v, "Impersonate user", &av)

			refs["Security.ImpersonateUser"] = envoyx.Ref{
				ResourceType: types.UserResourceType,
				Identifiers:  envoyx.MakeIdentifiers(av),
			}
			break

		case "permittedroles":
			refs = envoyx.MergeRefs(refs, roleSliceToRefs("Security.PermittedRoles", v))
			break

		case "prohibitedroles":
			refs = envoyx.MergeRefs(refs, roleSliceToRefs("Security.ProhibitedRoles", v))
			break

		case "forcedroles":
			refs = envoyx.MergeRefs(refs, roleSliceToRefs("Security.ForcedRoles", v))
			break

		}

		return nil
	})

	return
}

func roleSliceToRefs(k string, n *yaml.Node) (refs map[string]envoyx.Ref) {
	refs = make(map[string]envoyx.Ref, len(n.Content))

	i := 0
	y7s.EachSeq(n, func(n *yaml.Node) error {
		refs[fmt.Sprintf("%s.%d", k, i)] = envoyx.Ref{
			ResourceType: types.RoleResourceType,
			Identifiers:  envoyx.MakeIdentifiers(n.Value),
		}
		i++
		return nil
	})

	return
}
