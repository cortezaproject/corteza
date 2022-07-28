package store

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

type (
	ConnectorFn func(ctx context.Context, dsn string) (s Storer, err error)
)

var (
	registered = make(map[string]ConnectorFn)
)

// Connect returns store based on dsn from environment.
//
// If you are in development environment,
// you can use {version} with database name, that will be replaced with build version
// suppose build version is `20xx.x.x-dev-1` then database_name_{version} will be `database_name_20xx_x_x_dev_1`,
//
// IE. `BUILD_VERSION=20xx.x.x-dev-1` and
// `DB_DSN=corteza:corteza@tcp(localhost:3306)/corteza_{version}?collation=utf8mb4_general_ci`
// will be `DB_DSN=corteza:corteza@tcp(localhost:3306)/corteza_20xx_x_x_dev_1?collation=utf8mb4_general_ci`
func Connect(ctx context.Context, log *zap.Logger, dsn string, isDevelopment bool) (s Storer, err error) {
	if isDevelopment {
		if strings.Contains(dsn, "{version}") {
			log.Warn("You're using DB_DSN with {version}, It is still in EXPERIMENTAL phase")
			log.Warn("Should be used only for development mode")
			log.Warn("You may experience instability")
		}
		expr := regexp.MustCompile(`[.\-]+`)
		version := expr.ReplaceAllString(os.Getenv("BUILD_VERSION"), "_")
		dsn = strings.Replace(dsn, "{version}", version, 1)
	}

	var storeType = strings.SplitN(dsn, "://", 2)[0]
	if storeType == "" {
		// Backward compatibility
		storeType = "mysql"
	}

	if conn, ok := registered[storeType]; ok {
		return conn(ctx, dsn)
	} else {
		return nil, fmt.Errorf("unknown store type used: %q (check your database configuration)", storeType)
	}
}

// Register add on ore more store types and their connector fn
func Register(fn ConnectorFn, tt ...string) {
	for _, t := range tt {
		registered[t] = fn
	}
}
