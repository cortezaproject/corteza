package service

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/handle"
	"github.com/cortezaproject/corteza/server/pkg/label"
	"github.com/cortezaproject/corteza/server/store"
	"go.uber.org/zap"
)

type (
	triggerDefinition struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        triggerDefinitionAccessController
		log       *zap.Logger
	}

	triggerDefinitionAccessController interface {
		CanCreateTriggerDefinition(context.Context) bool
		CanSearchTriggerDefinitions(context.Context) bool
		CanReadTriggerDefinition(context.Context, *types.TriggerDefinition) bool
		CanUpdateTriggerDefinition(context.Context, *types.TriggerDefinition) bool
		CanDeleteTriggerDefinition(context.Context, *types.TriggerDefinition) bool
	}

	TriggerDefinitionService interface {
		Search(ctx context.Context, filter types.TriggerDefinitionFilter) (types.TriggerDefinitionSet, types.TriggerDefinitionFilter, error)
		LookupByID(ctx context.Context, id uint64) (*types.TriggerDefinition, error)
		LookupByHandle(ctx context.Context, handle string) (*types.TriggerDefinition, error)
		Create(ctx context.Context, new *types.TriggerDefinition) (*types.TriggerDefinition, error)
		Update(ctx context.Context, upd *types.TriggerDefinition) (*types.TriggerDefinition, error)
		DeleteByID(ctx context.Context, id uint64) error
		UndeleteByID(ctx context.Context, id uint64) error
	}
)

func TriggerDefinition(log *zap.Logger) *triggerDefinition {
	return &triggerDefinition{
		log:       log,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

func (svc triggerDefinition) Search(ctx context.Context, filter types.TriggerDefinitionFilter) (set types.TriggerDefinitionSet, f types.TriggerDefinitionFilter, err error) {
	var (
		aProps = &triggerDefinitionActionProps{filter: &filter}
	)

	filter.Check = func(res *types.TriggerDefinition) (bool, error) {
		if !svc.ac.CanReadTriggerDefinition(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchTriggerDefinitions(ctx) {
			return TriggerDefinitionErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.TriggerDefinition{}.LabelResourceKind(),
				filter.Labels,
			)

			if err != nil {
				return err
			}

			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if set, f, err = store.SearchAutomationTriggerDefinitions(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledTriggerDefinitions(set)...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, TriggerDefinitionActionSearch, err)
}

func (svc triggerDefinition) LookupByID(ctx context.Context, ID uint64) (def *types.TriggerDefinition, err error) {
	var (
		aProps = &triggerDefinitionActionProps{triggerDefinition: &types.TriggerDefinition{ID: ID}}
	)

	err = func() error {
		if def, err = loadTriggerDefinition(ctx, svc.store, ID); err != nil {
			return err
		}

		aProps.setTriggerDefinition(def)

		if !svc.ac.CanReadTriggerDefinition(ctx, def) {
			return TriggerDefinitionErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, def); err != nil {
			return err
		}

		return nil
	}()

	return def, svc.recordAction(ctx, aProps, TriggerDefinitionActionLookup, err)
}

func (svc triggerDefinition) LookupByHandle(ctx context.Context, h string) (def *types.TriggerDefinition, err error) {
	var (
		aProps = &triggerDefinitionActionProps{triggerDefinition: &types.TriggerDefinition{Handle: h}}
	)

	err = func() error {
		def, err = store.LookupAutomationTriggerDefinitionByHandle(ctx, svc.store, h)
		if err != nil {
			return err
		}

		aProps.setTriggerDefinition(def)

		if !svc.ac.CanReadTriggerDefinition(ctx, def) {
			return TriggerDefinitionErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, def); err != nil {
			return err
		}

		return nil
	}()

	return def, svc.recordAction(ctx, aProps, TriggerDefinitionActionLookup, err)
}

func (svc triggerDefinition) Create(ctx context.Context, new *types.TriggerDefinition) (def *types.TriggerDefinition, err error) {
	var (
		aProps = &triggerDefinitionActionProps{triggerDefinition: new}
		cUser  = internalAuth.GetIdentityFromContext(ctx).Identity()
	)

	err = func() (err error) {
		if !svc.ac.CanCreateTriggerDefinition(ctx) {
			return TriggerDefinitionErrNotAllowedToCreate()
		}

		if !handle.IsValid(new.Handle) {
			return TriggerDefinitionErrInvalidHandle()
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = cUser
		new.OwnedBy = cUser

		if err = store.CreateAutomationTriggerDefinition(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		def = new
		aProps.setNew(new)

		svc.reRegisterConstructs(def)

		return nil
	}()

	return def, svc.recordAction(ctx, aProps, TriggerDefinitionActionCreate, err)
}

func (svc triggerDefinition) Update(ctx context.Context, upd *types.TriggerDefinition) (def *types.TriggerDefinition, err error) {
	var (
		aProps = &triggerDefinitionActionProps{update: upd}
		cUser  = internalAuth.GetIdentityFromContext(ctx).Identity()
	)

	err = func() (err error) {
		if !handle.IsValid(upd.Handle) {
			return TriggerDefinitionErrInvalidHandle()
		}

		if def, err = loadTriggerDefinition(ctx, svc.store, upd.ID); err != nil {
			return
		}

		aProps.setTriggerDefinition(def)

		if !svc.ac.CanUpdateTriggerDefinition(ctx, def) {
			return TriggerDefinitionErrNotAllowedToUpdate()
		}

		if isStale(upd.UpdatedAt, def.UpdatedAt, def.CreatedAt) {
			return TriggerDefinitionErrStaleData()
		}

		if err = svc.uniqueCheck(ctx, upd); err != nil {
			return err
		}

		def.Handle = upd.Handle
		def.SkipEventBus = upd.SkipEventBus
		def.InputSchema = upd.InputSchema
		def.OutputSchema = upd.OutputSchema
		def.Meta = upd.Meta
		def.UpdatedAt = now()
		def.UpdatedBy = cUser

		if err = store.UpdateAutomationTriggerDefinition(ctx, svc.store, def); err != nil {
			return
		}

		if label.Changed(def.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}
			def.Labels = upd.Labels
		}

		svc.reRegisterConstructs(def)
		return nil
	}()

	return def, svc.recordAction(ctx, aProps, TriggerDefinitionActionUpdate, err)
}

func (svc triggerDefinition) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &triggerDefinitionActionProps{triggerDefinition: &types.TriggerDefinition{ID: ID}}
		def    *types.TriggerDefinition
	)

	err = func() (err error) {
		if def, err = loadTriggerDefinition(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setTriggerDefinition(def)

		if !svc.ac.CanDeleteTriggerDefinition(ctx, def) {
			return TriggerDefinitionErrNotAllowedToDelete()
		}

		def.DeletedAt = now()
		def.DeletedBy = internalAuth.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateAutomationTriggerDefinition(ctx, svc.store, def); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, aProps, TriggerDefinitionActionDelete, err)
}

func (svc triggerDefinition) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &triggerDefinitionActionProps{triggerDefinition: &types.TriggerDefinition{ID: ID}}
		def    *types.TriggerDefinition
	)

	err = func() (err error) {
		if def, err = loadTriggerDefinition(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setTriggerDefinition(def)

		if !svc.ac.CanDeleteTriggerDefinition(ctx, def) {
			return TriggerDefinitionErrNotAllowedToUndelete()
		}

		def.DeletedAt = nil
		def.DeletedBy = 0

		if err = store.UpdateAutomationTriggerDefinition(ctx, svc.store, def); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, aProps, TriggerDefinitionActionUndelete, err)
}

func (svc triggerDefinition) uniqueCheck(ctx context.Context, res *types.TriggerDefinition) (err error) {
	if res.Handle != "" {
		if e, _ := store.LookupAutomationTriggerDefinitionByHandle(ctx, svc.store, res.Handle); e != nil && e.ID != res.ID {
			return TriggerDefinitionErrHandleNotUnique()
		}
	}
	return nil
}

func (svc triggerDefinition) reRegisterConstructs(def *types.TriggerDefinition) {
	ConstructLibrary().AddTriggers(triggerDefToConstruct(def))
}

func triggerDefToConstruct(def *types.TriggerDefinition) types.ConstructTrigger {
	ref := def.Handle
	if ref == "" {
		ref = fmt.Sprintf("%d", def.ID)
	}

	props := make([]types.ConstructTriggerProperty, len(def.InputSchema))
	for i, p := range def.InputSchema {
		props[i] = types.ConstructTriggerProperty{
			Name: p.Name,
			Type: p.Type,
			Meta: types.ConstructTriggerPropertyMeta{
				Short: p.Description,
			},
		}
	}

	return types.ConstructTrigger{
		ResourceType: "automation:trigger-definition:" + ref,
		EventType:    "onAgentic",
		Meta: &types.ConstructTriggerMeta{
			Short:       def.Meta.Short,
			Description: def.Meta.Description,
			Icon:        def.Meta.Icon,
		},
		Properties: props,
	}
}

func validateAgenticInput(schema types.TriggerDefinitionSchema, input *expr.Vars) error {
	if input == nil {
		input = &expr.Vars{}
	}

	for _, p := range schema {
		val, err := input.Select(p.Name)
		if p.Required && (err != nil || val == nil) {
			return fmt.Errorf("missing required parameter %q", p.Name)
		}
		// @todo type validation if needed
	}

	return nil
}

func loadTriggerDefinition(ctx context.Context, s store.Storer, ID uint64) (res *types.TriggerDefinition, err error) {
	if ID == 0 {
		return nil, TriggerDefinitionErrInvalidID()
	}

	if res, err = store.LookupAutomationTriggerDefinitionByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, TriggerDefinitionErrNotFound()
	}

	return
}

// toLabeledTriggerDefinitions converts to []label.LabeledResource
func toLabeledTriggerDefinitions(set []*types.TriggerDefinition) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
