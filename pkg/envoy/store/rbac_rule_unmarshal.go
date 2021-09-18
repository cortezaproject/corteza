package store

import (
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/system/types"
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
		resource.NewRbacRule(rl.rule, rl.refRole.Identifiers.First(), rl.refRbacRes, rl.refRbacResource, rl.refPathRes...),
	)
}

func rbacResToRef(rr string) (*resource.Ref, error) {
	if rr == "" {
		return nil, nil
	}

	ref := &resource.Ref{}
	if ref.ResourceType = rbac.ResourceType(rr); ref.ResourceType != "" {
		return ref, nil
	}

	return nil, fmt.Errorf("invalid resource provided: %s", rr)
}
