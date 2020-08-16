package bulk

import (
	"context"
)

type (
	storeInterface interface {
		storeGeneratedInterfaces
	}

	Job interface {
		Do(ctx context.Context, s storeInterface) error
	}
)
