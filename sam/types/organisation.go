package types

import (
	"github.com/crusttech/crust/system/types"
)

type (
	// Organisations - Organisations represent a top-level grouping entity.
	// There may be many organisations defined in a single deployment.
	Organisation struct {
		types.Organisation
	}
)
