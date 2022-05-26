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
		connectionID     uint64
		label            string
		sensitivityLevel uint64

		connection Connection
		Defaults   ConnectionDefaults
	}

	ConnectionMeta struct {
		ConnectionDefaults

		SensitivityLevel uint64
		Label            string
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

		sensitivityLevels sensitivityLevelIndex
	}
)

const (
	DefaultConnectionID uint64 = 0
)

var (
	gSvc *service
)

// InitGlobalService initializes a fresh DAL where the given primary connection
func InitGlobalService(ctx context.Context, log *zap.Logger, inDev bool, cp ConnectionParams, cm ConnectionMeta, capabilities ...capabilities.Capability) (*service, error) {
	if gSvc == nil {
		log.Debug("initializing DAL service with primary connection", zap.Any("connection params", cp))

		gSvc = &service{
			connections: make(map[uint64]*connectionWrap),
			models:      make(map[uint64]ModelSet),
			primary:     nil,

			logger: log,
			inDev:  inDev,
		}

		var err error
		cw := &connectionWrap{
			Defaults:         cm.ConnectionDefaults,
			sensitivityLevel: cm.SensitivityLevel,
			label:            cm.Label,
		}
		cw.connection, err = connect(ctx, log, inDev, cp, capabilities...)
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
// meta

func (svc *service) Drivers() (drivers []Driver) {
	for _, d := range registeredDrivers {
		drivers = append(drivers, d)
	}

	return
}

func (svc *service) ReloadSensitivityLevels(levels SensitivityLevelSet) (err error) {
	svc.logger.Debug("reloading sensitivity levels", zap.Any("sensitivity levels", levels))
	newLevelIndex := svc.newSensitivityLevelIndex(levels)

	// Validate state after sensitivity level change
	if err = svc.validateNewSensitivityLevels(newLevelIndex); err != nil {
		return
	}

	// Replace old ones
	svc.sensitivityLevels = newLevelIndex

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //

// // // // // // // // // // // // // // // // // // // // // // // // //
// Connection management

// AddConnection adds a new connection to the DAL
func (svc *service) AddConnection(ctx context.Context, connectionID uint64, cp ConnectionParams, cm ConnectionMeta, capabilities ...capabilities.Capability) (err error) {
	svc.logger.Debug("adding new connection", zap.Uint64("connectionID", connectionID), zap.Any("connection params", cp))

	cw := &connectionWrap{
		connectionID:     connectionID,
		Defaults:         cm.ConnectionDefaults,
		sensitivityLevel: cm.SensitivityLevel,
		label:            cm.Label,
	}

	cw.connection, err = connect(ctx, svc.logger, svc.inDev, cp, capabilities...)
	if err != nil {
		return
	}
	svc.connections[connectionID] = cw
	return
}

// RemoveConnection removes the given connection from the DAL
func (svc *service) RemoveConnection(ctx context.Context, connectionID uint64) (err error) {
	svc.logger.Debug("removing connection", zap.Uint64("connectionID", connectionID))

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
func (svc *service) UpdateConnection(ctx context.Context, connectionID uint64, cp ConnectionParams, cm ConnectionMeta, capabilities ...capabilities.Capability) (err error) {
	svc.logger.Debug("updating connection", zap.Uint64("connectionID", connectionID))

	if err = svc.RemoveConnection(ctx, connectionID); err != nil {
		return
	}
	// @todo check sensitivity level against modules

	return svc.AddConnection(ctx, connectionID, cp, cm, capabilities...)
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
		err = fmt.Errorf("model not found")
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
	svc.logger.Debug("reloading models")

	// Clear up the old ones
	// @todo profile if manually removing nested pointers makes it faster
	svc.models = make(map[uint64]ModelSet)
	return svc.AddModel(ctx, models...)
}

// AddModel adds support for a new model
func (svc *service) AddModel(ctx context.Context, models ...*Model) (err error) {
	svc.logger.Debug("adding model", zap.Int("count", len(models)))

	var (
		cw *connectionWrap
	)

	for connectionID, models := range svc.modelByConnection(models) {
		cw, _, err = svc.getConnection(ctx, connectionID)
		if err != nil {
			return err
		}

		err = svc.registerModel(ctx, cw, connectionID, models)
		if err != nil {
			return
		}
	}

	return
}

// RemoveModel removes support for the given model
func (svc *service) RemoveModel(ctx context.Context, models ...*Model) (err error) {
	svc.logger.Debug("removing models", zap.Int("count", len(models)))

	// validation
	for _, model := range models {
		svc.logger.Debug("removing model", zap.String("resource type", model.ResourceType), zap.String("resource model", model.Resource))

		// Validate existence
		old := svc.GetModelByResource(model.ConnectionID, model.ResourceType, model.Resource)
		if old == nil {
			return fmt.Errorf("cannot remove model %s: model not found", model.Resource)
		}

		// Validate no leftover references
		// @todo we can probably expand on this quitea bit
		for _, registered := range svc.models {
			refs := registered.FilterByReferenced(model)
			if len(refs) > 0 {
				return fmt.Errorf("cannot remove model %s: referenced by other models", model.Resource)
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
			return fmt.Errorf("connection %d does not exist", connectionID)
		}

		// check if connection supports requested capabilities
		if !cw.connection.Can(cc...) {
			return fmt.Errorf("connection %d does not support requested capabilities %v", connectionID, capabilities.Set(cc).Diff(cw.connection.Capabilities()))
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

func (svc *service) registerModel(ctx context.Context, cw *connectionWrap, connectionID uint64, models ModelSet) (err error) {
	for _, model := range models {
		svc.logger.Debug("adding model for connection", zap.Uint64("connectionID", connectionID), zap.String("resource type", model.ResourceType), zap.String("resource model", model.Resource))

		existing := svc.GetModelByResource(connectionID, model.ResourceType, model.Resource)
		if existing != nil {
			return fmt.Errorf("cannot add model %s to store %d: model already exists", model.Resource, connectionID)
		}

		err = svc.registerModelToConnection(ctx, cw, model)
		if err != nil {
			return
		}

		svc.models[connectionID] = append(svc.models[connectionID], model)
	}

	return
}

func (svc *service) registerModelToConnection(ctx context.Context, cw *connectionWrap, model *Model) (err error) {
	available, err := cw.connection.Models(ctx)
	if err != nil {
		return err
	}

	// Check if already in there
	if existing := available.FindByResource(model.ResourceType, model.Resource); existing != nil {
		// Assert validity
		diff := existing.Diff(model)
		if len(diff) > 0 {
			return fmt.Errorf("cannot add model %d: model already exists for connection %d: models not compatible: %v", existing.ResourceID, cw.connectionID, diff)
		}

		return nil
	}

	// Validate model against connection
	{
		if !svc.sensitivityLevels.isSubset(model.SensitivityLevel, cw.sensitivityLevel) {
			return errModelHigherSensitivity(model.Label, cw.label)
		}

		for _, attr := range model.Attributes {
			if !svc.sensitivityLevels.isSubset(attr.SensitivityLevel, model.SensitivityLevel) {
				return errAttributeHigherSensitivity(model.Label, attr.Label)
			}
		}
	}

	// Try to add to store
	err = cw.connection.CreateModel(ctx, model)
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

func (svc *service) newSensitivityLevelIndex(levels SensitivityLevelSet) (out sensitivityLevelIndex) {
	out = sensitivityLevelIndex{
		byID:     make(map[uint64]int),
		byHandle: make(map[string]int),
		set:      make(SensitivityLevelSet, len(levels)),
	}

	for i, l := range levels {
		out.set[i] = l

		out.byID[l.ID] = i
		out.byHandle[l.Handle] = i
	}

	return
}

func (svc *service) validateNewSensitivityLevels(levels sensitivityLevelIndex) (err error) {
	err = func() (err error) {
		cIndex := make(map[uint64]*connectionWrap)

		// - connections
		for _, _c := range svc.connections {
			c := _c
			cIndex[c.connectionID] = c

			if !levels.includes(c.sensitivityLevel) {
				return fmt.Errorf("connection sensitivity level missing %d", c.sensitivityLevel)
			}
		}

		// - models
		for _, mm := range svc.models {
			for _, m := range mm {
				if !levels.includes(m.SensitivityLevel) {
					return fmt.Errorf("model sensitivity level missing %d", m.SensitivityLevel)
				}
				if !levels.isSubset(m.SensitivityLevel, cIndex[m.ConnectionID].sensitivityLevel) {
					return fmt.Errorf("model sensitivity level missing %d", m.SensitivityLevel)
				}

				for _, attr := range m.Attributes {
					if !levels.includes(attr.SensitivityLevel) {
						return fmt.Errorf("attribute sensitivity level missing %d", attr.SensitivityLevel)
					}
					if !levels.isSubset(attr.SensitivityLevel, m.SensitivityLevel) {
						return fmt.Errorf("attribute sensitivity level %d greater then model sensitivity level %d", attr.SensitivityLevel, m.SensitivityLevel)
					}
				}
			}
		}
		return
	}()

	if err != nil {
		return fmt.Errorf("cannot reload sensitivity levels: %v", err)
	}
	return
}
