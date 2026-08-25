package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/label"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/agentic/tcl"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD skeleton (FindByID, Search, Create, Update, DeleteByID and the
// toLabeledAgents helper) is generated in agent.gen.go from system/agent.cue.
//
// This file owns the struct, access-controller interface, constructor, the
// LLM-validator opt-in, the on<Op> custom bodies the generated FindByID /
// Search / Create / Update delegate to, the hand-written UndeleteByID (the
// generated undelete is disabled because it uses the non-standard
// CanDeleteAgent + AgentErrNotAllowedToDelete pairing), and the
// resource-specific helpers (Get, prepareTCL). The loadAgent helper is
// generated in agent.gen.go.

type (
	agentLLMValidator interface {
		ValidateTemperature(ctx context.Context, providerID uint64, model string, temperature *float64) error
	}

	agentServices struct {
		llm agentLLMValidator
		// namespaceID resolves a compose namespace reference. Injected rather
		// than imported: compose depends on system, so this package cannot
		// reach the other way.
		namespaceID func(ctx context.Context, ref string) (uint64, error)
	}
)

func Agent() *agent {
	return &agent{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		services:  &agentServices{},
	}
}

// WithNamespaceResolver supplies the compose namespace lookup used to scope a
// new agent's default grant.
func (svc *agent) WithNamespaceResolver(fn func(ctx context.Context, ref string) (uint64, error)) *agent {
	svc.services.namespaceID = fn
	return svc
}

func (svc *agent) WithLLMValidator(v agentLLMValidator) *agent {
	svc.services.llm = v
	return svc
}

func (svc *agent) Get(ctx context.Context, ID uint64) (*types.Agent, error) {
	return svc.FindByID(ctx, ID)
}

// onLookup is the custom body for the generated FindByID. The generated method
// owns the action-log scaffold + recordAction; the load, read access check and
// label load (which the standard template does not emit) live here.
func (svc *agent) guard(ctx context.Context, res *types.Agent) error {
	return guardProjectWritable(ctx, svc.store, res.ProjectID)
}

func (svc *agent) onLookup(ctx context.Context, ID uint64, aProps *agentActionProps) (a *types.Agent, err error) {
	if a, err = loadAgent(ctx, svc.store, ID); err != nil {
		return nil, err
	}

	aProps.setAgent(a)

	if !svc.ac.CanReadAgent(ctx, a) {
		return nil, AgentErrNotAllowedToRead()
	}

	if err = label.Load(ctx, svc.store, a); err != nil {
		return nil, err
	}

	return a, nil
}

// onCreate is the custom body for the generated Create. The generated method
// owns the action-log scaffold + recordAction + the CanCreateAgent check; the
// Revision / Status defaults, optional temperature validation, prepareTCL and
// persistence live here.
func (svc *agent) onCreate(ctx context.Context, new *types.Agent) (err error) {
	new.ID = nextID()
	new.CreatedAt = *now()
	new.Revision = 1

	if new.Status == "" {
		new.Status = "active"
	}

	if new.Execution.Model.Temperature != nil && svc.services.llm != nil {
		if err = svc.services.llm.ValidateTemperature(ctx, new.Execution.Model.LLMProviderID, new.Execution.Model.Model, new.Execution.Model.Temperature); err != nil {
			return
		}
	}

	prepareTCL(&new.Behavior)

	if err = svc.defaultAccess(ctx, new); err != nil {
		return
	}

	if err = store.CreateAgent(ctx, svc.store, new); err != nil {
		return
	}

	if err = label.Create(ctx, svc.store, new); err != nil {
		return
	}

	return nil
}

// defaultAccess gives a new agent something to do.
//
// Access is deny-by-default, so an agent created without a single grant can
// call nothing: it answers every question with a refusal, which reads as broken
// rather than as unconfigured. When the author has said which namespace the
// agent is for and granted nothing, the useful reading of that is "let it read
// this namespace" — the one posture that is immediately useful and cannot
// damage anything. Anything more than reading stays a deliberate act.
//
// Nothing is assumed when no namespace is named: there would be nothing to
// scope the grant to, and a grant with no scope is denied anyway.
func (svc *agent) defaultAccess(ctx context.Context, a *types.Agent) error {
	if len(a.Access.Tools) > 0 || len(a.Access.TAQs) > 0 || len(a.Access.Workflows) > 0 {
		return nil
	}

	if a.Access.Context.Namespace == "" {
		return nil
	}

	if svc.services.namespaceID == nil {
		return nil
	}

	ns, err := svc.services.namespaceID(ctx, a.Access.Context.Namespace)
	if err != nil || ns == 0 {
		// An unresolvable namespace is the caller's to fix; failing the create
		// over a convenience default would be worse than leaving it ungranted.
		return nil
	}

	a.Access.Tools = []types.AgentAccessTool{{
		Group:       "usage",
		MaxRisk:     "read",
		Description: "Read and aggregate data in this namespace",
		Allow:       []types.AgentAccessAllow{{NamespaceID: ns}},
	}}

	// The grant is only worth anything if the agent is told what it reaches.
	// Without the platform context there is no list of namespaces and modules,
	// so it guesses handles and reports that they do not exist. Set together or
	// not at all — this runs only for an agent that arrived with no access.
	a.Behavior.InjectSystemContext = true

	return nil
}

// onUpdate is the custom body for the generated Update. The generated method
// owns the action-log scaffold + recordAction; the update access check (on the
// incoming resource), the stale guard, the Revision bump, optional temperature
// validation, prepareTCL and the whole-record persistence live here.
func (svc *agent) onUpdate(ctx context.Context, s store.Storer, upd, res *types.Agent, aProps *agentActionProps, _ func() error, _ func() error) error {
	if !svc.ac.CanUpdateAgent(ctx, upd) {
		return AgentErrNotAllowedToUpdate()
	}

	upd.Revision = res.Revision + 1
	upd.UpdatedAt = now()
	upd.CreatedAt = res.CreatedAt
	upd.DeletedAt = res.DeletedAt
	// ProjectID is set at create and isn't part of the update payload; preserve
	// the existing owning project so an update never unlinks the agent from it
	// (otherwise it would vanish from the project-scoped resource graph).
	upd.ProjectID = res.ProjectID

	if upd.Execution.Model.Temperature != nil && svc.services.llm != nil {
		if err := svc.services.llm.ValidateTemperature(ctx, upd.Execution.Model.LLMProviderID, upd.Execution.Model.Model, upd.Execution.Model.Temperature); err != nil {
			return err
		}
	}

	prepareTCL(&upd.Behavior)

	if err := store.UpdateAgent(ctx, s, upd); err != nil {
		return err
	}

	// copy fields back so the caller (generated Update) returns the updated record
	*res = *upd

	return label.Update(ctx, s, upd)
}

func (svc *agent) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	err = func() (err error) {
		var a *types.Agent
		if a, err = loadAgent(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteAgent(ctx, a) {
			return AgentErrNotAllowedToDelete()
		}

		a.DeletedAt = nil
		if err = store.UpdateAgent(ctx, svc.store, a); err != nil {
			return
		}

		return nil
	}()

	return err
}

// onSearch is the custom body for the generated Search. The generated method
// owns the action-log scaffold + recordAction + the CanSearchAgents check; the
// Check predicate, label filtering and label load live here.
func (svc *agent) onSearch(ctx context.Context, filter types.AgentFilter, aProps *agentActionProps) (set types.AgentSet, f types.AgentFilter, err error) {
	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Agent) (bool, error) {
		if !svc.ac.CanReadAgent(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	if len(filter.Labels) > 0 {
		filter.LabeledIDs, err = label.Search(
			ctx,
			svc.store,
			types.Agent{}.LabelResourceKind(),
			filter.Labels,
		)

		if err != nil {
			return set, f, err
		}

		// labels specified but no labeled resources found
		if len(filter.LabeledIDs) == 0 {
			return set, f, nil
		}
	}

	if set, f, err = store.SearchAgents(ctx, svc.store, filter); err != nil {
		return set, f, err
	}

	if err = label.Load(ctx, svc.store, toLabeledAgents(set)...); err != nil {
		return set, f, err
	}

	return set, f, nil
}

// prepareTCL ensures TCL is properly initialized on the agent behavior:
// - if enabled and no articles selected, populate defaults
// - always merge hardwired articles back in (user cannot remove them)
// TCL is opt-in: only an explicit true enables it; nil (unset) means off.
func prepareTCL(b *types.AgentBehavior) {
	if b.TreatyCLEnabled == nil || !*b.TreatyCLEnabled {
		return
	}
	if len(b.TreatyCLArticles) == 0 {
		b.TreatyCLArticles = tcl.DefaultArticleIDs()
	} else {
		b.TreatyCLArticles = tcl.MergeWithHardwired(b.TreatyCLArticles)
	}
}
