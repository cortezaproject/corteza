package dal

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"go.uber.org/zap"
)

type (
	ConnectionWrap struct {
		connectionID uint64

		connection   Connection
		params       ConnectionParams
		meta         ConnectionMeta
		capabilities capabilities.Set
	}

	ConnectionMeta struct {
		SensitivityLevel uint64
		Label            string

		// When model does not specifiy the ident (table name for example), fallback to this
		// @todo we can lose "Default" prefix
		// @todo do we need a separate setting or can we get away with using just PartitionFormat
		DefaultModelIdent string

		// If model attribute(s) do not specify
		// @todo needs to be more explicit that this is for  JSON encode attributes
		// @todo we can lose "Default" prefix
		DefaultAttributeIdent string

		// If data is partitioned we fallback to this,
		// @todo we can lose "Default" prefix
		DefaultPartitionFormat string

		PartitionValidator string
	}

	service struct {
		connections map[uint64]*ConnectionWrap

		// Default connection ID
		// Can not be changed in the runtime, only set to value different than zero!
		defConnID uint64

		// Indexed by corresponding storeID
		models map[uint64]ModelSet

		logger *zap.Logger
		inDev  bool

		sensitivityLevels *sensitivityLevelIndex

		connectionIssues dalIssueIndex
		modelIssues      dalIssueIndex
	}
)

var (
	gSvc *service
)

func SetGlobal(svc *service) {
	gSvc = svc
}

// New creates a DAL service with the primary connection
//
// It needs an established and working connection to the primary store
func New(log *zap.Logger, inDev bool) (*service, error) {
	svc := &service{
		connections:       make(map[uint64]*ConnectionWrap),
		models:            make(map[uint64]ModelSet),
		sensitivityLevels: SensitivityLevelIndex(),

		logger: log,
		inDev:  inDev,

		connectionIssues: make(dalIssueIndex),
		modelIssues:      make(dalIssueIndex),
	}

	return svc, nil
}

// Service returns the global initialized DAL service
//
// Function will panic if DAL service is not set (via SetGlobal)
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

func MakeSensitivityLevel(ID uint64, level int, handle string) SensitivityLevel {
	return SensitivityLevel{
		ID:     ID,
		Level:  level,
		Handle: handle,
	}
}

func (svc *service) ReplaceSensitivityLevel(levels ...SensitivityLevel) (err error) {
	var (
		log = svc.logger.Named("sensitivity level")
	)

	log.Debug("replacing levels", zap.Any("levels", levels))

	if svc.sensitivityLevels == nil {
		svc.sensitivityLevels = SensitivityLevelIndex()
	}
	nx := svc.sensitivityLevels

	for _, l := range levels {
		log := log.With(zap.Uint64("ID", l.ID), zap.Int("level", l.Level), zap.String("handle", l.Handle))
		if nx.includes(l.ID) {
			log.Debug("found existing")
		} else {
			log.Debug("adding new")
		}
	}

	nx = svc.sensitivityLevels.with(levels...)

	// Validate state after sensitivity level change
	log.Debug("validating new levels")
	if err = svc.validateNewSensitivityLevels(nx); err != nil {
		return
	}

	// Replace the old one
	svc.sensitivityLevels = nx

	svc.logger.Debug("reloaded sensitivity levels")
	return
}

func (svc *service) RemoveSensitivityLevel(levelIDs ...uint64) (err error) {
	var (
		log = svc.logger.Named("sensitivity level")
	)

	log.Debug("removing levels", zap.Any("levels", levelIDs))

	levels := make(SensitivityLevelSet, len(levelIDs))
	for i, lID := range levelIDs {
		levels[i] = MakeSensitivityLevel(lID, i, strconv.FormatUint(lID, 10))
	}

	if svc.sensitivityLevels == nil {
		svc.sensitivityLevels = SensitivityLevelIndex()
	}
	nx := svc.sensitivityLevels

	for _, l := range levels {
		log := log.With(zap.Uint64("ID", l.ID))
		if !nx.includes(l.ID) {
			log.Debug("sensitivity level not found")
			return errSensitivityLevelRemoveNotFound(l.ID)
		}
	}

	nx = svc.sensitivityLevels.without(levels...)

	// Validate state after sensitivity level change
	log.Debug("validating new levels")
	if err = svc.validateNewSensitivityLevels(nx); err != nil {
		return
	}

	// Replace the old one
	svc.sensitivityLevels = nx

	svc.logger.Debug("removed sensitivity levels")
	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //

// // // // // // // // // // // // // // // // // // // // // // // // //
// Connection management

// MakeConnection makes and returns a new connection (wrap)
func MakeConnection(ID uint64, conn Connection, p ConnectionParams, m ConnectionMeta, cap ...capabilities.Capability) *ConnectionWrap {
	return &ConnectionWrap{
		connectionID: ID,
		connection:   conn,

		params:       p,
		meta:         m,
		capabilities: cap,
	}
}

// ReplaceConnection adds new or updates an existing connection
//
// We rely on the user to provide stable connection IDs and
// uses valid relations to these connections in the models.
//
// Is isDefault when adding a default connection. Service will then
// compensate and use proper IDs when models refer to connection with ID=0
func (svc *service) ReplaceConnection(ctx context.Context, cw *ConnectionWrap, isDefault bool) (err error) {
	// @todo lock/unlock
	var (
		ID      = cw.connectionID
		issues  = newIssueHelper().addConnection(ID)
		oldConn *ConnectionWrap

		log = svc.logger.Named("connection").With(
			zap.Uint64("ID", ID),
			zap.Any("params", cw.params),
			zap.Any("meta", cw.meta),
		)
	)

	if isDefault {
		if svc.defConnID == 0 {
			// default connection not set yet
			log.Debug("setting as default connection")
			svc.defConnID = ID
		} else if svc.defConnID != ID {
			// default connection set but ID is different.
			// this does not make any sense
			return fmt.Errorf("different ID for default connection detected (old: %d, new: %d)", svc.defConnID, ID)
		}
	}

	defer svc.updateIssues(issues)

	// Sensitivity level validations
	if !svc.sensitivityLevels.includes(cw.meta.SensitivityLevel) {
		issues.addConnectionIssue(ID, errConnectionCreateMissingSensitivityLevel(ID, cw.meta.SensitivityLevel))
	}

	if oldConn = svc.getConnectionByID(ID); oldConn != nil {
		// Connection exists, validate models and sensitivity levels and close and remove connection at the end
		log.Debug("found existing")

		// Check already registered models and their capabilities
		//
		// Defer the return till the end so we can get a nicer report of what all is wrong
		errored := false
		for _, model := range svc.models[ID] {
			log.Debug("validating model before connection is updated", zap.String("ident", model.Ident))

			// - capabilities
			if !model.Capabilities.IsSubset(cw.capabilities...) {
				issues.addConnectionIssue(ID, fmt.Errorf("cannot update connection %d: new connection does not support existing models", ID))
				errored = true
			}

			// - sensitivity levels
			if !svc.sensitivityLevels.isSubset(model.SensitivityLevel, cw.meta.SensitivityLevel) {
				issues.addConnectionIssue(ID, fmt.Errorf("cannot update connection %d: new connection sensitivity level does not support model %d", ID, model.ResourceID))
				errored = true
			}
		}

		// Don't update if meta bits are not ok
		if errored {
			log.Warn("update failed")
			return
		}

		// close old connection
		if cc, ok := oldConn.connection.(ConnectionCloser); ok {
			if err = cc.Close(ctx); err != nil {
				issues.addConnectionIssue(ID, err)
				return nil
			}

			log.Debug("disconnected")
		}

		svc.removeConnection(ID)
	}

	if cw.connection == nil {
		cw.connection, err = connect(ctx, svc.logger, svc.inDev, cw.params, cw.capabilities...)
		if err != nil {
			log.Warn("could not connect", zap.Error(err))
			issues.addConnectionIssue(ID, err)
		} else {
			log.Debug("connected")
		}
	} else {
		log.Debug("using preexisting connection")

	}

	svc.addConnection(cw)
	log.Debug("added")
	return nil
}

// RemoveConnection removes the given connection from the DAL
func (svc *service) RemoveConnection(ctx context.Context, ID uint64) (err error) {
	var (
		issues = newIssueHelper().addConnection(ID)
	)

	c := svc.getConnectionByID(ID)
	if c == nil {
		return errConnectionDeleteNotFound(ID)
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
	svc.removeConnection(ID)

	// Only if successful should we cleanup the issue registry
	svc.updateIssues(issues)

	svc.logger.Named("connection").Debug("deleted",
		zap.Uint64("ID", ID),
		zap.Any("meta", c.meta),
	)

	return nil
}

// // // // // // // // // // // // // // // // // // // // // // // // //

// // // // // // // // // // // // // // // // // // // // // // // // //
// DML

// Create stores new data (create compose record)
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

func (svc *service) storeOpPrep(ctx context.Context, mf ModelFilter, capabilities capabilities.Set) (model *Model, cw *ConnectionWrap, err error) {
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
	var (
		log       = svc.logger.Named("models")
		issues    = newIssueHelper()
		auxIssues = newIssueHelper()
	)

	svc.logger.Debug("creating", zap.Int("count", len(models)))

	defer svc.updateIssues(issues)

	// Validate models
	for _, model := range models {
		issues.addModel(model.ResourceID)

		if model.ConnectionID == 0 {
			// Replace model's connection ID with default one when zero
			model.ConnectionID = svc.defConnID
		}

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
			if !svc.sensitivityLevels.isSubset(model.SensitivityLevel, conn.meta.SensitivityLevel) {
				issues.addModelIssue(model.ConnectionID, model.ResourceID, errModelCreateGreaterSensitivityLevel(model.ConnectionID, model.ResourceID, model.SensitivityLevel, conn.meta.SensitivityLevel))
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

		log.Debug("validated",
			zap.Uint64("ID", model.ResourceID),
			zap.String("ident", model.Ident),
			zap.String("label", model.Label),
		)
	}

	// Add models to corresponding connections
	for connection, models := range svc.modelByConnection(models) {
		for _, model := range models {
			mLog := log.With(
				zap.Uint64("connectionID", model.ConnectionID),
				zap.Uint64("ID", model.ResourceID),
				zap.String("ident", model.Ident),
				zap.String("label", model.Label),
			)

			connectionIssues := svc.hasConnectionIssues(model.ConnectionID)
			modelIssues := svc.hasModelIssues(model.ConnectionID, model.ResourceID)

			if !modelIssues && !connectionIssues {
				mLog.Debug("adding model to connection")

				// Add model to connection
				auxIssues, err = svc.registerModelToConnection(ctx, connection, model)
				issues.mergeWith(auxIssues)
				if err != nil {
					return
				}
			} else {
				if connectionIssues {
					mLog.Warn("not adding to connection due to connection issues")
				}
				if modelIssues {
					mLog.Warn("not adding to connection due to model issues")
				}
			}

			// Add model to internal registry
			svc.models[model.ConnectionID] = append(svc.models[model.ConnectionID], model)
		}
	}

	log.Debug("done")

	return
}

// DeleteModel removes support for the model and deletes it from the connection
func (svc *service) DeleteModel(ctx context.Context, models ...*Model) (err error) {

	var (
		log    = svc.logger.Named("models")
		issues = newIssueHelper()
	)

	log.Debug("deleting", zap.Int("count", len(models)))

	defer svc.updateIssues(issues)

	// validation
	skip := make(map[uint64]bool)
	for _, model := range models {
		if model.ConnectionID == 0 {
			// Replace model's connection ID with default one when zero
			model.ConnectionID = svc.defConnID
		}

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
		mLog := log.With(
			zap.Uint64("connectionID", model.ConnectionID),
			zap.Uint64("ID", model.ResourceID),
			zap.String("ident", model.Ident),
			zap.String("label", model.Label),
		)

		if skip[model.ResourceID] {
			mLog.Debug("model does not exist; skipping")
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

		mLog.Debug("deleted")
	}

	return nil
}

func (svc *service) UpdateModel(ctx context.Context, old *Model, new *Model) (err error) {
	if old.ConnectionID == 0 {
		// Replace old model's connection ID with default one when zero
		old.ConnectionID = svc.defConnID
	}

	if new.ConnectionID == 0 {
		// Replace new model's connection ID with default one when zero
		new.ConnectionID = svc.defConnID
	}

	var (
		log = svc.logger.Named("models").With(
			zap.Uint64("connection", new.ConnectionID),
			zap.Uint64("model", new.ResourceID),
			zap.Uint64("ID", new.ResourceID),
			zap.String("ident", new.Ident),
			zap.String("label", new.Label),
		)

		conn *ConnectionWrap

		issues = newIssueHelper().addModel(old.ResourceID)
	)

	log.Debug("updating", zap.Uint64("model", old.ResourceID))

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
			if !svc.sensitivityLevels.isSubset(new.SensitivityLevel, conn.meta.SensitivityLevel) {
				issues.addModelIssue(new.ConnectionID, new.ResourceID, errModelUpdateGreaterSensitivityLevel(new.ConnectionID, new.ResourceID, new.SensitivityLevel, conn.meta.SensitivityLevel))
			}
		}

		// @note attribute check should be done in update model attribute so it's omitted here
	}

	// Update connection
	connectionIssues := svc.hasConnectionIssues(new.ConnectionID)
	modelIssues := svc.hasModelIssues(new.ConnectionID, new.ResourceID)

	if !modelIssues && !connectionIssues {
		log.Debug("updating connection's model")

		err = conn.connection.UpdateModel(ctx, old, new)
		if err != nil {
			issues.addModelIssue(new.ConnectionID, new.ResourceID, err)
		}
	} else {
		if connectionIssues {
			log.Warn("not updating connection's model due to connection issues")
		}
		if modelIssues {
			log.Warn("not updating connection's model due to model issues")
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

	log.Debug("updated")

	return
}

func (svc *service) UpdateModelAttribute(ctx context.Context, model *Model, old, new *Attribute, trans ...TransformationFunction) (err error) {
	svc.logger.Debug("updating model attribute", zap.Uint64("model", model.ResourceID))

	var (
		conn   *ConnectionWrap
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
	if connectionID == 0 {
		connectionID = svc.defConnID
	}

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

func (svc *service) removeConnection(connectionID uint64) {
	delete(svc.connections, connectionID)
}

func (svc *service) addConnection(cw *ConnectionWrap) {
	svc.connections[cw.connectionID] = cw
}

func (svc *service) getConnectionByID(connectionID uint64) (cw *ConnectionWrap) {
	return svc.connections[connectionID]
}

func (svc *service) getConnection(connectionID uint64, cc ...capabilities.Capability) (cw *ConnectionWrap, can capabilities.Set, err error) {
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
func (svc *service) modelByConnection(models ModelSet) (out map[*ConnectionWrap]ModelSet) {
	out = make(map[*ConnectionWrap]ModelSet)

	for _, model := range models {
		c := svc.getConnectionByID(model.ConnectionID)
		out[c] = append(out[c], model)
	}

	return
}

func (svc *service) registerModelToConnection(ctx context.Context, cw *ConnectionWrap, model *Model) (issues *issueHelper, err error) {
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
	if mf.ConnectionID == 0 {
		mf.ConnectionID = svc.defConnID
	}

	if mf.ResourceID > 0 {
		return svc.FindModelByResourceID(mf.ConnectionID, mf.ResourceID)
	}
	return svc.FindModelByResourceIdent(mf.ConnectionID, mf.ResourceType, mf.Resource)
}

func (svc *service) validateNewSensitivityLevels(levels *sensitivityLevelIndex) (err error) {
	err = func() (err error) {
		cIndex := make(map[uint64]*ConnectionWrap)

		// - connections
		for _, _c := range svc.connections {
			c := _c
			cIndex[c.connectionID] = c

			if !levels.includes(c.meta.SensitivityLevel) {
				return fmt.Errorf("connection sensitivity level missing %d", c.meta.SensitivityLevel)
			}
		}

		// - models
		for _, mm := range svc.models {
			for _, m := range mm {
				if !levels.includes(m.SensitivityLevel) {
					return fmt.Errorf("model sensitivity level missing %d", m.SensitivityLevel)
				}
				if !levels.isSubset(m.SensitivityLevel, cIndex[m.ConnectionID].meta.SensitivityLevel) {
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

func wrapError(pfx string, err error) error {
	return fmt.Errorf("%s: %v", pfx, err)
}
