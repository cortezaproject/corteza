package dml

import (
	"context"
	"fmt"

	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/ast"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/store"
)

// internalConn implements dal.Connection backed by compose records in the
// internal database. It is ephemeral: created at the start of a publish run
// and removed when the run completes. Write operations are not supported.
type internalConn struct {
	namespaceID uint64
	store       store.Storer
	recordSvc   ComposeRecordSvc

	// modules indexed by handle, populated lazily on first Models() call
	modulesByHandle map[string]*composeTypes.Module
}

// NewInternalConn returns a dal.Connection reading from the given namespace.
func NewInternalConn(namespaceID uint64, s store.Storer, r ComposeRecordSvc) *internalConn {
	return &internalConn{
		namespaceID: namespaceID,
		store:       s,
		recordSvc:   r,
	}
}

// Models loads compose modules for the namespace and returns them as dal.ModelSet.
// The model Ident is the module handle; ResourceID carries the module ID for Search.
func (c *internalConn) Models(ctx context.Context) (dal.ModelSet, error) {
	mm, _, err := store.SearchComposeModules(ctx, c.store, composeTypes.ModuleFilter{
		NamespaceID: c.namespaceID,
	})
	if err != nil {
		return nil, fmt.Errorf("internal dal: load modules for ns %d: %w", c.namespaceID, err)
	}

	c.modulesByHandle = make(map[string]*composeTypes.Module, len(mm))
	out := make(dal.ModelSet, 0, len(mm))
	for _, m := range mm {
		c.modulesByHandle[m.Handle] = m
		model := &dal.Model{
			Ident:      m.Handle,
			Label:      m.Name,
			ResourceID: m.ID,
		}
		out = append(out, model)
	}
	return out, nil
}

// Operations returns the read-only operations this connection supports.
func (c *internalConn) Operations() dal.OperationSet {
	return dal.SearchOperations()
}

// Can returns true only for read operations.
func (c *internalConn) Can(ops ...dal.Operation) bool {
	allowed := dal.SearchOperations()
	for _, op := range ops {
		if !allowed.IsSuperset(op) {
			return false
		}
	}
	return true
}

// RegisterModelCache satisfies the externalModelRegistrar interface required by
// dal.SearchExternalData. For internal connections the model is already cached
// via Models(); this is a no-op.
func (c *internalConn) RegisterModelCache(_ context.Context, _ ...*dal.Model) error {
	return nil
}

// Search returns an iterator over compose records for the module identified by
// model.Ident (the module handle).
func (c *internalConn) Search(ctx context.Context, m *dal.Model, _ filter.Filter) (dal.Iterator, error) {
	if c.modulesByHandle == nil {
		if _, err := c.Models(ctx); err != nil {
			return nil, err
		}
	}

	mod, ok := c.modulesByHandle[m.Ident]
	if !ok {
		return nil, fmt.Errorf("internal dal: module %q not found in namespace %d", m.Ident, c.namespaceID)
	}

	records, _, err := c.recordSvc.Search(ctx, composeTypes.RecordFilter{
		ModuleID:    mod.ID,
		NamespaceID: c.namespaceID,
	})
	if err != nil {
		return nil, fmt.Errorf("internal dal: search records for module %q: %w", m.Ident, err)
	}

	return &composeRecordIterator{records: records, pos: -1}, nil
}

// --- unsupported write methods ---

func (c *internalConn) Create(_ context.Context, _ *dal.Model, _ ...dal.ValueGetter) ([]map[string]any, error) {
	return nil, fmt.Errorf("internal dal: create not supported")
}
func (c *internalConn) Update(_ context.Context, _ *dal.Model, _ dal.ValueGetter) error {
	return fmt.Errorf("internal dal: update not supported")
}
func (c *internalConn) Lookup(_ context.Context, _ *dal.Model, _ dal.ValueGetter, _ dal.ValueSetter) error {
	return fmt.Errorf("internal dal: lookup not supported")
}
func (c *internalConn) Count(_ context.Context, _ *dal.Model, _ filter.Filter) (uint, error) {
	return 0, fmt.Errorf("internal dal: count not supported")
}
func (c *internalConn) Analyze(_ context.Context, _ *dal.Model) (map[string]dal.OpAnalysis, error) {
	return nil, fmt.Errorf("internal dal: analyze not supported")
}
func (c *internalConn) Aggregate(_ context.Context, _ *dal.Model, _ filter.Filter, _ []dal.AggregateAttr, _ []dal.AggregateAttr, _ *ast.ASTNode) (dal.Iterator, error) {
	return nil, fmt.Errorf("internal dal: aggregate not supported")
}
func (c *internalConn) Delete(_ context.Context, _ *dal.Model, _ dal.ValueGetter) error {
	return fmt.Errorf("internal dal: delete not supported")
}
func (c *internalConn) Truncate(_ context.Context, _ *dal.Model) error {
	return fmt.Errorf("internal dal: truncate not supported")
}
func (c *internalConn) CreateModel(_ context.Context, _ ...*dal.Model) error {
	return fmt.Errorf("internal dal: createModel not supported")
}
func (c *internalConn) DeleteModel(_ context.Context, _ ...*dal.Model) error {
	return fmt.Errorf("internal dal: deleteModel not supported")
}
func (c *internalConn) UpdateModel(_ context.Context, _ *dal.Model, _ *dal.Model) error {
	return fmt.Errorf("internal dal: updateModel not supported")
}
func (c *internalConn) AssertSchemaAlterations(_ context.Context, _ *dal.Model, _ ...*dal.Alteration) ([]*dal.Alteration, error) {
	return nil, nil
}
func (c *internalConn) ApplyAlteration(_ context.Context, _ *dal.Model, _ ...*dal.Alteration) []error {
	return nil
}

// composeRecordIterator adapts a composeTypes.RecordSet to dal.Iterator.
type composeRecordIterator struct {
	records composeTypes.RecordSet
	pos     int
	err     error
}

func (it *composeRecordIterator) Next(_ context.Context) bool {
	it.pos++
	return it.pos < len(it.records)
}

func (it *composeRecordIterator) Scan(dst dal.ValueSetter) error {
	rec := it.records[it.pos]
	for _, v := range rec.Values {
		if err := dst.SetValue(v.Name, uint(v.Place), v.Value); err != nil {
			return err
		}
	}
	return nil
}

func (it *composeRecordIterator) More(_ uint, _ dal.ValueGetter) error { return nil }
func (it *composeRecordIterator) BackCursor(_ dal.ValueGetter) (*filter.PagingCursor, error) {
	return nil, nil
}
func (it *composeRecordIterator) ForwardCursor(_ dal.ValueGetter) (*filter.PagingCursor, error) {
	return nil, nil
}
func (it *composeRecordIterator) Err() error  { return it.err }
func (it *composeRecordIterator) Close() error { return nil }
