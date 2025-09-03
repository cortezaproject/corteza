package rbac

import (
	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/spf13/cast"
)

func AllFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("isDescendantOf", isDescendantOf),
	}
}

func isDescendantOf(userID any, resourceOwner any) bool {
	if gRBAC == nil {
		return false
	}

	owner := id.MustNumID(cast.ToUint64(resourceOwner))
	if owner.IsZero() {
		return false
	}

	user := id.MustNumID(cast.ToUint64(cast.ToString(userID)))
	if user.IsZero() {
		return false
	}

	out := gRBAC.orgTree.IsAbove(owner, user)
	return out
}
