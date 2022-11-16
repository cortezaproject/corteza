package store

import (
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	rbacRule struct {
		cfg *EncoderConfig

		rule *rbac.Rule

		// point to the rbac rule
		refRbacResource string
		refRbacRes      *resource.Ref

		refPathRes []*resource.Ref

		refRole *resource.Ref
		role    *types.Role
	}
)

func rbacResourceErrUnidentifiable(ii resource.Identifiers) error {
	return fmt.Errorf("rbac resource unidentifiable %v", ii.StringSlice())
}
