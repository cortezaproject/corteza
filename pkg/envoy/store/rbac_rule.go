package store

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	rbacRule struct {
		cfg *EncoderConfig

		res  *resource.RbacRule
		rule *rbac.Rule

		relRole *types.Role
	}
)

func rbacResourceErrUnidentifiable(ii resource.Identifiers) error {
	return fmt.Errorf("rbac resource unidentifiable %v", ii.StringSlice())
}
