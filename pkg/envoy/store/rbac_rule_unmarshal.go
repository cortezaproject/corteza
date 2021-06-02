package store

import (
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

func newRbacRule(rl *rbac.Rule) *rbacRule {
	return &rbacRule{
		rule: rl,
	}
}

func (rl *rbacRule) MarshalEnvoy() ([]resource.Interface, error) {
	refRole := strconv.FormatUint(rl.rule.RoleID, 10)

	rl.rule.Resource = rbac.ResourceType(rl.rule.Resource)

	return envoy.CollectNodes(
		resource.NewRbacRule(rl.rule, refRole, rl.refRes, rl.refPathRes...),
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

	// When len is 1; only top-level defined (system, compose, ....)
	//if len(parts) == 1 {
	//	ref.ResourceType = rr
	//	return ref, nil
	//}

	//When len is 3; both levels defined; resource ref also provided
	//if len(parts) == 3 {
	//	ref.ResourceType = strings.Join(parts[0:2], rbacSep) + rbacSep
	//	if parts[2] != "*" {
	//		ref.Identifiers = resource.MakeIdentifiers(parts[2])
	//	}
	//	return ref, nil
	//}

	return nil, fmt.Errorf("invalid resource provided: %s", rr)
}
