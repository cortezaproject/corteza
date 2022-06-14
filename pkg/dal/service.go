package dal

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"go.uber.org/zap"
)

type (
	connectionWrap struct {
		connectionID     uint64
		label            string
		sensitivityLevel uint64

		connection Connection
		meta       ConnectionMeta
	}

	ConnectionMeta struct {
		SensitivityLevel uint64
		Label            string

		DefaultModelIdent     string
		DefaultAttributeIdent string

		DefaultPartitionFormat string

		PartitionValidator string
	}

	service struct {
		connections         map[uint64]*connectionWrap
		primaryConnectionID uint64

		// Indexed by corresponding storeID
		models map[uint64]ModelSet

		logger *zap.Logger
		inDev  bool

		sensitivityLevels sensitivityLevelIndex

		connectionIssues dalIssueIndex
		modelIssues      dalIssueIndex
	}
)

const (
	DefaultConnectionID uint64 = 0
)

var (
	gSvc *service
)

// InitGlobalService initializes a fresh DAL where the given primary connection
func InitGlobalService(ctx context.Context, log *zap.Logger, inDev bool, connectionID uint64, cp ConnectionParams, cm ConnectionMeta, capabilities ...capabilities.Capability) (_ *service, err error) {
	log.Debug("initializing DAL service with primary connection", zap.Any("connection params", cp))

	// To help prevent awkward issues due to globally shared resources
	if gSvc != nil {
		panic("cannot initialize global DAL service: already initialized")
	}

	if gSvc == nil {
		gSvc, err = New(ctx, log, inDev, connectionID, cp, cm, capabilities...)
		return gSvc, err
	}

	return gSvc, nil
}

func New(ctx context.Context, log *zap.Logger, inDev bool, connectionID uint64, cp ConnectionParams, cm ConnectionMeta, capabilities ...capabilities.Capability) (*service, error) {
	svc := &service{
		connections:         make(map[uint64]*connectionWrap),
		models:              make(map[uint64]ModelSet),
		primaryConnectionID: connectionID,

		logger: log,
		inDev:  inDev,

		connectionIssues: make(dalIssueIndex),
		modelIssues:      make(dalIssueIndex),
	}

	var err error
	cw := &connectionWrap{
		meta:             cm,
		sensitivityLevel: cm.SensitivityLevel,
		label:            cm.Label,
		connectionID:     connectionID,
	}
	cw.connection, err = connect(ctx, log, inDev, cp, capabilities...)
	if err != nil {
		return nil, err
	}

	svc.connections[connectionID] = cw

	return svc, nil
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

func (svc *service) ReloadSensitivityLevels(levels ...SensitivityLevel) (err error) {
	svc.logger.Debug("reloading sensitivity levels", zap.Any("sensitivity levels", levels))
	newLevelIndex := svc.newSensitivityLevelIndex(levels)

	// Validate state after sensitivity level change
	if err = svc.validateNewSensitivityLevels(newLevelIndex); err != nil {
		return
	}

	// Replace old ones
	svc.sensitivityLevels = newLevelIndex

	svc.logger.Debug("reloaded sensitivity levels")
	return
}

func (svc *service) CreateSensitivityLevel(levels ...SensitivityLevel) (err error) {
	svc.logger.Debug("creating sensitivity levels", zap.Any("sensitivity levels", levels))
	newIndex := svc.newAddedSensitivityLevelIndex(svc.sensitivityLevels, levels...)

	// Validate state after sensitivity level change
	if err = svc.validateNewSensitivityLevels(newIndex); err != nil {
		return
	}

	// Replace old ones
	svc.sensitivityLevels = newIndex
	svc.logger.Debug("created sensitivity levels")
	return
}

func (svc *service) UpdateSensitivityLevel(levels ...SensitivityLevel) (err error) {
	svc.logger.Debug("updating sensitivity levels", zap.Any("sensitivity levels", levels))
	newIndex := svc.newRemovedSensitivityLevelIndex(svc.sensitivityLevels, levels...)
	newIndex = svc.newAddedSensitivityLevelIndex(newIndex, levels...)

	// Validate state after sensitivity level change
	if err = svc.validateNewSensitivityLevels(newIndex); err != nil {
		return
	}

	// Replace old ones
	svc.sensitivityLevels = newIndex
	svc.logger.Debug("updated sensitivity levels")
	return
}

func (svc *service) DeleteSensitivityLevel(levels ...SensitivityLevel) (err error) {
	svc.logger.Debug("deleting sensitivity levels", zap.Any("sensitivity levels", levels))
	newIndex := svc.newRemovedSensitivityLevelIndex(svc.sensitivityLevels, levels...)

	// Validate state after sensitivity level change
	if err = svc.validateNewSensitivityLevels(newIndex); err != nil {
		return
	}

	// Replace old ones
	svc.sensitivityLevels = newIndex
	svc.logger.Debug("deleted sensitivity levels")

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //

// // // // // // // // // // // // // // // // // // // // // // // // //
// Connection management

// CreateConnection adds a new connection to the DAL
func (svc *service) CreateConnection(ctx context.Context, connectionID uint64, cp ConnectionParams, cm ConnectionMeta, capabilities ...capabilities.Capability) (err error) {
	svc.logger.Debug("creating connection", zap.Uint64("connectionID", connectionID), zap.Any("connection params", cp))

	var (
		issues = newIssueHelper().addConnection(connectionID)
	)
	defer svc.updateIssues(issues)

	// sensitivity levels
	if !svc.sensitivityLevels.includes(cm.SensitivityLevel) {
		issues.addConnectionIssue(connectionID, errConnectionCreateMissingSensitivityLevel(connectionID, cm.SensitivityLevel))
	}

	// Prepare connection bits
	cw := &connectionWrap{
		connectionID:     connectionID,
		meta:             cm,
		sensitivityLevel: cm.SensitivityLevel,
		label:            cm.Label,
	}
	if cw.connection, err = connect(ctx, svc.logger, svc.inDev, cp, capabilities...); err != nil {
		issues.addConnectionIssue(connectionID, errConnectionCreateConnectionFailed(connectionID, err))
	}

	svc.addConnection(cw)

	svc.logger.Debug("created connection")
	return nil
}

// DeleteConnection removes the given connection from the DAL
func (svc *service) DeleteConnection(ctx context.Context, connectionID uint64) (err error) {
	svc.logger.Debug("deleting connection", zap.Uint64("connectionID", connectionID))

	var (
		issues = newIssueHelper().addConnection(connectionID)
	)

	c := svc.getConnectionByID(connectionID)
	if c == nil {
		return errConnectionDeleteNotFound(connectionID)
	}

	// Potential cleanups
	if cc, ok := c.connection.(ConnectionCloser); ok {
		if err := cc.Close(ctx); err != nil {
			svc.logger.Error(errConnectionDeleteCloserFailed(c.connectionID, err).Error())
		}
	}

	// Remove from registry
	//
	// @todo this is temporary until a proper update function is prepared.
	// The primary connection must not be removable!
	svc.removeConnection(connectionID)

	// Only if successful should we cleanup the issue registry
	svc.updateIssues(issues)

	svc.logger.Debug("deleted connection")

	return nil
}

// UpdateConnection updates the given connection
func (svc *service) UpdateConnection(ctx context.Context, connectionID uint64, cp ConnectionParams, cm ConnectionMeta, capabilities ...capabilities.Capability) (err error) {
	svc.logger.Debug("updating connection", zap.Uint64("connectionID", connectionID))

	var (
		issues  = newIssueHelper().addConnection(connectionID)
		oldConn *connectionWrap
	)
	defer svc.updateIssues(issues)

	// Validation
	{
		// Check if connection exists
		oldConn = svc.getConnectionByID(connectionID)
		if oldConn == nil {
			issues.addConnectionIssue(connectionID, errConnectionUpdateNotFound(connectionID))
		}

		// sensitivity levels
		if !svc.sensitivityLevels.includes(cm.SensitivityLevel) {
			issues.addConnectionIssue(connectionID, errConnectionUpdateMissingSensitivityLevel(connectionID, cm.SensitivityLevel))
		}

		// Check already registered models and their capabilities
		//
		// Defer the return till the end so we can get a nicer report of what all is wrong
		errored := false
		for _, model := range svc.models[connectionID] {
			// - capabilities
			if !model.Capabilities.IsSubset(capabilities...) {
				issues.addConnectionIssue(connectionID, fmt.Errorf("cannot update connection %d: new connection does not support existing models", connectionID))
				errored = errored || true
			}
			// - sensitivity levels
			if !svc.sensitivityLevels.isSubset(model.SensitivityLevel, cm.SensitivityLevel) {
				issues.addConnectionIssue(connectionID, fmt.Errorf("cannot update connection %d: new connection sensitivity level does not support model %d", connectionID, model.ResourceID))
				errored = errored || true
			}
		}

		// Don't update if meta bits are not ok
		if errored {
			return
		}
	}

	// close old connection
	{
		if cc, ok := oldConn.connection.(ConnectionCloser); ok {
			if err = cc.Close(ctx); err != nil {
				issues.addConnectionIssue(connectionID, err)
				return nil
			}
		}
		svc.removeConnection(connectionID)
	}

	// open new connection
	{
		newConnection, err := connect(ctx, svc.logger, svc.inDev, cp, capabilities...)
		if err != nil {
			issues.addConnectionIssue(connectionID, err)
		}

		svc.addConnection(&connectionWrap{
			meta:             cm,
			sensitivityLevel: cm.SensitivityLevel,
			label:            cm.Label,
			connectionID:     connectionID,
			connection:       newConnection,
		})
	}

	svc.logger.Debug("updated connection")

	return nil
}

// // // // // // // // // // // // // // // // // // // // // // // // //

// // // // // // // // // // // // // // // // // // // // // // // // //
// DML

func (svc *service) Create(ctx context.Context, mf ModelFilter, capabilities capabilities.Set, rr ...ValueGetter) (err error) {
	if err = svc.canOpRecord(mf.ConnectionID, mf.ResourceID); err != nil {
		return wrapError("cannot create record", err)
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, capabilities)
	if err != nil {
		return wrapError("cannot create record", err)
	}

	return cw.connection.Create(ctx, model, rr...)
}

func (svc *service) Update(ctx context.Context, mf ModelFilter, capabilities capabilities.Set, rr ...ValueGetter) (err error) {
	if err = svc.canOpRecord(mf.ConnectionID, mf.ResourceID); err != nil {
		return wrapError("cannot update record", err)
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, capabilities)
	if err != nil {
		return wrapError("cannot update record", err)
	}

	for _, r := range rr {
		if err = cw.connection.Update(ctx, model, r); err != nil {
			return wrapError("cannot update record", err)
		}
	}

	return
}

func (svc *service) Search(ctx context.Context, mf ModelFilter, capabilities capabilities.Set, f filter.Filter) (iter Iterator, err error) {
	if err = svc.canOpRecord(mf.ConnectionID, mf.ResourceID); err != nil {
		err = wrapError("cannot search record", err)
		return
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, capabilities)
	if err != nil {
		err = wrapError("cannot search record", err)
		return
	}

	return cw.connection.Search(ctx, model, f)
}

func (svc *service) Lookup(ctx context.Context, mf ModelFilter, capabilities capabilities.Set, lookup ValueGetter, dst ValueSetter) (err error) {
	if err = svc.canOpRecord(mf.ConnectionID, mf.ResourceID); err != nil {
		return wrapError("cannot lookup record", err)
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, capabilities)
	if err != nil {
		return wrapError("cannot lookup record", err)
	}
	return cw.connection.Lookup(ctx, model, lookup, dst)
}

func (svc *service) Delete(ctx context.Context, mf ModelFilter, capabilities capabilities.Set, vv ...ValueGetter) (err error) {
	if err = svc.canOpRecord(mf.ConnectionID, mf.ResourceID); err != nil {
		return wrapError("cannot delete record", err)
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, capabilities)
	if err != nil {
		return wrapError("cannot delete record", err)
	}

	for _, v := range vv {
		if err = cw.connection.Delete(ctx, model, v); err != nil {
			return wrapError("cannot delete record", err)
		}
	}
	return
}

func (svc *service) Truncate(ctx context.Context, mf ModelFilter, capabilities capabilities.Set) (err error) {
	if err = svc.canOpRecord(mf.ConnectionID, mf.ResourceID); err != nil {
		return wrapError("cannot truncate record", err)
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, capabilities)
	if err != nil {
		return wrapError("cannot truncate record", err)
	}

	return cw.connection.Truncate(ctx, model)
}

func (svc *service) storeOpPrep(ctx context.Context, mf ModelFilter, capabilities capabilities.Set) (model *Model, cw *connectionWrap, err error) {
	model = svc.getModelByFilter(mf)
	if model == nil {
		err = errModelNotFound(mf.ResourceID)
		return
	}

	cw, _, err = svc.getConnection(model.ConnectionID, capabilities...)
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
	svc.logger.Debug("reloading models", zap.Int("count", len(models)))

	// Clear up the old ones
	// @todo profile if manually removing nested pointers makes it faster
	svc.models = make(map[uint64]ModelSet)
	svc.clearModelIssues()

	err = svc.CreateModel(ctx, models...)
	if err != nil {
		return
	}

	svc.logger.Debug("reloaded models")

	return
}

func (svc *service) SearchModels(ctx context.Context) (out ModelSet, err error) {
	out = make(ModelSet, 0, 100)
	for _, models := range svc.models {
		out = append(out, models...)
	}
	return
}

// AddModel adds support for a new model
func (svc *service) CreateModel(ctx context.Context, models ...*Model) (err error) {
	svc.logger.Debug("creating models", zap.Int("count", len(models)))

	var (
		issues    = newIssueHelper()
		auxIssues = newIssueHelper()
	)
	defer svc.updateIssues(issues)

	// Validate models
	for _, model := range models {
		svc.logger.Debug("validating model", zap.Uint64("ID", model.ResourceID))
		issues.addModel(model.ResourceID)

		// Assure the connection has no issues
		if svc.hasConnectionIssues(model.ConnectionID) {
			issues.addModelIssue(model.ConnectionID, model.ResourceID, errModelCreateProblematicConnection(model.ConnectionID, model.ResourceID))
		}

		// Check the connection exists
		conn := svc.getConnectionByID(model.ConnectionID)
		if conn == nil {
			issues.addModelIssue(model.ConnectionID, model.ResourceID, errModelCreateMissingConnection(model.ConnectionID, model.ResourceID))
		}

		// Check if model for the given resource already exists
		existing := svc.FindModelByResourceID(model.ConnectionID, model.ResourceID)
		if existing != nil {
			issues.addModelIssue(model.ConnectionID, model.ResourceID, errModelCreateDuplicate(model.ConnectionID, model.ResourceID))
		}

		// Check sensitivity levels
		// - model
		if !svc.sensitivityLevels.includes(model.SensitivityLevel) {
			issues.addModelIssue(model.ConnectionID, model.ResourceID, errModelCreateMissingSensitivityLevel(model.ConnectionID, model.ResourceID, model.SensitivityLevel))
		} else {
			// Only check if it is present
			if !svc.sensitivityLevels.isSubset(model.SensitivityLevel, conn.sensitivityLevel) {
				issues.addModelIssue(model.ConnectionID, model.ResourceID, errModelCreateGreaterSensitivityLevel(model.ConnectionID, model.ResourceID, model.SensitivityLevel, conn.sensitivityLevel))
			}
		}
		// - attributes
		for _, attr := range model.Attributes {
			if !svc.sensitivityLevels.includes(attr.SensitivityLevel) {
				issues.addModelIssue(model.ConnectionID, model.ResourceID, errModelCreateMissingAttributeSensitivityLevel(model.ConnectionID, model.ResourceID, attr.SensitivityLevel))
			} else {
				if !svc.sensitivityLevels.isSubset(attr.SensitivityLevel, model.SensitivityLevel) {
					issues.addModelIssue(model.ConnectionID, model.ResourceID, errModelCreateGreaterAttributeSensitivityLevel(model.ConnectionID, model.ResourceID, attr.SensitivityLevel, model.SensitivityLevel))
				}
			}
		}

		svc.logger.Debug("validated model")
	}

	// Add models to corresponding connections
	for connection, models := range svc.modelByConnection(models) {
		for _, model := range models {
			svc.logger.Debug("adding model", zap.Uint64("model", model.ResourceID))

			connectionIssues := svc.hasConnectionIssues(model.ConnectionID)
			modelIssues := svc.hasModelIssues(model.ConnectionID, model.ResourceID)

			if !modelIssues && !connectionIssues {
				svc.logger.Debug("adding model to connection", zap.Uint64("connection", model.ConnectionID), zap.Uint64("model", model.ResourceID))

				// Add model to connection
				auxIssues, err = svc.registerModelToConnection(ctx, connection, model)
				issues.mergeWith(auxIssues)
				if err != nil {
					return
				}
			} else {
				if connectionIssues {
					svc.logger.Warn("not adding to connection due to connection issues", zap.Uint64("connection", model.ConnectionID))
				}
				if modelIssues {
					svc.logger.Warn("not adding to connection due to model issues", zap.Uint64("model", model.ResourceID))
				}
			}

			// Add model to internal registry
			svc.models[model.ConnectionID] = append(svc.models[model.ConnectionID], model)
		}
	}

	svc.logger.Debug("created models")

	return
}

// DeleteModel removes support for the model and deletes it from the connection
func (svc *service) DeleteModel(ctx context.Context, models ...*Model) (err error) {
	svc.logger.Debug("deleting models", zap.Int("count", len(models)))

	var (
		issues = newIssueHelper()
	)
	defer svc.updateIssues(issues)

	// validation
	skip := make(map[uint64]bool)
	for _, model := range models {
		issues.addModel(model.ResourceID)

		// Validate existence
		old := svc.FindModelByResourceIdent(model.ConnectionID, model.ResourceType, model.Resource)
		if old == nil {
			skip[model.ResourceID] = true
			continue
		}

		// Validate no leftover references
		// @todo we can probably expand on this quitea bit
		// for _, registered := range svc.models {
		// 	refs := registered.FilterByReferenced(model)
		// 	if len(refs) > 0 {
		// 		return fmt.Errorf("cannot remove model %s: referenced by other models", model.Resource)
		// 	}
		// }
	}

	// Work
	for _, model := range models {
		svc.logger.Debug("deleting model", zap.Uint64("model", model.ResourceID))

		if skip[model.ResourceID] {
			svc.logger.Debug("model does not exist; skipping")
			continue
		}

		oldModels := svc.models[model.ConnectionID]
		svc.models[model.ConnectionID] = make(ModelSet, 0, len(oldModels))
		for _, o := range oldModels {
			if o.Resource == model.Resource {
				continue
			}

			svc.models[model.ConnectionID] = append(svc.models[model.ConnectionID], o)
		}

		// @todo should the underlying store be notified about this?
		// how should this be handled; a straight up delete doesn't sound sane to me
		// anymore

		svc.logger.Debug("deleted model")
	}

	return nil
}

func (svc *service) UpdateModel(ctx context.Context, old *Model, new *Model) (err error) {
	svc.logger.Debug("updating model", zap.Uint64("model", old.ResourceID))

	var (
		conn *connectionWrap

		issues = newIssueHelper().addModel(old.ResourceID)
	)
	defer svc.updateIssues(issues)

	// Validation
	{
		// Assure the connection has no issues
		if svc.hasConnectionIssues(old.ConnectionID) {
			issues.addModelIssue(old.ConnectionID, old.ResourceID, errModelUpdateProblematicConnection(old.ConnectionID, old.ResourceID))
		}
		// Check the connection exists
		conn = svc.getConnectionByID(old.ConnectionID)
		if conn == nil {
			issues.addModelIssue(old.ConnectionID, old.ResourceID, errModelUpdateMissingConnection(old.ConnectionID, old.ResourceID))
		}

		// Check if old one exists
		if tmp := svc.FindModelByResourceID(old.ConnectionID, old.ResourceID); tmp == nil {
			issues.addModelIssue(old.ConnectionID, old.ResourceID, errModelUpdateMissingOldModel(old.ConnectionID, old.ResourceID))
		}

		// Check if the new one can fit in
		// - if ident changed, check if it's duplicated
		if old.Ident != new.Ident {
			if tmp := svc.FindModelByIdent(new.ConnectionID, new.Ident); tmp == nil {
				issues.addModelIssue(old.ConnectionID, old.ResourceID, errModelUpdateDuplicate(new.ConnectionID, new.ResourceID))
			}
		}
		// - assure same connection
		//   @todo some migration between different connections
		if old.ConnectionID != new.ConnectionID {
			issues.addModelIssue(new.ConnectionID, new.ResourceID, errModelUpdateConnectionMissmatch(old.ConnectionID, old.ResourceID))
		}

		// Sensitivity levels
		// - model
		if !svc.sensitivityLevels.includes(new.SensitivityLevel) {
			issues.addModelIssue(new.ConnectionID, new.ResourceID, errModelUpdateMissingSensitivityLevel(new.ConnectionID, new.ResourceID, new.SensitivityLevel))
		} else {
			if !svc.sensitivityLevels.isSubset(new.SensitivityLevel, conn.sensitivityLevel) {
				issues.addModelIssue(new.ConnectionID, new.ResourceID, errModelUpdateGreaterSensitivityLevel(new.ConnectionID, new.ResourceID, new.SensitivityLevel, conn.sensitivityLevel))
			}
		}

		// @note attribute check should be done in update model attribute so it's omitted here
	}

	// Update connection
	connectionIssues := svc.hasConnectionIssues(new.ConnectionID)
	modelIssues := svc.hasModelIssues(new.ConnectionID, new.ResourceID)

	if !modelIssues && !connectionIssues {
		svc.logger.Debug("updating connection's model", zap.Uint64("connection", new.ConnectionID), zap.Uint64("model", new.ResourceID))

		err = conn.connection.UpdateModel(ctx, old, new)
		if err != nil {
			issues.addModelIssue(new.ConnectionID, new.ResourceID, err)
		}
	} else {
		if connectionIssues {
			svc.logger.Warn("not updating connection's model due to connection issues", zap.Uint64("connection", new.ConnectionID))
		}
		if modelIssues {
			svc.logger.Warn("not updating connection's model due to model issues", zap.Uint64("model", new.ResourceID))
		}
	}

	// Update registry
	ok := false
	for i, model := range svc.models[old.ConnectionID] {
		if model.ResourceID == old.ResourceID {
			svc.models[old.ConnectionID][i] = new
			ok = true
			break
		}
	}
	if !ok {
		svc.models[old.ConnectionID] = append(svc.models[old.ConnectionID], new)
	}

	svc.logger.Debug("updated model")

	return
}

func (svc *service) UpdateModelAttribute(ctx context.Context, model *Model, old, new *Attribute, trans ...TransformationFunction) (err error) {
	svc.logger.Debug("updating model attribute", zap.Uint64("model", model.ResourceID))

	var (
		conn   *connectionWrap
		issues = newIssueHelper().addModel(model.ResourceID)
	)
	defer svc.updateIssues(issues)

	// Validation
	{
		// Connection issues
		if svc.hasConnectionIssues(model.ConnectionID) {
			issues.addModelIssue(model.ConnectionID, model.ResourceID, errAttributeUpdateProblematicConnection(model.ConnectionID, model.ResourceID))
		}

		// Check if it exists
		model := svc.FindModelByResourceID(model.ConnectionID, model.ResourceID)
		if model == nil {
			issues.addModelIssue(model.ConnectionID, model.ResourceID, errAttributeUpdateMissingModel(model.ConnectionID, model.ResourceID))
		}

		// In case we're deleting it we can ignore this check
		if new != nil {
			if !svc.sensitivityLevels.includes(new.SensitivityLevel) {
				issues.addModelIssue(model.ConnectionID, model.ResourceID, errAttributeUpdateMissingSensitivityLevel(model.ConnectionID, model.ResourceID, new.SensitivityLevel))
			} else {
				if !svc.sensitivityLevels.isSubset(new.SensitivityLevel, model.SensitivityLevel) {
					issues.addModelIssue(model.ConnectionID, model.ResourceID, errAttributeUpdateGreaterSensitivityLevel(model.ConnectionID, model.ResourceID, new.SensitivityLevel, model.SensitivityLevel))
				}
			}
		}

		conn = svc.getConnectionByID(model.ConnectionID)
	}

	// Update attribute
	// Update connection
	connectionIssues := svc.hasConnectionIssues(model.ConnectionID)
	modelIssues := svc.hasModelIssues(model.ConnectionID, model.ResourceID)

	if !modelIssues && !connectionIssues {
		svc.logger.Debug("updating model attribute", zap.Uint64("connection", model.ConnectionID), zap.Uint64("model", model.ResourceID))

		err = conn.connection.UpdateModelAttribute(ctx, model, old, new, trans...)
		if err != nil {
			issues.addModelIssue(model.ConnectionID, model.ResourceID, err)
		}
	} else {
		if connectionIssues {
			svc.logger.Warn("not updating model attribute due to connection issues", zap.Uint64("connection", model.ConnectionID))
		}
		if modelIssues {
			svc.logger.Warn("not updating model attribute due to model issues", zap.Uint64("model", model.ResourceID))
		}
	}

	// Update registry
	if old == nil {
		// adding
		model.Attributes = append(model.Attributes, new)
	} else if new == nil {
		// removing
		model = svc.FindModelByResourceID(model.ConnectionID, model.ResourceID)
		nSet := make(AttributeSet, 0, len(model.Attributes))
		for _, attribute := range model.Attributes {
			if attribute.Ident != old.Ident {
				nSet = append(nSet, attribute)
			}
		}
		model.Attributes = nSet
	} else {
		// updating
		model = svc.FindModelByResourceID(model.ConnectionID, model.ResourceID)
		for i, attribute := range model.Attributes {
			if attribute.Ident == old.Ident {
				model.Attributes[i] = new
				break
			}
		}
	}

	svc.logger.Debug("updated model attribute")

	return
}

func (svc *service) ModelIdentFormatter(connectionID uint64) (f *IdentFormatter, err error) {
	c := svc.getConnectionByID(connectionID)

	if c == nil {
		err = errConnectionNotFound(connectionID)
		return
	}

	f = &IdentFormatter{
		defaultModelIdent:      c.meta.DefaultModelIdent,
		defaultAttributeIdent:  c.meta.DefaultAttributeIdent,
		defaultPartitionFormat: c.meta.DefaultPartitionFormat,
	}

	if c.meta.PartitionValidator != "" {
		f.partitionFormatValidator, err = expr.Parser().NewEvaluable(c.meta.PartitionValidator)
		if err != nil {
			return
		}
	}

	return
}

func (svc *service) FindModelByResourceID(connectionID uint64, resourceID uint64) *Model {
	return svc.models[connectionID].FindByResourceID(resourceID)
}

func (svc *service) FindModelByResourceIdent(connectionID uint64, resourceType, resourceIdent string) *Model {
	return svc.models[connectionID].FindByResourceIdent(resourceType, resourceIdent)
}

func (svc *service) FindModelByIdent(connectionID uint64, ident string) *Model {
	return svc.models[connectionID].FindByIdent(ident)
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utilities

func (svc *service) getConnectionByID(connectionID uint64) (cw *connectionWrap) {
	if connectionID == DefaultConnectionID || connectionID == svc.primaryConnectionID {
		return svc.connections[svc.primaryConnectionID]
	}

	return svc.connections[connectionID]
}

func (svc *service) getConnection(connectionID uint64, cc ...capabilities.Capability) (cw *connectionWrap, can capabilities.Set, err error) {
	err = func() error {
		// get the requested connection
		cw = svc.getConnectionByID(connectionID)
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
func (svc *service) modelByConnection(models ModelSet) (out map[*connectionWrap]ModelSet) {
	out = make(map[*connectionWrap]ModelSet)

	for _, model := range models {
		c := svc.getConnectionByID(model.ConnectionID)
		out[c] = append(out[c], model)
	}

	return
}

func (svc *service) registerModelToConnection(ctx context.Context, cw *connectionWrap, model *Model) (issues *issueHelper, err error) {
	issues = newIssueHelper()

	available, err := cw.connection.Models(ctx)
	if err != nil {
		issues.addModelIssue(model.ConnectionID, model.ResourceID, err)
		return issues, nil
	}

	// Check if already in there
	if existing := available.FindByResourceIdent(model.ResourceType, model.Resource); existing != nil {
		// Assert validity
		diff := existing.Diff(model)
		if len(diff) > 0 {
			issues.addModelIssue(model.ConnectionID, model.ResourceID, errModelCreateConnectionModelUnsupported(model.ConnectionID, model.ResourceID))
			return issues, nil
		}

		return
	}

	// Try to add to store
	err = cw.connection.CreateModel(ctx, model)
	if err != nil {
		issues.addModelIssue(model.ConnectionID, model.ResourceID, err)
		return issues, nil
	}

	return nil, nil
}

func (svc *service) getModelByFilter(mf ModelFilter) *Model {
	if mf.ResourceID > 0 {
		return svc.FindModelByResourceID(mf.ConnectionID, mf.ResourceID)
	}
	return svc.FindModelByResourceIdent(mf.ConnectionID, mf.ResourceType, mf.Resource)
}

func (svc *service) newAddedSensitivityLevelIndex(sli sensitivityLevelIndex, add ...SensitivityLevel) (out sensitivityLevelIndex) {
	newLevels := make(SensitivityLevelSet, 0, len(sli.set)+len(add))

	var (
		i = 0
		j = 0
	)

	for i < len(sli.set) {
		for j < len(add) {
			if sli.set[i].Level <= add[j].Level {
				newLevels = append(newLevels, sli.set[i])
				i++
			}
			if sli.set[i].Level > add[j].Level {
				newLevels = append(newLevels, add[j])
				j++
			}
		}
	}

	if j < len(add)-1 {
		newLevels = append(newLevels, add[j:]...)
	}

	return svc.newSensitivityLevelIndex(newLevels)
}

func (svc *service) newRemovedSensitivityLevelIndex(sli sensitivityLevelIndex, remove ...SensitivityLevel) (out sensitivityLevelIndex) {
	newLevels := make(SensitivityLevelSet, 0, len(sli.set)+len(remove))

	removeSet := SensitivityLevelSet(remove)

	for _, l := range sli.set {
		if !removeSet.includes(l.ID) {
			newLevels = append(newLevels, l)
		}
	}

	return svc.newSensitivityLevelIndex(newLevels)
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

func (svc *service) removeConnection(connectionID uint64) {
	if connectionID == DefaultConnectionID || connectionID == svc.primaryConnectionID {
		connectionID = svc.primaryConnectionID
	}

	delete(svc.connections, connectionID)
}

func (svc *service) addConnection(cw *connectionWrap) {
	svc.connections[cw.connectionID] = cw
}

func wrapError(pfx string, err error) error {
	return fmt.Errorf("%s: %v", pfx, err)
}
