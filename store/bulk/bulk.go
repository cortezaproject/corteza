package bulk

import (
	"context"
)

type (
	Job interface {
		Do(ctx context.Context, s storeInterface) error
	}
)
