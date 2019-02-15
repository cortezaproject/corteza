package organization

import (
	"context"

	"github.com/crusttech/crust/system/types"
)

func Crust() types.Organisation {
	return types.Organisation{ID: 1}
}

func GetFromContext(ctx context.Context) types.Organisation {
	if orgID, ok := ctx.Value("organizationID").(uint64); ok {
		return types.Organisation{ID: orgID}
	} else {
		return Crust()
	}
}
