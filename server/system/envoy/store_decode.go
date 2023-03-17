package envoy

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

func (d StoreDecoder) extendDecoder(ctx context.Context, s store.Storer, dl dal.FullService, rt string, nodes map[string]*envoyx.Node, f envoyx.ResourceFilter) (out envoyx.NodeSet, err error) {
	return
}

func (d StoreDecoder) extendedApigwRouteDecoder(ctx context.Context, s store.Storer, dl dal.FullService, f types.ApigwRouteFilter, base envoyx.NodeSet) (out envoyx.NodeSet, err error) {
	for _, b := range base {
		route := b.Resource.(*types.ApigwRoute)

		filters, err := d.decodeApigwFilter(ctx, s, dl, types.ApigwFilterFilter{
			RouteID: route.ID,
		})
		if err != nil {
			return nil, err
		}

		out = append(out, filters...)
	}

	return
}

func (d StoreDecoder) decodeAuthClientRefs(c *types.AuthClient) (refs map[string]envoyx.Ref) {
	refs = make(map[string]envoyx.Ref, 4)

	if c.Security.ImpersonateUser > 0 {
		refs["Security.ImpersonateUser"] = envoyx.Ref{
			ResourceType: types.UserResourceType,
			Identifiers:  envoyx.MakeIdentifiers(c.Security.ImpersonateUser),
		}
	}

	d.roleSliceToRefs(refs, "Security.PermittedRoles", c.Security.PermittedRoles)
	d.roleSliceToRefs(refs, "Security.ProhibitedRoles", c.Security.ProhibitedRoles)
	d.roleSliceToRefs(refs, "Security.ForcedRoles", c.Security.ForcedRoles)

	return
}

func (d StoreDecoder) decodeDalConnectionRefs(c *types.DalConnection) (refs map[string]envoyx.Ref) {
	if c.Config.Privacy.SensitivityLevelID == 0 {
		return
	}

	refs = map[string]envoyx.Ref{
		"Config.Privacy.SensitivityLevelID": {
			ResourceType: types.DalSensitivityLevelResourceType,
			Identifiers:  envoyx.MakeIdentifiers(c.Config.Privacy.SensitivityLevelID),
		},
	}

	return
}

func (d StoreDecoder) roleSliceToRefs(refs map[string]envoyx.Ref, k string, rr []string) {
	for i, r := range rr {
		refs[fmt.Sprintf("%s.%d.RoleID", k, i)] = envoyx.Ref{
			ResourceType: types.RoleResourceType,
			Identifiers:  envoyx.MakeIdentifiers(r),
		}
	}
}
