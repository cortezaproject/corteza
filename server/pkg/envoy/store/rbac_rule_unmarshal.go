package store

import (
	"strconv"

	"github.com/crusttech/human/server/pkg/envoy"
	"github.com/crusttech/human/server/pkg/envoy/resource"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/system/types"
)

func newRbacRule(rl *rbac.Rule) (*rbacRule, error) {
	res := rl.Resource
	_, ref, pp, err := resource.ParseRule(res)

	return &rbacRule{
		rule: rl,

		refRbacResource: res,
		refRbacRes:      ref,

		refPathRes: pp,

		refRole: resource.MakeRef(types.RoleResourceType, resource.MakeIdentifiers(strconv.FormatUint(rl.RoleID, 10))),
	}, err
}

func (rl *rbacRule) MarshalEnvoy() ([]resource.Interface, error) {
	return envoy.CollectNodes(
		resource.NewRbacRule(rl.rule, rl.refRole, rl.refRbacRes, rl.refRbacResource, rl.refPathRes...),
	)
}
