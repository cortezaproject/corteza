package scope

import "context"

// Services is the default service bag for Runtime[Services].
// Implementations live in their own packages and depend on this package — not the reverse.
//
// Tenant-level fields: RBAC, EventBus, Corredor, AgenticRuntime, LLM, DAL.
// Project-level fields: AutomationRunner, WorkflowRunner, Scheduler.
// Nil project-level fields on tenant Runtimes is intentional.
type Services struct {
	// Tenant-level
	RBAC          AccessControlService
	EventBus      EventBusService
	Corredor      CorredorService
	AgenticRunner AgenticRuntimeService
	LLM           LLMService
	DAL           DALService

	// Project-level only; nil on tenant Runtimes.
	AutomationRunner AutomationRunnerService
	WorkflowRunner   WorkflowRunnerService
	Scheduler        SchedulerService
}

// AccessControlService evaluates permissions and manages the user-group hierarchy.
// Covers pkg/rbac service + orgTree (UpdateUserGroups delegates to orgTree internally).
type AccessControlService interface {
	Can(ses interface{}, op string, res interface{}) bool
	Check(ses interface{}, op string, res interface{}) interface{}
	Grant(ctx context.Context, rules ...interface{}) error
	UpdateUserGroups(gm ...interface{}) error
	Watch(ctx context.Context)
}

// EventBusService dispatches and handles scoped events.
type EventBusService interface {
	Dispatch(ctx context.Context, ev interface{})
	WaitFor(ctx context.Context, ev interface{}) error
	Register(h interface{}, ops ...interface{}) uintptr
	Unregister(ptrs ...uintptr)
}

// CorredorService executes server-side scripts via the Corredor gRPC bridge.
type CorredorService interface {
	// Expand when corredor is mapped to runtime layers.
}

// AgenticRuntimeService executes agentic (LLM tool-loop) requests.
type AgenticRuntimeService interface {
	Run(ctx context.Context, req interface{}) (interface{}, error)
}

// LLMService sends prompts to configured LLM providers and returns responses.
type LLMService interface {
	Prompt(ctx context.Context, providerID uint64, model string, temperature *float64, outputTokens int, messages interface{}, tools interface{}) (interface{}, error)
	Chat(ctx context.Context, prompt string, history interface{}, tools interface{}, config interface{}) (interface{}, error)
	ListModels(ctx context.Context, providerID uint64) ([]string, error)
}

// DALService is the data-access layer: CRUD + pipeline execution over scoped connections.
type DALService interface {
	Create(ctx context.Context, mf interface{}, operations interface{}, rr ...interface{}) ([]map[string]any, error)
	Update(ctx context.Context, mf interface{}, operations interface{}, rr ...interface{}) error
	Search(ctx context.Context, mf interface{}, operations interface{}, f interface{}) (interface{}, error)
	Count(ctx context.Context, mf interface{}, operations interface{}, f interface{}) (uint, error)
	Run(ctx context.Context, pp interface{}) (interface{}, error)
}

// AutomationRunnerService manages ng_automation execution sessions within a project.
// Covers pkg/automation_exec runtimeManager + supervisor.
type AutomationRunnerService interface {
	Start(ctx context.Context, executableID interface{}, revision int, params interface{}) (interface{}, error)
	Stop(execID interface{}) error
}

// WorkflowRunnerService manages wfexec workflow execution sessions within a project.
type WorkflowRunnerService interface {
	// Expand when wfexec session lifecycle is mapped to runtime layers.
}

// SchedulerService submits timed jobs into the global worker pool on behalf of a project.
type SchedulerService interface {
	// Expand when scheduler is mapped to runtime layers.
}
