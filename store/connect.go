package store

import (
	"context"
	"fmt"
	"strings"
)

type (
	ConnectorFn func(ctx context.Context, dsn string) (s Storer, err error)
)

var (
	registered = make(map[string]ConnectorFn)
)

func Connect(ctx context.Context, dsn string) (s Storer, err error) {
	var storeType = strings.SplitN(dsn, "://", 2)[0]
	if storeType == "" {
		// Backward compatibility
		storeType = "mysql"
	}

	if conn, ok := registered[storeType]; ok {
		return conn(ctx, dsn)
	} else {
		return nil, fmt.Errorf("unknown store type used: %q (check your storage configuration)", storeType)
	}
}

// Register add on ore more store types and their connector fn
func Register(fn ConnectorFn, tt ...string) {
	for _, t := range tt {
		registered[t] = fn
	}
}
