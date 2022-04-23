package crs

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/crs/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	driver interface {
		// Capabilities returns all of the capabilities the driver is able to support
		Capabilities() capabilities.Set
		// Can determines if the driver is able to handle the connection
		Can(dsn string, capabilities ...capabilities.Capability) bool

		// Store connects and returns a store we can use
		Store(ctx context.Context, dsn string) (Store, error)
		// Close closes the connection to specified
		Close(ctx context.Context, dsn string) error
	}

	// Loader provides an interface for loading data from the underlying store
	Loader interface {
		More() bool
		Load(*data.Model, []Setter) (coppied int, err error)

		// -1 means unknown
		Total() int
		Cursor() any
		// ... do we need anything else here?
	}

	// Store provides an interface which CRS uses to interact with the underlying database
	Store interface {
		// DML
		CreateRecords(ctx context.Context, sch *data.Model, cc ...Getter) error
		// @note providing a model here would allow us to control what attributes we return; would we find this useful?
		SearchRecords(ctx context.Context, sch *data.Model, filter any) (Loader, error)
		// Rest are same difference
		// - LookupRecord
		// - UpdateRecords
		// - DeleteRecords
		// - TruncateRecords

		// DDL
		// Models returns all of the collections the given database already defines
		Models(context.Context) (data.ModelSet, error)

		// returns all attribute types that driver supports
		AttributeTypes() []data.AttributeType

		// AddModel requests the driver to support the specified collections
		AddModel(context.Context, *data.Model, ...*data.Model) error

		// // RemoveModel requests the driver to remove support for the specified collections
		// RemoveModel(context.Context, *data.Model, ...*data.Model) error

		// AlterModel requests the driver to alter the general collectio parameters
		AlterModel(ctx context.Context, old *data.Model, new *data.Model) error

		// AlterModelAttribute requests the driver to alter the specified attribute of the given collection
		AlterModelAttribute(ctx context.Context, sch *data.Model, old data.Attribute, new data.Attribute, trans ...func(*data.Model, data.Attribute, expr.TypedValue) (expr.TypedValue, bool, error)) error
	}

	// probably somewhere else (on the db level?
	Getter interface {
		GetValue(name string, pos int) (expr.TypedValue, error)
	}

	Setter interface {
		SetValue(name string, pos int, value expr.TypedValue) error
	}
)
