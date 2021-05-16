package store

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

const (
	rbacSep = ":"
)

func newRbacRule(rl *rbac.Rule) *rbacRule {
	return &rbacRule{
		rule: rl,
	}
}

func (rl *rbacRule) MarshalEnvoy() ([]resource.Interface, error) {
	// @todo RBACv2
	//refRole := strconv.FormatUint(rl.rule.RoleID, 10)
	//
	//refRes, err := rbacResToRef(rl.rule.Resource.String())
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Remove the identifier once we're finished with it
	//rl.rule.Resource = rl.rule.Resource.TrimID()
	//
	//return envoy.CollectNodes(
	//	resource.NewRbacRule(rl.rule, refRole, refRes),
	//)
	return nil, nil
}

func rbacResToRef(rr string) (*resource.Ref, error) {
	if rr == "" {
		return nil, nil
	}

	ref := &resource.Ref{}

	rr = strings.TrimSpace(rr)
	rr = strings.TrimRight(rr, rbacSep)

	parts := strings.Split(rr, rbacSep)

	// When len is 1; only top-level defined (system, compose, ....)
	if len(parts) == 1 {
		ref.ResourceType = rr
		return ref, nil
	}

	//When len is 3; both levels defined; resource ref also provided
	if len(parts) == 3 {
		ref.ResourceType = strings.Join(parts[0:2], rbacSep) + rbacSep
		if parts[2] != "*" {
			ref.Identifiers = resource.MakeIdentifiers(parts[2])
		}
		return ref, nil
	}

	return nil, fmt.Errorf("invalid resource provided: %s", rr)
}
