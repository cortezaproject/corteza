// This package is just so that I have an interface conforming driver

package noop

import (
	"context"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/crs"
	"github.com/cortezaproject/corteza-server/compose/crs/capabilities"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	noopDriver struct{}

	noopStore map[string]*types.Record
)

var (
	gStore = make(noopStore)
)

func Driver() *noopDriver {
	return &noopDriver{}
}

func (d noopDriver) Capabilities() capabilities.Set {
	return capabilities.FullCapabilities()
}

func (d noopDriver) Can(uri string, cc ...capabilities.Capability) bool {
	return strings.HasPrefix(uri, "noop://") && capabilities.Set(cc).IsSubset(d.Capabilities()...)
}

func (d noopDriver) Store(ctx context.Context, uri string) (str crs.Store, err error) {
	return &noopStore{}, nil
}

// ---

func (ds *noopStore) CreateRecords(ctx context.Context, sch *data.Model, cc ...crs.Getter) error {
	return nil
}

func (ds *noopStore) SearchRecords(ctx context.Context, sch *data.Model, filter any) (crs.Loader, error) {
	return nil, nil
}

func (ds *noopStore) Models(context.Context) (data.ModelSet, error) {
	return nil, nil
}

func (ds *noopStore) AddModel(context.Context, *data.Model, ...*data.Model) error {
	return nil
}

func (ds *noopStore) AlterModel(ctx context.Context, old *data.Model, new *data.Model) error {
	return nil
}

func (ds *noopStore) AlterModelAttribute(ctx context.Context, sch *data.Model, old data.Attribute, new data.Attribute, trans ...func(*data.Model, data.Attribute, expr.TypedValue) (expr.TypedValue, bool, error)) error {
	return nil
}
