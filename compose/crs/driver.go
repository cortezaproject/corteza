package crs

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/crs/capabilities"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"go.uber.org/zap"
)

type (
	StoreConnection interface {
		// ---

		// Meta
		// Capabilities returns all of the capabilities the given store supports
		Capabilities() capabilities.Set
		// Can returns true if this store can handle the given capabilities
		Can(capabilities ...capabilities.Capability) bool

		// ---

		// Connection stuff
		// Close closes the store connection allowing the driver to perform potential
		// cleanup operations
		Close(ctx context.Context) error

		// ---

		// DML
		// CreateRecords stores the given records into the underlying database
		//
		// @todo rename to Create and replace types.Record with Getter
		CreateRecords(ctx context.Context, m *data.Model, rr ...*types.Record) error
		// SearchRecords returns an iterator we can use to load data
		//
		// @todo rename to Search and replace filter with xyz
		SearchRecords(ctx context.Context, sch *data.Model, filter any) (Iterator, error)
		// ...
		// @todo other functions; they are same difference so omitting for now
		// - LookupRecord
		// - UpdateRecords
		// - DeleteRecords
		// - TruncateRecords

		// ---

		// DDL
		// Models returns all of the models the underlying database already supports
		//
		// This is useful when adding support for new modules since we can find out what
		// can work out of the box.
		Models(context.Context) (data.ModelSet, error)

		// // returns all attribute types that driver supports
		// AttributeTypes() []data.AttributeType

		// AddModel adds support for the given models to the underlying database
		//
		// The operation returns an error if any of the models already exists.
		AddModel(context.Context, *data.Model, ...*data.Model) error

		// RemoveModel removes support for the given model from the underlying database
		RemoveModel(context.Context, *data.Model, ...*data.Model) error

		// AlterModel requests for metadata changes to the existing model
		//
		// Only metadata (such as idents) are affected; attributes can not be changed here
		AlterModel(ctx context.Context, old *data.Model, new *data.Model) error

		// AlterModelAttribute requests for the model attribute change
		//
		// Specific operations require data transformations (type change).
		// Some basic ops. should be implemented on DB driver level, but greater controll can be
		// achieved via the trans functions.
		AlterModelAttribute(ctx context.Context, sch *data.Model, old data.Attribute, new data.Attribute, trans ...TransformationFunction) error
	}

	TransformationFunction func(*data.Model, data.Attribute, expr.TypedValue) (expr.TypedValue, bool, error)

	// Iterator provides an interface for loading data from the underlying store
	Iterator interface {
		Next(ctx context.Context) bool
		Scan(r *types.Record) (err error)
		Err() error
		Close() error

		// // -1 means unknown
		// Total() int
		// Cursor() any
		// // ... do we need anything else here?
	}

	// Store provides an interface which CRS uses to interact with the underlying database

	Getter interface {
		GetValue(name string, pos int) (any, error)
	}

	Setter interface {
		SetValue(name string, pos int, value any) error
	}

	ConnectorFn func(ctx context.Context, dsn string, cc ...capabilities.Capability) (StoreConnection, error)
)

var (
	registered = make(map[string]ConnectorFn)
)

// Register registers a new connector for the given DSN schema
//
// In case of a duplicate schema the latter overwrites the prior
func Register(fn ConnectorFn, tt ...string) {
	for _, t := range tt {
		registered[t] = fn
	}
}

// connect opens a new StoreConnection for the given CRS
func connect(ctx context.Context, log *zap.Logger, def crsDefiner, isDevelopment bool) (StoreConnection, error) {
	dsn := def.StoreDSN()

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
		return conn(ctx, dsn, def.Capabilities()...)
	} else {
		return nil, fmt.Errorf("unknown store type used: %q (check your storage configuration)", storeType)
	}
}
