package service

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/auth"
	intAuth "github.com/crusttech/human/server/pkg/auth"
	execTypes "github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/pkg/scope"
	sysAutomation "github.com/crusttech/human/server/system/automation"

	"github.com/crusttech/human/server/store"
	"go.uber.org/zap"
)

type (
	ngAutomationServices struct {
		mux      *sync.RWMutex
		eventbus ngAutomationEventTriggerHandler
		reg      map[uint64]map[uint64]uintptr
		log      *zap.Logger
		parser   expr.Parsable
	}

	ngAutomationAccessController interface {
		CanCreateNgAutomation(context.Context) bool
		CanSearchNgAutomations(context.Context) bool
		CanReadNgAutomation(context.Context, *types.NgAutomation) bool
		CanUpdateNgAutomation(context.Context, *types.NgAutomation) bool
		CanDeleteNgAutomation(context.Context, *types.NgAutomation) bool
		CanUndeleteNgAutomation(context.Context, *types.NgAutomation) bool

		Grant(ctx context.Context, rr ...*rbac.Rule) error

		ngAutomationExecController
	}

	ngAutomationExecController interface {
		// CanExecuteNgAutomation(context.Context, *types.NgAutomation) bool
	}

	ngAutomationEventTriggerHandler interface {
		Register(h eventbus.HandlerFn, ops ...eventbus.HandlerRegOp) uintptr
		Unregister(ptrs ...uintptr)
	}

	executionEngine interface {
		DeprecateExecutable(ctx context.Context, exeID id.ID, rev int) error
		Execute(ctx context.Context, exeID id.ID, rev int, params execTypes.ExecutionParams) (id.ID, error)
		ExecuteAndWait(ctx context.Context, exeID id.ID, rev int, params execTypes.ExecutionParams) (*execTypes.ExecutionResult, error)
		RegisterExecutable(ctx context.Context, exe execTypes.Executable) error
		RemoveExecutable(ctx context.Context, exeID id.ID, rev int) error

		GetExecutionTrace(ctx context.Context, exeID id.ID, rev int, execID id.ID) ([]execTypes.StackFrame, error)
		ListExecutions(ctx context.Context, exeID id.ID, rev int) ([]*execTypes.ExecutionResult, error)
		ListAllExecutions(ctx context.Context, f execTypes.ExecutionFilter) ([]*execTypes.ExecutionResult, error)
	}

	ngAutomationUpdateHandler func(ctx context.Context, ns *types.NgAutomation) (ngAutomationChanges, error)
	ngAutomationChanges       uint8
)

const (
	ngAutomationUnchanged     ngAutomationChanges = 0
	ngAutomationChanged       ngAutomationChanges = 1
	ngAutomationLabelsChanged ngAutomationChanges = 2
	ngAutomationDefChanged    ngAutomationChanges = 4
)

func NgAutomation(log *zap.Logger, corredorOpt options.CorredorOpt) *ngAutomation {
	return &ngAutomation{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
		services: &ngAutomationServices{
			log: log,

			// registry for event bus triggers
			// @todo can we get rid of this and keep track of event bus entries via
			// some reference/identifier?
			mux:      &sync.RWMutex{},
			reg:      map[uint64]map[uint64]uintptr{},
			eventbus: eventbus.Service(),
			parser:   expr.NewParser(),
		},
	}
}

// engine resolves the execution engine for the scope on ctx. NG automation is
// project-scoped, so each tenant/project combination gets its own isolated
// engine, lazily built by the registry on first access.
func (svc *ngAutomation) engine(ctx context.Context) (executionEngine, error) {
	sc := scope.GetScopeFromContext(ctx)

	rt, ok := scope.Default.Get(sc)
	if !ok {
		return nil, fmt.Errorf("automation runtime missing for tenant %d project %d", sc.TenantID, sc.ProjectID)
	}

	eng, ok := scope.Get[executionEngine](rt.Container())
	if !ok {
		return nil, fmt.Errorf("automation engine missing for tenant %d project %d", sc.TenantID, sc.ProjectID)
	}

	return eng, nil
}

func (svc *ngAutomation) FindByID(ctx context.Context, ngAutomationID uint64) (ngAtuomation *types.NgAutomation, err error) {
	var (
		wap = &ngAutomationActionProps{ngAutomation: &types.NgAutomation{ID: ngAutomationID}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if ngAtuomation, err = loadNgAutomation(ctx, s, ngAutomationID); err != nil {
			return err
		}

		if !svc.ac.CanReadNgAutomation(ctx, ngAtuomation) {
			return NgAutomationErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, ngAtuomation); err != nil {
			return err
		}

		return nil
	})

	return ngAtuomation, svc.recordAction(ctx, wap, NgAutomationActionLookup, err)
}

func (svc *ngAutomation) LookupByHandle(ctx context.Context, handle string) (*types.NgAutomation, error) {
	rr, _, err := svc.Search(ctx, types.NgAutomationFilter{Handle: handle})
	if err != nil {
		return nil, err
	}
	if len(rr) == 0 {
		return nil, NgAutomationErrNotFound()
	}
	return rr[0], nil
}

// Create adds new ngAutomation resource and saves it into store
// It updates service's cache
func (svc *ngAutomation) guard(ctx context.Context, res *types.NgAutomation) error {
	return guardProjectWritable(ctx, svc.store, res.ProjectID)
}

func (svc *ngAutomation) Create(ctx context.Context, new *types.NgAutomation) (automation *types.NgAutomation, err error) {
	var (
		wap   = &ngAutomationActionProps{ngAutomation: new}
		cUser = intAuth.GetIdentityFromContext(ctx).Identity()
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !svc.ac.CanCreateNgAutomation(ctx) {
			return NgAutomationErrNotAllowedToCreate()
		}

		if new.Meta.Short == "" {
			return NgAutomationErrMissingName()
		}

		if !handle.IsValid(new.Handle) {
			return NgAutomationErrInvalidHandle()
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		// @note triggers have an ID to simplify referencing
		triggers := make(types.NgAutomationTriggerSet, len(new.Triggers))
		for i := range triggers {
			t := new.Triggers[i]
			t.ID = nextID()
			triggers[i] = t
		}

		automation = &types.NgAutomation{
			ID:        nextID(),
			ProjectID: new.ProjectID,
			Handle:    new.Handle,
			Labels:    new.Labels,
			Meta:      new.Meta,
			Enabled:   new.Enabled,

			Scope:    new.Scope,
			Triggers: triggers,
			Steps:    new.Steps,
			Paths:    new.Paths,

			// @todo need to check against access control if current user can modify security descriptor
			RunAs:     new.RunAs,
			OwnedBy:   cUser,
			CreatedAt: *now(),
			CreatedBy: cUser,
		}

		wap.ngAutomation = automation

		res, exec, err := svc.procAutomation(ctx, automation)
		if err != nil {
			return
		}

		if len(res.Issues) == 0 {
			eng, err := svc.engine(ctx)
			if err != nil {
				return err
			}

			err = eng.RegisterExecutable(ctx, exec)
			if err != nil {
				return err
			}

			err = svc.registerAutomation(ctx, automation)
			if err != nil {
				return err
			}
		}

		if err = store.CreateAutomationNgAutomation(ctx, s, automation); err != nil {
			return
		}

		if err = label.Create(ctx, s, automation); err != nil {
			return
		}

		wap.setNew(automation)

		return
	})

	return automation, svc.recordAction(ctx, wap, NgAutomationActionCreate, err)
}

// Update modifies existing ngAutomation resource in the store
func (svc *ngAutomation) Update(ctx context.Context, upd *types.NgAutomation) (*types.NgAutomation, error) {
	return svc.updater(ctx, upd.ID, NgAutomationActionUpdate, func(ctx context.Context, res *types.NgAutomation) (ngAutomationChanges, error) {
		if upd.Meta == nil || upd.Meta.Short == "" {
			return ngAutomationUnchanged, NgAutomationErrMissingName()
		}

		if !svc.ac.CanUpdateNgAutomation(ctx, res) {
			return ngAutomationUnchanged, NgAutomationErrNotAllowedToUpdate()
		}

		handler := svc.handleUpdate(upd)
		return handler(ctx, res)
	})
}

func (svc *ngAutomation) DeleteByID(ctx context.Context, ngAutomationID uint64) error {
	return trim1st(svc.updater(ctx, ngAutomationID, NgAutomationActionDelete, func(ctx context.Context, res *types.NgAutomation) (ngAutomationChanges, error) {
		changes, err := svc.handleDelete(ctx, res)
		if err != nil {
			return ngAutomationUnchanged, err
		}

		return changes, err
	}))
}

func (svc *ngAutomation) UndeleteByID(ctx context.Context, ngAutomationID uint64) error {
	return trim1st(svc.updater(ctx, ngAutomationID, NgAutomationActionUndelete, func(ctx context.Context, res *types.NgAutomation) (ngAutomationChanges, error) {
		var (
			changes, err = svc.handleUndelete(ctx, res)
		)

		if err != nil {
			return ngAutomationUnchanged, err
		}

		return changes, err

	}))
}

func (svc *ngAutomation) onExec(ctx context.Context, _ *ngAutomationActionProps, automationID uint64, p types.NgAutomationExecParams) (executionID id.ID, err error) {
	atm, loadErr := loadNgAutomation(ctx, svc.store, automationID)
	if loadErr != nil {
		return id.Zero(), loadErr
	}

	p.Input, err = svc.injectIdentities(ctx, atm.RunAs, p.Input)
	if err != nil {
		return
	}

	eng, err := svc.engine(ctx)
	if err != nil {
		return
	}

	executionID, err = eng.Execute(ctx, id.MustNumID(automationID), 0, execTypes.ExecutionParams{
		EntryPoint:   p.EntryPoint,
		Input:        p.Input,
		EventType:    p.EventType,
		ResourceType: p.ResourceType,
	})
	if err != nil {
		return
	}

	return
}

func (svc *ngAutomation) onExecAndWait(ctx context.Context, _ *ngAutomationActionProps, automationID uint64, p types.NgAutomationExecParams) (out *execTypes.ExecutionResult, err error) {
	entryPoint := p.EntryPoint
	atm, loadErr := loadNgAutomation(ctx, svc.store, automationID)
	if loadErr != nil {
		return nil, loadErr
	}

	if entryPoint == "" {
		// In case no entry point, use the first trigger as the parent.
		if len(atm.Triggers) > 0 {
			entryPoint = atm.Triggers[0].Handle
		}
	}

	// Validate input against embedded schema if the trigger corresponds to the entrypoint
	for _, t := range atm.Triggers {
		if t.Handle == entryPoint || (entryPoint == "" && t.ID == atm.Triggers[0].ID) {
			if err := validateAgenticInput(t.InputSchema, p.Input); err != nil {
				return nil, err
			}
			break
		}
	}

	p.Input, err = svc.injectIdentities(ctx, atm.RunAs, p.Input)
	if err != nil {
		return
	}

	eng, err := svc.engine(ctx)
	if err != nil {
		return
	}

	out, err = eng.ExecuteAndWait(ctx, id.MustNumID(automationID), 0, execTypes.ExecutionParams{
		EntryPoint:   entryPoint,
		Input:        p.Input,
		EventType:    p.EventType,
		ResourceType: p.ResourceType,
	})
	if err != nil {
		return
	}

	return
}

// injectIdentities resolves the invoker (from ctx) and runner (from runAsID, defaulting
// to invoker) to full *system/types.User values and injects them as typed expr vars
// into the execution input scope.
func (svc *ngAutomation) injectIdentities(ctx context.Context, runAsID uint64, input *expr.Vars) (*expr.Vars, error) {
	sysCtx := intAuth.SetIdentityToContext(ctx, intAuth.ServiceUser())

	if input == nil {
		input = &expr.Vars{}
	}

	invokerIdentity := intAuth.GetIdentityFromContext(ctx)
	invokerUser, err := DefaultUser.FindByAny(sysCtx, invokerIdentity.Identity())
	if err != nil {
		return nil, fmt.Errorf("failed to resolve invoker user: %w", err)
	}

	invokerExpr, err := sysAutomation.NewUser(invokerUser)
	if err != nil {
		return nil, err
	}
	_ = input.AssignFieldValue("invoker", invokerExpr)

	runnerExpr := invokerExpr
	if runAsID > 0 {
		runnerUser, err := DefaultUser.FindByAny(sysCtx, runAsID)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve runner user: %w", err)
		}
		runnerExpr, err = sysAutomation.NewUser(runnerUser)
		if err != nil {
			return nil, err
		}
	}
	_ = input.AssignFieldValue("runner", runnerExpr)

	return input, nil
}

func (svc *ngAutomation) onGetExecutions(ctx context.Context, _ *ngAutomationActionProps, automationID uint64) (out []*execTypes.ExecutionResult, err error) {
	eng, err := svc.engine(ctx)
	if err != nil {
		return nil, err
	}

	out, err = eng.ListExecutions(ctx, id.MustNumID(automationID), 0)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (svc *ngAutomation) onGetExecutionTrace(ctx context.Context, _ *ngAutomationActionProps, exeID, executionID uint64, rev int) (out []execTypes.StackFrame, err error) {
	eng, err := svc.engine(ctx)
	if err != nil {
		return
	}

	out, err = eng.GetExecutionTrace(ctx, id.MustNumID(exeID), 0, id.MustNumID(executionID))
	if err != nil {
		return
	}

	return
}

func (svc *ngAutomation) onGetAllExecutions(ctx context.Context, _ *ngAutomationActionProps, f execTypes.ExecutionFilter) (out []*execTypes.ExecutionResult, err error) {
	eng, err := svc.engine(ctx)
	if err != nil {
		return nil, err
	}

	return eng.ListAllExecutions(ctx, f)
}

func (svc ngAutomation) uniqueCheck(ctx context.Context, res *types.NgAutomation) (err error) {
	// Scoped to the project, matching the store lookup. A handle is unique
	// WITHIN a project, not globally: a project revision branch copies its
	// automations into the draft, so the same handle deliberately exists in both
	// revisions at once (see system/service/project_revision_clone.go).
	if res.Handle != "" {
		if e, _ := store.LookupAutomationNgAutomationByProjectIDHandle(ctx, svc.store, res.ProjectID, res.Handle); e != nil && e.ID != res.ID {
			return NgAutomationErrHandleNotUnique()
		}
	}

	return nil
}

func (svc *ngAutomation) updater(ctx context.Context, ngAutomationID uint64, action func(...*ngAutomationActionProps) *ngAutomationAction, fn ngAutomationUpdateHandler) (*types.NgAutomation, error) {
	var (
		changes ngAutomationChanges
		old     *types.NgAutomation
		res     *types.NgAutomation
		aProps  = &ngAutomationActionProps{ngAutomation: &types.NgAutomation{ID: ngAutomationID}}
		err     error
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {

		res, err = loadNgAutomation(ctx, s, ngAutomationID)
		if err != nil {
			return
		}

		old = res.Clone()

		res.Issues = nil

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setNgAutomation(res)
		aProps.setUpdate(res)

		if changes, err = fn(ctx, res); err != nil {
			return err
		}

		// Triggers
		triggers := make(types.NgAutomationTriggerSet, len(res.Triggers))
		for i := range triggers {
			t := res.Triggers[i]

			if t.ID == 0 {
				t.ID = nextID()
			}

			triggers[i] = t
		}

		res.Triggers = triggers

		res, exec, err := svc.procAutomation(ctx, res)
		if err != nil {
			return
		}

		if len(res.Issues) == 0 {
			eng, err := svc.engine(ctx)
			if err != nil {
				return err
			}

			err = eng.RegisterExecutable(ctx, exec)
			if err != nil {
				return err
			}
		}

		err = svc.registerAutomation(ctx, res)
		if err != nil {
			return
		}

		if changes&ngAutomationChanged > 0 || len(res.Issues) > 0 {
			if err = store.UpdateAutomationNgAutomation(ctx, svc.store, res); err != nil {
				return err
			}
		}

		if changes&ngAutomationLabelsChanged > 0 {
			if err = label.Update(ctx, s, res); err != nil {
				return
			}
		}

		return
	})

	return res, svc.recordAction(ctx, aProps, action, err, old, res)
}

func (svc ngAutomation) handleUpdate(upd *types.NgAutomation) ngAutomationUpdateHandler {
	return func(ctx context.Context, res *types.NgAutomation) (changes ngAutomationChanges, err error) {
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ngAutomationUnchanged, NgAutomationErrStaleData()
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return ngAutomationUnchanged, NgAutomationErrInvalidHandle()
		}

		if err := svc.uniqueCheck(ctx, upd); err != nil {
			return ngAutomationUnchanged, err
		}

		if !svc.ac.CanUpdateNgAutomation(ctx, res) {
			return ngAutomationUnchanged, NgAutomationErrNotAllowedToUpdate()
		}

		if res.Handle != upd.Handle {
			changes |= ngAutomationChanged
			res.Handle = upd.Handle
		}

		if res.Enabled != upd.Enabled {
			changes |= ngAutomationChanged | ngAutomationDefChanged
			res.Enabled = upd.Enabled
		}

		if upd.Labels != nil {
			if label.Changed(res.Labels, upd.Labels) {
				changes |= ngAutomationLabelsChanged
				res.Labels = upd.Labels
			}
		}

		if upd.Meta != nil {
			if !reflect.DeepEqual(upd.Meta, res.Meta) {
				changes |= ngAutomationChanged
				res.Meta = upd.Meta
			}
		}

		if upd.Scope != nil {
			if !reflect.DeepEqual(upd.Scope, res.Scope) {
				changes |= ngAutomationChanged | ngAutomationDefChanged
				res.Scope = upd.Scope
			}
		}

		if upd.Triggers != nil {
			if !reflect.DeepEqual(upd.Triggers, res.Triggers) {
				changes |= ngAutomationChanged | ngAutomationDefChanged
				res.Triggers = upd.Triggers
			}
		}

		if upd.Steps != nil {
			if !reflect.DeepEqual(upd.Steps, res.Steps) {
				changes |= ngAutomationChanged | ngAutomationDefChanged
				res.Steps = upd.Steps
			}
		}

		if upd.Paths != nil {
			if !reflect.DeepEqual(upd.Paths, res.Paths) {
				changes |= ngAutomationChanged | ngAutomationDefChanged
				res.Paths = upd.Paths
			}
		}

		if res.RunAs != upd.RunAs {
			// @todo need to check against access control if current user can modify security descriptor
			changes |= ngAutomationChanged | ngAutomationDefChanged
			res.RunAs = upd.RunAs
		}

		if res.OwnedBy != upd.OwnedBy {
			// @todo need to check against access control if current user can modify owner
			changes |= ngAutomationChanged
			res.OwnedBy = upd.OwnedBy
		}

		if changes&ngAutomationChanged > 0 {
			res.UpdatedAt = now()
		}

		return
	}
}

func (svc ngAutomation) handleDelete(ctx context.Context, res *types.NgAutomation) (ngAutomationChanges, error) {
	if !svc.ac.CanDeleteNgAutomation(ctx, res) {
		return ngAutomationUnchanged, NgAutomationErrNotAllowedToDelete()
	}

	if res.DeletedAt != nil {
		// ngAutomation already deleted
		return ngAutomationUnchanged, nil
	}

	res.DeletedAt = now()
	return ngAutomationChanged, nil
}

func (svc ngAutomation) handleUndelete(ctx context.Context, res *types.NgAutomation) (ngAutomationChanges, error) {
	if !svc.ac.CanUndeleteNgAutomation(ctx, res) {
		return ngAutomationUnchanged, NgAutomationErrNotAllowedToUndelete()
	}

	if res.DeletedAt == nil {
		// ngAutomation not deleted
		return ngAutomationUnchanged, nil
	}

	res.DeletedAt = nil
	return ngAutomationChanged, nil
}

func (svc *ngAutomation) onLoad(ctx context.Context, _ *ngAutomationActionProps) error {
	var (
		set, _, err = store.SearchAutomationNgAutomations(ctx, svc.store, types.NgAutomationFilter{
			Deleted:  filter.StateExcluded,
			Disabled: filter.StateExcluded,
		})
	)
	if err != nil {
		return err
	}

	// @todo Load currently registers every automation into the engine for the
	//       scope on ctx (system scope at boot). Once automations are loaded
	//       per project, resolve the engine per automation's scope instead.
	eng, err := svc.engine(ctx)
	if err != nil {
		return err
	}

	for _, atn := range set {
		atm, exe, err := svc.procAutomation(ctx, atn)
		if err != nil {
			// @todo?
			return err
		}

		if len(atm.Issues) == 0 {
			err = eng.RegisterExecutable(ctx, exe)
			if err != nil {
				// @todo?
				return err
			}

			err = svc.registerAutomation(ctx, atm)
			if err != nil {
				// @todo?
				return err
			}
		}
	}

	return nil
}

func loadNgAutomation(ctx context.Context, s store.Storer, ngAutomationID uint64) (res *types.NgAutomation, err error) {
	if ngAutomationID == 0 {
		return nil, NgAutomationErrInvalidID()
	}

	// @todo :)
	if res, err = store.LookupAutomationNgAutomationByID(ctx, s, ngAutomationID); errors.IsNotFound(err) {
		return nil, NgAutomationErrNotFound()
	}

	return
}

func (svc *ngAutomation) procAutomation(ctx context.Context, atm *types.NgAutomation) (out *types.NgAutomation, exe execTypes.Executable, err error) {
	out = atm

	exe, issues := ConvertNgAutomation(ctx, svc, out)
	if len(issues) > 0 {
		out.Issues = types.NgAutomationIssueSet(issues)
	}

	// Returns context with identity set to service user
	//
	// Current user (identity in the context) might not have
	// sufficient privileges to load info about invoker and runner
	sysUserCtx := func() context.Context {
		return intAuth.SetIdentityToContext(ctx, intAuth.ServiceUser())
	}

	// @todo this might not be the smartest thing, users might get invalidated after
	//       we add cache them as workflow runners
	if out.RunAs > 0 {
		if exe.RunAs, err = DefaultUser.FindByAny(sysUserCtx(), out.RunAs); err != nil {
			out.Issues = append(out.Issues, &types.NgAutomationIssue{
				Code:     types.IssueCodeRunAsLoadFailed,
				Severity: types.NgAutomationSeverityError,
				Message:  fmt.Sprintf("failed to load run-as user %d: %v", out.RunAs, err),
			})
		} else if !exe.RunAs.Valid() {
			out.Issues = append(out.Issues, &types.NgAutomationIssue{
				Code:     types.IssueCodeRunAsInvalid,
				Severity: types.NgAutomationSeverityError,
				Message:  fmt.Sprintf("invalid user %d used for workflow run-as", out.RunAs),
			})
		}
	}

	return
}

func (svc *ngAutomation) unregisterAutomation(a *types.NgAutomation) {
	svc.services.mux.Lock()
	defer svc.services.mux.Unlock()

	if ptrs, ok := svc.services.reg[a.ID]; ok {
		for _, ptr := range ptrs {
			svc.services.eventbus.Unregister(ptr)
		}
		delete(svc.services.reg, a.ID)
	}
}

func (svc *ngAutomation) registerAutomation(ctx context.Context, a *types.NgAutomation) error {
	if !a.Enabled || a.DeletedAt != nil || len(a.Issues) > 0 || len(a.Triggers) == 0 {
		svc.unregisterAutomation(a)
		return nil
	}

	var runAs auth.Identifiable
	if a.RunAs > 0 {
		sysCtx := auth.SetIdentityToContext(ctx, auth.ServiceUser())
		var err error
		if runAs, err = DefaultUser.FindByAny(sysCtx, a.RunAs); err != nil {
			return fmt.Errorf("failed to load run-as user %d: %w", a.RunAs, err)
		}
		if !runAs.Valid() {
			return fmt.Errorf("invalid user %d used for automation run-as", a.RunAs)
		}
	}

	log := svc.services.log.With(logger.Uint64("automationID", a.ID))

	svc.services.mux.Lock()
	defer svc.services.mux.Unlock()

	if svc.services.reg[a.ID] == nil {
		svc.services.reg[a.ID] = make(map[uint64]uintptr)
	}

	for _, t := range a.Triggers {
		if !t.Enabled {
			continue
		}

		svc.registerTrigger(log, a, t)
	}

	return nil
}

func (svc *ngAutomation) registerTrigger(log *zap.Logger, a *types.NgAutomation, t *types.NgAutomationTrigger) {
	log = log.With(logger.Uint64("triggerID", t.ID))

	if t.ResourceType == "automation:trigger:agentic" {
		// EventBus is skipped for natively-executed agentic triggers
		if ptr := svc.services.reg[a.ID][t.ID]; ptr != 0 {
			svc.services.eventbus.Unregister(ptr)
			delete(svc.services.reg[a.ID], t.ID)
		}
		return
	}

	// Always unregister existing handler
	if ptr := svc.services.reg[a.ID][t.ID]; ptr != 0 {
		svc.services.eventbus.Unregister(ptr)
	}

	ops := []eventbus.HandlerRegOp{
		eventbus.On(t.EventType),
		eventbus.For(t.ResourceType),
	}

	for _, c := range t.Constraints {
		name, op, values, err := svc.prepConstraintBits(c)
		if err != nil {
			log.Debug("failed to prepare constraint for automation trigger",
				zap.Any("constraint", c),
				zap.Error(err),
			)
			continue
		}

		cnstr, err := eventbus.ConstraintMaker(name, op, values...)
		if err != nil {
			log.Debug("failed to make constraint for automation trigger",
				zap.Any("constraint", c),
				zap.Error(err),
			)
			continue
		}
		ops = append(ops, eventbus.Constraint(cnstr))
	}

	handlerFn := makeAutomationHandler(svc, a, t)
	svc.services.reg[a.ID][t.ID] = svc.services.eventbus.Register(handlerFn, ops...)

	log.Debug("trigger registered",
		zap.String("eventType", t.EventType),
		zap.String("resourceType", t.ResourceType),
		zap.Any("constraints", t.Constraints),
	)
}

// @todo expand when needed; will probably need some custom handler on the trigger
func (svc *ngAutomation) prepConstraintBits(c types.NgTriggerConstraint) (name string, _ string, values []string, err error) {
	if len(c.Values) == 0 {
		return c.Name, c.Op, nil, nil
	}

	typ := c.Values[0].Type
	for _, v := range c.Values {
		if v.Type != typ {
			err = fmt.Errorf("constraint values must be of the same type")
			return
		}

		values = append(values, v.Value)
	}

	switch typ {
	case "String":
		name = fmt.Sprintf("%s.%s", c.Name, "name")
	case "Handle":
		name = fmt.Sprintf("%s.%s", c.Name, "handle")
	case "ID":
		name = fmt.Sprintf("%s.%s", c.Name, "id")
	default:
		err = fmt.Errorf("unknown constraint type %s", typ)
		return
	}

	return name, c.Op, values, nil
}

func makeAutomationHandler(svc *ngAutomation, a *types.NgAutomation, t *types.NgAutomationTrigger) eventbus.HandlerFn {
	return func(ctx context.Context, ev eventbus.Event) error {
		var scope *expr.Vars
		if dec, ok := ev.(varsEncoder); ok {
			var err error
			if scope, err = dec.EncodeVars(); err != nil {
				return err
			}
		}

		_, err := svc.ExecAndWait(ctx, a.ID, types.NgAutomationExecParams{
			EntryPoint:   t.Handle,
			EventType:    t.EventType,
			ResourceType: t.ResourceType,
			Input:        t.Input.MustMerge(scope),
		})
		return err
	}
}

func validateAgenticInput(schema types.NgAutomationTriggerSchema, input *expr.Vars) error {
	if len(schema) == 0 {
		return nil
	}
	if input == nil {
		input = &expr.Vars{}
	}
	for _, p := range schema {
		if p.Required && !input.Has(p.Name) {
			return fmt.Errorf("missing required input: %s", p.Name)
		}
	}
	return nil
}
