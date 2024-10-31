package dal

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"go.uber.org/zap"
)

type (
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

	FullService interface {
		Drivers() (drivers []Driver)

		ReplaceSensitivityLevel(levels ...SensitivityLevel) (err error)
		RemoveSensitivityLevel(levelIDs ...uint64) (err error)

		GetConnectionByID(connectionID uint64) *ConnectionWrap

		ReplaceConnection(ctx context.Context, cw *ConnectionWrap, isDefault bool) (err error)
		RemoveConnection(ctx context.Context, ID uint64) (err error)

		SearchModels(ctx context.Context) (out ModelSet, err error)
		ReplaceModel(ctx context.Context, currentAlts []*Alteration, model *Model) (newAlts []*Alteration, err error)
		RemoveModel(ctx context.Context, connectionID, ID uint64) (err error)
		FindModelByResourceID(connectionID uint64, resourceID uint64) *Model
		FindModelByResourceIdent(connectionID uint64, resourceType, resourceIdent string) *Model
		FindModelByIdent(connectionID uint64, ident string) *Model

		Create(ctx context.Context, mf ModelRef, operations OperationSet, rr ...ValueGetter) (err error)
		Update(ctx context.Context, mf ModelRef, operations OperationSet, rr ...ValueGetter) (err error)
		Search(ctx context.Context, mf ModelRef, operations OperationSet, f filter.Filter) (iter Iterator, err error)
		Lookup(ctx context.Context, mf ModelRef, operations OperationSet, lookup ValueGetter, dst ValueSetter) (err error)
		Count(ctx context.Context, mf ModelRef, operations OperationSet, f filter.Filter) (uint, error)
		Delete(ctx context.Context, mf ModelRef, operations OperationSet, vv ...ValueGetter) (err error)
		Truncate(ctx context.Context, mf ModelRef, operations OperationSet) (err error)

		Run(ctx context.Context, pp Pipeline) (iter Iterator, err error)
		Dryrun(ctx context.Context, pp Pipeline) (err error)

		ApplyAlteration(ctx context.Context, alts ...*Alteration) (errs []error, err error)
		ReloadModel(ctx context.Context, currentAlts []*Alteration, model *Model) (newAlts []*Alteration, err error)

		SearchConnectionIssues(connectionID uint64) (out []Issue)
		SearchModelIssues(resourceID uint64) (out []Issue)
		SearchResourceIssues(resourceType, resource string) (out []Issue)
	}
)

var (
	gSvc *service

	// wrapper around id.Next() that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)

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

func Initialized() bool {
	return gSvc != nil
}

// Service returns the global initialized DAL service
//
// Function will panic if DAL service is not set (via SetGlobal)
func Service() *service {
	if gSvc == nil {
		panic("DAL global service not initialized: call dal.SetGlobal first")
	}

	return gSvc
}

func SetGlobal(svc *service, err error) {
	if err != nil {
		panic(err)
	}

	gSvc = svc
}

// Purge resets the service to the initial zero state
// @todo will probably need to change but for now this is ok
//
// Primarily used for testing reasons
func (svc *service) Purge(ctx context.Context) {
	nc := map[uint64]*ConnectionWrap{}
	nc[svc.defConnID] = svc.connections[svc.defConnID]

	svc.connections = nc
	svc.models = make(map[uint64]ModelSet)
	svc.sensitivityLevels = SensitivityLevelIndex()
	svc.connectionIssues = make(dalIssueIndex)
	svc.modelIssues = make(dalIssueIndex)
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// meta

// Drivers returns a set of drivers registered to the DAL service
//
// The driver outlines connection params and operations supported by the
// underlying system.
func (svc *service) Drivers() (drivers []Driver) {
	for _, d := range registeredDrivers {
		drivers = append(drivers, d)
	}

	return
}

// MakeSensitivityLevel prepares a new sensitivity level
func MakeSensitivityLevel(ID uint64, level int, handle string) SensitivityLevel {
	return SensitivityLevel{
		ID:     ID,
		Level:  level,
		Handle: handle,
	}
}

// ReplaceSensitivityLevel creates or updates the provided sensitivity levels
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
		log := log.With(logger.Uint64("ID", l.ID), zap.Int("level", l.Level), zap.String("handle", l.Handle))
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

// RemoveSensitivityLevel removes the provided sensitivity levels
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
		log := log.With(logger.Uint64("ID", l.ID))
		if !nx.includes(l.ID) {
			log.Debug("sensitivity level not found")
			continue
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

// InUseSensitivityLevel checks if and where the sensitivity level is being used
func (svc *service) InUseSensitivityLevel(levelID uint64) (usage SensitivityLevelUsage) {
	usage = SensitivityLevelUsage{}

	// - connections
	for _, c := range svc.connections {
		if levelID == c.Config.SensitivityLevelID {
			usage.connections = append(usage.connections, map[string]any{
				// @todo add when needed
				"ident": c.Config.Label,
			})
		}
	}

	// - models
	for _, mm := range svc.models {
		for _, m := range mm {
			if levelID == m.SensitivityLevelID {
				usage.modules = append(usage.modules, map[string]any{
					// @todo add when needed
					"ident": m.Ident,
				})
			}

			for _, attr := range m.Attributes {
				if levelID == attr.SensitivityLevelID {
					usage.fields = append(usage.fields, map[string]any{
						// @todo add when needed
						"ident": attr.Ident,
					})
				}
			}
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //

// // // // // // // // // // // // // // // // // // // // // // // // //
// Connection management

// MakeConnection makes and returns a new connection (wrap)
func MakeConnection(ID uint64, conn Connection, p ConnectionParams, c ConnectionConfig) *ConnectionWrap {
	return &ConnectionWrap{
		ID:         ID,
		connection: conn,

		params: p,
		Config: c,
	}
}

// ReplaceConnection adds new or updates an existing connection
//
// We rely on the user to provide stable connection IDs and
// uses valid relations to these connections in the models.
//
// Is isDefault when adding a default connection. Service will then
// compensate and use proper IDs when models refer to connection with ID=0
func (svc *service) ReplaceConnection(ctx context.Context, conn *ConnectionWrap, isDefault bool) (err error) {
	// @todo lock/unlock
	var (
		ID      = conn.ID
		issues  = newIssueHelper().addConnection(ID)
		oldConn *ConnectionWrap

		log = svc.logger.Named("connection").With(
			logger.Uint64("ID", ID),
			zap.Any("params", conn.params),
			zap.Any("config", conn.Config),
		)
	)

	if isDefault {
		svc.defConnID = ID

		// @note disabling this cause it removes some test-related boilerplate issues.
		// 			 Optimally we'd have it so please figure something out.
		// if svc.defConnID == 0 {
		// 	// default connection not set yet
		// 	log.Debug("setting as default connection")
		// 	svc.defConnID = ID
		// } else if svc.defConnID != ID {
		// 	// default connection set but ID is different.
		// 	// this does not make any sense
		// 	return fmt.Errorf("different ID for default connection detected (old: %d, new: %d)", svc.defConnID, ID)
		// }
	}

	defer svc.updateIssues(issues)

	// Sensitivity level validations
	if !svc.sensitivityLevels.includes(conn.Config.SensitivityLevelID) {
		issues.addConnectionIssue(ID, Issue{
			err: errConnectionCreateMissingSensitivityLevel(ID, conn.Config.SensitivityLevelID),
		})
	}

	if oldConn = svc.GetConnectionByID(ID); oldConn != nil {
		// Connection exists, validate models and sensitivity levels and close and remove connection at the end
		log.Debug("found existing")

		// Check already registered models and their operations
		//
		// Defer the return till the end so we can get a nicer report of what all is wrong
		errored := false
		for _, model := range svc.models[ID] {
			log.Debug("validating model before connection is updated", zap.String("ident", model.Ident))

			// - sensitivity levels
			if !svc.sensitivityLevels.isSubset(model.SensitivityLevelID, conn.Config.SensitivityLevelID) {
				issues.addConnectionIssue(ID, Issue{
					err: fmt.Errorf("cannot update connection %d: new connection sensitivity level does not support model %d", ID, model.ResourceID),
				})
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
				issues.addConnectionIssue(ID, Issue{
					err: err,
				})
				return nil
			}

			log.Debug("disconnected")
		}

		svc.removeConnection(ID)
	}

	if conn.connection == nil {
		conn.connection, err = connect(ctx, svc.logger, svc.inDev, conn.params)
		if err != nil {
			log.Warn("could not connect", zap.Error(err))
			issues.addConnectionIssue(ID, Issue{
				err: err,
			})
		} else {
			log.Debug("connected")
		}
	} else {
		log.Debug("using preexisting connection")
	}

	svc.addConnection(conn)
	log.Debug("added")
	return nil
}

// RemoveConnection removes the given connection from the DAL
func (svc *service) RemoveConnection(ctx context.Context, ID uint64) (err error) {
	var (
		issues = newIssueHelper().addConnection(ID)
	)

	c := svc.GetConnectionByID(ID)
	if c == nil {
		return errConnectionDeleteNotFound(ID)
	}

	// Potential cleanups
	if cc, ok := c.connection.(ConnectionCloser); ok {
		if err := cc.Close(ctx); err != nil {
			svc.logger.Error(errConnectionDeleteCloserFailed(c.ID, err).Error())
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
		logger.Uint64("ID", ID),
		zap.Any("config", c.Config),
	)

	return nil
}

// // // // // // // // // // // // // // // // // // // // // // // // //

// // // // // // // // // // // // // // // // // // // // // // // // //
// DML

// Create stores new data (create data entry)
func (svc *service) Create(ctx context.Context, mf ModelRef, operations OperationSet, rr ...ValueGetter) (err error) {
	if err = svc.canOpData(mf); err != nil {
		return fmt.Errorf("cannot create data entry: %w", err)
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, operations)
	if err != nil {
		return fmt.Errorf("cannot create data entry: %w", err)
	}

	return cw.connection.Create(ctx, model, rr...)
}

func (svc *service) Update(ctx context.Context, mf ModelRef, operations OperationSet, rr ...ValueGetter) (err error) {
	if err = svc.canOpData(mf); err != nil {
		return fmt.Errorf("cannot update data entry: %w", err)
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, operations)
	if err != nil {
		return fmt.Errorf("cannot update data entry: %w", err)
	}

	for _, r := range rr {
		if err = cw.connection.Update(ctx, model, r); err != nil {
			return fmt.Errorf("cannot update data entry: %w", err)
		}
	}

	return
}

func (svc *service) FindModel(mr ModelRef) *Model {
	return svc.getModelByRef(mr)
}

func (svc *service) Search(ctx context.Context, mf ModelRef, operations OperationSet, f filter.Filter) (iter Iterator, err error) {
	if err = svc.canOpData(mf); err != nil {
		err = fmt.Errorf("cannot search data entry: %w", err)
		return
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, operations)
	if err != nil {
		err = fmt.Errorf("cannot search data entry: %w", err)
		return
	}

	return cw.connection.Search(ctx, model, f)
}

func (svc *service) Count(ctx context.Context, mf ModelRef, operations OperationSet, f filter.Filter) (cnt uint, err error) {
	if err = svc.canOpData(mf); err != nil {
		err = fmt.Errorf("cannot count data entry: %w", err)
		return
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, operations)
	if err != nil {
		err = fmt.Errorf("cannot count data entry: %w", err)
		return
	}

	return cw.connection.Count(ctx, model, f)
}

// Run returns an iterator based on the provided Pipeline
// @todo consider moving the Search method to utilize this also
func (svc *service) Run(ctx context.Context, pp Pipeline) (iter Iterator, err error) {
	pp, err = svc.pipelinePrerun(ctx, pp)
	if err != nil {
		return
	}

	err = svc.analyzePipeline(ctx, pp)
	if err != nil {
		return
	}

	pp, err = svc.optimizePipeline(ctx, pp)
	if err != nil {
		return
	}

	return svc.run(ctx, pp.root(), false)
}

// Dryrun processes the Pipeline but doesn't return the iterator
// @todo consider reworking this; I'm not the biggest fan of this one
//
// The method is primarily used by system reports to obtain some metadata
func (svc *service) Dryrun(ctx context.Context, pp Pipeline) (err error) {
	// @note we don't need to do any optimization or analisis here
	pp, err = svc.pipelinePrerun(ctx, pp)
	if err != nil {
		return
	}

	_, err = svc.run(ctx, pp.root(), true)
	return
}

// run is the recursive counterpart to run
func (svc *service) run(ctx context.Context, s PipelineStep, dry bool) (it Iterator, err error) {
	switch s := s.(type) {
	case *Datasource:
		err = s.init(ctx)
		if err != nil {
			return
		}
		if dry {
			return nil, nil
		}
		return s.exec(ctx)

	case *Aggregate:
		it, err = svc.run(ctx, s.rel, dry)
		if err != nil {
			return
		}

		if dry {
			return nil, s.dryrun(ctx)
		}
		return s.iterator(ctx, it)

	case *Join:
		var left Iterator
		var right Iterator
		left, err = svc.run(ctx, s.relLeft, dry)
		if err != nil {
			return
		}
		right, err = svc.run(ctx, s.relRight, dry)
		if err != nil {
			return
		}

		if dry {
			return nil, s.dryrun(ctx)
		}
		return s.iterator(ctx, left, right)

	case *Link:
		var left Iterator
		var right Iterator
		left, err = svc.run(ctx, s.relLeft, dry)
		if err != nil {
			return
		}
		right, err = svc.run(ctx, s.relRight, dry)
		if err != nil {
			return
		}

		if dry {
			return nil, s.dryrun(ctx)
		}
		return s.iterator(ctx, left, right)
	}

	return nil, fmt.Errorf("unsupported step")
}

func collectAttributes(s PipelineStep) (out []AttributeMapping) {
	switch s := s.(type) {
	case *Datasource:
		return s.OutAttributes

	case *Aggregate:
		out = make([]AttributeMapping, 0, 24)
		for _, a := range s.Group {
			out = append(out, a.toSimpleAttr())
		}
		for _, a := range s.OutAttributes {
			out = append(out, a.toSimpleAttr())
		}
		return

	case *Join:
		return s.OutAttributes

	case *Link:
		panic("impossible state; link can not be nested")
	}

	return
}

func (svc *service) Lookup(ctx context.Context, mf ModelRef, operations OperationSet, lookup ValueGetter, dst ValueSetter) (err error) {
	if err = svc.canOpData(mf); err != nil {
		return fmt.Errorf("cannot lookup data entry: %w", err)
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, operations)
	if err != nil {
		return fmt.Errorf("cannot lookup data entry: %w", err)
	}
	return cw.connection.Lookup(ctx, model, lookup, dst)
}

func (svc *service) Delete(ctx context.Context, mf ModelRef, operations OperationSet, vv ...ValueGetter) (err error) {
	if err = svc.canOpData(mf); err != nil {
		return fmt.Errorf("cannot delete data entry: %w", err)
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, operations)
	if err != nil {
		return fmt.Errorf("cannot delete data entry: %w", err)
	}

	for _, v := range vv {
		if err = cw.connection.Delete(ctx, model, v); err != nil {
			return fmt.Errorf("cannot delete data entry: %w", err)
		}
	}
	return
}

func (svc *service) Truncate(ctx context.Context, mf ModelRef, operations OperationSet) (err error) {
	if err = svc.canOpData(mf); err != nil {
		return fmt.Errorf("cannot truncate data entry: %w", err)
	}

	model, cw, err := svc.storeOpPrep(ctx, mf, operations)
	if err != nil {
		return fmt.Errorf("cannot truncate data entry: %w", err)
	}

	return cw.connection.Truncate(ctx, model)
}

func (svc *service) storeOpPrep(ctx context.Context, mf ModelRef, operations OperationSet) (model *Model, cw *ConnectionWrap, err error) {
	model = svc.getModelByRef(mf)
	if model == nil {
		err = errModelNotFound(mf.ResourceID)
		return
	}

	cw, _, err = svc.getConnection(model.ConnectionID, operations...)
	if err != nil {
		return
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //

// // // // // // // // // // // // // // // // // // // // // // // // //
// DDL

// SearchModels returns a list of modules registered under DAL
//
// Primarily used for testing (for data truncate).
func (svc *service) SearchModels(ctx context.Context) (out ModelSet, err error) {
	out = make(ModelSet, 0, 100)
	for _, models := range svc.models {
		out = append(out, models...)
	}
	return
}

// ReplaceModel adds new or updates an existing model
//
// We rely on the user to provide stable and valid model definitions.
func (svc *service) ReplaceModel(ctx context.Context, currentAlts []*Alteration, model *Model) (newAlts []*Alteration, err error) {
	var (
		ID         = model.ResourceID
		connection = svc.GetConnectionByID(model.ConnectionID)
		oldModel   *Model
		issues     = newIssueHelper().addModel(ID)
		upd        bool

		modelIssues      bool
		connectionIssues bool

		log = svc.logger.Named("models").With(
			logger.Uint64("ID", ID),
			zap.String("ident", model.Ident),
			zap.Any("label", model.Label),
		)
	)

	defer svc.updateIssues(issues)

	if model.ConnectionID == 0 {
		model.ConnectionID = svc.defConnID
	}

	// Check if we're creating or updating the model
	if oldModel = svc.FindModel(model.ToFilter()); oldModel != nil {
		log.Debug("found existing")

		if oldModel.ConnectionID != model.ConnectionID {
			log.Warn("changed model connection, existing data potentially unavailable")
		}

		upd = true
	}

	// @todo consider adding some logging to validators
	svc.validateModel(issues, connection, model, oldModel)
	svc.validateAttributes(issues, model, model.Attributes...)

	connectionIssues = issues.hasConnectionIssues()
	if connectionIssues {
		log.Warn(
			"not adding to store due to connection issues",
			zap.Any("issues", svc.SearchConnectionIssues(model.ConnectionID)),
		)
	}

	modelIssues = issues.hasModelIssues()
	if modelIssues {
		log.Warn(
			"not adding to store due to model issues",
			zap.Any("issues", svc.SearchModelIssues(model.ResourceID)),
		)
	}

	// Remove the old model from the registry
	if oldModel != nil {
		svc.removeModelFromRegistry(oldModel)
		log.Debug(
			"removed old model from registry",
			logger.Uint64("connectionID", model.ConnectionID),
		)
	}

	// Add to registry
	// Models should be added to the registry regardless of issues
	svc.addModelToRegistry(model, upd)
	log.Debug(
		"added to registry",
		logger.Uint64("connectionID", model.ConnectionID),
	)

	// Can't continue with connection alterations if there are issues already
	if connectionIssues || modelIssues {
		return
	}

	newAlts, batchID, err := svc.getSchemaAlterations(ctx, connection, currentAlts, oldModel, model)
	if err != nil {
		return
	}

	if len(newAlts) > 0 {
		svc.setAlterationsModelIssue(issues, batchID, connection, model, newAlts)
		log.Info("not adding to store: alterations required", zap.Error(err))
		return
	}

	if upd {
		err = connection.connection.UpdateModel(ctx, oldModel, model)
	} else {
		err = connection.connection.CreateModel(ctx, model)
	}

	if err != nil {
		log.Error("failed with errors", zap.Error(err))
	}

	return
}

// ApplyAlteration updates the underlying schema with the requested changes
func (svc *service) ApplyAlteration(ctx context.Context, alts ...*Alteration) (errs []error, err error) {
	if len(alts) == 0 {
		return
	}

	altsIndexed := make(map[uint64]map[string]map[string][]*Alteration)
	for _, a := range alts {
		if _, ok := altsIndexed[a.ConnectionID]; !ok {
			altsIndexed[a.ConnectionID] = make(map[string]map[string][]*Alteration)
		}

		if _, ok := altsIndexed[a.ConnectionID][a.ResourceType]; !ok {
			altsIndexed[a.ConnectionID][a.ResourceType] = make(map[string][]*Alteration)
		}

		altsIndexed[a.ConnectionID][a.ResourceType][a.Resource] = append(altsIndexed[a.ConnectionID][a.ResourceType][a.Resource], a)
	}

	var auxErrs []error
	for _, byCon := range altsIndexed {
		for _, byRt := range byCon {
			for _, byR := range byRt {
				auxErrs, err = svc.applyAlteration(ctx, byR...)
				if err != nil {
					return
				}

				errs = append(errs, auxErrs...)
			}
		}
	}

	return
}

func (svc *service) applyAlteration(ctx context.Context, alts ...*Alteration) (errs []error, err error) {
	if len(alts) == 0 {
		return
	}

	var (
		c = svc.GetConnectionByID(0)

		connectionID = alts[0].ConnectionID
		resource     = alts[0].Resource
		resourceType = alts[0].ResourceType

		issues = newIssueHelper()
	)

	defer svc.updateIssues(issues)

	for _, alt := range alts {
		if alt.ConnectionID != connectionID {
			return nil, fmt.Errorf("alterations must be for the same connection")
		}

		if alt.ResourceType != resourceType {
			return nil, fmt.Errorf("alterations must be for the same resource type")
		}

		if alt.Resource != resource {
			return nil, fmt.Errorf("alterations must be for the same resource")
		}
	}

	connection := svc.GetConnectionByID(connectionID)
	if connection == nil {
		return nil, fmt.Errorf("connection not found")
	}

	model := svc.getModelByRef(ModelRef{Resource: resource, ResourceType: resourceType, ConnectionID: connectionID})
	if model == nil {
		return nil, fmt.Errorf("model not found")
	}

	if model.ConnectionID == c.ID && model.Ident == "compose_record" {
		err = fmt.Errorf("cannot apply alterations for default schema")
		return
	}

	issues = issues.addModel(model.ResourceID)

	// @todo consider adding some logging to validators
	svc.validateModel(issues, connection, model, model)
	svc.validateAttributes(issues, model, model.Attributes...)

	return connection.connection.ApplyAlteration(ctx, model, alts...), nil
}

func (svc *service) ReloadModel(ctx context.Context, currentAlts []*Alteration, model *Model) (newAlts []*Alteration, err error) {
	var (
		issues = newIssueHelper()

		log = svc.logger.Named("models").With(
			logger.Uint64("ID", model.ResourceID),
			zap.String("ident", model.Ident),
			zap.Any("label", model.Label),
		)
	)
	defer svc.updateIssues(issues)

	connection := svc.GetConnectionByID(model.ConnectionID)
	if connection == nil {
		err = fmt.Errorf("connection not found")
		return
	}

	issues = issues.addModel(model.ResourceID)

	// @todo consider adding some logging to validators
	svc.validateModel(issues, connection, model, model)
	svc.validateAttributes(issues, model, model.Attributes...)

	// If there are any issues at this stage, there is nothing for us to do
	if issues.hasConnectionIssues() || issues.hasModelIssues() {
		log.Warn(
			"not reloading due to issues",
			zap.Any("connection issues", svc.SearchConnectionIssues(model.ConnectionID)),
			zap.Any("model issues", svc.SearchModelIssues(model.ResourceID)),
		)
		return
	}

	// Re-evaluate schema alterations
	// oldModel nil will force it to re-check the entire thing
	newAlts, batchID, err := svc.getSchemaAlterations(ctx, connection, currentAlts, nil, model)
	if err != nil {
		return
	}

	if len(newAlts) > 0 {
		svc.setAlterationsModelIssue(issues, batchID, connection, model, newAlts)
		log.Info("not adding to store: alterations required", zap.Error(err))
		return
	}

	err = connection.connection.UpdateModel(ctx, model, model)
	if err != nil {
		log.Error("failed with errors", zap.Error(err))
	}

	return
}

// RemoveModel removes the given model from DAL
//
// @todo potentially add more interaction with the connection as in letting it know a model was removed.
func (svc *service) RemoveModel(ctx context.Context, connectionID, ID uint64) (err error) {
	var (
		old *Model

		log = svc.logger.Named("models").With(
			logger.Uint64("connectionID", connectionID),
			logger.Uint64("ID", ID),
		)
		issues = newIssueHelper().addModel(ID)
	)

	log.Debug("deleting")

	if connectionID == 0 {
		connectionID = svc.defConnID
	}

	defer svc.updateIssues(issues)

	// Check we have something to remove
	if old = svc.FindModelByResourceID(connectionID, ID); old == nil {
		return
	}

	// Validate no leftover references
	// @todo we can probably expand on this quitea bit
	// for _, registered := range svc.models {
	// 	refs := registered.FilterByReferenced(model)
	// 	if len(refs) > 0 {
	// 		return fmt.Errorf("cannot remove model %s: referenced by other models", model.Resource)
	// 	}
	// }

	// @todo should the underlying store be notified about this?
	// how should this be handled; a straight up delete doesn't sound sane to me
	// anymore

	svc.removeModelFromRegistry(old)

	log.Debug("removed")
	return nil
}

func (svc *service) validateModel(issues *issueHelper, c *ConnectionWrap, model, oldModel *Model) {
	// Connection ok?
	if svc.hasConnectionIssues(model.ConnectionID) {
		issues.addModelIssue(model.ResourceID, Issue{
			err: errModelCreateProblematicConnection(model.ConnectionID, model.ResourceID),
		})
	}
	// Connection exists?
	if c == nil {
		issues.addModelIssue(model.ResourceID, Issue{
			err: errModelCreateMissingConnection(model.ConnectionID, model.ResourceID),
		})
	}

	// If ident changed, check for duplicate
	if oldModel != nil && oldModel.Ident != model.Ident {
		if tmp := svc.FindModelByIdent(model.ConnectionID, model.Ident); tmp != nil {
			issues.addModelIssue(oldModel.ResourceID, Issue{
				err: errModelUpdateDuplicate(model.ConnectionID, model.ResourceID),
			})
		}
	}

	// Sensitivity level ok and valid?
	if !svc.sensitivityLevels.includes(model.SensitivityLevelID) {
		issues.addModelIssue(model.ResourceID, Issue{
			err: errModelCreateMissingSensitivityLevel(model.ConnectionID, model.ResourceID, model.SensitivityLevelID),
		})
	} else {
		// Only check if it is present
		if !svc.sensitivityLevels.isSubset(model.SensitivityLevelID, c.Config.SensitivityLevelID) {
			issues.addModelIssue(model.ResourceID, Issue{
				err: errModelCreateGreaterSensitivityLevel(model.ConnectionID, model.ResourceID, model.SensitivityLevelID, c.Config.SensitivityLevelID),
			})
		}
	}

	if c != nil && !checkIdent(model.Ident, c.Config.ModelIdentCheck...) {
		issues.addModelIssue(model.ResourceID, Issue{
			err: errModelCreateInvalidIdent(model.ConnectionID, model.ResourceID, model.Ident),
		})
	}
}

func (svc *service) validateAttributes(issues *issueHelper, model *Model, attr ...*Attribute) {
	for _, a := range attr {
		if !svc.sensitivityLevels.includes(a.SensitivityLevelID) {
			issues.addModelIssue(model.ResourceID, Issue{
				err: errModelCreateMissingAttributeSensitivityLevel(model.ConnectionID, model.ResourceID, a.SensitivityLevelID),
			})
		} else {
			if !svc.sensitivityLevels.isSubset(a.SensitivityLevelID, model.SensitivityLevelID) {
				issues.addModelIssue(model.ResourceID, Issue{
					err: errModelCreateGreaterAttributeSensitivityLevel(model.ConnectionID, model.ResourceID, a.SensitivityLevelID, model.SensitivityLevelID),
				})
			}
		}
	}
}

func (svc *service) addModelToRegistry(model *Model, upd bool) {
	if !upd {
		svc.models[model.ConnectionID] = append(svc.models[model.ConnectionID], model)
		return
	}

	ok := false
	for i, old := range svc.models[model.ConnectionID] {
		if old.ResourceID == model.ResourceID {
			svc.models[model.ConnectionID][i] = model
			ok = true
			break
		}
	}
	if !ok {
		svc.models[model.ConnectionID] = append(svc.models[model.ConnectionID], model)
	}
}

func (svc *service) removeModelFromRegistry(model *Model) {
	oldModels := svc.models[model.ConnectionID]
	svc.models[model.ConnectionID] = make(ModelSet, 0, len(oldModels))
	for _, o := range oldModels {
		if o.Resource == model.Resource {
			continue
		}

		svc.models[model.ConnectionID] = append(svc.models[model.ConnectionID], o)
	}
}

// FindModelByRefs returns the model with all of the given refs matching
//
// @note refs are primarily used for DAL pipelines where steps can reference models
//
//	by handles and slugs such as module and namespace.
func (svc *service) FindModelByRefs(connectionID uint64, refs map[string]any) *Model {
	if connectionID == 0 {
		connectionID = svc.defConnID
	}
	return svc.models[connectionID].FindByRefs(refs)
}

func (svc *service) FindModelByResourceID(connectionID uint64, resourceID uint64) *Model {
	if connectionID == 0 {
		connectionID = svc.defConnID
	}

	return svc.models[connectionID].FindByResourceID(resourceID)
}

func (svc *service) FindModelByResourceIdent(connectionID uint64, resourceType, resourceIdent string) *Model {
	if connectionID == 0 {
		connectionID = svc.defConnID
	}

	return svc.models[connectionID].FindByResourceIdent(resourceType, resourceIdent)
}

func (svc *service) FindModelByRef(ref ModelRef) *Model {
	connectionID := ref.ConnectionID
	if connectionID == 0 {
		connectionID = svc.defConnID
	}

	if ref.Refs != nil {
		return svc.FindModelByRefs(connectionID, ref.Refs)
	}

	if ref.ResourceID > 0 {
		return svc.models[connectionID].FindByResourceID(ref.ResourceID)
	}

	return svc.models[connectionID].FindByResourceIdent(ref.ResourceType, ref.Resource)
}

func (svc *service) FindModelByIdent(connectionID uint64, ident string) *Model {
	if connectionID == 0 {
		connectionID = svc.defConnID
	}

	return svc.models[connectionID].FindByIdent(ident)
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utilities

func (svc *service) removeConnection(connectionID uint64) {
	delete(svc.connections, connectionID)
}

func (svc *service) addConnection(cw *ConnectionWrap) {
	svc.connections[cw.ID] = cw
}

func (svc *service) GetConnectionByID(connectionID uint64) (cw *ConnectionWrap) {
	if connectionID == 0 {
		connectionID = svc.defConnID
	}

	return svc.connections[connectionID]
}

func (svc *service) getConnection(connectionID uint64, oo ...Operation) (cw *ConnectionWrap, can OperationSet, err error) {
	err = func() error {
		// get the requested connection
		cw = svc.GetConnectionByID(connectionID)
		if cw == nil {
			return errConnectionNotFound(connectionID)
		}

		// check if connection supports requested operations
		if !cw.connection.Can(oo...) {
			return fmt.Errorf("connection %d does not support requested operations %v", connectionID, OperationSet(oo).Diff(cw.connection.Operations()))
		}
		can = cw.connection.Operations()
		return nil
	}()

	if err != nil {
		err = fmt.Errorf("could not connect to %d: %v", connectionID, err)
		return
	}

	return
}

func (svc *service) getModelByRef(mr ModelRef) *Model {
	if mr.ConnectionID == 0 {
		mr.ConnectionID = svc.defConnID
	}

	if mr.Refs != nil {
		return svc.FindModelByRefs(mr.ConnectionID, mr.Refs)
	} else if mr.ResourceID > 0 {
		return svc.FindModelByResourceID(mr.ConnectionID, mr.ResourceID)
	}
	return svc.FindModelByResourceIdent(mr.ConnectionID, mr.ResourceType, mr.Resource)
}

func (svc *service) validateNewSensitivityLevels(levels *sensitivityLevelIndex) (err error) {
	err = func() (err error) {
		cIndex := make(map[uint64]*ConnectionWrap)

		// - connections
		for _, _c := range svc.connections {
			c := _c
			cIndex[c.ID] = c

			if !levels.includes(c.Config.SensitivityLevelID) {
				return fmt.Errorf("connection sensitivity level missing %d", c.Config.SensitivityLevelID)
			}
		}

		// - models
		for _, mm := range svc.models {
			for _, m := range mm {
				if !levels.includes(m.SensitivityLevelID) {
					return fmt.Errorf("model sensitivity level missing %d", m.SensitivityLevelID)
				}
				if !levels.isSubset(m.SensitivityLevelID, cIndex[m.ConnectionID].Config.SensitivityLevelID) {
					return fmt.Errorf("model sensitivity level missing %d", m.SensitivityLevelID)
				}

				for _, attr := range m.Attributes {
					if !levels.includes(attr.SensitivityLevelID) {
						return fmt.Errorf("attribute sensitivity level missing %d", attr.SensitivityLevelID)
					}
					if !levels.isSubset(attr.SensitivityLevelID, m.SensitivityLevelID) {
						return fmt.Errorf("attribute sensitivity level %d greater then model sensitivity level %d", attr.SensitivityLevelID, m.SensitivityLevelID)
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

// pipelinePrerun performs the common operations for both the Run and Dryrun
func (svc *service) pipelinePrerun(ctx context.Context, pp Pipeline) (_ Pipeline, err error) {
	err = pp.LinkSteps()
	if err != nil {
		return
	}
	err = svc.bindDatasourceConnections(ctx, pp)
	if err != nil {
		return
	}

	return pp, nil
}

// optimizePipeline runs optimization over the given Pipeline
func (svc *service) optimizePipeline(ctx context.Context, pp Pipeline) (_ Pipeline, err error) {
	pp, err = svc.optimizePipelineStructure(ctx, pp)
	if err != nil {
		return
	}

	// @todo add step-based optimization such as filter pushdown; omitting for now
	//       since we don't have any of it in place yet

	return pp, nil
}

// optimizePipelineStructure performs general pipeline structure optimizations
// such as restructuring and clobbering steps onto the datasource layer
func (svc *service) optimizePipelineStructure(ctx context.Context, pp Pipeline) (_ Pipeline, err error) {
	for _, o := range pipelineOptimizers {
		pp, err = o(pp)
		if err != nil {
			return nil, err
		}
	}

	return pp, nil
}

// bindDatasourceConnections special handles datasource pipeline steps and binds
// a DAL connection to them
func (svc *service) bindDatasourceConnections(ctx context.Context, pp Pipeline) error {
	for _, p := range pp {
		ds, ok := p.(*Datasource)
		if !ok {
			continue
		}

		ds.model = svc.FindModelByRef(ds.ModelRef)
		if ds.model == nil {
			return fmt.Errorf("model %v does not exist", ds.ModelRef)
		}

		ds.connection = svc.GetConnectionByID(ds.model.ConnectionID)
		if ds.connection == nil {
			return fmt.Errorf("connection %d does not exist", ds.model.ConnectionID)
		}
	}

	return nil
}

// analyzePipeline runs analysis over each step in the pipeline
//
// Step analysis hints to the optimizers as to how expensive specific operations
// are and the general dataset size involved.
func (svc *service) analyzePipeline(ctx context.Context, pp Pipeline) (err error) {
	for _, p := range pp {
		err = p.Analyze(ctx)
		if err != nil {
			return
		}
	}

	return
}

func (svc *service) mergeAlterations(base, added AlterationSet) (out AlterationSet) {
	var (
		batchID uint64
	)

	if len(base) > 0 {
		batchID = base[0].BatchID
	}

	for _, a := range base {
		a.BatchID = batchID
	}

	return base.Merge(added)
}

func (svc *service) getSchemaAlterations(ctx context.Context, connection *ConnectionWrap, currentAlts []*Alteration, oldModel, model *Model) (newAlts []*Alteration, batchID uint64, err error) {
	// - use the diff between the two models as a starting point to see what we should do to support the change
	df := oldModel.Diff(model)
	newAlts = df.Alterations()
	batchID = nextID()
	for _, a := range newAlts {
		a.BatchID = batchID
		a.Resource = model.Resource
		a.ResourceType = model.ResourceType
		a.ConnectionID = model.ConnectionID
	}

	// - merge stale and new alteration
	// @todo for now we're just using the newly calculated alterations as merging with
	//       existing ones is not that trivial and doesn't add much value.
	// @note this merging assumes the two sets are already ok, valid, and without any
	//       duplications.
	newAlts = svc.mergeAlterations(currentAlts, newAlts)
	// - run the alterations against the database to take the schema into consideration
	newAlts, err = connection.connection.AssertSchemaAlterations(ctx, model, newAlts...)
	if err != nil {
		return
	}

	// - set all of the alterations to the same batch ID
	for _, a := range newAlts {
		a.BatchID = batchID
	}

	return
}

func (svc *service) setAlterationsModelIssue(issues *issueHelper, batchID uint64, connection *ConnectionWrap, model *Model, alts []*Alteration) {
	if len(alts) > 0 {
		issues.addModelIssue(model.ResourceID, Issue{
			err: errModelRequiresAlteration(connection.ID, model.ResourceID, batchID),
			Meta: map[string]any{
				"batchID": strconv.FormatUint(batchID, 10),
			},
		})
	}
}
