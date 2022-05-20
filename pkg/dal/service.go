package dal

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"go.uber.org/zap"
)

type (
	connectionWrap struct {
		connection Connection
		Defaults   ConnectionDefaults
	}

	ConnectionDefaults struct {
		ModelIdent     string
		AttributeIdent string

		PartitionFormat string
	}

	service struct {
		connections map[uint64]*connectionWrap
		primary     *connectionWrap

		// Indexed by corresponding storeID
		models map[uint64]ModelSet

		logger *zap.Logger
		inDev  bool
	}
)

const (
	DefaultConnectionID uint64 = 0
)

var (
	gSvc *service
)

// InitGlobalService initializes a fresh DAL where the given primary connection
func InitGlobalService(ctx context.Context, log *zap.Logger, inDev bool, dsn string, dft ConnectionDefaults, capabilities ...capabilities.Capability) (*service, error) {
	if gSvc == nil {
		gSvc = &service{
			connections: make(map[uint64]*connectionWrap),
			models:      make(map[uint64]ModelSet),
			primary:     nil,

			logger: log,
			inDev:  inDev,
		}

		var err error
		cw := &connectionWrap{
			Defaults: dft,
		}
		cw.connection, err = connect(ctx, log, inDev, dsn, capabilities...)
		if err != nil {
			return nil, err
		}

		gSvc.primary = cw
	}

	return gSvc, nil
}

// Service returns the global initialized DAL service
//
// If InitGlobalService has not yet been called the function will panic
func Service() *service {
	if gSvc == nil {
		panic("DAL global service not initialized: call dal.InitGlobalService() first")
	}

	return gSvc
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Connection management

// AddConnection adds a new connection to the DAL
func (svc *service) AddConnection(ctx context.Context, connectionID uint64, dsn string, dft ConnectionDefaults, capabilities ...capabilities.Capability) (err error) {
	cw := &connectionWrap{
		Defaults: dft,
	}
	cw.connection, err = connect(ctx, svc.logger, svc.inDev, dsn, capabilities...)
	if err != nil {
		return
	}
	svc.connections[connectionID] = cw
	return
}

// RemoveConnection removes the given connection from the DAL
func (svc *service) RemoveConnection(ctx context.Context, connectionID uint64) (err error) {
	c := svc.connections[connectionID]
	if c == nil {
		return fmt.Errorf("can not remove connection %d: connection does not exist", connectionID)
	}

	// Potential cleanups
	if cc, ok := c.connection.(ConnectionCloser); ok {
		if err = cc.Close(ctx); err != nil {
			return err
		}
	}

	// Remove from registry
	delete(svc.connections, connectionID)

	return nil
}

// UpdateConnection updates the given connection
//
// @todo make this better; for now remove + add
func (svc *service) UpdateConnection(ctx context.Context, connectionID uint64, dsn string, dft ConnectionDefaults, capabilities ...capabilities.Capability) (err error) {
	if err = svc.RemoveConnection(ctx, connectionID); err != nil {
		return
	}

	return svc.AddConnection(ctx, connectionID, dsn, dft, capabilities...)
}

// ConnectionDefaultreturns the defaults we can use with this connection
func (svc *service) ConnectionDefaults(ctx context.Context, connectionID uint64) (dft ConnectionDefaults, err error) {
	wrap, _, err := svc.getConnection(ctx, connectionID)
	if err != nil {
		return
	}

	return wrap.Defaults, nil
}

// // // // // // // // // // // // // // // // // // // // // // // // //

// // // // // // // // // // // // // // // // // // // // // // // // //
// DML

func (svc *service) Create(ctx context.Context, mf ModelFilter, capabilities capabilities.Set, rr ...ValueGetter) (err error) {
	model, cw, err := svc.storeOpPrep(ctx, mf, capabilities)
	if err != nil {
		return
	}

	return cw.connection.Create(ctx, model, rr...)
}

func (svc *service) Update(ctx context.Context, mf ModelFilter, capabilities capabilities.Set, r ValueGetter) (err error) {
	model, cw, err := svc.storeOpPrep(ctx, mf, capabilities)
	if err != nil {
		return
	}

	return cw.connection.Update(ctx, model, r)
}

func (svc *service) Search(ctx context.Context, mf ModelFilter, capabilities capabilities.Set, f filter.Filter) (iter Iterator, err error) {
	model, cw, err := svc.storeOpPrep(ctx, mf, capabilities)
	if err != nil {
		return
	}

	return cw.connection.Search(ctx, model, f)
}

func (svc *service) Lookup(ctx context.Context, mf ModelFilter, capabilities capabilities.Set, lookup ValueGetter, dst ValueSetter) (err error) {
	model, cw, err := svc.storeOpPrep(ctx, mf, capabilities)
	if err != nil {
		return
	}
	return cw.connection.Lookup(ctx, model, lookup, dst)
}

func (svc *service) Delete(ctx context.Context, mf ModelFilter, capabilities capabilities.Set, pkv ValueGetter) (err error) {
	model, cw, err := svc.storeOpPrep(ctx, mf, capabilities)
	if err != nil {
		return
	}

	return cw.connection.Delete(ctx, model, pkv)
}
func (svc *service) Truncate(ctx context.Context, mf ModelFilter, capabilities capabilities.Set) (err error) {
	model, cw, err := svc.storeOpPrep(ctx, mf, capabilities)
	if err != nil {
		return
	}

	return cw.connection.Truncate(ctx, model)
}

func (svc *service) storeOpPrep(ctx context.Context, mf ModelFilter, capabilities capabilities.Set) (model *Model, cw *connectionWrap, err error) {
	model = svc.getModelByFilter(mf)
	if model == nil {
		err = fmt.Errorf("cannot perform operation: model not registered")
		return
	}

	cw, _, err = svc.getConnection(ctx, model.ConnectionID, capabilities...)
	if err != nil {
		return
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //

// // // // // // // // // // // // // // // // // // // // // // // // //
// DDL

// ReloadModel unregister old models and register the new ones
func (svc *service) ReloadModel(ctx context.Context, models ...*Model) (err error) {
	// Clear up the old ones
	// @todo profile if manually removing nested pointers makes it faster
	svc.models = make(map[uint64]ModelSet)
	return svc.AddModel(ctx, models...)
}

// AddModel adds support for a new model
func (svc *service) AddModel(ctx context.Context, models ...*Model) (err error) {
	var (
		cw *connectionWrap
	)

	for connectionID, models := range svc.modelByConnection(models) {
		cw, _, err = svc.getConnection(ctx, connectionID)
		if err != nil {
			return err
		}

		err = svc.registerModel(ctx, cw.connection, connectionID, models)
		if err != nil {
			return
		}
	}

	return
}

// RemoveModel removes support for the given model
func (svc *service) RemoveModel(ctx context.Context, models ...*Model) (err error) {
	// validation
	for _, model := range models {
		// Validate existence
		old := svc.GetModelByResource(model.ConnectionID, model.ResourceType, model.Resource)
		if old == nil {
			return fmt.Errorf("cannot remove module %s: not registered", model.Resource)
		}

		// Validate no leftover references
		// @todo we can probably expand on this quitea bit
		for _, registered := range svc.models {
			refs := registered.FilterByReferenced(model)
			if len(refs) > 0 {
				return fmt.Errorf("cannot remove module %s: referenced by other modules", model.Resource)
			}
		}
	}

	// Work
	for _, model := range models {
		oldModels := svc.models[model.ConnectionID]
		svc.models[model.ConnectionID] = make(ModelSet, 0, len(oldModels))
		for _, o := range oldModels {
			if o.Resource == model.Resource {
				continue
			}

			svc.models[model.ConnectionID] = append(svc.models[model.ConnectionID], o)
		}

		// @todo should the underlying store be notified about this?
	}

	return nil
}

// DeleteModel removes support for the model and deletes it from the connection
//
// @todo do we really want this?
func (svc *service) DeleteModel(ctx context.Context, models ...*Model) (err error) {
	panic("implement DeleteModel")
}

func (svc *service) UpdateModel(ctx context.Context, old *Model, new *Model) error {
	panic("implement UpdateModel")
}

func (svc *service) UpdateModelAttribute(ctx context.Context, sch *Model, old Attribute, new Attribute, trans ...TransformationFunction) error {
	panic("implement UpdateModelAttribute")
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utilities

func (svc *service) GetModelByID(connectionID uint64, id uint64) *Model {
	return svc.models[connectionID].FindByID(id)
}

func (svc *service) GetModelByResource(connectionID uint64, resType string, resource string) *Model {
	return svc.models[connectionID].FindByResource(resType, resource)
}

func (svc *service) getConnection(ctx context.Context, connectionID uint64, cc ...capabilities.Capability) (cw *connectionWrap, can capabilities.Set, err error) {
	err = func() error {
		// get the requested connection
		if connectionID == DefaultConnectionID {
			cw = svc.primary
		} else {
			cw = svc.connections[connectionID]
		}
		if cw == nil {
			return fmt.Errorf("could not get connection %d: store does not exist", connectionID)
		}

		// check if connection supports requested capabilities
		if !cw.connection.Can(cc...) {
			return fmt.Errorf("connection does not support requested capabilities: %v", capabilities.Set(cc).Diff(cw.connection.Capabilities()))
		}
		can = cw.connection.Capabilities()
		return nil
	}()

	if err != nil {
		err = fmt.Errorf("could not connect to %d: %v", connectionID, err)
		return
	}

	return
}

// modelByConnection maps the given models by their CRS
func (svc *service) modelByConnection(models ModelSet) (out map[uint64]ModelSet) {
	out = make(map[uint64]ModelSet)

	for _, model := range models {
		out[model.ConnectionID] = append(out[model.ConnectionID], model)
	}

	return
}

func (svc *service) registerModel(ctx context.Context, s Connection, storeID uint64, models ModelSet) (err error) {
	for _, model := range models {
		existing := svc.GetModelByResource(storeID, model.ResourceType, model.Resource)
		if existing != nil {
			return fmt.Errorf("cannot add model %s to store %d: already exists", model.Resource, storeID)
		}

		err = svc.registerModelToConnection(ctx, s, model)
		if err != nil {
			return
		}

		svc.models[storeID] = append(svc.models[storeID], model)
	}

	return
}

func (svc *service) registerModelToConnection(ctx context.Context, s Connection, model *Model) (err error) {
	available, err := s.Models(ctx)
	if err != nil {
		return err
	}

	// Check if already in there
	if existing := available.FindByResource(model.ResourceType, model.Resource); existing != nil {
		// Assert validity
		diff := existing.Diff(model)
		if len(diff) > 0 {
			return fmt.Errorf("model %s exists: model not compatible: %v", existing.Resource, diff)
		}

		return nil
	}

	// Try to add to store
	err = s.CreateModel(ctx, model)
	if err != nil {
		return
	}

	return nil
}

func (svc *service) getModelByFilter(mf ModelFilter) *Model {
	if mf.ResourceID > 0 {
		return svc.GetModelByID(mf.ConnectionID, mf.ResourceID)
	}
	return svc.GetModelByResource(mf.ConnectionID, mf.ResourceType, mf.Resource)
}
