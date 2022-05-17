package dal

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"go.uber.org/zap"
)

type (
	PKValues map[string]any

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
		CreateRecords(ctx context.Context, m *Model, rr ...ValueGetter) error

		//UpdateRecords(ctx context.Context, m *data.Model, rr ...ValueGetter) error
		//DeleteRecordsByPK(ctx context.Context, m *data.Model, rr ...ValueGetter) error
		//TruncateRecords(ctx context.Context, m *data.Model) error

		LookupRecord(context.Context, *Model, ValueGetter, ValueSetter) error

		SearchRecords(context.Context, *Model, filter.Filter) (Iterator, error)

		// ---

		// DDL

		// Models returns all the models the underlying database already supports
		//
		// This is useful when adding support for new modules since we can find out what
		// can work out of the box.
		Models(context.Context) (ModelSet, error)

		// // returns all attribute types that driver supports
		// AttributeTypes() []data.AttributeType

		// AddModel adds support for the given models to the underlying database
		//
		// The operation returns an error if any of the models already exists.
		AddModel(context.Context, *Model, ...*Model) error

		// RemoveModel removes support for the given model from the underlying database
		RemoveModel(context.Context, *Model, ...*Model) error

		// AlterModel requests for metadata changes to the existing model
		//
		// Only metadata (such as idents) are affected; attributes can not be changed here
		AlterModel(ctx context.Context, old *Model, new *Model) error

		// AlterModelAttribute requests for the model attribute change
		//
		// Specific operations require data transformations (type change).
		// Some basic ops. should be implemented on DB driver level, but greater controll can be
		// achieved via the trans functions.
		AlterModelAttribute(ctx context.Context, sch *Model, old Attribute, new Attribute, trans ...TransformationFunction) error
	}

	TransformationFunction func(*Model, Attribute, expr.TypedValue) (expr.TypedValue, bool, error)

	// Iterator provides an interface for loading data from the underlying store
	Iterator interface {
		Next(ctx context.Context) bool
		Err() error
		Scan(ValueSetter) error
		Close() error

		BackCursor(ValueGetter) (*filter.PagingCursor, error)
		ForwardCursor(ValueGetter) (*filter.PagingCursor, error)

		// // -1 means unknown
		// Total() int
		// Cursor() any
		// // ... do we need anything else here?
	}

	// Store provides an interface which CRS uses to interact with the underlying database

	ValueGetter interface {
		CountValues() map[string]uint
		GetValue(string, uint) (any, error)
	}

	ValueSetter interface {
		SetValue(string, uint, any) error
	}

	ConnectorFn func(ctx context.Context, dsn string, cc ...capabilities.Capability) (StoreConnection, error)
)

var (
	registered = make(map[string]ConnectorFn)
)

func (pkv PKValues) CountValues() map[string]uint {
	c := make(map[string]uint)
	for k := range pkv {
		c[k] = 1
	}

	return c
}

func (pkv PKValues) GetValue(key string, _ uint) (any, error) {
	if val, has := pkv[key]; has {
		return val, nil
	} else {
		return nil, fmt.Errorf("missing")
	}
}

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
